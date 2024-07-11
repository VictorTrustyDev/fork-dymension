package types

//goland:noinspection SpellCheckingInspection
import (
	httpclient "github.com/cometbft/cometbft/rpc/client/http"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	cosmostxtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypeslegacy "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	rollapptypes "github.com/dymensionxyz/dymension/v3/x/rollapp/types"
	sequencertypes "github.com/dymensionxyz/dymension/v3/x/sequencer/types"
	streamertypes "github.com/dymensionxyz/dymension/v3/x/streamer/types"
	rpctypes "github.com/evmos/ethermint/rpc/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type QueryClients struct {
	GrpcConnection          grpc1.ClientConn
	ClientQueryCtx          cosmosclient.Context
	TendermintRpcHttpClient *httpclient.HTTP
	Auth                    authtypes.QueryClient
	Bank                    banktypes.QueryClient
	Distribution            distributiontypes.QueryClient
	EVM                     evmtypes.QueryClient
	GovV1                   govtypesv1.QueryClient
	GovLegacy               govtypeslegacy.QueryClient
	IbcTransfer             ibctransfertypes.QueryClient
	RollApp                 rollapptypes.QueryClient
	Sequencer               sequencertypes.QueryClient
	Slashing                slashingtypes.QueryClient
	Staking                 stakingtypes.QueryClient
	Streamer                streamertypes.QueryClient
	ServiceClient           cosmostxtypes.ServiceClient
	Rpc                     *rpctypes.QueryClient
}
