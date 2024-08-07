package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO DymNS: remove this mock
type mockRollAppData struct {
	alias  string
	bech32 string
}

var mockRollAppsData = map[string]mockRollAppData{
	"nim_1122-1": {
		alias:  "nim",
		bech32: "nim",
	},
	"mande_18071918-1": {
		alias:  "mande",
		bech32: "mande",
	},
}

// end of mock

// IsRollAppId checks if the chain-id is a RollApp-Id.
func (k Keeper) IsRollAppId(ctx sdk.Context, chainId string) bool {
	_, found := k.rollappKeeper.GetRollapp(ctx, chainId)

	if !found {
		_, found = mockRollAppsData[chainId]
	}

	return found
}

// IsRollAppCreator returns true if the input bech32 address is the creator of the RollApp.
func (k Keeper) IsRollAppCreator(ctx sdk.Context, rollAppId, account string) bool {
	rollApp, found := k.rollappKeeper.GetRollapp(ctx, rollAppId)

	if !found {
		return false
	}

	return rollApp.Owner == account
}

// GetRollAppBech32Prefix returns the Bech32 prefix of the RollApp by the chain-id.
func (k Keeper) GetRollAppBech32Prefix(ctx sdk.Context, chainId string) (bech32Prefix string, found bool) {
	rollApp, found := k.rollappKeeper.GetRollapp(ctx, chainId)
	if found && len(rollApp.Bech32Prefix) > 0 {
		return rollApp.Bech32Prefix, true
	}

	if data, found := mockRollAppsData[chainId]; found {
		return data.bech32, len(data.bech32) > 0
	}

	return "", false
}
