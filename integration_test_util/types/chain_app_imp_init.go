package types

//goland:noinspection SpellCheckingInspection
import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/simapp"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cosmosed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	chainapp "github.com/dymensionxyz/dymension/v3/app"
	"github.com/dymensionxyz/dymension/v3/app/params"
	itutilutils "github.com/dymensionxyz/dymension/v3/integration_test_util/utils"
	"github.com/ethereum/go-ethereum/common"
	evmostypes "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"strings"
	"time"
)

var defaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   40000000, // 40m
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

const TendermintGovVotingPeriod = 5 * time.Second

func NewChainApp(chainCfg ChainConfig, testConfig TestConfig, encCfg params.EncodingConfig, db *MemDB, validatorAccounts TestAccounts, walletAccounts TestAccounts, genesisAccountBalance sdk.Coins, tempHolder *TemporaryHolder, logger log.Logger) (chainApp ChainApp, tendermintApp TendermintApp, validatorSet *tmtypes.ValidatorSet) {
	defaultNodeHome := chainapp.DefaultNodeHome
	moduleBasics := chainapp.ModuleBasics

	// create validator set
	var validators []*tmtypes.Validator
	for _, validatorAccount := range validatorAccounts {
		//goland:noinspection GoDeprecation
		pv := mock.PV{
			PrivKey: &cosmosed25519.PrivKey{
				Key: ed25519.NewKeyFromSeed(validatorAccount.PrivateKey.Key),
			},
		}
		pubKey, err := pv.GetPubKey()
		if err != nil {
			panic(err)
		}
		validators = append(validators, tmtypes.NewValidator(pubKey, 1))
	}
	valSet := tmtypes.NewValidatorSet(validators)

	// generate genesis accounts
	var genesisValidatorAccounts []authtypes.GenesisAccount
	var genesisWalletAccounts []authtypes.GenesisAccount
	var genesisBalances []banktypes.Balance
	var signingInfos []slashingtypes.SigningInfo
	for i, account := range append(validatorAccounts, walletAccounts...) {
		acc := &evmostypes.EthAccount{
			BaseAccount: authtypes.NewBaseAccount(account.GetCosmosAddress(), account.GetPubKey(), uint64(i), 0),
			CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
		}
		if account.Type == TestAccountTypeValidator {
			genesisValidatorAccounts = append(genesisValidatorAccounts, acc)

			signingInfos = append(signingInfos, slashingtypes.SigningInfo{
				Address: account.GetConsensusAddress().String(),
				ValidatorSigningInfo: slashingtypes.ValidatorSigningInfo{
					Address:             account.GetConsensusAddress().String(),
					StartHeight:         0,
					IndexOffset:         0,
					JailedUntil:         time.Time{},
					Tombstoned:          false,
					MissedBlocksCounter: 0,
				},
			})
		} else if account.Type == TestAccountTypeWallet {
			genesisWalletAccounts = append(genesisWalletAccounts, acc)
		} else {
			panic(fmt.Sprintf("unknown account type %d", account.Type))
		}
		genesisBalances = append(genesisBalances, banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   genesisAccountBalance,
		})
	}

	app := chainapp.New(
		logger,           // logger
		db.AsCometBFT(),  // db
		nil,              // trace store
		true,             // load latest
		map[int64]bool{}, // skipUpgradeHeights
		defaultNodeHome,  // homePath
		0,                // invCheckPeriod
		encCfg,           // encodingConfig
		disableCrisisInvariantsChecking{},
		baseapp.SetChainID(chainCfg.CosmosChainId), // baseAppOptions
	)

	// init chain must be called to stop deliverState from being nil
	genesisState := moduleBasics.DefaultGenesis(encCfg.Codec)

	genesisState = genesisStateWithValSet(chainCfg, testConfig, encCfg.Codec, genesisState, valSet, genesisValidatorAccounts, genesisWalletAccounts, genesisBalances, signingInfos)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	cai := &chainAppImp{
		app: app,
	}

	genesisDoc := tmtypes.GenesisDoc{
		GenesisTime:   time.Time{},
		ChainID:       chainCfg.CosmosChainId,
		InitialHeight: 0,
		ConsensusParams: &tmtypes.ConsensusParams{
			Block: tmtypes.BlockParams{
				MaxBytes: defaultConsensusParams.Block.MaxBytes,
				MaxGas:   defaultConsensusParams.Block.MaxGas,
			},
			Evidence: tmtypes.EvidenceParams{
				MaxAgeNumBlocks: defaultConsensusParams.Evidence.MaxAgeNumBlocks,
				MaxAgeDuration:  defaultConsensusParams.Evidence.MaxAgeDuration,
				MaxBytes:        defaultConsensusParams.Evidence.MaxBytes,
			},
			Validator: tmtypes.ValidatorParams{
				PubKeyTypes: defaultConsensusParams.Validator.PubKeyTypes,
			},
			Version: tmtypes.VersionParams{
				App: 0,
			},
		},
		Validators: make([]tmtypes.GenesisValidator, len(valSet.Validators)),
		AppHash:    nil,
		AppState:   stateBytes,
	}

	for i, validator := range valSet.Validators {
		genesisDoc.Validators[i] = tmtypes.GenesisValidator{
			Address: validator.Address,
			PubKey:  validator.PubKey,
			Power:   validator.VotingPower,
			Name:    "",
		}
	}
	tempHolder.CacheGenesisDoc(&genesisDoc)

	if chainCfg.DisableTendermint {
		app.InitChain(abci.RequestInitChain{
			ChainId: chainCfg.CosmosChainId,
			ConsensusParams: &tmproto.ConsensusParams{
				Block: &tmproto.BlockParams{
					MaxBytes: defaultConsensusParams.Block.MaxBytes,
					MaxGas:   defaultConsensusParams.Block.MaxGas,
				},
				Evidence:  defaultConsensusParams.Evidence,
				Validator: defaultConsensusParams.Validator,
			},
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
			InitialHeight: 0,
		})
		tendermintApp = nil
	} else {
		validator := validatorAccounts.Number(1)
		if validator.GetValidatorAddress().String() != sdk.ValAddress(validator.GetPubKey().Address()).String() {
			panic("validator address does not match")
		}
		node, rpcPort, tempFiles := itutilutils.StartTendermintNode(app, &genesisDoc, db, validator.GetTmPrivKey(), logger)
		for _, tempFile := range tempFiles {
			tempHolder.AddTempFile(tempFile)
		}
		tendermintApp = NewTendermintApp(node, rpcPort)
	}

	return cai, tendermintApp, valSet
}

