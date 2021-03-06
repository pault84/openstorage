package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/opsworks"
	sh "github.com/codeskyblue/go-sh"
	oexec "github.com/libopenstorage/openstorage/pkg/exec"
	"github.com/libopenstorage/openstorage/pkg/storageops"
	"github.com/portworx/sched-ops/task"
	"github.com/sirupsen/logrus"
)

const (
	awsDevicePrefix      = "/dev/sd"
	awsDevicePrefixWithX = "/dev/xvd"
	awsDevicePrefixWithH = "/dev/hd"
	awsDevicePrefixNvme  = "/dev/nvme"
)

type ec2Ops struct {
	instanceType string
	instance     string
	ec2          *ec2.EC2
	mutex        sync.Mutex
}

var (
	// ErrAWSEnvNotAvailable is the error type when aws credentials are not set
	ErrAWSEnvNotAvailable = fmt.Errorf("AWS credentials are not set in environment")
	nvmeCmd               = oexec.Which("nvme")
)

// NewEnvClient creates a new AWS storage ops instance using environment vars
func NewEnvClient() (storageops.Ops, error) {
	region, err := storageops.GetEnvValueStrict("AWS_REGION")
	if err != nil {
		return nil, err
	}

	instance, err := storageops.GetEnvValueStrict("AWS_INSTANCE_NAME")
	if err != nil {
		return nil, err
	}

	instanceType, err := storageops.GetEnvValueStrict("AWS_INSTANCE_TYPE")
	if err != nil {
		return nil, err
	}

	if _, err := credentials.NewEnvCredentials().Get(); err != nil {
		return nil, ErrAWSEnvNotAvailable
	}

	ec2 := ec2.New(
		session.New(
			&aws.Config{
				Region:      &region,
				Credentials: credentials.NewEnvCredentials(),
			},
		),
	)

	return NewEc2Storage(instance, instanceType, ec2), nil
}

// NewEc2Storage creates a new aws storage ops instance
func NewEc2Storage(instance, instanceType string, ec2 *ec2.EC2) storageops.Ops {
	return &ec2Ops{
		instance:     instance,
		instanceType: instanceType,
		ec2:          ec2,
	}
}

// nvmeInstanceTypes are list of instance types whose EBS volumes are exposed as NVMe block devices
var nvmeInstanceTypes = []string{"c5", "c5d", "i3.metal", "m5", "m5d", "r5", "r5d", "z1d"}

func (s *ec2Ops) filters(
	labels map[string]string,
	keys []string,
) []*ec2.Filter {
	if len(labels) == 0 {
		return nil
	}
	f := make([]*ec2.Filter, len(labels)+len(keys))
	i := 0
	for k, v := range labels {
		s := string("tag:") + k
		value := v
		f[i] = &ec2.Filter{Name: &s, Values: []*string{&value}}
		i++
	}
	for _, k := range keys {
		s := string("tag-key:") + k
		f[i] = &ec2.Filter{Name: &s}
		i++
	}
	return f
}

func (s *ec2Ops) tags(labels map[string]string) []*ec2.Tag {
	if len(labels) == 0 {
		return nil
	}
	t := make([]*ec2.Tag, len(labels))
	i := 0
	for k, v := range labels {
		key := k
		value := v
		t[i] = &ec2.Tag{Key: &key, Value: &value}
		i++
	}
	return t
}

func (s *ec2Ops) waitStatus(id string, desired string) error {
	request := &ec2.DescribeVolumesInput{VolumeIds: []*string{&id}}
	actual := ""

	_, err := task.DoRetryWithTimeout(
		func() (interface{}, bool, error) {
			awsVols, err := s.ec2.DescribeVolumes(request)
			if err != nil {
				return nil, true, err
			}

			if len(awsVols.Volumes) != 1 {
				return nil, true, fmt.Errorf("expected one volume %v got %v",
					id, len(awsVols.Volumes))
			}

			if awsVols.Volumes[0].State == nil {
				return nil, true, fmt.Errorf("Nil volume state for %v", id)
			}

			actual = *awsVols.Volumes[0].State
			if actual == desired {
				return nil, false, nil
			}

			return nil, true, fmt.Errorf(
				"Volume %v did not transition to %v current state %v",
				id, desired, actual)

		},
		storageops.ProviderOpsTimeout,
		storageops.ProviderOpsRetryInterval)

	return err

}

