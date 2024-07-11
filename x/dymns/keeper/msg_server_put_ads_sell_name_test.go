package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
func Test_msgServer_PutAdsSellName(t *testing.T) {
	now := time.Now().UTC()

	const daysProhibitSell = 30
	const daysOpenPurchaseOrderDuration = 7

	setupTest := func() (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
		dk, bk, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		moduleParams := dk.GetParams(ctx)
		moduleParams.Misc.DaysOpenPurchaseOrderDuration = daysProhibitSell
		moduleParams.Misc.DaysProhibitSell = daysOpenPurchaseOrderDuration
		moduleParams.Misc.GasCrudOpenPurchaseOrder = 20_000_000
		err := dk.SetParams(ctx, moduleParams)
		require.NoError(t, err)

		return dk, bk, ctx
	}

	t.Run("reject if message not pass validate basic", func(t *testing.T) {
		dk, _, ctx := setupTest()

		requireErrorFContains(t, func() error {
			_, err := dymnskeeper.NewMsgServerImpl(dk).PutAdsSellName(ctx, &dymnstypes.MsgPutAdsSellName{})
			return err
		}, dymnstypes.ErrValidationFailed.Error())
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	coin100 := dymnsutils.TestCoin(100)
	coin200 := dymnsutils.TestCoin(200)
	coin300 := dymnsutils.TestCoin(300)

	tests := []struct {
		name                    string
		withoutDymName          bool
		existingOpo             *dymnstypes.OpenPurchaseOrder
		dymNameExpiryOffsetDays int64
		customOwner             string
		customDymNameOwner      string
		minPrice                sdk.Coin
		sellPrice               *sdk.Coin
		wantErr                 bool
		wantErrContains         string
	}{
		{
			name:            "Dym-Name does not exists",
			withoutDymName:  true,
			minPrice:        coin100,
			wantErr:         true,
			wantErrContains: dymnstypes.ErrDymNameNotFound.Error(),
		},
		{
			name:               "wrong owner",
			customOwner:        owner,
			customDymNameOwner: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
			minPrice:           coin100,
			wantErr:            true,
			wantErrContains:    sdkerrors.ErrUnauthorized.Error(),
		},
		{
			name:                    "expired Dym-Name",
			withoutDymName:          false,
			existingOpo:             nil,
			dymNameExpiryOffsetDays: -1,
			minPrice:                coin100,
			wantErr:                 true,
			wantErrContains:         "Dym-Name is already expired",
		},
		{
			name: "existing active OPO, not finished",
			existingOpo: &dymnstypes.OpenPurchaseOrder{
				ExpireAt:  now.Add(time.Hour).Unix(),
				MinPrice:  coin100,
				SellPrice: &coin200,
			},
			minPrice:        coin100,
			wantErr:         true,
			wantErrContains: "an active Open-Purchase-Order already exists",
		},
		{
			name: "existing active OPO, expired",
			existingOpo: &dymnstypes.OpenPurchaseOrder{
				ExpireAt:  now.Add(-1 * time.Hour).Unix(),
				MinPrice:  coin100,
				SellPrice: &coin200,
			},
			minPrice:        coin100,
			wantErr:         true,
			wantErrContains: "an active expired/completed Open-Purchase-Order already exists ",
		},
		{
			name: "existing active OPO, not expired, completed",
			existingOpo: &dymnstypes.OpenPurchaseOrder{
				ExpireAt:  now.Add(time.Hour).Unix(),
				MinPrice:  coin100,
				SellPrice: &coin200,
				HighestBid: &dymnstypes.OpenPurchaseOrderBid{
					Bidder: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
					Price:  coin200,
				},
			},
			minPrice:        coin100,
			wantErr:         true,
			wantErrContains: "an active expired/completed Open-Purchase-Order already exists ",
		},
		{
			name: "existing active OPO, expired, completed",
			existingOpo: &dymnstypes.OpenPurchaseOrder{
				ExpireAt:  now.Add(-1 * time.Hour).Unix(),
				MinPrice:  coin100,
				SellPrice: &coin200,
				HighestBid: &dymnstypes.OpenPurchaseOrderBid{
					Bidder: "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
					Price:  coin200,
				},
			},
			minPrice:        coin100,
			wantErr:         true,
			wantErrContains: "an active expired/completed Open-Purchase-Order already exists",
		},
		{
			name:            "not allowed denom",
			minPrice:        sdk.NewInt64Coin("u"+params.BaseDenom, 100),
			wantErr:         true,
			wantErrContains: "only adym is allowed as price",
		},
		{
			name:                    "can not sell Dym-Name that almost expired",
			dymNameExpiryOffsetDays: daysProhibitSell + daysOpenPurchaseOrderDuration/2,
			minPrice:                coin100,
			wantErr:                 true,
			wantErrContains:         "days before Dym-Name expiry, can not sell",
		},
		{
			name:                    "successfully place ads for selling Dym-Name, without sell price",
			dymNameExpiryOffsetDays: 9999,
			minPrice:                coin100,
			sellPrice:               nil,
		},
		{
			name:                    "successfully place ads for selling Dym-Name, without sell price",
			dymNameExpiryOffsetDays: 9999,
			minPrice:                coin100,
			sellPrice:               dymnsutils.TestCoinP(0),
		},
		{
			name:                    "successfully place ads for selling Dym-Name, with sell price",
			dymNameExpiryOffsetDays: 9999,
			minPrice:                coin100,
			sellPrice:               &coin300,
		},
		{
			name:                    "successfully place ads for selling Dym-Name, with sell price equals to min-price",
			dymNameExpiryOffsetDays: 9999,
			minPrice:                coin100,
			sellPrice:               &coin100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, _, ctx := setupTest()

			recordName := "bonded-pool"

			useDymNameOwner := owner
			if tt.customDymNameOwner != "" {
				useDymNameOwner = tt.customDymNameOwner
			}
			useDymNameExpiry := ctx.BlockTime().Add(
				time.Hour * 24 * time.Duration(tt.dymNameExpiryOffsetDays),
			).Unix()

			if !tt.withoutDymName {
				dymName := dymnstypes.DymName{
					Name:       recordName,
					Owner:      useDymNameOwner,
					Controller: useDymNameOwner,
					ExpireAt:   useDymNameExpiry,
				}
				err := dk.SetDymName(ctx, dymName)
				require.NoError(t, err)
			}

			if tt.existingOpo != nil {
				tt.existingOpo.Name = recordName
				err := dk.SetOpenPurchaseOrder(ctx, *tt.existingOpo)
				require.NoError(t, err)
			}

			useOwner := owner
			if tt.customOwner != "" {
				useOwner = tt.customOwner
			}
			msg := &dymnstypes.MsgPutAdsSellName{
				Name:      recordName,
				MinPrice:  tt.minPrice,
				SellPrice: tt.sellPrice,
				Owner:     useOwner,
			}
			resp, err := dymnskeeper.NewMsgServerImpl(dk).PutAdsSellName(ctx, msg)
			moduleParams := dk.GetParams(ctx)

			defer func() {
				laterDymName := dk.GetDymName(ctx, recordName)
				if tt.withoutDymName {
					require.Nil(t, laterDymName)
					return
				}

				require.NotNil(t, laterDymName)
				require.Equal(t, dymnstypes.DymName{
					Name:       recordName,
					Owner:      useDymNameOwner,
					Controller: useDymNameOwner,
					ExpireAt:   useDymNameExpiry,
				}, *laterDymName, "Dym-Name record should not be changed in any case")
			}()

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)

				require.Nil(t, resp)

				opo := dk.GetOpenPurchaseOrder(ctx, recordName)
				if tt.existingOpo != nil {
					require.NotNil(t, opo)
					require.Equal(t, *tt.existingOpo, *opo)
				} else {
					require.Nil(t, opo)
				}

				require.Less(t,
					ctx.GasMeter().GasConsumed(), sdk.Gas(moduleParams.Misc.GasCrudOpenPurchaseOrder),
					"should not consume params gas on failed operation",
				)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			opo := dk.GetOpenPurchaseOrder(ctx, recordName)
			require.NotNil(t, opo)

			expectedOpo := dymnstypes.OpenPurchaseOrder{
				Name:       recordName,
				ExpireAt:   ctx.BlockTime().Unix() + int64(moduleParams.Misc.DaysOpenPurchaseOrderDuration)*86400,
				MinPrice:   msg.MinPrice,
				SellPrice:  msg.SellPrice,
				HighestBid: nil,
			}
			if !expectedOpo.HasSetSellPrice() {
				expectedOpo.SellPrice = nil
			}

			require.Nil(t, opo.HighestBid, "highest bid should not be set")

			require.Equal(t, expectedOpo, *opo)

			require.GreaterOrEqual(t,
				ctx.GasMeter().GasConsumed(), sdk.Gas(moduleParams.Misc.GasCrudOpenPurchaseOrder),
				"should consume params gas",
			)

			apoe := dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
			require.Contains(t, apoe.ExpiryByName, recordName)
			require.Equal(t, expectedOpo.ExpireAt, apoe.ExpiryByName[recordName])
		})
	}
}
