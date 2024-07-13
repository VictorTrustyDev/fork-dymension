package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	MaxDymNameLength = 20
)

const (
	// OpGasPutAds is the gas consumed for Dym-Name owner to put an ads for selling Dym-Name.
	OpGasPutAds sdk.Gas = 5_000_000
	// OpGasCloseAds is the gas consumed for Dym-Name owner to close the Dym-Name ads.
	OpGasCloseAds sdk.Gas = 5_000_000

	// OpGasBidAds is the gas consumed for bidding an ads for Dym-Name.
	OpGasBidAds sdk.Gas = 20_000_000

	// OpGasConfig is the gas consumed for updating Dym-Name configuration,
	// We charge this high amount of gas for extra permanent data
	// needed to be stored like reverse lookup record.
	// So we do not charge this fee on Delete operation.
	OpGasConfig sdk.Gas = 30_000_000
)
