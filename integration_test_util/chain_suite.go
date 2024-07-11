package integration_test_util

//goland:noinspection SpellCheckingInspection,GoSnakeCaseUsage
import (
	"fmt"
	cdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	tmstate "github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/store"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cometbft/cometbft/version"
	"github.com/cosmos/cosmos-sdk/baseapp"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmostxtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govlegacytypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	chainapp "github.com/dymensionxyz/dymension/v3/app"
	"github.com/dymensionxyz/dymension/v3/app/params"
	itutiltypes "github.com/dymensionxyz/dymension/v3/integration_test_util/types"
	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	sequencertypes "github.com/dymensionxyz/dymension/v3/x/sequencer/types"
	streamerkeeper "github.com/dymensionxyz/dymension/v3/x/streamer/keeper"
	streamertypes "github.com/dymensionxyz/dymension/v3/x/streamer/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethermint_hd "github.com/evmos/ethermint/crypto/hd"
	kvindexer "github.com/evmos/ethermint/indexer"
	rpcbackend "github.com/evmos/ethermint/rpc/backend"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/require"
	"math"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"
)

// ChainIntegrationTestSuite is a helper for Chain integration test.
type ChainIntegrationTestSuite struct {
	t                    *testing.T
	require              *require.Assertions
	muTest               sync.RWMutex
	mu                   sync.RWMutex
	ibcSuite             *ChainsIbcIntegrationTestSuite
	historicalContext    map[int64]sdk.Context
	tempHolder           *itutiltypes.TemporaryHolder
	logger               log.Logger
	EncodingConfig       params.EncodingConfig
	ChainConstantsConfig itutiltypes.ChainConstantConfig
	DB                   *itutiltypes.MemDB
	TendermintApp        itutiltypes.TendermintApp
	ChainApp             itutiltypes.ChainApp
	ValidatorSet         *tmtypes.ValidatorSet
	CurrentContext       sdk.Context
	ValidatorAccounts    itutiltypes.TestAccounts
	WalletAccounts       itutiltypes.TestAccounts
	ModuleAccounts       map[string]authtypes.ModuleAccountI
	QueryClients         *itutiltypes.QueryClients
	EvmTxIndexer         *kvindexer.KVIndexer
	RpcBackend           *rpcbackend.Backend
	EthSigner            ethtypes.Signer
	TestConfig           itutiltypes.TestConfig
}

// CreateChainIntegrationTestSuite initialize an integration test suite using default configuration.
func CreateChainIntegrationTestSuite(t *testing.T, r *require.Assertions) *ChainIntegrationTestSuite {
	return CreateChainIntegrationTestSuiteFromChainConfig(t, r, IntegrationTestChain1)
}

//goland:noinspection SpellCheckingInspection
var IntegrationTestChain1 = itutiltypes.ChainConfig{
	CosmosChainId:            "dymension_100-1",
	BaseDenom:                "udym",
	Bech32Prefix:             "dym",
	EvmChainId:               100,
	DisabledContractCreation: true,
}

//goland:noinspection SpellCheckingInspection
var IntegrationTestChain2 = itutiltypes.ChainConfig{
	CosmosChainId:            "froopyland_100-1",
	BaseDenom:                "udym",
	Bech32Prefix:             "dym",
	EvmChainId:               100,
	DisabledContractCreation: true,
}

