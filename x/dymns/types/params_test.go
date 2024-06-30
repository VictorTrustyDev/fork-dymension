package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParamKeyTable(t *testing.T) {
	m := ParamKeyTable()
	require.NotNil(t, m)
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, (&params).Validate())
}

func TestNewParams(t *testing.T) {
	params := NewParams(
		"a",
		PriceParams{
			PriceDenom: "b",
		},
		MiscParams{
			DaysOpenPurchaseOrderDuration:    888,
			DaysPreservedClosedPurchaseOrder: 88,
			GasCrudOpenPurchaseOrder:         8,
		},
	)
	require.Equal(t, "a", params.EpochIdentifier)
	require.Equal(t, "b", params.Price.PriceDenom)
	require.Equal(t, int32(888), params.Misc.DaysOpenPurchaseOrderDuration)
	require.Equal(t, int32(88), params.Misc.DaysPreservedClosedPurchaseOrder)
	require.Equal(t, int32(8), params.Misc.GasCrudOpenPurchaseOrder)
}

func TestDefaultPriceParams(t *testing.T) {
	require.NoError(t, DefaultPriceParams().Validate())
}

func TestDefaultMiscParams(t *testing.T) {
	require.NoError(t, DefaultMiscParams().Validate())
}

func TestParams_ParamSetPairs(t *testing.T) {
	params := DefaultParams()
	paramSetPairs := (&params).ParamSetPairs()
	require.Len(t, paramSetPairs, 3)
}

func TestParams_Validate(t *testing.T) {
	params := DefaultParams()
	require.NoError(t, (&params).Validate())

	params.EpochIdentifier = "--"
	require.Error(t, (&params).Validate())

	params = DefaultParams()
	params.Price.Price_1Letter = sdk.ZeroInt()
	require.Error(t, (&params).Validate())

	params = DefaultParams()
	params.Misc.DaysPreservedClosedPurchaseOrder = 0
	require.Error(t, (&params).Validate())
}

func TestPriceParams_Validate(t *testing.T) {
	validPriceParams := PriceParams{
		Price_1Letter:      sdk.NewInt(6),
		Price_2Letters:     sdk.NewInt(5),
		Price_3Letters:     sdk.NewInt(4),
		Price_4Letters:     sdk.NewInt(3),
		Price_5PlusLetters: sdk.NewInt(2),
		PriceDenom:         "adym",
	}

	require.NoError(t, validPriceParams.Validate())

	t.Run("price denom", func(t *testing.T) {
		m := validPriceParams
		m.PriceDenom = ""
		require.Error(t, m.Validate())

		m.PriceDenom = "--"
		require.Error(t, m.Validate())
	})

	type modifierPrice func(PriceParams, sdkmath.Int) PriceParams
	type swapPrice func(PriceParams) PriceParams

	testsInvalidPrice := []struct {
		name          string
		modifierPrice modifierPrice
		swapPrice     swapPrice
	}{
		{
			name:          "invalid 1 letter price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.Price_1Letter = v; return p },
			swapPrice: func(params PriceParams) PriceParams {
				backup := params.Price_1Letter
				params.Price_1Letter = params.Price_2Letters
				params.Price_2Letters = backup
				return params
			},
		},
		{
			name:          "invalid 2 letters price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.Price_2Letters = v; return p },
			swapPrice: func(params PriceParams) PriceParams {
				backup := params.Price_2Letters
				params.Price_2Letters = params.Price_3Letters
				params.Price_3Letters = backup
				return params
			},
		},
		{
			name:          "invalid 3 letters price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.Price_3Letters = v; return p },
			swapPrice: func(params PriceParams) PriceParams {
				backup := params.Price_3Letters
				params.Price_3Letters = params.Price_4Letters
				params.Price_4Letters = backup
				return params
			},
		},
		{
			name:          "invalid 4 letters price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.Price_4Letters = v; return p },
			swapPrice: func(params PriceParams) PriceParams {
				backup := params.Price_4Letters
				params.Price_4Letters = params.Price_5PlusLetters
				params.Price_5PlusLetters = backup
				return params
			},
		},
		{
			name:          "invalid 5+ letters price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.Price_5PlusLetters = v; return p },
		},
	}
	for _, tt := range testsInvalidPrice {
		t.Run(tt.name, func(t *testing.T) {
			err1 := tt.modifierPrice(validPriceParams, sdk.ZeroInt()).Validate()
			require.Error(t, err1)
			require.Contains(t, err1.Error(), "is zero")
			err2 := tt.modifierPrice(validPriceParams, sdk.NewInt(-1)).Validate()
			require.Error(t, err2)
			require.Contains(t, err2.Error(), "is negative")

			if tt.swapPrice != nil {
				err3 := tt.swapPrice(validPriceParams).Validate()
				require.Error(t, err3)
				require.Contains(t, err3.Error(), "must be greater than")
			}
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require.Error(t, validatePriceParams("hello world"))
		require.Error(t, validatePriceParams(&PriceParams{}), "not accept pointer")
	})
}

func TestMiscParams_Validate(t *testing.T) {
	var tests = []struct {
		name            string
		modifier        func(MiscParams) MiscParams
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "valid",
			modifier: func(p MiscParams) MiscParams { return p },
		},
		{
			name: "all = 1 is valid",
			modifier: func(p MiscParams) MiscParams {
				p.DaysOpenPurchaseOrderDuration = 1
				p.DaysPreservedClosedPurchaseOrder = 1
				p.GasCrudOpenPurchaseOrder = 1
				return p
			},
		},
		{
			name: "gas for CRUD = 0 is valid",
			modifier: func(p MiscParams) MiscParams {
				p.GasCrudOpenPurchaseOrder = 0
				return p
			},
		},
		{
			name:            "days OPO duration can not be zero",
			modifier:        func(p MiscParams) MiscParams { p.DaysOpenPurchaseOrderDuration = 0; return p },
			wantErr:         true,
			wantErrContains: "days OPO duration must be greater than 0",
		},
		{
			name:            "days OPO duration can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysOpenPurchaseOrderDuration = -1; return p },
			wantErr:         true,
			wantErrContains: "days OPO duration must be greater than 0",
		},
		{
			name:            "days preserved closed OPO duration can not be zero",
			modifier:        func(p MiscParams) MiscParams { p.DaysPreservedClosedPurchaseOrder = 0; return p },
			wantErr:         true,
			wantErrContains: "days preserved closed OPO must be greater than 0",
		},
		{
			name:            "days preserved closed OPO duration can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysPreservedClosedPurchaseOrder = -1; return p },
			wantErr:         true,
			wantErrContains: "days preserved closed OPO must be greater than 0",
		},
		{
			name:            "gas CRUD OPO can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.GasCrudOpenPurchaseOrder = -1; return p },
			wantErr:         true,
			wantErrContains: "gas for CRUD operations on OPO cannot be negative",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.modifier(DefaultMiscParams()).Validate()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require.Error(t, validateMiscParams("hello world"))
		require.Error(t, validateMiscParams(&MiscParams{}), "not accept pointer")
	})
}
func Test_validateEpochIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		i       interface{}
		wantErr bool
	}{
		{
			name: "'hour' is valid",
			i:    "hour",
		},
		{
			name: "'day' is valid",
			i:    "day",
		},
		{
			name: "'week' is valid",
			i:    "week",
		},
		{
			name:    "empty",
			i:       "",
			wantErr: true,
		},
		{
			name:    "not string",
			i:       1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Error(t, validateEpochIdentifier(tt.i))
			} else {
				require.NoError(t, validateEpochIdentifier(tt.i))
			}
		})
	}
}
