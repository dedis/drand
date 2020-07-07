// Code generated by protoc-gen-go. DO NOT EDIT.
// source: drand/control.proto

package drand

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

// SetupInfoPacket contains all information necessary to run an "automatic"
// setup phase where the designated leader acts as a coordinator as to what is
// the group file and when does the chain starts.
type SetupInfoPacket struct {
	Leader bool `protobuf:"varint,1,opt,name=leader,proto3" json:"leader,omitempty"`
	// LeaderAddress is only used by non-leader
	LeaderAddress string `protobuf:"bytes,2,opt,name=leader_address,json=leaderAddress,proto3" json:"leader_address,omitempty"`
	// LeaderTls is only used by non-leader
	LeaderTls bool `protobuf:"varint,3,opt,name=leader_tls,json=leaderTls,proto3" json:"leader_tls,omitempty"`
	// the expected number of nodes the group must have
	Nodes uint32 `protobuf:"varint,4,opt,name=nodes,proto3" json:"nodes,omitempty"`
	// the threshold to set to the group
	Threshold uint32 `protobuf:"varint,5,opt,name=threshold,proto3" json:"threshold,omitempty"`
	// timeout of the dkg - it is used for transitioning to the different phases of
	// the dkg (deal, responses and justifications if needed). Unit is in seconds.
	Timeout uint32 `protobuf:"varint,6,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// This field is used by the coordinator to set a genesis time or transition
	// time for the beacon to start. It normally takes time.Now() +
	// beacon_offset.  This offset MUST be superior to the time it takes to
	// run the DKG, even under "malicious case" when the dkg takes longer.
	// In such cases, the dkg takes 3 * timeout time to finish because of the
	// three phases: deal, responses and justifications.
	// XXX: should find a way to designate the time *after* the DKG - beacon
	// generation and dkg should be more separated.
	BeaconOffset uint32 `protobuf:"varint,7,opt,name=beacon_offset,json=beaconOffset,proto3" json:"beacon_offset,omitempty"`
	// dkg_offset is used to set the time for which nodes should start the DKG.
	// To avoid any concurrency / networking effect where nodes start the DKG
	// while some others still haven't received the group configuration, the
	// coordinator do this in two steps: first, send the group configuration to
	// every node, and then every node start at the specified time. This offset
	// is set to be sufficiently large such that with high confidence all nodes
	// received the group file by then.
	DkgOffset uint32 `protobuf:"varint,8,opt,name=dkg_offset,json=dkgOffset,proto3" json:"dkg_offset,omitempty"`
	// the secret used to authentify group members
	Secret               []byte   `protobuf:"bytes,9,opt,name=secret,proto3" json:"secret,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetupInfoPacket) Reset()         { *m = SetupInfoPacket{} }
func (m *SetupInfoPacket) String() string { return proto.CompactTextString(m) }
func (*SetupInfoPacket) ProtoMessage()    {}
func (*SetupInfoPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{0}
}

func (m *SetupInfoPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetupInfoPacket.Unmarshal(m, b)
}
func (m *SetupInfoPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetupInfoPacket.Marshal(b, m, deterministic)
}
func (m *SetupInfoPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetupInfoPacket.Merge(m, src)
}
func (m *SetupInfoPacket) XXX_Size() int {
	return xxx_messageInfo_SetupInfoPacket.Size(m)
}
func (m *SetupInfoPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_SetupInfoPacket.DiscardUnknown(m)
}

var xxx_messageInfo_SetupInfoPacket proto.InternalMessageInfo

func (m *SetupInfoPacket) GetLeader() bool {
	if m != nil {
		return m.Leader
	}
	return false
}

func (m *SetupInfoPacket) GetLeaderAddress() string {
	if m != nil {
		return m.LeaderAddress
	}
	return ""
}

func (m *SetupInfoPacket) GetLeaderTls() bool {
	if m != nil {
		return m.LeaderTls
	}
	return false
}

func (m *SetupInfoPacket) GetNodes() uint32 {
	if m != nil {
		return m.Nodes
	}
	return 0
}

func (m *SetupInfoPacket) GetThreshold() uint32 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

func (m *SetupInfoPacket) GetTimeout() uint32 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *SetupInfoPacket) GetBeaconOffset() uint32 {
	if m != nil {
		return m.BeaconOffset
	}
	return 0
}

