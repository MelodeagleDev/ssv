// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/bloxapp/ssv/ibft/consensus/validation/msgs.proto

package validation

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	v1alpha1 "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
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

type InputValue struct {
	// Types that are valid to be assigned to Data:
	//	*InputValue_AttestationData
	//	*InputValue_BeaconBlock
	Data                 isInputValue_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *InputValue) Reset()         { *m = InputValue{} }
func (m *InputValue) String() string { return proto.CompactTextString(m) }
func (*InputValue) ProtoMessage()    {}
func (*InputValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afacd051edad517, []int{0}
}

func (m *InputValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InputValue.Unmarshal(m, b)
}
func (m *InputValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InputValue.Marshal(b, m, deterministic)
}
func (m *InputValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InputValue.Merge(m, src)
}
func (m *InputValue) XXX_Size() int {
	return xxx_messageInfo_InputValue.Size(m)
}
func (m *InputValue) XXX_DiscardUnknown() {
	xxx_messageInfo_InputValue.DiscardUnknown(m)
}

var xxx_messageInfo_InputValue proto.InternalMessageInfo

type isInputValue_Data interface {
	isInputValue_Data()
}

type InputValue_AttestationData struct {
	AttestationData *v1alpha1.AttestationData `protobuf:"bytes,1,opt,name=attestation_data,json=attestationData,proto3,oneof"`
}

type InputValue_BeaconBlock struct {
	BeaconBlock *v1alpha1.BeaconBlock `protobuf:"bytes,2,opt,name=beacon_block,json=beaconBlock,proto3,oneof"`
}

func (*InputValue_AttestationData) isInputValue_Data() {}

func (*InputValue_BeaconBlock) isInputValue_Data() {}

func (m *InputValue) GetData() isInputValue_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *InputValue) GetAttestationData() *v1alpha1.AttestationData {
	if x, ok := m.GetData().(*InputValue_AttestationData); ok {
		return x.AttestationData
	}
	return nil
}

func (m *InputValue) GetBeaconBlock() *v1alpha1.BeaconBlock {
	if x, ok := m.GetData().(*InputValue_BeaconBlock); ok {
		return x.BeaconBlock
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*InputValue) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*InputValue_AttestationData)(nil),
		(*InputValue_BeaconBlock)(nil),
	}
}

func init() {
	proto.RegisterType((*InputValue)(nil), "validation.InputValue")
}

func init() {
	proto.RegisterFile("github.com/bloxapp/ssv/ibft/consensus/validation/msgs.proto", fileDescriptor_6afacd051edad517)
}

var fileDescriptor_6afacd051edad517 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x1b, 0x84, 0x3a, 0xb8, 0x48, 0xa0, 0x4c, 0x55, 0x07, 0x40, 0x1d, 0x10, 0x93, 0x4d,
	0xcb, 0xc8, 0x44, 0x84, 0x44, 0x59, 0x8b, 0xc4, 0xc0, 0x52, 0x9d, 0xd3, 0xa3, 0x89, 0x70, 0x7c,
	0x56, 0xef, 0x1c, 0xf1, 0x5c, 0x3c, 0x21, 0x8a, 0x2b, 0x94, 0x66, 0x60, 0x60, 0xb3, 0xff, 0xff,
	0xd3, 0x67, 0xdf, 0xa9, 0x87, 0x5d, 0x2d, 0x55, 0xb4, 0xba, 0xa4, 0xc6, 0x58, 0x47, 0x5f, 0x10,
	0x82, 0x61, 0x6e, 0x4d, 0x6d, 0x3f, 0xc4, 0x94, 0xe4, 0x19, 0x3d, 0x47, 0x36, 0x2d, 0xb8, 0x7a,
	0x0b, 0x52, 0x93, 0x37, 0x0d, 0xef, 0x58, 0x87, 0x3d, 0x09, 0xe5, 0xaa, 0x8f, 0x67, 0x97, 0x28,
	0x95, 0x69, 0x17, 0xe0, 0x42, 0x05, 0x0b, 0x03, 0x22, 0xc8, 0x92, 0x9a, 0x03, 0x3b, 0xbb, 0x1a,
	0xf4, 0x16, 0xa1, 0x24, 0xbf, 0xb1, 0x8e, 0xca, 0xcf, 0x03, 0x30, 0xff, 0xce, 0x94, 0x7a, 0xf1,
	0x21, 0xca, 0x1b, 0xb8, 0x88, 0xf9, 0xab, 0xba, 0x38, 0x92, 0x6c, 0xb6, 0x20, 0x30, 0xcd, 0xae,
	0xb3, 0xdb, 0xc9, 0xf2, 0x46, 0xa3, 0x54, 0xb8, 0xc7, 0xd8, 0x74, 0x07, 0xfd, 0xeb, 0xd4, 0x8f,
	0x3d, 0xfe, 0x04, 0x02, 0xab, 0xd1, 0xfa, 0x1c, 0x86, 0x51, 0xfe, 0xac, 0xce, 0x8e, 0x5f, 0x9e,
	0x9e, 0x24, 0xe1, 0xfc, 0x0f, 0x61, 0x91, 0xd0, 0xa2, 0x23, 0x57, 0xa3, 0xf5, 0xc4, 0xf6, 0xd7,
	0x62, 0xac, 0x4e, 0xbb, 0x1f, 0x15, 0xcb, 0xf7, 0xbb, 0xff, 0x2e, 0xd0, 0x8e, 0xd3, 0xbc, 0xf7,
	0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe0, 0xed, 0x83, 0xbb, 0x7b, 0x01, 0x00, 0x00,
}