func (s *ec2Ops) waitAttachmentStatus(
	volumeID string,
	desired string,
	timeout time.Duration,
) (*ec2.Volume, error) {
	id := volumeID
	request := &ec2.DescribeVolumesInput{VolumeIds: []*string{&id}}
	interval := 2 * time.Second
	logrus.Infof("Waiting for state transition to %q", desired)

	f := func() (interface{}, bool, error) {
		awsVols, err := s.ec2.DescribeVolumes(request)
		if err != nil {
			return nil, false, err
		}
		if len(awsVols.Volumes) != 1 {
			return nil, false, fmt.Errorf("expected one volume %v got %v",
				volumeID, len(awsVols.Volumes))
		}

		var actual string
		vol := awsVols.Volumes[0]
		awsAttachment := vol.Attachments
		if awsAttachment == nil || len(awsAttachment) == 0 {
			// We have encountered scenarios where AWS returns a nil attachment state
			// for a volume transitioning from detaching -> attaching.
			actual = ec2.VolumeAttachmentStateDetached
		} else {
			actual = *awsAttachment[0].State
		}
		if actual == desired {
			return vol, false, nil
		}
		return nil, true, fmt.Errorf("Volume %v failed to transition to  %v current state %v",
			volumeID, desired, actual)
	}

	outVol, err := task.DoRetryWithTimeout(f, timeout, interval)
	if err != nil {
		return nil, err
	}
	if vol, ok := outVol.(*ec2.Volume); ok {
		return vol, nil
	}
	return nil, storageops.NewStorageError(storageops.ErrVolInval,
		fmt.Sprintf("Invalid volume object for volume %s", volumeID), "")
}

func (s *ec2Ops) Name() string { return "aws" }

func (s *ec2Ops) InstanceID() string { return s.instance }

func (s *ec2Ops) ApplyTags(volumeID string, labels map[string]string) error {
	req := &ec2.CreateTagsInput{
		Resources: []*string{&volumeID},
		Tags:      s.tags(labels),
	}
	_, err := s.ec2.CreateTags(req)
	return err
}

func (s *ec2Ops) RemoveTags(volumeID string, labels map[string]string) error {
	req := &ec2.DeleteTagsInput{
		Resources: []*string{&volumeID},
		Tags:      s.tags(labels),
	}
	_, err := s.ec2.DeleteTags(req)
	return err
}

func (s *ec2Ops) matchTag(tag *ec2.Tag, match string) bool {
	return tag.Key != nil &&
		tag.Value != nil &&
		len(*tag.Key) != 0 &&
		len(*tag.Value) != 0 &&
		*tag.Key == match
}

func (s *ec2Ops) DeviceMappings() (map[string]string, error) {
	instance, err := s.describe()
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, d := range instance.BlockDeviceMappings {
		if d.DeviceName != nil && d.Ebs != nil && d.Ebs.VolumeId != nil {
			devName := *d.DeviceName
			// Skip the root device
			if devName == *instance.RootDeviceName {
				continue
			}

			devicePath, err := s.getActualDevicePath(devName, *d.Ebs.VolumeId)
			if err != nil {
				return nil, storageops.NewStorageError(
					storageops.ErrInvalidDevicePath,
					fmt.Sprintf("unable to get actual device path for %s. %v", devName, err),
					s.instance)
			}
			m[devicePath] = *d.Ebs.VolumeId
		}
	}
	return m, nil
}

// Describe current instance.
func (s *ec2Ops) Describe() (interface{}, error) {
	return s.describe()
}

