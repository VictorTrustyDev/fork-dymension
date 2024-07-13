package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k msgServer) CancelAdsSellName(goCtx context.Context, msg *dymnstypes.MsgCancelAdsSellName) (*dymnstypes.MsgCancelAdsSellNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateCancelAdsSellName(ctx, msg); err != nil {
		return nil, err
	}

	k.DeleteOpenPurchaseOrder(ctx, msg.Name)

	apoe := k.GetActiveOpenPurchaseOrdersExpiration(ctx)
	delete(apoe.ExpiryByName, msg.Name)
	if err := k.SetActiveOpenPurchaseOrdersExpiration(ctx, apoe); err != nil {
		return nil, err
	}

	minimumTxGas := dymnstypes.OpGasCloseAds
	if consumedGas := ctx.GasMeter().GasConsumed(); consumedGas < minimumTxGas {
		ctx.GasMeter().ConsumeGas(minimumTxGas-consumedGas, "CancelAdsSellName")
	}

	return &dymnstypes.MsgCancelAdsSellNameResponse{}, nil
}

func (k msgServer) validateCancelAdsSellName(ctx sdk.Context, msg *dymnstypes.MsgCancelAdsSellName) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return dymnstypes.ErrDymNameNotFound.Wrap(msg.Name)
	}

	if dymName.Owner != msg.Owner {
		return sdkerrors.ErrUnauthorized.Wrap("not the owner of the dym name")
	}

	opo := k.GetOpenPurchaseOrder(ctx, msg.Name)
	if opo == nil {
		return dymnstypes.ErrOpenPurchaseOrderNotFound.Wrap(msg.Name)
	}

	if opo.HasExpiredAtCtx(ctx) {
		return dymnstypes.ErrInvalidState.Wrap("cannot cancel an expired order")
	}

	if opo.HighestBid != nil {
		return dymnstypes.ErrInvalidState.Wrap("cannot cancel once bid placed")
	}

	return nil
}
