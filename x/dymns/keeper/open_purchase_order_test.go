package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func TestKeeper_GetSetDeleteOpenPurchaseOrder(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

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
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
		HighestBid: &dymnstypes.OpenPurchaseOrderBid{
			Bidder: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
			Price:  dymnsutils.TestCoin(200),
		},
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo1)
	require.NoError(t, err)
	t.Run("opo1 should be equals to original", func(t *testing.T) {
		require.Equal(t, opo1, *dk.GetOpenPurchaseOrder(ctx, opo1.Name))
	})
	t.Run("opo list should have length 1", func(t *testing.T) {
		require.Len(t, dk.GetAllOpenPurchaseOrders(ctx), 1)
	})
	t.Run("event should be fired on set purchase order", func(t *testing.T) {
		events := ctx.EventManager().Events()
		require.NotEmpty(t, events)

		for _, event := range events {
			if event.Type != dymnstypes.EventTypeDymNameOpenPurchaseOrder {
				continue
			}

			var actionName string
			for _, attr := range event.Attributes {
				if attr.Key == dymnstypes.AttributeKeyDymNameOpoActionName {
					actionName = attr.Value
				}
			}
			require.NotEmpty(t, actionName, "event attr action name could not be found")
			require.Equalf(t,
				actionName, dymnstypes.AttributeKeyDymNameOpoActionNameSet,
				"event attr action name should be `%s`", dymnstypes.AttributeKeyDymNameOpoActionNameSet,
			)
			return
		}

		t.Errorf("event %s not found", dymnstypes.EventTypeDymNameOpenPurchaseOrder)
	})

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:     dymName2.Name,
		ExpireAt: 1,
		MinPrice: dymnsutils.TestCoin(100),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo2)
	require.NoError(t, err)
	t.Run("opo2 should be equals to original", func(t *testing.T) {
		require.Equal(t, opo2, *dk.GetOpenPurchaseOrder(ctx, opo2.Name))
	})
	t.Run("opo list should have length 2", func(t *testing.T) {
		require.Len(t, dk.GetAllOpenPurchaseOrders(ctx), 2)
	})

	dk.DeleteOpenPurchaseOrder(ctx, opo1.Name)
	t.Run("event should be fired on delete purchase order", func(t *testing.T) {
		events := ctx.EventManager().Events()
		require.NotEmpty(t, events)

		for _, event := range events {
			if event.Type != dymnstypes.EventTypeDymNameOpenPurchaseOrder {
				continue
			}

			var actionName string
			for _, attr := range event.Attributes {
				if attr.Key == dymnstypes.AttributeKeyDymNameOpoActionName {
					actionName = attr.Value
				}
			}
			require.NotEmpty(t, actionName, "event attr action name could not be found")
			require.Equalf(t,
				actionName, dymnstypes.AttributeKeyDymNameOpoActionNameSet,
				"event attr action name should be `%s`", dymnstypes.AttributeKeyDymNameOpoActionNameDelete,
			)
			return
		}

		t.Errorf("event %s not found", dymnstypes.EventTypeDymNameOpenPurchaseOrder)
	})

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

	t.Run("omit Sell Price if not nil but zero", func(t *testing.T) {
		opo3 := dymnstypes.OpenPurchaseOrder{
			Name:      "hello",
			ExpireAt:  1,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(0),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo3)
		require.NoError(t, err)

		require.Nil(t, dk.GetOpenPurchaseOrder(ctx, opo3.Name).SellPrice)
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_MoveOpenPurchaseOrderToHistorical(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	// setting block time
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: time.Now().UTC(),
	})

	futureEpoch := ctx.BlockTime().Add(time.Hour).Unix()

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   futureEpoch,
	}
	err := dk.SetDymName(ctx, dymName1)
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

	opo11 := dymnstypes.OpenPurchaseOrder{
		Name:      dymName1.Name,
		ExpireAt:  1,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
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

	t.Run("has min expiry mapping", func(t *testing.T) {
		minExpiry, found := dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, opo11.Name)
		require.True(t, found)
		require.Equal(t, opo11.ExpireAt, minExpiry)
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
			dk.DeleteOpenPurchaseOrder(ctx, opo11.Name)
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
		MinPrice: dymnsutils.TestCoin(100),
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
		ExpireAt:  futureEpoch,
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
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
		Price:  dymnsutils.TestCoin(300),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo12)
	require.NoError(t, err)

	t.Run("should able to move finished OPO", func(t *testing.T) {
		err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo12.Name)
		require.NoError(t, err)

		list := dk.GetHistoricalOpenPurchaseOrders(ctx, opo12.Name)
		require.Len(t, list, 2, "should appended to historical")

		minExpiry, found := dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, opo12.Name)
		require.True(t, found)
		require.Equal(t, opo11.ExpireAt, minExpiry, "should keep the minimum")
		require.NotEqual(t, opo12.ExpireAt, minExpiry, "should keep the minimum")
	})

	t.Run("other records remaining as-is", func(t *testing.T) {
		require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name), 1)
	})
}

