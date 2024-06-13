// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/streamer/stream.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Stream is an object that stores and distributes yields to recipients who
// satisfy certain conditions. Currently streams support conditions around the
// duration for which a given denom is locked.
type Stream struct {
	// id is the unique ID of a Stream
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// distribute_to is the distr_info.
	DistributeTo *DistrInfo `protobuf:"bytes,2,opt,name=distribute_to,json=distributeTo,proto3" json:"distribute_to,omitempty"`
	// coins is the total amount of coins that have been in the stream
	// Can distribute multiple coin denoms
	Coins github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=coins,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"coins"`
	// start_time is the distribution start time
	StartTime time.Time `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time" yaml:"start_time"`
	// distr_epoch_identifier is what epoch type di-stribution will be triggered by
	// (day, week, etc.)
	DistrEpochIdentifier string `protobuf:"bytes,5,opt,name=distr_epoch_identifier,json=distrEpochIdentifier,proto3" json:"distr_epoch_identifier,omitempty" yaml:"distr_epoch_identifier"`
	// num_epochs_paid_over is the number of total epochs distribution will be
	// completed over
	NumEpochsPaidOver uint64 `protobuf:"varint,6,opt,name=num_epochs_paid_over,json=numEpochsPaidOver,proto3" json:"num_epochs_paid_over,omitempty"`
	// filled_epochs is the number of epochs distribution has been completed on
	// already
	FilledEpochs uint64 `protobuf:"varint,7,opt,name=filled_epochs,json=filledEpochs,proto3" json:"filled_epochs,omitempty"`
	// distributed_coins are coins that have been distributed already
	DistributedCoins github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,8,rep,name=distributed_coins,json=distributedCoins,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"distributed_coins"`
}

func (m *Stream) Reset()         { *m = Stream{} }
func (m *Stream) String() string { return proto.CompactTextString(m) }
func (*Stream) ProtoMessage()    {}
func (*Stream) Descriptor() ([]byte, []int) {
	return fileDescriptor_19586ad841c00cd9, []int{0}
}
func (m *Stream) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stream.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stream.Merge(m, src)
}
func (m *Stream) XXX_Size() int {
	return m.Size()
}
func (m *Stream) XXX_DiscardUnknown() {
	xxx_messageInfo_Stream.DiscardUnknown(m)
}

var xxx_messageInfo_Stream proto.InternalMessageInfo

func (m *Stream) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Stream) GetDistributeTo() *DistrInfo {
	if m != nil {
		return m.DistributeTo
	}
	return nil
}

func (m *Stream) GetCoins() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Coins
	}
	return nil
}

func (m *Stream) GetStartTime() time.Time {
	if m != nil {
		return m.StartTime
	}
	return time.Time{}
}

func (m *Stream) GetDistrEpochIdentifier() string {
	if m != nil {
		return m.DistrEpochIdentifier
	}
	return ""
}

func (m *Stream) GetNumEpochsPaidOver() uint64 {
	if m != nil {
		return m.NumEpochsPaidOver
	}
	return 0
}

func (m *Stream) GetFilledEpochs() uint64 {
	if m != nil {
		return m.FilledEpochs
	}
	return 0
}

func (m *Stream) GetDistributedCoins() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.DistributedCoins
	}
	return nil
}

func init() {
	proto.RegisterType((*Stream)(nil), "dymensionxyz.dymension.streamer.Stream")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/streamer/stream.proto", fileDescriptor_19586ad841c00cd9)
}

var fileDescriptor_19586ad841c00cd9 = []byte{
	// 501 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0xcf, 0x6e, 0xd3, 0x30,
	0x1c, 0x6e, 0xba, 0xae, 0x30, 0x6f, 0x43, 0x34, 0xaa, 0x50, 0xa8, 0xb4, 0xa4, 0x94, 0x4b, 0x84,
	0xc0, 0xde, 0x1f, 0x71, 0xe1, 0x58, 0xe0, 0xb0, 0x53, 0x51, 0x98, 0x04, 0xe2, 0x12, 0x39, 0xb5,
	0x93, 0x59, 0x34, 0x71, 0x14, 0x3b, 0x55, 0xcb, 0x53, 0xec, 0x39, 0x78, 0x92, 0xdd, 0xd8, 0x91,
	0x53, 0x87, 0xda, 0x37, 0xd8, 0x13, 0x20, 0xdb, 0x49, 0x5b, 0x21, 0xd0, 0x2e, 0x9c, 0x12, 0xff,
	0x7e, 0xdf, 0xf7, 0xfd, 0x7e, 0xdf, 0x67, 0x83, 0x97, 0x64, 0x9e, 0xd2, 0x4c, 0x30, 0x9e, 0xcd,
	0xe6, 0xdf, 0xd0, 0xfa, 0x80, 0x84, 0x2c, 0x28, 0x4e, 0x69, 0x51, 0xfd, 0xc0, 0xbc, 0xe0, 0x92,
	0xdb, 0xde, 0x36, 0x1a, 0xae, 0x0f, 0xb0, 0x46, 0xf7, 0xba, 0x09, 0x4f, 0xb8, 0xc6, 0x22, 0xf5,
	0x67, 0x68, 0x3d, 0x37, 0xe1, 0x3c, 0x99, 0x50, 0xa4, 0x4f, 0x51, 0x19, 0x23, 0x52, 0x16, 0x58,
	0x2a, 0xa2, 0xe9, 0x7b, 0x7f, 0xf6, 0x25, 0x4b, 0xa9, 0x90, 0x38, 0xcd, 0x6b, 0x81, 0x31, 0x17,
	0x29, 0x17, 0x28, 0xc2, 0x82, 0xa2, 0xe9, 0x49, 0x44, 0x25, 0x3e, 0x41, 0x63, 0xce, 0x6a, 0x81,
	0xe3, 0xfb, 0x5c, 0x10, 0x26, 0x64, 0x11, 0xb2, 0x2c, 0xae, 0x56, 0x1a, 0xfc, 0x68, 0x81, 0xf6,
	0x47, 0xdd, 0xb5, 0x1f, 0x81, 0x26, 0x23, 0x8e, 0xd5, 0xb7, 0xfc, 0x56, 0xd0, 0x64, 0xc4, 0x1e,
	0x81, 0x43, 0x0d, 0x67, 0x51, 0x29, 0x69, 0x28, 0xb9, 0xd3, 0xec, 0x5b, 0xfe, 0xfe, 0xe9, 0x0b,
	0x78, 0x8f, 0x79, 0xf8, 0x4e, 0xb1, 0xce, 0xb3, 0x98, 0x07, 0x07, 0x1b, 0x81, 0x0b, 0x6e, 0x63,
	0xb0, 0xab, 0x76, 0x15, 0xce, 0x4e, 0x7f, 0xc7, 0xdf, 0x3f, 0x7d, 0x0a, 0x8d, 0x1b, 0xa8, 0xdc,
	0xc0, 0xca, 0x0d, 0x7c, 0xcb, 0x59, 0x36, 0x3c, 0xbe, 0x5e, 0x78, 0x8d, 0xef, 0xb7, 0x9e, 0x9f,
	0x30, 0x79, 0x59, 0x46, 0x70, 0xcc, 0x53, 0x54, 0x59, 0x37, 0x9f, 0x57, 0x82, 0x7c, 0x45, 0x72,
	0x9e, 0x53, 0xa1, 0x09, 0x22, 0x30, 0xca, 0xf6, 0x67, 0x00, 0x84, 0xc4, 0x85, 0x0c, 0x55, 0x72,
	0x4e, 0x4b, 0x2f, 0xdc, 0x83, 0x26, 0x56, 0x58, 0xc7, 0x0a, 0x2f, 0xea, 0x58, 0x87, 0x47, 0x6a,
	0xd0, 0xdd, 0xc2, 0xeb, 0xcc, 0x71, 0x3a, 0x79, 0x33, 0xd8, 0x70, 0x07, 0x57, 0xb7, 0x9e, 0x15,
	0xec, 0xe9, 0x82, 0x82, 0xdb, 0x9f, 0xc0, 0x13, 0x13, 0x1e, 0xcd, 0xf9, 0xf8, 0x32, 0x64, 0x84,
	0x66, 0x92, 0xc5, 0x8c, 0x16, 0xce, 0x6e, 0xdf, 0xf2, 0xf7, 0x86, 0xcf, 0xee, 0x16, 0xde, 0x91,
	0x51, 0xf9, 0x3b, 0x6e, 0x10, 0x74, 0x75, 0xe3, 0xbd, 0xaa, 0x9f, 0xaf, 0xcb, 0x36, 0x02, 0xdd,
	0xac, 0x4c, 0x0d, 0x5c, 0x84, 0x39, 0x66, 0x24, 0xe4, 0x53, 0x5a, 0x38, 0x6d, 0x7d, 0x11, 0x9d,
	0xac, 0x4c, 0x35, 0x43, 0x7c, 0xc0, 0x8c, 0x8c, 0xa6, 0xb4, 0xb0, 0x9f, 0x83, 0xc3, 0x98, 0x4d,
	0x26, 0x94, 0x54, 0x1c, 0xe7, 0x81, 0x46, 0x1e, 0x98, 0xa2, 0x01, 0xdb, 0x33, 0xd0, 0xd9, 0x64,
	0x4f, 0x42, 0x93, 0xfb, 0xc3, 0xff, 0x9f, 0xfb, 0xe3, 0xad, 0x29, 0xba, 0x32, 0x1c, 0x5d, 0x2f,
	0x5d, 0xeb, 0x66, 0xe9, 0x5a, 0xbf, 0x96, 0xae, 0x75, 0xb5, 0x72, 0x1b, 0x37, 0x2b, 0xb7, 0xf1,
	0x73, 0xe5, 0x36, 0xbe, 0xbc, 0xde, 0x52, 0xfd, 0xc7, 0x43, 0x9d, 0x9e, 0xa1, 0xd9, 0xe6, 0xb5,
	0xea, 0x41, 0x51, 0x5b, 0xdf, 0xdb, 0xd9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8d, 0xe9, 0x2b,
	0xb1, 0xa3, 0x03, 0x00, 0x00,
}

func (m *Stream) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stream) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Stream) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DistributedCoins) > 0 {
		for iNdEx := len(m.DistributedCoins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DistributedCoins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintStream(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if m.FilledEpochs != 0 {
		i = encodeVarintStream(dAtA, i, uint64(m.FilledEpochs))
		i--
		dAtA[i] = 0x38
	}
	if m.NumEpochsPaidOver != 0 {
		i = encodeVarintStream(dAtA, i, uint64(m.NumEpochsPaidOver))
		i--
		dAtA[i] = 0x30
	}
	if len(m.DistrEpochIdentifier) > 0 {
		i -= len(m.DistrEpochIdentifier)
		copy(dAtA[i:], m.DistrEpochIdentifier)
		i = encodeVarintStream(dAtA, i, uint64(len(m.DistrEpochIdentifier)))
		i--
		dAtA[i] = 0x2a
	}
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.StartTime, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintStream(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if len(m.Coins) > 0 {
		for iNdEx := len(m.Coins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Coins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintStream(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.DistributeTo != nil {
		{
			size, err := m.DistributeTo.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintStream(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintStream(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintStream(dAtA []byte, offset int, v uint64) int {
	offset -= sovStream(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Stream) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovStream(uint64(m.Id))
	}
	if m.DistributeTo != nil {
		l = m.DistributeTo.Size()
		n += 1 + l + sovStream(uint64(l))
	}
	if len(m.Coins) > 0 {
		for _, e := range m.Coins {
			l = e.Size()
			n += 1 + l + sovStream(uint64(l))
		}
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.StartTime)
	n += 1 + l + sovStream(uint64(l))
	l = len(m.DistrEpochIdentifier)
	if l > 0 {
		n += 1 + l + sovStream(uint64(l))
	}
	if m.NumEpochsPaidOver != 0 {
		n += 1 + sovStream(uint64(m.NumEpochsPaidOver))
	}
	if m.FilledEpochs != 0 {
		n += 1 + sovStream(uint64(m.FilledEpochs))
	}
	if len(m.DistributedCoins) > 0 {
		for _, e := range m.DistributedCoins {
			l = e.Size()
			n += 1 + l + sovStream(uint64(l))
		}
	}
	return n
}

func sovStream(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStream(x uint64) (n int) {
	return sovStream(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Stream) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStream
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
			return fmt.Errorf("proto: Stream: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stream: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributeTo", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStream
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStream
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DistributeTo == nil {
				m.DistributeTo = &DistrInfo{}
			}
			if err := m.DistributeTo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Coins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStream
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStream
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Coins = append(m.Coins, types.Coin{})
			if err := m.Coins[len(m.Coins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStream
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStream
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.StartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistrEpochIdentifier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
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
				return ErrInvalidLengthStream
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStream
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DistrEpochIdentifier = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumEpochsPaidOver", wireType)
			}
			m.NumEpochsPaidOver = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumEpochsPaidOver |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FilledEpochs", wireType)
			}
			m.FilledEpochs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FilledEpochs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributedCoins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStream
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStream
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStream
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DistributedCoins = append(m.DistributedCoins, types.Coin{})
			if err := m.DistributedCoins[len(m.DistributedCoins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStream(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStream
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
func skipStream(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStream
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
					return 0, ErrIntOverflowStream
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
					return 0, ErrIntOverflowStream
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
				return 0, ErrInvalidLengthStream
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStream
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStream
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStream        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStream          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStream = fmt.Errorf("proto: unexpected end of group")
)