// CreateChainIntegrationTestSuiteFromChainConfig initialize an integration test suite from a given chain config.
func CreateChainIntegrationTestSuiteFromChainConfig(t *testing.T, r *require.Assertions, chainCfg itutiltypes.ChainConfig) *ChainIntegrationTestSuite {
	chainCfg.DisableTendermint = true
	// force disable Tendermint to save  because Dymension does not enable contract call at this moment and won't invoke RPC backend.

	//goland:noinspection SpellCheckingInspection
	const balancePerAccount = 2

	chainCfg.EvmChainIdBigInt = big.NewInt(chainCfg.EvmChainId)

	appEncodingCfg := chainapp.MakeEncodingConfig()

	//goland:noinspection SpellCheckingInspection
	testConfig := itutiltypes.TestConfig{
		SecondaryDenomUnits: []banktypes.DenomUnit{
			{
				Denom:    "utwo",
				Exponent: 6,
			},
			{
				Denom:    "uthree",
				Exponent: 8,
			},
		},
		InitBalanceAmount:        sdk.NewInt(int64(balancePerAccount * math.Pow10(18))),
		DefaultFeeAmount:         sdk.NewInt(int64(math.Pow10(16))),
		DisableTendermint:        chainCfg.DisableTendermint,
		DisabledContractCreation: chainCfg.DisabledContractCreation,
	}

	clientCtx := cosmosclient.Context{}.
		WithChainID(chainCfg.CosmosChainId).
		WithCodec(appEncodingCfg.Codec).
		WithInterfaceRegistry(appEncodingCfg.InterfaceRegistry).
		WithTxConfig(appEncodingCfg.TxConfig).
		WithLegacyAmino(appEncodingCfg.Amino).
		WithKeyringOptions(ethermint_hd.EthSecp256k1Option())

	tempHolder := itutiltypes.NewTemporaryHolder()

	// Setup assertions
	if r == nil {
		r = require.New(t)
	}

	// Setup Test accounts

	validatorAccounts := newValidatorAccounts(t)
	if chainCfg.DisableTendermint {
		// no-op
	} else {
		// test tendermint use only one validator
		validatorAccounts = []*itutiltypes.TestAccount{validatorAccounts.Number(1)}
	}

	walletAccounts := newWalletsAccounts(t)

	// Init database
	sharedDb := itutiltypes.WrapCometBFT(cdb.NewMemDB())
	evmIndexerDb := cdb.NewMemDB() // use dedicated db for EVM Tx-Indexer to prevent data corruption

	// Setup chain app
	genesisAccountBalance := sdk.NewCoins(
		sdk.NewCoin(chainCfg.BaseDenom, testConfig.InitBalanceAmount),
	)
	for _, secondaryDenomUnit := range testConfig.SecondaryDenomUnits {
		genesisAccountBalance = genesisAccountBalance.Add(
			sdk.NewCoin(secondaryDenomUnit.Denom, testConfig.InitBalanceAmount),
		)
	}
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger = log.NewFilter(logger, log.AllowError())
	app, tmApp, valSet := itutiltypes.NewChainApp(chainCfg, testConfig, appEncodingCfg, sharedDb, validatorAccounts, walletAccounts, genesisAccountBalance, tempHolder, logger)
	baseApp := app.BaseApp()

	header := createFirstBlockHeader(
		chainCfg.CosmosChainId,
		validatorAccounts.Number(1).GetConsensusAddress(),
	)
	ctx := baseApp.NewContext(false, header)

	evmParams := app.EvmKeeper().GetParams(ctx)
	evmParams.EvmDenom = chainCfg.BaseDenom
	err := app.EvmKeeper().SetParams(ctx, evmParams)
	require.NoError(t, err)

	// Setup validators
	for _, validatorAccount := range validatorAccounts {
		val, err := stakingtypes.NewValidator(
			validatorAccount.GetValidatorAddress(),
			validatorAccount.GetSdkPubKey(),
			stakingtypes.Description{},
		)
		require.NoError(t, err)

		val = stakingkeeper.TestingUpdateValidator(app.StakingKeeper(), ctx, val, true)

		refKeeper := reflect.ValueOf(app.StakingKeeper())
		if refKeeper.Kind() == reflect.Ptr {
			refKeeper = reflect.Indirect(refKeeper.Elem())
		}

		err = app.DistributionKeeper().Hooks().AfterValidatorCreated(ctx, val.GetOperator())

		err = app.StakingKeeper().SetValidatorByConsAddr(ctx, val)
		require.NoError(t, err)
	}

	result := &ChainIntegrationTestSuite{
		t:                 t,
		require:           r,
		muTest:            sync.RWMutex{},
		mu:                sync.RWMutex{},
		historicalContext: make(map[int64]sdk.Context),
		tempHolder:        tempHolder,
		logger:            logger,
		EncodingConfig:    appEncodingCfg,
		ChainConstantsConfig: itutiltypes.NewChainConstantConfig(
			chainCfg.CosmosChainId,
			chainCfg.BaseDenom,
		),
		DB:                sharedDb,
		ChainApp:          app,
		TendermintApp:     tmApp,
		ValidatorSet:      valSet,
		CurrentContext:    ctx,
		ValidatorAccounts: validatorAccounts,
		WalletAccounts:    walletAccounts,
		ModuleAccounts:    make(map[string]authtypes.ModuleAccountI),
		EvmTxIndexer:      kvindexer.NewKVIndexer(evmIndexerDb, log.NewNopLogger(), clientCtx),
		EthSigner:         ethtypes.LatestSignerForChainID(chainCfg.EvmChainIdBigInt),
		TestConfig:        testConfig,
	}

	if chainCfg.DisableTendermint {

		result.Commit() // Commit the initial block
	} else {
		time.Sleep(300 * time.Millisecond)
		result.Commit()
	}

	result.CreateAllQueryClientsAndRpcBackend()

	accounts, _ := result.QueryClients.Auth.ModuleAccounts(nil, &authtypes.QueryModuleAccountsRequest{})
	for _, acc := range accounts.Accounts {
		var account authtypes.AccountI
		err = appEncodingCfg.InterfaceRegistry.UnpackAny(acc, &account)
		require.NoError(t, err)
		moduleAccount, ok := account.(authtypes.ModuleAccountI)
		require.True(t, ok)
		result.ModuleAccounts[moduleAccount.GetName()] = moduleAccount
	}

	return result
}

