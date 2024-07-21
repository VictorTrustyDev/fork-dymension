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
	prefixDymName                                  = iota + 1
	prefixRvlDymNamesOwnedByAccount                // reverse lookup store
	prefixRvlConfiguredAddressToDymNamesInclude    // reverse lookup store
	prefixRvlCoinType60HexAddressToDymNamesInclude // reverse lookup store
	prefixSellOrder
	prefixActiveSellOrdersExpiration
	prefixHistoricalSellOrders
	prefixMinExpiryHistoricalSellOrders
)

var (
	// KeyPrefixDymName is the key prefix for the Dym-Name records
	KeyPrefixDymName = []byte{prefixDymName}

	// KeyPrefixRvlDymNamesOwnedByAccount is the key prefix for the reverse lookup for Dym-Names owned by an account
	KeyPrefixRvlDymNamesOwnedByAccount = []byte{prefixRvlDymNamesOwnedByAccount}

	KeyPrefixRvlConfiguredAddressToDymNamesInclude = []byte{prefixRvlConfiguredAddressToDymNamesInclude}

	KeyPrefixRvlCoinType60HexAddressToDymNamesInclude = []byte{prefixRvlCoinType60HexAddressToDymNamesInclude}

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
	return append(KeyPrefixRvlDymNamesOwnedByAccount, owner.Bytes()...)
}

// ConfiguredAddressToDymNamesIncludeRvlKey returns a key for reverse lookup for Dym-Names that contain the configured address
func ConfiguredAddressToDymNamesIncludeRvlKey(address string) []byte {
	return append(KeyPrefixRvlConfiguredAddressToDymNamesInclude, []byte(address)...)
}

// CoinType60HexAddressToDymNamesIncludeRvlKey returns a key for reverse lookup for Dym-Names that contain the 0x address (coin-type 60, secp256k1, ethereum address)
func CoinType60HexAddressToDymNamesIncludeRvlKey(coinType60AccAddr []byte) []byte {
	return append(KeyPrefixRvlCoinType60HexAddressToDymNamesInclude, coinType60AccAddr...)
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
