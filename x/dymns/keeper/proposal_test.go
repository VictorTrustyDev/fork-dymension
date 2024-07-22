package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	rollappkeeper "github.com/dymensionxyz/dymension/v3/x/rollapp/keeper"
	"github.com/stretchr/testify/require"
	"reflect"
	"sort"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestKeeper_MigrateChainIds(t *testing.T) {
	now := time.Now().UTC()
	const chainId = "dymension_1100-1"

	const addr1 = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const addr2 = "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	const addrCosmos3 = "cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a"

	tests := []struct {
		name                       string
		dymNames                   []dymnstypes.DymName
		replacement                []dymnstypes.MigrateChainId
		chainsAliasParams          map[string]dymnstypes.AliasesOfChainId
		chainsCoinType60Params     []string
		additionalSetup            func(ctx sdk.Context, dk dymnskeeper.Keeper, rk rollappkeeper.Keeper)
		wantErr                    bool
		wantErrContains            string
		wantDymNames               []dymnstypes.DymName
		wantChainsAliasParams      map[string]dymnstypes.AliasesOfChainId
		wantChainsCoinType60Params []string
	}{
		{
			name: "pass - can migrate",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-3",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
				},
			},
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "blumbus_111-2",
				},
			},
			chainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				"cosmoshub-3": {
					Aliases: []string{"cosmos"},
				},
			},
			chainsCoinType60Params: []string{"blumbus_111-1"},
			additionalSetup:        nil,
			wantErr:                false,
			wantDymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-4",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
				},
			},
			wantChainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				"cosmoshub-4": {
					Aliases: []string{"cosmos"},
				},
			},
			wantChainsCoinType60Params: []string{"blumbus_111-2"},
		},
		{
			name: "pass - can migrate params alias chain-id",
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "blumbus_111-2",
				},
			},
			chainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				"cosmoshub-3": {
					Aliases: []string{"cosmos"},
				},
				"blumbus_111-1": {
					Aliases: []string{"bb"},
				},
				chainId: {
					Aliases: []string{"dym"},
				},
			},
			wantErr: false,
			wantChainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				"cosmoshub-4": {
					Aliases: []string{"cosmos"},
				},
				"blumbus_111-2": {
					Aliases: []string{"bb"},
				},
				chainId: {
					Aliases: []string{"dym"},
				},
			},
		},
		{
			name: "pass - can migrate params coin type 60 chain-id",
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "blumbus_111-2",
				},
			},
			chainsCoinType60Params:     []string{"blumbus_111-1", "nim_1122-1"},
			wantErr:                    false,
			wantChainsCoinType60Params: []string{"blumbus_111-2", "nim_1122-1"},
		},
		{
			name: "pass - can Dym-Name",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-3",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-3",
							Path:    "",
							Value:   addrCosmos3,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "blumbus_111-1",
							Path:    "",
							Value:   addr2,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "froopyland_100-1",
							Path:    "",
							Value:   addr2,
						},
					},
				},
			},
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "blumbus_111-2",
				},
			},
			wantErr: false,
			wantDymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-4",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   addrCosmos3,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "blumbus_111-2",
							Path:    "",
							Value:   addr2,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "froopyland_100-1",
							Path:    "",
							Value:   addr2,
						},
					},
				},
			},
		},
		{
			name: "pass - ignore expired Dym-Name",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-3",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() - 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-3",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
			},
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
			wantErr: false,
			wantDymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-4",
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() - 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "cosmoshub-3", // keep
						Path:    "",
						Value:   addrCosmos3,
					}},
				},
			},
		},
		{
			name: "fail - should stop if can not migrate params",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_111-1",
						Path:    "",
						Value:   addr1,
					}},
				},
			},
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "dym", // collision with alias
				},
			},
			chainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				chainId: {
					Aliases: []string{"dym"}, // collision with new chain-id
				},
				"blumbus_111-1": {
					Aliases: []string{"bb"},
				},
			},
			chainsCoinType60Params: []string{"blumbus_111-1"},
			wantErr:                true,
			wantErrContains:        "chains params: alias: chain ID and alias must unique among all",
			wantDymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_111-1", // not updated
						Path:    "",
						Value:   addr1,
					}},
				},
			},
			wantChainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				// not changed
				chainId: {
					Aliases: []string{"dym"},
				},
				"blumbus_111-1": {
					Aliases: []string{"bb"},
				},
			},
			wantChainsCoinType60Params: []string{
				// not changed
				"blumbus_111-1",
			},
		},
		{
			name: "fail - should stop if new params does not valid",
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "blumbus_111-1",
					NewChainId:      "dym", // collision with alias
				},
			},
			chainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				chainId: {
					Aliases: []string{"dym"}, // collision with new chain-id
				},
				"blumbus_111-1": {
					Aliases: []string{"bb"},
				},
			},
			wantErr:         true,
			wantErrContains: "chains params: alias: chain ID and alias must unique among all",
			wantChainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				// not changed
				chainId: {
					Aliases: []string{"dym"},
				},
				"blumbus_111-1": {
					Aliases: []string{"bb"},
				},
			},
		},
		{
			name:     "pass - should complete even tho nothing to update",
			dymNames: nil,
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
			wantErr: false,
		},
		{
			name: "pass - skip migrate alias if new chain-id present, just remove",
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
			chainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				chainId: {
					Aliases: []string{"dym"},
				},
				"cosmoshub-3": {
					Aliases: []string{"cosmos3"},
				},
				"cosmoshub-4": {
					Aliases: []string{"cosmos4"},
				},
			},
			wantErr: false,
			wantChainsAliasParams: map[string]dymnstypes.AliasesOfChainId{
				chainId: {
					Aliases: []string{"dym"},
				},
				"cosmoshub-4": {
					Aliases: []string{"cosmos4"},
				},
			},
		},
		{
			name: "pass - skip migrate coin-type-60 chain-ids if new chain-id present, just remove",
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "nim_1122-1",
					NewChainId:      "nim_1122-2",
				},
			},
			chainsCoinType60Params:     []string{"blumbus_1122-1", "nim_1122-1", "nim_1122-2"},
			wantErr:                    false,
			wantChainsCoinType60Params: []string{"blumbus_1122-1", "nim_1122-2"},
		},
		{
			name: "pass - skip migrate Dym-Name if new record does not pass validation",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						// migrate this will cause non-unique config
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-3",
							Path:    "",
							Value:   addrCosmos3,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-3",
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
				{
					Name:       "c",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
			},
			replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
			wantErr: false,
			wantDymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-3", // keep
							Path:    "",
							Value:   addrCosmos3,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
				{
					Name:       "b",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4", // migrated
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
				{
					Name:       "c",
					Owner:      addr1,
					Controller: addr1,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   addrCosmos3,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, _, rk, ctx := testkeeper.DymNSKeeper(t)

			ctx = ctx.WithBlockTime(now).WithChainID(chainId)

			moduleParams := dk.GetParams(ctx)
			moduleParams.Chains.AliasesByChainId = tt.chainsAliasParams
			moduleParams.Chains.CoinType60ChainIds = tt.chainsCoinType60Params
			require.NoError(t, dk.SetParams(ctx, moduleParams))

			for _, dymName := range tt.dymNames {
				setDymNameWithFunctionsAfter(ctx, dymName, t, dk)
			}

			if tt.additionalSetup != nil {
				tt.additionalSetup(ctx, dk, rk)
			}

			err := dk.MigrateChainIds(ctx, tt.replacement)

			defer func() {
				laterModuleParams := dk.GetParams(ctx)
				if len(tt.wantChainsAliasParams) > 0 || len(laterModuleParams.Chains.AliasesByChainId) > 0 {
					if !reflect.DeepEqual(tt.wantChainsAliasParams, laterModuleParams.Chains.AliasesByChainId) {
						t.Errorf("alias: want %v, got %v", tt.wantChainsAliasParams, laterModuleParams.Chains.AliasesByChainId)
					}
				}
				if len(tt.wantChainsCoinType60Params) > 0 || len(laterModuleParams.Chains.CoinType60ChainIds) > 0 {
					sort.Strings(tt.wantChainsCoinType60Params)
					sort.Strings(laterModuleParams.Chains.CoinType60ChainIds)
					require.Equal(t, tt.wantChainsCoinType60Params, laterModuleParams.Chains.CoinType60ChainIds)
				}
			}()

			defer func() {
				for _, wantDymName := range tt.wantDymNames {
					laterDymName := dk.GetDymName(ctx, wantDymName.Name)
					require.NotNil(t, laterDymName)
					if !reflect.DeepEqual(wantDymName.Configs, laterDymName.Configs) {
						t.Errorf("dym name config: want %v, got %v", wantDymName.Configs, laterDymName.Configs)
					}
					if !reflect.DeepEqual(wantDymName, *laterDymName) {
						t.Errorf("dym name: want %v, got %v", wantDymName, *laterDymName)
					}
				}
			}()

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				return
			}

			require.NoError(t, err)
		})
	}
}
