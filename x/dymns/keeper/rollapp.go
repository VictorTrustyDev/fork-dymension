package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) IsRollAppId(ctx sdk.Context, chainId string) bool {
	_, found := k.rollappKeeper.GetRollapp(ctx, chainId)
	return found
}

func (k Keeper) GetRollAppIdByAlias(ctx sdk.Context, alias string) (rollAppId string, found bool) {
	// TODO DymNS: implement GetRollAppIdByAlias
	return "", false
}

func (k Keeper) GetRollAppBech32Prefix(ctx sdk.Context, chainId string) (bech32Prefix string, found bool) {
	// TODO DymNS: implement GetRollAppBech32Prefix
	return "", false
}