func (m *SetupInfoPacket) GetDkgOffset() uint32 {
	if m != nil {
		return m.DkgOffset
	}
	return 0
}

func (m *SetupInfoPacket) GetSecret() []byte {
	if m != nil {
		return m.Secret
	}
	return nil
}

type InitDKGPacket struct {
	Info    *SetupInfoPacket `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	Entropy *EntropyInfo     `protobuf:"bytes,2,opt,name=entropy,proto3" json:"entropy,omitempty"`
	// the period time of the beacon in seconds.
	// used only in a fresh dkg
	BeaconPeriod         uint32   `protobuf:"varint,3,opt,name=beacon_period,json=beaconPeriod,proto3" json:"beacon_period,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitDKGPacket) Reset()         { *m = InitDKGPacket{} }
func (m *InitDKGPacket) String() string { return proto.CompactTextString(m) }
func (*InitDKGPacket) ProtoMessage()    {}
func (*InitDKGPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{1}
}

func (m *InitDKGPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitDKGPacket.Unmarshal(m, b)
}
func (m *InitDKGPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitDKGPacket.Marshal(b, m, deterministic)
}
func (m *InitDKGPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitDKGPacket.Merge(m, src)
}
func (m *InitDKGPacket) XXX_Size() int {
	return xxx_messageInfo_InitDKGPacket.Size(m)
}
func (m *InitDKGPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_InitDKGPacket.DiscardUnknown(m)
}

var xxx_messageInfo_InitDKGPacket proto.InternalMessageInfo

func (m *InitDKGPacket) GetInfo() *SetupInfoPacket {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *InitDKGPacket) GetEntropy() *EntropyInfo {
	if m != nil {
		return m.Entropy
	}
	return nil
}

func (m *InitDKGPacket) GetBeaconPeriod() uint32 {
	if m != nil {
		return m.BeaconPeriod
	}
	return 0
}

// EntropyInfo contains information about external entropy sources
// can be optional
type EntropyInfo struct {
	// the path to the script to run that returns random bytes when called
	Script string `protobuf:"bytes,1,opt,name=script,proto3" json:"script,omitempty"`
	// do we only take this entropy source or mix it with /dev/urandom
	UserOnly             bool     `protobuf:"varint,10,opt,name=userOnly,proto3" json:"userOnly,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EntropyInfo) Reset()         { *m = EntropyInfo{} }
func (m *EntropyInfo) String() string { return proto.CompactTextString(m) }
func (*EntropyInfo) ProtoMessage()    {}
func (*EntropyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{2}
}

func (m *EntropyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntropyInfo.Unmarshal(m, b)
}
func (m *EntropyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntropyInfo.Marshal(b, m, deterministic)
}
func (m *EntropyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntropyInfo.Merge(m, src)
}
func (m *EntropyInfo) XXX_Size() int {
	return xxx_messageInfo_EntropyInfo.Size(m)
}
func (m *EntropyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EntropyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EntropyInfo proto.InternalMessageInfo

func (m *EntropyInfo) GetScript() string {
	if m != nil {
		return m.Script
	}
	return ""
}

func (m *EntropyInfo) GetUserOnly() bool {
	if m != nil {
		return m.UserOnly
	}
	return false
}

// ReshareRequest contains references to the old and new group to perform the
// resharing protocol.
type InitResharePacket struct {
	// Old group that needs to issue the shares for the new group
	// NOTE: It can be empty / nil. In that case, the drand node will try to
	// load the group he belongs to at the moment, if any, and use it as the old
	// group.
	Old                  *GroupInfo       `protobuf:"bytes,1,opt,name=old,proto3" json:"old,omitempty"`
	Info                 *SetupInfoPacket `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *InitResharePacket) Reset()         { *m = InitResharePacket{} }
func (m *InitResharePacket) String() string { return proto.CompactTextString(m) }
func (*InitResharePacket) ProtoMessage()    {}
func (*InitResharePacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{3}
}

func (m *InitResharePacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitResharePacket.Unmarshal(m, b)
}
func (m *InitResharePacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitResharePacket.Marshal(b, m, deterministic)
}
func (m *InitResharePacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitResharePacket.Merge(m, src)
}
func (m *InitResharePacket) XXX_Size() int {
	return xxx_messageInfo_InitResharePacket.Size(m)
}
func (m *InitResharePacket) XXX_DiscardUnknown() {
	xxx_messageInfo_InitResharePacket.DiscardUnknown(m)
}

var xxx_messageInfo_InitResharePacket proto.InternalMessageInfo

func (m *InitResharePacket) GetOld() *GroupInfo {
	if m != nil {
		return m.Old
	}
	return nil
}

func (m *InitResharePacket) GetInfo() *SetupInfoPacket {
	if m != nil {
		return m.Info
	}
	return nil
}

// GroupInfo holds the information to load a group information such as the nodes
// and the genesis etc. Currently only the loading of a group via filesystem is
// supported although the basis to support loading a group from a URI is setup.
// For example, for new nodes that wants to join a network, they could point to
// the URL that returns a group definition, for example at one of the currently
// running node.
type GroupInfo struct {
	// Types that are valid to be assigned to Location:
	//	*GroupInfo_Path
	//	*GroupInfo_Url
	Location             isGroupInfo_Location `protobuf_oneof:"location"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GroupInfo) Reset()         { *m = GroupInfo{} }
func (m *GroupInfo) String() string { return proto.CompactTextString(m) }
func (*GroupInfo) ProtoMessage()    {}
func (*GroupInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{4}
}

func (m *GroupInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupInfo.Unmarshal(m, b)
}
func (m *GroupInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupInfo.Marshal(b, m, deterministic)
}
func (m *GroupInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupInfo.Merge(m, src)
}
func (m *GroupInfo) XXX_Size() int {
	return xxx_messageInfo_GroupInfo.Size(m)
}
func (m *GroupInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupInfo.DiscardUnknown(m)
}

var xxx_messageInfo_GroupInfo proto.InternalMessageInfo

type isGroupInfo_Location interface {
	isGroupInfo_Location()
}

type GroupInfo_Path struct {
	Path string `protobuf:"bytes,1,opt,name=path,proto3,oneof"`
}

type GroupInfo_Url struct {
	Url string `protobuf:"bytes,2,opt,name=url,proto3,oneof"`
}

func (*GroupInfo_Path) isGroupInfo_Location() {}

func (*GroupInfo_Url) isGroupInfo_Location() {}

func (m *GroupInfo) GetLocation() isGroupInfo_Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *GroupInfo) GetPath() string {
	if x, ok := m.GetLocation().(*GroupInfo_Path); ok {
		return x.Path
	}
	return ""
}

func (m *GroupInfo) GetUrl() string {
	if x, ok := m.GetLocation().(*GroupInfo_Url); ok {
		return x.Url
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*GroupInfo) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*GroupInfo_Path)(nil),
		(*GroupInfo_Url)(nil),
	}
}

// ShareRequest requests the private share of a drand node
type ShareRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShareRequest) Reset()         { *m = ShareRequest{} }
func (m *ShareRequest) String() string { return proto.CompactTextString(m) }
func (*ShareRequest) ProtoMessage()    {}
func (*ShareRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{5}
}

func (m *ShareRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShareRequest.Unmarshal(m, b)
}
func (m *ShareRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShareRequest.Marshal(b, m, deterministic)
}
func (m *ShareRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShareRequest.Merge(m, src)
}
func (m *ShareRequest) XXX_Size() int {
	return xxx_messageInfo_ShareRequest.Size(m)
}
func (m *ShareRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShareRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShareRequest proto.InternalMessageInfo

// ShareResponse holds the private share of a drand node
type ShareResponse struct {
	Index                uint32   `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Share                []byte   `protobuf:"bytes,3,opt,name=share,proto3" json:"share,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShareResponse) Reset()         { *m = ShareResponse{} }
func (m *ShareResponse) String() string { return proto.CompactTextString(m) }
func (*ShareResponse) ProtoMessage()    {}
func (*ShareResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{6}
}

func (m *ShareResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShareResponse.Unmarshal(m, b)
}
func (m *ShareResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShareResponse.Marshal(b, m, deterministic)
}
func (m *ShareResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShareResponse.Merge(m, src)
}
func (m *ShareResponse) XXX_Size() int {
	return xxx_messageInfo_ShareResponse.Size(m)
}
func (m *ShareResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShareResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShareResponse proto.InternalMessageInfo

func (m *ShareResponse) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ShareResponse) GetShare() []byte {
	if m != nil {
		return m.Share
	}
	return nil
}

type Ping struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{7}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

type Pong struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pong) Reset()         { *m = Pong{} }
func (m *Pong) String() string { return proto.CompactTextString(m) }
func (*Pong) ProtoMessage()    {}
func (*Pong) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{8}
}

func (m *Pong) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pong.Unmarshal(m, b)
}
func (m *Pong) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pong.Marshal(b, m, deterministic)
}
func (m *Pong) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pong.Merge(m, src)
}
func (m *Pong) XXX_Size() int {
	return xxx_messageInfo_Pong.Size(m)
}
func (m *Pong) XXX_DiscardUnknown() {
	xxx_messageInfo_Pong.DiscardUnknown(m)
}

var xxx_messageInfo_Pong proto.InternalMessageInfo

// PublicKeyRequest requests the public key of a drand node
type PublicKeyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicKeyRequest) Reset()         { *m = PublicKeyRequest{} }
func (m *PublicKeyRequest) String() string { return proto.CompactTextString(m) }
func (*PublicKeyRequest) ProtoMessage()    {}
func (*PublicKeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{9}
}

func (m *PublicKeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicKeyRequest.Unmarshal(m, b)
}
func (m *PublicKeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicKeyRequest.Marshal(b, m, deterministic)
}
func (m *PublicKeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicKeyRequest.Merge(m, src)
}
func (m *PublicKeyRequest) XXX_Size() int {
	return xxx_messageInfo_PublicKeyRequest.Size(m)
}
func (m *PublicKeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicKeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PublicKeyRequest proto.InternalMessageInfo

// PublicKeyResponse holds the public key of a drand node
type PublicKeyResponse struct {
	PubKey               []byte   `protobuf:"bytes,2,opt,name=pubKey,proto3" json:"pubKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicKeyResponse) Reset()         { *m = PublicKeyResponse{} }
func (m *PublicKeyResponse) String() string { return proto.CompactTextString(m) }
func (*PublicKeyResponse) ProtoMessage()    {}
func (*PublicKeyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{10}
}

func (m *PublicKeyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicKeyResponse.Unmarshal(m, b)
}
func (m *PublicKeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicKeyResponse.Marshal(b, m, deterministic)
}
func (m *PublicKeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicKeyResponse.Merge(m, src)
}
func (m *PublicKeyResponse) XXX_Size() int {
	return xxx_messageInfo_PublicKeyResponse.Size(m)
}
func (m *PublicKeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicKeyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PublicKeyResponse proto.InternalMessageInfo

func (m *PublicKeyResponse) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

// PrivateKeyRequest requests the private key of a drand node
type PrivateKeyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateKeyRequest) Reset()         { *m = PrivateKeyRequest{} }
func (m *PrivateKeyRequest) String() string { return proto.CompactTextString(m) }
func (*PrivateKeyRequest) ProtoMessage()    {}
func (*PrivateKeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{11}
}

func (m *PrivateKeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrivateKeyRequest.Unmarshal(m, b)
}
func (m *PrivateKeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrivateKeyRequest.Marshal(b, m, deterministic)
}
func (m *PrivateKeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivateKeyRequest.Merge(m, src)
}
func (m *PrivateKeyRequest) XXX_Size() int {
	return xxx_messageInfo_PrivateKeyRequest.Size(m)
}
func (m *PrivateKeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivateKeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PrivateKeyRequest proto.InternalMessageInfo

// PrivateKeyResponse holds the private key of a drand node
type PrivateKeyResponse struct {
	PriKey               []byte   `protobuf:"bytes,2,opt,name=priKey,proto3" json:"priKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrivateKeyResponse) Reset()         { *m = PrivateKeyResponse{} }
func (m *PrivateKeyResponse) String() string { return proto.CompactTextString(m) }
func (*PrivateKeyResponse) ProtoMessage()    {}
func (*PrivateKeyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{12}
}

func (m *PrivateKeyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrivateKeyResponse.Unmarshal(m, b)
}
func (m *PrivateKeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrivateKeyResponse.Marshal(b, m, deterministic)
}
func (m *PrivateKeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrivateKeyResponse.Merge(m, src)
}
func (m *PrivateKeyResponse) XXX_Size() int {
	return xxx_messageInfo_PrivateKeyResponse.Size(m)
}
func (m *PrivateKeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PrivateKeyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PrivateKeyResponse proto.InternalMessageInfo

func (m *PrivateKeyResponse) GetPriKey() []byte {
	if m != nil {
		return m.PriKey
	}
	return nil
}

// CokeyRequest requests the collective key of a drand node
type CokeyRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CokeyRequest) Reset()         { *m = CokeyRequest{} }
func (m *CokeyRequest) String() string { return proto.CompactTextString(m) }
func (*CokeyRequest) ProtoMessage()    {}
func (*CokeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{13}
}

func (m *CokeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CokeyRequest.Unmarshal(m, b)
}
func (m *CokeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CokeyRequest.Marshal(b, m, deterministic)
}
func (m *CokeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CokeyRequest.Merge(m, src)
}
func (m *CokeyRequest) XXX_Size() int {
	return xxx_messageInfo_CokeyRequest.Size(m)
}
func (m *CokeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CokeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CokeyRequest proto.InternalMessageInfo

// CokeyResponse holds the collective key of a drand node
type CokeyResponse struct {
	CoKey                []byte   `protobuf:"bytes,2,opt,name=coKey,proto3" json:"coKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CokeyResponse) Reset()         { *m = CokeyResponse{} }
func (m *CokeyResponse) String() string { return proto.CompactTextString(m) }
func (*CokeyResponse) ProtoMessage()    {}
func (*CokeyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{14}
}

func (m *CokeyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CokeyResponse.Unmarshal(m, b)
}
func (m *CokeyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CokeyResponse.Marshal(b, m, deterministic)
}
func (m *CokeyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CokeyResponse.Merge(m, src)
}
func (m *CokeyResponse) XXX_Size() int {
	return xxx_messageInfo_CokeyResponse.Size(m)
}
func (m *CokeyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CokeyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CokeyResponse proto.InternalMessageInfo

func (m *CokeyResponse) GetCoKey() []byte {
	if m != nil {
		return m.CoKey
	}
	return nil
}

type GroupTOMLResponse struct {
	// TOML-encoded group file
	GroupToml            string   `protobuf:"bytes,1,opt,name=group_toml,json=groupToml,proto3" json:"group_toml,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GroupTOMLResponse) Reset()         { *m = GroupTOMLResponse{} }
func (m *GroupTOMLResponse) String() string { return proto.CompactTextString(m) }
func (*GroupTOMLResponse) ProtoMessage()    {}
func (*GroupTOMLResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{15}
}

func (m *GroupTOMLResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GroupTOMLResponse.Unmarshal(m, b)
}
func (m *GroupTOMLResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GroupTOMLResponse.Marshal(b, m, deterministic)
}
func (m *GroupTOMLResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GroupTOMLResponse.Merge(m, src)
}
func (m *GroupTOMLResponse) XXX_Size() int {
	return xxx_messageInfo_GroupTOMLResponse.Size(m)
}
func (m *GroupTOMLResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GroupTOMLResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GroupTOMLResponse proto.InternalMessageInfo

func (m *GroupTOMLResponse) GetGroupToml() string {
	if m != nil {
		return m.GroupToml
	}
	return ""
}

type ShutdownRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShutdownRequest) Reset()         { *m = ShutdownRequest{} }
func (m *ShutdownRequest) String() string { return proto.CompactTextString(m) }
func (*ShutdownRequest) ProtoMessage()    {}
func (*ShutdownRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{16}
}

func (m *ShutdownRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShutdownRequest.Unmarshal(m, b)
}
func (m *ShutdownRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShutdownRequest.Marshal(b, m, deterministic)
}
func (m *ShutdownRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShutdownRequest.Merge(m, src)
}
func (m *ShutdownRequest) XXX_Size() int {
	return xxx_messageInfo_ShutdownRequest.Size(m)
}
func (m *ShutdownRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShutdownRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShutdownRequest proto.InternalMessageInfo

type ShutdownResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShutdownResponse) Reset()         { *m = ShutdownResponse{} }
func (m *ShutdownResponse) String() string { return proto.CompactTextString(m) }
func (*ShutdownResponse) ProtoMessage()    {}
func (*ShutdownResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dd5961950a69ad7, []int{17}
}

func (m *ShutdownResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShutdownResponse.Unmarshal(m, b)
}
func (m *ShutdownResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShutdownResponse.Marshal(b, m, deterministic)
}
func (m *ShutdownResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShutdownResponse.Merge(m, src)
}
func (m *ShutdownResponse) XXX_Size() int {
	return xxx_messageInfo_ShutdownResponse.Size(m)
}
func (m *ShutdownResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShutdownResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShutdownResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SetupInfoPacket)(nil), "drand.SetupInfoPacket")
	proto.RegisterType((*InitDKGPacket)(nil), "drand.InitDKGPacket")
	proto.RegisterType((*EntropyInfo)(nil), "drand.EntropyInfo")
	proto.RegisterType((*InitResharePacket)(nil), "drand.InitResharePacket")
	proto.RegisterType((*GroupInfo)(nil), "drand.GroupInfo")
	proto.RegisterType((*ShareRequest)(nil), "drand.ShareRequest")
	proto.RegisterType((*ShareResponse)(nil), "drand.ShareResponse")
	proto.RegisterType((*Ping)(nil), "drand.Ping")
	proto.RegisterType((*Pong)(nil), "drand.Pong")
	proto.RegisterType((*PublicKeyRequest)(nil), "drand.PublicKeyRequest")
	proto.RegisterType((*PublicKeyResponse)(nil), "drand.PublicKeyResponse")
	proto.RegisterType((*PrivateKeyRequest)(nil), "drand.PrivateKeyRequest")
	proto.RegisterType((*PrivateKeyResponse)(nil), "drand.PrivateKeyResponse")
	proto.RegisterType((*CokeyRequest)(nil), "drand.CokeyRequest")
	proto.RegisterType((*CokeyResponse)(nil), "drand.CokeyResponse")
	proto.RegisterType((*GroupTOMLResponse)(nil), "drand.GroupTOMLResponse")
	proto.RegisterType((*ShutdownRequest)(nil), "drand.ShutdownRequest")
	proto.RegisterType((*ShutdownResponse)(nil), "drand.ShutdownResponse")
}

func init() {
	proto.RegisterFile("drand/control.proto", fileDescriptor_2dd5961950a69ad7)
}

var fileDescriptor_2dd5961950a69ad7 = []byte{
	// 769 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x55, 0x5d, 0x6f, 0xda, 0x48,
	0x14, 0x85, 0x84, 0x2f, 0x5f, 0x20, 0x09, 0x13, 0x44, 0xbc, 0xd6, 0x46, 0x42, 0x5e, 0x65, 0x17,
	0xed, 0x46, 0x44, 0x62, 0x77, 0xfb, 0x52, 0xb5, 0x6a, 0x42, 0xdb, 0x24, 0x4a, 0xab, 0xa0, 0x49,
	0x9e, 0xfa, 0x12, 0x19, 0x7b, 0x00, 0x0b, 0x33, 0xe3, 0x8e, 0xc7, 0x6d, 0xf3, 0x27, 0xfa, 0x5e,
	0xf5, 0xcf, 0x56, 0xf3, 0x61, 0x63, 0x92, 0x46, 0x7d, 0x01, 0xce, 0xb9, 0x73, 0xaf, 0xef, 0x39,
	0xf7, 0x0e, 0x86, 0xfd, 0x80, 0x7b, 0x34, 0x38, 0xf1, 0x19, 0x15, 0x9c, 0x45, 0xc3, 0x98, 0x33,
	0xc1, 0x50, 0x55, 0x91, 0x0e, 0xca, 0x62, 0xab, 0x15, 0xa3, 0x3a, 0xe4, 0x7e, 0xdb, 0x82, 0xdd,
	0x1b, 0x22, 0xd2, 0xf8, 0x92, 0xce, 0xd8, 0xc4, 0xf3, 0x97, 0x44, 0xa0, 0x1e, 0xd4, 0x22, 0xe2,
	0x05, 0x84, 0xdb, 0xe5, 0x7e, 0x79, 0xd0, 0xc0, 0x06, 0xa1, 0x23, 0xd8, 0xd1, 0xbf, 0xee, 0xbc,
	0x20, 0xe0, 0x24, 0x49, 0xec, 0xad, 0x7e, 0x79, 0x60, 0xe1, 0xb6, 0x66, 0x4f, 0x35, 0x89, 0x0e,
	0x01, 0xcc, 0x31, 0x11, 0x25, 0xf6, 0xb6, 0x2a, 0x61, 0x69, 0xe6, 0x36, 0x4a, 0x50, 0x17, 0xaa,
	0x94, 0x05, 0x24, 0xb1, 0x2b, 0xfd, 0xf2, 0xa0, 0x8d, 0x35, 0x40, 0xbf, 0x83, 0x25, 0x16, 0x9c,
	0x24, 0x0b, 0x16, 0x05, 0x76, 0x55, 0x45, 0xd6, 0x04, 0xb2, 0xa1, 0x2e, 0xc2, 0x15, 0x61, 0xa9,
	0xb0, 0x6b, 0x2a, 0x96, 0x41, 0xf4, 0x07, 0xb4, 0xa7, 0xc4, 0xf3, 0x19, 0xbd, 0x63, 0xb3, 0x59,
	0x42, 0x84, 0x5d, 0x57, 0xf1, 0x96, 0x26, 0xaf, 0x15, 0x27, 0x3b, 0x0a, 0x96, 0xf3, 0xec, 0x44,
	0x43, 0x57, 0x0f, 0x96, 0x73, 0x13, 0xee, 0x41, 0x2d, 0x21, 0x3e, 0x27, 0xc2, 0xb6, 0xfa, 0xe5,
	0x41, 0x0b, 0x1b, 0xe4, 0x7e, 0x2d, 0x43, 0xfb, 0x92, 0x86, 0xe2, 0xf5, 0xd5, 0xb9, 0x71, 0xe6,
	0x6f, 0xa8, 0x84, 0x74, 0xc6, 0x94, 0x2f, 0xcd, 0x51, 0x6f, 0xa8, 0x0c, 0x1d, 0x3e, 0xf0, 0x0f,
	0xab, 0x33, 0xe8, 0x18, 0xea, 0x44, 0x0e, 0x21, 0xbe, 0x57, 0x36, 0x35, 0x47, 0xc8, 0x1c, 0x7f,
	0xa3, 0x59, 0x99, 0x80, 0xb3, 0x23, 0x05, 0x1d, 0x31, 0xe1, 0x21, 0x0b, 0x94, 0x6f, 0xb9, 0x8e,
	0x89, 0xe2, 0xdc, 0x53, 0x68, 0x16, 0x92, 0x55, 0xdf, 0x3e, 0x0f, 0x63, 0xa1, 0xfa, 0xb1, 0xb0,
	0x41, 0xc8, 0x81, 0x46, 0x9a, 0x10, 0x7e, 0x4d, 0xa3, 0x7b, 0x1b, 0x94, 0xfd, 0x39, 0x76, 0x7d,
	0xe8, 0x48, 0x49, 0x98, 0x24, 0x0b, 0x8f, 0x13, 0x23, 0xcb, 0x85, 0x6d, 0x69, 0xbb, 0x56, 0xb5,
	0x67, 0xda, 0x3c, 0xe7, 0x4c, 0xab, 0xc2, 0x32, 0x98, 0x4b, 0xdf, 0xfa, 0xb5, 0x74, 0xf7, 0x14,
	0xac, 0x3c, 0x1b, 0x75, 0xa1, 0x12, 0x7b, 0x62, 0xa1, 0x7b, 0xbc, 0x28, 0x61, 0x85, 0x10, 0x82,
	0xed, 0x94, 0x47, 0x7a, 0x81, 0x2e, 0x4a, 0x58, 0x82, 0x33, 0x80, 0x46, 0xc4, 0x7c, 0x4f, 0x84,
	0x8c, 0xba, 0x3b, 0xd0, 0xba, 0x91, 0x1d, 0x62, 0xf2, 0x31, 0x25, 0x89, 0x70, 0x9f, 0x43, 0xdb,
	0xe0, 0x24, 0x66, 0x34, 0x21, 0x72, 0x8d, 0x42, 0x1a, 0x90, 0x2f, 0xaa, 0x44, 0x1b, 0x6b, 0x20,
	0x59, 0x25, 0x4c, 0xd9, 0xd7, 0xc2, 0x1a, 0xb8, 0x35, 0xa8, 0x4c, 0x42, 0x3a, 0x57, 0xdf, 0x8c,
	0xce, 0x5d, 0x04, 0x7b, 0x93, 0x74, 0x1a, 0x85, 0xfe, 0x15, 0xb9, 0xcf, 0x1e, 0xf0, 0x0f, 0x74,
	0x0a, 0x9c, 0x79, 0x48, 0x0f, 0x6a, 0x71, 0x3a, 0xbd, 0x22, 0x7a, 0x84, 0x2d, 0x6c, 0x90, 0xbb,
	0x0f, 0x9d, 0x09, 0x0f, 0x3f, 0x79, 0x82, 0x14, 0x2a, 0x1c, 0x03, 0x2a, 0x92, 0x85, 0x12, 0x3c,
	0x2c, 0x96, 0x50, 0x48, 0x0a, 0x1c, 0xb3, 0xe5, 0x3a, 0xfb, 0x08, 0xda, 0x06, 0xaf, 0x05, 0xfa,
	0x6c, 0x9d, 0xa7, 0x81, 0x3b, 0x82, 0x8e, 0xb2, 0xf6, 0xf6, 0xfa, 0xfd, 0xbb, 0xfc, 0xe8, 0x21,
	0xc0, 0x5c, 0x92, 0x77, 0x82, 0xad, 0x22, 0xb3, 0x0c, 0x96, 0x62, 0x6e, 0xd9, 0x2a, 0x72, 0x3b,
	0xb0, 0x7b, 0xb3, 0x48, 0x45, 0xc0, 0x3e, 0xd3, 0xec, 0x69, 0x08, 0xf6, 0xd6, 0x94, 0xae, 0x32,
	0xfa, 0x5e, 0x81, 0xfa, 0x58, 0xff, 0x6f, 0xa0, 0x3f, 0xa1, 0x21, 0x1d, 0x93, 0x6e, 0xa1, 0xa6,
	0x99, 0xb5, 0x24, 0x9c, 0x1c, 0x48, 0x1f, 0x4b, 0xe8, 0x7f, 0xa8, 0x9b, 0x1b, 0x82, 0xba, 0x26,
	0xb2, 0x71, 0x63, 0x1c, 0x54, 0xdc, 0x26, 0xcd, 0xb9, 0x25, 0xf4, 0x02, 0x9a, 0x85, 0x2d, 0x44,
	0x76, 0x21, 0x75, 0x63, 0x33, 0x9f, 0x48, 0xff, 0x0f, 0xaa, 0x6a, 0x19, 0xd0, 0x7e, 0xb6, 0x86,
	0x85, 0x55, 0x71, 0xba, 0x9b, 0xa4, 0x56, 0xe7, 0x96, 0xd0, 0x2b, 0xb0, 0xf2, 0x09, 0xa3, 0x83,
	0x4c, 0xc7, 0x83, 0x3d, 0x70, 0xec, 0xc7, 0x81, 0xbc, 0xc2, 0x18, 0x60, 0x3d, 0xe1, 0xbc, 0xeb,
	0x47, 0x9b, 0xe0, 0xfc, 0xf6, 0x93, 0x48, 0x5e, 0xe4, 0x25, 0x58, 0xe3, 0x85, 0x17, 0x52, 0x75,
	0x39, 0xb2, 0x36, 0x72, 0x26, 0x2b, 0xd1, 0x7b, 0x18, 0xc8, 0xc5, 0x3f, 0x33, 0x97, 0xeb, 0x6d,
	0x18, 0xad, 0x0d, 0x50, 0x4c, 0x96, 0xfb, 0x94, 0xe7, 0x8d, 0x6c, 0xe4, 0x28, 0xbf, 0xbe, 0x9b,
	0x6b, 0xe1, 0x1c, 0x3c, 0xe2, 0xb3, 0xb6, 0xcf, 0xfe, 0xfa, 0x70, 0x34, 0x0f, 0xc5, 0x22, 0x9d,
	0x0e, 0x7d, 0xb6, 0x3a, 0xd1, 0x6f, 0x12, 0xfd, 0xa9, 0x5e, 0x24, 0xd3, 0x74, 0xa6, 0xe1, 0xb4,
	0xa6, 0xf0, 0xbf, 0x3f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xef, 0x06, 0x8c, 0x36, 0x8a, 0x06, 0x00,
	0x00,
}
