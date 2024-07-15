package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "dymns"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName
)

// prefix bytes for the DymNS persistent store
const (
	prefixDymName                                 = iota + 1
	prefixRvlDymNameOwnedByAccount                // reverse lookup store
	prefixRvlConfiguredAddressToDymNameInclude    // reverse lookup store
	prefixRvlCoinType60HexAddressToDymNameInclude // reverse lookup store
	prefixSellOrder
	prefixActiveSellOrdersExpiration
	prefixHistoricalSellOrders
	prefixMinExpiryHistoricalSellOrders
)

var (
	// KeyPrefixDymName is the key prefix for the Dym-Name records
	KeyPrefixDymName = []byte{prefixDymName}

	// KeyPrefixRvlDymNameOwnedByAccount is the key prefix for the reverse lookup for Dym-Names owned by an account
	KeyPrefixRvlDymNameOwnedByAccount = []byte{prefixRvlDymNameOwnedByAccount}

	KeyPrefixRvlConfiguredAddressToDymNameInclude = []byte{prefixRvlConfiguredAddressToDymNameInclude}

	KeyPrefixRvlCoinType60HexAddressToDymNameInclude = []byte{prefixRvlCoinType60HexAddressToDymNameInclude}

	// KeyPrefixSellOrder is the key prefix for the active Sell-Order records
	KeyPrefixSellOrder = []byte{prefixSellOrder}

	// KeyPrefixHistoricalSellOrders is the key prefix for the historical Sell-Order records
	KeyPrefixHistoricalSellOrders = []byte{prefixHistoricalSellOrders}

	// KeyPrefixMinExpiryHistoricalSellOrders is the key prefix for the lowest expiry among the historical Sell-Order records of each specific Dym-Name
	KeyPrefixMinExpiryHistoricalSellOrders = []byte{prefixMinExpiryHistoricalSellOrders}
)

var (
	KeyActiveSellOrdersExpiration = []byte{prefixActiveSellOrdersExpiration}
)

// DymNameKey returns a key for specific Dym-Name
func DymNameKey(name string) []byte {
	return append(KeyPrefixDymName, []byte(name)...)
}

// DymNamesOwnedByAccountRvlKey returns a key for reverse lookup for Dym-Names owned by an account
func DymNamesOwnedByAccountRvlKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixRvlDymNameOwnedByAccount, owner.Bytes()...)
}

func ConfiguredAddressToDymNameIncludeRvlKey(address string) []byte {
	return append(KeyPrefixRvlConfiguredAddressToDymNameInclude, []byte(address)...)
}

func CoinType60HexAddressToDymNameIncludeRvlKey(coinType60AccAddr sdk.AccAddress) []byte {
	return append(KeyPrefixRvlCoinType60HexAddressToDymNameInclude, []byte(coinType60AccAddr)...)
}

// SellOrderKey returns a key for the active Sell-Order of the Dym-Name
func SellOrderKey(dymName string) []byte {
	return append(KeyPrefixSellOrder, []byte(dymName)...)
}

// HistoricalSellOrdersKey returns a key for the historical Sell-Orders of the Dym-Name
func HistoricalSellOrdersKey(dymName string) []byte {
	return append(KeyPrefixHistoricalSellOrders, []byte(dymName)...)
}

// MinExpiryHistoricalSellOrdersKey returns a key for lowest expiry among the historical Sell-Orders
// of the Dym-Name
func MinExpiryHistoricalSellOrdersKey(dymName string) []byte {
	return append(KeyPrefixMinExpiryHistoricalSellOrders, []byte(dymName)...)
}
