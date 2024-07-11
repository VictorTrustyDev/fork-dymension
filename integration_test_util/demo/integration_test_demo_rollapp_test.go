package demo

//goland:noinspection SpellCheckingInspection
import (
	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"
)

//goland:noinspection SpellCheckingInspection

func (suite *DemoTestSuite) Test_Create_RollApps() {
	creator := suite.CITS.WalletAccounts.Number(1)

	for r := 1; r <= 5; r++ {
		rollApp := suite.CITS.GenSampleRollApp(creator)

		_, _, err := suite.CITS.DeliverTx(suite.Ctx(), creator, nil, &rollapptypes.MsgCreateRollapp{
			Creator:               rollApp.Creator,
			RollappId:             rollApp.RollappId,
			MaxSequencers:         rollApp.MaxSequencers,
			PermissionedAddresses: rollApp.PermissionedAddresses,
		})
		suite.Require().NoError(err)
		suite.Commit()

		res, err := suite.CITS.QueryClients.RollApp.RollappAll(suite.Ctx(), &rollapptypes.QueryAllRollappRequest{})
		suite.Require().NoError(err)
		suite.Require().NotNil(res)
		suite.Len(res.Rollapp, r)
	}
}

func (suite *DemoTestSuite) Test_QC_RollApps_At_Different_Blocks() {
	creator := suite.CITS.WalletAccounts.Number(1)

	wantHistoricalRecords := make(map[int64]int)

	wantHistoricalRecords[suite.CITS.GetLatestBlockHeight()] = 0

	suite.Commit()

	for r := 1; r <= 5; r++ {
		rollApp := suite.CITS.GenSampleRollApp(creator)

		_, _, err := suite.CITS.DeliverTx(suite.Ctx(), creator, nil, &rollapptypes.MsgCreateRollapp{
			Creator:       rollApp.Creator,
			RollappId:     rollApp.RollappId,
			MaxSequencers: rollApp.MaxSequencers,
		})
		suite.Require().NoError(err)
		suite.Commit()

		wantHistoricalRecords[suite.CITS.GetLatestBlockHeight()] = r
		suite.Commit() // shift one block to keep state

		res, err := suite.CITS.QueryClients.RollApp.RollappAll(suite.Ctx(), &rollapptypes.QueryAllRollappRequest{})
		suite.Require().NoError(err)
		suite.Require().NotNil(res)
		suite.Len(res.Rollapp, r)
	}

	for height, wantCount := range wantHistoricalRecords {
		res, err := suite.CITS.QueryClientsAt(height).RollApp.RollappAll(suite.Ctx(), &rollapptypes.QueryAllRollappRequest{})
		suite.Require().NoError(err)
		suite.Require().NotNil(res)
		suite.Lenf(res.Rollapp, wantCount, "want %d RollApps at height %d but got %d", wantCount, height, len(res.Rollapp))
	}
}
