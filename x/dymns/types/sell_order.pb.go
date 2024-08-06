// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dymensionxyz/dymension/dymns/sell_order.proto

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

// SellOrder defines a sell order, placed by owner, to sell a Dym-Name/Alias.
// Sell-Order has an expiry date.
// After expiry date, if no one has placed a bid, this Sell-Order will be closed, no change.
// If there is a bid, the highest bid will win, and the Dym-Name/Alias ownership will be transferred to the winner.
// If the bid matches the sell price, the Dym-Name/Alias ownership will be transferred to the bidder immediately.
// SO will be moved to historical records after completion/expiry.
type SellOrder struct {
	// goods_id is the Dym-Name/Alias being opened to be sold.
	GoodsId string `protobuf:"bytes,1,opt,name=goods_id,json=goodsId,proto3" json:"goods_id,omitempty"`
	// type of the order, is Dym-Name or Alias.
	Type OrderType `protobuf:"varint,2,opt,name=type,proto3,enum=dymensionxyz.dymension.dymns.OrderType" json:"type,omitempty"`
	// expire_at is the last effective date of this SO
	ExpireAt int64 `protobuf:"varint,3,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
	// min_price is the minimum price that the owner is willing to accept for the goods.
	MinPrice types.Coin `protobuf:"bytes,4,opt,name=min_price,json=minPrice,proto3" json:"min_price"`
	// sell_price is the price that the owner is willing to sell the Dym-Name/Alias for,
	// the SO will be closed when the price is met.
	// If the sell price is zero, the SO will be closed when the expire_at is reached and the highest bidder wins.
	SellPrice *types.Coin `protobuf:"bytes,5,opt,name=sell_price,json=sellPrice,proto3" json:"sell_price,omitempty"`
	// highest_bid is the highest bid on the SO, if any. Price must be greater than or equal to the min_price.
	HighestBid *SellOrderBid `protobuf:"bytes,6,opt,name=highest_bid,json=highestBid,proto3" json:"highest_bid,omitempty"`
}

func (m *SellOrder) Reset()         { *m = SellOrder{} }
func (m *SellOrder) String() string { return proto.CompactTextString(m) }
func (*SellOrder) ProtoMessage()    {}
func (*SellOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6763b7e2ceeacb, []int{0}
}
func (m *SellOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SellOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SellOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SellOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SellOrder.Merge(m, src)
}
func (m *SellOrder) XXX_Size() int {
	return m.Size()
}
func (m *SellOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_SellOrder.DiscardUnknown(m)
}

var xxx_messageInfo_SellOrder proto.InternalMessageInfo

func (m *SellOrder) GetGoodsId() string {
	if m != nil {
		return m.GoodsId
	}
	return ""
}

func (m *SellOrder) GetType() OrderType {
	if m != nil {
		return m.Type
	}
	return OrderType_OT_UNKNOWN
}

func (m *SellOrder) GetExpireAt() int64 {
	if m != nil {
		return m.ExpireAt
	}
	return 0
}

func (m *SellOrder) GetMinPrice() types.Coin {
	if m != nil {
		return m.MinPrice
	}
	return types.Coin{}
}

func (m *SellOrder) GetSellPrice() *types.Coin {
	if m != nil {
		return m.SellPrice
	}
	return nil
}

func (m *SellOrder) GetHighestBid() *SellOrderBid {
	if m != nil {
		return m.HighestBid
	}
	return nil
}

// ActiveSellOrdersExpiration contains list of active SOs, store expiration date mapped by goods identity.
// Used by hook to find out expired SO instead of iterating through all records.
type ActiveSellOrdersExpiration struct {
	Records []ActiveSellOrdersExpirationRecord `protobuf:"bytes,1,rep,name=records,proto3" json:"records"`
}

func (m *ActiveSellOrdersExpiration) Reset()         { *m = ActiveSellOrdersExpiration{} }
func (m *ActiveSellOrdersExpiration) String() string { return proto.CompactTextString(m) }
func (*ActiveSellOrdersExpiration) ProtoMessage()    {}
func (*ActiveSellOrdersExpiration) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6763b7e2ceeacb, []int{1}
}
func (m *ActiveSellOrdersExpiration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ActiveSellOrdersExpiration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ActiveSellOrdersExpiration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ActiveSellOrdersExpiration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ActiveSellOrdersExpiration.Merge(m, src)
}
func (m *ActiveSellOrdersExpiration) XXX_Size() int {
	return m.Size()
}
func (m *ActiveSellOrdersExpiration) XXX_DiscardUnknown() {
	xxx_messageInfo_ActiveSellOrdersExpiration.DiscardUnknown(m)
}

var xxx_messageInfo_ActiveSellOrdersExpiration proto.InternalMessageInfo

func (m *ActiveSellOrdersExpiration) GetRecords() []ActiveSellOrdersExpirationRecord {
	if m != nil {
		return m.Records
	}
	return nil
}

// ActiveSellOrdersExpirationRecord contains the expiration date of an active Sell-Order.
type ActiveSellOrdersExpirationRecord struct {
	// goods_id is the Dym-Name/Alias being opened to be sold.
	GoodsId string `protobuf:"bytes,1,opt,name=goods_id,json=goodsId,proto3" json:"goods_id,omitempty"`
	// expire_at is the last effective date of this Sell-Order.
	ExpireAt int64 `protobuf:"varint,2,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty"`
}

func (m *ActiveSellOrdersExpirationRecord) Reset()         { *m = ActiveSellOrdersExpirationRecord{} }
func (m *ActiveSellOrdersExpirationRecord) String() string { return proto.CompactTextString(m) }
func (*ActiveSellOrdersExpirationRecord) ProtoMessage()    {}
func (*ActiveSellOrdersExpirationRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6763b7e2ceeacb, []int{2}
}
func (m *ActiveSellOrdersExpirationRecord) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ActiveSellOrdersExpirationRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ActiveSellOrdersExpirationRecord.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ActiveSellOrdersExpirationRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ActiveSellOrdersExpirationRecord.Merge(m, src)
}
func (m *ActiveSellOrdersExpirationRecord) XXX_Size() int {
	return m.Size()
}
func (m *ActiveSellOrdersExpirationRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_ActiveSellOrdersExpirationRecord.DiscardUnknown(m)
}

var xxx_messageInfo_ActiveSellOrdersExpirationRecord proto.InternalMessageInfo

func (m *ActiveSellOrdersExpirationRecord) GetGoodsId() string {
	if m != nil {
		return m.GoodsId
	}
	return ""
}

func (m *ActiveSellOrdersExpirationRecord) GetExpireAt() int64 {
	if m != nil {
		return m.ExpireAt
	}
	return 0
}

// SellOrderBid defines a bid placed by an account on a Sell-Order.
type SellOrderBid struct {
	// bidder is the account address of the account which placed the bid.
	Bidder string `protobuf:"bytes,1,opt,name=bidder,proto3" json:"bidder,omitempty"`
	// price is the amount of coin offered by the bidder.
	Price types.Coin `protobuf:"bytes,2,opt,name=price,proto3" json:"price"`
}

func (m *SellOrderBid) Reset()         { *m = SellOrderBid{} }
func (m *SellOrderBid) String() string { return proto.CompactTextString(m) }
func (*SellOrderBid) ProtoMessage()    {}
func (*SellOrderBid) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6763b7e2ceeacb, []int{3}
}
func (m *SellOrderBid) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SellOrderBid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SellOrderBid.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SellOrderBid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SellOrderBid.Merge(m, src)
}
func (m *SellOrderBid) XXX_Size() int {
	return m.Size()
}
func (m *SellOrderBid) XXX_DiscardUnknown() {
	xxx_messageInfo_SellOrderBid.DiscardUnknown(m)
}

var xxx_messageInfo_SellOrderBid proto.InternalMessageInfo

func (m *SellOrderBid) GetBidder() string {
	if m != nil {
		return m.Bidder
	}
	return ""
}

func (m *SellOrderBid) GetPrice() types.Coin {
	if m != nil {
		return m.Price
	}
	return types.Coin{}
}

// HistoricalSellOrders contains list of closed Sell-Orders of the same goods.
type HistoricalSellOrders struct {
	// sell_orders is list of closed Sell-Orders of the same goods.
	SellOrders []SellOrder `protobuf:"bytes,1,rep,name=sell_orders,json=sellOrders,proto3" json:"sell_orders"`
}

func (m *HistoricalSellOrders) Reset()         { *m = HistoricalSellOrders{} }
func (m *HistoricalSellOrders) String() string { return proto.CompactTextString(m) }
func (*HistoricalSellOrders) ProtoMessage()    {}
func (*HistoricalSellOrders) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6763b7e2ceeacb, []int{4}
}
func (m *HistoricalSellOrders) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HistoricalSellOrders) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HistoricalSellOrders.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HistoricalSellOrders) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HistoricalSellOrders.Merge(m, src)
}
func (m *HistoricalSellOrders) XXX_Size() int {
	return m.Size()
}
func (m *HistoricalSellOrders) XXX_DiscardUnknown() {
	xxx_messageInfo_HistoricalSellOrders.DiscardUnknown(m)
}

var xxx_messageInfo_HistoricalSellOrders proto.InternalMessageInfo

func (m *HistoricalSellOrders) GetSellOrders() []SellOrder {
	if m != nil {
		return m.SellOrders
	}
	return nil
}

func init() {
	proto.RegisterType((*SellOrder)(nil), "dymensionxyz.dymension.dymns.SellOrder")
	proto.RegisterType((*ActiveSellOrdersExpiration)(nil), "dymensionxyz.dymension.dymns.ActiveSellOrdersExpiration")
	proto.RegisterType((*ActiveSellOrdersExpirationRecord)(nil), "dymensionxyz.dymension.dymns.ActiveSellOrdersExpirationRecord")
	proto.RegisterType((*SellOrderBid)(nil), "dymensionxyz.dymension.dymns.SellOrderBid")
	proto.RegisterType((*HistoricalSellOrders)(nil), "dymensionxyz.dymension.dymns.HistoricalSellOrders")
}

func init() {
	proto.RegisterFile("dymensionxyz/dymension/dymns/sell_order.proto", fileDescriptor_1b6763b7e2ceeacb)
}

var fileDescriptor_1b6763b7e2ceeacb = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x41, 0x8b, 0xd3, 0x40,
	0x18, 0xed, 0xb4, 0xdd, 0x6e, 0x33, 0x15, 0x0f, 0xc3, 0x22, 0xd9, 0x2a, 0x31, 0xe4, 0xb2, 0x51,
	0x70, 0xc2, 0x76, 0x11, 0x04, 0x45, 0xd8, 0x8a, 0xa0, 0x28, 0x2a, 0xd1, 0xd3, 0x82, 0x86, 0x24,
	0x33, 0xa6, 0x83, 0x49, 0x26, 0xcc, 0x8c, 0xa5, 0x15, 0x7f, 0x84, 0x7f, 0x4a, 0xd8, 0xe3, 0x1e,
	0x3d, 0x89, 0xb4, 0x7f, 0x44, 0x32, 0x93, 0xed, 0xee, 0x1e, 0x36, 0xdb, 0xdb, 0x7c, 0xc9, 0x7b,
	0x6f, 0xbe, 0xf7, 0xbd, 0xf9, 0xe0, 0x23, 0xb2, 0x2c, 0x68, 0x29, 0x19, 0x2f, 0x17, 0xcb, 0x1f,
	0xc1, 0xa6, 0xa8, 0x4f, 0xa5, 0x0c, 0x24, 0xcd, 0xf3, 0x88, 0x0b, 0x42, 0x05, 0xae, 0x04, 0x57,
	0x1c, 0xdd, 0xbb, 0x0c, 0xc7, 0x9b, 0x02, 0x6b, 0xf8, 0x78, 0x2f, 0xe3, 0x19, 0xd7, 0xc0, 0xa0,
	0x3e, 0x19, 0xce, 0xd8, 0x49, 0xb9, 0x2c, 0xb8, 0x0c, 0x92, 0x58, 0xd2, 0x60, 0x7e, 0x98, 0x50,
	0x15, 0x1f, 0x06, 0x29, 0x67, 0x65, 0xf3, 0xff, 0x41, 0x6b, 0x0b, 0x45, 0x2c, 0xbe, 0x51, 0x65,
	0xa0, 0xde, 0xef, 0x2e, 0xb4, 0x3e, 0xd2, 0x3c, 0x7f, 0x5f, 0xb7, 0x84, 0xf6, 0xe1, 0x30, 0xe3,
	0x9c, 0xc8, 0x88, 0x11, 0x1b, 0xb8, 0xc0, 0xb7, 0xc2, 0x5d, 0x5d, 0xbf, 0x26, 0xe8, 0x29, 0xec,
	0xab, 0x65, 0x45, 0xed, 0xae, 0x0b, 0xfc, 0xdb, 0x93, 0x03, 0xdc, 0xd6, 0x36, 0xd6, 0x6a, 0x9f,
	0x96, 0x15, 0x0d, 0x35, 0x09, 0xdd, 0x85, 0x16, 0x5d, 0x54, 0x4c, 0xd0, 0x28, 0x56, 0x76, 0xcf,
	0x05, 0x7e, 0x2f, 0x1c, 0x9a, 0x0f, 0xc7, 0x0a, 0x3d, 0x83, 0x56, 0xc1, 0xca, 0xa8, 0x12, 0x2c,
	0xa5, 0x76, 0xdf, 0x05, 0xfe, 0x68, 0xb2, 0x8f, 0x8d, 0x43, 0x5c, 0x3b, 0xc4, 0x8d, 0x43, 0xfc,
	0x82, 0xb3, 0x72, 0xda, 0x3f, 0xfd, 0x7b, 0xbf, 0x13, 0x0e, 0x0b, 0x56, 0x7e, 0xa8, 0x09, 0xe8,
	0x09, 0x84, 0x7a, 0xa6, 0x86, 0xbe, 0x73, 0x03, 0x3d, 0xb4, 0x6a, 0xb0, 0x61, 0xbe, 0x81, 0xa3,
	0x19, 0xcb, 0x66, 0x54, 0xaa, 0x28, 0x61, 0xc4, 0x1e, 0x68, 0xea, 0xc3, 0x76, 0x63, 0x9b, 0x51,
	0x4d, 0x19, 0x09, 0x61, 0x43, 0x9f, 0x32, 0xe2, 0xfd, 0x84, 0xe3, 0xe3, 0x54, 0xb1, 0x39, 0xdd,
	0x20, 0xe4, 0xcb, 0xda, 0x60, 0xac, 0x18, 0x2f, 0xd1, 0x17, 0xb8, 0x2b, 0x68, 0xca, 0x05, 0x91,
	0x36, 0x70, 0x7b, 0xfe, 0x68, 0xf2, 0xbc, 0xfd, 0x9a, 0xeb, 0xa5, 0x42, 0x2d, 0xd3, 0x4c, 0xe1,
	0x5c, 0xd4, 0x3b, 0x81, 0xee, 0x4d, 0x94, 0xb6, 0x6c, 0xaf, 0xc4, 0xd3, 0xbd, 0x1a, 0x8f, 0xf7,
	0x19, 0xde, 0xba, 0xec, 0x1a, 0xdd, 0x81, 0x83, 0x84, 0x11, 0x42, 0x45, 0xa3, 0xd2, 0x54, 0xe8,
	0x31, 0xdc, 0x31, 0x19, 0x74, 0xb7, 0x8b, 0xd0, 0xa0, 0xbd, 0xaf, 0x70, 0xef, 0x15, 0x93, 0x8a,
	0x0b, 0x96, 0xc6, 0xf9, 0x45, 0xfb, 0xe8, 0x1d, 0x1c, 0x5d, 0xec, 0xca, 0xf9, 0xd8, 0x0e, 0xb6,
	0x4d, 0xc7, 0x5c, 0xa1, 0x5f, 0x86, 0xd1, 0x9b, 0xbe, 0x3d, 0x5d, 0x39, 0xe0, 0x6c, 0xe5, 0x80,
	0x7f, 0x2b, 0x07, 0xfc, 0x5a, 0x3b, 0x9d, 0xb3, 0xb5, 0xd3, 0xf9, 0xb3, 0x76, 0x3a, 0x27, 0x93,
	0x8c, 0xa9, 0xd9, 0xf7, 0x04, 0xa7, 0xbc, 0x08, 0xae, 0x59, 0x9c, 0xf9, 0x51, 0xb0, 0x68, 0xb6,
	0xa7, 0x7e, 0xcf, 0x32, 0x19, 0xe8, 0xed, 0x39, 0xfa, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x9e,
	0xdb, 0xc4, 0xed, 0x03, 0x00, 0x00,
}

func (m *SellOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SellOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SellOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.HighestBid != nil {
		{
			size, err := m.HighestBid.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSellOrder(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.SellPrice != nil {
		{
			size, err := m.SellPrice.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSellOrder(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	{
		size, err := m.MinPrice.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSellOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.ExpireAt != 0 {
		i = encodeVarintSellOrder(dAtA, i, uint64(m.ExpireAt))
		i--
		dAtA[i] = 0x18
	}
	if m.Type != 0 {
		i = encodeVarintSellOrder(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if len(m.GoodsId) > 0 {
		i -= len(m.GoodsId)
		copy(dAtA[i:], m.GoodsId)
		i = encodeVarintSellOrder(dAtA, i, uint64(len(m.GoodsId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ActiveSellOrdersExpiration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ActiveSellOrdersExpiration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ActiveSellOrdersExpiration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Records) > 0 {
		for iNdEx := len(m.Records) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Records[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSellOrder(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ActiveSellOrdersExpirationRecord) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ActiveSellOrdersExpirationRecord) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ActiveSellOrdersExpirationRecord) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ExpireAt != 0 {
		i = encodeVarintSellOrder(dAtA, i, uint64(m.ExpireAt))
		i--
		dAtA[i] = 0x10
	}
	if len(m.GoodsId) > 0 {
		i -= len(m.GoodsId)
		copy(dAtA[i:], m.GoodsId)
		i = encodeVarintSellOrder(dAtA, i, uint64(len(m.GoodsId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SellOrderBid) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SellOrderBid) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SellOrderBid) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Price.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintSellOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Bidder) > 0 {
		i -= len(m.Bidder)
		copy(dAtA[i:], m.Bidder)
		i = encodeVarintSellOrder(dAtA, i, uint64(len(m.Bidder)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *HistoricalSellOrders) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HistoricalSellOrders) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HistoricalSellOrders) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SellOrders) > 0 {
		for iNdEx := len(m.SellOrders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SellOrders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSellOrder(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintSellOrder(dAtA []byte, offset int, v uint64) int {
	offset -= sovSellOrder(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SellOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GoodsId)
	if l > 0 {
		n += 1 + l + sovSellOrder(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovSellOrder(uint64(m.Type))
	}
	if m.ExpireAt != 0 {
		n += 1 + sovSellOrder(uint64(m.ExpireAt))
	}
	l = m.MinPrice.Size()
	n += 1 + l + sovSellOrder(uint64(l))
	if m.SellPrice != nil {
		l = m.SellPrice.Size()
		n += 1 + l + sovSellOrder(uint64(l))
	}
	if m.HighestBid != nil {
		l = m.HighestBid.Size()
		n += 1 + l + sovSellOrder(uint64(l))
	}
	return n
}

func (m *ActiveSellOrdersExpiration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Records) > 0 {
		for _, e := range m.Records {
			l = e.Size()
			n += 1 + l + sovSellOrder(uint64(l))
		}
	}
	return n
}

func (m *ActiveSellOrdersExpirationRecord) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GoodsId)
	if l > 0 {
		n += 1 + l + sovSellOrder(uint64(l))
	}
	if m.ExpireAt != 0 {
		n += 1 + sovSellOrder(uint64(m.ExpireAt))
	}
	return n
}

func (m *SellOrderBid) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Bidder)
	if l > 0 {
		n += 1 + l + sovSellOrder(uint64(l))
	}
	l = m.Price.Size()
	n += 1 + l + sovSellOrder(uint64(l))
	return n
}

func (m *HistoricalSellOrders) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SellOrders) > 0 {
		for _, e := range m.SellOrders {
			l = e.Size()
			n += 1 + l + sovSellOrder(uint64(l))
		}
	}
	return n
}

func sovSellOrder(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSellOrder(x uint64) (n int) {
	return sovSellOrder(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SellOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSellOrder
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
			return fmt.Errorf("proto: SellOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SellOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GoodsId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GoodsId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= OrderType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpireAt", wireType)
			}
			m.ExpireAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpireAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SellPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SellPrice == nil {
				m.SellPrice = &types.Coin{}
			}
			if err := m.SellPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HighestBid", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HighestBid == nil {
				m.HighestBid = &SellOrderBid{}
			}
			if err := m.HighestBid.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSellOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSellOrder
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
func (m *ActiveSellOrdersExpiration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSellOrder
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
			return fmt.Errorf("proto: ActiveSellOrdersExpiration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ActiveSellOrdersExpiration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Records", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Records = append(m.Records, ActiveSellOrdersExpirationRecord{})
			if err := m.Records[len(m.Records)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSellOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSellOrder
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
func (m *ActiveSellOrdersExpirationRecord) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSellOrder
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
			return fmt.Errorf("proto: ActiveSellOrdersExpirationRecord: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ActiveSellOrdersExpirationRecord: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GoodsId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GoodsId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpireAt", wireType)
			}
			m.ExpireAt = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpireAt |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSellOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSellOrder
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
func (m *SellOrderBid) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSellOrder
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
			return fmt.Errorf("proto: SellOrderBid: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SellOrderBid: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bidder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bidder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSellOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSellOrder
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
func (m *HistoricalSellOrders) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSellOrder
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
			return fmt.Errorf("proto: HistoricalSellOrders: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HistoricalSellOrders: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SellOrders", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSellOrder
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
				return ErrInvalidLengthSellOrder
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSellOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SellOrders = append(m.SellOrders, SellOrder{})
			if err := m.SellOrders[len(m.SellOrders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSellOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSellOrder
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
func skipSellOrder(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSellOrder
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
					return 0, ErrIntOverflowSellOrder
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
					return 0, ErrIntOverflowSellOrder
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
				return 0, ErrInvalidLengthSellOrder
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSellOrder
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSellOrder
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSellOrder        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSellOrder          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSellOrder = fmt.Errorf("proto: unexpected end of group")
)
