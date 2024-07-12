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
	prefixDymName = iota + 1
	prefixDymNameOwnedByAccount
	prefixOpenPurchaseOrder
	prefixActiveOpenPurchaseOrdersExpiration
	prefixHistoricalOpenPurchaseOrders
	prefixMinExpiryHistoricalOpenPurchaseOrders
)

var (
	KeyPrefixDymName                               = []byte{prefixDymName}
	KeyPrefixDymNameOwnedByAccount                 = []byte{prefixDymNameOwnedByAccount}
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

// DymNamesOwnedByAccountKey returns a key for Dym-Names owned by an account
func DymNamesOwnedByAccountKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixDymNameOwnedByAccount, owner.Bytes()...)
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