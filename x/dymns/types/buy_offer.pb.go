// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/dymns/buy_offer.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
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

// BuyOffer defines an offer to buy a Dym-Name, placed by buyer.
// Buyer will need to deposit the offer amount to the module account.
// When the owner of the Dym-Name accepts the offer, deposited amount will be transferred to the owner.
// When the buyer cancels the offer, deposited amount will be refunded to the buyer.
type BuyOffer struct {
	// id is the unique identifier of the offer. Generated by the module.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name of the Dym-Name willing to buy.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// type is the type of the order, is Dym-Name for current version.
	Type MarketOrderType `protobuf:"varint,3,opt,name=type,proto3,enum=dymensionxyz.dymension.dymns.MarketOrderType" json:"type,omitempty"`
	// buyer is bech32 address of the account which placed the offer.
	Buyer string `protobuf:"bytes,4,opt,name=buyer,proto3" json:"buyer,omitempty"`
	// offer_price is the amount of coins that buyer willing to pay for the Dym-Name.
	// This amount is deposited to the module account upon placing the offer.
	OfferPrice types.Coin `protobuf:"bytes,5,opt,name=offer_price,json=offerPrice,proto3" json:"offer_price"`
	// counterparty_offer_price is the price that the Dym-Name owner is willing to sell the Dym-Name for.
	// This is used for counterparty price negotiation and for information only.
	// The transaction can only be executed when the owner accepts the offer with exact offer_price.
	CounterpartyOfferPrice *types.Coin `protobuf:"bytes,6,opt,name=counterparty_offer_price,json=counterpartyOfferPrice,proto3" json:"counterparty_offer_price,omitempty"`
}

func (m *BuyOffer) Reset()         { *m = BuyOffer{} }
func (m *BuyOffer) String() string { return proto.CompactTextString(m) }
func (*BuyOffer) ProtoMessage()    {}
func (*BuyOffer) Descriptor() ([]byte, []int) {
	return fileDescriptor_cce8233ba07ff78c, []int{0}
}
func (m *BuyOffer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BuyOffer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BuyOffer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BuyOffer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuyOffer.Merge(m, src)
}
func (m *BuyOffer) XXX_Size() int {
	return m.Size()
}
func (m *BuyOffer) XXX_DiscardUnknown() {
	xxx_messageInfo_BuyOffer.DiscardUnknown(m)
}

var xxx_messageInfo_BuyOffer proto.InternalMessageInfo

func (m *BuyOffer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BuyOffer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BuyOffer) GetType() MarketOrderType {
	if m != nil {
		return m.Type
	}
	return MarketOrderType_MOT_UNKNOWN
}

func (m *BuyOffer) GetBuyer() string {
	if m != nil {
		return m.Buyer
	}
	return ""
}

func (m *BuyOffer) GetOfferPrice() types.Coin {
	if m != nil {
		return m.OfferPrice
	}
	return types.Coin{}
}

func (m *BuyOffer) GetCounterpartyOfferPrice() *types.Coin {
	if m != nil {
		return m.CounterpartyOfferPrice
	}
	return nil
}

// ReverseLookupOfferIds contains a list of offer-ids for reverse lookup.
type ReverseLookupOfferIds struct {
	// offer_ids is a list of offer-ids of the Buy-Orders linked to the reverse-lookup record.
	OfferIds []string `protobuf:"bytes,1,rep,name=offer_ids,json=offerIds,proto3" json:"offer_ids,omitempty"`
}

func (m *ReverseLookupOfferIds) Reset()         { *m = ReverseLookupOfferIds{} }
func (m *ReverseLookupOfferIds) String() string { return proto.CompactTextString(m) }
func (*ReverseLookupOfferIds) ProtoMessage()    {}
func (*ReverseLookupOfferIds) Descriptor() ([]byte, []int) {
	return fileDescriptor_cce8233ba07ff78c, []int{1}
}
func (m *ReverseLookupOfferIds) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReverseLookupOfferIds) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReverseLookupOfferIds.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReverseLookupOfferIds) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReverseLookupOfferIds.Merge(m, src)
}
func (m *ReverseLookupOfferIds) XXX_Size() int {
	return m.Size()
}
func (m *ReverseLookupOfferIds) XXX_DiscardUnknown() {
	xxx_messageInfo_ReverseLookupOfferIds.DiscardUnknown(m)
}

var xxx_messageInfo_ReverseLookupOfferIds proto.InternalMessageInfo

