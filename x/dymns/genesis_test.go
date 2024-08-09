package dymns_test

import (
	"testing"
	"time"

	"github.com/dymensionxyz/dymension/v3/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	"github.com/dymensionxyz/dymension/v3/x/dymns"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
func TestExportThenInitGenesis(t *testing.T) {
	now := time.Now().UTC()

	oldKeeper, _, _, oldCtx := testkeeper.DymNSKeeper(t)
	oldCtx = oldCtx.WithBlockTime(now)
	moduleParams := oldKeeper.GetParams(oldCtx)
	moduleParams.Misc.GracePeriodDuration = 30 * 24 * time.Hour
	require.NoError(t, oldKeeper.SetParams(oldCtx, moduleParams))

	// Setup genesis state
	owner1 := sample.AccAddress()
	owner2 := sample.AccAddress()
	anotherAccount := sample.AccAddress()

	bidder1 := sample.AccAddress()
	bidder2 := sample.AccAddress()
	bidder3 := sample.AccAddress()

	buyer1 := sample.AccAddress()
	buyer2 := sample.AccAddress()
	buyer3 := sample.AccAddress()
	buyer4 := sample.AccAddress()
	buyer5 := sample.AccAddress()

	dymName1 := dymnstypes.DymName{
		Name:       "my-name",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   now.Add(time.Hour).Unix(),
		Configs: []dymnstypes.DymNameConfig{
			{
				Type:  dymnstypes.DymNameConfigType_DCT_NAME,
				Path:  "pseudo",
				Value: anotherAccount,
			},
		},
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName1))

	dymName2 := dymnstypes.DymName{
		Name:       "light",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   now.Add(time.Hour).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName2))

	dymName3JustExpired := dymnstypes.DymName{
		Name:       "just-expired",
		Owner:      owner2,
		Controller: owner2,
		ExpireAt:   now.Add(-time.Second).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName3JustExpired))

	dymName4LongExpired := dymnstypes.DymName{
		Name:       "long-expired",
		Owner:      owner1,
		Controller: owner1,
		ExpireAt:   now.Add(-moduleParams.Misc.GracePeriodDuration - time.Second).Unix(),
	}
	require.NoError(t, oldKeeper.SetDymName(oldCtx, dymName4LongExpired))

	so1 := dymnstypes.SellOrder{
		AssetId:   dymName1.Name,
		AssetType: dymnstypes.TypeName,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
		HighestBid: &dymnstypes.SellOrderBid{
			Bidder: bidder1,
			Price:  dymnsutils.TestCoin(200),
		},
	}
	require.NoError(t, oldKeeper.SetSellOrder(oldCtx, so1))

	so2 := dymnstypes.SellOrder{
		AssetId:   dymName2.Name,
		AssetType: dymnstypes.TypeName,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(900),
		HighestBid: &dymnstypes.SellOrderBid{
			Bidder: bidder2,
			Price:  dymnsutils.TestCoin(800),
		},
	}
	require.NoError(t, oldKeeper.SetSellOrder(oldCtx, so2))

	so3 := dymnstypes.SellOrder{
		AssetId:   dymName3JustExpired.Name,
		AssetType: dymnstypes.TypeName,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(200),
	}
	require.NoError(t, oldKeeper.SetSellOrder(oldCtx, so3))

	so4 := dymnstypes.SellOrder{
		AssetId:   "alias",
		AssetType: dymnstypes.TypeAlias,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(900),
		HighestBid: &dymnstypes.SellOrderBid{
			Bidder: bidder3,
			Price:  dymnsutils.TestCoin(777),
			Params: []string{"rollapp_1-1"},
		},
	}
	require.NoError(t, oldKeeper.SetSellOrder(oldCtx, so4))

	so5 := dymnstypes.SellOrder{
		AssetId:   "cosmos",
		AssetType: dymnstypes.TypeAlias,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
	}
	require.NoError(t, oldKeeper.SetSellOrder(oldCtx, so5))

	offer1 := dymnstypes.BuyOrder{
		Id:         "101",
		AssetId:    dymName1.Name,
		AssetType:  dymnstypes.TypeName,
		Buyer:      buyer1,
		OfferPrice: dymnsutils.TestCoin(100),
	}
	require.NoError(t, oldKeeper.SetBuyOrder(oldCtx, offer1))

	offer2 := dymnstypes.BuyOrder{
		Id:         "102",
		AssetId:    dymName2.Name,
		AssetType:  dymnstypes.TypeName,
		Buyer:      buyer2,
		OfferPrice: dymnsutils.TestCoin(200),
	}
	require.NoError(t, oldKeeper.SetBuyOrder(oldCtx, offer2))

	offer3OfExpired := dymnstypes.BuyOrder{
		Id:         "103",
		AssetId:    dymName3JustExpired.Name,
		AssetType:  dymnstypes.TypeName,
		Buyer:      buyer3,
		OfferPrice: dymnsutils.TestCoin(300),
	}
	require.NoError(t, oldKeeper.SetBuyOrder(oldCtx, offer3OfExpired))

	offer4 := dymnstypes.BuyOrder{
		Id:         "204",
		AssetId:    "cosmos",
		AssetType:  dymnstypes.TypeAlias,
		Params:     []string{"rollapp_2-2"},
		Buyer:      buyer4,
		OfferPrice: dymnsutils.TestCoin(333),
	}
	require.NoError(t, oldKeeper.SetBuyOrder(oldCtx, offer4))

	offer5 := dymnstypes.BuyOrder{
		Id:         "205",
		AssetId:    "alias",
		AssetType:  dymnstypes.TypeAlias,
		Params:     []string{"rollapp_3-3"},
		Buyer:      buyer5,
		OfferPrice: dymnsutils.TestCoin(555),
	}
	require.NoError(t, oldKeeper.SetBuyOrder(oldCtx, offer5))

	// Export genesis state
	genState := dymns.ExportGenesis(oldCtx, oldKeeper)

	t.Run("params should be exported correctly", func(t *testing.T) {
		require.Equal(t, oldKeeper.GetParams(oldCtx), genState.Params)
	})

	t.Run("dym-names should be exported correctly", func(t *testing.T) {
		require.Len(t, genState.DymNames, 3)
		require.Contains(t, genState.DymNames, dymName1)
		require.Contains(t, genState.DymNames, dymName2)

		// Expired Dym-Names
		// which less than grace period should be included
		require.Contains(t, genState.DymNames, dymName3JustExpired)
		// which passed grace period should not be included
		require.NotContains(t, genState.DymNames, dymName4LongExpired)
	})

	t.Run("sell orders's non-refunded bids should be exported correctly", func(t *testing.T) {
		require.Len(t, genState.SellOrderBids, 3)
		require.Contains(t, genState.SellOrderBids, *so1.HighestBid)
		require.Contains(t, genState.SellOrderBids, *so2.HighestBid)
		require.Contains(t, genState.SellOrderBids, *so4.HighestBid)
		// Expired sell order should not be exported
	})

	t.Run("buy offers should be exported correctly", func(t *testing.T) {
		require.Len(t, genState.BuyOrders, 5)
		require.Contains(t, genState.BuyOrders, offer1)
		require.Contains(t, genState.BuyOrders, offer2)
		require.Contains(
			t, genState.BuyOrders, offer3OfExpired,
			"offer should be exported even if the dym-name is expired",
		)
		require.Contains(t, genState.BuyOrders, offer4)
		require.Contains(t, genState.BuyOrders, offer5)
	})

	// Init genesis state

	genState.Params.Misc.BeginEpochHookIdentifier = "week" // Change the epoch identifier to test if it is imported correctly
	genState.Params.Misc.SellOrderDuration = 9999 * time.Hour
	genState.Params.PreservedRegistration.ExpirationEpoch = 8888

	newDymNsKeeper, newBankKeeper, _, newCtx := testkeeper.DymNSKeeper(t)
	newCtx = newCtx.WithBlockTime(now)

	dymns.InitGenesis(newCtx, newDymNsKeeper, *genState)

	t.Run("params should be imported correctly", func(t *testing.T) {
		importedParams := newDymNsKeeper.GetParams(newCtx)
		require.Equal(t, genState.Params, importedParams)
		require.Equal(t, "week", importedParams.Misc.BeginEpochHookIdentifier)
		require.Equal(t, 9999*time.Hour, importedParams.Misc.SellOrderDuration)
		require.Equal(t, int64(8888), importedParams.PreservedRegistration.ExpirationEpoch)
	})

	t.Run("Dym-Names should be imported correctly", func(t *testing.T) {
		require.Len(t, newDymNsKeeper.GetAllNonExpiredDymNames(newCtx), 2)
		require.Len(t, newDymNsKeeper.GetAllDymNames(newCtx), 3)

		require.Equal(t, &dymName1, newDymNsKeeper.GetDymName(newCtx, dymName1.Name))
		require.Equal(t, &dymName2, newDymNsKeeper.GetDymName(newCtx, dymName2.Name))
		require.Equal(t, &dymName3JustExpired, newDymNsKeeper.GetDymName(newCtx, dymName3JustExpired.Name))
		require.Nil(t, newDymNsKeeper.GetDymName(newCtx, dymName4LongExpired.Name))
	})

	t.Run("reverse lookup should be created correctly", func(t *testing.T) {
		owned, err := newDymNsKeeper.GetDymNamesOwnedBy(newCtx, owner1)
		require.NoError(t, err)
		require.Len(t, owned, 1)

		owned, err = newDymNsKeeper.GetDymNamesOwnedBy(newCtx, owner2)
		require.NoError(t, err)
		require.Len(t, owned, 1)

		names, err := newDymNsKeeper.GetDymNamesContainsConfiguredAddress(newCtx, owner1)
		require.NoError(t, err)
		require.Len(t, names, 1)

		names, err = newDymNsKeeper.GetDymNamesContainsConfiguredAddress(newCtx, owner2)
		require.NoError(t, err)
		require.Len(t, names, 1)

		names, err = newDymNsKeeper.GetDymNamesContainsConfiguredAddress(newCtx, anotherAccount)
		require.NoError(t, err)
		require.Len(t, names, 1)

		names, err = newDymNsKeeper.GetDymNamesContainsFallbackAddress(newCtx, sdk.MustAccAddressFromBech32(owner1).Bytes())
		require.NoError(t, err)
		require.Len(t, names, 1)

		names, err = newDymNsKeeper.GetDymNamesContainsFallbackAddress(newCtx, sdk.MustAccAddressFromBech32(owner2).Bytes())
		require.NoError(t, err)
		require.Len(t, names, 1)

		names, err = newDymNsKeeper.GetDymNamesContainsFallbackAddress(newCtx, sdk.MustAccAddressFromBech32(anotherAccount).Bytes())
		require.NoError(t, err)
		require.Empty(t, names, 0)
	})

	t.Run("sell orders's non-refunded bids should be refunded correctly", func(t *testing.T) {
		require.Equal(t,
			dymnsutils.TestCoin(200),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(bidder1), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(800),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(bidder2), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(777),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(bidder3), params.BaseDenom),
		)
	})

	t.Run("non-refunded buy-offers should be refunded correctly", func(t *testing.T) {
		require.Equal(t,
			dymnsutils.TestCoin(100),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(buyer1), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(200),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(buyer2), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(300),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(buyer3), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(333),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(buyer4), params.BaseDenom),
		)
		require.Equal(t,
			dymnsutils.TestCoin(555),
			newBankKeeper.GetBalance(newCtx, sdk.MustAccAddressFromBech32(buyer5), params.BaseDenom),
		)
	})

	// Init genesis state but with invalid input
	newDymNsKeeper, newBankKeeper, _, newCtx = testkeeper.DymNSKeeper(t)

	t.Run("fail - invalid params", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.Params{
					Price: dymnstypes.PriceParams{}, // empty
					Misc:  dymnstypes.MiscParams{},  // empty
				},
			})
		})
	})

	t.Run("fail - invalid dym-name", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.DefaultParams(),
				DymNames: []dymnstypes.DymName{
					{}, // empty content
				},
			})
		})
	})

	t.Run("fail - invalid highest bid", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.DefaultParams(),
				SellOrderBids: []dymnstypes.SellOrderBid{
					{}, // empty content
				},
			})
		})
	})

	t.Run("fail - invalid offer", func(t *testing.T) {
		require.Panics(t, func() {
			dymns.InitGenesis(newCtx, newDymNsKeeper, dymnstypes.GenesisState{
				Params: dymnstypes.DefaultParams(),
				BuyOrders: []dymnstypes.BuyOrder{
					{}, // empty content
				},
			})
		})
	})
}
