package dymns_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	"github.com/dymensionxyz/dymension/v3/x/dymns"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestExportThenInitGenesis(t *testing.T) {
	oldKeeper, _, _, oldCtx := testkeeper.DymNSKeeper(t)

	// Setup genesis state
	owner1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	owner2 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	bidder1 := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	bidder2 := "dym1wl60kvsq5c4wa600h7rnez8dguk5lpnqp4u0y2"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName1))

	dymName2 := dymnstypes.DymName{
		Name:       "not-bonded-pool",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(time.Hour).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName2))

	dymName3Expired := dymnstypes.DymName{
		Name:       "not-bonded-pool2",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   time.Now().UTC().Add(-time.Hour).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName3Expired))

	opo1 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
		HighestBid: &dymnstypes.OpenPurchaseOrderBid{
			Bidder: bidder1,
			Price:  dymnsutils.TestCoin(200),
		},
	}
	require.NoError(t, oldKeeper.SetOpenPurchaseOrder(oldCtx, opo1))

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName2.Name,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(900),
		HighestBid: &dymnstypes.OpenPurchaseOrderBid{
			Bidder: bidder2,
			Price:  dymnsutils.TestCoin(800),
		},
	}
	require.NoError(t, oldKeeper.SetOpenPurchaseOrder(oldCtx, opo2))

	opo3 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName3Expired.Name,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(200),
	}
	require.NoError(t, oldKeeper.SetOpenPurchaseOrder(oldCtx, opo3))

	// Export genesis state
	genState := dymns.ExportGenesis(oldCtx, oldKeeper)

	t.Run("params should be exported correctly", func(t *testing.T) {
		require.Equal(t, oldKeeper.GetParams(oldCtx), genState.Params)
	})

	t.Run("dym-names should be exported correctly", func(t *testing.T) {
		require.Len(t, genState.DymNames, 2)
		require.Contains(t, genState.DymNames, dymName1)
		require.Contains(t, genState.DymNames, dymName2)
		// Expired dym-name should not be exported
	})

	t.Run("open purchase orders's non-refunded bids should be exported correctly", func(t *testing.T) {
		require.Len(t, genState.OpenPurchaseOrderBids, 2)
		require.Contains(t, genState.OpenPurchaseOrderBids, *opo1.HighestBid)
		require.Contains(t, genState.OpenPurchaseOrderBids, *opo2.HighestBid)
		// Expired open purchase order should not be exported
	})

	// Init genesis state

	genState.Params.Misc.BeginEpochHookIdentifier = "week" // Change the epoch identifier to test if it is imported correctly
	genState.Params.Misc.DaysOpenPurchaseOrderDuration = 9999

	newDymNsKeeper, newBankKeeper, _, newCtx := testkeeper.DymNSKeeper(t)
	dymns.InitGenesis(newCtx, newDymNsKeeper, *genState)

	t.Run("params should be imported correctly", func(t *testing.T) {
		importedParams := newDymNsKeeper.GetParams(newCtx)
		require.Equal(t, genState.Params, importedParams)
		require.Equal(t, "week", importedParams.Misc.BeginEpochHookIdentifier)
		require.Equal(t, int32(9999), importedParams.Misc.DaysOpenPurchaseOrderDuration)
	})

	t.Run("dym-names should be imported correctly", func(t *testing.T) {
		// Expired dym-name should not be imported
		require.Len(t, newDymNsKeeper.GetAllNonExpiredDymNames(newCtx, time.Now().Unix()), 2)
		require.Equal(t, &dymName1, newDymNsKeeper.GetDymName(newCtx, dymName1.Name))
		require.Equal(t, &dymName2, newDymNsKeeper.GetDymName(newCtx, dymName2.Name))
	})

	t.Run("open purchase orders's non-refunded bids should be refunded correctly", func(t *testing.T) {
		require.Equal(t,
			dymnsutils.TestCoin(200),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(bidder1), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(800),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(bidder2), params.BaseDenom),
		)
	})

	// Init genesis state but with invalid input
	newDymNsKeeper, newBankKeeper, _, newCtx = testkeeper.DymNSKeeper(t)

	t.Run("invalid params", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.Params{
					Price: dymnstypes.PriceParams{}, // empty
					Misc:  dymnstypes.MiscParams{},  // empty
				},
			})
		})
	})

	t.Run("invalid dym-name", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.DefaultParams(),
				DymNames: []dymnstypes.DymName{
					{}, // empty content
				},
			})
		})
	})

	t.Run("invalid highest bid", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.DefaultParams(),
				OpenPurchaseOrderBids: []dymnstypes.OpenPurchaseOrderBid{
					{}, // empty content
				},
			})
		})
	})
}
