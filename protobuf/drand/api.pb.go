// Code generated by protoc-gen-go. DO NOT EDIT.
// source: drand/api.proto

package drand

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

// PublicRandRequest requests a public random value that has been generated in a
// unbiasable way and verifiable.
type PublicRandRequest struct {
	// round uniquely identifies a beacon. If round == 0 (or unspecified), then
	// the response will contain the last.
	Round                uint64   `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicRandRequest) Reset()         { *m = PublicRandRequest{} }
func (m *PublicRandRequest) String() string { return proto.CompactTextString(m) }
func (*PublicRandRequest) ProtoMessage()    {}
func (*PublicRandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{0}
}

func (m *PublicRandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicRandRequest.Unmarshal(m, b)
}
func (m *PublicRandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicRandRequest.Marshal(b, m, deterministic)
}
func (m *PublicRandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicRandRequest.Merge(m, src)
}
func (m *PublicRandRequest) XXX_Size() int {
	return xxx_messageInfo_PublicRandRequest.Size(m)
}
func (m *PublicRandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicRandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PublicRandRequest proto.InternalMessageInfo

func (m *PublicRandRequest) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

// PublicRandResponse holds a signature which is the random value. It can be
// verified thanks to the distributed public key of the nodes that have ran the
// DKG protocol and is unbiasable. The randomness can be verified using the BLS
// verification routine with the message "round || previous_rand".
type PublicRandResponse struct {
	Round             uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	Signature         []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	PreviousSignature []byte `protobuf:"bytes,3,opt,name=previous_signature,json=previousSignature,proto3" json:"previous_signature,omitempty"`
	// randomness is simply there to demonstrate - it is the hash of the
	// signature. It should be computed locally.
	Randomness           []byte   `protobuf:"bytes,4,opt,name=randomness,proto3" json:"randomness,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicRandResponse) Reset()         { *m = PublicRandResponse{} }
func (m *PublicRandResponse) String() string { return proto.CompactTextString(m) }
func (*PublicRandResponse) ProtoMessage()    {}
func (*PublicRandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{1}
}

func (m *PublicRandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicRandResponse.Unmarshal(m, b)
}
func (m *PublicRandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicRandResponse.Marshal(b, m, deterministic)
}
func (m *PublicRandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicRandResponse.Merge(m, src)
}
func (m *PublicRandResponse) XXX_Size() int {
	return xxx_messageInfo_PublicRandResponse.Size(m)
}
func (m *PublicRandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicRandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PublicRandResponse proto.InternalMessageInfo

func (m *PublicRandResponse) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *PublicRandResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PublicRandResponse) GetPreviousSignature() []byte {
	if m != nil {
		return m.PreviousSignature
	}
	return nil
}

func (m *PublicRandResponse) GetRandomness() []byte {
	if m != nil {
		return m.Randomness
	}
	return nil
}

// PrivateRandRequest is the message to send when requesting a private random
// value.
type PrivateRandRequest struct {
	// Request is the ECIES encryption of an ephemereal public key towards which
	// to encrypt the private randomness. The format of the bytes is denoted by
	// the ECIES encryption used by drand.
	Request              []byte   `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateRandRequest) Reset()         { *m = PrivateRandRequest{} }
func (m *PrivateRandRequest) String() string { return proto.CompactTextString(m) }
func (*PrivateRandRequest) ProtoMessage()    {}
func (*PrivateRandRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{2}
}

func (m *PrivateRandRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrivateRandRequest.Unmarshal(m, b)
}
func (m *PrivateRandRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrivateRandRequest.Marshal(b, m, deterministic)
}
func (m *PrivateRandRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivateRandRequest.Merge(m, src)
}
func (m *PrivateRandRequest) XXX_Size() int {
	return xxx_messageInfo_PrivateRandRequest.Size(m)
}
func (m *PrivateRandRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivateRandRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PrivateRandRequest proto.InternalMessageInfo

func (m *PrivateRandRequest) GetRequest() []byte {
	if m != nil {
		return m.Request
	}
	return nil
}

type PrivateRandResponse struct {
	// Responses is the ECIES encryption of the private randomness using the
	// ephemereal public key sent in the request.  The format of the bytes is
	// denoted by the ECIES  encryption used by drand.
	Response             []byte   `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateRandResponse) Reset()         { *m = PrivateRandResponse{} }
func (m *PrivateRandResponse) String() string { return proto.CompactTextString(m) }
func (*PrivateRandResponse) ProtoMessage()    {}
func (*PrivateRandResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{3}
}

func (m *PrivateRandResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrivateRandResponse.Unmarshal(m, b)
}
func (m *PrivateRandResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrivateRandResponse.Marshal(b, m, deterministic)
}
func (m *PrivateRandResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivateRandResponse.Merge(m, src)
}
func (m *PrivateRandResponse) XXX_Size() int {
	return xxx_messageInfo_PrivateRandResponse.Size(m)
}
func (m *PrivateRandResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivateRandResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PrivateRandResponse proto.InternalMessageInfo

func (m *PrivateRandResponse) GetResponse() []byte {
	if m != nil {
		return m.Response
	}
	return nil
}

type HomeRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HomeRequest) Reset()         { *m = HomeRequest{} }
func (m *HomeRequest) String() string { return proto.CompactTextString(m) }
func (*HomeRequest) ProtoMessage()    {}
func (*HomeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{4}
}

func (m *HomeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HomeRequest.Unmarshal(m, b)
}
func (m *HomeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HomeRequest.Marshal(b, m, deterministic)
}
func (m *HomeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HomeRequest.Merge(m, src)
}
func (m *HomeRequest) XXX_Size() int {
	return xxx_messageInfo_HomeRequest.Size(m)
}
func (m *HomeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HomeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HomeRequest proto.InternalMessageInfo

type HomeResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HomeResponse) Reset()         { *m = HomeResponse{} }
func (m *HomeResponse) String() string { return proto.CompactTextString(m) }
func (*HomeResponse) ProtoMessage()    {}
func (*HomeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c0cff3fc81cf7d79, []int{5}
}

func (m *HomeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HomeResponse.Unmarshal(m, b)
}
func (m *HomeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HomeResponse.Marshal(b, m, deterministic)
}
func (m *HomeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HomeResponse.Merge(m, src)
}
func (m *HomeResponse) XXX_Size() int {
	return xxx_messageInfo_HomeResponse.Size(m)
}
func (m *HomeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HomeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HomeResponse proto.InternalMessageInfo

func (m *HomeResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*PublicRandRequest)(nil), "drand.PublicRandRequest")
	proto.RegisterType((*PublicRandResponse)(nil), "drand.PublicRandResponse")
	proto.RegisterType((*PrivateRandRequest)(nil), "drand.PrivateRandRequest")
	proto.RegisterType((*PrivateRandResponse)(nil), "drand.PrivateRandResponse")
	proto.RegisterType((*HomeRequest)(nil), "drand.HomeRequest")
	proto.RegisterType((*HomeResponse)(nil), "drand.HomeResponse")
}

