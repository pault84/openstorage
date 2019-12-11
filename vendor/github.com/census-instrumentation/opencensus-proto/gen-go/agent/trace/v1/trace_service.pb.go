// Code generated by protoc-gen-go. DO NOT EDIT.
// source: opencensus/proto/agent/trace/v1/trace_service.proto

package v1

import (
	fmt "fmt"
	v1 "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	v12 "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	v11 "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CurrentLibraryConfig struct {
	// This is required only in the first message on the stream or if the
	// previous sent CurrentLibraryConfig message has a different Node (e.g.
	// when the same RPC is used to configure multiple Applications).
	Node *v1.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	// Current configuration.
	Config               *v11.TraceConfig `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CurrentLibraryConfig) Reset()         { *m = CurrentLibraryConfig{} }
func (m *CurrentLibraryConfig) String() string { return proto.CompactTextString(m) }
func (*CurrentLibraryConfig) ProtoMessage()    {}
func (*CurrentLibraryConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_7027f99caf7ac6a5, []int{0}
}

func (m *CurrentLibraryConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CurrentLibraryConfig.Unmarshal(m, b)
}
func (m *CurrentLibraryConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CurrentLibraryConfig.Marshal(b, m, deterministic)
}
func (m *CurrentLibraryConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurrentLibraryConfig.Merge(m, src)
}
func (m *CurrentLibraryConfig) XXX_Size() int {
	return xxx_messageInfo_CurrentLibraryConfig.Size(m)
}
func (m *CurrentLibraryConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CurrentLibraryConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CurrentLibraryConfig proto.InternalMessageInfo

func (m *CurrentLibraryConfig) GetNode() *v1.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *CurrentLibraryConfig) GetConfig() *v11.TraceConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

type UpdatedLibraryConfig struct {
	// This field is ignored when the RPC is used to configure only one Application.
	// This is required only in the first message on the stream or if the
	// previous sent UpdatedLibraryConfig message has a different Node.
	Node *v1.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	// Requested updated configuration.
	Config               *v11.TraceConfig `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *UpdatedLibraryConfig) Reset()         { *m = UpdatedLibraryConfig{} }
func (m *UpdatedLibraryConfig) String() string { return proto.CompactTextString(m) }
func (*UpdatedLibraryConfig) ProtoMessage()    {}
func (*UpdatedLibraryConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_7027f99caf7ac6a5, []int{1}
}

func (m *UpdatedLibraryConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatedLibraryConfig.Unmarshal(m, b)
}
func (m *UpdatedLibraryConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatedLibraryConfig.Marshal(b, m, deterministic)
}
func (m *UpdatedLibraryConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatedLibraryConfig.Merge(m, src)
}
func (m *UpdatedLibraryConfig) XXX_Size() int {
	return xxx_messageInfo_UpdatedLibraryConfig.Size(m)
}
func (m *UpdatedLibraryConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatedLibraryConfig.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatedLibraryConfig proto.InternalMessageInfo

func (m *UpdatedLibraryConfig) GetNode() *v1.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *UpdatedLibraryConfig) GetConfig() *v11.TraceConfig {
	if m != nil {
		return m.Config
	}
	return nil
}

