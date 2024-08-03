package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	rollappkeeper "github.com/dymensionxyz/dymension/v3/x/rollapp/keeper"

	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func TestGetSetParams(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	params := dymnstypes.DefaultParams()

	err := dk.SetParams(ctx, params)
	require.NoError(t, err)

	require.Equal(t, params, dk.GetParams(ctx))

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Misc.BeginEpochHookIdentifier = ""
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Price.PriceDenom = ""
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
			{
				ChainId: "@",
				Aliases: nil,
			},
		}
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Misc.GracePeriodDuration = -999 * time.Hour
		require.Error(t, dk.SetParams(ctx, params))
	})
}

func TestKeeper_CheckChainIsCoinType60ByChainId(t *testing.T) {
	dk, _, rk, ctx := testkeeper.DymNSKeeper(t)

	const chainIdInjective = "injective-1"

	params := dk.GetParams(ctx)

	t.Run("roll-app is coin-type 60", func(t *testing.T) {
		rollApp1 := rollapptypes.Rollapp{
			RollappId: "ra_1-1",
			Creator:   testAddr(0).bech32(),
		}
		rk.SetRollapp(ctx, rollApp1)

		require.True(t, dk.CheckChainIsCoinType60ByChainId(ctx, rollApp1.RollappId))
	})

	t.Run("chain-id in params is coin-type 60", func(t *testing.T) {
		params.Chains.CoinType60ChainIds = []string{chainIdInjective}
		err := dk.SetParams(ctx, params)
		require.NoError(t, err)

		require.True(t, dk.CheckChainIsCoinType60ByChainId(ctx, chainIdInjective))
	})

	t.Run("otherwise not coin-type 60", func(t *testing.T) {
		require.False(t, dk.CheckChainIsCoinType60ByChainId(ctx, "cosmoshub-4"))
	})

	t.Run("chain-id not in params is not coin-type 60 regardless actual", func(t *testing.T) {
		params.Chains.CoinType60ChainIds = nil
		err := dk.SetParams(ctx, params)
		require.NoError(t, err)

		require.False(t, dk.CheckChainIsCoinType60ByChainId(ctx, chainIdInjective))
	})
}

func TestKeeper_CanUseAliasForNewRegistration(t *testing.T) {
	tests := []struct {
		name            string
		alias           string
		preSetup        func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper)
		wantErr         bool
		wantErrContains string
		want            bool
	}{
		{
			name:    "pass - can check",
			alias:   "a",
			wantErr: false,
			want:    true,
		},
		{
			name:            "fail - reject bad alias",
			alias:           "@",
			wantErr:         true,
			wantErrContains: "alias candidate: invalid argument",
		},
		{
			name:  "pass - returns as free if neither in Params or Roll-App",
			alias: "free",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "dymension_1100-1",
						Aliases: []string{"dym"},
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)

				rk.SetRollapp(ctx, rollapptypes.Rollapp{
					RollappId: "rollapp_1-1",
					Creator:   testAddr(1).bech32(),
				})
				err = dk.SetAliasForRollAppId(ctx, "rollapp_1-1", "ra")
				require.NoError(t, err)
			},
			wantErr: false,
			want:    true,
		},
		{
			name:    "pass - returns as free if no params, no Roll-App",
			alias:   "free",
			wantErr: false,
			want:    true,
		},
		{
			name:  "pass - returns as NOT free if reserved in Params",
			alias: "dymension",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "dymension_1100-1",
						Aliases: []string{"dym", "dymension"},
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)

				rk.SetRollapp(ctx, rollapptypes.Rollapp{
					RollappId: "rollapp_1-1",
					Creator:   testAddr(1).bech32(),
				})
				err = dk.SetAliasForRollAppId(ctx, "rollapp_1-1", "ra")
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if reserved in Params as chain-id, without alias",
			alias: "zeta",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "zeta",
						Aliases: nil,
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if reserved in RollApp",
			alias: "ra",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "dymension_1100-1",
						Aliases: []string{"dym", "dymension"},
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)

				rk.SetRollapp(ctx, rollapptypes.Rollapp{
					RollappId: "rollapp_1-1",
					Creator:   testAddr(1).bech32(),
				})
				err = dk.SetAliasForRollAppId(ctx, "rollapp_1-1", "ra")
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if reserved in both Params and RollApp",
			alias: "dym",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "dymension_1100-1",
						Aliases: []string{"dym", "dymension"},
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)

				rk.SetRollapp(ctx, rollapptypes.Rollapp{
					RollappId: "dymension_1-1",
					Creator:   testAddr(1).bech32(),
				})
				err = dk.SetAliasForRollAppId(ctx, "dymension_1-1", "dym")
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if it is a Chain-ID in params mapping",
			alias: "bridge",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.AliasesOfChainIds = []dymnstypes.AliasesOfChainId{
					{
						ChainId: "bridge",
						Aliases: []string{"b"},
					},
				}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if it is a Chain-ID in params list Coin-Type-60",
			alias: "bridge",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				params := dk.GetParams(ctx)
				params.Chains.CoinType60ChainIds = []string{"bridge"}
				err := dk.SetParams(ctx, params)
				require.NoError(t, err)
			},
			wantErr: false,
			want:    false,
		},
		{
			name:  "pass - returns as NOT free if it is a RollApp-ID",
			alias: "bridge",
			preSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper) {
				rk.SetRollapp(ctx, rollapptypes.Rollapp{
					RollappId: "bridge",
					Creator:   testAddr(1).bech32(),
				})
				err := dk.SetAliasForRollAppId(ctx, "bridge", "b")
				require.NoError(t, err)

				require.True(t, dk.IsRollAppId(ctx, "bridge"))
			},
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, _, rk, ctx := testkeeper.DymNSKeeper(t)

			if tt.preSetup != nil {
				tt.preSetup(ctx, dk, rk)
			}

			can, err := dk.CanUseAliasForNewRegistration(ctx, tt.alias)
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)

				require.False(t, can)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, can)
		})
	}
}
