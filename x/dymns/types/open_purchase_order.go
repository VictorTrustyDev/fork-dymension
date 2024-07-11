package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

// Identity returns the unique identity of the OPO
func (m *OpenPurchaseOrder) Identity() string {
	return fmt.Sprintf("%s|%d", m.Name, m.ExpireAt)
}

// HasSetSellPrice returns true if the sell price is set
func (m *OpenPurchaseOrder) HasSetSellPrice() bool {
	return m.SellPrice != nil && !m.SellPrice.Amount.IsNil() && !m.SellPrice.IsZero()
}

// HasExpiredAtCtx returns true if the OPO has expired at given context
func (m *OpenPurchaseOrder) HasExpiredAtCtx(ctx sdk.Context) bool {
	return m.HasExpired(ctx.BlockTime().Unix())
}

// HasExpired returns true if the OPO has expired at given epoch
func (m *OpenPurchaseOrder) HasExpired(nowEpoch int64) bool {
	return m.ExpireAt < nowEpoch
}

// HasFinishedAtCtx returns true if the OPO has expired or completed at given context
func (m *OpenPurchaseOrder) HasFinishedAtCtx(ctx sdk.Context) bool {
	return m.HasFinished(ctx.BlockTime().Unix())
}

// HasFinished returns true if the OPO has expired or completed at given epoch
func (m *OpenPurchaseOrder) HasFinished(nowEpoch int64) bool {
	if m.HasExpired(nowEpoch) {
		return true
	}

	if !m.HasSetSellPrice() {
		// when no sell price is set, must wait until completed auction
		return false
	}

	// complete condition: bid >= sell price

	if m.HighestBid == nil {
		return false
	}

	return m.HighestBid.Price.IsGTE(*m.SellPrice)
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

		if m.SellPrice.Denom != m.MinPrice.Denom {
			return ErrValidationFailed.Wrap("OPO sell price denom is different from min price denom")
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

func (m *HistoricalOpenPurchaseOrders) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("historical OPOs is nil")
	}

	if len(m.OpenPurchaseOrders) > 0 {
		name := m.OpenPurchaseOrders[0].Name
		var uniqueIdentity = make(map[string]bool)
		for _, opo := range m.OpenPurchaseOrders {
			if err := opo.Validate(); err != nil {
				return err
			}

			if opo.Name != name {
				return ErrValidationFailed.Wrapf("historical OPOs have different Dym-Name, expected only %s but got %s", name, opo.Name)
			}

			if _, duplicated := uniqueIdentity[opo.Identity()]; duplicated {
				return ErrValidationFailed.Wrapf("historical OPO is not unique: %s", opo.Identity())
			}
			uniqueIdentity[opo.Identity()] = true
		}
	}

	return nil
}

func (m OpenPurchaseOrder) GetSdkEvent(actionName string) sdk.Event {
	var sellPrice sdk.Coin
	if m.HasSetSellPrice() {
		sellPrice = *m.SellPrice
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
		sdk.NewAttribute(AttributeKeyDymNameOpoActionName, actionName),
	)
}
