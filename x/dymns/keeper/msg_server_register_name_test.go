package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func Test_msgServer_RegisterName(t *testing.T) {
	now := time.Now().UTC()

	var denom = dymnsutils.TestCoin(0).Denom
	const firstYearPrice1L = 6
	const firstYearPrice2L = 5
	const firstYearPrice3L = 4
	const firstYearPrice4L = 3
	const firstYearPrice5PlusL = 2
	const extendsPrice = 1
	const gracePeriod = 30

	setupTest := func() (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
		dk, bk, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		moduleParams := dk.GetParams(ctx)
		moduleParams.Price.Price_1Letter = sdk.NewInt(firstYearPrice1L)
		moduleParams.Price.Price_2Letters = sdk.NewInt(firstYearPrice2L)
		moduleParams.Price.Price_3Letters = sdk.NewInt(firstYearPrice3L)
		moduleParams.Price.Price_4Letters = sdk.NewInt(firstYearPrice4L)
		moduleParams.Price.Price_5PlusLetters = sdk.NewInt(firstYearPrice5PlusL)
		moduleParams.Price.PriceExtends = sdk.NewInt(extendsPrice)
		moduleParams.Price.PriceDenom = denom
		moduleParams.Misc.DaysGracePeriod = gracePeriod
		err := dk.SetParams(ctx, moduleParams)
		require.NoError(t, err)

		return dk, bk, ctx
	}

	t.Run("reject if message not pass validate basic", func(t *testing.T) {
		dk, _, ctx := setupTest()

		requireErrorFContains(t, func() error {
			_, err := dymnskeeper.NewMsgServerImpl(dk).RegisterName(ctx, &dymnstypes.MsgRegisterName{})
			return err
		}, dymnstypes.ErrValidationFailed.Error())
	})

	buyer := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	previousOwner := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	var originalModuleBalance int64 = 88

	tests := []struct {
		name                    string
		buyer                   string
		originalBalance         int64
		duration                int32
		customDymName           string
		existingDymName         *dymnstypes.DymName
		setupHistoricalData     bool
		wantLaterDymName        *dymnstypes.DymName
		wantErr                 bool
		wantErrContains         string
		wantLaterBalance        int64
		wantPruneHistoricalData bool
	}{
		{
			name:            "not allow to takeover a non-expired Dym-Name",
			buyer:           buyer,
			originalBalance: 1,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Add(time.Hour).Unix(),
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Add(time.Hour).Unix(),
			},
			wantErr:          true,
			wantErrContains:  sdkerrors.ErrUnauthorized.Error(),
			wantLaterBalance: 1,
		},
		{
			name:            "not allow to takeover an expired Dym-Name which in grace period",
			buyer:           buyer,
			originalBalance: 1,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			wantErr:          true,
			wantErrContains:  "can be taken over after",
			wantLaterBalance: 1,
		},
		{
			name:             "not enough balance to pay for the Dym-Name",
			buyer:            buyer,
			originalBalance:  1,
			duration:         2,
			wantErr:          true,
			wantErrContains:  "insufficient funds",
			wantLaterBalance: 1,
		},
		{
			name:            "can register, new Dym-Name",
			buyer:           buyer,
			originalBalance: firstYearPrice5PlusL + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 5+ letters, multiple years",
			buyer:           buyer,
			originalBalance: firstYearPrice5PlusL + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 5+ letters, 1 year",
			buyer:           buyer,
			originalBalance: firstYearPrice5PlusL + 3,
			duration:        1,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 4 letters, multiple years",
			buyer:           buyer,
			customDymName:   "abcd",
			originalBalance: firstYearPrice4L + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 4 letters, 1 year",
			buyer:           buyer,
			customDymName:   "abcd",
			originalBalance: firstYearPrice4L + 3,
			duration:        1,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 3 letters, multiple years",
			buyer:           buyer,
			customDymName:   "abc",
			originalBalance: firstYearPrice3L + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 3 letters, 1 year",
			buyer:           buyer,
			customDymName:   "abc",
			originalBalance: firstYearPrice3L + 3,
			duration:        1,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 2 letters, multiple years",
			buyer:           buyer,
			customDymName:   "ab",
			originalBalance: firstYearPrice2L + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 2 letters, 1 year",
			buyer:           buyer,
			customDymName:   "ab",
			originalBalance: firstYearPrice2L + 3,
			duration:        1,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 1 letter, multiple years",
			buyer:           buyer,
			customDymName:   "a",
			originalBalance: firstYearPrice1L + extendsPrice + 3,
			duration:        2,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "deduct balance for new Dym-Name, 1 letter, 1 year",
			buyer:           buyer,
			customDymName:   "a",
			originalBalance: firstYearPrice1L + 3,
			duration:        1,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "can extend owned Dym-Name, not expired",
			buyer:           buyer,
			originalBalance: extendsPrice*2 + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1 + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "when extend owned non-expired Dym-Name, keep config and historical data",
			buyer:           buyer,
			originalBalance: extendsPrice*2 + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: buyer,
				}},
			},
			setupHistoricalData: true,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1 + 86400*365*2,
				Configs: []dymnstypes.DymNameConfig{{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: buyer,
				}},
			},
			wantLaterBalance:        3,
			wantPruneHistoricalData: false,
		},
		{
			name:            "can renew owned Dym-Name, expired",
			buyer:           buyer,
			originalBalance: extendsPrice*2 + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   1,
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "when renew previously-owned expired Dym-Name, reset config",
			buyer:           buyer,
			originalBalance: extendsPrice*2 + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   5,
				Configs: []dymnstypes.DymNameConfig{{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: buyer,
				}},
			},
			setupHistoricalData: true,
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
				Configs:    nil,
			},
			wantLaterBalance:        3,
			wantPruneHistoricalData: true,
		},
		{
			name:            "can take over an expired Dym-Name after grace period has passed",
			buyer:           buyer,
			originalBalance: firstYearPrice5PlusL + extendsPrice + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   1,
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "take over an expired when ownership changed, reset config",
			buyer:           buyer,
			originalBalance: firstYearPrice5PlusL + extendsPrice + 3,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   1,
				Configs: []dymnstypes.DymNameConfig{{
					Type:  dymnstypes.DymNameConfigType_NAME,
					Value: buyer,
				}},
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 86400*365*2,
				Configs:    nil,
			},
			wantLaterBalance: 3,
		},
		{
			name:            "not enough balance to take over an expired Dym-Name after grace period has passed",
			buyer:           buyer,
			originalBalance: 1,
			duration:        2,
			existingDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   3,
			},
			wantLaterDymName: &dymnstypes.DymName{
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   3,
			},
			wantErr:          true,
			wantErrContains:  "insufficient funds",
			wantLaterBalance: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, bk, ctx := setupTest()

			err := bk.MintCoins(ctx, dymnstypes.ModuleName, dymnsutils.TestCoins(originalModuleBalance))
			require.NoError(t, err)

			if tt.originalBalance > 0 {
				coin := dymnsutils.TestCoins(tt.originalBalance)
				err := bk.MintCoins(ctx, dymnstypes.ModuleName, coin)
				require.NoError(t, err)
				err = bk.SendCoinsFromModuleToAccount(
					ctx,
					dymnstypes.ModuleName, sdk.MustAccAddressFromBech32(tt.buyer),
					coin,
				)
				require.NoError(t, err)
			}

			useRecordName := "bonded-pool"
			if tt.customDymName != "" {
				useRecordName = tt.customDymName
			}

			if tt.existingDymName != nil {
				tt.existingDymName.Name = useRecordName
				err := dk.SetDymName(ctx, *tt.existingDymName)
				require.NoError(t, err)

				if tt.setupHistoricalData {
					opo1 := dymnstypes.OpenPurchaseOrder{
						Name:     useRecordName,
						ExpireAt: now.Unix() - 1,
						MinPrice: dymnsutils.TestCoin(1),
					}
					err := dk.SetOpenPurchaseOrder(ctx, opo1)
					require.NoError(t, err)

					err = dk.MoveOpenPurchaseOrderToHistorical(ctx, useRecordName)
					require.NoError(t, err)

					opo2 := dymnstypes.OpenPurchaseOrder{
						Name:      useRecordName,
						ExpireAt:  tt.existingDymName.ExpireAt - 1,
						MinPrice:  dymnsutils.TestCoin(1),
						SellPrice: dymnsutils.TestCoinP(2),
						HighestBid: &dymnstypes.OpenPurchaseOrderBid{
							Bidder: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
							Price:  dymnsutils.TestCoin(2),
						},
					}
					err = dk.SetOpenPurchaseOrder(ctx, opo2)
					require.NoError(t, err)

					require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, useRecordName), 1)
				}
			} else {
				require.False(t, tt.setupHistoricalData, "bad setup testcase")
			}
			if tt.wantLaterDymName != nil {
				tt.wantLaterDymName.Name = useRecordName
			}

			resp, err := dymnskeeper.NewMsgServerImpl(dk).RegisterName(ctx, &dymnstypes.MsgRegisterName{
				Name:     useRecordName,
				Duration: tt.duration,
				Owner:    tt.buyer,
			})
			laterDymName := dk.GetDymName(ctx, useRecordName)

			defer func() {
				laterBalance := bk.GetBalance(ctx, sdk.MustAccAddressFromBech32(tt.buyer), denom).Amount.Int64()
				require.Equal(t, tt.wantLaterBalance, laterBalance)
			}()

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)

				require.Nil(t, resp)

				defer func() {
					laterModuleBalance := bk.GetBalance(
						ctx,
						authtypes.NewModuleAddress(dymnstypes.ModuleName), denom,
					).Amount.Int64()
					require.Equal(t, originalModuleBalance, laterModuleBalance, "module account balance should not be changed")
				}()

				if tt.existingDymName != nil {
					require.Equal(t, *tt.existingDymName, *laterDymName, "should not change existing record")
					require.NotNil(t, tt.wantLaterDymName, "bad setup testcase")
					require.Equal(t, *tt.wantLaterDymName, *laterDymName)
				} else {
					require.Nil(t, laterDymName)
					require.Nil(t, tt.wantLaterDymName, "bad setup testcase")
				}

				if tt.setupHistoricalData {
					require.NotNil(t, dk.GetOpenPurchaseOrder(ctx, useRecordName), "open purchase order must be kept")
					require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, useRecordName), 1, "historical data must be kept")
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			defer func() {
				laterModuleBalance := bk.GetBalance(
					ctx,
					authtypes.NewModuleAddress(dymnstypes.ModuleName), denom,
				).Amount.Int64()
				require.Equal(t, originalModuleBalance, laterModuleBalance, "token should be burned")
			}()

			require.NotNil(t, laterDymName)
			require.NotNil(t, tt.wantLaterDymName, "bad setup testcase")
			require.Equal(t, *tt.wantLaterDymName, *laterDymName)

			if tt.setupHistoricalData {
				if tt.wantPruneHistoricalData {
					require.Nil(t, dk.GetOpenPurchaseOrder(ctx, useRecordName), "open purchase order must be pruned")
					require.Empty(t, dk.GetHistoricalOpenPurchaseOrders(ctx, useRecordName), "historical data must be pruned")

					if tt.existingDymName.Owner != laterDymName.Owner {
						ownedByPreviousOwner, err := dk.GetDymNamesOwnedBy(ctx, tt.existingDymName.Owner, now.Unix())
						require.NoError(t, err)
						require.Empty(t, ownedByPreviousOwner, "reverse mapping should be removed")
					}
				} else {
					require.NotNil(t, dk.GetOpenPurchaseOrder(ctx, useRecordName), "open purchase order must be kept")
					require.Len(t, dk.GetHistoricalOpenPurchaseOrders(ctx, useRecordName), 1, "historical data must be kept")
				}
			} else {
				require.False(t, tt.wantPruneHistoricalData, "bad setup testcase")
			}

			ownedByBuyer, err := dk.GetDymNamesOwnedBy(ctx, tt.buyer, now.Unix())
			require.NoError(t, err)
			require.Len(t, ownedByBuyer, 1, "reverse mapping should be set")
		})
	}
}
