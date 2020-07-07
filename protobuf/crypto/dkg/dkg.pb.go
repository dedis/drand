// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crypto/dkg/dkg.proto

package dkg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Packet is a wrapper around the three different types of DKG messages
type Packet struct {
	// Types that are valid to be assigned to Bundle:
	//	*Packet_Deal
	//	*Packet_Response
	//	*Packet_Justification
	Bundle               isPacket_Bundle `protobuf_oneof:"Bundle"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Packet) Reset()         { *m = Packet{} }
func (m *Packet) String() string { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()    {}
func (*Packet) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{0}
}

func (m *Packet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Packet.Unmarshal(m, b)
}
func (m *Packet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Packet.Marshal(b, m, deterministic)
}
func (m *Packet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Packet.Merge(m, src)
}
func (m *Packet) XXX_Size() int {
	return xxx_messageInfo_Packet.Size(m)
}
func (m *Packet) XXX_DiscardUnknown() {
	xxx_messageInfo_Packet.DiscardUnknown(m)
}

var xxx_messageInfo_Packet proto.InternalMessageInfo

type isPacket_Bundle interface {
	isPacket_Bundle()
}

type Packet_Deal struct {
	Deal *DealBundle `protobuf:"bytes,1,opt,name=deal,proto3,oneof"`
}

type Packet_Response struct {
	Response *ResponseBundle `protobuf:"bytes,2,opt,name=response,proto3,oneof"`
}

type Packet_Justification struct {
	Justification *JustificationBundle `protobuf:"bytes,3,opt,name=justification,proto3,oneof"`
}

func (*Packet_Deal) isPacket_Bundle() {}

func (*Packet_Response) isPacket_Bundle() {}

func (*Packet_Justification) isPacket_Bundle() {}

func (m *Packet) GetBundle() isPacket_Bundle {
	if m != nil {
		return m.Bundle
	}
	return nil
}

func (m *Packet) GetDeal() *DealBundle {
	if x, ok := m.GetBundle().(*Packet_Deal); ok {
		return x.Deal
	}
	return nil
}

func (m *Packet) GetResponse() *ResponseBundle {
	if x, ok := m.GetBundle().(*Packet_Response); ok {
		return x.Response
	}
	return nil
}

func (m *Packet) GetJustification() *JustificationBundle {
	if x, ok := m.GetBundle().(*Packet_Justification); ok {
		return x.Justification
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Packet) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Packet_Deal)(nil),
		(*Packet_Response)(nil),
		(*Packet_Justification)(nil),
	}
}

// DealBundle is a packet issued by a dealer that contains each individual
// deals, as well as the coefficients of the public polynomial he used.
type DealBundle struct {
	// Index of the dealer that issues these deals
	DealerIndex uint32 `protobuf:"varint,1,opt,name=dealer_index,json=dealerIndex,proto3" json:"dealer_index,omitempty"`
	// Coefficients of the public polynomial that is created from the
	// private polynomial from which the shares are derived.
	Commits [][]byte `protobuf:"bytes,2,rep,name=commits,proto3" json:"commits,omitempty"`
	// list of deals for each individual share holders.
	Deals []*Deal `protobuf:"bytes,3,rep,name=deals,proto3" json:"deals,omitempty"`
	// session identifier of the protocol run
	SessionId []byte `protobuf:"bytes,4,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// signature over the hash of the deal
	Signature            []byte   `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DealBundle) Reset()         { *m = DealBundle{} }
func (m *DealBundle) String() string { return proto.CompactTextString(m) }
func (*DealBundle) ProtoMessage()    {}
func (*DealBundle) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{1}
}

func (m *DealBundle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DealBundle.Unmarshal(m, b)
}
func (m *DealBundle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DealBundle.Marshal(b, m, deterministic)
}
func (m *DealBundle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DealBundle.Merge(m, src)
}
func (m *DealBundle) XXX_Size() int {
	return xxx_messageInfo_DealBundle.Size(m)
}
func (m *DealBundle) XXX_DiscardUnknown() {
	xxx_messageInfo_DealBundle.DiscardUnknown(m)
}

var xxx_messageInfo_DealBundle proto.InternalMessageInfo

func (m *DealBundle) GetDealerIndex() uint32 {
	if m != nil {
		return m.DealerIndex
	}
	return 0
}

func (m *DealBundle) GetCommits() [][]byte {
	if m != nil {
		return m.Commits
	}
	return nil
}

func (m *DealBundle) GetDeals() []*Deal {
	if m != nil {
		return m.Deals
	}
	return nil
}

func (m *DealBundle) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

func (m *DealBundle) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// Deal contains a share for a participant.
type Deal struct {
	ShareIndex uint32 `protobuf:"varint,1,opt,name=share_index,json=shareIndex,proto3" json:"share_index,omitempty"`
	// encryption of the share using ECIES
	EncryptedShare       []byte   `protobuf:"bytes,2,opt,name=encrypted_share,json=encryptedShare,proto3" json:"encrypted_share,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Deal) Reset()         { *m = Deal{} }
func (m *Deal) String() string { return proto.CompactTextString(m) }
func (*Deal) ProtoMessage()    {}
func (*Deal) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{2}
}

