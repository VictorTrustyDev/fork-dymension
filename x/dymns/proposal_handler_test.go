package dymns_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	"github.com/dymensionxyz/dymension/v3/x/dymns"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ProposalHandler(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	proposalHandler := dymns.NewDymNsProposalHandler(dk)

	t.Run("pass - can process proposal", func(t *testing.T) {
		err := proposalHandler(ctx, &dymnstypes.MigrateChainIdsProposal{
			Title:       "T",
			Description: "D",
			Replacement: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("fail - can not process unknown proposal", func(t *testing.T) {
		//goland:noinspection GoDeprecation
		err := proposalHandler(ctx, &distrtypes.CommunityPoolSpendProposal{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unrecognized dymns proposal content type")
	})
}

func Test_ProposalHandler_MigrateChainIdsProposal(t *testing.T) {
	tests := []struct {
		name            string
		additionalSetup func(ctx sdk.Context, dk dymnskeeper.Keeper)
		proposal        dymnstypes.MigrateChainIdsProposal
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "pass - migration successfully",
			proposal: dymnstypes.MigrateChainIdsProposal{
				Title:       "T",
				Description: "D",
				Replacement: []dymnstypes.MigrateChainId{
					{
						PreviousChainId: "cosmoshub-3",
						NewChainId:      "cosmoshub-4",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail - reject invalid proposal content",
			proposal: dymnstypes.MigrateChainIdsProposal{
				Title:       "T",
				Description: "D",
				Replacement: []dymnstypes.MigrateChainId{
					{
						PreviousChainId: "",
						NewChainId:      "cosmoshub-4",
					},
				},
			},
			wantErr:         true,
			wantErrContains: "previous chain id cannot be empty",
		},
		{
			name: "fail - returns error if migration failed",
			additionalSetup: func(ctx sdk.Context, dk dymnskeeper.Keeper) {
				moduleParams := dk.GetParams(ctx)
				moduleParams.Chains.AliasesByChainId = map[string]dymnstypes.AliasesOfChainId{
					"dymension_1100-1": {
						Aliases: []string{"dym"},
					},
					"blumbus_111-1": {
						Aliases: []string{"bb"},
					},
				}
				require.NoError(t, dk.SetParams(ctx, moduleParams))
			},
			proposal: dymnstypes.MigrateChainIdsProposal{
				Title:       "T",
				Description: "D",
				Replacement: []dymnstypes.MigrateChainId{
					{
						PreviousChainId: "blumbus_111-1",
						NewChainId:      "dym", // collision with alias of dymension_1100-1
					},
				},
			},
			wantErr:         true,
			wantErrContains: "chain ID and alias must unique among all",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, _, _, ctx := testkeeper.DymNSKeeper(t)

			if tt.additionalSetup != nil {
				tt.additionalSetup(ctx, dk)
			}

			proposalHandler := dymns.NewDymNsProposalHandler(dk)

			err := proposalHandler(ctx, &tt.proposal)
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
