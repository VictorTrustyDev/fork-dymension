package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

// HasCounterpartyOfferPrice returns true if the offer has a raise-offer request from the Dym-Name owner.
func (m *OfferToBuy) HasCounterpartyOfferPrice() bool {
	return m.CounterpartyOfferPrice != nil && !m.CounterpartyOfferPrice.Amount.IsNil() && !m.CounterpartyOfferPrice.IsZero()
}

// Validate performs basic validation for the OfferToBuy.
func (m *OfferToBuy) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("offer is nil")
	}

	if m.Id == "" {
		return ErrValidationFailed.Wrap("ID of offer is empty")
	}

	if !dymnsutils.IsValidBuyNameOfferId(m.Id) {
		return ErrValidationFailed.Wrap("ID of offer is not a valid offer id")
	}

	if m.Name == "" {
		return ErrValidationFailed.Wrap("Dym-Name of offer is empty")
	}

	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("Dym-Name of offer is not a valid dym name")
	}

	if !dymnsutils.IsValidBech32AccountAddress(m.Buyer, true) {
		return ErrValidationFailed.Wrap("buyer is not a valid bech32 account address")
	}

	if m.OfferPrice.Amount.IsNil() || m.OfferPrice.IsZero() {
		return ErrValidationFailed.Wrap("offer price is zero")
	} else if m.OfferPrice.IsNegative() {
		return ErrValidationFailed.Wrap("offer price is negative")
	} else if err := m.OfferPrice.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("offer price is invalid: %v", err)
	}

	if m.HasCounterpartyOfferPrice() {
		if m.CounterpartyOfferPrice.IsNegative() {
			return ErrValidationFailed.Wrap("counterparty offer price is negative")
		} else if err := m.CounterpartyOfferPrice.Validate(); err != nil {
			return ErrValidationFailed.Wrapf("counterparty offer price is invalid: %v", err)
		}

		if m.CounterpartyOfferPrice.Denom != m.OfferPrice.Denom {
			return ErrValidationFailed.Wrap("counterparty offer price denom is different from offer price denom")
		}
	}

	return nil
}

// GetSdkEvent returns the sdk event contains information of Offer-To-Buy record.
// Fired when Offer-To-Buy record is set into store.
func (m OfferToBuy) GetSdkEvent(actionName string) sdk.Event {
	var attrCounterpartyOfferPrice sdk.Attribute
	if m.CounterpartyOfferPrice != nil {
		attrCounterpartyOfferPrice = sdk.NewAttribute(AttributeKeyOtbCounterpartyOfferPrice, m.CounterpartyOfferPrice.String())
	} else {
		attrCounterpartyOfferPrice = sdk.NewAttribute(AttributeKeyOtbCounterpartyOfferPrice, "")
	}

	return sdk.NewEvent(
		EventTypeOfferToBuy,
		sdk.NewAttribute(AttributeKeyOtbId, m.Id),
		sdk.NewAttribute(AttributeKeyOtbName, m.Name),
		sdk.NewAttribute(AttributeKeyOtbOfferPrice, m.OfferPrice.String()),
		attrCounterpartyOfferPrice,
		sdk.NewAttribute(AttributeKeyOtbActionName, actionName),
	)
}