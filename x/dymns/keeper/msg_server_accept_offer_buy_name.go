package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// AcceptOfferBuyName is message handler,
// handles accepting an Offer-To-Buy or raising the amount for negotiation, performed by the owner of the Dym-Name.
func (k msgServer) AcceptOfferBuyName(goCtx context.Context, msg *dymnstypes.MsgAcceptOfferBuyName) (*dymnstypes.MsgAcceptOfferBuyNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	offer, dymName, err := k.validateAcceptOffer(ctx, msg)
	if err != nil {
		return nil, err
	}

	var accepted bool

	if msg.MinAccept.IsLT(offer.OfferPrice) {
		panic("min-accept is less than offer price")
	} else if msg.MinAccept.IsEqual(offer.OfferPrice) {
		accepted = true

		// take the offer
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			dymnstypes.ModuleName,
			sdk.MustAccAddressFromBech32(dymName.Owner),
			sdk.Coins{offer.OfferPrice},
		); err != nil {
			return nil, err
		}

		if err := k.removeOfferToBuy(ctx, *offer); err != nil {
			return nil, err
		}

		if err := k.transferOwnership(ctx, *dymName, offer.Buyer); err != nil {
			return nil, err
		}
	} else {
		accepted = false

		offer.CounterpartyOfferPrice = &msg.MinAccept
		if err := k.SetOfferToBuy(ctx, *offer); err != nil {
			return nil, err
		}
	}

	consumeMinimumGas(ctx, dymnstypes.OpGasUpdateOffer, "AcceptOfferBuyName")

	return &dymnstypes.MsgAcceptOfferBuyNameResponse{
		Accepted: accepted,
	}, nil
}

// validateAcceptOffer handles validation for the message handled by AcceptOfferBuyName
func (k msgServer) validateAcceptOffer(ctx sdk.Context, msg *dymnstypes.MsgAcceptOfferBuyName) (*dymnstypes.OfferToBuy, *dymnstypes.DymName, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, nil, err
	}

	offer := k.GetOfferToBuy(ctx, msg.OfferId)
	if offer == nil {
		return nil, nil, errorsmod.Wrapf(gerrc.ErrNotFound, "Offer-To-Buy: %s", msg.OfferId)
	}

	dymName := k.GetDymNameWithExpirationCheck(ctx, offer.Name)
	if dymName == nil {
		return nil, nil, errorsmod.Wrapf(gerrc.ErrNotFound, "Dym-Name: %s", offer.Name)
	}

	if dymName.Owner != msg.Owner {
		return nil, nil, errorsmod.Wrapf(gerrc.ErrPermissionDenied, "not the owner of the Dym-Name")
	}

	params := k.GetParams(ctx)

	if dymName.IsProhibitedTradingAt(ctx.BlockTime(), params.Misc.ProhibitSellDuration) {
		return nil, nil, errorsmod.Wrapf(gerrc.ErrFailedPrecondition,
			"duration before Dym-Name expiry, prohibited to sell: %s",
			params.Misc.ProhibitSellDuration,
		)
	}

	if offer.Buyer == msg.Owner {
		return nil, nil, errorsmod.Wrapf(gerrc.ErrPermissionDenied, "cannot accept own offer")
	}

	if msg.MinAccept.Denom != offer.OfferPrice.Denom {
		return nil, nil, errorsmod.Wrapf(
			gerrc.ErrInvalidArgument,
			"denom must be the same as the offer price: %s", offer.OfferPrice.Denom,
		)
	}

	if msg.MinAccept.IsLT(offer.OfferPrice) {
		return nil, nil, errorsmod.Wrapf(
			gerrc.ErrInvalidArgument,
			"amount must be greater than or equals to the offer price: %s", offer.OfferPrice.Denom,
		)
	}

	return offer, dymName, nil
}