func (suite *ChainIntegrationTestSuite) T() *testing.T {
	suite.muTest.RLock()
	defer suite.muTest.RUnlock()
	return suite.t
}

func (suite *ChainIntegrationTestSuite) Require() *require.Assertions {
	suite.muTest.RLock()
	defer suite.muTest.RUnlock()
	return suite.require
}

// Cleanup cleans up the ChainIntegrationTestSuite.
// This method should be called after each test or suite, depends on the tactic you shut down the Integration chain.
func (suite *ChainIntegrationTestSuite) Cleanup() {
	if suite == nil {
		return
	}

	if suite.HasTendermint() {
		suite.TendermintApp.Shutdown()
	}

	if suite.tempHolder != nil {
		if tempFiles, anyTemp := suite.tempHolder.GetTempFiles(); anyTemp {
			for _, file := range tempFiles {
				err := os.RemoveAll(file)
				if err != nil {
					fmt.Println("Failed to remove temp file", file)
					fmt.Println(err)
				}
			}
		}
	}
}

// BaseApp returns the BaseApp instance of the Integrated chain.
func (suite *ChainIntegrationTestSuite) BaseApp() *baseapp.BaseApp {
	return suite.ChainApp.BaseApp()
}

// CreateAllQueryClientsAndRpcBackend creates all query clients and RPC backend instance at recent block height.
// This method should be called after each commit to refresh the query clients.
func (suite *ChainIntegrationTestSuite) CreateAllQueryClientsAndRpcBackend() {
	suite.QueryClients = suite.QueryClientsAt(0)
	suite.RpcBackend = suite.RpcBackendAt(0)
}

