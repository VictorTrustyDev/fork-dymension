package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"time"
)

func (k msgServer) PutAdsSellName(goCtx context.Context, msg *dymnstypes.MsgPutAdsSellName) (*dymnstypes.MsgPutAdsSellNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName, params, err := k.validatePutAdsSellName(ctx, msg)
	if err != nil {
		return nil, err
	}

	opo := msg.ToOpenPurchaseOrder()
	opo.ExpireAt = ctx.BlockTime().Add(
		24 * time.Hour * time.Duration(params.Misc.DaysOpenPurchaseOrderDuration),
	).Unix()

	if err := opo.Validate(); err != nil {
		panic(errors.Wrap(err, "un-expected invalid state of created OPO"))
	}

	prohibitSellingAfterEpoch := time.Unix(dymName.ExpireAt, 0).Add(
		-1 * (24 * time.Hour * time.Duration(params.Misc.DaysProhibitSell)),
	).Unix()

	if opo.ExpireAt > prohibitSellingAfterEpoch {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf(
			"%d days before Dym-Name expiry, can not sell",
			params.Misc.DaysProhibitSell,
		)
	}

	if err := k.SetOpenPurchaseOrder(ctx, opo); err != nil {
		return nil, err
	}

	apoe := k.GetActiveOpenPurchaseOrdersExpiration(ctx)
	apoe.ExpiryByName[opo.Name] = opo.ExpireAt
	if err := k.SetActiveOpenPurchaseOrdersExpiration(ctx, apoe); err != nil {
		return nil, err
	}

	minimumTxGas := sdk.Gas(params.Misc.GasCrudOpenPurchaseOrder)
	if consumedGas := ctx.GasMeter().GasConsumed(); consumedGas < minimumTxGas {
		ctx.GasMeter().ConsumeGas(minimumTxGas-consumedGas, "PutAdsSellName")
	}

	return &dymnstypes.MsgPutAdsSellNameResponse{}, nil
}

func (k msgServer) validatePutAdsSellName(ctx sdk.Context, msg *dymnstypes.MsgPutAdsSellName) (*dymnstypes.DymName, *dymnstypes.Params, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return nil, nil, dymnstypes.ErrDymNameNotFound.Wrap(msg.Name)
	}

	if dymName.Owner != msg.Owner {
		return nil, nil, sdkerrors.ErrUnauthorized
	}

	if dymName.IsExpiredAt(ctx.BlockTime()) {
		return nil, nil, sdkerrors.ErrUnauthorized.Wrap("Dym-Name is already expired")
	}

	existingActiveOpo := k.GetOpenPurchaseOrder(ctx, dymName.Name)
	if existingActiveOpo != nil {
		if existingActiveOpo.HasFinishedAtCtx(ctx) {
			return nil, nil, sdkerrors.ErrConflict.Wrap(
				"an active expired/completed Open-Purchase-Order already exists for the Dym-Name, must wait until processed",
			)
		}
		return nil, nil, sdkerrors.ErrConflict.Wrap("an active Open-Purchase-Order already exists for the Dym-Name")
	}

	params := k.GetParams(ctx)

	if msg.MinPrice.Denom != params.Price.PriceDenom {
		return nil, nil, sdkerrors.ErrInvalidRequest.Wrapf(
			"only %s is allowed as price",
			params.Price.PriceDenom,
		)
	}

	return dymName, &params, nil
}
