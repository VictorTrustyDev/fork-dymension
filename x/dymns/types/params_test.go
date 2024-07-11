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
		PriceParams{
			PriceDenom: "a",
		},
		AliasParams{
			ByChainId: map[string]AliasesOfChainId{
				"dymension_1100-1": {
					Aliases: []string{"dym", "dymension"},
				},
			},
		},
		MiscParams{
			BeginEpochHookIdentifier:         "b",
			EndEpochHookIdentifier:           "c",
			DaysGracePeriod:                  30,
			DaysOpenPurchaseOrderDuration:    888,
			DaysPreservedClosedPurchaseOrder: 88,
			DaysProhibitSell:                 90,
			GasCrudOpenPurchaseOrder:         8,
		},
	)
	require.Equal(t, "a", params.Price.PriceDenom)
	require.Len(t, params.Alias.ByChainId, 1)
	require.Len(t, params.Alias.ByChainId["dymension_1100-1"].Aliases, 2)
	require.Equal(t, AliasesOfChainId{Aliases: []string{"dym", "dymension"}}, params.Alias.ByChainId["dymension_1100-1"])
	require.Equal(t, "b", params.Misc.BeginEpochHookIdentifier)
	require.Equal(t, "c", params.Misc.EndEpochHookIdentifier)
	require.Equal(t, int32(30), params.Misc.DaysGracePeriod)
	require.Equal(t, int32(888), params.Misc.DaysOpenPurchaseOrderDuration)
	require.Equal(t, int32(88), params.Misc.DaysPreservedClosedPurchaseOrder)
	require.Equal(t, int32(90), params.Misc.DaysProhibitSell)
	require.Equal(t, int32(8), params.Misc.GasCrudOpenPurchaseOrder)
}

func TestDefaultPriceParams(t *testing.T) {
	priceParams := DefaultPriceParams()
	require.NoError(t, priceParams.Validate())

	t.Run("ensure setting is correct", func(t *testing.T) {
		i, ok := sdk.NewIntFromString("5" + "000000000000000000")
		require.True(t, ok)
		require.Equal(t, i, priceParams.Price_5PlusLetters)
	})

	t.Run("ensure price setting is at least 1 DYM", func(t *testing.T) {
		oneDym, ok := sdk.NewIntFromString("1" + "000000000000000000")
		require.True(t, ok)
		if oneDym.GT(priceParams.Price_5PlusLetters) {
			require.Fail(t, "price should be at least 1 DYM")
		}
		if oneDym.GT(priceParams.PriceExtends) {
			require.Fail(t, "price should be at least 1 DYM")
		}
	})
}

func TestDefaultAliasParams(t *testing.T) {
	require.NoError(t, DefaultAliasParams().Validate())
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
		PriceExtends:       sdk.NewInt(2),
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
		{
			name:          "invalid yearly extends price",
			modifierPrice: func(p PriceParams, v sdkmath.Int) PriceParams { p.PriceExtends = v; return p },
			swapPrice: func(params PriceParams) PriceParams {
				params.PriceExtends = params.Price_5PlusLetters.Add(params.Price_5PlusLetters)
				return params
			},
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
				require.Contains(t, err3.Error(), "must be greater")
			}
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require.Error(t, validatePriceParams("hello world"))
		require.Error(t, validatePriceParams(&PriceParams{}), "not accept pointer")
	})
}

func TestAliasParams_Validate(t *testing.T) {
	var tests = []struct {
		name            string
		modifier        func(AliasParams) AliasParams
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "default is valid",
			modifier: func(p AliasParams) AliasParams { return p },
		},
		{
			name: "empty is valid",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = nil
				return p
			},
		},
		{
			name: "empty alias of chain is valid",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: nil},
				}
				return p
			},
		},
		{
			name: "empty alias of chain is valid",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: nil},
				}
				return p
			},
		},
		{
			name: "valid and correct alias",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym"}},
					"blumbus_100-1":    {Aliases: []string{"bb", "blumbus"}},
				}
				return p
			},
		},
		{
			name: "chain_id and alias must be unique among all, case alias & alias",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym"}},
					"blumbus_100-1":    {Aliases: []string{"dym", "blumbus"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "chain ID and alias must unique among all",
		},
		{
			name: "chain_id and alias must be unique among all, case chain-id & alias",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym", "dymension"}},
					"blumbus_100-1":    {Aliases: []string{"blumbus", "cosmoshub"}},
					"cosmoshub":        {Aliases: []string{"cosmos"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "chain ID and alias must unique among all",
		},
		{
			name: "reject if chain-id format is bad",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension@":    {Aliases: []string{"dym"}},
					"blumbus_100-1": {Aliases: []string{"blumbus"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "is not well-formed",
		},
		{
			name: "reject if chain-id format is bad",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"d": {Aliases: []string{"dym"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "must be at least 3 characters",
		},
		{
			name: "reject if alias format is bad",
			modifier: func(p AliasParams) AliasParams {
				p.ByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym-dym"}},
					"blumbus_100-1":    {Aliases: []string{"blumbus"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "is not well-formed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.modifier(DefaultAliasParams()).Validate()
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
		require.Error(t, validateAliasParams("hello world"))
		require.Error(t, validateAliasParams(&AliasParams{}), "not accept pointer")
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
			name:     "default is valid",
			modifier: func(p MiscParams) MiscParams { return p },
		},
		{
			name: "all = 1 is valid",
			modifier: func(p MiscParams) MiscParams {
				p.DaysGracePeriod = 1
				p.DaysOpenPurchaseOrderDuration = 1
				p.DaysPreservedClosedPurchaseOrder = 1
				p.GasCrudOpenPurchaseOrder = 1
				return p
			},
		},
		{
			name: "epoch hour is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "hour"
				p.EndEpochHookIdentifier = "hour"
				return p
			},
		},
		{
			name: "epoch day is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "day"
				p.EndEpochHookIdentifier = "day"
				return p
			},
		},
		{
			name: "epoch week is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "week"
				p.EndEpochHookIdentifier = "week"
				return p
			},
		},
		{
			name: "other epoch is invalid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "invalid"
				p.EndEpochHookIdentifier = "invalid"
				return p
			},
			wantErr:         true,
			wantErrContains: "invalid epoch identifier: invalid",
		},
		{
			name: "grace period = 0 is valid",
			modifier: func(p MiscParams) MiscParams {
				p.DaysGracePeriod = 0
				return p
			},
		},
		{
			name:            "grace period can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysGracePeriod = -1; return p },
			wantErr:         true,
			wantErrContains: "days grace period cannot be negative",
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
			name:            "days prohibit sell can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysProhibitSell = -1; return p },
			wantErr:         true,
			wantErrContains: "prohibit sell must be at least 7 days",
		},
		{
			name:            "days prohibit sell can not be lower than 7",
			modifier:        func(p MiscParams) MiscParams { p.DaysProhibitSell = 6; return p },
			wantErr:         true,
			wantErrContains: "prohibit sell must be at least 7 days",
		},
		{
			name: "gas CRUD OPO = 0 is valid",
			modifier: func(p MiscParams) MiscParams {
				p.GasCrudOpenPurchaseOrder = 0
				return p
			},
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
