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
func Test_msgServer_PurchaseName(t *testing.T) {
	now := time.Now().UTC()
	futureEpoch := now.Unix() + 1

	setupTest := func() (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
		dk, bk, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		moduleParams := dk.GetParams(ctx)
		moduleParams.Misc.GasCrudOpenPurchaseOrder = 20_000_000
		err := dk.SetParams(ctx, moduleParams)
		require.NoError(t, err)

		return dk, bk, ctx
	}

	t.Run("reject if message not pass validate basic", func(t *testing.T) {
		dk, _, ctx := setupTest()

		requireErrorFContains(t, func() error {
			_, err := dymnskeeper.NewMsgServerImpl(dk).PurchaseName(ctx, &dymnstypes.MsgPurchaseName{})
			return err
		}, dymnstypes.ErrValidationFailed.Error())
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	buyer := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	previousBidder := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	originalDymNameExpiry := futureEpoch
	dymName := dymnstypes.DymName{
		Name:       "bonded-pool",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   originalDymNameExpiry,
	}

	var ownerOriginalBalance int64 = 1000
	var buyerOriginalBalance int64 = 500
	var previousBidderOriginalBalance int64 = 400
	var minPrice int64 = 100
	tests := []struct {
		name                           string
		withoutDymName                 bool
		withoutOpo                     bool
		expiredOpo                     bool
		sellPrice                      int64
		previousBid                    int64
		skipPreMintModuleAccount       bool
		overrideBuyerOriginalBalance   int64
		customBuyer                    string
		newBid                         int64
		customBidDenom                 string
		wantOwnershipChanged           bool
		wantErr                        bool
		wantErrContains                string
		wantOwnerBalanceLater          int64
		wantBuyerBalanceLater          int64
		wantPreviousBidderBalanceLater int64
	}{
		{
			name:                           "fail - Dym-Name does not exists, OPO does not exists",
			withoutDymName:                 true,
			withoutOpo:                     true,
			newBid:                         100,
			wantErr:                        true,
			wantErrContains:                dymnstypes.ErrDymNameNotFound.Error(),
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - Dym-Name does not exists, OPO exists",
			withoutDymName:                 true,
			withoutOpo:                     false,
			newBid:                         100,
			wantErr:                        true,
			wantErrContains:                dymnstypes.ErrDymNameNotFound.Error(),
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - Dym-Name exists, OPO does not exists",
			withoutDymName:                 false,
			withoutOpo:                     true,
			newBid:                         100,
			wantErr:                        true,
			wantErrContains:                dymnstypes.ErrOpenPurchaseOrderNotFound.Error(),
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - self-purchase is not allowed",
			customBuyer:                    owner,
			newBid:                         100,
			wantErr:                        true,
			wantErrContains:                "cannot purchase your own dym name",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase an expired order, no bid",
			expiredOpo:                     true,
			newBid:                         100,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, expired, with bid, without sell price",
			expiredOpo:                     true,
			sellPrice:                      0,
			previousBid:                    200,
			newBid:                         300,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, expired, with sell price, with bid under sell price",
			expiredOpo:                     true,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         300,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, expired, with sell price, with bid = sell price",
			expiredOpo:                     true,
			sellPrice:                      300,
			previousBid:                    300,
			newBid:                         300,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, not expired, fail because previous bid matches sell price",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    300,
			newBid:                         300,
			wantErr:                        true,
			wantErrContains:                "cannot purchase a completed order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase order, not expired, fail because lower than previous bid",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200 - 1,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase order, not expired, fail because equals to previous bid",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, expired, bid equals to previous bid",
			expiredOpo:                     true,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - purchase a completed order, expired, bid lower than previous bid",
			expiredOpo:                     true,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200 - 1,
			wantErr:                        true,
			wantErrContains:                "cannot purchase an expired order",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - mis-match denom",
			expiredOpo:                     false,
			newBid:                         200,
			customBidDenom:                 "u" + params.BaseDenom,
			wantErr:                        true,
			wantErrContains:                "offer denom does not match the order denom",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer lower than min-price",
			expiredOpo:                     false,
			newBid:                         minPrice - 1,
			wantErr:                        true,
			wantErrContains:                "offer is lower than minimum price",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer higher than sell-price",
			expiredOpo:                     false,
			sellPrice:                      300,
			newBid:                         300 + 1,
			wantErr:                        true,
			wantErrContains:                "offer is higher than sell price",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer equals to previous bid, no sell price",
			expiredOpo:                     false,
			previousBid:                    200,
			newBid:                         200,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer lower than previous bid, no sell price",
			expiredOpo:                     false,
			previousBid:                    200,
			newBid:                         200 - 1,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer equals to previous bid, has sell price",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - offer lower than previous bid, has sell price",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    200,
			newBid:                         200 - 1,
			wantErr:                        true,
			wantErrContains:                "new offer must be higher than current highest bid",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "success - place bid, = min price, no previous bid, no sell price",
			expiredOpo:                     false,
			newBid:                         minPrice,
			wantOwnershipChanged:           false,
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance - minPrice,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "success - place bid, greater than previous bid, no sell price",
			expiredOpo:                     false,
			previousBid:                    minPrice,
			newBid:                         minPrice + 1,
			wantOwnershipChanged:           false,
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance - (minPrice + 1),
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance + minPrice, // refund
		},
		{
			name:                           "fail - failed to refund previous bid",
			expiredOpo:                     false,
			previousBid:                    minPrice,
			skipPreMintModuleAccount:       true,
			newBid:                         minPrice + 1,
			wantOwnershipChanged:           false,
			wantErr:                        true,
			wantErrContains:                "insufficient funds",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "fail - insufficient buyer funds",
			expiredOpo:                     false,
			overrideBuyerOriginalBalance:   1,
			newBid:                         minPrice + 1,
			wantOwnershipChanged:           false,
			wantErr:                        true,
			wantErrContains:                "insufficient funds",
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          1,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance,
		},
		{
			name:                           "success - place bid, greater than previous bid, under sell price",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    minPrice,
			newBid:                         300 - 1,
			wantOwnershipChanged:           false,
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance - (300 - 1),
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance + minPrice, // refund
		},
		{
			name:                           "success - place bid, greater than previous bid, equals sell price, transfer ownership",
			expiredOpo:                     false,
			sellPrice:                      300,
			previousBid:                    minPrice,
			newBid:                         300,
			wantOwnershipChanged:           true,
			wantOwnerBalanceLater:          ownerOriginalBalance + 300,
			wantBuyerBalanceLater:          buyerOriginalBalance - 300,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance + minPrice, // refund
		},
		{
			name:                           "refund previous bidder",
			expiredOpo:                     false,
			previousBid:                    minPrice,
			newBid:                         200,
			wantOwnershipChanged:           false,
			wantOwnerBalanceLater:          ownerOriginalBalance,
			wantBuyerBalanceLater:          buyerOriginalBalance - 200,
			wantPreviousBidderBalanceLater: previousBidderOriginalBalance + minPrice, // refund
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup execution context
			dk, bk, ctx := setupTest()

			useOwnerOriginalBalance := ownerOriginalBalance
			useBuyerOriginalBalance := buyerOriginalBalance
			if tt.overrideBuyerOriginalBalance > 0 {
				useBuyerOriginalBalance = tt.overrideBuyerOriginalBalance
			}
			usePreviousBidderOriginalBalance := previousBidderOriginalBalance

			err := bk.MintCoins(ctx,
				dymnstypes.ModuleName,
				dymnsutils.TestCoins(
					useOwnerOriginalBalance+useBuyerOriginalBalance+usePreviousBidderOriginalBalance,
				),
			)
			require.NoError(t, err)
			err = bk.SendCoinsFromModuleToAccount(ctx,
				dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(owner),
				dymnsutils.TestCoins(useOwnerOriginalBalance),
			)
			require.NoError(t, err)
			err = bk.SendCoinsFromModuleToAccount(ctx,
				dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(buyer),
				dymnsutils.TestCoins(useBuyerOriginalBalance),
			)
			require.NoError(t, err)
			err = bk.SendCoinsFromModuleToAccount(ctx,
				dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(previousBidder),
				dymnsutils.TestCoins(usePreviousBidderOriginalBalance),
			)
			require.NoError(t, err)

			dymName.Configs = []dymnstypes.DymNameConfig{
				{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: owner,
				},
			}

			if !tt.withoutDymName {
				err = dk.SetDymName(ctx, dymName)
				require.NoError(t, err)
			}

			opo := dymnstypes.OpenPurchaseOrder{
				Name:     dymName.Name,
				MinPrice: dymnsutils.TestCoin(minPrice),
			}

			if tt.expiredOpo {
				opo.ExpireAt = now.Unix() - 1
			} else {
				opo.ExpireAt = futureEpoch
			}

			require.GreaterOrEqual(t, tt.sellPrice, int64(0), "bad setup")
			if tt.sellPrice > 0 {
				opo.SellPrice = dymnsutils.TestCoinP(tt.sellPrice)
			}

			require.GreaterOrEqual(t, tt.previousBid, int64(0), "bad setup")
			if tt.previousBid > 0 {
				opo.HighestBid = &dymnstypes.OpenPurchaseOrderBid{
					Bidder: previousBidder,
					Price:  dymnsutils.TestCoin(tt.previousBid),
				}

				// mint coin to module account because we charged bidder before update OPO
				if !tt.skipPreMintModuleAccount {
					err = bk.MintCoins(ctx, dymnstypes.ModuleName, sdk.NewCoins(opo.HighestBid.Price))
					require.NoError(t, err)
				}
			}

			if !tt.withoutOpo {
				err = dk.SetOpenPurchaseOrder(ctx, opo)
				require.NoError(t, err)
			}

			// test

			require.Greater(t, tt.newBid, int64(0), "mis-configured test case")
			useBuyer := buyer
			if tt.customBuyer != "" {
				useBuyer = tt.customBuyer
			}
			useDenom := params.BaseDenom
			if tt.customBidDenom != "" {
				useDenom = tt.customBidDenom
			}
			resp, errPurchaseName := dymnskeeper.NewMsgServerImpl(dk).PurchaseName(ctx, &dymnstypes.MsgPurchaseName{
				Name:  dymName.Name,
				Offer: sdk.NewInt64Coin(useDenom, tt.newBid),
				Buyer: useBuyer,
			})
			laterDymName := dk.GetDymName(ctx, dymName.Name)
			if !tt.withoutDymName {
				require.NotNil(t, laterDymName)
				require.Equal(t, dymName.Name, laterDymName.Name, "name should not be changed")
				require.Equal(t, originalDymNameExpiry, laterDymName.ExpireAt, "expiry should not be changed")
			}

			laterOpo := dk.GetOpenPurchaseOrder(ctx, dymName.Name)
			historicalOpo := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName.Name)
			laterOwnerBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(owner), params.BaseDenom)
			laterBuyerBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(buyer), params.BaseDenom)
			laterPreviousBidderBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(previousBidder), params.BaseDenom)
			laterDymNamesOwnedByOwner, err := dk.GetDymNamesOwnedBy(ctx, owner, now.Unix())
			require.NoError(t, err)
			laterDymNamesOwnedByBuyer, err := dk.GetDymNamesOwnedBy(ctx, buyer, now.Unix())
			require.NoError(t, err)
			laterDymNamesOwnedByPreviousBidder, err := dk.GetDymNamesOwnedBy(ctx, previousBidder, now.Unix())
			require.NoError(t, err)

			require.Equal(t, tt.wantOwnerBalanceLater, laterOwnerBalance.Amount.Int64(), "owner balance mis-match")
			require.Equal(t, tt.wantBuyerBalanceLater, laterBuyerBalance.Amount.Int64(), "buyer balance mis-match")
			require.Equal(t, tt.wantPreviousBidderBalanceLater, laterPreviousBidderBalance.Amount.Int64(), "previous bidder balance mis-match")

			require.Empty(t, laterDymNamesOwnedByPreviousBidder, "no reverse record should be made for previous bidder")

			moduleParams := dk.GetParams(ctx)
			if tt.wantErr {
				require.Error(t, errPurchaseName, "action should be failed")
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Contains(t, errPurchaseName.Error(), tt.wantErrContains)
				require.Nil(t, resp)

				require.False(t, tt.wantOwnershipChanged, "mis-configured test case")

				require.Less(t,
					ctx.GasMeter().GasConsumed(), sdk.Gas(moduleParams.Misc.GasCrudOpenPurchaseOrder),
					"should not consume params gas on failed operation",
				)
			} else {
				require.NoError(t, errPurchaseName, "action should be successful")
				require.NotNil(t, resp)

				require.GreaterOrEqual(t,
					ctx.GasMeter().GasConsumed(), sdk.Gas(moduleParams.Misc.GasCrudOpenPurchaseOrder),
					"should consume params gas",
				)
			}

			if tt.wantOwnershipChanged {
				if tt.withoutDymName {
					t.Errorf("mis-configured test case")
					return
				}
				if tt.withoutOpo {
					t.Errorf("mis-configured test case")
					return
				}

				require.Nil(t, laterOpo, "OPO should be deleted")
				require.Len(t, historicalOpo, 1, "OPO should be moved to historical")

				require.Equal(t, buyer, laterDymName.Owner, "ownership should be changed")
				require.Equal(t, buyer, laterDymName.Controller, "controller should be changed")
				require.Empty(t, laterDymName.Configs, "configs should be cleared")
				require.Empty(t, laterDymNamesOwnedByOwner, "reverse record should be removed")
				require.Len(t, laterDymNamesOwnedByBuyer, 1, "reverse record should be added")
			} else {
				if tt.withoutDymName {
					require.Nil(t, laterDymName)
					require.Empty(t, laterDymNamesOwnedByOwner)
					require.Empty(t, laterDymNamesOwnedByBuyer)
				} else {
					require.Equal(t, owner, laterDymName.Owner, "ownership should not be changed")
					require.Equal(t, owner, laterDymName.Controller, "controller should not be changed")
					require.NotEmpty(t, laterDymName.Configs, "configs should be kept")
					require.Equal(t, dymName.Configs, laterDymName.Configs, "configs not be changed")
					require.Len(t, laterDymNamesOwnedByOwner, 1, "reverse record should be kept")
					require.Empty(t, laterDymNamesOwnedByBuyer, "reverse record should not be added")
				}

				if tt.withoutOpo {
					require.Nil(t, laterOpo)
					require.Empty(t, historicalOpo)
				} else {
					require.NotNil(t, laterOpo, "OPO should not be deleted")
					require.Empty(t, historicalOpo, "OPO should not be moved to historical")
				}
			}
		})
	}
}