type ExportTraceServiceRequest struct {
	// This is required only in the first message on the stream or if the
	// previous sent ExportTraceServiceRequest message has a different Node (e.g.
	// when the same RPC is used to send Spans from multiple Applications).
	Node *v1.Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	// A list of Spans that belong to the last received Node.
	Spans []*v11.Span `protobuf:"bytes,2,rep,name=spans,proto3" json:"spans,omitempty"`
	// The resource for the spans in this message that do not have an explicit
	// resource set.
	// If unset, the most recently set resource in the RPC stream applies. It is
	// valid to never be set within a stream, e.g. when no resource info is known.
	Resource             *v12.Resource `protobuf:"bytes,3,opt,name=resource,proto3" json:"resource,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ExportTraceServiceRequest) Reset()         { *m = ExportTraceServiceRequest{} }
func (m *ExportTraceServiceRequest) String() string { return proto.CompactTextString(m) }
func (*ExportTraceServiceRequest) ProtoMessage()    {}
func (*ExportTraceServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7027f99caf7ac6a5, []int{2}
}

func (m *ExportTraceServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportTraceServiceRequest.Unmarshal(m, b)
}
func (m *ExportTraceServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportTraceServiceRequest.Marshal(b, m, deterministic)
}
func (m *ExportTraceServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportTraceServiceRequest.Merge(m, src)
}
func (m *ExportTraceServiceRequest) XXX_Size() int {
	return xxx_messageInfo_ExportTraceServiceRequest.Size(m)
}
func (m *ExportTraceServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportTraceServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExportTraceServiceRequest proto.InternalMessageInfo

func (m *ExportTraceServiceRequest) GetNode() *v1.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *ExportTraceServiceRequest) GetSpans() []*v11.Span {
	if m != nil {
		return m.Spans
	}
	return nil
}

func (m *ExportTraceServiceRequest) GetResource() *v12.Resource {
	if m != nil {
		return m.Resource
	}
	return nil
}

type ExportTraceServiceResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExportTraceServiceResponse) Reset()         { *m = ExportTraceServiceResponse{} }
func (m *ExportTraceServiceResponse) String() string { return proto.CompactTextString(m) }
func (*ExportTraceServiceResponse) ProtoMessage()    {}
func (*ExportTraceServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7027f99caf7ac6a5, []int{3}
}

func (m *ExportTraceServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportTraceServiceResponse.Unmarshal(m, b)
}
func (m *ExportTraceServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportTraceServiceResponse.Marshal(b, m, deterministic)
}
func (m *ExportTraceServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportTraceServiceResponse.Merge(m, src)
}
func (m *ExportTraceServiceResponse) XXX_Size() int {
	return xxx_messageInfo_ExportTraceServiceResponse.Size(m)
}
func (m *ExportTraceServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportTraceServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExportTraceServiceResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CurrentLibraryConfig)(nil), "opencensus.proto.agent.trace.v1.CurrentLibraryConfig")
	proto.RegisterType((*UpdatedLibraryConfig)(nil), "opencensus.proto.agent.trace.v1.UpdatedLibraryConfig")
	proto.RegisterType((*ExportTraceServiceRequest)(nil), "opencensus.proto.agent.trace.v1.ExportTraceServiceRequest")
	proto.RegisterType((*ExportTraceServiceResponse)(nil), "opencensus.proto.agent.trace.v1.ExportTraceServiceResponse")
}

func init() {
	proto.RegisterFile("opencensus/proto/agent/trace/v1/trace_service.proto", fileDescriptor_7027f99caf7ac6a5)
}

var fileDescriptor_7027f99caf7ac6a5 = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0xbf, 0x6b, 0xdb, 0x40,
	0x14, 0xee, 0xd9, 0xad, 0x28, 0xe7, 0x2e, 0x15, 0x1d, 0x54, 0x51, 0xb0, 0x11, 0xb4, 0x18, 0x5a,
	0x9d, 0x2a, 0x1b, 0x2f, 0x2e, 0x74, 0xb0, 0x29, 0x74, 0x28, 0xc5, 0xc8, 0xed, 0x92, 0xc5, 0xc8,
	0xd2, 0x8b, 0xa2, 0xc1, 0x77, 0xca, 0xdd, 0x49, 0x24, 0x90, 0x2d, 0x43, 0xf6, 0x0c, 0xf9, 0xc3,
	0xf2, 0x17, 0x05, 0xdd, 0xc9, 0x3f, 0x12, 0x5b, 0x11, 0x24, 0x4b, 0xb6, 0x87, 0xde, 0xf7, 0x7d,
	0xf7, 0xbd, 0x7b, 0xdf, 0x09, 0x0f, 0x59, 0x06, 0x34, 0x02, 0x2a, 0x72, 0xe1, 0x65, 0x9c, 0x49,
	0xe6, 0x85, 0x09, 0x50, 0xe9, 0x49, 0x1e, 0x46, 0xe0, 0x15, 0xbe, 0x2e, 0x16, 0x02, 0x78, 0x91,
	0x46, 0x40, 0x14, 0xc4, 0xec, 0x6e, 0x49, 0xfa, 0x0b, 0x51, 0x24, 0xa2, 0xb0, 0xa4, 0xf0, 0x6d,
	0xb7, 0x46, 0x35, 0x62, 0xab, 0x15, 0xa3, 0xa5, 0xac, 0xae, 0x34, 0xdb, 0xfe, 0xba, 0x07, 0xe7,
	0x20, 0x58, 0xce, 0xb5, 0x83, 0x75, 0x5d, 0x81, 0x3f, 0xef, 0x81, 0xef, 0x7b, 0xad, 0x60, 0xdf,
	0x1a, 0x60, 0x8b, 0x88, 0xd1, 0xe3, 0x34, 0xd1, 0x68, 0xe7, 0x1a, 0xe1, 0x0f, 0xd3, 0x9c, 0x73,
	0xa0, 0xf2, 0x4f, 0xba, 0xe4, 0x21, 0x3f, 0x9f, 0xaa, 0xb6, 0x39, 0xc6, 0xaf, 0x29, 0x8b, 0xc1,
	0x42, 0x3d, 0xd4, 0xef, 0x0c, 0xbe, 0x90, 0x9a, 0xc9, 0xab, 0x71, 0x0a, 0x9f, 0xfc, 0x65, 0x31,
	0x04, 0x8a, 0x63, 0xfe, 0xc4, 0x86, 0x3e, 0xc4, 0x6a, 0xd5, 0xb1, 0xd7, 0x37, 0x46, 0xfe, 0x95,
	0x85, 0x3e, 0x33, 0xa8, 0x58, 0xca, 0xd4, 0xff, 0x2c, 0x0e, 0x25, 0xc4, 0x2f, 0xc7, 0xd4, 0x2d,
	0xc2, 0x1f, 0x7f, 0x9d, 0x65, 0x8c, 0x4b, 0xd5, 0x9d, 0xeb, 0x60, 0x04, 0x70, 0x9a, 0x83, 0x90,
	0xcf, 0x72, 0x36, 0xc2, 0x6f, 0x44, 0x16, 0x52, 0x61, 0xb5, 0x7a, 0xed, 0x7e, 0x67, 0xd0, 0x7d,
	0xc4, 0xd8, 0x3c, 0x0b, 0x69, 0xa0, 0xd1, 0xe6, 0x04, 0xbf, 0x5d, 0x27, 0xc4, 0x6a, 0xd7, 0x1d,
	0xbb, 0xc9, 0x50, 0xe1, 0x93, 0xa0, 0xaa, 0x83, 0x0d, 0xcf, 0xf9, 0x84, 0xed, 0x43, 0x33, 0x89,
	0x8c, 0x51, 0x01, 0x83, 0x9b, 0x16, 0x7e, 0xb7, 0xdb, 0x30, 0x2f, 0xb0, 0x51, 0x6d, 0x62, 0x44,
	0x1a, 0x9e, 0x02, 0x39, 0x94, 0x2a, 0xbb, 0x99, 0x76, 0x68, 0xef, 0xce, 0xab, 0x3e, 0xfa, 0x8e,
	0xcc, 0x2b, 0x84, 0x0d, 0xed, 0xd6, 0x1c, 0x37, 0xea, 0xd4, 0xae, 0xca, 0xfe, 0xf1, 0x24, 0xae,
	0xbe, 0x12, 0xed, 0x64, 0x72, 0x89, 0xb0, 0x93, 0xb2, 0x26, 0x9d, 0xc9, 0xfb, 0x5d, 0x89, 0x59,
	0x89, 0x98, 0xa1, 0xa3, 0xdf, 0x49, 0x2a, 0x4f, 0xf2, 0x65, 0x19, 0x05, 0x4f, 0x93, 0xdd, 0x94,
	0x0a, 0xc9, 0xf3, 0x15, 0x50, 0x19, 0xca, 0x94, 0x51, 0x6f, 0xab, 0xeb, 0xea, 0x17, 0x9c, 0x00,
	0x75, 0x93, 0x87, 0x7f, 0xa8, 0xa5, 0xa1, 0x9a, 0xc3, 0xbb, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcf,
	0x9c, 0x9b, 0xf7, 0xcb, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TraceServiceClient is the client API for TraceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TraceServiceClient interface {
	// After initialization, this RPC must be kept alive for the entire life of
	// the application. The agent pushes configs down to applications via a
	// stream.
	Config(ctx context.Context, opts ...grpc.CallOption) (TraceService_ConfigClient, error)
	// For performance reasons, it is recommended to keep this RPC
	// alive for the entire life of the application.
	Export(ctx context.Context, opts ...grpc.CallOption) (TraceService_ExportClient, error)
}

type traceServiceClient struct {
	cc *grpc.ClientConn
}

func NewTraceServiceClient(cc *grpc.ClientConn) TraceServiceClient {
	return &traceServiceClient{cc}
}

func (c *traceServiceClient) Config(ctx context.Context, opts ...grpc.CallOption) (TraceService_ConfigClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TraceService_serviceDesc.Streams[0], "/opencensus.proto.agent.trace.v1.TraceService/Config", opts...)
	if err != nil {
		return nil, err
	}
	x := &traceServiceConfigClient{stream}
	return x, nil
}

type TraceService_ConfigClient interface {
	Send(*CurrentLibraryConfig) error
	Recv() (*UpdatedLibraryConfig, error)
	grpc.ClientStream
}

type traceServiceConfigClient struct {
	grpc.ClientStream
}

func (x *traceServiceConfigClient) Send(m *CurrentLibraryConfig) error {
	return x.ClientStream.SendMsg(m)
}

func (x *traceServiceConfigClient) Recv() (*UpdatedLibraryConfig, error) {
	m := new(UpdatedLibraryConfig)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *traceServiceClient) Export(ctx context.Context, opts ...grpc.CallOption) (TraceService_ExportClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TraceService_serviceDesc.Streams[1], "/opencensus.proto.agent.trace.v1.TraceService/Export", opts...)
	if err != nil {
		return nil, err
	}
	x := &traceServiceExportClient{stream}
	return x, nil
}

type TraceService_ExportClient interface {
	Send(*ExportTraceServiceRequest) error
	Recv() (*ExportTraceServiceResponse, error)
	grpc.ClientStream
}

type traceServiceExportClient struct {
	grpc.ClientStream
}

func (x *traceServiceExportClient) Send(m *ExportTraceServiceRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *traceServiceExportClient) Recv() (*ExportTraceServiceResponse, error) {
	m := new(ExportTraceServiceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TraceServiceServer is the server API for TraceService service.
type TraceServiceServer interface {
	// After initialization, this RPC must be kept alive for the entire life of
	// the application. The agent pushes configs down to applications via a
	// stream.
	Config(TraceService_ConfigServer) error
	// For performance reasons, it is recommended to keep this RPC
	// alive for the entire life of the application.
	Export(TraceService_ExportServer) error
}

func RegisterTraceServiceServer(s *grpc.Server, srv TraceServiceServer) {
	s.RegisterService(&_TraceService_serviceDesc, srv)
}

func _TraceService_Config_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TraceServiceServer).Config(&traceServiceConfigServer{stream})
}

type TraceService_ConfigServer interface {
	Send(*UpdatedLibraryConfig) error
	Recv() (*CurrentLibraryConfig, error)
	grpc.ServerStream
}

type traceServiceConfigServer struct {
	grpc.ServerStream
}

func (x *traceServiceConfigServer) Send(m *UpdatedLibraryConfig) error {
	return x.ServerStream.SendMsg(m)
}

func (x *traceServiceConfigServer) Recv() (*CurrentLibraryConfig, error) {
	m := new(CurrentLibraryConfig)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _TraceService_Export_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TraceServiceServer).Export(&traceServiceExportServer{stream})
}

type TraceService_ExportServer interface {
	Send(*ExportTraceServiceResponse) error
	Recv() (*ExportTraceServiceRequest, error)
	grpc.ServerStream
}

type traceServiceExportServer struct {
	grpc.ServerStream
}

func (x *traceServiceExportServer) Send(m *ExportTraceServiceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *traceServiceExportServer) Recv() (*ExportTraceServiceRequest, error) {
	m := new(ExportTraceServiceRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _TraceService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "opencensus.proto.agent.trace.v1.TraceService",
	HandlerType: (*TraceServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Config",
			Handler:       _TraceService_Config_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Export",
			Handler:       _TraceService_Export_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "opencensus/proto/agent/trace/v1/trace_service.proto",
}