func (s *ec2Ops) describe() (*ec2.Instance, error) {
	request := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{&s.instance},
	}
	out, err := s.ec2.DescribeInstances(request)
	if err != nil {
		return nil, err
	}
	if len(out.Reservations) != 1 {
		return nil, fmt.Errorf("DescribeInstances(%v) returned %v reservations, expect 1",
			s.instance, len(out.Reservations))
	}
	if len(out.Reservations[0].Instances) != 1 {
		return nil, fmt.Errorf("DescribeInstances(%v) returned %v Reservations, expect 1",
			s.instance, len(out.Reservations[0].Instances))
	}
	return out.Reservations[0].Instances[0], nil
}

func (s *ec2Ops) getPrefixFromRootDeviceName(rootDeviceName string) (string, error) {
	devPrefix := awsDevicePrefix
	if !strings.HasPrefix(rootDeviceName, devPrefix) {
		devPrefix = awsDevicePrefixWithX
		if !strings.HasPrefix(rootDeviceName, devPrefix) {
			devPrefix = awsDevicePrefixWithH
			if !strings.HasPrefix(rootDeviceName, devPrefix) {
				return "", fmt.Errorf("unknown prefix type on root device: %s",
					rootDeviceName)
			}
		}
	}
	return devPrefix, nil
}

// getParentDevice returns the parent device of the given device path
// by following the symbolic link. It is expected that the input device
// path exists
func (s *ec2Ops) getParentDevice(ipDevPath string) (string, error) {
	// Check if the path is a symbolic link
	var parentDevPath string
	fi, err := os.Lstat(ipDevPath)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		// input device path is a symbolic link
		// get the parent device
		output, err := filepath.EvalSymlinks(ipDevPath)
		if err != nil {
			return "", fmt.Errorf("failed to read symlink due to: %v", err)
		}
		parentDevPath = strings.TrimSpace(string(output))
	} else {
		parentDevPath = ipDevPath
	}
	return parentDevPath, nil
}

// getActualDevicePath does an os.Stat on the provided devicePath.
// If not found it will try all the different devicePrefixes provided by AWS
// such as /dev/sd and /dev/xvd and return the devicePath which is found
// or return an error
func (s *ec2Ops) getActualDevicePath(ipDevicePath, volumeID string) (string, error) {
	letter := ipDevicePath[len(ipDevicePath)-1:]
	devicePath := awsDevicePrefix + letter
	if _, err := os.Stat(devicePath); err == nil {
		return s.getParentDevice(devicePath)
	}
	devicePath = awsDevicePrefixWithX + letter
	if _, err := os.Stat(devicePath); err == nil {
		return s.getParentDevice(devicePath)
	}

	devicePath = awsDevicePrefixWithH + letter
	if _, err := os.Stat(devicePath); err == nil {
		return s.getParentDevice(devicePath)
	}

	// Check if the EBS volumes are exposed as NVMe drives
	found := false
	for _, instancePrefix := range nvmeInstanceTypes {
		if strings.HasPrefix(s.instanceType, instancePrefix) {
			found = true
			break
		}
	}

	if !found {
		return "", fmt.Errorf("unable to map volume %v with block device mapping %v to an"+
			" actual device path on the host", volumeID, ipDevicePath)
	}

	devicePath, err := s.getNvmeDeviceFromVolumeID(volumeID)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(devicePath); err != nil {
		return "", err
	}
	return devicePath, nil

}