func (m *Deal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Deal.Unmarshal(m, b)
}
func (m *Deal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Deal.Marshal(b, m, deterministic)
}
func (m *Deal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deal.Merge(m, src)
}
func (m *Deal) XXX_Size() int {
	return xxx_messageInfo_Deal.Size(m)
}
func (m *Deal) XXX_DiscardUnknown() {
	xxx_messageInfo_Deal.DiscardUnknown(m)
}

var xxx_messageInfo_Deal proto.InternalMessageInfo

func (m *Deal) GetShareIndex() uint32 {
	if m != nil {
		return m.ShareIndex
	}
	return 0
}

func (m *Deal) GetEncryptedShare() []byte {
	if m != nil {
		return m.EncryptedShare
	}
	return nil
}

// ResponseBundle is a packet issued by a share holder that contains all the
// responses (complaint and/or success) to broadcast.
type ResponseBundle struct {
	ShareIndex uint32      `protobuf:"varint,1,opt,name=share_index,json=shareIndex,proto3" json:"share_index,omitempty"`
	Responses  []*Response `protobuf:"bytes,2,rep,name=responses,proto3" json:"responses,omitempty"`
	// session identifier of the protocol run
	SessionId []byte `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// signature over the hash of the response
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseBundle) Reset()         { *m = ResponseBundle{} }
func (m *ResponseBundle) String() string { return proto.CompactTextString(m) }
func (*ResponseBundle) ProtoMessage()    {}
func (*ResponseBundle) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{3}
}

func (m *ResponseBundle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseBundle.Unmarshal(m, b)
}
func (m *ResponseBundle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseBundle.Marshal(b, m, deterministic)
}
func (m *ResponseBundle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseBundle.Merge(m, src)
}
func (m *ResponseBundle) XXX_Size() int {
	return xxx_messageInfo_ResponseBundle.Size(m)
}
func (m *ResponseBundle) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseBundle.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseBundle proto.InternalMessageInfo

func (m *ResponseBundle) GetShareIndex() uint32 {
	if m != nil {
		return m.ShareIndex
	}
	return 0
}

func (m *ResponseBundle) GetResponses() []*Response {
	if m != nil {
		return m.Responses
	}
	return nil
}

func (m *ResponseBundle) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

func (m *ResponseBundle) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// Response holds the response that a participant broadcast after having
// received a deal.
type Response struct {
	// index of the dealer for which this response is for
	DealerIndex uint32 `protobuf:"varint,1,opt,name=dealer_index,json=dealerIndex,proto3" json:"dealer_index,omitempty"`
	// Status represents a complaint if set to false, a success if set to
	// true.
	Status               bool     `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{4}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetDealerIndex() uint32 {
	if m != nil {
		return m.DealerIndex
	}
	return 0
}

func (m *Response) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

// JustificationBundle is a packet that holds all justifications a dealer must
// produce
type JustificationBundle struct {
	DealerIndex    uint32           `protobuf:"varint,1,opt,name=dealer_index,json=dealerIndex,proto3" json:"dealer_index,omitempty"`
	Justifications []*Justification `protobuf:"bytes,2,rep,name=justifications,proto3" json:"justifications,omitempty"`
	// session identifier of the protocol run
	SessionId []byte `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	// signature over the hash of the justification
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JustificationBundle) Reset()         { *m = JustificationBundle{} }
func (m *JustificationBundle) String() string { return proto.CompactTextString(m) }
func (*JustificationBundle) ProtoMessage()    {}
func (*JustificationBundle) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{5}
}

func (m *JustificationBundle) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JustificationBundle.Unmarshal(m, b)
}
func (m *JustificationBundle) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JustificationBundle.Marshal(b, m, deterministic)
}
func (m *JustificationBundle) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JustificationBundle.Merge(m, src)
}
func (m *JustificationBundle) XXX_Size() int {
	return xxx_messageInfo_JustificationBundle.Size(m)
}
func (m *JustificationBundle) XXX_DiscardUnknown() {
	xxx_messageInfo_JustificationBundle.DiscardUnknown(m)
}

var xxx_messageInfo_JustificationBundle proto.InternalMessageInfo

func (m *JustificationBundle) GetDealerIndex() uint32 {
	if m != nil {
		return m.DealerIndex
	}
	return 0
}

func (m *JustificationBundle) GetJustifications() []*Justification {
	if m != nil {
		return m.Justifications
	}
	return nil
}

func (m *JustificationBundle) GetSessionId() []byte {
	if m != nil {
		return m.SessionId
	}
	return nil
}

func (m *JustificationBundle) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// Justification holds the justification from a dealer after a participant
// issued a complaint response because of a supposedly invalid deal.
type Justification struct {
	// represents for who share holder this justification is
	ShareIndex uint32 `protobuf:"varint,1,opt,name=share_index,json=shareIndex,proto3" json:"share_index,omitempty"`
	// plaintext share so everyone can see it correct
	Share                []byte   `protobuf:"bytes,2,opt,name=share,proto3" json:"share,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Justification) Reset()         { *m = Justification{} }
func (m *Justification) String() string { return proto.CompactTextString(m) }
func (*Justification) ProtoMessage()    {}
func (*Justification) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cd2862d3a18e91b, []int{6}
}

func (m *Justification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Justification.Unmarshal(m, b)
}
func (m *Justification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Justification.Marshal(b, m, deterministic)
}
func (m *Justification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Justification.Merge(m, src)
}
func (m *Justification) XXX_Size() int {
	return xxx_messageInfo_Justification.Size(m)
}
func (m *Justification) XXX_DiscardUnknown() {
	xxx_messageInfo_Justification.DiscardUnknown(m)
}

var xxx_messageInfo_Justification proto.InternalMessageInfo

func (m *Justification) GetShareIndex() uint32 {
	if m != nil {
		return m.ShareIndex
	}
	return 0
}

func (m *Justification) GetShare() []byte {
	if m != nil {
		return m.Share
	}
	return nil
}

func init() {
	proto.RegisterType((*Packet)(nil), "dkg.Packet")
	proto.RegisterType((*DealBundle)(nil), "dkg.DealBundle")
	proto.RegisterType((*Deal)(nil), "dkg.Deal")
	proto.RegisterType((*ResponseBundle)(nil), "dkg.ResponseBundle")
	proto.RegisterType((*Response)(nil), "dkg.Response")
	proto.RegisterType((*JustificationBundle)(nil), "dkg.JustificationBundle")
	proto.RegisterType((*Justification)(nil), "dkg.Justification")
}

func init() {
	proto.RegisterFile("crypto/dkg/dkg.proto", fileDescriptor_2cd2862d3a18e91b)
}

var fileDescriptor_2cd2862d3a18e91b = []byte{
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcf, 0x8e, 0xd3, 0x30,
	0x10, 0xc6, 0xc9, 0xa6, 0x0d, 0xc9, 0x24, 0xe9, 0x4a, 0xde, 0x15, 0xf2, 0x01, 0xb4, 0x21, 0x12,
	0x22, 0x02, 0xd4, 0x8a, 0x72, 0xe3, 0x84, 0x2a, 0x40, 0x2c, 0xa7, 0x95, 0xb9, 0x71, 0xa9, 0xdc,
	0xd8, 0x9b, 0x35, 0x6d, 0xed, 0x2a, 0x76, 0x24, 0x78, 0x16, 0x1e, 0x80, 0x1b, 0x12, 0x6f, 0x88,
	0xe2, 0xfc, 0xe9, 0xa6, 0x20, 0xb6, 0x12, 0x87, 0x44, 0xf2, 0x6f, 0xbe, 0x19, 0xcd, 0xe7, 0x19,
	0xc3, 0x79, 0x5e, 0x7e, 0xdb, 0x19, 0x35, 0x63, 0xeb, 0xa2, 0xfe, 0xa6, 0xbb, 0x52, 0x19, 0x85,
	0x5c, 0xb6, 0x2e, 0xd2, 0x9f, 0x0e, 0x78, 0x57, 0x34, 0x5f, 0x73, 0x83, 0x9e, 0xc0, 0x88, 0x71,
	0xba, 0xc1, 0x4e, 0xe2, 0x64, 0xe1, 0xfc, 0x74, 0x5a, 0x2b, 0xdf, 0x72, 0xba, 0x59, 0x54, 0x92,
	0x6d, 0xf8, 0x87, 0x7b, 0xc4, 0x86, 0xd1, 0x4b, 0xf0, 0x4b, 0xae, 0x77, 0x4a, 0x6a, 0x8e, 0x4f,
	0xac, 0xf4, 0xcc, 0x4a, 0x49, 0x0b, 0x7b, 0x79, 0x2f, 0x43, 0x6f, 0x20, 0xfe, 0x52, 0x69, 0x23,
	0xae, 0x45, 0x4e, 0x8d, 0x50, 0x12, 0xbb, 0x36, 0x0f, 0xdb, 0xbc, 0x8f, 0xb7, 0x23, 0x7d, 0xf2,
	0x30, 0x61, 0xe1, 0x83, 0xd7, 0x84, 0xd2, 0x1f, 0x0e, 0xc0, 0xbe, 0x2b, 0xf4, 0x18, 0xa2, 0xba,
	0x2b, 0x5e, 0x2e, 0x85, 0x64, 0xfc, 0xab, 0x6d, 0x3e, 0x26, 0x61, 0xc3, 0x2e, 0x6b, 0x84, 0x30,
	0xdc, 0xcf, 0xd5, 0x76, 0x2b, 0x8c, 0xc6, 0x27, 0x89, 0x9b, 0x45, 0xa4, 0x3b, 0xa2, 0x0b, 0x18,
	0xd7, 0x42, 0x8d, 0xdd, 0xc4, 0xcd, 0xc2, 0x79, 0xd0, 0x5b, 0x26, 0x0d, 0x47, 0x8f, 0x00, 0x34,
	0xd7, 0x5a, 0x28, 0xb9, 0x14, 0x0c, 0x8f, 0x12, 0x27, 0x8b, 0x48, 0xd0, 0x92, 0x4b, 0x86, 0x1e,
	0x42, 0xa0, 0x45, 0x21, 0xa9, 0xa9, 0x4a, 0x8e, 0xc7, 0x6d, 0xb4, 0x03, 0xe9, 0x15, 0x8c, 0xea,
	0x5a, 0xe8, 0x02, 0x42, 0x7d, 0x43, 0x4b, 0x3e, 0xe8, 0x10, 0x2c, 0x6a, 0x1a, 0x7c, 0x0a, 0xa7,
	0x5c, 0xda, 0x11, 0x71, 0xb6, 0xb4, 0xdc, 0x5e, 0x6c, 0x44, 0x26, 0x3d, 0xfe, 0x54, 0xd3, 0xf4,
	0xbb, 0x03, 0x93, 0xe1, 0x35, 0xdf, 0x5d, 0xfc, 0x39, 0x04, 0xdd, 0x1c, 0x1a, 0xff, 0xe1, 0x3c,
	0x1e, 0xcc, 0x8b, 0xec, 0xe3, 0x07, 0x7e, 0xdd, 0x7f, 0xfa, 0x1d, 0x1d, 0xfa, 0x7d, 0x07, 0x7e,
	0x57, 0xf3, 0x98, 0xb1, 0x3c, 0x00, 0x4f, 0x1b, 0x6a, 0x2a, 0x6d, 0xcd, 0xfa, 0xa4, 0x3d, 0xa5,
	0xbf, 0x1c, 0x38, 0xfb, 0xcb, 0x4e, 0x1c, 0x53, 0xf2, 0x35, 0x4c, 0x06, 0x6b, 0xd3, 0x19, 0x46,
	0x7f, 0x2e, 0x1a, 0x39, 0x50, 0xfe, 0x9f, 0xf5, 0xf7, 0x10, 0x0f, 0xaa, 0xdf, 0x3d, 0x96, 0x73,
	0x18, 0xdf, 0x9e, 0x74, 0x73, 0x58, 0xbc, 0xf8, 0xfc, 0xac, 0x10, 0xe6, 0xa6, 0x5a, 0x4d, 0x73,
	0xb5, 0x9d, 0xb1, 0x92, 0x4a, 0xd6, 0xfe, 0xed, 0x93, 0x5d, 0x55, 0xd7, 0xb3, 0xfd, 0x53, 0x5e,
	0x79, 0x16, 0xbe, 0xfa, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xde, 0x23, 0x27, 0x2a, 0xdf, 0x03, 0x00,
	0x00,
}
