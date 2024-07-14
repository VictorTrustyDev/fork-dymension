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
	prefixDymName                  = iota + 1
	prefixRvlDymNameOwnedByAccount // reverse lookup store
	//prefixRvlConfiguredAddressToDymNameInclude // reverse lookup store
	prefixOpenPurchaseOrder
	prefixActiveOpenPurchaseOrdersExpiration
	prefixHistoricalOpenPurchaseOrders
	prefixMinExpiryHistoricalOpenPurchaseOrders
)

var (
	KeyPrefixDymName                  = []byte{prefixDymName}
	KeyPrefixRvlDymNameOwnedByAccount = []byte{prefixRvlDymNameOwnedByAccount}
	//KeyPrefixRvlConfiguredAddressToDymNameInclude  = []byte{prefixRvlConfiguredAddressToDymNameInclude}
	KeyPrefixOpenPurchaseOrder                     = []byte{prefixOpenPurchaseOrder}
	KeyPrefixHistoricalOpenPurchaseOrders          = []byte{prefixHistoricalOpenPurchaseOrders}
	KeyPrefixMinExpiryHistoricalOpenPurchaseOrders = []byte{prefixMinExpiryHistoricalOpenPurchaseOrders}
)

var (
	KeyActiveOpenPurchaseOrdersExpiration = []byte{prefixActiveOpenPurchaseOrdersExpiration}
)

// DymNameKey returns a key for specific Dym-Name
func DymNameKey(name string) []byte {
	return append(KeyPrefixDymName, []byte(name)...)
}

// DymNamesOwnedByAccountRvlKey returns a key for reverse lookup for Dym-Names owned by an account
func DymNamesOwnedByAccountRvlKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixRvlDymNameOwnedByAccount, owner.Bytes()...)
}

// OpenPurchaseOrderKey returns a key for open purchase order for a Dym-Name
func OpenPurchaseOrderKey(dymName string) []byte {
	return append(KeyPrefixOpenPurchaseOrder, []byte(dymName)...)
}

// HistoricalOpenPurchaseOrdersKey returns a key for historical open purchase orders for a Dym-Name
func HistoricalOpenPurchaseOrdersKey(dymName string) []byte {
	return append(KeyPrefixHistoricalOpenPurchaseOrders, []byte(dymName)...)
}

// MinExpiryHistoricalOpenPurchaseOrdersKey returns a key for lowest expiry among historical open purchase orders of a Dym-Name
func MinExpiryHistoricalOpenPurchaseOrdersKey(dymName string) []byte {
	return append(KeyPrefixMinExpiryHistoricalOpenPurchaseOrders, []byte(dymName)...)
}
