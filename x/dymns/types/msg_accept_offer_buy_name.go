package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

var _ sdk.Msg = &MsgAcceptOfferBuyName{}

func (m *MsgAcceptOfferBuyName) ValidateBasic() error {
	if !dymnsutils.IsValidBuyNameOfferId(m.OfferId) {
		return ErrValidationFailed.Wrap("offer id is not a valid buy name offer id")
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrValidationFailed.Wrap("owner is not a valid bech32 account address")
	}

	if !m.MinAccept.IsValid() {
		return ErrValidationFailed.Wrap("invalid min-accept amount")
	} else if !m.MinAccept.IsPositive() {
		return ErrValidationFailed.Wrap("min-accept amount must be positive")
	}

	return nil
}

func (m *MsgAcceptOfferBuyName) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (m *MsgAcceptOfferBuyName) Route() string {
	return RouterKey
}

func (m *MsgAcceptOfferBuyName) Type() string {
	return TypeMsgAcceptOfferBuyName
}

func (m *MsgAcceptOfferBuyName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}
