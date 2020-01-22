// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crypto/dkg/dkg.proto

package dkg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import vss "github.com/drand/drand/protobuf/crypto/vss"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Packet is a wrapper around the three different types of DKG messages
type Packet struct {
	Deal                 *Deal          `protobuf:"bytes,1,opt,name=deal,proto3" json:"deal,omitempty"`
	Response             *Response      `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
	Justification        *Justification `protobuf:"bytes,3,opt,name=justification,proto3" json:"justification,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_dkg_4ed6c9b8ccb6a245, []int{0}
}
func (m *Packet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet.Unmarshal(m, b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
}
func (dst *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(dst, src)
}
func (m *Packet) XXX_Size() int {
	return xxx_messageInfo_Packet.Size(m)
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

func (m *Packet) GetDeal() *Deal {
	if m != nil {
		return m.Deal
	}
	return nil
}

func (m *Packet) GetResponse() *Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *Packet) GetJustification() *Justification {
	if m != nil {
		return m.Justification
	}
	return nil
}

// Deal contains a share for a participant.
type Deal struct {
	// index of the dealer, the issuer of the share
	Index uint32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// encrypted version of the deal
	Deal *vss.EncryptedDeal `protobuf:"bytes,2,opt,name=deal,proto3" json:"deal,omitempty"`
	// signature of the whole deal
	// NOTE: this is almost duplicated data, since the vss deal already includes
	// a signature. However it does not include the index of the dealer that
	// issue this deal, so another one is required. Best would be to merge vss
	// and dkg so we could use only one field of signature. For future work...
	// :)
	Signature            []byte   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Deal) Reset()         { *m = Deal{} }
func (m *Deal) String() string { return proto.CompactTextString(m) }
func (*Deal) ProtoMessage()    {}
func (*Deal) Descriptor() ([]byte, []int) {
	return fileDescriptor_dkg_4ed6c9b8ccb6a245, []int{1}
}
func (m *Deal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Deal.Unmarshal(m, b)
}
func (m *Deal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Deal.Marshal(b, m, deterministic)
}
func (dst *Deal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deal.Merge(dst, src)
}
func (m *Deal) XXX_Size() int {
	return xxx_messageInfo_Deal.Size(m)
}
func (m *Deal) XXX_DiscardUnknown() {
	xxx_messageInfo_Deal.DiscardUnknown(m)
}

var xxx_messageInfo_Deal proto.InternalMessageInfo

func (m *Deal) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Deal) GetDeal() *vss.EncryptedDeal {
	if m != nil {
		return m.Deal
	}
	return nil
}

func (m *Deal) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// Response holds the response that a participant broadcast after having
// received a deal.
type Response struct {
	// index of the dealer for which this response is for
	Index uint32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// response from the participant which received a deal
	Response             *vss.Response `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_dkg_4ed6c9b8ccb6a245, []int{2}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Response) GetResponse() *vss.Response {
	if m != nil {
		return m.Response
	}
	return nil
}

// Justification holds the justification from a dealer after a participant
// issued a complaint response because of a supposedly invalid deal.
type Justification struct {
	// index of the dealer who is issuing this justification
	Index uint32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	// justification from the dealer
	Justification        *vss.Justification `protobuf:"bytes,2,opt,name=justification,proto3" json:"justification,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Justification) Reset()         { *m = Justification{} }
func (m *Justification) String() string { return proto.CompactTextString(m) }
func (*Justification) ProtoMessage()    {}
func (*Justification) Descriptor() ([]byte, []int) {
	return fileDescriptor_dkg_4ed6c9b8ccb6a245, []int{3}
}
func (m *Justification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Justification.Unmarshal(m, b)
}
func (m *Justification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Justification.Marshal(b, m, deterministic)
}
func (dst *Justification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Justification.Merge(dst, src)
}
func (m *Justification) XXX_Size() int {
	return xxx_messageInfo_Justification.Size(m)
}
func (m *Justification) XXX_DiscardUnknown() {
	xxx_messageInfo_Justification.DiscardUnknown(m)
}

var xxx_messageInfo_Justification proto.InternalMessageInfo

func (m *Justification) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Justification) GetJustification() *vss.Justification {
	if m != nil {
		return m.Justification
	}
	return nil
}

func init() {
	proto.RegisterType((*Packet)(nil), "dkg.Packet")
	proto.RegisterType((*Deal)(nil), "dkg.Deal")
	proto.RegisterType((*Response)(nil), "dkg.Response")
	proto.RegisterType((*Justification)(nil), "dkg.Justification")
}

func init() { proto.RegisterFile("github.com/drand/drand/protobuf/crypto/dkg/dkg.proto", fileDescriptor_dkg_4ed6c9b8ccb6a245) }

var fileDescriptor_dkg_4ed6c9b8ccb6a245 = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x49, 0x2e, 0xaa, 0x2c,
	0x28, 0xc9, 0xd7, 0x4f, 0xc9, 0x4e, 0x07, 0x61, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0xe6,
	0x94, 0xec, 0x74, 0x29, 0x98, 0x54, 0x59, 0x71, 0x31, 0x08, 0x43, 0xa4, 0x94, 0x7a, 0x18, 0xb9,
	0xd8, 0x02, 0x12, 0x93, 0xb3, 0x53, 0x4b, 0x84, 0x64, 0xb9, 0x58, 0x52, 0x52, 0x13, 0x73, 0x24,
	0x18, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0x38, 0xf5, 0x40, 0xfa, 0x5d, 0x52, 0x13, 0x73, 0x82, 0xc0,
	0xc2, 0x42, 0x9a, 0x5c, 0x1c, 0x45, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x12, 0x4c, 0x60,
	0x25, 0xbc, 0x60, 0x25, 0x41, 0x50, 0xc1, 0x20, 0xb8, 0xb4, 0x90, 0x05, 0x17, 0x6f, 0x56, 0x69,
	0x71, 0x49, 0x66, 0x5a, 0x66, 0x72, 0x62, 0x49, 0x66, 0x7e, 0x9e, 0x04, 0x33, 0x58, 0xbd, 0x10,
	0x58, 0xbd, 0x17, 0xb2, 0x4c, 0x10, 0xaa, 0x42, 0xa5, 0x24, 0x2e, 0x16, 0x90, 0x95, 0x42, 0x22,
	0x5c, 0xac, 0x99, 0x79, 0x29, 0xa9, 0x15, 0x60, 0xc7, 0xf0, 0x06, 0x41, 0x38, 0x42, 0x6a, 0x50,
	0x17, 0x32, 0x41, 0x8d, 0x03, 0x79, 0xc3, 0x35, 0x0f, 0xec, 0xaf, 0xd4, 0x14, 0x24, 0xa7, 0xca,
	0x70, 0x71, 0x16, 0x67, 0xa6, 0xe7, 0x25, 0x96, 0x94, 0x16, 0xa5, 0x82, 0xed, 0xe6, 0x09, 0x42,
	0x08, 0x28, 0x79, 0x73, 0x71, 0xc0, 0xdc, 0x8c, 0xc3, 0x1e, 0x6c, 0x5e, 0x05, 0xd9, 0x85, 0xe9,
	0x55, 0xa5, 0x78, 0x2e, 0x5e, 0x14, 0x0f, 0xe1, 0x30, 0x11, 0x23, 0x44, 0x90, 0xbd, 0x80, 0x2f,
	0x44, 0x9c, 0x58, 0xa3, 0x40, 0xb1, 0x97, 0xc4, 0x06, 0x8e, 0x2e, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x01, 0x82, 0x40, 0x8e, 0xe1, 0x01, 0x00, 0x00,
}