// ContextAt returns the context at a given context block height.
func (suite *ChainIntegrationTestSuite) ContextAt(height int64) sdk.Context {
	if height == 0 {
		height = suite.GetLatestBlockHeight()
	}

	if ctx, found := suite.historicalContext[height]; found {
		return ctx
	}

	qCtx, err := suite.createAppQueryContext(height, false)
	suite.Require().NoError(err)

	return qCtx
}

// createAppQueryContext returns the query context at a given context block height.
// Used as a helper method to create query context to adapt with older version of Cosmos-SDK BaseApp,
// which does not expose CreateQueryContext method.
func (suite *ChainIntegrationTestSuite) createAppQueryContext(height int64, prove bool) (sdk.Context, error) {
	if height < 0 {
		panic("height cannot be negative")
	}

	// use custom query multistore if provided
	var qms sdk.MultiStore
	if qms == nil {
		qms = sdk.MultiStore(suite.BaseApp().CommitMultiStore())
	}

	lastBlockHeight := qms.LatestVersion()
	if height > lastBlockHeight {
		return sdk.Context{},
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidHeight,
				"cannot query with height in the future; please provide a valid height",
			)
	}

	// when a client did not provide a query height, manually inject the latest
	if height == 0 {
		height = lastBlockHeight
	}

	if height <= 1 && prove {
		return sdk.Context{},
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			)
	}

	cacheMS, err := qms.CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{},
			sdkerrors.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"failed to load state at height %d; %s (latest height: %d)", height, err, lastBlockHeight,
			)
	}

	// branch the commit-multistore for safety
	ctx := sdk.NewContext(
		cacheMS, suite.CurrentContext.BlockHeader(), true, suite.BaseApp().Logger(),
	).WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(suite.NewBaseCoin(0)))).WithBlockHeight(height)

	return ctx, nil
}

