package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/x/sequencer/types"
)

func (suite *SequencerTestSuite) TestDecreaseBond() {
	suite.SetupTest()
	bondDenom := types.DefaultParams().MinBond.Denom
	rollappId, pk := suite.CreateDefaultRollapp()
	// setup a default sequencer with has minBond + 20token
	defaultSequencerAddress := suite.CreateSequencerWithBond(suite.Ctx, rollappId, bond.AddAmount(sdk.NewInt(20)), pk)
	// setup an unbonded sequencer
	unbondedPk := ed25519.GenPrivKey().PubKey()
	unbondedSequencerAddress := suite.CreateSequencer(suite.Ctx, rollappId, unbondedPk)
	unbondedSequencer, _ := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, unbondedSequencerAddress)
	unbondedSequencer.Status = types.Unbonded
	suite.App.SequencerKeeper.UpdateSequencer(suite.Ctx, &unbondedSequencer, unbondedSequencer.Status)
	// setup a jailed sequencer
	jailedPk := ed25519.GenPrivKey().PubKey()
	jailedSequencerAddress := suite.CreateSequencer(suite.Ctx, rollappId, jailedPk)
	jailedSequencer, _ := suite.App.SequencerKeeper.GetSequencer(suite.Ctx, jailedSequencerAddress)
	jailedSequencer.Jailed = true
	suite.App.SequencerKeeper.UpdateSequencer(suite.Ctx, &jailedSequencer, jailedSequencer.Status)

	testCase := []struct {
		name        string
		msg         types.MsgDecreaseBond
		expectedErr error
	}{
		{
			name: "invalid sequencer",
			msg: types.MsgDecreaseBond{
				Creator:        "invalid_address",
				DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
			},
			expectedErr: types.ErrUnknownSequencer,
		},
		{
			name: "sequencer is not bonded",
			msg: types.MsgDecreaseBond{
				Creator:        unbondedSequencerAddress,
				DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
			},
			expectedErr: types.ErrInvalidSequencerStatus,
		},
		{
			name: "decreased bond value to less than minimum bond value",
			msg: types.MsgDecreaseBond{
				Creator:        defaultSequencerAddress,
				DecreaseAmount: sdk.NewInt64Coin(bondDenom, 100),
			},
			expectedErr: types.ErrInsufficientBond,
		},
		{
			name: "trying to decrease more bond than they have tokens bonded",
			msg: types.MsgDecreaseBond{
				Creator:        defaultSequencerAddress,
				DecreaseAmount: bond.AddAmount(sdk.NewInt(30)),
			},
			expectedErr: types.ErrInsufficientBond,
		},
		{
			name: "valid decrease bond",
			msg: types.MsgDecreaseBond{
				Creator:        defaultSequencerAddress,
				DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
			},
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			resp, err := suite.msgServer.DecreaseBond(suite.Ctx, &tc.msg)
			if tc.expectedErr != nil {
				suite.Require().ErrorIs(err, tc.expectedErr)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)
				expectedCompletionTime := suite.Ctx.BlockHeader().Time.Add(suite.App.SequencerKeeper.UnbondingTime(suite.Ctx))
				suite.Require().Equal(expectedCompletionTime, resp.CompletionTime)
				// check if the unbonding is set correctly
				bondReductionIDs := suite.App.SequencerKeeper.GetMatureDecreasingBondIDs(suite.Ctx, expectedCompletionTime)
				suite.Require().Len(bondReductionIDs, 1)
				bondReduction, found := suite.App.SequencerKeeper.GetBondReduction(suite.Ctx, bondReductionIDs[0])
				suite.Require().True(found)
				suite.Require().Equal(tc.msg.Creator, bondReduction.SequencerAddress)
				suite.Require().Equal(tc.msg.DecreaseAmount, bondReduction.DecreaseBondAmount)
			}
		})
	}
}

func (suite *SequencerTestSuite) TestDecreaseBond_BondDecreaseInProgress() {
	suite.SetupTest()
	bondDenom := types.DefaultParams().MinBond.Denom
	rollappId, pk := suite.CreateDefaultRollapp()
	// setup a default sequencer with has minBond + 20token
	defaultSequencerAddress := suite.CreateSequencerWithBond(suite.Ctx, rollappId, bond.AddAmount(sdk.NewInt(20)), pk)
	// decrease the bond of the sequencer
	_, err := suite.msgServer.DecreaseBond(suite.Ctx, &types.MsgDecreaseBond{
		Creator:        defaultSequencerAddress,
		DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
	})
	suite.Require().NoError(err)
	// try to decrease the bond again - should be fine as still not below minbond
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1).WithBlockTime(suite.Ctx.BlockTime().Add(10))
	_, err = suite.msgServer.DecreaseBond(suite.Ctx, &types.MsgDecreaseBond{
		Creator:        defaultSequencerAddress,
		DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
	})
	suite.Require().NoError(err)
	// try to decrease the bond again - should err as below minbond
	suite.Ctx = suite.Ctx.WithBlockHeight(suite.Ctx.BlockHeight() + 1).WithBlockTime(suite.Ctx.BlockTime().Add(10))
	_, err = suite.msgServer.DecreaseBond(suite.Ctx, &types.MsgDecreaseBond{
		Creator:        defaultSequencerAddress,
		DecreaseAmount: sdk.NewInt64Coin(bondDenom, 10),
	})
	suite.Require().ErrorIs(err, types.ErrInsufficientBond)
}
