package keeper

import (
	"context"
	"time"

	"github.com/dymensionxyz/gerr-cosmos/gerrc"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// PlaceSellOrder is message handler,
// handles creating a Sell-Order that advertise a Dym-Name/Alias is for sale, performed by the owner.
func (k msgServer) PlaceSellOrder(goCtx context.Context, msg *dymnstypes.MsgPlaceSellOrder) (*dymnstypes.MsgPlaceSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)

	var resp *dymnstypes.MsgPlaceSellOrderResponse
	var err error

	if msg.AssetType == dymnstypes.TypeName {
		resp, err = k.processPlaceSellOrderWithAssetTypeDymName(ctx, msg, params)
	} else if msg.AssetType == dymnstypes.TypeAlias {
		resp, err = k.processPlaceSellOrderWithAssetTypeAlias(ctx, msg, params)
	} else {
		err = errorsmod.Wrapf(gerrc.ErrInvalidArgument, "invalid asset type: %s", msg.AssetType)
	}
	if err != nil {
		return nil, err
	}

	consumeMinimumGas(ctx, dymnstypes.OpGasPlaceSellOrder, "PlaceSellOrder")

	return resp, nil
}

// processPlaceSellOrderWithAssetTypeDymName handles the message handled by PlaceSellOrder, type Dym-Name.
func (k msgServer) processPlaceSellOrderWithAssetTypeDymName(
	ctx sdk.Context,
	msg *dymnstypes.MsgPlaceSellOrder, params dymnstypes.Params,
) (*dymnstypes.MsgPlaceSellOrderResponse, error) {
	if !params.Misc.EnableTradingName {
		return nil, errorsmod.Wrapf(gerrc.ErrFailedPrecondition, "trading of Dym-Name is disabled")
	}

	dymName, err := k.validatePlaceSellOrderWithAssetTypeDymName(ctx, msg, params)
	if err != nil {
		return nil, err
	}

	so := msg.ToSellOrder()
	so.ExpireAt = ctx.BlockTime().Add(params.Misc.SellOrderDuration).Unix()

	if err := so.Validate(); err != nil {
		panic(errorsmod.Wrap(err, "un-expected invalid state of created SO"))
	}

	if dymName.IsProhibitedTradingAt(time.Unix(so.ExpireAt, 0), params.Misc.ProhibitSellDuration) {
		return nil, errorsmod.Wrapf(gerrc.ErrFailedPrecondition,
			"duration before Dym-Name expiry, prohibited to sell: %s",
			params.Misc.ProhibitSellDuration,
		)
	}

	if err := k.SetSellOrder(ctx, so); err != nil {
		return nil, err
	}

	aSoe := k.GetActiveSellOrdersExpiration(ctx, so.AssetType)
	aSoe.Add(so.AssetId, so.ExpireAt)
	if err := k.SetActiveSellOrdersExpiration(ctx, aSoe, so.AssetType); err != nil {
		return nil, err
	}

	return &dymnstypes.MsgPlaceSellOrderResponse{}, nil
}

// validatePlaceSellOrderWithAssetTypeDymName handles validation for message handled by PlaceSellOrder, type Dym-Name.
func (k msgServer) validatePlaceSellOrderWithAssetTypeDymName(
	ctx sdk.Context,
	msg *dymnstypes.MsgPlaceSellOrder, params dymnstypes.Params,
) (*dymnstypes.DymName, error) {
	dymName := k.GetDymName(ctx, msg.AssetId)
	if dymName == nil {
		return nil, errorsmod.Wrapf(gerrc.ErrNotFound, "Dym-Name: %s", msg.AssetId)
	}

	if dymName.Owner != msg.Owner {
		return nil, errorsmod.Wrap(gerrc.ErrPermissionDenied, "not the owner of the Dym-Name")
	}

	if dymName.IsExpiredAtCtx(ctx) {
		return nil, errorsmod.Wrap(gerrc.ErrUnauthenticated, "Dym-Name is already expired")
	}

	existingActiveSo := k.GetSellOrder(ctx, dymName.Name, msg.AssetType)
	if existingActiveSo != nil {
		if existingActiveSo.HasFinishedAtCtx(ctx) {
			return nil, errorsmod.Wrap(
				gerrc.ErrAlreadyExists,
				"an active expired/completed Sell-Order already exists for the Dym-Name, must wait until processed",
			)
		}
		return nil, errorsmod.Wrap(gerrc.ErrAlreadyExists, "an active Sell-Order already exists for the Dym-Name")
	}

	if msg.MinPrice.Denom != params.Price.PriceDenom {
		return nil, errorsmod.Wrapf(
			gerrc.ErrInvalidArgument,
			"the only denom allowed as price: %s", params.Price.PriceDenom,
		)
	}

	return dymName, nil
}

