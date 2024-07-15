package types

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSellOrder_GetIdentity(t *testing.T) {
	m := &SellOrder{
		Name:     "aabb",
		ExpireAt: 1234,
	}
	require.Equal(t, "aabb|1234", m.GetIdentity())
}

func TestSellOrder_HasSetSellPrice(t *testing.T) {
	require.False(t, (&SellOrder{
		SellPrice: nil,
	}).HasSetSellPrice())
	require.False(t, (&SellOrder{
		SellPrice: &sdk.Coin{},
	}).HasSetSellPrice())
	require.False(t, (&SellOrder{
		SellPrice: dymnsutils.TestCoinP(0),
	}).HasSetSellPrice())
	require.True(t, (&SellOrder{
		SellPrice: dymnsutils.TestCoinP(1),
	}).HasSetSellPrice())
}

func TestSellOrder_HasExpiredAtCtx(t *testing.T) {
	var epoch int64 = 2
	ctx := sdk.Context{}.WithBlockHeader(tmproto.Header{Time: time.Unix(2, 0)})
	require.True(t, (&SellOrder{
		ExpireAt: epoch - 1,
	}).HasExpiredAtCtx(ctx))
	require.False(t, (&SellOrder{
		ExpireAt: epoch + 1,
	}).HasExpiredAtCtx(ctx))
	require.False(t, (&SellOrder{
		ExpireAt: epoch,
	}).HasExpiredAtCtx(ctx), "SO expires after expires at")
}

func TestSellOrder_HasExpired(t *testing.T) {
	var epoch int64 = 2
	require.True(t, (&SellOrder{
		ExpireAt: epoch - 1,
	}).HasExpired(epoch))
	require.False(t, (&SellOrder{
		ExpireAt: epoch + 1,
	}).HasExpired(epoch))
	require.False(t, (&SellOrder{
		ExpireAt: epoch,
	}).HasExpired(epoch), "SO expires after expires at")
}

func TestSellOrder_HasFinished(t *testing.T) {
	oneCoin := dymnsutils.TestCoin(1)
	threeCoin := dymnsutils.TestCoin(3)
	zeroCoin := dymnsutils.TestCoin(0)

	nowEpoch := time.Now().Unix()

	tests := []struct {
		name       string
		ExpireAt   int64
		SellPrice  *sdk.Coin
		HighestBid *SellOrderBid
		want       bool
	}{
		{
			name:       "expired, without sell-price, without bid",
			ExpireAt:   nowEpoch - 1,
			SellPrice:  &zeroCoin,
			HighestBid: nil,
			want:       true,
		},
		{
			name:       "expired, without sell-price, without bid",
			ExpireAt:   nowEpoch - 1,
			SellPrice:  nil,
			HighestBid: nil,
			want:       true,
		},
		{
			name:       "expired, + sell-price, without bid",
			ExpireAt:   nowEpoch - 1,
			SellPrice:  &threeCoin,
			HighestBid: nil,
			want:       true,
		},
		{
			name:      "expired, + sell-price, + bid (under sell-price)",
			ExpireAt:  nowEpoch - 1,
			SellPrice: &threeCoin,
			HighestBid: &SellOrderBid{
				Bidder: "x",
				Price:  oneCoin,
			},
			want: true,
		},
		{
			name:      "expired, + sell-price, + bid (= sell-price)",
			ExpireAt:  nowEpoch - 1,
			SellPrice: &threeCoin,
			HighestBid: &SellOrderBid{
				Bidder: "x",
				Price:  threeCoin,
			},
			want: true,
		},
		{
			name:       "not expired, without sell-price, without bid",
			ExpireAt:   nowEpoch + 1,
			SellPrice:  &zeroCoin,
			HighestBid: nil,
			want:       false,
		},
		{
			name:       "not expired, without sell-price, without bid",
			ExpireAt:   nowEpoch + 1,
			SellPrice:  nil,
			HighestBid: nil,
			want:       false,
		},
		{
			name:       "not expired, + sell-price, without bid",
			ExpireAt:   nowEpoch + 1,
			SellPrice:  &threeCoin,
			HighestBid: nil,
			want:       false,
		},
		{
			name:      "not expired, + sell-price, + bid (under sell-price)",
			ExpireAt:  nowEpoch + 1,
			SellPrice: &threeCoin,
			HighestBid: &SellOrderBid{
				Bidder: "x",
				Price:  oneCoin,
			},
			want: false,
		},
		{
			name:      "not expired, + sell-price, + bid (= sell-price)",
			ExpireAt:  nowEpoch + 1,
			SellPrice: &threeCoin,
			HighestBid: &SellOrderBid{
				Bidder: "x",
				Price:  threeCoin,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SellOrder{
				Name:       "a",
				ExpireAt:   tt.ExpireAt,
				MinPrice:   oneCoin,
				SellPrice:  tt.SellPrice,
				HighestBid: tt.HighestBid,
			}

			require.Equal(t, tt.want, m.HasFinishedAtCtx(
				sdk.Context{}.
					WithBlockHeader(
						tmproto.Header{
							Time: time.Unix(nowEpoch, 0),
						},
					)),
			)
			require.Equal(t, tt.want, m.HasFinished(nowEpoch))
		})
	}
}

