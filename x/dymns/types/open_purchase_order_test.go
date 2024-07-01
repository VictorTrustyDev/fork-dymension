package types

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOpenPurchaseOrder_Identity(t *testing.T) {
	m := &OpenPurchaseOrder{
		Name:     "aabb",
		ExpireAt: 1234,
	}
	require.Equal(t, "aabb|1234", m.Identity())
}

func TestOpenPurchaseOrder_HasSetSellPrice(t *testing.T) {
	require.False(t, (&OpenPurchaseOrder{
		SellPrice: sdk.Coin{},
	}).HasSetSellPrice())
	require.False(t, (&OpenPurchaseOrder{
		SellPrice: sdk.NewCoin("adym", sdk.ZeroInt()),
	}).HasSetSellPrice())
	require.True(t, (&OpenPurchaseOrder{
		SellPrice: sdk.NewCoin("adym", sdk.OneInt()),
	}).HasSetSellPrice())
}

func TestOpenPurchaseOrder_HasExpiredAtCtx(t *testing.T) {
	var epoch int64 = 2
	ctx := sdk.Context{}.WithBlockHeader(tmproto.Header{Time: time.Unix(2, 0)})
	require.True(t, (&OpenPurchaseOrder{
		ExpireAt: epoch - 1,
	}).HasExpiredAtCtx(ctx))
	require.False(t, (&OpenPurchaseOrder{
		ExpireAt: epoch + 1,
	}).HasExpiredAtCtx(ctx))
	require.False(t, (&OpenPurchaseOrder{
		ExpireAt: epoch,
	}).HasExpiredAtCtx(ctx), "OPO expires after expires at")
}

func TestOpenPurchaseOrder_HasExpired(t *testing.T) {
	var epoch int64 = 2
	require.True(t, (&OpenPurchaseOrder{
		ExpireAt: epoch - 1,
	}).HasExpired(epoch))
	require.False(t, (&OpenPurchaseOrder{
		ExpireAt: epoch + 1,
	}).HasExpired(epoch))
	require.False(t, (&OpenPurchaseOrder{
		ExpireAt: epoch,
	}).HasExpired(epoch), "OPO expires after expires at")
}

