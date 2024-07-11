package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgPutAdsSellName_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		Name            string
		MinPrice        sdk.Coin
		SellPrice       *sdk.Coin
		Owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "valid open purchase order",
			Name:      "bonded-pool",
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(1),
			Owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:      "valid open purchase order without bid",
			Name:      "bonded-pool",
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(1),
			Owner:     "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:     "valid open purchase order without setting sell price",
			Name:     "bonded-pool",
			MinPrice: dymnsutils.TestCoin(1),
			Owner:    "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "empty name",
			Name:            "",
			MinPrice:        dymnsutils.TestCoin(1),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "bad name",
			Name:            "-bonded-pool",
			MinPrice:        dymnsutils.TestCoin(1),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "min price is zero",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(0),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO min price is zero",
		},
		{
			name:            "min price is empty",
			Name:            "bonded-pool",
			MinPrice:        sdk.Coin{},
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO min price is zero",
		},
		{
			name: "min price is negative",
			Name: "bonded-pool",
			MinPrice: sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(-1),
			},
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO min price is negative",
		},
		{
			name: "min price is invalid",
			Name: "bonded-pool",
			MinPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO min price is invalid",
		},
		{
			name:            "sell price is negative",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(1),
			SellPrice:       dymnsutils.TestCoinP(-1),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO sell price is negative",
		},
		{
			name:     "sell price is invalid",
			Name:     "bonded-pool",
			MinPrice: dymnsutils.TestCoin(1),
			SellPrice: &sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO sell price is invalid",
		},
		{
			name:            "sell price is less than min price",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(2),
			SellPrice:       dymnsutils.TestCoinP(1),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO sell price is less than min price",
		},
		{
			name:            "sell price denom must match min price denom",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(2),
			SellPrice:       dymnsutils.TestCoin2P(sdk.NewCoin("u"+params.BaseDenom, sdk.OneInt())),
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "OPO sell price denom is different from min price denom",
		},
		{
			name:            "missing owner",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(2),
			Owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "invalid owner",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(2),
			Owner:           "dym1fl48vsnmsdzcv85",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "owner must be dym1",
			Name:            "bonded-pool",
			MinPrice:        dymnsutils.TestCoin(2),
			Owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgPutAdsSellName{
				Name:      tt.Name,
				MinPrice:  tt.MinPrice,
				SellPrice: tt.SellPrice,
				Owner:     tt.Owner,
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

func TestMsgPutAdsSellName_ToOpenPurchaseOrder(t *testing.T) {
	validMinPrice := dymnsutils.TestCoin(1)
	validSellPrice := dymnsutils.TestCoin(1)

	tests := []struct {
		name      string
		Name      string
		MinPrice  sdk.Coin
		SellPrice *sdk.Coin
		Owner     string
		want      OpenPurchaseOrder
	}{
		{
			name:      "valid",
			Name:      "a",
			MinPrice:  validMinPrice,
			SellPrice: &validSellPrice,
			Owner:     "",
			want: OpenPurchaseOrder{
				Name:      "a",
				MinPrice:  validMinPrice,
				SellPrice: &validSellPrice,
			},
		},
		{
			name:      "valid without sell price",
			Name:      "a",
			MinPrice:  validMinPrice,
			SellPrice: nil,
			Owner:     "",
			want: OpenPurchaseOrder{
				Name:      "a",
				MinPrice:  validMinPrice,
				SellPrice: nil,
			},
		},
		{
			name:      "valid without sell price, auto omit zero sell price",
			Name:      "a",
			MinPrice:  validMinPrice,
			SellPrice: dymnsutils.TestCoin2P(sdk.NewCoin(validMinPrice.Denom, sdk.ZeroInt())),
			Owner:     "",
			want: OpenPurchaseOrder{
				Name:      "a",
				MinPrice:  validMinPrice,
				SellPrice: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgPutAdsSellName{
				Name:      tt.Name,
				MinPrice:  tt.MinPrice,
				SellPrice: tt.SellPrice,
				Owner:     tt.Owner,
			}

			opo := m.ToOpenPurchaseOrder()
			require.Equal(t, tt.want, opo)
			require.Zero(t, opo.ExpireAt)
			require.Nil(t, opo.HighestBid)
		})
	}
}
