package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, dk.SetDymName(ctx, dymNameA))

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      addr1,
		Controller: addr2,
		ExpireAt:   now.Unix() + 1,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameB))

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
	require.NoError(t, dk.SetDymName(ctx, dymNameCExpired))

	dymNameD := dymnstypes.DymName{
		Name:       "d",
		Owner:      addr3,
		Controller: addr3,
		ExpireAt:   now.Unix() + 1,
	}
	require.NoError(t, dk.SetDymName(ctx, dymNameD))

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
func Test_queryServer_OpenPurchaseOrder(t *testing.T) {
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
	err := dk.SetOpenPurchaseOrder(ctx, dymnstypes.OpenPurchaseOrder{
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
	resp, err := queryServer.OpenPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryOpenPurchaseOrderRequest{
		DymName: dymNameA.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Result.Name == dymNameA.Name)

	t.Run("returns error code not found", func(t *testing.T) {
		resp, err := queryServer.OpenPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryOpenPurchaseOrderRequest{
			DymName: dymNameB.Name,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no active Open Purchase Order")
		require.Nil(t, resp)
	})

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.OpenPurchaseOrder(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject invalid request", func(t *testing.T) {
		resp, err := queryServer.OpenPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryOpenPurchaseOrderRequest{
			DymName: "$$$",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

//goland:noinspection SpellCheckingInspection
func Test_queryServer_HistoricalPurchaseOrder(t *testing.T) {
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
		err := dk.SetOpenPurchaseOrder(ctx, dymnstypes.OpenPurchaseOrder{
			Name:      dymNameA.Name,
			ExpireAt:  now.Unix() + r,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(200),
			HighestBid: &dymnstypes.OpenPurchaseOrderBid{
				Bidder: addr3,
				Price:  dymnsutils.TestCoin(200),
			},
		})
		require.NoError(t, err)
		err = dk.MoveOpenPurchaseOrderToHistorical(ctx, dymNameA.Name)
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
		err := dk.SetOpenPurchaseOrder(ctx, dymnstypes.OpenPurchaseOrder{
			Name:      dymNameB.Name,
			ExpireAt:  now.Unix() + r,
			MinPrice:  dymnsutils.TestCoin(100),
			SellPrice: dymnsutils.TestCoinP(300),
			HighestBid: &dymnstypes.OpenPurchaseOrderBid{
				Bidder: addr3,
				Price:  dymnsutils.TestCoin(300),
			},
		})
		require.NoError(t, err)
		err = dk.MoveOpenPurchaseOrderToHistorical(ctx, dymNameB.Name)
		require.NoError(t, err)
	}

	queryServer := dymnskeeper.NewQueryServerImpl(dk)
	resp, err := queryServer.HistoricalPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalPurchaseOrderRequest{
		DymName: dymNameA.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Result, 5)

	resp, err = queryServer.HistoricalPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalPurchaseOrderRequest{
		DymName: dymNameB.Name,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Result, 3)

	t.Run("returns empty for non-exists Dym-Name", func(t *testing.T) {
		resp, err := queryServer.HistoricalPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalPurchaseOrderRequest{
			DymName: "not-exists",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Empty(t, resp.Result)
	})

	t.Run("reject nil request", func(t *testing.T) {
		resp, err := queryServer.HistoricalPurchaseOrder(sdk.WrapSDKContext(ctx), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("reject invalid request", func(t *testing.T) {
		resp, err := queryServer.HistoricalPurchaseOrder(sdk.WrapSDKContext(ctx), &dymnstypes.QueryHistoricalPurchaseOrderRequest{
			DymName: "$$$",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
