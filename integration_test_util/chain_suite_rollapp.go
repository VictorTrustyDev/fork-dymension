package integration_test_util

//goland:noinspection SpellCheckingInspection
import (
	"fmt"
	"github.com/cometbft/cometbft/libs/rand"
	itutiltypes "github.com/dymensionxyz/dymension/v3/integration_test_util/types"
	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"
)

// GenSampleRollApp generates a sample rollapp with non-existing chain-id.
func (suite *ChainIntegrationTestSuite) GenSampleRollApp(creator *itutiltypes.TestAccount) rollapptypes.Rollapp {
	for {
		seed := int32(rand.Uint16()) + 101
		rollAppId := fmt.Sprintf("rollapp_%d-1", seed)
		_, err := suite.ChainApp.RollAppKeeper().Rollapp(suite.CurrentContext, &rollapptypes.QueryGetRollappRequest{
			RollappId: rollAppId,
		})
		if err == nil {
			continue
		}

		baseDenom := fmt.Sprintf("arax%d", seed)
		return rollapptypes.Rollapp{
			RollappId:             rollAppId,
			Creator:               creator.GetCosmosAddress().String(),
			Version:               0,
			MaxSequencers:         1,
			PermissionedAddresses: nil,
			RegisteredDenoms:      []string{baseDenom},
		}
	}
}
