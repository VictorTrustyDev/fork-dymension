package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dymensionxyz/dymension/v3/x/incentives/types"
)

// TestDistributeToRollappGauges tests distributing rewards to rollapp gauges.
func (suite *KeeperTestSuite) TestDistributeToRollappGauges() {
	oneKRewardCoins := sdk.Coins{sdk.NewInt64Coin(defaultRewardDenom, 1000)}
	testCases := []struct {
		name        string
		rewards     sdk.Coins
		noSequencer bool
	}{
		{
			name:    "rollapp gauge with sequencer",
			rewards: oneKRewardCoins,
		},
		{
			name:    "rollapp gauge with no rewards",
			rewards: sdk.Coins{},
		},
		{
			name:        "rollapp gauge with no sequencer",
			rewards:     oneKRewardCoins,
			noSequencer: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			pubkey := ed25519.GenPrivKey().PubKey()
			addr := sdk.AccAddress(pubkey.Address())

			// create rollapp and check rollapp gauge created
			rollapp := suite.CreateDefaultRollapp(addr)
			res, err := suite.querier.RollappGauges(sdk.WrapSDKContext(suite.Ctx), &types.GaugesRequest{})
			suite.Require().NoError(err)
			suite.Require().Len(res.Data, 1)

			gaugeId := res.Data[0].Id

			var proposerAddr sdk.AccAddress
			if !tc.noSequencer {
				err = suite.CreateSequencer(suite.Ctx, rollapp, pubkey)
				suite.Require().NoError(err)
				proposerAddr = addr
			}

			if tc.rewards.Len() > 0 {
				suite.AddToGauge(tc.rewards, gaugeId)
			}

			gauge, err := suite.App.IncentivesKeeper.GetGaugeByID(suite.Ctx, gaugeId)
			suite.Require().NoError(err)
			_, err = suite.App.IncentivesKeeper.Distribute(suite.Ctx, []types.Gauge{*gauge})
			suite.Require().NoError(err)
			// check expected rewards against actual rewards received
			if !proposerAddr.Empty() {
				bal := suite.App.BankKeeper.GetAllBalances(suite.Ctx, proposerAddr)
				suite.Require().Equal(tc.rewards.String(), bal.String())
			}
		})
	}
}
