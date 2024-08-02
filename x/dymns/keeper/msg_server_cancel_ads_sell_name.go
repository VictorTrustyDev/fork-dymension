package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// CancelAdsSellName is message handler,
// handles canceling Sell-Order, performed by the owner.
// This will stop the advertisement and remove the Dym-Name from the market.
// Can only be performed if the Dym-Name is not in any offer.
func (k msgServer) CancelAdsSellName(goCtx context.Context, msg *dymnstypes.MsgCancelAdsSellName) (*dymnstypes.MsgCancelAdsSellNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateCancelAdsSellName(ctx, msg); err != nil {
		return nil, err
	}

	k.DeleteSellOrder(ctx, msg.Name)

	aSoe := k.GetActiveSellOrdersExpiration(ctx)
	aSoe.Remove(msg.Name)
	if err := k.SetActiveSellOrdersExpiration(ctx, aSoe); err != nil {
		return nil, err
	}

	consumeMinimumGas(ctx, dymnstypes.OpGasCloseAds, "CancelAdsSellName")

	return &dymnstypes.MsgCancelAdsSellNameResponse{}, nil
}

// validateCancelAdsSellName handles validation for the message handled by CancelAdsSellName.
func (k msgServer) validateCancelAdsSellName(ctx sdk.Context, msg *dymnstypes.MsgCancelAdsSellName) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return errorsmod.Wrapf(gerrc.ErrNotFound, "Dym-Name: %s", msg.Name)
	}

	if dymName.Owner != msg.Owner {
		return errorsmod.Wrap(gerrc.ErrPermissionDenied, "not the owner of the Dym-Name")
	}

	so := k.GetSellOrder(ctx, msg.Name)
	if so == nil {
		return errorsmod.Wrapf(gerrc.ErrNotFound, "Sell-Order: %s", msg.Name)
	}

	if so.HasExpiredAtCtx(ctx) {
		return errorsmod.Wrap(gerrc.ErrFailedPrecondition, "cannot cancel an expired order")
	}

	if so.HighestBid != nil {
		return errorsmod.Wrap(gerrc.ErrFailedPrecondition, "cannot cancel once bid placed")
	}

	return nil
}
