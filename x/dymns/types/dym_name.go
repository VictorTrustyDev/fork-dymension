package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/ethereum/go-ethereum/common"
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
		return ErrValidationFailed.Wrapf("owner is not a valid bech32 account address: %s", m.Owner)
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

	if m.Value != strings.ToLower(m.Value) {
		return ErrValidationFailed.Wrap("dym name config value must be lowercase")
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

func (m *ReverseLookupDymNames) Validate() error {
	if m == nil {
		return ErrValidationFailed.Wrap("reverse lookup record is nil")
	}

	for _, name := range m.DymNames {
		if !dymnsutils.IsValidDymName(name) {
			return ErrValidationFailed.Wrapf("invalid dym name: %s", name)
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

// IsDefaultNameConfig checks if the config is a default name config, satisfy the following conditions:
//   - Type is NAME
//   - ChainId is empty (means host chain)
//   - Path is empty (means root Dym-Name)
func (m DymNameConfig) IsDefaultNameConfig() bool {
	return m.Type == DymNameConfigType_NAME &&
		m.ChainId == "" &&
		m.Path == ""
}

func (m DymNameConfig) IsDelete() bool {
	return m.Value == ""
}

func (m *DymName) GetAddressesForReverseMapping() (
	configuredAddressesToConfigs map[string][]DymNameConfig,
	hexAddressesToConfigs map[string][]DymNameConfig,
) {
	if err := m.Validate(); err != nil {
		// should validate before calling this method
		panic(err)
	}

	configuredAddressesToConfigs = make(map[string][]DymNameConfig)
	hexAddressesToConfigs = make(map[string][]DymNameConfig)

	addConfiguredAddress := func(address string, config DymNameConfig) {
		configuredAddressesToConfigs[address] = append(configuredAddressesToConfigs[address], config)
	}

	addHexAddress := func(accAddr sdk.AccAddress, config DymNameConfig) {
		var strAddr string
		if len(accAddr.Bytes()) == 32 { // Interchain Account
			strAddr = common.BytesToHash(accAddr.Bytes()).String()
		} else {
			strAddr = common.BytesToAddress(accAddr.Bytes()).String()
		}
		strAddr = strings.ToLower(strAddr)
		hexAddressesToConfigs[strAddr] = append(hexAddressesToConfigs[strAddr], config)
	}

	var nameConfigs []DymNameConfig
	for _, config := range m.Configs {
		if config.Type == DymNameConfigType_NAME {
			nameConfigs = append(nameConfigs, config)
		}
	}

	var defaultConfig *DymNameConfig
	for i, config := range nameConfigs {
		if config.IsDefaultNameConfig() {
			if config.Value == "" {
				config.Value = m.Owner
				nameConfigs[i] = config
			}

			defaultConfig = &config
			break
		}
	}

	if defaultConfig == nil {
		// add a fake record to be used to generate default address
		nameConfigs = append(nameConfigs, DymNameConfig{
			Type:    DymNameConfigType_NAME,
			ChainId: "",
			Path:    "",
			Value:   m.Owner,
		})
	}

	for _, config := range nameConfigs {
		if config.Value == "" {
			continue
		}

		// just a friendly reminder, in the current implementation,
		// config value is always a bech32 account address

		if config.IsDefaultNameConfig() {

			accAddr, err := sdk.AccAddressFromBech32(config.Value)
			if err != nil {
				// should not happen as configuration should be validated before calling this method
				panic(err)
			}

			addConfiguredAddress(config.Value, config)
			addHexAddress(accAddr, config)

			continue
		}

		addConfiguredAddress(config.Value, config)

		// note: this config is not a default config, is not a part of fallback mechanism,
		// so we don't add hex address for this config
	}

	return
}

func (m ReverseLookupDymNames) Distinct() ReverseLookupDymNames {
	var uniqueDymNames = make(map[string]bool)
	for _, name := range m.DymNames {
		uniqueDymNames[name] = true
	}
	distinctDymNames := make([]string, 0, len(uniqueDymNames))
	for name := range uniqueDymNames {
		distinctDymNames = append(distinctDymNames, name)
	}
	return ReverseLookupDymNames{
		DymNames: distinctDymNames,
	}
}

func (m ReverseLookupDymNames) Combine(other ReverseLookupDymNames) ReverseLookupDymNames {
	return ReverseLookupDymNames{
		DymNames: append(m.DymNames, other.DymNames...),
	}.Distinct()
}

func (m ReverseLookupDymNames) Exclude(toBeExcluded ReverseLookupDymNames) ReverseLookupDymNames {
	if len(toBeExcluded.DymNames) == 0 {
		return m
	}

	var excludedDymNames = make(map[string]bool)
	for _, name := range toBeExcluded.DymNames {
		excludedDymNames[name] = true
	}

	var filteredDymNames = make([]string, 0, len(m.DymNames))
	for _, name := range m.DymNames {
		if !excludedDymNames[name] {
			filteredDymNames = append(filteredDymNames, name)
		}
	}

	return ReverseLookupDymNames{
		DymNames: filteredDymNames,
	}.Distinct()
}
