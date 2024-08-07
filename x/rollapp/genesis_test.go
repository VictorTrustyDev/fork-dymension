package rollapp_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	"github.com/dymensionxyz/dymension/v3/testutil/nullify"
	"github.com/dymensionxyz/dymension/v3/x/rollapp"
	"github.com/dymensionxyz/dymension/v3/x/rollapp/types"
)

func TestInitExportGenesis(t *testing.T) {
	const (
		rollappID1 = "rollapp_1234-1"
		rollappID2 = "rollupp_1235-2"
	)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		RollappList: []types.Rollapp{
			{
				RollappId: rollappID1,
			},
			{
				RollappId: rollappID2,
			},
		},
		StateInfoList: []types.StateInfo{
			{
				StateInfoIndex: types.StateInfoIndex{
					RollappId: rollappID1,
					Index:     0,
				},
			},
			{
				StateInfoIndex: types.StateInfoIndex{
					RollappId: rollappID2,
					Index:     1,
				},
			},
		},
		LatestStateInfoIndexList: []types.StateInfoIndex{
			{
				RollappId: rollappID1,
			},
			{
				RollappId: rollappID2,
			},
		},
		BlockHeightToFinalizationQueueList: []types.BlockHeightToFinalizationQueue{
			{
				CreationHeight: 0,
			},
			{
				CreationHeight: 1,
			},
		},
	}

	k, ctx := keepertest.RollappKeeper(t)
	rollapp.InitGenesis(ctx, *k, genesisState)
	got := rollapp.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(genesisState)
	nullify.Fill(*got)

	require.ElementsMatch(t, genesisState.RollappList, got.RollappList)
	require.ElementsMatch(t, genesisState.StateInfoList, got.StateInfoList)
	require.ElementsMatch(t, genesisState.LatestStateInfoIndexList, got.LatestStateInfoIndexList)
	require.ElementsMatch(t, genesisState.BlockHeightToFinalizationQueueList, got.BlockHeightToFinalizationQueueList)
}
