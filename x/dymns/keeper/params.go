package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) dymnstypes.Params {
	return dymnstypes.NewParams(
		k.PriceParams(ctx),
		k.AliasParams(ctx),
		k.MiscParams(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params dymnstypes.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	k.paramstore.SetParamSet(ctx, &params)
	return nil
}

func (k Keeper) PriceParams(ctx sdk.Context) (res dymnstypes.PriceParams) {
	k.paramstore.Get(ctx, dymnstypes.KeyPriceParams, &res)
	return
}

func (k Keeper) AliasParams(ctx sdk.Context) (res dymnstypes.AliasParams) {
	k.paramstore.Get(ctx, dymnstypes.KeyAliasParams, &res)
	return
}

func (k Keeper) MiscParams(ctx sdk.Context) (res dymnstypes.MiscParams) {
	k.paramstore.Get(ctx, dymnstypes.KeyMiscParams, &res)
	return
}
