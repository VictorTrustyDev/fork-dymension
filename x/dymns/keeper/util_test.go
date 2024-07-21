package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func setDymNameWithFunctionsAfter(ctx sdk.Context, dymName dymnstypes.DymName, t *testing.T, dk dymnskeeper.Keeper) {
	require.NoError(t, dk.SetDymName(ctx, dymName))
	require.NoError(t, dk.AfterDymNameOwnerChanged(ctx, dymName.Name))
	require.NoError(t, dk.AfterDymNameConfigChanged(ctx, dymName.Name))
}

func requireErrorContains(t *testing.T, err error, contains string) {
	require.Error(t, err)
	require.NotEmpty(t, contains, "mis-configured test")
	require.Contains(t, err.Error(), contains)
}

func requireErrorFContains(t *testing.T, f func() error, contains string) {
	requireErrorContains(t, f(), contains)
}
