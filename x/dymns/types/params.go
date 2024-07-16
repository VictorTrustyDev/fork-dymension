package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	// KeyPriceParams is the key for the price params
	KeyPriceParams = []byte("PriceParams")

	// KeyChainsParams is the key for the chains params
	KeyChainsParams = []byte("ChainsParams")

	// KeyMiscParams is the key for the misc params
	KeyMiscParams = []byte("MiscParams")
)

const (
	defaultBeginEpochHookIdentifier = "day" // less-frequently for cleanup
	defaultEndEpochHookIdentifier   = "hour"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// TODO DymNS: I'm not really familiar with this kind of params update via GOV, so please test with care.

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultPriceParams(),
		DefaultChainsParams(),
		DefaultMiscParams(),
	)
}

// DefaultPriceParams returns a default set of price parameters
func DefaultPriceParams() PriceParams {
	return PriceParams{
		Price_1Letter:      sdk.NewInt(5000 /* DYM */).MulRaw(1e18),
		Price_2Letters:     sdk.NewInt(2500 /* DYM */).MulRaw(1e18),
		Price_3Letters:     sdk.NewInt(1000 /* DYM */).MulRaw(1e18),
		Price_4Letters:     sdk.NewInt(100 /* DYM */).MulRaw(1e18),
		Price_5PlusLetters: sdk.NewInt(5 /* DYM */).MulRaw(1e18),
		PriceExtends:       sdk.NewInt(5 /* DYM */).MulRaw(1e18),
		PriceDenom:         params.BaseDenom,
	}
}

// DefaultChainsParams returns a default set of chains configuration
func DefaultChainsParams() ChainsParams {
	//goland:noinspection SpellCheckingInspection
	return ChainsParams{
		AliasesByChainId: make(map[string]AliasesOfChainId),
		CoinType60ChainIds: []string{
			"evmos_9001-2",       // Evmos Mainnet
			"evmos_9001-3",       // Evmos Mainnet Re-Genesis
			"injective-1",        // Injective Mainnet
			"cronosmainnet_25-1", // Cronos Mainnet
		},
	}
}

// DefaultMiscParams returns a default set of misc parameters
func DefaultMiscParams() MiscParams {
	return MiscParams{
		BeginEpochHookIdentifier:     defaultBeginEpochHookIdentifier,
		EndEpochHookIdentifier:       defaultEndEpochHookIdentifier,
		DaysGracePeriod:              30,
		DaysSellOrderDuration:        3,
		DaysPreservedClosedSellOrder: 7,
		DaysProhibitSell:             30,
	}
}

func NewParams(price PriceParams, chains ChainsParams, misc MiscParams) Params {
	return Params{
		Price:  price,
		Chains: chains,
		Misc:   misc,
	}
}

// ParamSetPairs get the params.ParamSet
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPriceParams, &m.Price, validatePriceParams),
		paramtypes.NewParamSetPair(KeyChainsParams, &m.Chains, validateChainsParams),
		paramtypes.NewParamSetPair(KeyMiscParams, &m.Misc, validateMiscParams),
	}
}

func (m *Params) Validate() error {
	if err := m.Price.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("price params: %v", err)
	}
	if err := m.Chains.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("chains params: %v", err)
	}
	if err := m.Misc.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("misc params: %v", err)
	}
	return nil
}

func (m PriceParams) Validate() error {
	return validatePriceParams(m)
}

func (m ChainsParams) Validate() error {
	return validateChainsParams(m)
}

func (m MiscParams) Validate() error {
	return validateMiscParams(m)
}

func validateEpochIdentifier(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("epoch identifier cannot be empty")
	}
	switch v {
	case "hour", "day", "week":
	default:
		return fmt.Errorf("invalid epoch identifier: %s", v)
	}
	return nil
}

