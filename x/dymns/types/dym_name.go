package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"strings"
	"time"
)

func (m *DymName) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("dym name is nil")
	}
	if m.Name == "" {
		return ErrValidationFailed.Wrap("name is empty")
	}
	if !dymnsutils.IsValidDymName(m.Name) {
		return ErrValidationFailed.Wrap("name is not a valid dym name")
	}
	if m.Owner == "" {
		return ErrValidationFailed.Wrap("owner is empty")
	}
	if !dymnsutils.IsValidBech32AccountAddress(m.Owner, true) {
		return ErrValidationFailed.Wrap("owner is not a valid bech32 account address")
	}
	if m.Controller == "" {
		return ErrValidationFailed.Wrap("controller is empty")
	}
	if !dymnsutils.IsValidBech32AccountAddress(m.Controller, true) {
		return ErrValidationFailed.Wrap("controller is not a valid bech32 account address")
	}
	if m.ExpireAt == 0 {
		return ErrValidationFailed.Wrap("expire at is empty")
	}

	var uniqueConfig = make(map[string]bool)
	for _, config := range m.Configs {
		if err := config.Validate(); err != nil {
			return err
		}

		configIdentity := config.GetIdentity()
		if _, duplicated := uniqueConfig[configIdentity]; duplicated {
			return ErrValidationFailed.Wrapf("dym name config is not unique: %s", configIdentity)
		}
		uniqueConfig[configIdentity] = true
	}
	return nil
}

func (m *DymNameConfig) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("dym name config is nil")
	}

	if m.ChainId == "" {
		// ok to be empty
	} else if !dymnsutils.IsValidChainIdFormat(m.ChainId) {
		return ErrValidationFailed.Wrap("dym name config chain id must be a valid chain id format")
	}

	if m.Path == "" {
		// ok to be empty
	} else if !dymnsutils.IsValidSubDymName(m.Path) {
		return ErrValidationFailed.Wrap("dym name config path must be a valid dym name")
	}

	if m.Type == DymNameConfigType_NAME {
		if !m.IsDelete() && !dymnsutils.IsValidBech32AccountAddress(m.Value, false) {
			return ErrValidationFailed.Wrap("dym name config value must be a valid bech32 account address")
		}
	} else {
		return ErrValidationFailed.Wrapf("dym name config type is not %s - the only supported at this moment", DymNameConfigType_NAME.String())
	}

	return nil
}

func (m *OwnedDymNames) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("owned dym name is nil")
	}

	for _, name := range m.DymNames {
		if !dymnsutils.IsValidDymName(name) {
			return ErrValidationFailed.Wrapf("owned dym name is not a valid dym name: %s", name)
		}
	}

	return nil
}

func (m DymName) IsExpiredAt(anchor time.Time) bool {
	return m.IsExpiredAtEpoch(anchor.UTC().Unix())
}

func (m DymName) IsExpiredAtEpoch(epochUTC int64) bool {
	return m.ExpireAt < epochUTC
}

func (m DymName) GetSdkEvent() sdk.Event {
	return sdk.NewEvent(
		EventTypeSetDymName,
		sdk.NewAttribute(AttributeKeyDymName, m.Name),
		sdk.NewAttribute(AttributeKeyDymNameOwner, m.Owner),
		sdk.NewAttribute(AttributeKeyDymNameController, m.Controller),
		sdk.NewAttribute(AttributeKeyDymNameExpiryEpoch, fmt.Sprintf("%d", m.ExpireAt)),
		sdk.NewAttribute(AttributeKeyDymNameConfigCount, fmt.Sprintf("%d", len(m.Configs))),
	)
}

func (m DymNameConfig) GetIdentity() string {
	return strings.ToLower(fmt.Sprintf("%s|%s|%s", m.Type, m.ChainId, m.Path))
}

func (m DymNameConfig) IsDelete() bool {
	return m.Value == ""
}
