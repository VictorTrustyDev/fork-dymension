package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

// GetIdentity returns the unique identity of the SO
func (m *SellOrder) GetIdentity() string {
	return fmt.Sprintf("%s|%d", m.Name, m.ExpireAt)
}

// HasSetSellPrice returns true if the sell price is set
func (m *SellOrder) HasSetSellPrice() bool {
	return m.SellPrice != nil && !m.SellPrice.Amount.IsNil() && !m.SellPrice.IsZero()
}

// HasExpiredAtCtx returns true if the SO has expired at given context
func (m *SellOrder) HasExpiredAtCtx(ctx sdk.Context) bool {
	return m.HasExpired(ctx.BlockTime().Unix())
}

// HasExpired returns true if the SO has expired at given epoch
func (m *SellOrder) HasExpired(nowEpoch int64) bool {
	return m.ExpireAt < nowEpoch
}

// HasFinishedAtCtx returns true if the SO has expired or completed at given context
func (m *SellOrder) HasFinishedAtCtx(ctx sdk.Context) bool {
	return m.HasFinished(ctx.BlockTime().Unix())
}

// HasFinished returns true if the SO has expired or completed at given epoch
func (m *SellOrder) HasFinished(nowEpoch int64) bool {
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

func (m *SellOrder) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("SO is nil")
	}

	if m.Name == "" {
		return ErrValidationFailed.Wrap("Dym-Name of SO is empty")
	}

	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("Dym-Name of SO is not a valid dym name")
	}

	if m.ExpireAt == 0 {
		return ErrValidationFailed.Wrap("SO expiry is empty")
	}

	if m.MinPrice.Amount.IsNil() || m.MinPrice.IsZero() {
		return ErrValidationFailed.Wrap("SO min price is zero")
	} else if m.MinPrice.IsNegative() {
		return ErrValidationFailed.Wrap("SO min price is negative")
	} else if err := m.MinPrice.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("SO min price is invalid: %v", err)
	}

	if m.HasSetSellPrice() {
		if m.SellPrice.IsNegative() {
			return ErrValidationFailed.Wrap("SO sell price is negative")
		} else if err := m.SellPrice.Validate(); err != nil {
			return ErrValidationFailed.Wrapf("SO sell price is invalid: %v", err)
		}

		if m.SellPrice.Denom != m.MinPrice.Denom {
			return ErrValidationFailed.Wrap("SO sell price denom is different from min price denom")
		}

		if m.SellPrice.IsLT(m.MinPrice) {
			return ErrValidationFailed.Wrap("SO sell price is less than min price")
		}
	} else {
		// allowed
	}

	if m.HighestBid == nil {
		// valid, means no bid yet
	} else if err := m.HighestBid.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("SO highest bid is invalid: %v", err)
	} else if m.HighestBid.Price.IsLT(m.MinPrice) {
		return ErrValidationFailed.Wrap("SO highest bid price is less than min price")
	} else if m.HasSetSellPrice() && m.SellPrice.IsLT(m.HighestBid.Price) {
		return ErrValidationFailed.Wrap("SO sell price is less than highest bid price")
	}

	return nil
}

func (m *SellOrderBid) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("SO bid is nil")
	}

	if m.Bidder == "" {
		return ErrValidationFailed.Wrap("SO bidder is empty")
	}

	if !dymnsutils.IsValidBech32AccountAddress(m.Bidder, true) {
		return ErrValidationFailed.Wrap("SO bidder is not a valid bech32 account address")
	}

	if m.Price.Amount.IsNil() || m.Price.IsZero() {
		return ErrValidationFailed.Wrap("SO bid price is zero")
	} else if m.Price.IsNegative() {
		return ErrValidationFailed.Wrap("SO bid price is negative")
	} else if err := m.Price.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("SO bid price is invalid: %v", err)
	}

	return nil
}

func (m *HistoricalSellOrders) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("historical SOs is nil")
	}

	if len(m.SellOrders) > 0 {
		name := m.SellOrders[0].Name
		var uniqueIdentity = make(map[string]bool)
		for _, so := range m.SellOrders {
			if err := so.Validate(); err != nil {
				return err
			}

			if so.Name != name {
				return ErrValidationFailed.Wrapf("historical SOs have different Dym-Name, expected only %s but got %s", name, so.Name)
			}

			if _, duplicated := uniqueIdentity[so.GetIdentity()]; duplicated {
				return ErrValidationFailed.Wrapf("historical SO is not unique: %s", so.GetIdentity())
			}
			uniqueIdentity[so.GetIdentity()] = true
		}
	}

	return nil
}

func (m SellOrder) GetSdkEvent(actionName string) sdk.Event {
	var sellPrice sdk.Coin
	if m.HasSetSellPrice() {
		sellPrice = *m.SellPrice
	} else {
		sellPrice = sdk.NewCoin(m.MinPrice.Denom, sdk.ZeroInt())
	}

	var attrHighestBidder, attrHighestBidPrice sdk.Attribute
	if m.HighestBid != nil {
		attrHighestBidder = sdk.NewAttribute(AttributeKeyDymNameSoHighestBidder, m.HighestBid.Bidder)
		attrHighestBidPrice = sdk.NewAttribute(AttributeKeyDymNameSoHighestBidPrice, m.HighestBid.Price.String())
	} else {
		attrHighestBidder = sdk.NewAttribute(AttributeKeyDymNameSoHighestBidder, "")
		attrHighestBidPrice = sdk.NewAttribute(AttributeKeyDymNameSoHighestBidPrice, sdk.NewCoin(m.MinPrice.Denom, sdk.ZeroInt()).String())
	}

	return sdk.NewEvent(
		EventTypeDymNameSellOrder,
		sdk.NewAttribute(AttributeKeyDymNameSoName, m.Name),
		sdk.NewAttribute(AttributeKeyDymNameSoExpiryEpoch, fmt.Sprintf("%d", m.ExpireAt)),
		sdk.NewAttribute(AttributeKeyDymNameSoMinPrice, m.MinPrice.String()),
		sdk.NewAttribute(AttributeKeyDymNameSoSellPrice, sellPrice.String()),
		attrHighestBidder,
		attrHighestBidPrice,
		sdk.NewAttribute(AttributeKeyDymNameSoActionName, actionName),
	)
}
