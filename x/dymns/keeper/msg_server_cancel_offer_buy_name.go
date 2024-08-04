package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// CancelOfferBuyName is message handler,
// handles canceling a Buy-Offer, performed by the buyer who placed the offer.
func (k msgServer) CancelOfferBuyName(goCtx context.Context, msg *dymnstypes.MsgCancelOfferBuyName) (*dymnstypes.MsgCancelOfferBuyNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	offer, err := k.validateCancelOffer(ctx, msg)
	if err != nil {
		return nil, err
	}

	if err := k.RefundOffer(ctx, *offer); err != nil {
		return nil, err
	}

	if err := k.removeBuyOffer(ctx, *offer); err != nil {
		return nil, err
	}

	consumeMinimumGas(ctx, dymnstypes.OpGasCloseBuyOffer, "CancelOfferBuyName")

	return &dymnstypes.MsgCancelOfferBuyNameResponse{}, nil
}

// validateCancelOffer handles validation for the message handled by CancelOfferBuyName.
func (k msgServer) validateCancelOffer(ctx sdk.Context, msg *dymnstypes.MsgCancelOfferBuyName) (*dymnstypes.BuyOffer, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	offer := k.GetBuyOffer(ctx, msg.OfferId)
	if offer == nil {
		return nil, errorsmod.Wrapf(gerrc.ErrNotFound, "Buy-Offer ID: %s", msg.OfferId)
	}

	if offer.Buyer != msg.Buyer {
		return nil, errorsmod.Wrap(gerrc.ErrPermissionDenied, "not the owner of the offer")
	}

	return offer, nil
}

// removeBuyOffer removes the Buy-Offer from the store and the reverse mappings.
func (k msgServer) removeBuyOffer(ctx sdk.Context, offer dymnstypes.BuyOffer) error {
	k.DeleteBuyOffer(ctx, offer.Id)

	err := k.RemoveReverseMappingBuyerToBuyOffer(ctx, offer.Buyer, offer.Id)
	if err != nil {
		return err
	}

	err = k.RemoveReverseMappingDymNameToBuyOffer(ctx, offer.Name, offer.Id)
	if err != nil {
		return err
	}

	return nil
}
