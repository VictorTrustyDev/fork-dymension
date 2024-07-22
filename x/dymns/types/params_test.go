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
		ChainsParams{
			AliasesByChainId: map[string]AliasesOfChainId{
				"dymension_1100-1": {
					Aliases: []string{"dym", "dymension"},
				},
			},
			CoinType60ChainIds: []string{"injective-1"},
		},
		MiscParams{
			BeginEpochHookIdentifier:     "b",
			EndEpochHookIdentifier:       "c",
			DaysGracePeriod:              30,
			DaysSellOrderDuration:        888,
			DaysPreservedClosedSellOrder: 88,
			DaysProhibitSell:             90,
		},
	)
	require.Equal(t, "a", params.Price.PriceDenom)
	require.Len(t, params.Chains.AliasesByChainId, 1)
	require.Len(t, params.Chains.AliasesByChainId["dymension_1100-1"].Aliases, 2)
	require.Equal(t, AliasesOfChainId{Aliases: []string{"dym", "dymension"}}, params.Chains.AliasesByChainId["dymension_1100-1"])
	require.Len(t, params.Chains.CoinType60ChainIds, 1)
	require.Equal(t, params.Chains.CoinType60ChainIds[0], "injective-1")
	require.Equal(t, "b", params.Misc.BeginEpochHookIdentifier)
	require.Equal(t, "c", params.Misc.EndEpochHookIdentifier)
	require.Equal(t, int32(30), params.Misc.DaysGracePeriod)
	require.Equal(t, int32(888), params.Misc.DaysSellOrderDuration)
	require.Equal(t, int32(88), params.Misc.DaysPreservedClosedSellOrder)
	require.Equal(t, int32(90), params.Misc.DaysProhibitSell)
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

func TestDefaultChainsParams(t *testing.T) {
	require.NoError(t, DefaultChainsParams().Validate())
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
	params.Chains.CoinType60ChainIds = []string{"invalid@"}
	require.Error(t, (&params).Validate())

	params = DefaultParams()
	params.Misc.DaysPreservedClosedSellOrder = 0
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

//goland:noinspection SpellCheckingInspection
func TestChainsParams_Validate(t *testing.T) {
	var tests = []struct {
		name            string
		modifier        func(params ChainsParams) ChainsParams
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "default is valid",
			modifier: func(p ChainsParams) ChainsParams { return p },
		},
		{
			name: "alias: empty is valid",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = nil
				return p
			},
		},
		{
			name: "coin-type-60-chains: empty is valid",
			modifier: func(p ChainsParams) ChainsParams {
				p.CoinType60ChainIds = nil
				return p
			},
		},
		{
			name: "alias: empty alias of chain is valid",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: nil},
				}
				return p
			},
		},
		{
			name: "alias: valid and correct alias",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym"}},
					"blumbus_100-1":    {Aliases: []string{"bb", "blumbus"}},
				}
				return p
			},
		},
		{
			name: "coin-type-60-chains: valid and correct alias",
			modifier: func(p ChainsParams) ChainsParams {
				p.CoinType60ChainIds = []string{"injective-1", "cronosmainnet_25-1"}
				return p
			},
		},
		{
			name: "alias: chain_id and alias must be unique among all, case alias & alias",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
					"dymension_1100-1": {Aliases: []string{"dym"}},
					"blumbus_100-1":    {Aliases: []string{"dym", "blumbus"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "chain ID and alias must unique among all",
		},
		{
			name: "coin-type-60-chains: chain_id must be unique among all",
			modifier: func(p ChainsParams) ChainsParams {
				p.CoinType60ChainIds = []string{"injective-1", "cronosmainnet_25-1", "injective-1"}
				return p
			},
			wantErr:         true,
			wantErrContains: "chain ID is not unique",
		},
		{
			name: "alias: chain_id and alias must be unique among all, case chain-id & alias",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
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
			name: "alias: reject if chain-id format is bad",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
					"dymension@":    {Aliases: []string{"dym"}},
					"blumbus_100-1": {Aliases: []string{"blumbus"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "is not well-formed",
		},
		{
			name: "coin-type-60-chains: reject if chain-id format is bad",
			modifier: func(p ChainsParams) ChainsParams {
				p.CoinType60ChainIds = []string{"injective@"}
				return p
			},
			wantErr:         true,
			wantErrContains: "is not well-formed",
		},
		{
			name: "coin-type-60-chains: reject if chain-id format is bad",
			modifier: func(p ChainsParams) ChainsParams {
				p.CoinType60ChainIds = []string{"in"}
				return p
			},
			wantErr:         true,
			wantErrContains: "must be at least 3 characters",
		},
		{
			name: "alias: reject if chain-id format is bad",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
					"d": {Aliases: []string{"dym"}},
				}
				return p
			},
			wantErr:         true,
			wantErrContains: "must be at least 3 characters",
		},
		{
			name: "alias: reject if alias format is bad",
			modifier: func(p ChainsParams) ChainsParams {
				p.AliasesByChainId = map[string]AliasesOfChainId{
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
			err := tt.modifier(DefaultChainsParams()).Validate()
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
		require.Error(t, validateChainsParams("hello world"))
		require.Error(t, validateChainsParams(&ChainsParams{}), "not accept pointer")
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
				p.DaysSellOrderDuration = 1
				p.DaysPreservedClosedSellOrder = 1
				return p
			},
		},
		{
			name: "epoch hour is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "hour"
				return p
			},
		},
		{
			name: "epoch hour is valid",
			modifier: func(p MiscParams) MiscParams {
				p.EndEpochHookIdentifier = "hour"
				return p
			},
		},
		{
			name: "epoch day is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "day"
				return p
			},
		},
		{
			name: "epoch day is valid",
			modifier: func(p MiscParams) MiscParams {
				p.EndEpochHookIdentifier = "day"
				return p
			},
		},
		{
			name: "epoch week is valid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "week"
				return p
			},
		},
		{
			name: "epoch week is valid",
			modifier: func(p MiscParams) MiscParams {
				p.EndEpochHookIdentifier = "week"
				return p
			},
		},
		{
			name: "other epoch is invalid",
			modifier: func(p MiscParams) MiscParams {
				p.BeginEpochHookIdentifier = "invalid"
				return p
			},
			wantErr:         true,
			wantErrContains: "invalid epoch identifier: invalid",
		},
		{
			name: "other epoch is invalid",
			modifier: func(p MiscParams) MiscParams {
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
			name:            "days SO duration can not be zero",
			modifier:        func(p MiscParams) MiscParams { p.DaysSellOrderDuration = 0; return p },
			wantErr:         true,
			wantErrContains: "Sell Orders duration must be greater than 0 day",
		},
		{
			name:            "days SO duration can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysSellOrderDuration = -1; return p },
			wantErr:         true,
			wantErrContains: "Sell Orders duration must be greater than 0 day",
		},
		{
			name:            "days preserved closed SO duration can not be zero",
			modifier:        func(p MiscParams) MiscParams { p.DaysPreservedClosedSellOrder = 0; return p },
			wantErr:         true,
			wantErrContains: "preserved closed Sell Orders must be greater than 0 day",
		},
		{
			name:            "days preserved closed SO duration can not be negative",
			modifier:        func(p MiscParams) MiscParams { p.DaysPreservedClosedSellOrder = -1; return p },
			wantErr:         true,
			wantErrContains: "preserved closed Sell Orders must be greater than 0 day",
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
