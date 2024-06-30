package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	prefixOpenPurchaseOrderDymNameByDymName
)

var (
	KeyPrefixDymName                    = []byte{prefixDymName}
	KeyPrefixDymNameOwnedByAccount      = []byte{prefixDymNameOwnedByAccount}
	KeyPrefixOpenPurchaseOrderByDymName = []byte{prefixOpenPurchaseOrderDymNameByDymName}
)

// DymNameKey returns a key for specific Dym-Name
func DymNameKey(name string) []byte {
	return append(KeyPrefixDymName, crypto.Keccak256([]byte(name))...)
}

// DymNamesOwnedByAccountKey returns a key for Dym-Names owned by an account
func DymNamesOwnedByAccountKey(owner sdk.AccAddress) []byte {
	return append(KeyPrefixDymNameOwnedByAccount, owner.Bytes()...)
}

// OpenPurchaseOrderKey returns a key for open purchase order for a Dym-Name
func OpenPurchaseOrderKey(dymName string) []byte {
	return append(KeyPrefixOpenPurchaseOrderByDymName, crypto.Keccak256([]byte(dymName))...)
}
