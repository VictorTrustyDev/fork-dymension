package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func Test_epochHooks_BeforeEpochStart(t *testing.T) {
	now := time.Now().UTC()
	const daysKeepHistorical = 1
	require.Greater(t, daysKeepHistorical, 0, "mis-configured test case")

	setupTest := func() (dymnskeeper.Keeper, sdk.Context) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)

		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		params := dk.GetParams(ctx)
		params.Misc.DaysPreservedClosedPurchaseOrder = daysKeepHistorical
		err := dk.SetParams(ctx, params)
		require.NoError(t, err)

		return dk, ctx
	}

	t.Run("should do something even nothing to do", func(t *testing.T) {
		dk, ctx := setupTest()

		params := dk.GetParams(ctx)

		originalGas := ctx.GasMeter().GasConsumed()

		err := dk.GetEpochHooks().BeforeEpochStart(ctx, params.Misc.BeginEpochHookIdentifier, 1)
		require.NoError(t, err)

		// gas should be changed because it should at least reading the params to check epoch identifier
		require.Less(t, originalGas, ctx.GasMeter().GasConsumed(), "should do something")
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Add(365 * 24 * time.Hour).Unix(),
	}

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Unix(),
	}

	dymNameC := dymnstypes.DymName{
		Name:       "c",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Add(-365 * 24 * time.Hour).Unix(),
	}

	dymNameD := dymnstypes.DymName{
		Name:       "d",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   1,
	}

	getEpochWithOffset := func(offset int64) int64 {
		return now.Unix() + offset
	}
	genOpo := func(
		dymName dymnstypes.DymName, offsetExpiry int64,
	) dymnstypes.OpenPurchaseOrder {
		return dymnstypes.OpenPurchaseOrder{
			Name:     dymName.Name,
			ExpireAt: getEpochWithOffset(offsetExpiry),
			MinPrice: dymnsutils.TestCoin(100),
		}
	}

	type testSuite struct {
		t   *testing.T
		dk  dymnskeeper.Keeper
		ctx sdk.Context
	}

	nts := func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) testSuite {
		return testSuite{
			t:   t,
			dk:  dk,
			ctx: ctx,
		}
	}

	requireDymNameNotChanged := func(dymName dymnstypes.DymName, ts testSuite) {
		laterDymName := ts.dk.GetDymName(ts.ctx, dymName.Name)
		require.NotNil(t, laterDymName)

		require.Equal(t, dymName, *laterDymName, "nothing changed")
	}

	requireNoActiveOPO := func(dymName dymnstypes.DymName, ts testSuite) {
		opo := ts.dk.GetOpenPurchaseOrder(ts.ctx, dymName.Name)
		require.Nil(t, opo)
	}

	requireActiveOPO := func(dymName dymnstypes.DymName, ts testSuite) {
		opo := ts.dk.GetOpenPurchaseOrder(ts.ctx, dymName.Name)
		require.NotNil(t, opo)
	}

	requireHistoricalOPOs := func(dymName dymnstypes.DymName, wantCount int, ts testSuite) {
		historicalPOs := ts.dk.GetHistoricalOpenPurchaseOrders(ts.ctx, dymName.Name)
		require.Lenf(t, historicalPOs, wantCount, "should have %d historical OPOs", wantCount)
	}

	tests := []struct {
		name                   string
		dymNames               []dymnstypes.DymName
		historicalOPOs         []dymnstypes.OpenPurchaseOrder
		activeOPOs             []dymnstypes.OpenPurchaseOrder
		minExpiryByDymName     map[string]int64
		customEpochIdentifier  string
		wantErr                bool
		wantErrContains        string
		wantMinExpiryByDymName map[string]int64
		preHookTestFunc        func(*testing.T, dymnskeeper.Keeper, sdk.Context)
		afterHookTestFunc      func(*testing.T, dymnskeeper.Keeper, sdk.Context)
	}{
		{
			name:     "simple cleanup",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -daysKeepHistorical*86400-1),
			},
			activeOPOs: nil,
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-daysKeepHistorical*86400 - 1),
			},
			wantErr:                false,
			wantMinExpiryByDymName: nil,
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 0, ts)

				requireDymNameNotChanged(dymNameB, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)

				requireDymNameNotChanged(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)

				requireDymNameNotChanged(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
		},
		{
			name:     "mis-match epoch will clean nothing",
			dymNames: []dymnstypes.DymName{dymNameA},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -daysKeepHistorical*86400-1),
			},
			activeOPOs: nil,
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-daysKeepHistorical*86400 - 1),
			},
			customEpochIdentifier: "not-match",
			wantErr:               false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-daysKeepHistorical*86400 - 1),
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
			},
		},
		{
			name:     "simple cleanup, with active OPO",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -daysKeepHistorical*86400-1),
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-daysKeepHistorical*86400 - 1),
			},
			wantErr:                false,
			wantMinExpiryByDymName: nil,
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 0, ts)
				requireActiveOPO(dymNameA, ts)

				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
		},
		{
			name:           "simple cleanup, no historical record to prune",
			dymNames:       []dymnstypes.DymName{dymNameA},
			historicalOPOs: nil,
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1),
			},
			minExpiryByDymName:     nil,
			wantErr:                false,
			wantMinExpiryByDymName: nil,
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 0, ts)
				requireActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 0, ts)
				requireActiveOPO(dymNameA, ts)
			},
		},
		{
			name:     "simple cleanup, nothing to prune",
			dymNames: []dymnstypes.DymName{dymNameA},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -1),
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-1),
			},
			wantErr: false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-1),
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)
			},
		},
		{
			name:     "cleanup multiple Historical OPO, all need to prune",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -(daysKeepHistorical+0)*86400-1),
				genOpo(dymNameA, -(daysKeepHistorical+2)*86400-1),
				genOpo(dymNameA, -(daysKeepHistorical+1)*86400-1),
				genOpo(dymNameC, -(daysKeepHistorical+3)*86400-1),
				genOpo(dymNameC, -(daysKeepHistorical+5)*86400-1),
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameC, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-(daysKeepHistorical+2)*86400 - 1),
				dymNameC.Name: getEpochWithOffset(-(daysKeepHistorical+5)*86400 - 1),
			},
			wantErr:                false,
			wantMinExpiryByDymName: nil,
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 3, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 2, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 0, ts)
				requireNoActiveOPO(dymNameA, ts)

				requireDymNameNotChanged(dymNameB, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireNoActiveOPO(dymNameB, ts)

				requireDymNameNotChanged(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireActiveOPO(dymNameC, ts)

				requireDymNameNotChanged(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
				requireNoActiveOPO(dymNameD, ts)
			},
		},
		{
			name:     "cleanup multiple Historical OPO, some need to prune while some not",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -(daysKeepHistorical+0)*86400-1),
				genOpo(dymNameA, -(daysKeepHistorical+2)*86400-1),
				genOpo(dymNameA, -9),
				genOpo(dymNameC, -(daysKeepHistorical+3)*86400-1),
				genOpo(dymNameC, -(daysKeepHistorical+5)*86400-1),
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1), genOpo(dymNameB, +1), genOpo(dymNameC, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-(daysKeepHistorical+2)*86400 - 1),
				dymNameC.Name: getEpochWithOffset(-(daysKeepHistorical+5)*86400 - 1),
			},
			wantErr: false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-9),
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 3, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireHistoricalOPOs(dymNameC, 2, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)

				requireDymNameNotChanged(dymNameB, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireActiveOPO(dymNameB, ts)

				requireDymNameNotChanged(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireActiveOPO(dymNameC, ts)

				requireDymNameNotChanged(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
				requireNoActiveOPO(dymNameD, ts)
			},
		},
		{
			name:     "should update min expiry correctly",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, -9),
				genOpo(dymNameA, -(daysKeepHistorical+2)*86400-1),
				genOpo(dymNameA, -10),
			},
			activeOPOs: nil,
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-(daysKeepHistorical+2)*86400 - 1),
			},
			wantErr: false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-10),
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 3, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 2, ts)
			},
		},
		{
			name:     "mixed cleanup",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				// Dym-Name A has some historical OPO, some need to prune, some not
				genOpo(dymNameA, -(daysKeepHistorical+0)*86400-1),
				genOpo(dymNameA, -(daysKeepHistorical+2)*86400-1),
				genOpo(dymNameA, -9),
				// Dym-Name B has some historical OPO, no need to prune
				genOpo(dymNameB, -8),
				genOpo(dymNameB, -7),
				// Dym-Name C has some historical OPO, all need to prune
				genOpo(dymNameC, -(daysKeepHistorical+3)*86400-1),
				genOpo(dymNameC, -(daysKeepHistorical+5)*86400-1),
				// Dym-Name D has no historical OPO
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1), genOpo(dymNameB, +1), genOpo(dymNameC, +1), genOpo(dymNameD, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-(daysKeepHistorical+2)*86400 - 1),
				dymNameB.Name: getEpochWithOffset(-8),
				dymNameC.Name: getEpochWithOffset(-(daysKeepHistorical+5)*86400 - 1),
			},
			wantErr: false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-9),
				dymNameB.Name: getEpochWithOffset(-8),
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 3, ts)
				requireHistoricalOPOs(dymNameB, 2, ts)
				requireHistoricalOPOs(dymNameC, 2, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)

				requireDymNameNotChanged(dymNameB, ts)
				requireHistoricalOPOs(dymNameB, 2, ts)
				requireActiveOPO(dymNameB, ts)

				requireDymNameNotChanged(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireActiveOPO(dymNameC, ts)

				requireDymNameNotChanged(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
				requireActiveOPO(dymNameD, ts)
			},
		},
		{
			name:           "case no historical OPO but has min expiry",
			dymNames:       []dymnstypes.DymName{dymNameA},
			historicalOPOs: nil,
			activeOPOs:     nil,
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: 1, // incorrect state: no historical OPO but has min expiry
			},
			wantErr:                false,
			wantMinExpiryByDymName: nil,
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 0, ts)
			},
		},
		{
			name:     "mixed cleanup with incorrect state",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			historicalOPOs: []dymnstypes.OpenPurchaseOrder{
				// Dym-Name A has some OPO, no need to prune
				genOpo(dymNameA, -9),
				// Dym-Name D has no historical OPO
			},
			activeOPOs: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, +1), genOpo(dymNameB, +1), genOpo(dymNameC, +1), genOpo(dymNameD, +1),
			},
			minExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-daysKeepHistorical*86400 - 1), // incorrect state: has historical OPO, no need to prune but min-expiry indicates need to prune
				dymNameD.Name: 1,                                                 // incorrect state: no historical OPO but has min expiry
			},
			wantErr: false,
			wantMinExpiryByDymName: map[string]int64{
				dymNameA.Name: getEpochWithOffset(-9), // corrected value
				// incorrect of Dym-Name D was removed
			},
			preHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireHistoricalOPOs(dymNameA, 1, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, ctx sdk.Context) {
				ts := nts(t, dk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)
				requireActiveOPO(dymNameA, ts)

				requireDymNameNotChanged(dymNameB, ts)
				requireHistoricalOPOs(dymNameB, 0, ts)
				requireActiveOPO(dymNameB, ts)

				requireDymNameNotChanged(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 0, ts)
				requireActiveOPO(dymNameC, ts)

				requireDymNameNotChanged(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 0, ts)
				requireActiveOPO(dymNameD, ts)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.preHookTestFunc, "mis-configured test case")
			require.NotNil(t, tt.afterHookTestFunc, "mis-configured test case")

			dk, ctx := setupTest()

			for _, dymName := range tt.dymNames {
				err := dk.SetDymName(ctx, dymName)
				require.NoError(t, err)
			}

			for _, opo := range tt.historicalOPOs {
				err := dk.SetOpenPurchaseOrder(ctx, opo)
				require.NoError(t, err)
				err = dk.MoveOpenPurchaseOrderToHistorical(ctx, opo.Name)
				require.NoError(t, err)
			}

			for _, opo := range tt.activeOPOs {
				err := dk.SetOpenPurchaseOrder(ctx, opo)
				require.NoError(t, err)
			}

			meh := dk.GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx)
			if len(meh) > 0 {
				// clear existing records to simulate cases of malformed state
				for dymName := range meh {
					dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, 0)
				}
			}
			if len(tt.minExpiryByDymName) > 0 {
				for dymName, minExpiry := range tt.minExpiryByDymName {
					dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, minExpiry)
				}
			}

			tt.preHookTestFunc(t, dk, ctx)

			moduleParams := dk.GetParams(ctx)
			useEpochIdentifier := moduleParams.Misc.BeginEpochHookIdentifier
			if tt.customEpochIdentifier != "" {
				useEpochIdentifier = tt.customEpochIdentifier
			}
			err := dk.GetEpochHooks().BeforeEpochStart(ctx, useEpochIdentifier, 1)

			defer func() {
				if t.Failed() {
					return
				}

				tt.afterHookTestFunc(t, dk, ctx)

				meh := dk.GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx)
				if len(tt.wantMinExpiryByDymName) == 0 {
					require.Empty(t, meh)
				} else if !reflect.DeepEqual(tt.wantMinExpiryByDymName, meh) {
					t.Errorf("want map %v, got %v", tt.wantMinExpiryByDymName, meh)
				}
			}()

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				return
			}

			require.NoError(t, err)
		})
	}
}

