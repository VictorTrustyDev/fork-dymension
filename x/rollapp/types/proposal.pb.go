// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/rollapp/proposal.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SubmitFraudProposal struct {
	Title       string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// The rollapp id
	RollappId string `protobuf:"bytes,3,opt,name=rollapp_id,json=rollappId,proto3" json:"rollapp_id,omitempty"`
	// The ibc client id of the rollapp
	IbcClientId string `protobuf:"bytes,4,opt,name=ibc_client_id,json=ibcClientId,proto3" json:"ibc_client_id,omitempty"`
	// The height of the fraudelent block
	FraudelentHeight uint64 `protobuf:"varint,5,opt,name=fraudelent_height,json=fraudelentHeight,proto3" json:"fraudelent_height,omitempty"`
	// The address of the fraudelent sequencer
	FraudelentSequencerAddress string `protobuf:"bytes,6,opt,name=fraudelent_sequencer_address,json=fraudelentSequencerAddress,proto3" json:"fraudelent_sequencer_address,omitempty"`
}

func (m *SubmitFraudProposal) Reset()      { *m = SubmitFraudProposal{} }
func (*SubmitFraudProposal) ProtoMessage() {}
func (*SubmitFraudProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_23c44e927b26bbf5, []int{0}
}
func (m *SubmitFraudProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubmitFraudProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubmitFraudProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubmitFraudProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitFraudProposal.Merge(m, src)
}
func (m *SubmitFraudProposal) XXX_Size() int {
	return m.Size()
}
func (m *SubmitFraudProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitFraudProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitFraudProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SubmitFraudProposal)(nil), "dymensionxyz.dymension.rollapp.SubmitFraudProposal")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/rollapp/proposal.proto", fileDescriptor_23c44e927b26bbf5)
}

var fileDescriptor_23c44e927b26bbf5 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xbf, 0x4e, 0xeb, 0x30,
	0x18, 0xc5, 0xe3, 0xde, 0xb6, 0x52, 0x7d, 0x41, 0x02, 0xd3, 0x21, 0xaa, 0xc0, 0xad, 0x3a, 0x55,
	0x42, 0x24, 0x43, 0x99, 0x98, 0xf8, 0x23, 0x21, 0xba, 0x20, 0xd4, 0x6e, 0x2c, 0x55, 0x12, 0x9b,
	0xd4, 0x92, 0x1b, 0x1b, 0xdb, 0x41, 0x2d, 0x4f, 0xc0, 0x82, 0xc4, 0xc8, 0xd8, 0xc7, 0x61, 0xec,
	0xc8, 0x88, 0xda, 0x85, 0xc7, 0x40, 0x71, 0x42, 0x9b, 0x85, 0xcd, 0xe7, 0x7c, 0x3f, 0x1f, 0x5b,
	0xdf, 0x81, 0x27, 0x64, 0x3e, 0xa5, 0x89, 0x66, 0x22, 0x99, 0xcd, 0x9f, 0xfd, 0x8d, 0xf0, 0x95,
	0xe0, 0x3c, 0x90, 0xd2, 0x97, 0x4a, 0x48, 0xa1, 0x03, 0xee, 0x49, 0x25, 0x8c, 0x40, 0xb8, 0x8c,
	0x7b, 0x1b, 0xe1, 0x15, 0x78, 0xab, 0x19, 0x8b, 0x58, 0x58, 0xd4, 0xcf, 0x4e, 0xf9, 0xad, 0xee,
	0x6b, 0x05, 0x1e, 0x8c, 0xd2, 0x70, 0xca, 0xcc, 0xb5, 0x0a, 0x52, 0x72, 0x57, 0x64, 0xa2, 0x26,
	0xac, 0x19, 0x66, 0x38, 0x75, 0x41, 0x07, 0xf4, 0x1a, 0xc3, 0x5c, 0xa0, 0x0e, 0xfc, 0x4f, 0xa8,
	0x8e, 0x14, 0x93, 0x86, 0x89, 0xc4, 0xad, 0xd8, 0x59, 0xd9, 0x42, 0x47, 0x10, 0x16, 0x0f, 0x8e,
	0x19, 0x71, 0xff, 0x59, 0xa0, 0x51, 0x38, 0x03, 0x82, 0xba, 0x70, 0x97, 0x85, 0xd1, 0x38, 0xe2,
	0x8c, 0x26, 0x26, 0x23, 0xaa, 0x79, 0x04, 0x0b, 0xa3, 0x2b, 0xeb, 0x0d, 0x08, 0x3a, 0x86, 0xfb,
	0x0f, 0xd9, 0x5f, 0x28, 0xcf, 0x98, 0x09, 0x65, 0xf1, 0xc4, 0xb8, 0xb5, 0x0e, 0xe8, 0x55, 0x87,
	0x7b, 0xdb, 0xc1, 0x8d, 0xf5, 0xd1, 0x39, 0x3c, 0x2c, 0xc1, 0x9a, 0x3e, 0xa6, 0x34, 0x89, 0xa8,
	0x1a, 0x07, 0x84, 0x28, 0xaa, 0xb5, 0x5b, 0xb7, 0xf9, 0xad, 0x2d, 0x33, 0xfa, 0x45, 0x2e, 0x72,
	0xe2, 0x6c, 0xe7, 0x65, 0xd1, 0x76, 0xde, 0x17, 0x6d, 0xe7, 0x7b, 0xd1, 0x06, 0x97, 0xb7, 0x1f,
	0x2b, 0x0c, 0x96, 0x2b, 0x0c, 0xbe, 0x56, 0x18, 0xbc, 0xad, 0xb1, 0xb3, 0x5c, 0x63, 0xe7, 0x73,
	0x8d, 0x9d, 0xfb, 0xd3, 0x98, 0x99, 0x49, 0x1a, 0x7a, 0x91, 0x98, 0xfa, 0x7f, 0x34, 0xf3, 0xd4,
	0xf7, 0x67, 0x9b, 0x7a, 0xcc, 0x5c, 0x52, 0x1d, 0xd6, 0xed, 0x9a, 0xfb, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x04, 0x60, 0x11, 0xe4, 0xcd, 0x01, 0x00, 0x00,
}

