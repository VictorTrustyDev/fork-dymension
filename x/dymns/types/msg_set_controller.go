package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

var (
	_ sdk.Msg = &MsgSetController{}
)

func (m *MsgSetController) ValidateBasic() error {
	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("name is not a valid dym name")
	}

	if _, err := sdk.AccAddressFromBech32(m.Controller); err != nil {
		return ErrValidationFailed.Wrap("controller is not a valid bech32 account address")
	}

	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return ErrValidationFailed.Wrap("owner is not a valid bech32 account address")
	}

	return nil
}

func (m *MsgSetController) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

func (m *MsgSetController) Route() string {
	return RouterKey
}

func (m *MsgSetController) Type() string {
	return TypeSetController
}

func (m *MsgSetController) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}