func TestKeeper_GetAndDeleteHistoricalOpenPurchaseOrders(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	now := time.Now().UTC()
	futureEpoch := now.Unix() + 1

	// setting block time
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: now,
	})

	//goland:noinspection SpellCheckingInspection
	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymName1 := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   futureEpoch,
	}
	err := dk.SetDymName(ctx, dymName1)
	require.NoError(t, err)

	dymName2 := dymnstypes.DymName{
		Name:       "owned-by-1",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   futureEpoch,
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
		MinPrice:  dymnsutils.TestCoin(100),
		SellPrice: dymnsutils.TestCoinP(300),
	}
	err = dk.SetOpenPurchaseOrder(ctx, opo11)
	require.NoError(t, err)
	err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo11.Name)
	require.NoError(t, err)

	opo2 := dymnstypes.OpenPurchaseOrder{
		Name:     dymName2.Name,
		ExpireAt: 1,
		MinPrice: dymnsutils.TestCoin(100),
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

	t.Run("delete", func(t *testing.T) {
		dk.DeleteHistoricalOpenPurchaseOrders(ctx, dymName1.Name)
		require.Empty(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName1.Name))

		list2 := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name)
		require.Len(t, list2, 2)

		dk.DeleteHistoricalOpenPurchaseOrders(ctx, dymName2.Name)
		require.Empty(t, dk.GetHistoricalOpenPurchaseOrders(ctx, dymName2.Name))
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeeper_CompletePurchaseOrder(t *testing.T) {
	now := time.Now().UTC()
	futureEpoch := now.Unix() + 1

	setupTest := func() (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
		dk, bk, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		return dk, bk, ctx
	}

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	buyer := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"

	originalDymNameExpiry := futureEpoch
	dymName := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   originalDymNameExpiry,
	}

	t.Run("Dym-Name not found", func(t *testing.T) {
		dk, _, ctx := setupTest()

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, "non-exists"), dymnstypes.ErrDymNameNotFound.Error())
	})

	t.Run("OPO not found", func(t *testing.T) {
		dk, _, ctx := setupTest()

		err := dk.SetDymName(ctx, dymName)
		require.NoError(t, err)

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, dymName.Name), dymnstypes.ErrOpenPurchaseOrderNotFound.Error())
	})

	t.Run("OPO not yet completed, no bidder", func(t *testing.T) {
		dk, _, ctx := setupTest()

		err := dk.SetDymName(ctx, dymName)
		require.NoError(t, err)

		opo := dymnstypes.OpenPurchaseOrder{
			Name:     dymName.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo)
		require.NoError(t, err)

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, dymName.Name), "Open-Purchase-Order has not finished yet")
	})

	t.Run("OPO has bidder but not yet completed", func(t *testing.T) {
		dk, _, ctx := setupTest()

		err := dk.SetDymName(ctx, dymName)
		require.NoError(t, err)

		opo := dymnstypes.OpenPurchaseOrder{
			Name:      dymName.Name,
			ExpireAt:  futureEpoch,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
			HighestBid: &dymnstypes.OpenPurchaseOrderBid{
				Bidder: buyer,
				Price:  dymnsutils.TestCoin(200), // lower than sell price
			},
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo)
		require.NoError(t, err)

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, dymName.Name), "Open-Purchase-Order has not finished yet")
	})

	t.Run("OPO expired without bidder", func(t *testing.T) {
		dk, _, ctx := setupTest()

		err := dk.SetDymName(ctx, dymName)
		require.NoError(t, err)

		opo := dymnstypes.OpenPurchaseOrder{
			Name:      dymName.Name,
			ExpireAt:  now.Unix() - 1,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo)
		require.NoError(t, err)

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, dymName.Name), "no bid placed")
	})

	t.Run("OPO without sell price, with bid, finished by expiry", func(t *testing.T) {
		dk, _, ctx := setupTest()

		err := dk.SetDymName(ctx, dymName)
		require.NoError(t, err)

		opo := dymnstypes.OpenPurchaseOrder{
			Name:     dymName.Name,
			ExpireAt: futureEpoch,
			MinPrice: dymnsutils.TestCoin(100),
			HighestBid: &dymnstypes.OpenPurchaseOrderBid{
				Bidder: buyer,
				Price:  dymnsutils.TestCoin(200),
			},
		}
		err = dk.SetOpenPurchaseOrder(ctx, opo)
		require.NoError(t, err)

		requireErrorContains(t, dk.CompletePurchaseOrder(ctx, dymName.Name), "Open-Purchase-Order has not finished yet")
	})

	var ownerOriginalBalance int64 = 1000
	var buyerOriginalBalance int64 = 500
	tests := []struct {
		name                  string
		expiredOpo            bool
		sellPrice             int64
		bid                   int64
		wantErr               bool
		wantErrContains       string
		wantOwnerBalanceLater int64
	}{
		{
			name:                  "completed, expired, no sell price",
			expiredOpo:            true,
			sellPrice:             0,
			bid:                   200,
			wantErr:               false,
			wantOwnerBalanceLater: ownerOriginalBalance + 200,
		},
		{
			name:                  "completed, expired, under sell price",
			expiredOpo:            true,
			sellPrice:             300,
			bid:                   200,
			wantErr:               false,
			wantOwnerBalanceLater: ownerOriginalBalance + 200,
		},
		{
			name:                  "completed, expired, equals sell price",
			expiredOpo:            true,
			sellPrice:             300,
			bid:                   300,
			wantErr:               false,
			wantOwnerBalanceLater: ownerOriginalBalance + 300,
		},
		{
			name:                  "completed by sell-price met, not expired",
			expiredOpo:            false,
			sellPrice:             300,
			bid:                   300,
			wantErr:               false,
			wantOwnerBalanceLater: ownerOriginalBalance + 300,
		},
		{
			name:                  "expired without bid, no sell price",
			expiredOpo:            true,
			sellPrice:             0,
			bid:                   0,
			wantErr:               true,
			wantErrContains:       "no bid placed",
			wantOwnerBalanceLater: ownerOriginalBalance,
		},
		{
			name:                  "expired without bid, with sell price",
			expiredOpo:            true,
			sellPrice:             300,
			bid:                   0,
			wantErr:               true,
			wantErrContains:       "no bid placed",
			wantOwnerBalanceLater: ownerOriginalBalance,
		},
		{
			name:                  "not expired but bid under sell price",
			expiredOpo:            false,
			sellPrice:             300,
			bid:                   200,
			wantErr:               true,
			wantErrContains:       "Open-Purchase-Order has not finished yet",
			wantOwnerBalanceLater: ownerOriginalBalance,
		},
		{
			name:                  "not expired has bid, no sell price",
			expiredOpo:            false,
			sellPrice:             0,
			bid:                   200,
			wantErr:               true,
			wantErrContains:       "Open-Purchase-Order has not finished yet",
			wantOwnerBalanceLater: ownerOriginalBalance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup execution context
			dk, bk, ctx := setupTest()

			err := bk.MintCoins(ctx,
				dymnstypes.ModuleName,
				dymnsutils.TestCoins(ownerOriginalBalance+buyerOriginalBalance),
			)
			require.NoError(t, err)
			err = bk.SendCoinsFromModuleToAccount(ctx,
				dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(owner),
				dymnsutils.TestCoins(ownerOriginalBalance),
			)
			require.NoError(t, err)
			err = bk.SendCoinsFromModuleToAccount(ctx,
				dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(buyer),
				dymnsutils.TestCoins(buyerOriginalBalance),
			)
			require.NoError(t, err)

			dymName.Configs = []dymnstypes.DymNameConfig{
				{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: owner,
				},
			}
			err = dk.SetDymName(ctx, dymName)
			require.NoError(t, err)

			opo := dymnstypes.OpenPurchaseOrder{
				Name:     dymName.Name,
				MinPrice: dymnsutils.TestCoin(100),
			}

			if tt.expiredOpo {
				opo.ExpireAt = now.Unix() - 1
			} else {
				opo.ExpireAt = futureEpoch
			}

			require.GreaterOrEqual(t, tt.sellPrice, int64(0), "bad setup")
			opo.SellPrice = dymnsutils.TestCoinP(tt.sellPrice)

			require.GreaterOrEqual(t, tt.bid, int64(0), "bad setup")
			if tt.bid > 0 {
				opo.HighestBid = &dymnstypes.OpenPurchaseOrderBid{
					Bidder: buyer,
					Price:  dymnsutils.TestCoin(tt.bid),
				}

				// mint coin to module account because we charged buyer before update OPO
				err = bk.MintCoins(ctx, dymnstypes.ModuleName, sdk.NewCoins(opo.HighestBid.Price))
				require.NoError(t, err)
			}
			err = dk.SetOpenPurchaseOrder(ctx, opo)
			require.NoError(t, err)

			// test

			errCompletePurchaseOrder := dk.CompletePurchaseOrder(ctx, dymName.Name)
			laterDymName := dk.GetDymName(ctx, dymName.Name)
			require.NotNil(t, laterDymName)
			laterOpo := dk.GetOpenPurchaseOrder(ctx, dymName.Name)
			historicalOpo := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName.Name)
			laterOwnerBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(owner), params.BaseDenom)
			laterBuyerBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(buyer), params.BaseDenom)
			laterDymNamesOwnedByOwner, err := dk.GetDymNamesOwnedBy(ctx, owner, now.Unix())
			require.NoError(t, err)
			laterDymNamesOwnedByBuyer, err := dk.GetDymNamesOwnedBy(ctx, buyer, now.Unix())
			require.NoError(t, err)

			require.Equal(t, dymName.Name, laterDymName.Name, "name should not be changed")
			require.Equal(t, originalDymNameExpiry, laterDymName.ExpireAt, "expiry should not be changed")

			if tt.wantErr {
				require.Error(t, errCompletePurchaseOrder, "action should be failed")
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Contains(t, errCompletePurchaseOrder.Error(), tt.wantErrContains)

				require.NotNil(t, laterOpo, "OPO should not be deleted")
				require.Empty(t, historicalOpo, "OPO should not be moved to historical")

				require.Equal(t, owner, laterDymName.Owner, "ownership should not be changed")
				require.Equal(t, owner, laterDymName.Controller, "controller should not be changed")
				require.NotEmpty(t, laterDymName.Configs, "configs should be kept")
				require.Equal(t, dymName.Configs, laterDymName.Configs, "configs not be changed")
				require.Len(t, laterDymNamesOwnedByOwner, 1, "reverse record should be kept")
				require.Empty(t, laterDymNamesOwnedByBuyer, "reverse record should not be added")

				require.Equal(t, ownerOriginalBalance, laterOwnerBalance.Amount.Int64(), "owner balance should not be changed")
				require.Equal(t, tt.wantOwnerBalanceLater, laterOwnerBalance.Amount.Int64(), "owner balance mis-match")
				require.Equal(t, buyerOriginalBalance, laterBuyerBalance.Amount.Int64(), "buyer balance should not be changed")
				return
			}

			require.NoError(t, errCompletePurchaseOrder, "action should be successful")

			require.Nil(t, laterOpo, "OPO should be deleted")
			require.Len(t, historicalOpo, 1, "OPO should be moved to historical")

			require.Equal(t, buyer, laterDymName.Owner, "ownership should be changed")
			require.Equal(t, buyer, laterDymName.Controller, "controller should be changed")
			require.Empty(t, laterDymName.Configs, "configs should be cleared")
			require.Empty(t, laterDymNamesOwnedByOwner, "reverse record should be removed")
			require.Len(t, laterDymNamesOwnedByBuyer, 1, "reverse record should be added")

			require.Equal(t, tt.wantOwnerBalanceLater, laterOwnerBalance.Amount.Int64(), "owner balance mis-match")
			require.Equal(t, buyerOriginalBalance, laterBuyerBalance.Amount.Int64(), "buyer balance should not be changed")
		})
	}
}

