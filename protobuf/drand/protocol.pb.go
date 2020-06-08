// Code generated by protoc-gen-go. DO NOT EDIT.
// source: drand/protocol.proto

package drand

import (
	context "context"
	dkg "github.com/drand/drand/protobuf/crypto/dkg"
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

type IdentityRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentityRequest) Reset()         { *m = IdentityRequest{} }
func (m *IdentityRequest) String() string { return proto.CompactTextString(m) }
func (*IdentityRequest) ProtoMessage()    {}
func (*IdentityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{0}
}

func (m *IdentityRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdentityRequest.Unmarshal(m, b)
}
func (m *IdentityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdentityRequest.Marshal(b, m, deterministic)
}
func (m *IdentityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityRequest.Merge(m, src)
}
func (m *IdentityRequest) XXX_Size() int {
	return xxx_messageInfo_IdentityRequest.Size(m)
}
func (m *IdentityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityRequest proto.InternalMessageInfo

// SignalDKGPacket is the packet nodes send to a coordinator that collects all
// keys and setups the group and sends them back to the nodes such that they can
// start the DKG automatically.
type SignalDKGPacket struct {
	Node        *Identity `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	SecretProof string    `protobuf:"bytes,2,opt,name=secret_proof,json=secretProof,proto3" json:"secret_proof,omitempty"`
	// In resharing cases, previous_group_hash is the hash of the previous group.
	// It is to make sure the nodes build on top of the correct previous group.
	PreviousGroupHash    []byte   `protobuf:"bytes,3,opt,name=previous_group_hash,json=previousGroupHash,proto3" json:"previous_group_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignalDKGPacket) Reset()         { *m = SignalDKGPacket{} }
func (m *SignalDKGPacket) String() string { return proto.CompactTextString(m) }
func (*SignalDKGPacket) ProtoMessage()    {}
func (*SignalDKGPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{1}
}

func (m *SignalDKGPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignalDKGPacket.Unmarshal(m, b)
}
func (m *SignalDKGPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignalDKGPacket.Marshal(b, m, deterministic)
}
func (m *SignalDKGPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalDKGPacket.Merge(m, src)
}
func (m *SignalDKGPacket) XXX_Size() int {
	return xxx_messageInfo_SignalDKGPacket.Size(m)
}
func (m *SignalDKGPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalDKGPacket.DiscardUnknown(m)
}

var xxx_messageInfo_SignalDKGPacket proto.InternalMessageInfo

func (m *SignalDKGPacket) GetNode() *Identity {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *SignalDKGPacket) GetSecretProof() string {
	if m != nil {
		return m.SecretProof
	}
	return ""
}

func (m *SignalDKGPacket) GetPreviousGroupHash() []byte {
	if m != nil {
		return m.PreviousGroupHash
	}
	return nil
}

// PushDKGInfor is the packet the coordinator sends that contains the group over
// which to run the DKG on, the secret proof (to prove it's he's part of the
// expected group, and it's not a random packet) and as well the time at which
// every node should start the DKG.
type DKGInfoPacket struct {
	NewGroup    *GroupPacket `protobuf:"bytes,1,opt,name=new_group,json=newGroup,proto3" json:"new_group,omitempty"`
	SecretProof string       `protobuf:"bytes,2,opt,name=secret_proof,json=secretProof,proto3" json:"secret_proof,omitempty"`
	// timeout in seconds
	DkgTimeout uint32 `protobuf:"varint,3,opt,name=dkg_timeout,json=dkgTimeout,proto3" json:"dkg_timeout,omitempty"`
	// signature from the coordinator to prove he is the one sending that group
	// file.
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DKGInfoPacket) Reset()         { *m = DKGInfoPacket{} }
func (m *DKGInfoPacket) String() string { return proto.CompactTextString(m) }
func (*DKGInfoPacket) ProtoMessage()    {}
func (*DKGInfoPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{2}
}

func (m *DKGInfoPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DKGInfoPacket.Unmarshal(m, b)
}
func (m *DKGInfoPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DKGInfoPacket.Marshal(b, m, deterministic)
}
func (m *DKGInfoPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DKGInfoPacket.Merge(m, src)
}
func (m *DKGInfoPacket) XXX_Size() int {
	return xxx_messageInfo_DKGInfoPacket.Size(m)
}
func (m *DKGInfoPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_DKGInfoPacket.DiscardUnknown(m)
}

var xxx_messageInfo_DKGInfoPacket proto.InternalMessageInfo

func (m *DKGInfoPacket) GetNewGroup() *GroupPacket {
	if m != nil {
		return m.NewGroup
	}
	return nil
}

func (m *DKGInfoPacket) GetSecretProof() string {
	if m != nil {
		return m.SecretProof
	}
	return ""
}

func (m *DKGInfoPacket) GetDkgTimeout() uint32 {
	if m != nil {
		return m.DkgTimeout
	}
	return 0
}

func (m *DKGInfoPacket) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type PartialBeaconPacket struct {
	// Round is the round for which the beacon will be created from the partial
	// signatures
	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	// signature of the previous round - could be removed at some point but now
	// is used to verify the signature even before accessing the store
	PreviousSig []byte `protobuf:"bytes,2,opt,name=previous_sig,json=previousSig,proto3" json:"previous_sig,omitempty"`
	// partial signature - a threshold of them needs to be aggregated to produce
	// the final beacon at the given round.
	PartialSig           []byte   `protobuf:"bytes,3,opt,name=partial_sig,json=partialSig,proto3" json:"partial_sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PartialBeaconPacket) Reset()         { *m = PartialBeaconPacket{} }
func (m *PartialBeaconPacket) String() string { return proto.CompactTextString(m) }
func (*PartialBeaconPacket) ProtoMessage()    {}
func (*PartialBeaconPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{3}
}

func (m *PartialBeaconPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PartialBeaconPacket.Unmarshal(m, b)
}
func (m *PartialBeaconPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PartialBeaconPacket.Marshal(b, m, deterministic)
}
func (m *PartialBeaconPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PartialBeaconPacket.Merge(m, src)
}
func (m *PartialBeaconPacket) XXX_Size() int {
	return xxx_messageInfo_PartialBeaconPacket.Size(m)
}
func (m *PartialBeaconPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_PartialBeaconPacket.DiscardUnknown(m)
}

var xxx_messageInfo_PartialBeaconPacket proto.InternalMessageInfo

func (m *PartialBeaconPacket) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *PartialBeaconPacket) GetPreviousSig() []byte {
	if m != nil {
		return m.PreviousSig
	}
	return nil
}

func (m *PartialBeaconPacket) GetPartialSig() []byte {
	if m != nil {
		return m.PartialSig
	}
	return nil
}

type DKGPacket struct {
	Dkg                  *dkg.Packet `protobuf:"bytes,1,opt,name=dkg,proto3" json:"dkg,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *DKGPacket) Reset()         { *m = DKGPacket{} }
func (m *DKGPacket) String() string { return proto.CompactTextString(m) }
func (*DKGPacket) ProtoMessage()    {}
func (*DKGPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{4}
}

func (m *DKGPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DKGPacket.Unmarshal(m, b)
}
func (m *DKGPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DKGPacket.Marshal(b, m, deterministic)
}
func (m *DKGPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DKGPacket.Merge(m, src)
}
func (m *DKGPacket) XXX_Size() int {
	return xxx_messageInfo_DKGPacket.Size(m)
}
func (m *DKGPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_DKGPacket.DiscardUnknown(m)
}

var xxx_messageInfo_DKGPacket proto.InternalMessageInfo

func (m *DKGPacket) GetDkg() *dkg.Packet {
	if m != nil {
		return m.Dkg
	}
	return nil
}

// Reshare is a wrapper around a Setup packet for resharing operation that
// serves two purposes: - indicate to non-leader old nodes that they should
// generate and send their deals - indicate to which new group are we resharing.
// drand should keep a list of new ready-to-operate groups allowed.
type ResharePacket struct {
	Dkg                  *dkg.Packet `protobuf:"bytes,1,opt,name=dkg,proto3" json:"dkg,omitempty"`
	GroupHash            string      `protobuf:"bytes,2,opt,name=group_hash,json=groupHash,proto3" json:"group_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ResharePacket) Reset()         { *m = ResharePacket{} }
func (m *ResharePacket) String() string { return proto.CompactTextString(m) }
func (*ResharePacket) ProtoMessage()    {}
func (*ResharePacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{5}
}

func (m *ResharePacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResharePacket.Unmarshal(m, b)
}
func (m *ResharePacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResharePacket.Marshal(b, m, deterministic)
}
func (m *ResharePacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResharePacket.Merge(m, src)
}
func (m *ResharePacket) XXX_Size() int {
	return xxx_messageInfo_ResharePacket.Size(m)
}
func (m *ResharePacket) XXX_DiscardUnknown() {
	xxx_messageInfo_ResharePacket.DiscardUnknown(m)
}

var xxx_messageInfo_ResharePacket proto.InternalMessageInfo

func (m *ResharePacket) GetDkg() *dkg.Packet {
	if m != nil {
		return m.Dkg
	}
	return nil
}

func (m *ResharePacket) GetGroupHash() string {
	if m != nil {
		return m.GroupHash
	}
	return ""
}

// SyncRequest is from a node that needs to sync up with the current head of the
// chain
type SyncRequest struct {
	FromRound            uint64   `protobuf:"varint,1,opt,name=from_round,json=fromRound,proto3" json:"from_round,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncRequest) Reset()         { *m = SyncRequest{} }
func (m *SyncRequest) String() string { return proto.CompactTextString(m) }
func (*SyncRequest) ProtoMessage()    {}
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{6}
}

func (m *SyncRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncRequest.Unmarshal(m, b)
}
func (m *SyncRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncRequest.Marshal(b, m, deterministic)
}
func (m *SyncRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncRequest.Merge(m, src)
}
func (m *SyncRequest) XXX_Size() int {
	return xxx_messageInfo_SyncRequest.Size(m)
}
func (m *SyncRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncRequest proto.InternalMessageInfo

func (m *SyncRequest) GetFromRound() uint64 {
	if m != nil {
		return m.FromRound
	}
	return 0
}

type BeaconPacket struct {
	PreviousSig          []byte   `protobuf:"bytes,1,opt,name=previous_sig,json=previousSig,proto3" json:"previous_sig,omitempty"`
	Round                uint64   `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Signature            []byte   `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BeaconPacket) Reset()         { *m = BeaconPacket{} }
func (m *BeaconPacket) String() string { return proto.CompactTextString(m) }
func (*BeaconPacket) ProtoMessage()    {}
func (*BeaconPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_e344a98fea1e2f3a, []int{7}
}

func (m *BeaconPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BeaconPacket.Unmarshal(m, b)
}
func (m *BeaconPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BeaconPacket.Marshal(b, m, deterministic)
}
func (m *BeaconPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BeaconPacket.Merge(m, src)
}
func (m *BeaconPacket) XXX_Size() int {
	return xxx_messageInfo_BeaconPacket.Size(m)
}
func (m *BeaconPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_BeaconPacket.DiscardUnknown(m)
}

var xxx_messageInfo_BeaconPacket proto.InternalMessageInfo

func (m *BeaconPacket) GetPreviousSig() []byte {
	if m != nil {
		return m.PreviousSig
	}
	return nil
}

func (m *BeaconPacket) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *BeaconPacket) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*IdentityRequest)(nil), "drand.IdentityRequest")
	proto.RegisterType((*SignalDKGPacket)(nil), "drand.SignalDKGPacket")
	proto.RegisterType((*DKGInfoPacket)(nil), "drand.DKGInfoPacket")
	proto.RegisterType((*PartialBeaconPacket)(nil), "drand.PartialBeaconPacket")
	proto.RegisterType((*DKGPacket)(nil), "drand.DKGPacket")
	proto.RegisterType((*ResharePacket)(nil), "drand.ResharePacket")
	proto.RegisterType((*SyncRequest)(nil), "drand.SyncRequest")
	proto.RegisterType((*BeaconPacket)(nil), "drand.BeaconPacket")
}

