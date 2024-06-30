package utils

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/dymensionxyz/dymension/v3/app/params"
)

func IsValidBech32AccountAddress(address string, matchAccountAddressBech32Prefix bool) bool {
	hrp, _, err := bech32.DecodeAndConvert(address)

	if err != nil {
		return false
	}

	return !matchAccountAddressBech32Prefix || hrp == params.AccountAddressPrefix
}
