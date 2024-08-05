package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
)

func TestMsgPlaceSellOrder_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		goodsId         string
		orderType       MarketOrderType
		minPrice        sdk.Coin
		sellPrice       *sdk.Coin
		owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "pass - (Name) valid sell order",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  dymnsutils.TestCoin(1),
			sellPrice: dymnsutils.TestCoinP(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "pass - (Alias) valid sell order",
			goodsId:   "alias",
			orderType: MarketOrderType_MOT_ALIAS,
			minPrice:  dymnsutils.TestCoin(1),
			sellPrice: dymnsutils.TestCoinP(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "pass - (Name) valid sell order without bid",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  dymnsutils.TestCoin(1),
			sellPrice: dymnsutils.TestCoinP(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "pass - (Alias) valid sell order without bid",
			goodsId:   "alias",
			orderType: MarketOrderType_MOT_ALIAS,
			minPrice:  dymnsutils.TestCoin(1),
			sellPrice: dymnsutils.TestCoinP(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "pass - (Name) valid sell order without setting sell price",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  dymnsutils.TestCoin(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "pass - (Alias) valid sell order without setting sell price",
			goodsId:   "alias",
			orderType: MarketOrderType_MOT_ALIAS,
			minPrice:  dymnsutils.TestCoin(1),
			owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "fail - (Name) empty name",
			goodsId:         "",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "fail - (Alias) empty alias",
			goodsId:         "",
			orderType:       MarketOrderType_MOT_ALIAS,
			minPrice:        dymnsutils.TestCoin(1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "alias is not a valid alias",
		},
		{
			name:            "fail - (Name) bad name",
			goodsId:         "-my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "fail - (Alias) bad alias",
			goodsId:         "bad-alias",
			orderType:       MarketOrderType_MOT_ALIAS,
			minPrice:        dymnsutils.TestCoin(1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "alias is not a valid alias",
		},
		{
			name:            "fail - min price is zero",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(0),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO min price is zero",
		},
		{
			name:            "fail - min price is empty",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        sdk.Coin{},
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO min price is zero",
		},
		{
			name:      "fail - min price is negative",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice: sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(-1),
			},
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO min price is negative",
		},
		{
			name:      "fail - min price is invalid",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO min price is invalid",
		},
		{
			name:            "fail - sell price is negative",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(1),
			sellPrice:       dymnsutils.TestCoinP(-1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO sell price is negative",
		},
		{
			name:      "fail - sell price is invalid",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  dymnsutils.TestCoin(1),
			sellPrice: &sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO sell price is invalid",
		},
		{
			name:            "fail - sell price is less than min price",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(2),
			sellPrice:       dymnsutils.TestCoinP(1),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO sell price is less than min price",
		},
		{
			name:            "fail - sell price denom must match min price denom",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(2),
			sellPrice:       dymnsutils.TestCoin2P(sdk.NewCoin("u"+params.BaseDenom, sdk.OneInt())),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "SO sell price denom is different from min price denom",
		},
		{
			name:            "fail - missing owner",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(2),
			owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "fail - invalid owner",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(2),
			owner:           "dym1fl48vsnmsdzcv85",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "fail - owner must be dym1",
			goodsId:         "my-name",
			orderType:       MarketOrderType_MOT_DYM_NAME,
			minPrice:        dymnsutils.TestCoin(2),
			owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "fail - reject unknown order type",
			goodsId:         "goods",
			orderType:       MarketOrderType_MOT_UNKNOWN,
			minPrice:        dymnsutils.TestCoin(2),
			owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "invalid order type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgPlaceSellOrder{
				GoodsId:   tt.goodsId,
				OrderType: tt.orderType,
				MinPrice:  tt.minPrice,
				SellPrice: tt.sellPrice,
				Owner:     tt.owner,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgPlaceSellOrder_ToSellOrder(t *testing.T) {
	validMinPrice := dymnsutils.TestCoin(1)
	validSellPrice := dymnsutils.TestCoin(1)

	tests := []struct {
		name      string
		goodsId   string
		orderType MarketOrderType
		minPrice  sdk.Coin
		sellPrice *sdk.Coin
		Owner     string
		want      SellOrder
	}{
		{
			name:      "normal Dym-Name sell order",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  validMinPrice,
			sellPrice: &validSellPrice,
			Owner:     "",
			want: SellOrder{
				GoodsId:   "my-name",
				Type:      MarketOrderType_MOT_DYM_NAME,
				MinPrice:  validMinPrice,
				SellPrice: &validSellPrice,
			},
		},
		{
			name:      "normal Alias sell order",
			goodsId:   "alias",
			orderType: MarketOrderType_MOT_ALIAS,
			minPrice:  validMinPrice,
			sellPrice: &validSellPrice,
			Owner:     "",
			want: SellOrder{
				GoodsId:   "alias",
				Type:      MarketOrderType_MOT_ALIAS,
				MinPrice:  validMinPrice,
				SellPrice: &validSellPrice,
			},
		},
		{
			name:      "without sell price",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  validMinPrice,
			sellPrice: nil,
			Owner:     "",
			want: SellOrder{
				GoodsId:   "my-name",
				Type:      MarketOrderType_MOT_DYM_NAME,
				MinPrice:  validMinPrice,
				SellPrice: nil,
			},
		},
		{
			name:      "without sell price, auto omit zero sell price",
			goodsId:   "my-name",
			orderType: MarketOrderType_MOT_DYM_NAME,
			minPrice:  validMinPrice,
			sellPrice: dymnsutils.TestCoin2P(sdk.NewCoin(validMinPrice.Denom, sdk.ZeroInt())),
			Owner:     "",
			want: SellOrder{
				GoodsId:   "my-name",
				Type:      MarketOrderType_MOT_DYM_NAME,
				MinPrice:  validMinPrice,
				SellPrice: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgPlaceSellOrder{
				GoodsId:   tt.goodsId,
				OrderType: tt.orderType,
				MinPrice:  tt.minPrice,
				SellPrice: tt.sellPrice,
				Owner:     tt.Owner,
			}

			so := m.ToSellOrder()
			require.Equal(t, tt.want, so)
			require.Zero(t, so.ExpireAt)
			require.Nil(t, so.HighestBid)
		})
	}
}
