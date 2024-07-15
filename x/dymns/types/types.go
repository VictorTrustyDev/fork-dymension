package types

const (
	EventTypeSetDymName            = ModuleName + "_name"
	AttributeKeyDymName            = "name"
	AttributeKeyDymNameOwner       = "owner"
	AttributeKeyDymNameController  = "controller"
	AttributeKeyDymNameExpiryEpoch = "expiry_epoch"
	AttributeKeyDymNameConfigCount = "cfg_count"
)

const (
	EventTypeDymNameRefundBid       = ModuleName + "_bid_refund"
	AttributeKeyDymNameRefundBidder = "bidder"
	AttributeKeyDymNameRefundAmount = "amount"
)

const (
	EventTypeDymNameSellOrder            = ModuleName + "_so"
	AttributeKeyDymNameSoActionName      = "action"
	AttributeKeyDymNameSoName            = "name"
	AttributeKeyDymNameSoExpiryEpoch     = "expiry_epoch"
	AttributeKeyDymNameSoMinPrice        = "min_price"
	AttributeKeyDymNameSoSellPrice       = "sell_price"
	AttributeKeyDymNameSoHighestBidder   = "highest_bidder"
	AttributeKeyDymNameSoHighestBidPrice = "highest_bid_price"
)

const (
	AttributeKeyDymNameSoActionNameSet    = "set"
	AttributeKeyDymNameSoActionNameDelete = "delete"
)
