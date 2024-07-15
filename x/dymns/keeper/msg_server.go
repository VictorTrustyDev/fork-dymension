package keeper

import (
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

var _ dymnstypes.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) dymnstypes.MsgServer {
	return &msgServer{Keeper: keeper}
}