func TestOpenPurchaseOrder_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*OpenPurchaseOrder)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		Name            string
		ExpireAt        int64
		MinPrice        sdk.Coin
		SellPrice       sdk.Coin
		HighestBid      *OpenPurchaseOrderBid
		wantErr         bool
		wantErrContains string
	}{
		{
			name:      "valid open purchase order",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().UTC().Unix(),
			MinPrice:  sdk.NewCoin("adym", sdk.OneInt()),
			SellPrice: sdk.NewCoin("adym", sdk.OneInt()),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  sdk.NewCoin("adym", sdk.OneInt()),
			},
		},
		{
			name:      "valid open purchase order without bid",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().UTC().Unix(),
			MinPrice:  sdk.NewCoin("adym", sdk.OneInt()),
			SellPrice: sdk.NewCoin("adym", sdk.OneInt()),
		},
		{
			name:     "valid open purchase order without setting sell price",
			Name:     "bonded-pool",
			ExpireAt: time.Now().UTC().Unix(),
			MinPrice: sdk.NewCoin("adym", sdk.OneInt()),
		},
		{
			name:            "empty name",
			Name:            "",
			ExpireAt:        time.Now().UTC().Unix(),
			MinPrice:        sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "Dym-Name of OPO is empty",
		},
		{
			name:            "bad name",
			Name:            "-bonded-pool",
			ExpireAt:        time.Now().UTC().Unix(),
			MinPrice:        sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "Dym-Name of OPO is not a valid dym name",
		},
		{
			name:            "empty time",
			Name:            "bonded-pool",
			ExpireAt:        0,
			MinPrice:        sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "OPO expiry is empty",
		},
		{
			name:            "min price is zero",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().UTC().Unix(),
			MinPrice:        sdk.NewCoin("adym", sdk.ZeroInt()),
			wantErr:         true,
			wantErrContains: "OPO min price is zero",
		},
		{
			name:            "min price is empty",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().UTC().Unix(),
			MinPrice:        sdk.Coin{},
			wantErr:         true,
			wantErrContains: "OPO min price is zero",
		},
		{
			name:     "min price is negative",
			Name:     "bonded-pool",
			ExpireAt: time.Now().UTC().Unix(),
			MinPrice: sdk.Coin{
				Denom:  "adym",
				Amount: sdk.NewInt(-1),
			},
			wantErr:         true,
			wantErrContains: "OPO min price is negative",
		},
		{
			name:     "min price is invalid",
			Name:     "bonded-pool",
			ExpireAt: time.Now().UTC().Unix(),
			MinPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "OPO min price is invalid",
		},
		{
			name:     "sell price is negative",
			Name:     "bonded-pool",
			ExpireAt: time.Now().UTC().Unix(),
			MinPrice: sdk.NewCoin("adym", sdk.OneInt()),
			SellPrice: sdk.Coin{
				Denom:  "adym",
				Amount: sdk.NewInt(-1),
			},
			wantErr:         true,
			wantErrContains: "OPO sell price is negative",
		},
		{
			name:     "sell price is invalid",
			Name:     "bonded-pool",
			ExpireAt: time.Now().UTC().Unix(),
			MinPrice: sdk.NewCoin("adym", sdk.OneInt()),
			SellPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "OPO sell price is invalid",
		},
		{
			name:            "sell price is less than min price",
			Name:            "bonded-pool",
			ExpireAt:        time.Now().UTC().Unix(),
			MinPrice:        sdk.NewCoin("adym", sdk.NewInt(2)),
			SellPrice:       sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "OPO sell price is less than min price",
		},
		{
			name:      "invalid highest bid",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().UTC().Unix(),
			MinPrice:  sdk.NewCoin("adym", sdk.OneInt()),
			SellPrice: sdk.NewCoin("adym", sdk.OneInt()),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "0x1",
				Price:  sdk.NewCoin("adym", sdk.OneInt()),
			},
			wantErr:         true,
			wantErrContains: "OPO bidder is not a valid bech32 account address",
		},
		{
			name:      "highest bid < min price",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().UTC().Unix(),
			MinPrice:  sdk.NewCoin("adym", sdk.NewInt(2)),
			SellPrice: sdk.NewCoin("adym", sdk.NewInt(3)),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  sdk.NewCoin("adym", sdk.OneInt()),
			},
			wantErr:         true,
			wantErrContains: "OPO highest bid price is less than min price",
		},
		{
			name:      "highest bid > sell price",
			Name:      "bonded-pool",
			ExpireAt:  time.Now().UTC().Unix(),
			MinPrice:  sdk.NewCoin("adym", sdk.NewInt(2)),
			SellPrice: sdk.NewCoin("adym", sdk.NewInt(3)),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Price:  sdk.NewCoin("adym", sdk.NewInt(4)),
			},
			wantErr:         true,
			wantErrContains: "OPO sell price is less than highest bid price",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &OpenPurchaseOrder{
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

func TestOpenPurchaseOrderBid_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*OpenPurchaseOrderBid)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		Bidder          string
		Price           sdk.Coin
		wantErr         bool
		wantErrContains string
	}{
		{
			name:   "valid open purchase order bid",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:  sdk.NewCoin("adym", sdk.OneInt()),
		},
		{
			name:            "empty bidder",
			Bidder:          "",
			Price:           sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "OPO bidder is empty",
		},
		{
			name:            "bad bidder",
			Bidder:          "0x1",
			Price:           sdk.NewCoin("adym", sdk.OneInt()),
			wantErr:         true,
			wantErrContains: "OPO bidder is not a valid bech32 account address",
		},
		{
			name:            "zero price",
			Bidder:          "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:           sdk.NewCoin("adym", sdk.ZeroInt()),
			wantErr:         true,
			wantErrContains: "OPO bid price is zero",
		},
		{
			name:            "zero price",
			Bidder:          "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price:           sdk.Coin{},
			wantErr:         true,
			wantErrContains: "OPO bid price is zero",
		},
		{
			name:   "negative price",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price: sdk.Coin{
				Denom:  "adym",
				Amount: sdk.NewInt(-1),
			},
			wantErr:         true,
			wantErrContains: "OPO bid price is negative",
		},
		{
			name:   "invalid price",
			Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Price: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "OPO bid price is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &OpenPurchaseOrderBid{
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

func TestHistoricalOpenPurchaseOrders_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*HistoricalOpenPurchaseOrders)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name               string
		OpenPurchaseOrders []OpenPurchaseOrder
		wantErr            bool
		wantErrContains    string
	}{
		{
			name: "valid",
			OpenPurchaseOrders: []OpenPurchaseOrder{
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
					SellPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
				{
					Name:     "a",
					ExpireAt: 2,
					MinPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
			},
		},
		{
			name:               "allow empty",
			OpenPurchaseOrders: []OpenPurchaseOrder{},
		},
		{
			name: "reject if OPO element is invalid",
			OpenPurchaseOrders: []OpenPurchaseOrder{
				{
					Name:     "a",
					ExpireAt: 1,
					MinPrice: sdk.NewCoin(params.BaseDenom, sdk.ZeroInt()), // invalid
				},
				{
					Name:     "a",
					ExpireAt: 2,
					MinPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
			},
			wantErr:         true,
			wantErrContains: "OPO min price is zero",
		},
		{
			name: "reject if duplicated OPO",
			OpenPurchaseOrders: []OpenPurchaseOrder{
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
					SellPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
				{
					Name:      "a",
					ExpireAt:  1,
					MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
					SellPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
			},
			wantErr:         true,
			wantErrContains: "historical OPO is not unique",
		},
		{
			name: "reject if OPO element has different Dym-Name",
			OpenPurchaseOrders: []OpenPurchaseOrder{
				{
					Name:      "aaa",
					ExpireAt:  1,
					MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
					SellPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
				{
					Name:     "bbb",
					ExpireAt: 2,
					MinPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
				},
			},
			wantErr:         true,
			wantErrContains: "historical OPOs have different Dym-Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &HistoricalOpenPurchaseOrders{
				OpenPurchaseOrders: tt.OpenPurchaseOrders,
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

func TestOpenPurchaseOrder_GetSdkEvent(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		event := OpenPurchaseOrder{
			Name:      "a",
			ExpireAt:  123456,
			MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
			SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(3)),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "d",
				Price:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(2)),
			},
		}.GetSdkEvent()
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameOpenPurchaseOrder, event.Type)
		require.Len(t, event.Attributes, 6)
		require.Equal(t, AttributeKeyDymNameOpoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameOpoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameOpoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameOpoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "3"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidder, event.Attributes[4].Key)
		require.Equal(t, "d", event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "2"+params.BaseDenom, event.Attributes[5].Value)
	})

	t.Run("no sell-price", func(t *testing.T) {
		event := OpenPurchaseOrder{
			Name:     "a",
			ExpireAt: 123456,
			MinPrice: sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
			HighestBid: &OpenPurchaseOrderBid{
				Bidder: "d",
				Price:  sdk.NewCoin(params.BaseDenom, sdk.NewInt(2)),
			},
		}.GetSdkEvent()
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameOpenPurchaseOrder, event.Type)
		require.Len(t, event.Attributes, 6)
		require.Equal(t, AttributeKeyDymNameOpoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameOpoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameOpoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameOpoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "0"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidder, event.Attributes[4].Key)
		require.Equal(t, "d", event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "2"+params.BaseDenom, event.Attributes[5].Value)
	})
	t.Run("no highest bid", func(t *testing.T) {
		event := OpenPurchaseOrder{
			Name:      "a",
			ExpireAt:  123456,
			MinPrice:  sdk.NewCoin(params.BaseDenom, sdk.OneInt()),
			SellPrice: sdk.NewCoin(params.BaseDenom, sdk.NewInt(3)),
		}.GetSdkEvent()
		require.NotNil(t, event)
		require.Equal(t, EventTypeDymNameOpenPurchaseOrder, event.Type)
		require.Len(t, event.Attributes, 6)
		require.Equal(t, AttributeKeyDymNameOpoName, event.Attributes[0].Key)
		require.Equal(t, "a", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyDymNameOpoExpiryEpoch, event.Attributes[1].Key)
		require.Equal(t, "123456", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyDymNameOpoMinPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyDymNameOpoSellPrice, event.Attributes[3].Key)
		require.Equal(t, "3"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidder, event.Attributes[4].Key)
		require.Empty(t, event.Attributes[4].Value)
		require.Equal(t, AttributeKeyDymNameOpoHighestBidPrice, event.Attributes[5].Key)
		require.Equal(t, "0"+params.BaseDenom, event.Attributes[5].Value)
	})
}
