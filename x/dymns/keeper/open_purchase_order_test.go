package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetSetDeleteOpenPurchaseOrder(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	t.Run("reject invalid opo", func(t *testing.T) {
		err := dk.SetOpenPurchaseOrder(ctx, dymnstypes.OpenPurchaseOrder{})
		require.Error(t, err)
	})

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   1,
	}
	err := dk.SetDymName(ctx, dymName1)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "name2",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   1,
	}
	err = dk.SetDymName(ctx, dymName2)
	require.NoError(t, err)

	opo1 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  1,
		MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
		SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(300)),
		HighestBid: &dymnstypes.OpenPurchaseOrderBid{
			Bidder: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
			Price:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(200)),
		},
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo1)
	require.NoError(t, err)
	t.Run("opo1 should be equals to original", func(t *testing.T) {
		require.EqualValues(t, opo1, *dk.GetOpenPurchaseOrder(ctx, opo1.Name))
	})
	t.Run("opo list should have length 1", func(t *testing.T) {
		require.Len(t, dk.GetAllOpenPurchaseOrders(ctx), 1)
	})
	t.Run("event should be fired", func(t *testing.T) {
		events := ctx.EventManager().Events()
		require.NotEmpty(t, events)

		for _, event := range events {
			if event.Type == dymnstypes.EventTypeDymNameOpenPurchaseOrder {
				return
			}
		}

		t.Errorf("event %s not found", dymnstypes.EventTypeDymNameOpenPurchaseOrder)
	})

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName2.Name,
		ExpireAt:  1,
		MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
		SellPrice: sdk.NewCoin(params.BaseDenom, sdk.ZeroInt()),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo2)
	require.NoError(t, err)
	t.Run("opo2 should be equals to original", func(t *testing.T) {
		require.EqualValues(t, opo2, *dk.GetOpenPurchaseOrder(ctx, opo2.Name))
	})
	t.Run("opo list should have length 2", func(t *testing.T) {
		require.Len(t, dk.GetAllOpenPurchaseOrders(ctx), 2)
	})

	dk.DeleteOpenPurchaseOrder(ctx, opo1)
	t.Run("opo1 should be nil", func(t *testing.T) {
		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, opo1.Name))
	})
	t.Run("opo list should have length 1", func(t *testing.T) {
		list := dk.GetAllOpenPurchaseOrders(ctx)
		require.Len(t, list, 1)
		require.Equal(t, opo2.Name, list[0].Name)
	})

	t.Run("non-exists returns nil", func(t *testing.T) {
		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, "non-exists"))
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_MoveOpenPurchaseOrderToHistorical(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	epoch := time.Now().UTC().Add(time.Hour).Unix()

	// setting block time
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Now().UTC(),
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   epoch,
	}
	err := dk.SetDymName(ctx, dymName1)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "owned-by-1",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   epoch,
	}
	err = dk.SetDymName(ctx, dymName2)
	require.NoError(t, err)

	dymNames := dk.GetAllNonExpiredDymNames(ctx, time.Now().UTC().Unix())
	require.Len(t, dymNames, 2)

	opo11 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  1,
		MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
		SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(300)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo11)
	require.NoError(t, err)

	t.Run("should able to move", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo11.Name)
		require.NoError(t, err)
	})

	t.Run("moved OPO should be removed from active", func(t *testing.T) {
		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, opo11.Name))
	})

	t.Run("should not move non-exists", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, "non-exists")
		require.Error(t, err)
		require.Contains(t, err.Error(), dymnstypes.ErrOpenPurchaseOrderNotFound.Error())
	})

	t.Run("should able to move a duplicated without error", func(t *testing.T) {
		err = dk.SetOpenPurchaseOrder(ctx, opo11)
		require.NoError(t, err)

		defer func() {
			dk.DeleteOpenPurchaseOrder(ctx, opo11)
		}()

		err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo11.Name)
		require.NoError(t, err)

		list := dk.GetHistoricalOpenPurchaseOrders(ctx, opo11.Name)
		require.Len(t, list, 1, "do not persist duplicated historical OPO")
	})

	t.Run("other records remaining as-is", func(t *testing.T) {
		require.Empty(t, dk.GetOpenPurchaseOrder(ctx, dymName2.Name))
	})

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:     dymName2.Name,
		ExpireAt: 1,
		MinPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo2)
	require.NoError(t, err)

	t.Run("should able to move", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo2.Name)
		require.NoError(t, err)
	})

	t.Run("other records remaining as-is", func(t *testing.T) {
		require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName1.Name), 1)
		require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name), 1)
	})

	opo12 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  epoch,
		MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
		SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(300)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo12)
	require.NoError(t, err)
	t.Run("should not move yet finished OPO", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo12.Name)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has not expired yet")
	})

	opo12.HighestBid = &dymnstypes.OpenPurchaseOrderBid{
		Bidder: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
		Price:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(300)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo12)
	require.NoError(t, err)

	t.Run("should able to move finished OPO", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo12.Name)
		require.NoError(t, err)

		list := dk.GetHistoricalOpenPurchaseOrders(ctx, opo12.Name)
		require.Len(t, list, 2, "should appended to historical")
	})

	t.Run("other records remaining as-is", func(t *testing.T) {
		require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name), 1)
	})
}

func TestKeeper_GetHistoricalOpenPurchaseOrders(t *testing.T) {
	dk, _, ctx := testkeeper.DymNSKeeper(t)

	epoch := time.Now().UTC().Add(time.Hour).Unix()

	// setting block time
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Now().UTC(),
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   epoch,
	}
	err := dk.SetDymName(ctx, dymName1)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "owned-by-1",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   epoch,
	}
	err = dk.SetDymName(ctx, dymName2)
	require.NoError(t, err)

	t.Run("getting non-exists should returns empty", func(t *testing.T) {
		require.Empty(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName1.Name))
		require.Empty(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name))
	})

	opo11 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  1,
		MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
		SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(300)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo11)
	require.NoError(t, err)
	err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo11.Name)
	require.NoError(t, err)

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:     dymName2.Name,
		ExpireAt: 1,
		MinPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(100)),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo2)
	require.NoError(t, err)
	err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo2.Name)
	require.NoError(t, err)

	opo2.ExpireAt++
	err = dk.SetOpenPurchaseOrder(ctx, opo2)
	require.NoError(t, err)
	err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo2.Name)
	require.NoError(t, err)

	t.Run("fetch correctly", func(t *testing.T) {
		list1 := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName1.Name)
		require.Len(t, list1, 1)
		list2 := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name)
		require.Len(t, list2, 2)
		require.Equal(t, opo2.Name, list2[0].Name)
		require.Equal(t, opo2.Name, list2[1].Name)
		require.Equal(t, int64(1), list2[0].ExpireAt)
		require.Equal(t, int64(2), list2[1].ExpireAt)
	})
}