func init() {
	proto.RegisterFile("drand/protocol.proto", fileDescriptor_e344a98fea1e2f3a)
}

var fileDescriptor_e344a98fea1e2f3a = []byte{
	// 556 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4f, 0x6b, 0x13, 0x41,
	0x1c, 0x65, 0x9b, 0x56, 0xb3, 0xbf, 0x4d, 0x88, 0x9d, 0x04, 0x09, 0x8b, 0xc5, 0xb8, 0x5e, 0x42,
	0x91, 0xa4, 0x56, 0x28, 0x08, 0x9e, 0x6a, 0x35, 0x96, 0x22, 0x84, 0x8d, 0x27, 0x2f, 0x61, 0xdc,
	0x9d, 0xcc, 0x0e, 0x49, 0x66, 0xd6, 0xd9, 0x59, 0x4b, 0xee, 0x5e, 0xfc, 0x1a, 0x7e, 0x52, 0x99,
	0x3f, 0x49, 0x36, 0x6b, 0x0f, 0x3d, 0x04, 0xb2, 0xef, 0xf7, 0x87, 0x37, 0xef, 0xbd, 0x19, 0xe8,
	0xa5, 0x12, 0xf3, 0x74, 0x9c, 0x4b, 0xa1, 0x44, 0x22, 0x56, 0x23, 0xf3, 0x07, 0x9d, 0x18, 0x34,
	0xec, 0x25, 0x72, 0x93, 0x2b, 0x31, 0x4e, 0x97, 0x54, 0xff, 0x6c, 0x31, 0x44, 0x76, 0x24, 0x11,
	0xeb, 0xb5, 0xe0, 0x16, 0x8b, 0x4e, 0xa1, 0x73, 0x9b, 0x12, 0xae, 0x98, 0xda, 0xc4, 0xe4, 0x67,
	0x49, 0x0a, 0x15, 0xfd, 0xf1, 0xa0, 0x33, 0x63, 0x94, 0xe3, 0xd5, 0xcd, 0xdd, 0x64, 0x8a, 0x93,
	0x25, 0x51, 0xe8, 0x35, 0x1c, 0x73, 0x91, 0x92, 0xbe, 0x37, 0xf0, 0x86, 0xc1, 0x65, 0x67, 0x64,
	0x36, 0x8d, 0x76, 0x93, 0xa6, 0x88, 0x5e, 0x41, 0xab, 0x20, 0x89, 0x24, 0x6a, 0x9e, 0x4b, 0x21,
	0x16, 0xfd, 0xa3, 0x81, 0x37, 0xf4, 0xe3, 0xc0, 0x62, 0x53, 0x0d, 0xa1, 0x11, 0x74, 0x73, 0x49,
	0x7e, 0x31, 0x51, 0x16, 0x73, 0x2a, 0x45, 0x99, 0xcf, 0x33, 0x5c, 0x64, 0xfd, 0xc6, 0xc0, 0x1b,
	0xb6, 0xe2, 0xd3, 0x6d, 0x69, 0xa2, 0x2b, 0x5f, 0x70, 0x91, 0x45, 0x7f, 0x3d, 0x68, 0xdf, 0xdc,
	0x4d, 0x6e, 0xf9, 0x42, 0x38, 0x26, 0x63, 0xf0, 0x39, 0xb9, 0xb7, 0xc3, 0x8e, 0x0e, 0x72, 0x74,
	0xcc, 0x98, 0x6d, 0x8b, 0x9b, 0x9c, 0xdc, 0x9b, 0xef, 0xc7, 0xb0, 0x7a, 0x09, 0x41, 0xba, 0xa4,
	0x73, 0xc5, 0xd6, 0x44, 0x94, 0xca, 0xb0, 0x69, 0xc7, 0x90, 0x2e, 0xe9, 0x37, 0x8b, 0xa0, 0x17,
	0xe0, 0x17, 0x5a, 0x11, 0x55, 0x4a, 0xd2, 0x3f, 0x36, 0x64, 0xf7, 0x40, 0x24, 0xa0, 0x3b, 0xc5,
	0x52, 0x31, 0xbc, 0xba, 0x26, 0x38, 0x11, 0xdc, 0x31, 0xed, 0xc1, 0x89, 0x14, 0x25, 0x4f, 0x0d,
	0xcb, 0xe3, 0xd8, 0x7e, 0x68, 0x3a, 0x3b, 0x05, 0x0a, 0x46, 0x0d, 0x9d, 0x56, 0x1c, 0x6c, 0xb1,
	0x19, 0xa3, 0x9a, 0x4e, 0x6e, 0xf7, 0x99, 0x0e, 0x2b, 0x0e, 0x38, 0x68, 0xc6, 0x68, 0x74, 0x0e,
	0xfe, 0xde, 0x9a, 0x33, 0x68, 0xa4, 0x4b, 0xea, 0xa4, 0x08, 0x46, 0xda, 0x6e, 0xa7, 0x81, 0xc6,
	0xa3, 0xaf, 0xd0, 0x8e, 0x49, 0x91, 0x61, 0x49, 0x1e, 0xd5, 0x8f, 0xce, 0x00, 0x2a, 0xc6, 0x58,
	0xb1, 0x7c, 0xba, 0x33, 0xe4, 0x0d, 0x04, 0xb3, 0x0d, 0x4f, 0x5c, 0x56, 0x74, 0xf7, 0x42, 0x8a,
	0xf5, 0xbc, 0x7a, 0x50, 0x5f, 0x23, 0xb1, 0x06, 0x22, 0x02, 0xad, 0x03, 0x49, 0xea, 0x87, 0xf7,
	0xfe, 0x3f, 0xfc, 0x4e, 0xb5, 0xa3, 0xaa, 0x6a, 0x07, 0x06, 0x34, 0x6a, 0x06, 0x5c, 0xfe, 0x6e,
	0x40, 0x73, 0xea, 0x2e, 0x02, 0xba, 0x82, 0x60, 0x42, 0xd4, 0x36, 0x9a, 0xe8, 0x79, 0x3d, 0xab,
	0x96, 0x79, 0x58, 0xcf, 0x30, 0xfa, 0x00, 0xbd, 0x4a, 0xea, 0xa5, 0x62, 0x09, 0xcb, 0x31, 0x57,
	0xbb, 0x05, 0xb5, 0x2b, 0x11, 0xb6, 0x1c, 0xfe, 0x69, 0x9d, 0xab, 0x0d, 0x7a, 0x0b, 0xc1, 0xb4,
	0x2c, 0x32, 0x97, 0x55, 0xd4, 0x73, 0xc5, 0x83, 0xec, 0xd6, 0x46, 0xce, 0xa1, 0xf9, 0x59, 0x12,
	0x33, 0x83, 0x9e, 0xed, 0xfb, 0x1f, 0xec, 0xbd, 0x00, 0x70, 0x2e, 0xea, 0xee, 0xed, 0xf6, 0x03,
	0x63, 0x6b, 0x13, 0xef, 0xa1, 0x7d, 0x10, 0x4a, 0x14, 0xba, 0xf2, 0x03, 0x51, 0xad, 0x8d, 0x5e,
	0x81, 0xaf, 0x3d, 0xfe, 0x98, 0x61, 0xc6, 0xd1, 0xf6, 0x72, 0x55, 0x5c, 0x0f, 0xbb, 0x0e, 0xab,
	0xee, 0xb8, 0xf0, 0xae, 0x9f, 0x7e, 0xb7, 0xcf, 0xcf, 0x8f, 0x27, 0xe6, 0x6d, 0x79, 0xf7, 0x2f,
	0x00, 0x00, 0xff, 0xff, 0x89, 0x95, 0xeb, 0x93, 0xa4, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProtocolClient is the client API for Protocol service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProtocolClient interface {
	// GetIdentity returns the identity of the drand node
	GetIdentity(ctx context.Context, in *IdentityRequest, opts ...grpc.CallOption) (*Identity, error)
	// SignalDKGParticipant is called by non-coordinators nodes that sends their
	// public keys and secret proof they have to the coordinator so that he can
	// create the group.
	SignalDKGParticipant(ctx context.Context, in *SignalDKGPacket, opts ...grpc.CallOption) (*Empty, error)
	// PushDKGInfo is called by the coordinator to push the group he created
	// from all received keys and as well other information such as the time of
	// starting the DKG.
	PushDKGInfo(ctx context.Context, in *DKGInfoPacket, opts ...grpc.CallOption) (*Empty, error)
	// Setup is doing the DKG setup phase
	FreshDKG(ctx context.Context, in *DKGPacket, opts ...grpc.CallOption) (*Empty, error)
	// Reshare performs the resharing phase
	ReshareDKG(ctx context.Context, in *ResharePacket, opts ...grpc.CallOption) (*Empty, error)
	// PartialBeacon sends its partial beacon to another node
	PartialBeacon(ctx context.Context, in *PartialBeaconPacket, opts ...grpc.CallOption) (*Empty, error)
	SyncChain(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (Protocol_SyncChainClient, error)
}

type protocolClient struct {
	cc grpc.ClientConnInterface
}

func NewProtocolClient(cc grpc.ClientConnInterface) ProtocolClient {
	return &protocolClient{cc}
}

func (c *protocolClient) GetIdentity(ctx context.Context, in *IdentityRequest, opts ...grpc.CallOption) (*Identity, error) {
	out := new(Identity)
	err := c.cc.Invoke(ctx, "/drand.Protocol/GetIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) SignalDKGParticipant(ctx context.Context, in *SignalDKGPacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/drand.Protocol/SignalDKGParticipant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) PushDKGInfo(ctx context.Context, in *DKGInfoPacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/drand.Protocol/PushDKGInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) FreshDKG(ctx context.Context, in *DKGPacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/drand.Protocol/FreshDKG", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) ReshareDKG(ctx context.Context, in *ResharePacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/drand.Protocol/ReshareDKG", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) PartialBeacon(ctx context.Context, in *PartialBeaconPacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/drand.Protocol/PartialBeacon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protocolClient) SyncChain(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (Protocol_SyncChainClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Protocol_serviceDesc.Streams[0], "/drand.Protocol/SyncChain", opts...)
	if err != nil {
		return nil, err
	}
	x := &protocolSyncChainClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Protocol_SyncChainClient interface {
	Recv() (*BeaconPacket, error)
	grpc.ClientStream
}

type protocolSyncChainClient struct {
	grpc.ClientStream
}

func (x *protocolSyncChainClient) Recv() (*BeaconPacket, error) {
	m := new(BeaconPacket)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ProtocolServer is the server API for Protocol service.
type ProtocolServer interface {
	// GetIdentity returns the identity of the drand node
	GetIdentity(context.Context, *IdentityRequest) (*Identity, error)
	// SignalDKGParticipant is called by non-coordinators nodes that sends their
	// public keys and secret proof they have to the coordinator so that he can
	// create the group.
	SignalDKGParticipant(context.Context, *SignalDKGPacket) (*Empty, error)
	// PushDKGInfo is called by the coordinator to push the group he created
	// from all received keys and as well other information such as the time of
	// starting the DKG.
	PushDKGInfo(context.Context, *DKGInfoPacket) (*Empty, error)
	// Setup is doing the DKG setup phase
	FreshDKG(context.Context, *DKGPacket) (*Empty, error)
	// Reshare performs the resharing phase
	ReshareDKG(context.Context, *ResharePacket) (*Empty, error)
	// PartialBeacon sends its partial beacon to another node
	PartialBeacon(context.Context, *PartialBeaconPacket) (*Empty, error)
	SyncChain(*SyncRequest, Protocol_SyncChainServer) error
}

// UnimplementedProtocolServer can be embedded to have forward compatible implementations.
type UnimplementedProtocolServer struct {
}

func (*UnimplementedProtocolServer) GetIdentity(ctx context.Context, req *IdentityRequest) (*Identity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdentity not implemented")
}
func (*UnimplementedProtocolServer) SignalDKGParticipant(ctx context.Context, req *SignalDKGPacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignalDKGParticipant not implemented")
}
func (*UnimplementedProtocolServer) PushDKGInfo(ctx context.Context, req *DKGInfoPacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushDKGInfo not implemented")
}
func (*UnimplementedProtocolServer) FreshDKG(ctx context.Context, req *DKGPacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FreshDKG not implemented")
}
func (*UnimplementedProtocolServer) ReshareDKG(ctx context.Context, req *ResharePacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReshareDKG not implemented")
}
func (*UnimplementedProtocolServer) PartialBeacon(ctx context.Context, req *PartialBeaconPacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PartialBeacon not implemented")
}
func (*UnimplementedProtocolServer) SyncChain(req *SyncRequest, srv Protocol_SyncChainServer) error {
	return status.Errorf(codes.Unimplemented, "method SyncChain not implemented")
}

func RegisterProtocolServer(s *grpc.Server, srv ProtocolServer) {
	s.RegisterService(&_Protocol_serviceDesc, srv)
}

func _Protocol_GetIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdentityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).GetIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/GetIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).GetIdentity(ctx, req.(*IdentityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_SignalDKGParticipant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignalDKGPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).SignalDKGParticipant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/SignalDKGParticipant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).SignalDKGParticipant(ctx, req.(*SignalDKGPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_PushDKGInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DKGInfoPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).PushDKGInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/PushDKGInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).PushDKGInfo(ctx, req.(*DKGInfoPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_FreshDKG_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DKGPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).FreshDKG(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/FreshDKG",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).FreshDKG(ctx, req.(*DKGPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_ReshareDKG_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResharePacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).ReshareDKG(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/ReshareDKG",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).ReshareDKG(ctx, req.(*ResharePacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_PartialBeacon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartialBeaconPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtocolServer).PartialBeacon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/drand.Protocol/PartialBeacon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtocolServer).PartialBeacon(ctx, req.(*PartialBeaconPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Protocol_SyncChain_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SyncRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ProtocolServer).SyncChain(m, &protocolSyncChainServer{stream})
}

type Protocol_SyncChainServer interface {
	Send(*BeaconPacket) error
	grpc.ServerStream
}

type protocolSyncChainServer struct {
	grpc.ServerStream
}

func (x *protocolSyncChainServer) Send(m *BeaconPacket) error {
	return x.ServerStream.SendMsg(m)
}

var _Protocol_serviceDesc = grpc.ServiceDesc{
	ServiceName: "drand.Protocol",
	HandlerType: (*ProtocolServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIdentity",
			Handler:    _Protocol_GetIdentity_Handler,
		},
		{
			MethodName: "SignalDKGParticipant",
			Handler:    _Protocol_SignalDKGParticipant_Handler,
		},
		{
			MethodName: "PushDKGInfo",
			Handler:    _Protocol_PushDKGInfo_Handler,
		},
		{
			MethodName: "FreshDKG",
			Handler:    _Protocol_FreshDKG_Handler,
		},
		{
			MethodName: "ReshareDKG",
			Handler:    _Protocol_ReshareDKG_Handler,
		},
		{
			MethodName: "PartialBeacon",
			Handler:    _Protocol_PartialBeacon_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SyncChain",
			Handler:       _Protocol_SyncChain_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "drand/protocol.proto",
}