func (this *SubmitFraudProposal) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SubmitFraudProposal)
	if !ok {
		that2, ok := that.(SubmitFraudProposal)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.RollappId != that1.RollappId {
		return false
	}
	if this.IbcClientId != that1.IbcClientId {
		return false
	}
	if this.FraudelentHeight != that1.FraudelentHeight {
		return false
	}
	if this.FraudelentSequencerAddress != that1.FraudelentSequencerAddress {
		return false
	}
	return true
}
func (m *SubmitFraudProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubmitFraudProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubmitFraudProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FraudelentSequencerAddress) > 0 {
		i -= len(m.FraudelentSequencerAddress)
		copy(dAtA[i:], m.FraudelentSequencerAddress)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.FraudelentSequencerAddress)))
		i--
		dAtA[i] = 0x32
	}
	if m.FraudelentHeight != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.FraudelentHeight))
		i--
		dAtA[i] = 0x28
	}
	if len(m.IbcClientId) > 0 {
		i -= len(m.IbcClientId)
		copy(dAtA[i:], m.IbcClientId)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.IbcClientId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.RollappId) > 0 {
		i -= len(m.RollappId)
		copy(dAtA[i:], m.RollappId)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.RollappId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SubmitFraudProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.RollappId)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.IbcClientId)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.FraudelentHeight != 0 {
		n += 1 + sovProposal(uint64(m.FraudelentHeight))
	}
	l = len(m.FraudelentSequencerAddress)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	return n
}

func sovProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SubmitFraudProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SubmitFraudProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubmitFraudProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RollappId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RollappId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FraudelentHeight", wireType)
			}
			m.FraudelentHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FraudelentHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FraudelentSequencerAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FraudelentSequencerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposal = fmt.Errorf("proto: unexpected end of group")
)
