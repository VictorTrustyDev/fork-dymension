package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

var (
	_ sdk.Msg = &MsgPutAdsSellName{}
)

func (m *MsgPutAdsSellName) ValidateBasic() error {
	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("name is not a valid dym name")
	}

	opo := m.ToOpenPurchaseOrder()

	// put a dummy expire at to validate, as zero expire at is invalid,
	// and we don't have context of time at this point
	opo.ExpireAt = 1

	if err := opo.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("invalid order: %v", err)
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrValidationFailed.Wrap("owner is not a valid bech32 account address")
	}

	return nil
}

func (m *MsgPutAdsSellName) ToOpenPurchaseOrder() OpenPurchaseOrder {
	opo := OpenPurchaseOrder{
		Name:      m.Name,
		MinPrice:  m.MinPrice,
		SellPrice: m.SellPrice,
	}

	if !opo.HasSetSellPrice() {
		opo.SellPrice = nil
	}

	return opo
}

func (m *MsgPutAdsSellName) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (m *MsgPutAdsSellName) Route() string {
	return RouterKey
}

func (m *MsgPutAdsSellName) Type() string {
	return TypePutAdsSellName
}

func (m *MsgPutAdsSellName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}
