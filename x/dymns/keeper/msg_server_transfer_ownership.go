package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k msgServer) TransferOwnership(goCtx context.Context, msg *dymnstypes.MsgTransferOwnership) (*dymnstypes.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName, err := k.validateTransferOwnership(ctx, msg)
	if err != nil {
		return nil, err
	}

	if err := k.PruneDymName(ctx, dymName.Name); err != nil {
		return nil, err
	}

	newDymNameRecord := dymnstypes.DymName{
		Name:       dymName.Name,
		Owner:      msg.NewOwner,     // transfer ownership
		Controller: msg.NewOwner,     // transfer controller
		ExpireAt:   dymName.ExpireAt, // keep the same expiration date
		Configs:    nil,              // clear configs
	}

	if err := k.SetDymName(ctx, newDymNameRecord); err != nil {
		return nil, err
	}

	return &dymnstypes.MsgTransferOwnershipResponse{}, nil
}

func (k msgServer) validateTransferOwnership(ctx sdk.Context, msg *dymnstypes.MsgTransferOwnership) (*dymnstypes.DymName, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return nil, dymnstypes.ErrDymNameNotFound.Wrap(msg.Name)
	}

	if dymName.Owner != msg.Owner {
		return nil, sdkerrors.ErrUnauthorized
	}

	if dymName.IsExpiredAt(ctx.BlockTime()) {
		return nil, sdkerrors.ErrUnauthorized.Wrap("Dym-Name is already expired")
	}

	opo := k.GetOpenPurchaseOrder(ctx, msg.Name)
	if opo != nil {
		// by ignoring OPO, can fall into case that OPO not completed/lost funds of bidder,...

		return nil, sdkerrors.ErrInvalidRequest.Wrap("can not transfer ownership while there is an Open Purchase Order")
	}

	return dymName, nil
}
