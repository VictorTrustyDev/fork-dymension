package cli

import (
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_parseMigrateChainIdsProposal(t *testing.T) {
	testCases := []struct {
		name         string
		metadataFile string
		wantErr      bool
		want         []dymnstypes.MigrateChainId
	}{
		{
			name:         "fail - invalid file name",
			metadataFile: "",
			wantErr:      true,
		},
		{
			name:         "fail - invalid content",
			metadataFile: "test_proposals/mcid_invalid_update_chain_id_proposal_test.json",
			wantErr:      true,
		},
		{
			name:         "pass - update single",
			metadataFile: "test_proposals/mcid_update_single_chain_id_proposal_test.json",
			wantErr:      false,
			want: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
			},
		},
		{
			name:         "pass - update multiple",
			metadataFile: "test_proposals/mcid_update_multiple_chain_ids_proposal_test.json",
			wantErr:      false,
			want: []dymnstypes.MigrateChainId{
				{
					PreviousChainId: "cosmoshub-3",
					NewChainId:      "cosmoshub-4",
				},
				{
					PreviousChainId: "columbus-4",
					NewChainId:      "columbus-5",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proposal, err := parseMigrateChainIdsProposal(types.AminoCdc, tc.metadataFile)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, proposal.Replacement, len(tc.want))
			if !reflect.DeepEqual(tc.want, proposal.Replacement) {
				t.Errorf("expected: %v, got: %v", tc.want, proposal.Replacement)
			}
		})
	}
}
