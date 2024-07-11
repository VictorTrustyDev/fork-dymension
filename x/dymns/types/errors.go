package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
)

var (
	ErrValidationFailed          = errorsmod.Register(ModuleName, 1, "validation failed")
	ErrInvalidOwner              = errorsmod.Register(ModuleName, 2, "invalid owner")
	ErrInvalidState              = errorsmod.Register(ModuleName, 3, "invalid state")
	ErrDymNameNotFound           = errorsmod.Register(ModuleName, 4, "Dym-Name could not be found")
	ErrOpenPurchaseOrderNotFound = errorsmod.Register(ModuleName, 5, "open purchase order could not be found")
	ErrGracePeriod               = errorsmod.Register(ModuleName, 6, "expired Dym-Name still in grace period")
	ErrBadDymNameAddress         = errorsmod.Register(ModuleName, 7, "bad format Dym-Name address")
	ErrDymNameTooLong            = errorsmod.Register(ModuleName, 8, fmt.Sprintf("Dym-Name is too long, maximum %d characters", MaxDymNameLength))
)