//goland:noinspection SpellCheckingInspection
func Test_epochHooks_AfterEpochEnd(t *testing.T) {
	now := time.Now().UTC()

	setupTest := func() (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
		dk, bk, _, ctx := testkeeper.DymNSKeeper(t)

		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})

		return dk, bk, ctx
	}

	t.Run("should do something even nothing to do", func(t *testing.T) {
		dk, _, ctx := setupTest()

		params := dk.GetParams(ctx)

		originalGas := ctx.GasMeter().GasConsumed()

		err := dk.GetEpochHooks().AfterEpochEnd(ctx, params.Misc.EndEpochHookIdentifier, 1)
		require.NoError(t, err)

		// gas should be changed because it should at least reading the params to check epoch identifier
		require.Less(t, originalGas, ctx.GasMeter().GasConsumed(), "should do something")
	})

	owner := "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	bidder := "dym1ysjlrjcankjpmpxxzk27mvzhv25e266r80p5pv"
	dymNsModuleAccAddr := authtypes.NewModuleAddress(dymnstypes.ModuleName)

	dymNameA := dymnstypes.DymName{
		Name:       "a",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Unix() + 1,
	}

	dymNameB := dymnstypes.DymName{
		Name:       "b",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Unix() + 1,
	}

	dymNameC := dymnstypes.DymName{
		Name:       "c",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Unix() + 1,
	}

	dymNameD := dymnstypes.DymName{
		Name:       "d",
		Owner:      owner,
		Controller: owner,
		ExpireAt:   now.Unix() + 1,
	}

	coin100 := dymnsutils.TestCoin(100)
	coin200 := dymnsutils.TestCoin(200)
	denom := dymnsutils.TestCoin(0).Denom

	opoExpiredEpoch := now.Unix() - 1
	opoNotExpiredEpoch := now.Unix() + 1

	const opoExpired = true
	const opoNotExpired = false
	genOpo := func(
		dymName dymnstypes.DymName,
		expired bool, sellPrice *sdk.Coin, highestBid *dymnstypes.OpenPurchaseOrderBid,
	) dymnstypes.OpenPurchaseOrder {
		return dymnstypes.OpenPurchaseOrder{
			Name: dymName.Name,
			ExpireAt: func() int64 {
				if expired {
					return opoExpiredEpoch
				}
				return opoNotExpiredEpoch
			}(),
			MinPrice:   coin100,
			SellPrice:  sellPrice,
			HighestBid: highestBid,
		}
	}

	type testSuite struct {
		t   *testing.T
		dk  dymnskeeper.Keeper
		bk  dymnskeeper.BankKeeper
		ctx sdk.Context
	}

	nts := func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) testSuite {
		return testSuite{
			t:   t,
			dk:  dk,
			bk:  bk,
			ctx: ctx,
		}
	}

	requireOwnerChanged := func(dymName dymnstypes.DymName, newOwner string, ts testSuite) {
		require.NotEmpty(t, newOwner, "mis-configured test case")

		laterDymName := ts.dk.GetDymName(ts.ctx, dymName.Name)
		require.NotNil(t, laterDymName)

		require.Equal(t, newOwner, laterDymName.Owner, "ownership must be transferred")
		require.Equal(t, newOwner, laterDymName.Controller, "controller must be changed")
		require.Equal(t, dymName.ExpireAt, laterDymName.ExpireAt, "expiry must not be changed")
		require.Empty(t, laterDymName.Configs, "configs must be cleared")
	}

	requireDymNameNotChanged := func(dymName dymnstypes.DymName, ts testSuite) {
		laterDymName := ts.dk.GetDymName(ts.ctx, dymName.Name)
		require.NotNil(t, laterDymName)

		require.Equal(t, dymName, *laterDymName, "nothing changed")
	}

	requireNoActiveOPO := func(dymName dymnstypes.DymName, ts testSuite) {
		opo := ts.dk.GetOpenPurchaseOrder(ts.ctx, dymName.Name)
		require.Nil(t, opo)
	}

	requireHistoricalOPOs := func(dymName dymnstypes.DymName, wantCount int, ts testSuite) {
		historicalPOs := ts.dk.GetHistoricalOpenPurchaseOrders(ts.ctx, dymName.Name)
		require.Lenf(t, historicalPOs, wantCount, "should have %d historical OPOs", wantCount)
	}

	requireModuleBalance := func(wantAmount int64, ts testSuite) {
		moduleBalance := ts.bk.GetBalance(ts.ctx, dymNsModuleAccAddr, denom)
		require.NotNil(t, moduleBalance)

		require.Equalf(t, wantAmount, moduleBalance.Amount.Int64(), "module balance should be %d", wantAmount)
	}

	requireAccountBalance := func(bech32Addr string, wantAmount int64, ts testSuite) {
		accountBalance := ts.bk.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bech32Addr), denom)
		require.NotNil(t, accountBalance)

		require.Equalf(t, wantAmount, accountBalance.Amount.Int64(), "account balance should be %d", wantAmount)
	}

	tests := []struct {
		name                   string
		dymNames               []dymnstypes.DymName
		opos                   []dymnstypes.OpenPurchaseOrder
		mapExpiryByDymName     map[string]int64
		preMintModuleBalance   int64
		customEpochIdentifier  string
		wantErr                bool
		wantErrContains        string
		wantMapExpiryByDymName map[string]int64
		afterHookTestFunc      func(*testing.T, dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context)
	}{
		{
			name:     "simple process expired OPO",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			opos:     []dymnstypes.OpenPurchaseOrder{genOpo(dymNameA, opoExpired, &coin200, nil)},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
			},
			preMintModuleBalance:   200,
			wantErr:                false,
			wantMapExpiryByDymName: nil,
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				ts := nts(t, dk, bk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireNoActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)

				requireModuleBalance(200, ts)

				requireAccountBalance(dymNameA.Owner, 0, ts)
			},
		},
		{
			name:     "simple process expired & completed OPO",
			dymNames: []dymnstypes.DymName{dymNameA},
			opos: []dymnstypes.OpenPurchaseOrder{genOpo(dymNameA, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
				Bidder: bidder,
				Price:  coin200,
			})},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
			},
			preMintModuleBalance:   200,
			wantErr:                false,
			wantMapExpiryByDymName: nil,
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				ts := nts(t, dk, bk, ctx)

				requireOwnerChanged(dymNameA, bidder, ts)
				requireNoActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)

				requireModuleBalance(0, ts) // should transfered to previous owner

				requireAccountBalance(dymNameA.Owner, 200, ts) // previous owner should earn from bid
			},
		},
		{
			name:     "simple process expired & completed OPO, match by min price",
			dymNames: []dymnstypes.DymName{dymNameA},
			opos: []dymnstypes.OpenPurchaseOrder{genOpo(dymNameA, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
				Bidder: bidder,
				Price:  coin100,
			})},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
			},
			preMintModuleBalance:   250,
			wantErr:                false,
			wantMapExpiryByDymName: nil,
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				ts := nts(t, dk, bk, ctx)

				requireOwnerChanged(dymNameA, bidder, ts)
				requireNoActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)

				requireModuleBalance(150, ts) // 100 should transfered to previous owner

				requireAccountBalance(dymNameA.Owner, 100, ts) // previous owner should earn from bid
			},
		},
		{
			name:     "process multiple - mixed OPOs",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			opos: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, opoExpired, nil, nil),
				genOpo(dymNameB, opoNotExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					// not completed
					Bidder: bidder,
					Price:  coin100,
				}),
				genOpo(dymNameC, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					Bidder: bidder,
					Price:  coin200,
				}),
				genOpo(dymNameD, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					// completed by min price
					Bidder: bidder,
					Price:  coin100,
				}),
			},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
				dymNameB.Name: opoNotExpiredEpoch,
				dymNameC.Name: opoExpiredEpoch,
				dymNameD.Name: opoExpiredEpoch,
			},
			preMintModuleBalance: 450,
			wantErr:              false,
			wantMapExpiryByDymName: map[string]int64{
				dymNameB.Name: opoNotExpiredEpoch,
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				ts := nts(t, dk, bk, ctx)

				// OPO for Dym-Name A is expired without any bid/winner
				requireDymNameNotChanged(dymNameA, ts)
				requireNoActiveOPO(dymNameA, ts)
				requireHistoricalOPOs(dymNameA, 1, ts)

				// OPO for Dym-Name B not yet finished
				requireDymNameNotChanged(dymNameB, ts)
				opoB := ts.dk.GetOpenPurchaseOrder(ts.ctx, dymNameB.Name)
				require.NotNil(t, opoB)
				requireHistoricalOPOs(dymNameB, 0, ts)

				// OPO for Dym-Name C is completed with winner
				requireOwnerChanged(dymNameC, bidder, ts)
				requireNoActiveOPO(dymNameC, ts)
				requireHistoricalOPOs(dymNameC, 1, ts)

				// OPO for Dym-Name D is completed with winner
				requireOwnerChanged(dymNameD, bidder, ts)
				requireNoActiveOPO(dymNameD, ts)
				requireHistoricalOPOs(dymNameD, 1, ts)

				requireModuleBalance(150, ts)

				requireAccountBalance(owner, 300, ts) // price from 2 completed OPO
			},
		},
		{
			name:     "should do nothing if invalid epoch identifier",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB, dymNameC, dymNameD},
			opos: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, opoExpired, nil, nil),
				genOpo(dymNameB, opoNotExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					// not completed
					Bidder: bidder,
					Price:  coin100,
				}),
				genOpo(dymNameC, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					Bidder: bidder,
					Price:  coin200,
				}),
				genOpo(dymNameD, opoExpired, &coin200, &dymnstypes.OpenPurchaseOrderBid{
					// completed by min price
					Bidder: bidder,
					Price:  coin100,
				}),
			},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
				dymNameB.Name: opoNotExpiredEpoch,
				dymNameC.Name: opoExpiredEpoch,
				dymNameD.Name: opoExpiredEpoch,
			},
			preMintModuleBalance:  450,
			customEpochIdentifier: "another",
			wantErr:               false,
			wantMapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
				dymNameB.Name: opoNotExpiredEpoch,
				dymNameC.Name: opoExpiredEpoch,
				dymNameD.Name: opoExpiredEpoch,
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				ts := nts(t, dk, bk, ctx)

				requireDymNameNotChanged(dymNameA, ts)
				requireDymNameNotChanged(dymNameB, ts)
				requireDymNameNotChanged(dymNameC, ts)
				requireDymNameNotChanged(dymNameD, ts)

				requireModuleBalance(450, ts)

				requireAccountBalance(owner, 0, ts)
			},
		},
		{
			name:     "should remove expiry reference to non-exists OPO",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB},
			opos: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, opoExpired, nil, nil),
				// no OPO for Dym-Name B
			},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
				dymNameB.Name: opoExpiredEpoch, // no OPO for Dym-Name B but still have reference
			},
			wantErr:                false,
			wantMapExpiryByDymName: map[string]int64{
				// removed reference to Dym-Name A because of processed
				// removed reference to Dym-Name B because OPO not exists
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				// no extra test
			},
		},
		{
			name:     "update expiry if in-correct",
			dymNames: []dymnstypes.DymName{dymNameA, dymNameB},
			opos: []dymnstypes.OpenPurchaseOrder{
				genOpo(dymNameA, opoExpired, nil, nil),
				genOpo(dymNameB, opoNotExpired, nil, nil), // OPO not expired
			},
			mapExpiryByDymName: map[string]int64{
				dymNameA.Name: opoExpiredEpoch,
				dymNameB.Name: opoExpiredEpoch, // incorrect, OPO not expired
			},
			wantErr: false,
			wantMapExpiryByDymName: map[string]int64{
				dymNameB.Name: opoNotExpiredEpoch, // updated
			},
			afterHookTestFunc: func(t *testing.T, dk dymnskeeper.Keeper, bk dymnskeeper.BankKeeper, ctx sdk.Context) {
				// no extra test
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.afterHookTestFunc, "mis-configured test case")

			dk, bk, ctx := setupTest()

			if tt.preMintModuleBalance > 0 {
				err := bk.MintCoins(ctx, dymnstypes.ModuleName, dymnsutils.TestCoins(tt.preMintModuleBalance))
				require.NoError(t, err)
			}

			err := dk.SetActiveOpenPurchaseOrdersExpiration(ctx, dymnstypes.ActiveOpenPurchaseOrdersExpiration{
				ExpiryByName: tt.mapExpiryByDymName,
			})
			require.NoError(t, err)

			for _, dymName := range tt.dymNames {
				err = dk.SetDymName(ctx, dymName)
				require.NoError(t, err)
			}

			for _, opo := range tt.opos {
				err = dk.SetOpenPurchaseOrder(ctx, opo)
				require.NoError(t, err)
			}

			moduleParams := dk.GetParams(ctx)

			useEpochIdentifier := moduleParams.Misc.EndEpochHookIdentifier
			if tt.customEpochIdentifier != "" {
				useEpochIdentifier = tt.customEpochIdentifier
			}

			err = dk.GetEpochHooks().AfterEpochEnd(ctx, useEpochIdentifier, 1)

			defer func() {
				if t.Failed() {
					return
				}

				tt.afterHookTestFunc(t, dk, bk, ctx)

				aope := dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
				if len(tt.wantMapExpiryByDymName) == 0 {
					require.Empty(t, aope.ExpiryByName)
				} else if !reflect.DeepEqual(tt.wantMapExpiryByDymName, aope.ExpiryByName) {
					t.Errorf("want AOPE map %v, got %v", tt.wantMapExpiryByDymName, aope.ExpiryByName)
				}
			}()

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)

				return
			}

			require.NoError(t, err)
		})
	}
}
