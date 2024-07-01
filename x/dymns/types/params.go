package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (

	// KeyEpochIdentifier is the key for the epoch identifier
	KeyEpochIdentifier = []byte("EpochIdentifier")

	// KeyPriceParams is the key for the price params
	KeyPriceParams = []byte("PriceParams")

	// KeyMiscParams is the key for the misc params
	KeyMiscParams = []byte("MiscParams")
)

const (
	defaultEpochIdentifier = "hour"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(defaultEpochIdentifier, DefaultPriceParams(), DefaultMiscParams())
}

// DefaultPriceParams returns a default set of price parameters
func DefaultPriceParams() PriceParams {
	return PriceParams{
		Price_1Letter:      sdk.NewInt(5000 /* DYM */).MulRaw(1e18),
		Price_2Letters:     sdk.NewInt(2500 /* DYM */).MulRaw(1e18),
		Price_3Letters:     sdk.NewInt(1000 /* DYM */).MulRaw(1e18),
		Price_4Letters:     sdk.NewInt(100 /* DYM */).MulRaw(1e18),
		Price_5PlusLetters: sdk.NewInt(3 /* DYM */).MulRaw(1e18),
		PriceDenom:         params.BaseDenom,
	}
}

func DefaultMiscParams() MiscParams {
	return MiscParams{
		// TODO DymNS: add days when create new OPO
		DaysOpenPurchaseOrderDuration: 3,
		// TODO DymNS: prune historical data
		DaysPreservedClosedPurchaseOrder: 7,
		// TODO DymNS: add gas for CRUD operations on OPO
		GasCrudOpenPurchaseOrder: 5_000_000,
	}
}

func NewParams(epochIdentifier string, price PriceParams, misc MiscParams) Params {
	return Params{
		EpochIdentifier: epochIdentifier,
		Price:           price,
		Misc:            misc,
	}
}

// ParamSetPairs get the params.ParamSet
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &m.EpochIdentifier, validateEpochIdentifier),
		paramtypes.NewParamSetPair(KeyPriceParams, &m.Price, validatePriceParams),
		paramtypes.NewParamSetPair(KeyMiscParams, &m.Misc, validateMiscParams),
	}
}

func (m *Params) Validate() error {
	if err := validateEpochIdentifier(m.EpochIdentifier); err != nil {
		return ErrValidationFailed.Wrapf("epoch identifier: %v", err)
	}
	if err := m.Price.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("price params: %v", err)
	}
	if err := m.Misc.Validate(); err != nil {
		return ErrValidationFailed.Wrapf("misc params: %v", err)
	}
	return nil
}

func (m PriceParams) Validate() error {
	return validatePriceParams(m)
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

func validateMiscParams(i interface{}) error {
	m, ok := i.(MiscParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if m.DaysOpenPurchaseOrderDuration < 1 {
		return ErrValidationFailed.Wrap("days OPO duration must be greater than 0")
	}

	if m.DaysPreservedClosedPurchaseOrder < 1 {
		return ErrValidationFailed.Wrap("days preserved closed OPO must be greater than 0")
	}

	if m.GasCrudOpenPurchaseOrder < 0 {
		return ErrValidationFailed.Wrap("gas for CRUD operations on OPO cannot be negative")
	}

	return nil
}
