package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrValidationFailed = errorsmod.Register(ModuleName, 1, "validation failed")
	ErrInvalidOwner     = errorsmod.Register(ModuleName, 2, "invalid owner")
)
