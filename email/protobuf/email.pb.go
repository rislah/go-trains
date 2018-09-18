// Code generated by protoc-gen-go. DO NOT EDIT.
// source: email.proto

package email

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type Email struct {
	To                   string   `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	Verificationid       string   `protobuf:"bytes,2,opt,name=verificationid,proto3" json:"verificationid,omitempty"`
	Passwordresetid      string   `protobuf:"bytes,3,opt,name=passwordresetid,proto3" json:"passwordresetid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Email) Reset()         { *m = Email{} }
func (m *Email) String() string { return proto.CompactTextString(m) }
func (*Email) ProtoMessage()    {}
func (*Email) Descriptor() ([]byte, []int) {
	return fileDescriptor_6175298cb4ed6faa, []int{0}
}
func (m *Email) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Email.Unmarshal(m, b)
}
func (m *Email) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Email.Marshal(b, m, deterministic)
}
func (dst *Email) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Email.Merge(dst, src)
}
func (m *Email) XXX_Size() int {
	return xxx_messageInfo_Email.Size(m)
}
func (m *Email) XXX_DiscardUnknown() {
	xxx_messageInfo_Email.DiscardUnknown(m)
}

var xxx_messageInfo_Email proto.InternalMessageInfo

func (m *Email) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Email) GetVerificationid() string {
	if m != nil {
		return m.Verificationid
	}
	return ""
}

func (m *Email) GetPasswordresetid() string {
	if m != nil {
		return m.Passwordresetid
	}
	return ""
}

type EmailResponse struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmailResponse) Reset()         { *m = EmailResponse{} }
func (m *EmailResponse) String() string { return proto.CompactTextString(m) }
func (*EmailResponse) ProtoMessage()    {}
func (*EmailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6175298cb4ed6faa, []int{1}
}
func (m *EmailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmailResponse.Unmarshal(m, b)
}
func (m *EmailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmailResponse.Marshal(b, m, deterministic)
}
func (dst *EmailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmailResponse.Merge(dst, src)
}
func (m *EmailResponse) XXX_Size() int {
	return xxx_messageInfo_EmailResponse.Size(m)
}
func (m *EmailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmailResponse proto.InternalMessageInfo

func (m *EmailResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Email)(nil), "Email")
	proto.RegisterType((*EmailResponse)(nil), "EmailResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EmailServiceClient is the client API for EmailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EmailServiceClient interface {
	SendEmailVerification(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailServiceClient struct {
	cc *grpc.ClientConn
}

func NewEmailServiceClient(cc *grpc.ClientConn) EmailServiceClient {
	return &emailServiceClient{cc}
}

func (c *emailServiceClient) SendEmailVerification(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/EmailService/SendEmailVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailServiceServer is the server API for EmailService service.
type EmailServiceServer interface {
	SendEmailVerification(context.Context, *Email) (*EmailResponse, error)
}

func RegisterEmailServiceServer(s *grpc.Server, srv EmailServiceServer) {
	s.RegisterService(&_EmailService_serviceDesc, srv)
}

func _EmailService_SendEmailVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendEmailVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/EmailService/SendEmailVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendEmailVerification(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

var _EmailService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "EmailService",
	HandlerType: (*EmailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmailVerification",
			Handler:    _EmailService_SendEmailVerification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "email.proto",
}

func init() { proto.RegisterFile("email.proto", fileDescriptor_6175298cb4ed6faa) }

var fileDescriptor_6175298cb4ed6faa = []byte{
	// 176 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0xc1, 0x0a, 0x82, 0x40,
	0x10, 0x86, 0x53, 0x49, 0x68, 0x2a, 0x8b, 0x85, 0x40, 0x3a, 0x95, 0x87, 0xf0, 0x24, 0x54, 0x4f,
	0xd0, 0xa1, 0x17, 0x50, 0xe8, 0x6e, 0xee, 0x14, 0x03, 0xe9, 0x2c, 0xbb, 0x8b, 0xbd, 0x7e, 0x38,
	0x74, 0x28, 0x6f, 0x33, 0xdf, 0xfc, 0xcc, 0xc7, 0x0f, 0x73, 0x6c, 0x6b, 0x7a, 0x15, 0xc6, 0xb2,
	0xe7, 0x8c, 0x60, 0x7a, 0x1d, 0x56, 0x95, 0x40, 0xe8, 0x39, 0x0d, 0x76, 0x41, 0x3e, 0x2b, 0x43,
	0xcf, 0xea, 0x00, 0x49, 0x8f, 0x96, 0x1e, 0xd4, 0xd4, 0x9e, 0xb8, 0x23, 0x9d, 0x86, 0x72, 0x1b,
	0x51, 0x95, 0xc3, 0xca, 0xd4, 0xce, 0xbd, 0xd9, 0x6a, 0x8b, 0x0e, 0x3d, 0xe9, 0x34, 0x92, 0xe0,
	0x18, 0x67, 0x7b, 0x58, 0x8a, 0xaa, 0x44, 0x67, 0xb8, 0x73, 0xa8, 0xd6, 0x10, 0xb5, 0xee, 0xf9,
	0x75, 0x0e, 0xe3, 0xe9, 0x02, 0x0b, 0x89, 0x54, 0x68, 0x7b, 0x6a, 0x50, 0x1d, 0x61, 0x53, 0x61,
	0xa7, 0x85, 0xdd, 0x7e, 0xbc, 0x2a, 0x2e, 0x84, 0x6d, 0x93, 0xe2, 0xef, 0x65, 0x36, 0xb9, 0xc7,
	0xd2, 0xeb, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x02, 0xb0, 0x22, 0xa9, 0xe6, 0x00, 0x00, 0x00,
}