// processPlaceSellOrderWithAssetTypeAlias handles the message handled by PlaceSellOrder, type Alias.
func (k msgServer) processPlaceSellOrderWithAssetTypeAlias(
	ctx sdk.Context,
	msg *dymnstypes.MsgPlaceSellOrder, params dymnstypes.Params,
) (*dymnstypes.MsgPlaceSellOrderResponse, error) {
	if !params.Misc.EnableTradingAlias {
		return nil, errorsmod.Wrapf(gerrc.ErrFailedPrecondition, "trading of Alias is disabled")
	}

	err := k.validatePlaceSellOrderWithAssetTypeAlias(ctx, msg, params)
	if err != nil {
		return nil, err
	}

	so := msg.ToSellOrder()
	so.ExpireAt = ctx.BlockTime().Add(params.Misc.SellOrderDuration).Unix()

	if err := so.Validate(); err != nil {
		panic(errorsmod.Wrap(err, "un-expected invalid state of created SO"))
	}

	if err := k.SetSellOrder(ctx, so); err != nil {
		return nil, err
	}

	aSoe := k.GetActiveSellOrdersExpiration(ctx, so.AssetType)
	aSoe.Add(so.AssetId, so.ExpireAt)
	if err := k.SetActiveSellOrdersExpiration(ctx, aSoe, so.AssetType); err != nil {
		return nil, err
	}

	return &dymnstypes.MsgPlaceSellOrderResponse{}, nil
}

// validatePlaceSellOrderWithAssetTypeAlias handles validation for message handled by PlaceSellOrder, type Alias.
func (k msgServer) validatePlaceSellOrderWithAssetTypeAlias(
	ctx sdk.Context,
	msg *dymnstypes.MsgPlaceSellOrder, params dymnstypes.Params,
) error {
	alias := msg.AssetId

	if k.IsAliasPresentsInParamsAsAliasOrChainId(ctx, msg.AssetId) {
		return errorsmod.Wrapf(gerrc.ErrPermissionDenied,
			"prohibited to trade aliases which is reserved for chain-id or alias in module params: %s", msg.AssetId,
		)
	}

	sourceRollAppId, found := k.GetRollAppIdByAlias(ctx, alias)
	if !found {
		return errorsmod.Wrapf(gerrc.ErrNotFound, "alias: %s", alias)
	}

	if !k.IsRollAppCreator(ctx, sourceRollAppId, msg.Owner) {
		return errorsmod.Wrapf(gerrc.ErrPermissionDenied,
			"not the owner of the RollApp using the alias: %s", sourceRollAppId,
		)
	}

	existingActiveSo := k.GetSellOrder(ctx, alias, msg.AssetType)
	if existingActiveSo != nil {
		if existingActiveSo.HasFinishedAtCtx(ctx) {
			return errorsmod.Wrap(
				gerrc.ErrAlreadyExists,
				"an active expired/completed Sell-Order already exists for the Alias, must wait until processed",
			)
		}
		return errorsmod.Wrap(gerrc.ErrAlreadyExists, "an active Sell-Order already exists for the Alias")
	}

	if msg.MinPrice.Denom != params.Price.PriceDenom {
		return errorsmod.Wrapf(
			gerrc.ErrInvalidArgument,
			"the only denom allowed as price: %s", params.Price.PriceDenom,
		)
	}

	return nil
}