func TestSellOrder_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*SellOrder)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		Name            string
		ExpireAt        int64
		MinPrice        sdk.Coin
		SellPrice       *sdk.Coin
		HighestBid      *SellOrderBid
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "valid sell order",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().Unix(),
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(1),
			HighestBid: &SellOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  dymnsutils.TestCoin(1),
			},
		},
		{
			name:      "valid sell order without bid",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().Unix(),
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(1),
		},
		{
			name:     "valid sell order without setting sell price",
			Name:     "bonded-pool",
			ExpireAt: time.Now().Unix(),
			MinPrice: dymnsutils.TestCoin(1),
		},
		{
			name:            "empty name",
			Name:            "",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Dym-Name of SO is empty",
		},
		{
			name:            "bad name",
			Name:            "-bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Dym-Name of SO is not a valid dym name",
		},
		{
			name:            "empty time",
			Name:            "bonded-pool",
			ExpireAt:        0,
			MinPrice:        dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "SO expiry is empty",
		},
		{
			name:            "min price is zero",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(0),
			wantErr:         true,
			wantErrContains: "SO min price is zero",
		},
		{
			name:            "min price is empty",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        sdk.Coin{},
			wantErr:         true,
			wantErrContains: "SO min price is zero",
		},
		{
			name:            "min price is negative",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(-1),
			wantErr:         true,
			wantErrContains: "SO min price is negative",
		},
		{
			name:     "min price is invalid",
			Name:     "bonded-pool",
			ExpireAt: time.Now().Unix(),
			MinPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "SO min price is invalid",
		},
		{
			name:            "sell price is negative",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(1),
			SellPrice:       dymnsutils.TestCoinP(-1),
			wantErr:         true,
			wantErrContains: "SO sell price is negative",
		},
		{
			name:     "sell price is invalid",
			Name:     "bonded-pool",
			ExpireAt: time.Now().Unix(),
			MinPrice: dymnsutils.TestCoin(1),
			SellPrice: &sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "SO sell price is invalid",
		},
		{
			name:            "sell price is less than min price",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(2),
			SellPrice:       dymnsutils.TestCoinP(1),
			wantErr:         true,
			wantErrContains: "SO sell price is less than min price",
		},
		{
			name:            "sell price denom must match min price denom",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().Unix(),
			MinPrice:        dymnsutils.TestCoin(1),
			SellPrice:       dymnsutils.TestCoin2P(sdk.NewInt64Coin("u"+params.BaseDenom, 2)),
			wantErr:         true,
			wantErrContains: "SO sell price denom is different from min price denom",
		},
		{
			name:      "invalid highest bid",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().Unix(),
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(1),
			HighestBid: &SellOrderBid{
				Bidder: "0x1",
				Price:  dymnsutils.TestCoin(1),
			},
			wantErr:         true,
			wantErrContains: "SO bidder is not a valid bech32 account address",
		},
		{
			name:      "highest bid < min price",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().Unix(),
			MinPrice:  dymnsutils.TestCoin(2),
			SellPrice: dymnsutils.TestCoinP(3),
			HighestBid: &SellOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  dymnsutils.TestCoin(1),
			},
			wantErr:         true,
			wantErrContains: "SO highest bid price is less than min price",
		},
		{
			name:      "highest bid > sell price",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().Unix(),
			MinPrice:  dymnsutils.TestCoin(2),
			SellPrice: dymnsutils.TestCoinP(3),
			HighestBid: &SellOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  dymnsutils.TestCoin(4),
			},
			wantErr:         true,
			wantErrContains: "SO sell price is less than highest bid price",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SellOrder{
				Name:       tt.Name,
				ExpireAt:   tt.ExpireAt,
				MinPrice:   tt.MinPrice,
				SellPrice:  tt.SellPrice,
				HighestBid: tt.HighestBid,
			}

			err := m.Validate()
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

func TestSellOrderBid_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*SellOrderBid)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		Bidder          string
		Price           sdk.Coin
		wantErr         bool
		wantErrContains string
	}{
		{
			name:   "valid sell order bid",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:  dymnsutils.TestCoin(1),
		},
		{
			name:            "empty bidder",
			Bidder:          "",
			Price:           dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "SO bidder is empty",
		},
		{
			name:            "bad bidder",
			Bidder:          "0x1",
			Price:           dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "SO bidder is not a valid bech32 account address",
		},
		{
			name:            "zero price",
			Bidder:          "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:           dymnsutils.TestCoin(0),
			wantErr:         true,
			wantErrContains: "SO bid price is zero",
		},
		{
			name:            "zero price",
			Bidder:          "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:           sdk.Coin{},
			wantErr:         true,
			wantErrContains: "SO bid price is zero",
		},
		{
			name:   "negative price",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price: sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(-1),
			},
			wantErr:         true,
			wantErrContains: "SO bid price is negative",
		},
		{
			name:   "invalid price",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "SO bid price is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SellOrderBid{
				Bidder: tt.Bidder,
				Price:  tt.Price,
			}
			err := m.Validate()
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