func (s *ec2Ops) getNvmeDeviceFromVolumeID(volumeID string) (string, error) {
	// We will use nvme list command to find nvme device mappings
	// A typical output of nvme list looks like this
	// # nvme list
	// Node             SN                   Model                                    Namespace Usage                      Format           FW Rev
	// ---------------- -------------------- ---------------------------------------- --------- -------------------------- ---------------- --------
	// /dev/nvme0n1     vol00fd6f8c30dc619f4 Amazon Elastic Block Store               1           0.00   B / 137.44  GB    512   B +  0 B   1.0
	// /dev/nvme1n1     vol044e12c8c0af45b3d Amazon Elastic Block Store               1           0.00   B / 107.37  GB    512   B +  0 B   1.0
	trimmedVolumeID := strings.Replace(volumeID, "-", "", 1)
	out, err := sh.Command(nvmeCmd, "list").Command("grep", trimmedVolumeID).Command("awk", "{print $1}").Output()
	if err != nil {
		return "", fmt.Errorf("unable to map %v volume to an nvme device: %v", volumeID, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *ec2Ops) FreeDevices(
	blockDeviceMappings []interface{},
	rootDeviceName string,
) ([]string, error) {
	initial := []byte("fghijklmnop")
	devPrefix := awsDevicePrefix
	for _, d := range blockDeviceMappings {
		dev := d.(*ec2.InstanceBlockDeviceMapping)

		if dev.DeviceName == nil {
			return nil, fmt.Errorf("Nil device name")
		}
		devName := *dev.DeviceName

		if devName == rootDeviceName {
			continue
		}
		if !strings.HasPrefix(devName, devPrefix) {
			devPrefix = awsDevicePrefixWithX
			if !strings.HasPrefix(devName, devPrefix) {
				devPrefix = awsDevicePrefixWithH
				if !strings.HasPrefix(devName, devPrefix) {
					return nil, fmt.Errorf("bad device name %q", devName)
				}
			}
		}
		letter := devName[len(devPrefix):]

		// Reset devPrefix for next devices
		devPrefix = awsDevicePrefix

		// AWS instances can have the following device names
		// /dev/xvd[b-c][a-z]
		if len(letter) == 1 {
			index := letter[0] - 'f'
			if index > ('p' - 'f') {
				continue
			}
			initial[index] = '0'
		} else if len(letter) == 2 {
			// We do not attach EBS volumes with "/dev/xvdc[a-z]" formats
			continue
		} else {
			return nil, fmt.Errorf("cannot parse device name %q", devName)
		}
	}

	// Set the prefix to the same one used as the root drive
	// The reason we do this is based on the virtualization type AWS might attach
	// the device "sda" at /dev/sda OR /dev/xvda. So we look at how the root device
	// is attached and use that prefix
	devPrefix, err := s.getPrefixFromRootDeviceName(rootDeviceName)
	if err != nil {
		return nil, err
	}

	free := make([]string, len(initial))
	count := 0
	for _, b := range initial {
		if b != '0' {
			free[count] = devPrefix + string(b)
			count++
		}
	}
	if count == 0 {
		return nil, fmt.Errorf("No more free devices")
	}
	return free[:count], nil
}

func (s *ec2Ops) rollbackCreate(id string, createErr error) error {
	logrus.Warnf("Rollback create volume %v, Error %v", id, createErr)
	err := s.Delete(id)
	if err != nil {
		logrus.Warnf("Rollback failed volume %v, Error %v", id, err)
	}
	return createErr
}

func (s *ec2Ops) refreshVol(id *string) (*ec2.Volume, error) {
	vols, err := s.Inspect([]*string{id})
	if err != nil {
		return nil, err
	}

	if len(vols) != 1 {
		return nil, fmt.Errorf("failed to get vol: %s."+
			"Found: %d volumes on inspecting", *id, len(vols))
	}

	resp, ok := vols[0].(*ec2.Volume)
	if !ok {
		return nil, storageops.NewStorageError(storageops.ErrVolInval,
			fmt.Sprintf("Invalid volume returned by inspect API for vol: %s", *id),
			"")
	}

	return resp, nil
}

func (s *ec2Ops) deleted(v *ec2.Volume) bool {
	return *v.State == ec2.VolumeStateDeleting ||
		*v.State == ec2.VolumeStateDeleted
}

func (s *ec2Ops) available(v *ec2.Volume) bool {
	return *v.State == ec2.VolumeStateAvailable
}

func (s *ec2Ops) GetDeviceID(vol interface{}) (string, error) {
	if d, ok := vol.(*ec2.Volume); ok {
		return *d.VolumeId, nil
	} else if d, ok := vol.(*ec2.Snapshot); ok {
		return *d.SnapshotId, nil
	} else {
		return "", fmt.Errorf("invalid type: %v given to GetDeviceID", vol)
	}
}

func (s *ec2Ops) Inspect(volumeIds []*string) ([]interface{}, error) {
	req := &ec2.DescribeVolumesInput{VolumeIds: volumeIds}
	resp, err := s.ec2.DescribeVolumes(req)
	if err != nil {
		return nil, err
	}

	var awsVols = make([]interface{}, len(resp.Volumes))
	for i, v := range resp.Volumes {
		awsVols[i] = v
	}

	return awsVols, nil
}

func (s *ec2Ops) Tags(volumeID string) (map[string]string, error) {
	vol, err := s.refreshVol(&volumeID)
	if err != nil {
		return nil, err
	}

	labels := make(map[string]string)
	for _, tag := range vol.Tags {
		labels[*tag.Key] = *tag.Value
	}
	return labels, nil
}

func (s *ec2Ops) Enumerate(
	volumeIds []*string,
	labels map[string]string,
	setIdentifier string,
) (map[string][]interface{}, error) {
	sets := make(map[string][]interface{})

	// Enumerate all volumes that have same labels.
	f := s.filters(labels, nil)
	req := &ec2.DescribeVolumesInput{Filters: f, VolumeIds: volumeIds}
	awsVols, err := s.ec2.DescribeVolumes(req)
	if err != nil {
		return nil, err
	}

	// Volume sets are identified by volumes with the same setIdentifer.
	found := false
	for _, vol := range awsVols.Volumes {
		if s.deleted(vol) {
			continue
		}
		if len(setIdentifier) == 0 {
			storageops.AddElementToMap(sets, vol, storageops.SetIdentifierNone)
		} else {
			found = false
			for _, tag := range vol.Tags {
				if s.matchTag(tag, setIdentifier) {
					storageops.AddElementToMap(sets, vol, *tag.Value)
					found = true
					break
				}
			}
			if !found {
				storageops.AddElementToMap(sets, vol, storageops.SetIdentifierNone)
			}
		}
	}

	return sets, nil
}

func (s *ec2Ops) Create(
	v interface{},
	labels map[string]string,
) (interface{}, error) {
	vol, ok := v.(*ec2.Volume)
	if !ok {
		return nil, storageops.NewStorageError(storageops.ErrVolInval,
			"Invalid volume template given", "")
	}

	req := &ec2.CreateVolumeInput{
		AvailabilityZone: vol.AvailabilityZone,
		Encrypted:        vol.Encrypted,
		KmsKeyId:         vol.KmsKeyId,
		Size:             vol.Size,
		VolumeType:       vol.VolumeType,
		SnapshotId:       vol.SnapshotId,
	}
	if *vol.VolumeType == opsworks.VolumeTypeIo1 {
		req.Iops = vol.Iops
	}

	resp, err := s.ec2.CreateVolume(req)
	if err != nil {
		return nil, err
	}
	if err = s.waitStatus(
		*resp.VolumeId,
		ec2.VolumeStateAvailable,
	); err != nil {
		return nil, s.rollbackCreate(*resp.VolumeId, err)
	}
	if len(labels) > 0 {
		if err = s.ApplyTags(*resp.VolumeId, labels); err != nil {
			return nil, s.rollbackCreate(*resp.VolumeId, err)
		}
	}

	return s.refreshVol(resp.VolumeId)
}

func (s *ec2Ops) DeleteFrom(id, _ string) error {
	return s.Delete(id)
}

func (s *ec2Ops) Delete(id string) error {
	req := &ec2.DeleteVolumeInput{VolumeId: &id}
	_, err := s.ec2.DeleteVolume(req)
	return err
}

func (s *ec2Ops) Attach(volumeID string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	self, err := s.describe()
	if err != nil {
		return "", err
	}

	var blockDeviceMappings = make([]interface{}, len(self.BlockDeviceMappings))
	for i, b := range self.BlockDeviceMappings {
		blockDeviceMappings[i] = b
	}

	devices, err := s.FreeDevices(blockDeviceMappings, *self.RootDeviceName)
	if err != nil {
		return "", err
	}
	req := &ec2.AttachVolumeInput{
		Device:     &devices[0],
		InstanceId: &s.instance,
		VolumeId:   &volumeID,
	}
	if _, err := s.ec2.AttachVolume(req); err != nil {
		return "", err
	}
	vol, err := s.waitAttachmentStatus(
		volumeID,
		ec2.VolumeAttachmentStateAttached,
		time.Minute,
	)
	if err != nil {
		return "", err
	}
	return s.DevicePath(*vol.VolumeId)
}

func (s *ec2Ops) Detach(volumeID string) error {
	return s.detachInternal(volumeID, s.instance)
}

func (s *ec2Ops) DetachFrom(volumeID, instanceName string) error {
	return s.detachInternal(volumeID, instanceName)
}

func (s *ec2Ops) detachInternal(volumeID, instanceName string) error {
	force := false
	req := &ec2.DetachVolumeInput{
		InstanceId: &instanceName,
		VolumeId:   &volumeID,
		Force:      &force,
	}
	if _, err := s.ec2.DetachVolume(req); err != nil {
		return err
	}
	_, err := s.waitAttachmentStatus(volumeID,
		ec2.VolumeAttachmentStateDetached,
		time.Minute,
	)
	return err
}

func (s *ec2Ops) Snapshot(
	volumeID string,
	readonly bool,
) (interface{}, error) {
	request := &ec2.CreateSnapshotInput{
		VolumeId: &volumeID,
	}
	return s.ec2.CreateSnapshot(request)
}

func (s *ec2Ops) SnapshotDelete(snapID string) error {
	request := &ec2.DeleteSnapshotInput{
		SnapshotId: &snapID,
	}

	_, err := s.ec2.DeleteSnapshot(request)
	return err
}

func (s *ec2Ops) DevicePath(volumeID string) (string, error) {
	vol, err := s.refreshVol(&volumeID)
	if err != nil {
		return "", err
	}

	if vol.Attachments == nil || len(vol.Attachments) == 0 {
		return "", storageops.NewStorageError(storageops.ErrVolDetached,
			"Volume is detached", *vol.VolumeId)
	}
	if vol.Attachments[0].InstanceId == nil {
		return "", storageops.NewStorageError(storageops.ErrVolInval,
			"Unable to determine volume instance attachment", "")
	}
	if s.instance != *vol.Attachments[0].InstanceId {
		return "", storageops.NewStorageError(storageops.ErrVolAttachedOnRemoteNode,
			fmt.Sprintf("Volume attached on %q current instance %q",
				*vol.Attachments[0].InstanceId, s.instance),
			*vol.Attachments[0].InstanceId)

	}
	if vol.Attachments[0].State == nil {
		return "", storageops.NewStorageError(storageops.ErrVolInval,
			"Unable to determine volume attachment state", "")
	}
	if *vol.Attachments[0].State != ec2.VolumeAttachmentStateAttached {
		return "", storageops.NewStorageError(storageops.ErrVolInval,
			fmt.Sprintf("Invalid state %q, volume is not attached",
				*vol.Attachments[0].State), "")
	}
	if vol.Attachments[0].Device == nil {
		return "", storageops.NewStorageError(storageops.ErrVolInval,
			"Unable to determine volume attachment path", "")
	}
	devicePath, err := s.getActualDevicePath(*vol.Attachments[0].Device, volumeID)
	if err != nil {
		return "", storageops.NewStorageError(storageops.ErrVolInval,
			err.Error(), "")
	}
	return devicePath, nil
}
