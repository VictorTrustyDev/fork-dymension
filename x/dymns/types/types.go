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
	EventTypeDymNameOpenPurchaseOrder     = ModuleName + "_opo"
	AttributeKeyDymNameOpoActionName      = "action"
	AttributeKeyDymNameOpoName            = "name"
	AttributeKeyDymNameOpoExpiryEpoch     = "expiry_epoch"
	AttributeKeyDymNameOpoMinPrice        = "min_price"
	AttributeKeyDymNameOpoSellPrice       = "sell_price"
	AttributeKeyDymNameOpoHighestBidder   = "highest_bidder"
	AttributeKeyDymNameOpoHighestBidPrice = "highest_bid_price"
)

const (
	AttributeKeyDymNameOpoActionNameSet    = "set"
	AttributeKeyDymNameOpoActionNameDelete = "delete"
)
