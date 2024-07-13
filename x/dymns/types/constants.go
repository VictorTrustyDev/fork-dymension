package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	MaxDymNameLength = 20
)

const (
	OpGasPutAds   sdk.Gas = 5_000_000
	OpGasBidAds   sdk.Gas = 20_000_000
	OpGasCloseAds sdk.Gas = 5_000_000

	OpGasConfig sdk.Gas = 30_000_000
)