func (m *ReverseLookupOfferIds) GetOfferIds() []string {
	if m != nil {
		return m.OfferIds
	}
	return nil
}

func init() {
	proto.RegisterType((*BuyOffer)(nil), "dymensionxyz.dymension.dymns.BuyOffer")
	proto.RegisterType((*ReverseLookupOfferIds)(nil), "dymensionxyz.dymension.dymns.ReverseLookupOfferIds")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/dymns/buy_offer.proto", fileDescriptor_cce8233ba07ff78c)
}

var fileDescriptor_cce8233ba07ff78c = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xcb, 0x6e, 0xe2, 0x30,
	0x14, 0x8d, 0x43, 0x40, 0x60, 0x24, 0x16, 0x16, 0x33, 0xca, 0x30, 0xa3, 0x4c, 0xc4, 0x2a, 0x95,
	0x5a, 0x5b, 0x40, 0x3f, 0xa0, 0xa5, 0xab, 0x4a, 0x54, 0x54, 0x69, 0x57, 0xdd, 0xa0, 0x3c, 0x0c,
	0xb5, 0x50, 0xe2, 0xc8, 0x4e, 0x10, 0xe9, 0x57, 0xf4, 0x0b, 0xfa, 0x3d, 0x2c, 0x59, 0x76, 0x55,
	0x55, 0xf0, 0x23, 0x55, 0x9c, 0x08, 0xb1, 0x29, 0xdd, 0xdd, 0xeb, 0x7b, 0xce, 0xf1, 0xbd, 0xe7,
	0xc0, 0xf3, 0x30, 0x8f, 0x68, 0x2c, 0x19, 0x8f, 0xd7, 0xf9, 0x0b, 0x39, 0x34, 0x45, 0x15, 0x4b,
	0xe2, 0x67, 0xf9, 0x8c, 0xcf, 0xe7, 0x54, 0xe0, 0x44, 0xf0, 0x94, 0xa3, 0x7f, 0xc7, 0x68, 0x7c,
	0x68, 0xb0, 0x42, 0xf7, 0xba, 0x0b, 0xbe, 0xe0, 0x0a, 0x48, 0x8a, 0xaa, 0xe4, 0xf4, 0xac, 0x80,
	0xcb, 0x88, 0x4b, 0xe2, 0x7b, 0x92, 0x92, 0xd5, 0xc0, 0xa7, 0xa9, 0x37, 0x20, 0x01, 0x67, 0x71,
	0x35, 0x3f, 0x3b, 0xb9, 0x41, 0xe4, 0x89, 0x25, 0x4d, 0x4b, 0x68, 0xff, 0x4d, 0x87, 0xcd, 0x71,
	0x96, 0x4f, 0x8b, 0x8d, 0x50, 0x07, 0xea, 0x2c, 0x34, 0x81, 0x0d, 0x9c, 0x96, 0xab, 0xb3, 0x10,
	0x21, 0x68, 0xc4, 0x5e, 0x44, 0x4d, 0x5d, 0xbd, 0xa8, 0x1a, 0x5d, 0x43, 0x23, 0xcd, 0x13, 0x6a,
	0xd6, 0x6c, 0xe0, 0x74, 0x86, 0x17, 0xf8, 0xd4, 0xfa, 0xf8, 0x4e, 0x7d, 0x35, 0x15, 0x21, 0x15,
	0x8f, 0x79, 0x42, 0x5d, 0x45, 0x45, 0x5d, 0x58, 0xf7, 0xb3, 0x9c, 0x0a, 0xd3, 0x50, 0xba, 0x65,
	0x83, 0xae, 0x60, 0x5b, 0xf9, 0x32, 0x4b, 0x04, 0x0b, 0xa8, 0x59, 0xb7, 0x81, 0xd3, 0x1e, 0xfe,
	0xc1, 0xe5, 0xa9, 0xb8, 0x38, 0x15, 0x57, 0xa7, 0xe2, 0x1b, 0xce, 0xe2, 0xb1, 0xb1, 0xf9, 0xf8,
	0xaf, 0xb9, 0x50, 0x71, 0xee, 0x0b, 0x0a, 0x7a, 0x80, 0x66, 0xc0, 0xb3, 0x38, 0xa5, 0x22, 0xf1,
	0x44, 0x5a, 0xd9, 0x5c, 0xc9, 0x35, 0x7e, 0x90, 0x73, 0x7f, 0x1f, 0x53, 0xa7, 0x07, 0xd1, 0xfe,
	0x25, 0xfc, 0xe5, 0xd2, 0x15, 0x15, 0x92, 0x4e, 0x38, 0x5f, 0x66, 0x89, 0x1a, 0xdd, 0x86, 0x12,
	0xfd, 0x85, 0xad, 0xf2, 0x03, 0x16, 0x4a, 0x13, 0xd8, 0x35, 0xa7, 0xe5, 0x36, 0x79, 0x35, 0x1c,
	0x4f, 0x36, 0x3b, 0x0b, 0x6c, 0x77, 0x16, 0xf8, 0xdc, 0x59, 0xe0, 0x75, 0x6f, 0x69, 0xdb, 0xbd,
	0xa5, 0xbd, 0xef, 0x2d, 0xed, 0x69, 0xb8, 0x60, 0xe9, 0x73, 0xe6, 0xe3, 0x80, 0x47, 0xe4, 0x9b,
	0x98, 0x56, 0x23, 0xb2, 0xae, 0xb2, 0x2a, 0xfc, 0x92, 0x7e, 0x43, 0x65, 0x35, 0xfa, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0xf4, 0xea, 0xf7, 0xcc, 0x5a, 0x02, 0x00, 0x00,
}

