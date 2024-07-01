package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrValidationFailed          = errorsmod.Register(ModuleName, 1, "validation failed")
	ErrInvalidOwner              = errorsmod.Register(ModuleName, 2, "invalid owner")
	ErrInvalidState              = errorsmod.Register(ModuleName, 3, "invalid state")
	ErrOpenPurchaseOrderNotFound = errorsmod.Register(ModuleName, 4, "open purchase order not found")
)
