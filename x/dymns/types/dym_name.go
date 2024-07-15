package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

func (m DymNameConfig) IsDelete() bool {
	return m.Value == ""
}

func (m *DymName) GetAddressReverseMappingRecords() (
	configuredAddressesToDymNames map[string]ReverseLookupDymNames,
	coinType60HexAddressesToDymNames map[string]ReverseLookupDymNames,
) {
	// TODO DymNS: call this method from the keeper, to cleanup before modify the Dym-Name configuration
	// TODO DymNS: call this method from the keeper, to update after modify the Dym-Name configuration
	if err := m.Validate(); err != nil {
		// should validate before calling this method
		panic(err)
	}

	configuredAddressesToDymNames = make(map[string]ReverseLookupDymNames)
	coinType60HexAddressesToDymNames = make(map[string]ReverseLookupDymNames)

	addConfiguredAddressToDymNames := func(address string, dymName string) {
		existing, _ := configuredAddressesToDymNames[address]
		configuredAddressesToDymNames[address] = existing.Combine(ReverseLookupDymNames{
			DymNames: []string{
				dymName,
			},
		})
	}

	addCoinType60HexAddressToDymNames := func(accAddr sdk.AccAddress, dymName string) {
		var strAddr string
		if len(accAddr.Bytes()) == 32 { // Interchain Account
			strAddr = common.BytesToHash(accAddr.Bytes()).String()
		} else {
			strAddr = common.BytesToAddress(accAddr.Bytes()).String()
		}
		strAddr = strings.ToLower(strAddr)
		existing, _ := coinType60HexAddressesToDymNames[strAddr]
		coinType60HexAddressesToDymNames[strAddr] = existing.Combine(ReverseLookupDymNames{
			DymNames: []string{
				dymName,
			},
		})
	}

	var nameConfigs []DymNameConfig
	for _, config := range m.Configs {
		if config.Type == DymNameConfigType_NAME {
			nameConfigs = append(nameConfigs, config)
		}
	}

	var defaultConfig *DymNameConfig
	for _, config := range nameConfigs {
		if config.ChainId == "" && config.Path == "" {
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

		if !dymnsutils.IsValidBech32AccountAddress(config.Value, false) {
			// should not happen as configuration should be validated before calling this method.
			// But code still be kept to be aware of the possibility of future changes.
			panic("current implementation only accept bech32 account address")
		}

		if config.ChainId == "" && config.Path == "" {
			// default config

			accAddr, err := sdk.AccAddressFromBech32(config.Value)
			if err != nil {
				// should not happen as configuration should be validated before calling this method
				panic(err)
			}

			addConfiguredAddressToDymNames(config.Value, m.Name)
			addCoinType60HexAddressToDymNames(accAddr, m.Name)

			continue
		}

		_, bz, err := bech32.DecodeAndConvert(config.Value)
		if err != nil {
			// should not happen as configuration should be validated before calling this method
			// But code still be kept to be aware of the possibility of future changes.
			panic(err)
		}

		addConfiguredAddressToDymNames(config.Value, m.Name)
		addCoinType60HexAddressToDymNames(bz, m.Name)
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