func TestHistoricalSellOrders_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*HistoricalSellOrders)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		SellOrders      []SellOrder
		wantErr         bool
		wantErrContains string
	}{
		{
			name: "valid",
			SellOrders: []SellOrder{
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  dymnsutils.TestCoin(1),
					SellPrice: dymnsutils.TestCoinP(1),
				},
				{
					Name:     "a",
					ExpireAt: 2,
					MinPrice: dymnsutils.TestCoin(1),
				},
			},
		},
		{
			name:       "allow empty",
			SellOrders: []SellOrder{},
		},
		{
			name: "reject if SO element is invalid",
			SellOrders: []SellOrder{
				{
					Name:     "a",
					ExpireAt: 1,
					MinPrice: dymnsutils.TestCoin(0), // invalid
				},
				{
					Name:     "a",
					ExpireAt: 2,
					MinPrice: dymnsutils.TestCoin(1),
				},
			},
			wantErr:         true,
			wantErrContains: "SO min price is zero",
		},
		{
			name: "reject if duplicated SO",
			SellOrders: []SellOrder{
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  dymnsutils.TestCoin(1),
					SellPrice: dymnsutils.TestCoinP(1),
				},
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  dymnsutils.TestCoin(1),
					SellPrice: dymnsutils.TestCoinP(1),
				},
			},
			wantErr:         true,
			wantErrContains: "historical SO is not unique",
		},
		{
			name: "reject if SO element has different Dym-Name",
			SellOrders: []SellOrder{
				{
					Name:      "aaa",
					ExpireAt:  1,
					MinPrice:  dymnsutils.TestCoin(1),
					SellPrice: dymnsutils.TestCoinP(1),
				},
				{
					Name:     "bbb",
					ExpireAt: 2,
					MinPrice: dymnsutils.TestCoin(1),
				},
			},
			wantErr:         true,
			wantErrContains: "historical SOs have different Dym-Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &HistoricalSellOrders{
				SellOrders: tt.SellOrders,
			}

			err := m.Validate()
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

func TestSellOrder_GetSdkEvent(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		event := SellOrder{
			Name:      "a",
			ExpireAt:  123456,
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(3),
			HighestBid: &SellOrderBid{
				Bidder: "d",
				Price:  dymnsutils.TestCoin(2),
			},
		}.GetSdkEvent("action-name")
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameSellOrder, event.Type)
		require.Len(t, event.Attributes, 7)
		require.Equal(t, AttributeKeyDymNameSoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameSoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameSoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameSoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "3"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidder, event.Attributes[4].Key)
		require.Equal(t, "d", event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "2"+params.BaseDenom, event.Attributes[5].Value)
		require.Equal(t, AttributeKeyDymNameSoActionName, event.Attributes[6].Key)
		require.Equal(t, "action-name", event.Attributes[6].Value)
	})

	t.Run("no sell-price", func(t *testing.T) {
		event := SellOrder{
			Name:     "a",
			ExpireAt: 123456,
			MinPrice: dymnsutils.TestCoin(1),
			HighestBid: &SellOrderBid{
				Bidder: "d",
				Price:  dymnsutils.TestCoin(2),
			},
		}.GetSdkEvent("action-name")
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameSellOrder, event.Type)
		require.Len(t, event.Attributes, 7)
		require.Equal(t, AttributeKeyDymNameSoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameSoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameSoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameSoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "0"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidder, event.Attributes[4].Key)
		require.Equal(t, "d", event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "2"+params.BaseDenom, event.Attributes[5].Value)
		require.Equal(t, AttributeKeyDymNameSoActionName, event.Attributes[6].Key)
		require.Equal(t, "action-name", event.Attributes[6].Value)
	})
	t.Run("no highest bid", func(t *testing.T) {
		event := SellOrder{
			Name:      "a",
			ExpireAt:  123456,
			MinPrice:  dymnsutils.TestCoin(1),
			SellPrice: dymnsutils.TestCoinP(3),
		}.GetSdkEvent("action-name")
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameSellOrder, event.Type)
		require.Len(t, event.Attributes, 7)
		require.Equal(t, AttributeKeyDymNameSoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameSoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameSoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameSoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "3"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidder, event.Attributes[4].Key)
		require.Empty(t, event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameSoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "0"+params.BaseDenom, event.Attributes[5].Value)
		require.Equal(t, AttributeKeyDymNameSoActionName, event.Attributes[6].Key)
		require.Equal(t, "action-name", event.Attributes[6].Value)
	})
}