func genesisStateWithValSet(chainCfg ChainConfig, testConfig TestConfig, codec codec.Codec, genesisState map[string]json.RawMessage, valSet *tmtypes.ValidatorSet, genesisValidatorAccounts []authtypes.GenesisAccount, genesisWalletAccounts []authtypes.GenesisAccount, balances []banktypes.Balance, signingInfos []slashingtypes.SigningInfo) simapp.GenesisState {
	genesisAccounts := append(genesisValidatorAccounts, genesisWalletAccounts...)

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genesisAccounts)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for i, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			panic(err)
		}
		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			panic(err)
		}
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.OneInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genesisValidatorAccounts[i].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

		totalSupply = totalSupply.Add(sdk.NewCoin(chainCfg.BaseDenom, bondAmt))
	}

	// set validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = chainCfg.BaseDenom
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(chainCfg.BaseDenom, bondAmt.MulRaw(int64(len(validators))))},
	})

	// update total supply
	baseDenomDisplay := strings.ToUpper(chainCfg.BaseDenom[1:])
	denomMetadata := []banktypes.Metadata{
		{
			Description: "Base denom metadata",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    chainCfg.BaseDenom,
					Exponent: 0,
				},
				{
					Denom:    baseDenomDisplay,
					Exponent: 18,
				},
			},
			Base:    chainCfg.BaseDenom,
			Display: baseDenomDisplay,
			Name:    baseDenomDisplay,
			Symbol:  baseDenomDisplay,
		},
	}
	for _, secondaryDenomUnit := range testConfig.SecondaryDenomUnits {
		secondDenomDisplay := strings.ToUpper(secondaryDenomUnit.Denom[1:])
		denomMetadata = append(denomMetadata, banktypes.Metadata{
			Description: "Second denom metadata",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    secondaryDenomUnit.Denom,
					Exponent: 0,
				},
				{
					Denom:    secondDenomDisplay,
					Exponent: secondaryDenomUnit.Exponent,
				},
			},
			Base:    secondaryDenomUnit.Denom,
			Display: secondDenomDisplay,
			Name:    secondDenomDisplay,
			Symbol:  secondDenomDisplay,
		},
		)
	}

	{
		bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, denomMetadata, []banktypes.SendEnabled{})
		genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)
	}

	{
		// x/feemarket
		feeMarketGenesis := feemarkettypes.DefaultGenesisState()
		if feeMarketGenesis != nil {
			genesisState[feemarkettypes.ModuleName] = codec.MustMarshalJSON(feeMarketGenesis)
		}
	}

	{
		// x/evm
		var evmGenesis *evmtypes.GenesisState
		evmGenesis = evmtypes.DefaultGenesisState()
		if evmGenesis != nil {
			evmGenesis.Params.EvmDenom = chainCfg.BaseDenom
			evmGenesis.Params.EnableCall = true
			evmGenesis.Params.EnableCreate = false
			genesisState[evmtypes.ModuleName] = codec.MustMarshalJSON(evmGenesis)
		}
	}

	{
		// x/gov
		var govGenesis *govv1types.GenesisState
		govGenesis = govv1types.DefaultGenesisState()
		if govGenesis != nil {
			govGenesis.Params.MinDeposit = []sdk.Coin{sdk.NewCoin(chainCfg.BaseDenom, sdkmath.NewIntFromUint64(2))}
			var votingPeriod time.Duration
			if chainCfg.DisableTendermint {
				votingPeriod = 30 * time.Minute
			} else {
				// due to tendermint block time not configurable time jumping, we need to set a low voting period
				votingPeriod = TendermintGovVotingPeriod
			}
			govGenesis.Params.VotingPeriod = &votingPeriod
			genesisState[govtypes.ModuleName] = codec.MustMarshalJSON(govGenesis)
		}
	}

	{
		// x/slashing
		var slashingGenesis *slashingtypes.GenesisState
		slashingGenesis = slashingtypes.DefaultGenesisState()
		if slashingGenesis != nil {
			slashingGenesis.SigningInfos = signingInfos
			genesisState[slashingtypes.ModuleName] = codec.MustMarshalJSON(slashingGenesis)
		}
	}

	return genesisState
}

type disableCrisisInvariantsChecking struct{}

func (d disableCrisisInvariantsChecking) Get(o string) interface{} {
	if o == crisis.FlagSkipGenesisInvariants {
		return true
	}
	return nil
}
