package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

var (
	_ sdk.Msg = &MsgPurchaseName{}
)

func (m *MsgPurchaseName) ValidateBasic() error {
	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("name is not a valid dym name")
	}

	if !m.Offer.IsValid() {
		return ErrValidationFailed.Wrap("invalid offer")
	} else if !m.Offer.IsPositive() {
		return ErrValidationFailed.Wrap("offer must be positive")
	}

	if _, err := sdk.AccAddressFromBech32(m.Buyer); err != nil {
		return ErrValidationFailed.Wrap("buyer is not a valid bech32 account address")
	}

	return nil
}

func (m *MsgPurchaseName) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(m.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (m *MsgPurchaseName) Route() string {
	return RouterKey
}

func (m *MsgPurchaseName) Type() string {
	return TypePurchaseName
}

func (m *MsgPurchaseName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}