func init() {
	proto.RegisterFile("drand/api.proto", fileDescriptor_c0cff3fc81cf7d79)
}

var fileDescriptor_c0cff3fc81cf7d79 = []byte{
	// 348 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x3b, 0x4f, 0xc3, 0x30,
	0x10, 0x56, 0x4a, 0x1f, 0xf4, 0x5a, 0x04, 0xbd, 0xa2, 0x12, 0x22, 0x84, 0xaa, 0x0c, 0xa8, 0x0c,
	0xa4, 0x3c, 0x56, 0x16, 0x1e, 0x03, 0xdd, 0xaa, 0x74, 0x63, 0x41, 0x6e, 0x63, 0x20, 0x82, 0xd8,
	0xc1, 0x76, 0xfa, 0x57, 0xf8, 0x51, 0xfc, 0x29, 0x54, 0xdb, 0x4d, 0x52, 0x9a, 0x89, 0xcd, 0xdf,
	0xe3, 0xce, 0xe7, 0xcf, 0x07, 0xfb, 0x91, 0x20, 0x2c, 0x1a, 0x93, 0x34, 0x0e, 0x52, 0xc1, 0x15,
	0xc7, 0x86, 0x26, 0x3c, 0x34, 0xfc, 0x82, 0x27, 0x09, 0x67, 0x46, 0xf2, 0xcf, 0xa1, 0x37, 0xcd,
	0xe6, 0x9f, 0xf1, 0x22, 0x24, 0x2c, 0x0a, 0xe9, 0x57, 0x46, 0xa5, 0xc2, 0x43, 0x68, 0x08, 0x9e,
	0xb1, 0xc8, 0x75, 0x86, 0xce, 0xa8, 0x1e, 0x1a, 0xe0, 0x7f, 0x3b, 0x80, 0x65, 0xaf, 0x4c, 0x39,
	0x93, 0xb4, 0xda, 0x8c, 0x27, 0xd0, 0x96, 0xf1, 0x1b, 0x23, 0x2a, 0x13, 0xd4, 0xad, 0x0d, 0x9d,
	0x51, 0x37, 0x2c, 0x08, 0xbc, 0x00, 0x4c, 0x05, 0x5d, 0xc6, 0x3c, 0x93, 0x2f, 0x85, 0x6d, 0x47,
	0xdb, 0x7a, 0x6b, 0x65, 0x96, 0xdb, 0x4f, 0x01, 0x56, 0x93, 0xf3, 0x84, 0x51, 0x29, 0xdd, 0xba,
	0xb6, 0x95, 0x18, 0x3f, 0x00, 0x9c, 0x8a, 0x78, 0x49, 0x14, 0x2d, 0xbf, 0xc2, 0x85, 0x96, 0x30,
	0x47, 0x3d, 0x5a, 0x37, 0x5c, 0x43, 0xff, 0x0a, 0xfa, 0x1b, 0x7e, 0xfb, 0x12, 0x0f, 0x76, 0x85,
	0x3d, 0xdb, 0x8a, 0x1c, 0xfb, 0x7b, 0xd0, 0x79, 0xe2, 0x09, 0xb5, 0xbd, 0xfd, 0x33, 0xe8, 0x1a,
	0x68, 0x4b, 0x07, 0xd0, 0x94, 0x8a, 0xa8, 0x4c, 0xea, 0xc2, 0x76, 0x68, 0xd1, 0xf5, 0x4f, 0x0d,
	0x9a, 0x26, 0x33, 0xbc, 0x03, 0x28, 0xd2, 0x43, 0x37, 0xd0, 0x9f, 0x11, 0x6c, 0x85, 0xef, 0x1d,
	0x57, 0x28, 0xf6, 0x96, 0x09, 0x1c, 0x14, 0xec, 0x4c, 0x09, 0x4a, 0x92, 0x7f, 0x35, 0xba, 0x74,
	0xf0, 0x11, 0x3a, 0xa5, 0x08, 0x30, 0xf7, 0x6e, 0xc5, 0xe8, 0x79, 0x55, 0x92, 0x1d, 0xe8, 0x16,
	0xda, 0x0f, 0xef, 0x24, 0x66, 0x13, 0xf6, 0xca, 0xf1, 0xc8, 0x1a, 0x73, 0x66, 0xdd, 0x61, 0xf0,
	0x57, 0x98, 0x92, 0xc5, 0x07, 0x55, 0x38, 0x86, 0xfa, 0x2a, 0x44, 0x44, 0xab, 0x97, 0x02, 0xf6,
	0xfa, 0x1b, 0x9c, 0xb9, 0xee, 0xbe, 0xf5, 0x6c, 0x36, 0x79, 0xde, 0xd4, 0xcb, 0x7b, 0xf3, 0x1b,
	0x00, 0x00, 0xff, 0xff, 0xcb, 0xb8, 0x17, 0x8f, 0xea, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PublicClient is the client API for Public service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PublicClient interface {
	// PublicRand is the method that returns the publicly verifiable randomness
	// generated by the drand network.
	PublicRand(ctx context.Context, in *PublicRandRequest, opts ...grpc.CallOption) (*PublicRandResponse, error)
	PublicRandStream(ctx context.Context, in *PublicRandRequest, opts ...grpc.CallOption) (Public_PublicRandStreamClient, error)
	// PrivateRand is the method that returns the private randomness generated
	// by the drand node only.
	PrivateRand(ctx context.Context, in *PrivateRandRequest, opts ...grpc.CallOption) (*PrivateRandResponse, error)
	// ChainInfo returns the information related to the chain this node
	// participates to
	ChainInfo(ctx context.Context, in *ChainInfoRequest, opts ...grpc.CallOption) (*ChainInfoPacket, error)
	// Home is a simple endpoint
	Home(ctx context.Context, in *HomeRequest, opts ...grpc.CallOption) (*HomeResponse, error)
}

type publicClient struct {
	cc grpc.ClientConnInterface
}

func NewPublicClient(cc grpc.ClientConnInterface) PublicClient {
	return &publicClient{cc}
}

func (c *publicClient) PublicRand(ctx context.Context, in *PublicRandRequest, opts ...grpc.CallOption) (*PublicRandResponse, error) {
	out := new(PublicRandResponse)
	err := c.cc.Invoke(ctx, "/drand.Public/PublicRand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicClient) PublicRandStream(ctx context.Context, in *PublicRandRequest, opts ...grpc.CallOption) (Public_PublicRandStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Public_serviceDesc.Streams[0], "/drand.Public/PublicRandStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &publicPublicRandStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Public_PublicRandStreamClient interface {
	Recv() (*PublicRandResponse, error)
	grpc.ClientStream
}

type publicPublicRandStreamClient struct {
	grpc.ClientStream
}

func (x *publicPublicRandStreamClient) Recv() (*PublicRandResponse, error) {
	m := new(PublicRandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *publicClient) PrivateRand(ctx context.Context, in *PrivateRandRequest, opts ...grpc.CallOption) (*PrivateRandResponse, error) {
	out := new(PrivateRandResponse)
	err := c.cc.Invoke(ctx, "/drand.Public/PrivateRand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicClient) ChainInfo(ctx context.Context, in *ChainInfoRequest, opts ...grpc.CallOption) (*ChainInfoPacket, error) {
	out := new(ChainInfoPacket)
	err := c.cc.Invoke(ctx, "/drand.Public/ChainInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publicClient) Home(ctx context.Context, in *HomeRequest, opts ...grpc.CallOption) (*HomeResponse, error) {
	out := new(HomeResponse)
	err := c.cc.Invoke(ctx, "/drand.Public/Home", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublicServer is the server API for Public service.
type PublicServer interface {
	// PublicRand is the method that returns the publicly verifiable randomness
	// generated by the drand network.
	PublicRand(context.Context, *PublicRandRequest) (*PublicRandResponse, error)
	PublicRandStream(*PublicRandRequest, Public_PublicRandStreamServer) error
	// PrivateRand is the method that returns the private randomness generated
	// by the drand node only.
	PrivateRand(context.Context, *PrivateRandRequest) (*PrivateRandResponse, error)
	// ChainInfo returns the information related to the chain this node
	// participates to
	ChainInfo(context.Context, *ChainInfoRequest) (*ChainInfoPacket, error)
	// Home is a simple endpoint
	Home(context.Context, *HomeRequest) (*HomeResponse, error)
}

// UnimplementedPublicServer can be embedded to have forward compatible implementations.
type UnimplementedPublicServer struct {
}

func (*UnimplementedPublicServer) PublicRand(ctx context.Context, req *PublicRandRequest) (*PublicRandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublicRand not implemented")
}
func (*UnimplementedPublicServer) PublicRandStream(req *PublicRandRequest, srv Public_PublicRandStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PublicRandStream not implemented")
}
func (*UnimplementedPublicServer) PrivateRand(ctx context.Context, req *PrivateRandRequest) (*PrivateRandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrivateRand not implemented")
}
func (*UnimplementedPublicServer) ChainInfo(ctx context.Context, req *ChainInfoRequest) (*ChainInfoPacket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChainInfo not implemented")
}
func (*UnimplementedPublicServer) Home(ctx context.Context, req *HomeRequest) (*HomeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Home not implemented")
}

func RegisterPublicServer(s *grpc.Server, srv PublicServer) {
	s.RegisterService(&_Public_serviceDesc, srv)
}

func _Public_PublicRand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicRandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicServer).PublicRand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Public/PublicRand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicServer).PublicRand(ctx, req.(*PublicRandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Public_PublicRandStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PublicRandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PublicServer).PublicRandStream(m, &publicPublicRandStreamServer{stream})
}

type Public_PublicRandStreamServer interface {
	Send(*PublicRandResponse) error
	grpc.ServerStream
}

type publicPublicRandStreamServer struct {
	grpc.ServerStream
}

func (x *publicPublicRandStreamServer) Send(m *PublicRandResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Public_PrivateRand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrivateRandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicServer).PrivateRand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Public/PrivateRand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicServer).PrivateRand(ctx, req.(*PrivateRandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Public_ChainInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChainInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicServer).ChainInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Public/ChainInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicServer).ChainInfo(ctx, req.(*ChainInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Public_Home_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HomeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublicServer).Home(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Public/Home",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublicServer).Home(ctx, req.(*HomeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Public_serviceDesc = grpc.ServiceDesc{
	ServiceName: "drand.Public",
	HandlerType: (*PublicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublicRand",
			Handler:    _Public_PublicRand_Handler,
		},
		{
			MethodName: "PrivateRand",
			Handler:    _Public_PrivateRand_Handler,
		},
		{
			MethodName: "ChainInfo",
			Handler:    _Public_ChainInfo_Handler,
		},
		{
			MethodName: "Home",
			Handler:    _Public_Home_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PublicRandStream",
			Handler:       _Public_PublicRandStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "drand/api.proto",
}
