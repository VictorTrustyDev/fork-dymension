package keeper

import (
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"

	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func DymNSKeeper(t testing.TB) (dymnskeeper.Keeper, dymnskeeper.BankKeeper, sdk.Context) {
	dymNsStoreKey := sdk.NewKVStoreKey(dymnstypes.StoreKey)
	dymNsMemStoreKey := storetypes.NewMemoryStoreKey(dymnstypes.MemStoreKey)

	authStoreKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	bankStoreKey := sdk.NewKVStoreKey(banktypes.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(dymNsStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(dymNsMemStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(authStoreKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(bankStoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		dymnstypes.Amino,
		dymNsStoreKey,
		dymNsMemStoreKey,
		"DymNSParams",
	)

	authKeeper := authkeeper.NewAccountKeeper(
		cdc,
		authStoreKey,
		authtypes.ProtoBaseAccount,
		map[string][]string{
			banktypes.ModuleName:  {authtypes.Minter, authtypes.Burner},
			dymnstypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		},
		params.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	authtypes.RegisterInterfaces(registry)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		bankStoreKey,
		authKeeper,
		map[string]bool{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	banktypes.RegisterInterfaces(registry)

	k := dymnskeeper.NewKeeper(cdc,
		dymNsStoreKey,
		paramsSubspace,
		bankKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, dymnstypes.DefaultParams())

	return k, bankKeeper, ctx
}
