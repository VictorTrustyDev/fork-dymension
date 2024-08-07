// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/rollapp/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
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

// Params defines the parameters for the module.
type Params struct {
	// dispute_period_in_blocks the number of blocks it takes
	// to change a status of a state from received to finalized.
	// during that period, any user could submit fraud proof
	DisputePeriodInBlocks uint64 `protobuf:"varint,1,opt,name=dispute_period_in_blocks,json=disputePeriodInBlocks,proto3" json:"dispute_period_in_blocks,omitempty" yaml:"dispute_period_in_blocks"`
	// The time (num hub blocks) a sequencer has to post a block, before he will be slashed
	LivenessSlashBlocks uint64 `protobuf:"varint,4,opt,name=liveness_slash_blocks,json=livenessSlashBlocks,proto3" json:"liveness_slash_blocks,omitempty" yaml:"liveness_slash_blocks"`
	// The min gap (num hub blocks) between a sequence of slashes if the sequencer continues to be down
	LivenessSlashInterval uint64 `protobuf:"varint,5,opt,name=liveness_slash_interval,json=livenessSlashInterval,proto3" json:"liveness_slash_interval,omitempty" yaml:"liveness_slash_interval"`
	// The time (num hub blocks) a sequencer can be down after which he will be jailed rather than slashed
	LivenessJailBlocks uint64 `protobuf:"varint,6,opt,name=liveness_jail_blocks,json=livenessJailBlocks,proto3" json:"liveness_jail_blocks,omitempty" yaml:"liveness_jail_blocks"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_75a44aa904ae1ba5, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetDisputePeriodInBlocks() uint64 {
	if m != nil {
		return m.DisputePeriodInBlocks
	}
	return 0
}

func (m *Params) GetLivenessSlashBlocks() uint64 {
	if m != nil {
		return m.LivenessSlashBlocks
	}
	return 0
}

func (m *Params) GetLivenessSlashInterval() uint64 {
	if m != nil {
		return m.LivenessSlashInterval
	}
	return 0
}

func (m *Params) GetLivenessJailBlocks() uint64 {
	if m != nil {
		return m.LivenessJailBlocks
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "dymensionxyz.dymension.rollapp.Params")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/rollapp/params.proto", fileDescriptor_75a44aa904ae1ba5)
}

var fileDescriptor_75a44aa904ae1ba5 = []byte{
	// 367 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x6b, 0xe2, 0x40,
	0x14, 0xc7, 0x13, 0xcd, 0x8a, 0xe4, 0x24, 0x59, 0x65, 0xc5, 0x5d, 0x26, 0x92, 0xbd, 0x2c, 0x2c,
	0x64, 0x10, 0x7b, 0xf2, 0xe8, 0x4d, 0x0f, 0xc5, 0xa6, 0x3d, 0x49, 0x21, 0x4c, 0xe2, 0xa0, 0xd3,
	0x4e, 0x32, 0x21, 0x13, 0x83, 0xe9, 0x5f, 0xd1, 0x63, 0x8f, 0xfd, 0x73, 0x7a, 0xf4, 0xd8, 0x53,
	0x28, 0xfa, 0x1f, 0xe4, 0x5e, 0x28, 0x4e, 0x7e, 0x60, 0x45, 0x6f, 0x79, 0xdf, 0xf7, 0x79, 0x9f,
	0x3c, 0x98, 0xa7, 0xfe, 0x5f, 0x24, 0x1e, 0xf6, 0x39, 0x61, 0xfe, 0x26, 0x79, 0x82, 0x55, 0x01,
	0x43, 0x46, 0x29, 0x0a, 0x02, 0x18, 0xa0, 0x10, 0x79, 0xdc, 0x0c, 0x42, 0x16, 0x31, 0x0d, 0x1c,
	0xc3, 0x66, 0x55, 0x98, 0x05, 0xdc, 0x6b, 0x2f, 0xd9, 0x92, 0x09, 0x14, 0x1e, 0xbe, 0xf2, 0xa9,
	0x1e, 0x70, 0x19, 0xf7, 0x18, 0x87, 0x0e, 0xe2, 0x18, 0xc6, 0x03, 0x07, 0x47, 0x68, 0x00, 0x5d,
	0x46, 0xfc, 0xbc, 0x6f, 0x7c, 0xd6, 0xd4, 0xc6, 0x4c, 0xfc, 0x46, 0xbb, 0x57, 0xbb, 0x0b, 0xc2,
	0x83, 0x75, 0x84, 0xed, 0x00, 0x87, 0x84, 0x2d, 0x6c, 0xe2, 0xdb, 0x0e, 0x65, 0xee, 0x23, 0xef,
	0xca, 0x7d, 0xf9, 0x9f, 0x32, 0xfe, 0x9b, 0xa5, 0xba, 0x9e, 0x20, 0x8f, 0x8e, 0x8c, 0x4b, 0xa4,
	0x61, 0x75, 0x8a, 0xd6, 0x4c, 0x74, 0x26, 0xfe, 0x58, 0xe4, 0xda, 0x9d, 0xda, 0xa1, 0x24, 0xc6,
	0x3e, 0xe6, 0xdc, 0xe6, 0x14, 0xf1, 0x55, 0xa9, 0x56, 0x84, 0xba, 0x9f, 0xa5, 0xfa, 0x9f, 0x5c,
	0x7d, 0x16, 0x33, 0xac, 0x9f, 0x65, 0x7e, 0x7b, 0x88, 0x0b, 0xeb, 0x5c, 0xfd, 0x75, 0x82, 0x13,
	0x3f, 0xc2, 0x61, 0x8c, 0x68, 0xf7, 0x87, 0xf0, 0x1a, 0x59, 0xaa, 0x83, 0xb3, 0xde, 0x12, 0x34,
	0xac, 0xce, 0x37, 0xf3, 0xa4, 0xc8, 0xb5, 0x1b, 0xb5, 0x5d, 0x8d, 0x3c, 0x20, 0x42, 0xcb, 0x85,
	0x1b, 0x42, 0xac, 0x67, 0xa9, 0xfe, 0xfb, 0x44, 0x7c, 0x44, 0x19, 0x96, 0x56, 0xc6, 0x53, 0x44,
	0x68, 0xbe, 0xee, 0x48, 0x79, 0x79, 0xd5, 0xa5, 0xa9, 0xd2, 0xac, 0xb5, 0xea, 0x53, 0xa5, 0x59,
	0x6f, 0x29, 0xe3, 0xeb, 0xb7, 0x1d, 0x90, 0xb7, 0x3b, 0x20, 0x7f, 0xec, 0x80, 0xfc, 0xbc, 0x07,
	0xd2, 0x76, 0x0f, 0xa4, 0xf7, 0x3d, 0x90, 0xe6, 0x57, 0x4b, 0x12, 0xad, 0xd6, 0x8e, 0xe9, 0x32,
	0x0f, 0x5e, 0xb8, 0x93, 0x78, 0x08, 0x37, 0xd5, 0xb1, 0x44, 0x49, 0x80, 0xb9, 0xd3, 0x10, 0xcf,
	0x3a, 0xfc, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x67, 0x4f, 0x85, 0x3f, 0x5b, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LivenessJailBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.LivenessJailBlocks))
		i--
		dAtA[i] = 0x30
	}
	if m.LivenessSlashInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.LivenessSlashInterval))
		i--
		dAtA[i] = 0x28
	}
	if m.LivenessSlashBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.LivenessSlashBlocks))
		i--
		dAtA[i] = 0x20
	}
	if m.DisputePeriodInBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.DisputePeriodInBlocks))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DisputePeriodInBlocks != 0 {
		n += 1 + sovParams(uint64(m.DisputePeriodInBlocks))
	}
	if m.LivenessSlashBlocks != 0 {
		n += 1 + sovParams(uint64(m.LivenessSlashBlocks))
	}
	if m.LivenessSlashInterval != 0 {
		n += 1 + sovParams(uint64(m.LivenessSlashInterval))
	}
	if m.LivenessJailBlocks != 0 {
		n += 1 + sovParams(uint64(m.LivenessJailBlocks))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DisputePeriodInBlocks", wireType)
			}
			m.DisputePeriodInBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DisputePeriodInBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LivenessSlashBlocks", wireType)
			}
			m.LivenessSlashBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LivenessSlashBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LivenessSlashInterval", wireType)
			}
			m.LivenessSlashInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LivenessSlashInterval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LivenessJailBlocks", wireType)
			}
			m.LivenessJailBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LivenessJailBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