func (m *BuyOffer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BuyOffer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BuyOffer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CounterpartyOfferPrice != nil {
		{
			size, err := m.CounterpartyOfferPrice.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBuyOffer(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	{
		size, err := m.OfferPrice.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintBuyOffer(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.Buyer) > 0 {
		i -= len(m.Buyer)
		copy(dAtA[i:], m.Buyer)
		i = encodeVarintBuyOffer(dAtA, i, uint64(len(m.Buyer)))
		i--
		dAtA[i] = 0x22
	}
	if m.Type != 0 {
		i = encodeVarintBuyOffer(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintBuyOffer(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintBuyOffer(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ReverseLookupOfferIds) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReverseLookupOfferIds) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReverseLookupOfferIds) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OfferIds) > 0 {
		for iNdEx := len(m.OfferIds) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.OfferIds[iNdEx])
			copy(dAtA[i:], m.OfferIds[iNdEx])
			i = encodeVarintBuyOffer(dAtA, i, uint64(len(m.OfferIds[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintBuyOffer(dAtA []byte, offset int, v uint64) int {
	offset -= sovBuyOffer(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BuyOffer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovBuyOffer(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovBuyOffer(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovBuyOffer(uint64(m.Type))
	}
	l = len(m.Buyer)
	if l > 0 {
		n += 1 + l + sovBuyOffer(uint64(l))
	}
	l = m.OfferPrice.Size()
	n += 1 + l + sovBuyOffer(uint64(l))
	if m.CounterpartyOfferPrice != nil {
		l = m.CounterpartyOfferPrice.Size()
		n += 1 + l + sovBuyOffer(uint64(l))
	}
	return n
}

func (m *ReverseLookupOfferIds) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.OfferIds) > 0 {
		for _, s := range m.OfferIds {
			l = len(s)
			n += 1 + l + sovBuyOffer(uint64(l))
		}
	}
	return n
}

func sovBuyOffer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozBuyOffer(x uint64) (n int) {
	return sovBuyOffer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BuyOffer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBuyOffer
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
			return fmt.Errorf("proto: BuyOffer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BuyOffer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= MarketOrderType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Buyer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Buyer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfferPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OfferPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CounterpartyOfferPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CounterpartyOfferPrice == nil {
				m.CounterpartyOfferPrice = &types.Coin{}
			}
			if err := m.CounterpartyOfferPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBuyOffer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBuyOffer
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
func (m *ReverseLookupOfferIds) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBuyOffer
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
			return fmt.Errorf("proto: ReverseLookupOfferIds: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReverseLookupOfferIds: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfferIds", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBuyOffer
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
				return ErrInvalidLengthBuyOffer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthBuyOffer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OfferIds = append(m.OfferIds, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBuyOffer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthBuyOffer
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
func skipBuyOffer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowBuyOffer
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
					return 0, ErrIntOverflowBuyOffer
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
					return 0, ErrIntOverflowBuyOffer
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
				return 0, ErrInvalidLengthBuyOffer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupBuyOffer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthBuyOffer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthBuyOffer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowBuyOffer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupBuyOffer = fmt.Errorf("proto: unexpected end of group")
)