package types

//goland:noinspection SpellCheckingInspection
import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	delayedackkeeper "github.com/dymensionxyz/dymension/v3/x/delayedack/keeper"
	rollappkeeper "github.com/dymensionxyz/dymension/v3/x/rollapp/keeper"
	sequencerkeeper "github.com/dymensionxyz/dymension/v3/x/sequencer/keeper"
	streamerkeeper "github.com/dymensionxyz/dymension/v3/x/streamer/keeper"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
)

type ChainApp interface {
	App() abci.Application
	BaseApp() *baseapp.BaseApp
	IbcTestingApp() ibctesting.TestingApp
	InterfaceRegistry() codectypes.InterfaceRegistry

	// Keepers

	AccountKeeper() *authkeeper.AccountKeeper
	BankKeeper() bankkeeper.Keeper
	DelayedAckKeeper() *delayedackkeeper.Keeper
	DistributionKeeper() distributionkeeper.Keeper
	EvmKeeper() *evmkeeper.Keeper
	FeeMarketKeeper() *feemarketkeeper.Keeper
	GovKeeper() *govkeeper.Keeper
	IbcTransferKeeper() *ibctransferkeeper.Keeper
	IbcKeeper() *ibckeeper.Keeper
	RollAppKeeper() *rollappkeeper.Keeper
	SequencerKeeper() *sequencerkeeper.Keeper
	SlashingKeeper() *slashingkeeper.Keeper
	StakingKeeper() *stakingkeeper.Keeper
	StreamerKeeper() *streamerkeeper.Keeper

	// Tx

	FundAccount(ctx sdk.Context, account *TestAccount, amounts sdk.Coins) error
}
