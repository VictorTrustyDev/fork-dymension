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

// prefix bytes for the DymNS persistent store.
const (
	prefixDymName                               = iota + 1
	prefixRvlDymNamesOwnedByAccount             // reverse lookup store
	prefixRvlConfiguredAddressToDymNamesInclude // reverse lookup store
	prefixRvlFallbackAddressToDymNamesInclude   // reverse lookup store
	prefixSellOrder
	prefixActiveSellOrdersExpiration
	prefixHistoricalSellOrders
	prefixMinExpiryHistoricalSellOrders
	prefixCountBuyOrders
	prefixBuyOrder
	prefixRvlBuyerToBuyOrderIds   // reverse lookup store
	prefixRvlGoodsIdToBuyOrderIds // reverse lookup store
	prefixRollAppIdToAliases
	prefixRvlAliasToRollAppId // reverse lookup store
)

const (
	// partialStoreOrderTypeDymName is a part of the store key prefix for the SellOrder records of Dym-Name
	partialStoreOrderTypeDymName = iota

	// partialStoreOrderTypeAlias is a part of the store key prefix for the SellOrder records of Alias
	partialStoreOrderTypeAlias
)

var (
	// KeyPrefixDymName is the key prefix for the DymName records
	KeyPrefixDymName = []byte{prefixDymName}

	// KeyPrefixRvlDymNamesOwnedByAccount is the key prefix for the reverse lookup for Dym-Names owned by an account
	KeyPrefixRvlDymNamesOwnedByAccount = []byte{prefixRvlDymNamesOwnedByAccount}

	// KeyPrefixRvlConfiguredAddressToDymNamesInclude is the key prefix for the reverse lookup for Dym-Names that contain the configured address (bech32)
	KeyPrefixRvlConfiguredAddressToDymNamesInclude = []byte{prefixRvlConfiguredAddressToDymNamesInclude}

	// KeyPrefixRvlFallbackAddressToDymNamesInclude is the key prefix for the reverse lookup address for Dym-Names using fallback mechanism
	KeyPrefixRvlFallbackAddressToDymNamesInclude = []byte{prefixRvlFallbackAddressToDymNamesInclude}

	// KeyPrefixSellOrder is the key prefix for the active SellOrder records of both type DymName/Alias
	KeyPrefixSellOrder = []byte{prefixSellOrder}

	// KeyPrefixDymNameSellOrder is the key prefix for the active SellOrder records of type DymName
	KeyPrefixDymNameSellOrder = []byte{prefixSellOrder, partialStoreOrderTypeDymName}

	// KeyPrefixDymNameHistoricalSellOrders is the key prefix for the historical SellOrder records of type DymName
	KeyPrefixDymNameHistoricalSellOrders = []byte{prefixHistoricalSellOrders, partialStoreOrderTypeDymName}

	// KeyPrefixMinExpiryDymNameHistoricalSellOrders is the key prefix for the lowest expiry among the historical SellOrder records of each specific Dym-Name
	KeyPrefixMinExpiryDymNameHistoricalSellOrders = []byte{prefixMinExpiryHistoricalSellOrders, partialStoreOrderTypeDymName}

	// KeyPrefixAliasSellOrder is the key prefix for the active SellOrder records of type Alias
	KeyPrefixAliasSellOrder = []byte{prefixSellOrder, partialStoreOrderTypeAlias}

	// KeyPrefixAliasHistoricalSellOrders is the key prefix for the historical SellOrder records of type Alias
	KeyPrefixAliasHistoricalSellOrders = []byte{prefixHistoricalSellOrders, partialStoreOrderTypeAlias}

	// KeyPrefixMinExpiryAliasHistoricalSellOrders is the key prefix for the lowest expiry among the historical SellOrder records of each specific Alias
	KeyPrefixMinExpiryAliasHistoricalSellOrders = []byte{prefixMinExpiryHistoricalSellOrders, partialStoreOrderTypeAlias}

	// KeyPrefixBuyOrder is the key prefix for the active BuyOrder records regardless order type DymName/Alias
	KeyPrefixBuyOrder = []byte{prefixBuyOrder}

	// KeyPrefixRvlBuyerToBuyOrderIds is the key prefix for the reverse lookup for BuyOrder IDs by the buyer
	KeyPrefixRvlBuyerToBuyOrderIds = []byte{prefixRvlBuyerToBuyOrderIds}

	// KeyPrefixRvlDymNameToBuyOrderIds is the key prefix for the reverse lookup for BuyOrder IDs by the DymName
	KeyPrefixRvlDymNameToBuyOrderIds = []byte{prefixRvlGoodsIdToBuyOrderIds, partialStoreOrderTypeDymName}

	// KeyPrefixRvlAliasToBuyOrderIds is the key prefix for the reverse lookup for BuyOrder IDs by the Alias
	KeyPrefixRvlAliasToBuyOrderIds = []byte{prefixRvlGoodsIdToBuyOrderIds, partialStoreOrderTypeAlias}

	// KeyPrefixRollAppIdToAliases is the key prefix for the Roll-App ID to Alias records
	KeyPrefixRollAppIdToAliases = []byte{prefixRollAppIdToAliases}

	// KeyPrefixRvlAliasToRollAppId is the key prefix for the reverse lookup for Alias to Roll-App ID records
	KeyPrefixRvlAliasToRollAppId = []byte{prefixRvlAliasToRollAppId}
)

