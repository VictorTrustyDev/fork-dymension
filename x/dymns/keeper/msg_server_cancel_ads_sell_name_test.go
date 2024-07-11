package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func Test_msgServer_CancelAdsSellName(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	msgServer := dymnskeeper.NewMsgServerImpl(dk)

	// setting block time
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Now().UTC(),
	})

	moduleParams := dk.GetParams(ctx)
	moduleParams.Misc.GasCrudOpenPurchaseOrder = 20_000_000
	err := dk.SetParams(ctx, moduleParams)
	require.NoError(t, err)

	futureEpoch := ctx.BlockTime().Add(time.Hour).Unix()

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	bidder := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   futureEpoch,
	}
	err = dk.SetDymName(ctx, dymName1)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "owned-by-1",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   futureEpoch,
	}
	err = dk.SetDymName(ctx, dymName2)
	require.NoError(t, err)

	dymNames := dk.GetAllNonExpiredDymNames(ctx, time.Now().Unix())
	require.Len(t, dymNames, 2)

	t.Run("do not process message that not pass basic validation", func(t *testing.T) {
		requireErrorFContains(t, func() error {
			resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
				Name:  "aaa",
				Owner: "dym1xxx",
			})

			require.Nil(t, resp)

			return err
		}, dymnstypes.ErrValidationFailed.Error())
	})

	t.Run("do not process message that refer to non-existing Dym-Name", func(t *testing.T) {
		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  "aaa",
			Owner: owner,
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), dymnstypes.ErrDymNameNotFound.Error())
	})

	t.Run("do not process that owner does not match", func(t *testing.T) {
		opo11 := dymnstypes.OpenPurchaseOrder{
			Name:      dymName1.Name,
			ExpireAt:  futureEpoch,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
		}()

		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  opo11.Name,
			Owner: "dym1ysjlrjcankjpmpxxzk27mvzhv25e266r80p5pv",
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "not the owner of the dym name")
	})

	t.Run("do not process for Dym-Name that does not have any OPO", func(t *testing.T) {
		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  dymName1.Name,
			Owner: owner,
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), dymnstypes.ErrOpenPurchaseOrderNotFound.Error())
	})

	t.Run("can not cancel expired", func(t *testing.T) {
		opo11 := dymnstypes.OpenPurchaseOrder{
			Name:      dymName1.Name,
			ExpireAt:  1,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
		}()

		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  opo11.Name,
			Owner: owner,
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "cannot cancel an expired order")
	})

	t.Run("can not cancel once bid placed", func(t *testing.T) {
		opo11 := dymnstypes.OpenPurchaseOrder{
			Name:     dymName1.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
			HighestBid: &dymnstypes.OpenPurchaseOrderBid{
				Bidder: bidder,
				Price:  dymnsutils.TestCoin(300),
			},
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
		}()

		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  opo11.Name,
			Owner: owner,
		})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "cannot cancel once bid placed")
	})

	t.Run("can will remove the active OPO expiration mapping record", func(t *testing.T) {
		apoe := dk.GetActiveOpenPurchaseOrdersExpiration(ctx)

		opo11 := dymnstypes.OpenPurchaseOrder{
			Name:     dymName1.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)
		apoe.ExpiryByName[opo11.Name] = opo11.ExpireAt

		opo12 := dymnstypes.OpenPurchaseOrder{
			Name:     dymName2.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo12)
		require.NoError(t, err)
		apoe.ExpiryByName[opo12.Name] = opo12.ExpireAt

		err = dk.SetActiveOpenPurchaseOrdersExpiration(ctx, apoe)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
			dk.DeleteOpenPurchaseOrder(ctx, opo12.Name)
		}()

		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  opo11.Name,
			Owner: owner,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, opo11.Name), "OPO should be removed from active")

		apoe = dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
		require.NotContains(t, apoe.ExpiryByName, opo11.Name)
		require.Contains(t, apoe.ExpiryByName, opo12.Name)
	})

	t.Run("can cancel if statisfied conditions", func(t *testing.T) {
		opo11 := dymnstypes.OpenPurchaseOrder{
			Name:     dymName1.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)

		opo12 := dymnstypes.OpenPurchaseOrder{
			Name:     dymName2.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo12)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
			dk.DeleteOpenPurchaseOrder(ctx, opo12.Name)
		}()

		resp, err := msgServer.CancelAdsSellName(sdk.WrapSDKContext(ctx), &dymnstypes.MsgCancelAdsSellName{
			Name:  opo11.Name,
			Owner: owner,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, opo11.Name), "OPO should be removed from active")
		require.NotNil(t, dk.GetOpenPurchaseOrder(ctx, dymName2.Name), "other records remaining as-is")

		list := dk.GetHistoricalOpenPurchaseOrders(ctx, opo11.Name)
		require.Empty(t, list, "no historical record should be added")

		require.GreaterOrEqual(t,
			ctx.GasMeter().GasConsumed(), sdk.Gas(moduleParams.Misc.GasCrudOpenPurchaseOrder),
			"should consume params gas",
		)
	})
}