// QueryClientsAt returns the list of query client instance that connects to store data at a given context block height.
func (suite *ChainIntegrationTestSuite) QueryClientsAt(height int64) *itutiltypes.QueryClients {
	var sdkContext sdk.Context
	if suite.HasTendermint() {
		if height == 0 {
			height = suite.GetLatestBlockHeight()
		}
		sdkContext = suite.CurrentContext
		if height > 0 {
			var err error
			sdkContext, err = suite.createAppQueryContext(height, false)
			suite.Require().NoError(err)
		}
	} else if height == 0 || height == suite.GetLatestBlockHeight() {
		// latest block
		sdkContext = suite.CurrentContext
	} else {
		var err error
		sdkContext, err = suite.createAppQueryContext(height, false)
		suite.Require().NoError(err)
	}

	queryHelper := baseapp.NewQueryServerTestHelper(sdkContext, suite.ChainApp.InterfaceRegistry())

	authtypes.RegisterQueryServer(queryHelper, suite.ChainApp.AccountKeeper())
	authQueryClient := authtypes.NewQueryClient(queryHelper)

	banktypes.RegisterQueryServer(queryHelper, suite.ChainApp.BankKeeper())
	bankQueryClient := banktypes.NewQueryClient(queryHelper)

	distributiontypes.RegisterQueryServer(queryHelper, distributionkeeper.Querier{Keeper: suite.ChainApp.DistributionKeeper()})
	distributionQueryClient := distributiontypes.NewQueryClient(queryHelper)

	evmtypes.RegisterQueryServer(queryHelper, suite.ChainApp.EvmKeeper())
	evmQueryClient := evmtypes.NewQueryClient(queryHelper)

	feemarkettypes.RegisterQueryServer(queryHelper, suite.ChainApp.FeeMarketKeeper())
	feeMarketQueryClient := feemarkettypes.NewQueryClient(queryHelper)

	govv1types.RegisterQueryServer(queryHelper, suite.ChainApp.GovKeeper())
	govV1QueryClient := govv1types.NewQueryClient(queryHelper)

	govlegacytypes.RegisterQueryServer(queryHelper, govkeeper.NewLegacyQueryServer(suite.ChainApp.GovKeeper()))
	govLegacyQueryClient := govlegacytypes.NewQueryClient(queryHelper)

	ibctransfertypes.RegisterQueryServer(queryHelper, suite.ChainApp.IbcTransferKeeper())
	ibcTransferQueryClient := ibctransfertypes.NewQueryClient(queryHelper)

	rollapptypes.RegisterQueryServer(queryHelper, suite.ChainApp.RollAppKeeper())
	rollappQueryClient := rollapptypes.NewQueryClient(queryHelper)

	sequencertypes.RegisterQueryServer(queryHelper, suite.ChainApp.SequencerKeeper())
	sequencerQueryClient := sequencertypes.NewQueryClient(queryHelper)

	slashingtypes.RegisterQueryServer(queryHelper, suite.ChainApp.SlashingKeeper())
	slashingQueryClient := slashingtypes.NewQueryClient(queryHelper)

	stakingtypes.RegisterQueryServer(queryHelper, stakingkeeper.Querier{Keeper: suite.ChainApp.StakingKeeper()})
	stakingQueryClient := stakingtypes.NewQueryClient(queryHelper)

	streamertypes.RegisterQueryServer(queryHelper, streamerkeeper.NewQuerier(*suite.ChainApp.StreamerKeeper()))
	streamerQueryClient := streamertypes.NewQueryClient(queryHelper)

	serviceClient := cosmostxtypes.NewServiceClient(queryHelper)

	rpcQueryClient := rpctypes.QueryClient{
		ServiceClient: serviceClient,
		QueryClient:   evmQueryClient,
		FeeMarket:     feeMarketQueryClient,
	}

	var tendermintRpcHttpClient *httpclient.HTTP
	if suite.HasTendermint() {
		rpcAddr26657, supported := suite.TendermintApp.GetRpcAddr()
		suite.Require().True(supported)

		httpClient26657, err := jsonrpcclient.DefaultHTTPClient(rpcAddr26657)
		suite.Require().NoError(err)

		tendermintRpcHttpClient, err = httpclient.NewWithClient(rpcAddr26657, "/websocket", httpClient26657)
		suite.Require().NoError(err)

		err = tendermintRpcHttpClient.Start()
		suite.Require().NoError(err)
	}

	clientQueryCtx := cosmosclient.Context{}.
		WithChainID(suite.ChainConstantsConfig.GetCosmosChainID()).
		WithCodec(suite.EncodingConfig.Codec).
		WithInterfaceRegistry(suite.EncodingConfig.InterfaceRegistry).
		WithTxConfig(suite.EncodingConfig.TxConfig).
		WithLegacyAmino(suite.EncodingConfig.Amino).
		WithKeyringOptions(ethermint_hd.EthSecp256k1Option())

	if height > 0 {
		clientQueryCtx = clientQueryCtx.WithHeight(height)
	}

	if suite.HasTendermint() {
		clientQueryCtx = clientQueryCtx.WithClient(tendermintRpcHttpClient)
	}

	cosmostxtypes.RegisterServiceServer(
		queryHelper,
		authtx.NewTxServer(clientQueryCtx, suite.BaseApp().Simulate, suite.ChainApp.InterfaceRegistry()),
	)

	return &itutiltypes.QueryClients{
		GrpcConnection:          queryHelper,
		ClientQueryCtx:          clientQueryCtx,
		TendermintRpcHttpClient: tendermintRpcHttpClient,
		Auth:                    authQueryClient,
		Bank:                    bankQueryClient,
		Distribution:            distributionQueryClient,
		EVM:                     evmQueryClient,
		GovV1:                   govV1QueryClient,
		GovLegacy:               govLegacyQueryClient,
		IbcTransfer:             ibcTransferQueryClient,
		RollApp:                 rollappQueryClient,
		Sequencer:               sequencerQueryClient,
		Slashing:                slashingQueryClient,
		Staking:                 stakingQueryClient,
		Streamer:                streamerQueryClient,
		ServiceClient:           serviceClient,
		Rpc:                     &rpcQueryClient,
	}
}