func validatePriceParams(i interface{}) error {
	m, ok := i.(PriceParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	validatePrice := func(price sdkmath.Int, letterDesc string) error {
		if price.IsNil() || price.IsZero() {
			return ErrValidationFailed.Wrapf("%s Dym-Name price is zero", letterDesc)
		} else if price.IsNegative() {
			return ErrValidationFailed.Wrapf("%s Dym-Name price is negative", letterDesc)
		}
		return nil
	}

	if err := validatePrice(m.Price_1Letter, "1 letter"); err != nil {
		return err
	}

	if err := validatePrice(m.Price_2Letters, "2 letters"); err != nil {
		return err
	}

	if err := validatePrice(m.Price_3Letters, "3 letters"); err != nil {
		return err
	}

	if err := validatePrice(m.Price_4Letters, "4 letters"); err != nil {
		return err
	}

	if err := validatePrice(m.Price_5PlusLetters, "5+ letters"); err != nil {
		return err
	}

	if err := validatePrice(m.PriceExtends, "yearly extends"); err != nil {
		return err
	}

	if m.Price_1Letter.LTE(m.Price_2Letters) {
		return ErrValidationFailed.Wrap("1 letter price must be greater than 2 letters price")
	}

	if m.Price_2Letters.LTE(m.Price_3Letters) {
		return ErrValidationFailed.Wrap("2 letters price must be greater than 3 letters price")
	}

	if m.Price_3Letters.LTE(m.Price_4Letters) {
		return ErrValidationFailed.Wrap("3 letters price must be greater than 4 letters price")
	}

	if m.Price_4Letters.LTE(m.Price_5PlusLetters) {
		return ErrValidationFailed.Wrap("4 letters price must be greater than 5+ letters price")
	}

	if m.Price_5PlusLetters.LT(m.PriceExtends) {
		return ErrValidationFailed.Wrap("5 letters price must be greater or equals to yearly extend price")
	}

	if m.PriceDenom == "" {
		return ErrValidationFailed.Wrap("Dym-Name price denom cannot be empty")
	}

	if err := (sdk.Coin{
		Denom:  m.PriceDenom,
		Amount: sdk.ZeroInt(),
	}).Validate(); err != nil {
		return ErrValidationFailed.Wrapf("invalid Dym-Name price denom: %s", err)
	}

	return nil
}

func validateChainsParams(i interface{}) error {
	m, ok := i.(ChainsParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	uniqueChainIdAliasAmongAliasConfig := make(map[string]bool)
	for chainID, aliases := range m.AliasesByChainId {
		if len(chainID) < 3 {
			return ErrValidationFailed.Wrapf("alias: chain ID `%s` must be at least 3 characters", chainID)
		}

		if !dymnsutils.IsValidChainIdFormat(chainID) {
			return ErrValidationFailed.Wrapf("alias: chain ID `%s` is not well-formed", chainID)
		}

		if _, ok := uniqueChainIdAliasAmongAliasConfig[chainID]; ok {
			return ErrValidationFailed.Wrapf("alias: chain ID and alias must unique among all, found duplicated '%s'", chainID)
		}
		uniqueChainIdAliasAmongAliasConfig[chainID] = true

		for _, alias := range aliases.Aliases {
			if !dymnsutils.IsValidAlias(alias) {
				return ErrValidationFailed.Wrapf("alias `%s` is not well-formed", alias)
			}

			if _, ok := uniqueChainIdAliasAmongAliasConfig[alias]; ok {
				return ErrValidationFailed.Wrapf("alias: chain ID and alias must unique among all, found duplicated '%s'", alias)
			}
			uniqueChainIdAliasAmongAliasConfig[alias] = true
		}
	}

	uniqueChainIdAmongCoinType60ChainsConfig := make(map[string]bool)
	for _, chainID := range m.CoinType60ChainIds {
		if len(chainID) < 3 {
			return ErrValidationFailed.Wrapf("coin-type-60 chains: chain ID `%s` must be at least 3 characters", chainID)
		}

		if !dymnsutils.IsValidChainIdFormat(chainID) {
			return ErrValidationFailed.Wrapf("coin-type-60 chains: chain ID `%s` is not well-formed", chainID)
		}

		if _, ok := uniqueChainIdAmongCoinType60ChainsConfig[chainID]; ok {
			return ErrValidationFailed.Wrapf("coin-type-60 chains: chain ID is not unique '%s'", chainID)
		}
		uniqueChainIdAmongCoinType60ChainsConfig[chainID] = true
	}

	return nil
}

func validateMiscParams(i interface{}) error {
	m, ok := i.(MiscParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := validateEpochIdentifier(m.BeginEpochHookIdentifier); err != nil {
		return ErrValidationFailed.Wrapf("begin epoch hook identifier: %v", err)
	}

	if err := validateEpochIdentifier(m.EndEpochHookIdentifier); err != nil {
		return ErrValidationFailed.Wrapf("end epoch hook identifier: %v", err)
	}

	if m.DaysGracePeriod < 0 {
		return ErrValidationFailed.Wrap("days grace period cannot be negative")
	}

	if m.DaysSellOrderDuration < 1 {
		return ErrValidationFailed.Wrap("Sell Orders duration must be greater than 0 day")
	}

	if m.DaysPreservedClosedSellOrder < 1 {
		return ErrValidationFailed.Wrap("preserved closed Sell Orders must be greater than 0 day")
	}

	if m.DaysProhibitSell < 7 {
		return ErrValidationFailed.Wrap("prohibit sell must be at least 7 days")
	}

	return nil
}
