package utils

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/dymensionxyz/dymension/v3/app/params"
	"regexp"
	"strings"
)

func IsValidBech32AccountAddress(address string, matchAccountAddressBech32Prefix bool) bool {
	hrp, bz, err := bech32.DecodeAndConvert(address)

	if err != nil {
		return false
	}

	bytesCount := len(bz)
	if bytesCount != 20 && bytesCount != 32 /*32 bytes is interchain account*/ {
		return false
	}

	return !matchAccountAddressBech32Prefix || hrp == params.AccountAddressPrefix
}

var pattern0xHex = regexp.MustCompile(`^0x[a-f\d]+$`)

func IsValid0xAddress(address string) bool {
	length := len(address)
	if length != 42 && length != 66 /*32 bytes is interchain account*/ {
		return false
	}

	address = strings.ToLower(address)
	if !pattern0xHex.MatchString(address) {
		return false
	}

	return true
}
