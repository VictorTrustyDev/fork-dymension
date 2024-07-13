package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func TestGetSetParams(t *testing.T) {
	dk, _, _, ctx := testkeeper.DymNSKeeper(t)
	params := dymnstypes.DefaultParams()

	err := dk.SetParams(ctx, params)
	require.NoError(t, err)

	require.Equal(t, params, dk.GetParams(ctx))

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Misc.BeginEpochHookIdentifier = ""
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Price.PriceDenom = ""
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Alias.ByChainId = map[string]dymnstypes.AliasesOfChainId{
			"@": {},
		}
		require.Error(t, dk.SetParams(ctx, params))
	})

	t.Run("can not set invalid params", func(t *testing.T) {
		params := dymnstypes.DefaultParams()
		params.Misc.DaysGracePeriod = -1
		require.Error(t, dk.SetParams(ctx, params))
	})
}
