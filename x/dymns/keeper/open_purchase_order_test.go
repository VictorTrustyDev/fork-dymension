package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
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