func TestKeeper_GetSetActiveOpenPurchaseOrdersExpiration(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	t.Run("get", func(t *testing.T) {
		aope := dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
		require.Empty(t, aope.ExpiryByName, "default map must be empty")
		require.NotNil(t, aope.ExpiryByName, "map must be initialized")
	})

	t.Run("set", func(t *testing.T) {
		aope := dymnstypes.ActiveOpenPurchaseOrdersExpiration{
			ExpiryByName: map[string]int64{
				"hello": 123,
				"world": 456,
			},
		}
		err := dk.SetActiveOpenPurchaseOrdersExpiration(ctx, aope)
		require.NoError(t, err)

		aope = dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
		require.Len(t, aope.ExpiryByName, 2)
		require.Equal(t, int64(123), aope.ExpiryByName["hello"])
		require.Equal(t, int64(456), aope.ExpiryByName["world"])
	})
}

func TestKeeper_GetSetMinExpiryHistoricalOpenPurchaseOrder(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "hello", 123)
	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "world", 456)

	min, found := dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, "hello")
	require.True(t, found)
	require.Equal(t, int64(123), min)

	min, found = dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, "world")
	require.True(t, found)
	require.Equal(t, int64(456), min)

	min, found = dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, "non-exists")
	require.False(t, found)
	require.Zero(t, min)

	t.Run("set zero means delete", func(t *testing.T) {
		dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "hello", 0)

		min, found = dk.GetMinExpiryHistoricalOpenPurchaseOrder(ctx, "hello")
		require.False(t, found)
		require.Zero(t, min)
	})
}

func TestKeeper_GetMinExpiryOfAllHistoricalOpenPurchaseOrders(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "one", 1)
	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "two", 22)
	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "three", 333)

	mapped := dk.GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx)
	require.Len(t, mapped, 3)
	require.Equal(t, int64(1), mapped["one"])
	require.Equal(t, int64(22), mapped["two"])
	require.Equal(t, int64(333), mapped["three"])

	dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, "three", 0)
	mapped = dk.GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx)
	require.Len(t, mapped, 2)
	require.Equal(t, int64(1), mapped["one"])
	require.Equal(t, int64(22), mapped["two"])
}