// RpcBackendAt returns the RPC-backend instance at a given context block height.
func (suite *ChainIntegrationTestSuite) RpcBackendAt(height int64) *rpcbackend.Backend {
	queryClients := suite.QueryClientsAt(height)
	rpcServerCtx := server.NewDefaultContext()

	rpcBackend := rpcbackend.NewBackend(rpcServerCtx, rpcServerCtx.Logger, queryClients.ClientQueryCtx, false, suite.EvmTxIndexer)

	// override the query client with the mock query client, for changing query context
	getFieldQueryClient := func() reflect.Value {
		return reflect.Indirect(reflect.ValueOf(rpcBackend).Elem()).FieldByName("queryClient")
	}
	fieldQueryClient := getFieldQueryClient()
	reflect.NewAt(fieldQueryClient.Type(), unsafe.Pointer(fieldQueryClient.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(queryClients.Rpc))

	return rpcBackend
}

// GetLatestBlockHeight returns the most recent block height.
func (suite *ChainIntegrationTestSuite) GetLatestBlockHeight() int64 {
	if suite.HasTendermint() {
		// because Tendermint auto-commit blocks so the CurrentContext property might out-dated
		return suite.BaseApp().LastBlockHeight()
	}

	return suite.CurrentContext.BlockHeight()
}

// WaitNextBlockOrCommit returns the most recent block height beside the following logic:
//
// - When Tendermint is Enabled, it waits for the next block to be committed before returning result.
//
// - When Tendermint is Disabled, it triggers commit block and starts a new block with an updated context.
//
// USE-CASE for this: you want to submit one or multiple txs and have sometime to know the executed block,
// while Tendermint auto commit blocks.
func (suite *ChainIntegrationTestSuite) WaitNextBlockOrCommit() int64 {
	if !suite.HasTendermint() {
		suite.Commit()
		return suite.GetLatestBlockHeight()
	}

	oldHeight := suite.GetLatestBlockHeight()
	var currentHeight int64
	for {
		currentHeight = suite.GetLatestBlockHeight()
		if currentHeight > oldHeight {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return currentHeight
}

// Commit commits and starts a new block with an updated context.
func (suite *ChainIntegrationTestSuite) Commit() {
	if suite.ibcSuite != nil { // ibc-connected chains must be committed together
		suite.ibcSuite.CommitAllChains()
	} else {
		suite.commitAndBeginBlockAfter(1 * time.Hour)
	}
}

// ibcSuiteCommit is a helper function to commit with custom block time equals to IBC setup
func (suite *ChainIntegrationTestSuite) ibcSuiteCommit() {
	suite.commitAndBeginBlockAfter(5 * time.Second)
}

// commitAndBeginBlockAfter commits a block at a given time.
func (suite *ChainIntegrationTestSuite) commitAndBeginBlockAfter(t time.Duration) {
	suite.mu.Lock()
	defer suite.mu.Unlock()

	defer func() {
		suite.CreateAllQueryClientsAndRpcBackend()
	}()

	var newCtx sdk.Context
	var newValSet *tmtypes.ValidatorSet

	if suite.HasTendermint() {
		// awaiting next block generated by tendermint
		originalHeight := suite.GetLatestBlockHeight()
		var latestHeight int64
		for {
			time.Sleep(10 * time.Millisecond)
			latestHeight = suite.GetLatestBlockHeight()
			if latestHeight > originalHeight {
				break
			}
		}

		blockStore, stateStore := suite.GetBlockStoreAndStateStore()

		tmBlk := blockStore.LoadBlock(latestHeight)
		valSet, err := stateStore.LoadValidators(latestHeight)
		suite.Require().NoErrorf(err, "failed to load validator set for block %d", latestHeight)

		header := tmBlk.Header.ToProto()
		ctx := suite.createNewContext(suite.CurrentContext, *header)
		suite.triggerEvmIndexer(latestHeight, blockStore, stateStore) // trigger EVM Tx-Indexer indexing data to latest

		newCtx = ctx
		newValSet = valSet
	} else {
		// manually commit block and move to next
		backupContext := suite.CurrentContext

		nextCtx, nextValSet, err := suite.commitAndCreateNewCtx(suite.CurrentContext, t, suite.ValidatorSet)
		suite.Require().NoError(err)
		suite.Require().Equalf(suite.CurrentContext.BlockHeight()+1, nextCtx.BlockHeight(), "next block height must be increased by 1")

		suite.historicalContext[backupContext.BlockHeight()] = backupContext

		newCtx = nextCtx
		newValSet = nextValSet
	}

	suite.CurrentContext = newCtx
	suite.ValidatorSet = newValSet
}

// GetIbcTimeoutHeight returns a timeout height for IBC packet, based on recent block, plus the offset.
func (suite *ChainIntegrationTestSuite) GetIbcTimeoutHeight(offsetHeight int64) ibcclienttypes.Height {
	chainId := suite.ChainConstantsConfig.GetCosmosChainID()
	idx := strings.LastIndex(chainId, "-")
	rev := chainId[idx+1:]
	revInt, err := strconv.ParseUint(rev, 10, 64)
	suite.Require().NoError(err)
	return ibcclienttypes.NewHeight(revInt, uint64(suite.GetLatestBlockHeight()+offsetHeight))
}

// triggerEvmIndexer indexes EVM txs from blockStore and stateStore, upto latestHeight.
func (suite *ChainIntegrationTestSuite) triggerEvmIndexer(latestHeight int64, blockStore *store.BlockStore, stateStore tmstate.Store) {
	suite.Require().NotZero(latestHeight)
	suite.Require().NotNil(blockStore)
	suite.Require().NotNil(stateStore)

	lastIndexedHeight, err := suite.EvmTxIndexer.LastIndexedBlock()
	suite.Require().NoError(err)

	if lastIndexedHeight >= latestHeight {
		return
	}

	if lastIndexedHeight < 0 {
		lastIndexedHeight = 0
	}

	var ch int64
	for ch = lastIndexedHeight + 1; ch <= latestHeight; ch++ {
		tmBlk := blockStore.LoadBlock(ch)
		tmAbciResponse, err := stateStore.LoadABCIResponses(ch)
		suite.Require().NoErrorf(err, "failed to load abci response for block %d", ch)
		err = suite.EvmTxIndexer.IndexBlock(tmBlk, tmAbciResponse.DeliverTxs)
		suite.Require().NoErrorf(err, "failed to index block %d", ch)
	}
}

// GetBlockStoreAndStateStore returns blockStore and stateStore if Tendermint is Enabled.
//
// WARN: if Tendermint is Disabled, the call will panic.
func (suite *ChainIntegrationTestSuite) GetBlockStoreAndStateStore() (*store.BlockStore, tmstate.Store) {
	suite.EnsureTendermint()
	blockStore := store.NewBlockStore(suite.DB)
	stateStore := tmstate.NewStore(suite.DB, tmstate.StoreOptions{
		DiscardABCIResponses: false,
	})
	return blockStore, stateStore
}

// createFirstBlockHeader creates a new Tendermint header, with context 1, for testing purposes.
func createFirstBlockHeader(
	chainID string,
	proposer sdk.ConsAddress,
) tmproto.Header {
	//goland:noinspection SpellCheckingInspection
	return tmproto.Header{
		ChainID:         chainID,
		Height:          1,
		Time:            time.Now().UTC(),
		ValidatorsHash:  nil,
		AppHash:         nil,
		ProposerAddress: proposer.Bytes(),
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		DataHash:           tmhash.Sum([]byte("data")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
	}
}
