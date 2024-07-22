package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func Test_queryServer_Params(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)

	params := dk.GetParams(ctx)
	params.Misc.DaysProhibitSell++
	err := dk.SetParams(ctx, params)
	require.NoError(t, err)

	queryServer := dymnskeeper.NewQueryServerImpl(dk)

	resp, err := queryServer.Params(sdk.WrapSDKContext(ctx), &dymnstypes.QueryParamsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, params, resp.Params)
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_DymName(t *testing.T) {
	t.Run("dym name not found", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		queryServer := dymnskeeper.NewQueryServerImpl(dk)
		resp, err := queryServer.DymName(sdk.WrapSDKContext(ctx), &dymnstypes.QueryDymNameRequest{
			DymName: "not-exists",
		})
		require.NoError(t, err)
		require.Nil(t, resp.DymName)
	})

	now := time.Now().UTC()

	tests := []struct {
		name        string
		dymName     *dymnstypes.DymName
		queryName   string
		wantDymName *dymnstypes.DymName
	}{
		{
			name: "correct record",
			dymName: &dymnstypes.DymName{
				Name:       "bonded-pool",
				Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					},
				},
			},
			queryName: "bonded-pool",
			wantDymName: &dymnstypes.DymName{
				Name:       "bonded-pool",
				Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					},
				},
			},
		},
		{
			name: "NOT expired record only",
			dymName: &dymnstypes.DymName{
				Name:       "bonded-pool",
				Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				ExpireAt:   now.Unix() + 99,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					},
				},
			},
			queryName: "bonded-pool",
			wantDymName: &dymnstypes.DymName{
				Name:       "bonded-pool",
				Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				ExpireAt:   now.Unix() + 99,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					},
				},
			},
		},
		{
			name: "return nil for expired record",
			dymName: &dymnstypes.DymName{
				Name:       "bonded-pool",
				Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				ExpireAt:   now.Unix() - 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					},
				},
			},
			queryName:   "bonded-pool",
			wantDymName: nil,
		},
		{
			name:        "return nil if not found",
			dymName:     nil,
			queryName:   "non-exists",
			wantDymName: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, _, _, ctx := testkeeper.DymNSKeeper(t)

			ctx = ctx.WithBlockTime(now)

			if tt.dymName != nil {
				err := dk.SetDymName(ctx, *tt.dymName)
				require.NoError(t, err)
			}

			queryServer := dymnskeeper.NewQueryServerImpl(dk)
			resp, err := queryServer.DymName(sdk.WrapSDKContext(ctx), &dymnstypes.QueryDymNameRequest{
				DymName: tt.queryName,
			})
			require.NoError(t, err, "should never returns error")
			require.NotNil(t, resp, "should never returns nil response")

			if tt.wantDymName == nil {
				require.Nil(t, resp.DymName)
				return
			}

			require.NotNil(t, resp.DymName)
			require.Equal(t, tt.wantDymName, resp.DymName)
		})
	}

	t.Run("reject nil request", func(t *testing.T) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		queryServer := dymnskeeper.NewQueryServerImpl(dk)
		resp, err := queryServer.DymName(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_ResolveDymNameAddresses(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"

	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: now,
	}).WithChainID(chainId)

	addr1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	addr2 := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	addr3 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: addr1,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameA))

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: addr2,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameB))

	dymNameC := dymnstypes.DymName{
		Name:       "c",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: addr3,
		}},
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameC))

	dymNameD := dymnstypes.DymName{
		Name:       "d",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
		Configs: []dymnstypes.DymNameConfig{
			{
				Type:  dymnstypes.DymNameConfigType_NAME,
				Path:  "sub",
				Value: addr3,
			},
			{
				Type:    dymnstypes.DymNameConfigType_NAME,
				ChainId: "blumbus_111-1",
				Path:    "",
				Value:   addr3,
			},
		},
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameD))

	queryServer := dymnskeeper.NewQueryServerImpl(dk)

	resp, err := queryServer.ResolveDymNameAddresses(sdk.WrapSDKContext(ctx), &dymnstypes.QueryResolveDymNameAddressesRequest{
		Addresses: []string{
			"a.dymension_1100-1",
			"b.dymension_1100-1",
			"c.dymension_1100-1",
			"a.blumbus_111-1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.ResolvedAddresses, 4)

	require.Equal(t, addr1, resp.ResolvedAddresses[0].ResolvedAddress)
	require.Equal(t, addr2, resp.ResolvedAddresses[1].ResolvedAddress)
	require.Equal(t, addr3, resp.ResolvedAddresses[2].ResolvedAddress)
	require.Empty(t, resp.ResolvedAddresses[3].ResolvedAddress)
	require.NotEmpty(t, resp.ResolvedAddresses[3].Error)

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.ResolveDymNameAddresses(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject empty request", func(t *testing.T) {
		resp, err := queryServer.ResolveDymNameAddresses(
			sdk.WrapSDKContext(ctx),
			&dymnstypes.QueryResolveDymNameAddressesRequest{},
		)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("resolves default to owner if no config of default (without sub-name)", func(t *testing.T) {
		resp, err := queryServer.ResolveDymNameAddresses(
			sdk.WrapSDKContext(ctx),
			&dymnstypes.QueryResolveDymNameAddressesRequest{
				Addresses: []string{"d.dymension_1100-1", "d.blumbus_111-1"},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.ResolvedAddresses, 2)
		require.Equal(t, addr1, resp.ResolvedAddresses[0].ResolvedAddress)
		require.Equal(t, addr3, resp.ResolvedAddresses[1].ResolvedAddress)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_DymNamesOwnedByAccount(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"

	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: now,
	}).WithChainID(chainId)

	addr1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	addr2 := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	addr3 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: addr1,
		}},
	}
	setDymNameWithFunctionsAfter(ctx, dymNameA, t, dk)

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
	}
	setDymNameWithFunctionsAfter(ctx, dymNameB, t, dk)

	dymNameCExpired := dymnstypes.DymName{
		Name:       "c",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() - 1,
		Configs: []dymnstypes.DymNameConfig{{
			Type:  dymnstypes.DymNameConfigType_NAME,
			Value: addr3,
		}},
	}
	setDymNameWithFunctionsAfter(ctx, dymNameCExpired, t, dk)

	dymNameD := dymnstypes.DymName{
		Name:       "d",
		Owner:      addr3,
		Controller: addr3,
		ExpireAt:   now.Unix() + 1,
	}
	setDymNameWithFunctionsAfter(ctx, dymNameD, t, dk)

	queryServer := dymnskeeper.NewQueryServerImpl(dk)
	resp, err := queryServer.DymNamesOwnedByAccount(sdk.WrapSDKContext(ctx), &dymnstypes.QueryDymNamesOwnedByAccountRequest{
		Owner: addr1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.DymNames, 2)
	require.True(t, resp.DymNames[0].Name == dymNameA.Name || resp.DymNames[1].Name == dymNameA.Name)
	require.True(t, resp.DymNames[0].Name == dymNameB.Name || resp.DymNames[1].Name == dymNameB.Name)

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.DymNamesOwnedByAccount(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject invalid request", func(t *testing.T) {
		resp, err := queryServer.DymNamesOwnedByAccount(sdk.WrapSDKContext(ctx), &dymnstypes.QueryDymNamesOwnedByAccountRequest{
			Owner: "x",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_SellOrder(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"

	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: now,
	}).WithChainID(chainId)

	addr1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	addr2 := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameA))
	err := dk.SetSellOrder(ctx, dymnstypes.SellOrder{
		Name:     dymNameA.Name,
		ExpireAt: now.Unix() + 1,
		MinPrice: dymnsutils.TestCoin(100),
	})
	require.NoError(t, err)

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameB))

	queryServer := dymnskeeper.NewQueryServerImpl(dk)
	resp, err := queryServer.SellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QuerySellOrderRequest{
		DymName: dymNameA.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Result.Name == dymNameA.Name)

	t.Run("returns error code not found", func(t *testing.T) {
		resp, err := queryServer.SellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QuerySellOrderRequest{
			DymName: dymNameB.Name,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no active Sell Order")
		require.Nil(t, resp)
	})

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.SellOrder(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject invalid request", func(t *testing.T) {
		resp, err := queryServer.SellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QuerySellOrderRequest{
			DymName: "$$$",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_HistoricalSellOrder(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"

	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	ctx = ctx.WithBlockHeader(tmproto.Header{
		Time: now,
	}).WithChainID(chainId)

	addr1 := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	addr2 := "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	addr3 := "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 100,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameA))
	for r := int64(1); r <= 5; r++ {
		err := dk.SetSellOrder(ctx, dymnstypes.SellOrder{
			Name:      dymNameA.Name,
			ExpireAt:  now.Unix() + r,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(200),
			HighestBid: &dymnstypes.SellOrderBid{
				Bidder: addr3,
				Price:  dymnsutils.TestCoin(200),
			},
		})
		require.NoError(t, err)
		err = dk.MoveSellOrderToHistorical(ctx, dymNameA.Name)
		require.NoError(t, err)
	}

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 100,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameB))
	for r := int64(1); r <= 3; r++ {
		err := dk.SetSellOrder(ctx, dymnstypes.SellOrder{
			Name:      dymNameB.Name,
			ExpireAt:  now.Unix() + r,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
			HighestBid: &dymnstypes.SellOrderBid{
				Bidder: addr3,
				Price:  dymnsutils.TestCoin(300),
			},
		})
		require.NoError(t, err)
		err = dk.MoveSellOrderToHistorical(ctx, dymNameB.Name)
		require.NoError(t, err)
	}

	queryServer := dymnskeeper.NewQueryServerImpl(dk)
	resp, err := queryServer.HistoricalSellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalSellOrderRequest{
		DymName: dymNameA.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Result, 5)

	resp, err = queryServer.HistoricalSellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalSellOrderRequest{
		DymName: dymNameB.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Result, 3)

	t.Run("returns empty for non-exists Dym-Name", func(t *testing.T) {
		resp, err := queryServer.HistoricalSellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalSellOrderRequest{
			DymName: "not-exists",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Empty(t, resp.Result)
	})

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.HistoricalSellOrder(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject invalid request", func(t *testing.T) {
		resp, err := queryServer.HistoricalSellOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalSellOrderRequest{
			DymName: "$$$",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_EstimateRegisterName(t *testing.T) {
	now := time.Now()

	const denom = "atom"
	const price1L int64 = 9
	const price2L int64 = 8
	const price3L int64 = 7
	const price4L int64 = 6
	const price5PlusL int64 = 5
	const extendsPrice int64 = 4

	setupTest := func() (dymnskeeper.Keeper, sdk.Context) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		params := dymnstypes.DefaultParams()
		params.Price.PriceDenom = denom
		params.Price.Price_1Letter = sdk.NewInt(price1L)
		params.Price.Price_2Letters = sdk.NewInt(price2L)
		params.Price.Price_3Letters = sdk.NewInt(price3L)
		params.Price.Price_4Letters = sdk.NewInt(price4L)
		params.Price.Price_5PlusLetters = sdk.NewInt(price5PlusL)
		params.Price.PriceExtends = sdk.NewInt(extendsPrice)
		params.Misc.DaysGracePeriod = 1
		err := dk.SetParams(ctx, params)
		require.NoError(t, err)

		return dk, ctx
	}

	const buyer = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const previousOwner = "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"

	tests := []struct {
		name               string
		dymName            string
		existingDymName    *dymnstypes.DymName
		newOwner           string
		duration           int64
		wantErr            bool
		wantErrContains    string
		wantFirstYearPrice int64
		wantExtendPrice    int64
	}{
		{
			name:               "new registration, 1 letter, 1 year",
			dymName:            "a",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    0,
		},
		{
			name:               "new registration, empty buyer",
			dymName:            "a",
			existingDymName:    nil,
			newOwner:           "",
			duration:           1,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    0,
		},
		{
			name:               "new registration, 1 letter, 2 years",
			dymName:            "a",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           2,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    extendsPrice,
		},
		{
			name:               "new registration, 1 letter, N years",
			dymName:            "a",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           99,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    extendsPrice * (99 - 1),
		},
		{
			name:               "new registration, 6 letters, 1 year",
			dymName:            "abcdef",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    0,
		},
		{
			name:               "new registration, 6 letters, 2 years",
			dymName:            "abcdef",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           2,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    extendsPrice,
		},
		{
			name:               "new registration, 5+ letters, N years",
			dymName:            "abcdef",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           99,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    extendsPrice * (99 - 1),
		},
		{
			name:    "extends same owner, 1 letter, 1 year",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice,
		},
		{
			name:    "extends same owner, 1 letter, 2 years",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           2,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:    "extends same owner, 1 letter, N years",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           99,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice * 99,
		},
		{
			name:    "extends same owner, 6 letters, 1 year",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice,
		},
		{
			name:    "extends same owner, 6 letters, 2 years",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           2,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:    "extends same owner, 5+ letters, N years",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:           buyer,
			duration:           99,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice * 99,
		},
		{
			name:    "extends expired, same owner, 5+ letters, 2 years",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           buyer,
			duration:           2,
			wantFirstYearPrice: 0,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:    "extends expired, empty buyer, treat as take over",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      buyer,
				Controller: buyer,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           "",
			duration:           2,
			wantFirstYearPrice: 5,
			wantExtendPrice:    extendsPrice,
		},
		{
			name:    "take-over, 1 letter, 1 year",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    0,
		},
		{
			name:    "take-over, 1 letter, 3 years",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:    "take-over, 6 letters, 1 year",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           buyer,
			duration:           1,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    0,
		},
		{
			name:    "take-over, 6 letters, 3 years",
			dymName: "abcdef",
			existingDymName: &dymnstypes.DymName{
				Name:       "abcdef",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1,
			},
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:               "new registration, 2 letters",
			dymName:            "aa",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price2L,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:               "new registration, 3 letters",
			dymName:            "aaa",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price3L,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:               "new registration, 4 letters",
			dymName:            "aaaa",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price4L,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:               "new registration, 5 letters",
			dymName:            "aaaaa",
			existingDymName:    nil,
			newOwner:           buyer,
			duration:           3,
			wantFirstYearPrice: price5PlusL,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:            "reject invalid Dym-Name",
			dymName:         "-a-",
			existingDymName: nil,
			newOwner:        buyer,
			duration:        2,
			wantErr:         true,
			wantErrContains: "invalid dym name",
		},
		{
			name:            "reject invalid duration",
			dymName:         "a",
			existingDymName: nil,
			newOwner:        buyer,
			duration:        0,
			wantErr:         true,
			wantErrContains: "duration must be at least 1 year",
		},
		{
			name:    "reject estimation for Dym-Name owned by another and not expired",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:        buyer,
			duration:        1,
			wantErr:         true,
			wantErrContains: "you are not the owner",
		},
		{
			name:    "reject estimation for Dym-Name owned by another and not expired, empty buyer",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() + 1,
			},
			newOwner:        "",
			duration:        1,
			wantErr:         true,
			wantErrContains: "you are not the owner",
		},
		{
			name:    "allow estimation for take-over, regardless grace period",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1, // still in grace period
			},
			newOwner:           buyer,
			duration:           3,
			wantErr:            false,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    extendsPrice * 2,
		},
		{
			name:    "allow estimation for take-over, regardless grace period, empty buyer",
			dymName: "a",
			existingDymName: &dymnstypes.DymName{
				Name:       "a",
				Owner:      previousOwner,
				Controller: previousOwner,
				ExpireAt:   now.Unix() - 1, // still in grace period
			},
			newOwner:           "",
			duration:           3,
			wantErr:            false,
			wantFirstYearPrice: price1L,
			wantExtendPrice:    extendsPrice * 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, ctx := setupTest()

			require.Positive(t, dk.MiscParams(ctx).DaysGracePeriod, "bad setup, must have grace period")

			if tt.existingDymName != nil {
				err := dk.SetDymName(ctx, *tt.existingDymName)
				require.NoError(t, err)
			}

			queryServer := dymnskeeper.NewQueryServerImpl(dk)

			resp, err := queryServer.EstimateRegisterName(sdk.WrapSDKContext(ctx), &dymnstypes.QueryEstimateRegisterNameRequest{
				Name:     tt.dymName,
				Duration: tt.duration,
				Owner:    tt.newOwner,
			})

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			require.Equal(t, tt.wantFirstYearPrice, resp.FirstYearPrice.Amount.Int64())
			require.Equal(t, tt.wantExtendPrice, resp.ExtendPrice.Amount.Int64())
			require.Equal(
				t,
				tt.wantFirstYearPrice+tt.wantExtendPrice,
				resp.TotalPrice.Amount.Int64(),
				"total price must be equals to sum of first year and extend price",
			)
			require.Equal(t, denom, resp.FirstYearPrice.Denom)
			require.Equal(t, denom, resp.ExtendPrice.Denom)
			require.Equal(t, denom, resp.TotalPrice.Denom)
		})
	}
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_ReverseResolveAddress(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"
	const nimChainId = "nim_1122-1"

	setupTest := func() (dymnskeeper.Keeper, sdk.Context) {
		dk, _, rk, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		}).WithChainID(chainId)

		moduleParams := dk.GetParams(ctx)
		moduleParams.Chains.AliasesByChainId = map[string]dymnstypes.AliasesOfChainId{
			chainId: {
				Aliases: []string{"dym"},
			},
			nimChainId: {
				Aliases: []string{"nim"},
			},
		}
		err := dk.SetParams(ctx, moduleParams)
		require.NoError(t, err)

		// add rollapp to enable 0x address reverse mapping for this chain
		rk.SetRollapp(ctx, rollapptypes.Rollapp{
			RollappId: nimChainId,
			Creator:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		})

		return dk, ctx
	}

	t.Run("reject nil request", func(t *testing.T) {
		dk, ctx := setupTest()
		queryServer := dymnskeeper.NewQueryServerImpl(dk)

		resp, err := queryServer.ReverseResolveAddress(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject empty request", func(t *testing.T) {
		dk, ctx := setupTest()
		queryServer := dymnskeeper.NewQueryServerImpl(dk)

		resp, err := queryServer.ReverseResolveAddress(sdk.WrapSDKContext(ctx), &dymnstypes.QueryReverseResolveAddressRequest{
			Addresses: []string{},
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	const owner = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	owner0x := common.BytesToAddress(sdk.MustAccAddressFromBech32(owner)).Hex()
	const anotherAcc = "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d"
	anotherAcc0x := common.BytesToAddress(sdk.MustAccAddressFromBech32(anotherAcc)).Hex()
	const ica = "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"
	ica0x := common.BytesToHash(sdk.MustAccAddressFromBech32(ica)).Hex()

	const cosmosAcc = "cosmos18wvvwfmq77a6d8tza4h5sfuy2yj3jj88yqg82a"
	_, cosmosAcc0xBz, err := bech32.DecodeAndConvert(cosmosAcc)
	require.NoError(t, err)
	cosmosAcc0x := common.BytesToAddress(cosmosAcc0xBz).Hex()

	tests := []struct {
		name            string
		dymNames        []dymnstypes.DymName
		addresses       []string
		wantErr         bool
		wantErrContains string
		wantResult      map[string]dymnstypes.ReverseResolveAddressResult
	}{
		{
			name: "pass - mixed addresses type",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
				},
			},
			addresses: []string{owner, owner0x},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				owner: {
					Candidates: []string{"a@dym"},
				},
				owner0x: {
					Candidates: []string{"a@dym"},
				},
			},
		},
		{
			name: "pass - multiple addresses",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "another.account",
							Value:   anotherAcc,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   cosmosAcc,
						},
					},
				},
			},
			addresses: []string{
				owner,
				anotherAcc,
				cosmosAcc,
			},
			wantErr: false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				owner: {
					Candidates: []string{"a@dym"},
				},
				anotherAcc: {
					Candidates: []string{"another.account.a@dym"},
				},
				cosmosAcc: {
					Candidates: []string{"a@cosmoshub-4"},
				},
			},
		},
		{
			name: "pass - multi-level sub-name",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "a.b.c.d",
							Value:   owner,
						},
					},
				},
			},
			addresses: []string{owner},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				owner: {
					Candidates: []string{"a@dym", "a.b.c.d.a@dym"},
				},
			},
		},
		{
			name: "pass - each address match multiple result",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "a.b.c.d",
							Value:   owner,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "another",
							Value:   anotherAcc,
						},
					},
				},
				{
					Name:       "b",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "e.f.g.h",
							Value:   owner,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "another",
							Value:   anotherAcc,
						},
					},
				},
			},
			addresses: []string{owner, anotherAcc0x},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				owner: {
					Candidates: []string{"a@dym", "b@dym", "a.b.c.d.a@dym", "e.f.g.h.b@dym"},
				},
				anotherAcc0x: {
					Candidates: []string{"another.a@dym", "another.b@dym"},
				},
			},
		},
		{
			name: "pass - alias not mapped if no alias",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   cosmosAcc,
						},
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: nimChainId,
							Path:    "",
							Value:   owner,
						},
					},
				},
			},
			addresses: []string{cosmosAcc, owner},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				cosmosAcc: {
					Candidates: []string{"a@cosmoshub-4"},
				},
				owner: {
					Candidates: []string{"a@dym", "a@nim"},
				},
			},
		},
		{
			name: "pass - support ICA address",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "",
							Path:    "ica",
							Value:   ica,
						},
					},
				},
			},
			addresses: []string{ica, ica0x},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				ica: {
					Candidates: []string{"ica.a@dym"},
				},
				ica0x: {
					Candidates: []string{"ica.a@dym"},
				},
			},
		},
		{
			name: "pass - chains not coin-type-60 should not support reverse-resolve 0x address",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
					Configs: []dymnstypes.DymNameConfig{
						{
							Type:    dymnstypes.DymNameConfigType_NAME,
							ChainId: "cosmoshub-4",
							Path:    "",
							Value:   cosmosAcc,
						},
					},
				},
			},
			addresses: []string{cosmosAcc, cosmosAcc0x},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				cosmosAcc: {
					Candidates: []string{"a@cosmoshub-4"},
				},
				cosmosAcc0x: {
					Candidates: []string{},
				},
			},
		},
		{
			name: "pass - returns empty for non-reverse-resolvable address",
			dymNames: []dymnstypes.DymName{
				{
					Name:       "a",
					Owner:      owner,
					Controller: owner,
					ExpireAt:   now.Unix() + 1,
				},
			},
			addresses: []string{anotherAcc, anotherAcc0x},
			wantErr:   false,
			wantResult: map[string]dymnstypes.ReverseResolveAddressResult{
				anotherAcc: {
					Candidates: []string{},
				},
				anotherAcc0x: {
					Candidates: []string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, ctx := setupTest()

			for _, dymName := range tt.dymNames {
				setDymNameWithFunctionsAfter(ctx, dymName, t, dk)
			}

			queryServer := dymnskeeper.NewQueryServerImpl(dk)

			resp, err := queryServer.ReverseResolveAddress(sdk.WrapSDKContext(ctx), &dymnstypes.QueryReverseResolveAddressRequest{
				Addresses: tt.addresses,
			})

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			if !reflect.DeepEqual(tt.wantResult, resp.Result) {
				t.Errorf("got = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}