var (
	KeyActiveSellOrdersExpirationOfDymName = []byte{prefixActiveSellOrdersExpiration, partialStoreOrderTypeDymName}

	KeyActiveSellOrdersExpirationOfAlias = []byte{prefixActiveSellOrdersExpiration, partialStoreOrderTypeAlias}

	// KeyCountBuyOrders is the key for the count of all-time buy orders
	KeyCountBuyOrders = []byte{prefixCountBuyOrders}
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

// FallbackAddressToDymNamesIncludeRvlKey returns the key for the reverse lookup address for Dym-Names using fallback mechanism
func FallbackAddressToDymNamesIncludeRvlKey(fallbackAddr FallbackAddress) []byte {
	return append(KeyPrefixRvlFallbackAddressToDymNamesInclude, fallbackAddr...)
}

// SellOrderKey returns a key for the active Sell-Order of the Dym-Name/Alias
func SellOrderKey(goodsId string, orderType OrderType) []byte {
	switch orderType {
	case NameOrder:
		return append(KeyPrefixDymNameSellOrder, []byte(goodsId)...)
	case AliasOrder:
		return append(KeyPrefixAliasSellOrder, []byte(goodsId)...)
	default:
		panic("invalid order type: " + orderType.FriendlyString())
	}
}

// HistoricalSellOrdersKey returns a key for the historical Sell-Orders of the Dym-Name/Alias
func HistoricalSellOrdersKey(goodsId string, orderType OrderType) []byte {
	switch orderType {
	case NameOrder:
		return append(KeyPrefixDymNameHistoricalSellOrders, []byte(goodsId)...)
	case AliasOrder:
		return append(KeyPrefixAliasHistoricalSellOrders, []byte(goodsId)...)
	default:
		panic("invalid order type: " + orderType.FriendlyString())
	}
}

// MinExpiryHistoricalSellOrdersKey returns a key for lowest expiry among the historical Sell-Orders of the Dym-Name/Alias
func MinExpiryHistoricalSellOrdersKey(goodsId string, orderType OrderType) []byte {
	switch orderType {
	case NameOrder:
		return append(KeyPrefixMinExpiryDymNameHistoricalSellOrders, []byte(goodsId)...)
	case AliasOrder:
		return append(KeyPrefixMinExpiryAliasHistoricalSellOrders, []byte(goodsId)...)
	default:
		panic("invalid order type: " + orderType.FriendlyString())
	}
}

// BuyOrderKey returns a key for the active Buy-Order of the Dym-Name/Alias
func BuyOrderKey(orderId string) []byte {
	return append(KeyPrefixBuyOrder, []byte(orderId)...)
}

// BuyerToOrderIdsRvlKey returns a key for reverse lookup for Buy-Order IDs by the buyer
func BuyerToOrderIdsRvlKey(bzHexAddr []byte) []byte {
	return append(KeyPrefixRvlBuyerToBuyOrderIds, bzHexAddr...)
}

// DymNameToBuyOrderIdsRvlKey returns a key for reverse lookup for Buy-Order IDs by the Dym-Name
func DymNameToBuyOrderIdsRvlKey(dymName string) []byte {
	return append(KeyPrefixRvlDymNameToBuyOrderIds, []byte(dymName)...)
}

// AliasToBuyOrderIdsRvlKey returns a key for reverse lookup for Buy-Order IDs by the Alias
func AliasToBuyOrderIdsRvlKey(alias string) []byte {
	return append(KeyPrefixRvlAliasToBuyOrderIds, []byte(alias)...)
}

// RollAppIdToAliasesKey returns a key for the Roll-App ID to list of alias records
func RollAppIdToAliasesKey(rollAppId string) []byte {
	return append(KeyPrefixRollAppIdToAliases, []byte(rollAppId)...)
}

// AliasToRollAppIdRvlKey returns a key for reverse lookup for Alias to Roll-App ID records
func AliasToRollAppIdRvlKey(alias string) []byte {
	return append(KeyPrefixRvlAliasToRollAppId, []byte(alias)...)
}
