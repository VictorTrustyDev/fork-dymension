package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

func (m *OpenPurchaseOrder) HasSetSellPrice() bool {
	return !m.SellPrice.Amount.IsNil() && !m.SellPrice.IsZero()
}

func (m *OpenPurchaseOrder) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("OPO is nil")
	}

	if m.Name == "" {
		return ErrValidationFailed.Wrap("Dym-Name of OPO is empty")
	}

	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("Dym-Name of OPO is not a valid dym name")
	}

	if m.ExpireAt == 0 {
		return ErrValidationFailed.Wrap("OPO expiry is empty")
	}

	if m.MinPrice.Amount.IsNil() || m.MinPrice.IsZero() {
		return ErrValidationFailed.Wrap("OPO min price is zero")
	} else if m.MinPrice.IsNegative() {
		return ErrValidationFailed.Wrap("OPO min price is negative")
	} else if err := m.MinPrice.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("OPO min price is invalid: %v", err)
	}

	if m.HasSetSellPrice() {
		if m.SellPrice.IsNegative() {
			return ErrValidationFailed.Wrap("OPO sell price is negative")
		} else if err := m.SellPrice.Validate(); err != nil {
			return ErrValidationFailed.Wrapf("OPO sell price is invalid: %v", err)
		}

		if m.SellPrice.IsLT(m.MinPrice) {
			return ErrValidationFailed.Wrap("OPO sell price is less than min price")
		}
	} else {
		// allowed
	}

	if m.HighestBid == nil {
		// valid, means no bid yet
	} else if err := m.HighestBid.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("OPO highest bid is invalid: %v", err)
	} else if m.HighestBid.Price.IsLT(m.MinPrice) {
		return ErrValidationFailed.Wrap("OPO highest bid price is less than min price")
	} else if m.HasSetSellPrice() && m.SellPrice.IsLT(m.HighestBid.Price) {
		return ErrValidationFailed.Wrap("OPO sell price is less than highest bid price")
	}

	return nil
}

func (m *OpenPurchaseOrderBid) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("OPO bid is nil")
	}

	if m.Bidder == "" {
		return ErrValidationFailed.Wrap("OPO bidder is empty")
	}

	if !dymnsutils.IsValidBech32AccountAddress(m.Bidder, true) {
		return ErrValidationFailed.Wrap("OPO bidder is not a valid bech32 account address")
	}

	if m.Price.Amount.IsNil() || m.Price.IsZero() {
		return ErrValidationFailed.Wrap("OPO bid price is zero")
	} else if m.Price.IsNegative() {
		return ErrValidationFailed.Wrap("OPO bid price is negative")
	} else if err := m.Price.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("OPO bid price is invalid: %v", err)
	}

	return nil
}

func (m OpenPurchaseOrder) GetSdkEvent() sdk.Event {
	var sellPrice sdk.Coin
	if m.HasSetSellPrice() {
		sellPrice = m.SellPrice
	} else {
		sellPrice = sdk.NewCoin(m.MinPrice.Denom, sdk.ZeroInt())
	}

	var attrHighestBidder, attrHighestBidPrice sdk.Attribute
	if m.HighestBid != nil {
		attrHighestBidder = sdk.NewAttribute(AttributeKeyDymNameOpoHighestBidder, m.HighestBid.Bidder)
		attrHighestBidPrice = sdk.NewAttribute(AttributeKeyDymNameOpoHighestBidPrice, m.HighestBid.Price.String())
	} else {
		attrHighestBidder = sdk.NewAttribute(AttributeKeyDymNameOpoHighestBidder, "")
		attrHighestBidPrice = sdk.NewAttribute(AttributeKeyDymNameOpoHighestBidPrice, sdk.NewCoin(m.MinPrice.Denom, sdk.ZeroInt()).String())
	}

	return sdk.NewEvent(
		EventTypeDymNameOpenPurchaseOrder,
		sdk.NewAttribute(AttributeKeyDymNameOpoName, m.Name),
		sdk.NewAttribute(AttributeKeyDymNameOpoExpiryEpoch, fmt.Sprintf("%d", m.ExpireAt)),
		sdk.NewAttribute(AttributeKeyDymNameOpoMinPrice, m.MinPrice.String()),
		sdk.NewAttribute(AttributeKeyDymNameOpoSellPrice, sellPrice.String()),
		attrHighestBidder,
		attrHighestBidPrice,
	)
}
