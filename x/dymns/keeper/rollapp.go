package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) IsRollAppId(ctx sdk.Context, chainId string) bool {
	_, found := k.rollappKeeper.GetRollapp(ctx, chainId)
	return found
}

func (k Keeper) GetRollAppIdByAlias(ctx sdk.Context, alias string) (rollAppId string, found bool) {
	// TODO DymNS: implement Get RollApp-Id By Alias
	return "", false
}

func (k Keeper) GetAliasByRollAppId(ctx sdk.Context, chainId string) (alias string, found bool) {
	_, exists := k.rollappKeeper.GetRollapp(ctx, chainId)
	if !exists {
		return "", false
	}

	// TODO DymNS: implement Get Alias by RollApp-Id
	return "", false
}

func (k Keeper) GetRollAppBech32Prefix(ctx sdk.Context, chainId string) (bech32Prefix string, found bool) {
	// TODO DymNS: implement Get RollApp Bech32 Prefix
	return "", false
}
