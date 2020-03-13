// Code generated by protoc-gen-go. DO NOT EDIT.
// source: plugin.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ChannelMessage struct {
	Channel              string   `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelMessage) Reset()         { *m = ChannelMessage{} }
func (m *ChannelMessage) String() string { return proto.CompactTextString(m) }
func (*ChannelMessage) ProtoMessage()    {}
func (*ChannelMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{0}
}

func (m *ChannelMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelMessage.Unmarshal(m, b)
}
func (m *ChannelMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelMessage.Marshal(b, m, deterministic)
}
func (m *ChannelMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelMessage.Merge(m, src)
}
func (m *ChannelMessage) XXX_Size() int {
	return xxx_messageInfo_ChannelMessage.Size(m)
}
func (m *ChannelMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelMessage proto.InternalMessageInfo

func (m *ChannelMessage) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *ChannelMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type RawMessage struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawMessage) Reset()         { *m = RawMessage{} }
func (m *RawMessage) String() string { return proto.CompactTextString(m) }
func (*RawMessage) ProtoMessage()    {}
func (*RawMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{1}
}

func (m *RawMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawMessage.Unmarshal(m, b)
}
func (m *RawMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawMessage.Marshal(b, m, deterministic)
}
func (m *RawMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawMessage.Merge(m, src)
}
func (m *RawMessage) XXX_Size() int {
	return xxx_messageInfo_RawMessage.Size(m)
}
func (m *RawMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_RawMessage.DiscardUnknown(m)
}

var xxx_messageInfo_RawMessage proto.InternalMessageInfo

func (m *RawMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Error struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{2}
}

func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (m *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(m, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ChannelMessage)(nil), "rpc.ChannelMessage")
	proto.RegisterType((*RawMessage)(nil), "rpc.RawMessage")
	proto.RegisterType((*Error)(nil), "rpc.Error")
}

func init() {
	proto.RegisterFile("plugin.proto", fileDescriptor_22a625af4bc1cc87)
}

var fileDescriptor_22a625af4bc1cc87 = []byte{
	// 177 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xc8, 0x29, 0x4d,
	0xcf, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x72, 0xe1,
	0xe2, 0x73, 0xce, 0x48, 0xcc, 0xcb, 0x4b, 0xcd, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15,
	0x92, 0xe0, 0x62, 0x4f, 0x86, 0x88, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x20,
	0x99, 0x5c, 0x88, 0x22, 0x09, 0x26, 0x88, 0x0c, 0x94, 0xab, 0xa4, 0xc6, 0xc5, 0x15, 0x94, 0x58,
	0x8e, 0x64, 0x02, 0x4c, 0x1d, 0x23, 0xaa, 0x3a, 0x45, 0x2e, 0x56, 0xd7, 0xa2, 0xa2, 0xfc, 0x22,
	0xdc, 0x4a, 0x8c, 0xca, 0xb8, 0x38, 0x3d, 0x83, 0x9c, 0x03, 0xc0, 0x0e, 0x15, 0xb2, 0xe0, 0x12,
	0x2e, 0x4e, 0xcd, 0x4b, 0x41, 0x72, 0x21, 0xd8, 0x02, 0x61, 0xbd, 0xa2, 0x82, 0x64, 0x3d, 0x54,
	0x77, 0x4b, 0x71, 0x81, 0x05, 0xc1, 0xc6, 0x2b, 0x31, 0x08, 0xe9, 0x73, 0xf1, 0x81, 0x74, 0x22,
	0xb9, 0x8a, 0x1f, 0x2c, 0x8f, 0x10, 0x40, 0xd5, 0xe0, 0xc4, 0x1a, 0x05, 0x0a, 0x8f, 0x24, 0x36,
	0x70, 0xd8, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4e, 0xde, 0x33, 0x01, 0x2b, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// IRCPluginClient is the client API for IRCPlugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IRCPluginClient interface {
	SendChannelMesssage(ctx context.Context, in *ChannelMessage, opts ...grpc.CallOption) (*Error, error)
	SendRawMessage(ctx context.Context, in *RawMessage, opts ...grpc.CallOption) (*Error, error)
}

type iRCPluginClient struct {
	cc grpc.ClientConnInterface
}

func NewIRCPluginClient(cc grpc.ClientConnInterface) IRCPluginClient {
	return &iRCPluginClient{cc}
}

func (c *iRCPluginClient) SendChannelMesssage(ctx context.Context, in *ChannelMessage, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/rpc.IRCPlugin/sendChannelMesssage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iRCPluginClient) SendRawMessage(ctx context.Context, in *RawMessage, opts ...grpc.CallOption) (*Error, error) {
	out := new(Error)
	err := c.cc.Invoke(ctx, "/rpc.IRCPlugin/sendRawMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IRCPluginServer is the server API for IRCPlugin service.
type IRCPluginServer interface {
	SendChannelMesssage(context.Context, *ChannelMessage) (*Error, error)
	SendRawMessage(context.Context, *RawMessage) (*Error, error)
}

// UnimplementedIRCPluginServer can be embedded to have forward compatible implementations.
type UnimplementedIRCPluginServer struct {
}

func (*UnimplementedIRCPluginServer) SendChannelMesssage(ctx context.Context, req *ChannelMessage) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChannelMesssage not implemented")
}
func (*UnimplementedIRCPluginServer) SendRawMessage(ctx context.Context, req *RawMessage) (*Error, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRawMessage not implemented")
}

func RegisterIRCPluginServer(s *grpc.Server, srv IRCPluginServer) {
	s.RegisterService(&_IRCPlugin_serviceDesc, srv)
}

func _IRCPlugin_SendChannelMesssage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChannelMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IRCPluginServer).SendChannelMesssage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.IRCPlugin/SendChannelMesssage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IRCPluginServer).SendChannelMesssage(ctx, req.(*ChannelMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _IRCPlugin_SendRawMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IRCPluginServer).SendRawMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.IRCPlugin/SendRawMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IRCPluginServer).SendRawMessage(ctx, req.(*RawMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _IRCPlugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.IRCPlugin",
	HandlerType: (*IRCPluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendChannelMesssage",
			Handler:    _IRCPlugin_SendChannelMesssage_Handler,
		},
		{
			MethodName: "sendRawMessage",
			Handler:    _IRCPlugin_SendRawMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin.proto",
}