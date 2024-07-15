package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var _ dymnstypes.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) dymnstypes.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (q queryServer) Params(goCtx context.Context, _ *dymnstypes.QueryParamsRequest) (*dymnstypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := q.GetParams(ctx)

	return &dymnstypes.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) DymName(goCtx context.Context, req *dymnstypes.QueryDymNameRequest) (*dymnstypes.QueryDymNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName := q.GetDymNameWithExpirationCheck(ctx, req.DymName, getEpochFromContextOrNow(ctx))

	return &dymnstypes.QueryDymNameResponse{DymName: dymName}, nil
}

func (q queryServer) ResolveDymNameAddresses(goCtx context.Context, req *dymnstypes.QueryResolveDymNameAddressesRequest) (*dymnstypes.QueryResolveDymNameAddressesResponse, error) {
	if req == nil || len(req.Addresses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// Legacy TODO DymNS: should be aware in Documents and UI implementation
	// there is a phishing vector like this: dym1.....@dym
	// With the current implementation, it is limited to 20 characters per name/sub-name
	// so it is easier to recognize the phishing vector: dym1234.5678@dym

	ctx := sdk.UnwrapSDKContext(goCtx)

	var result []dymnstypes.ResultDymNameAddress
	for _, address := range req.Addresses {
		resolvedAddress, err := q.ResolveByDymNameAddress(ctx, address)

		r := dymnstypes.ResultDymNameAddress{
			Address: address,
		}

		if err != nil {
			r.Error = err.Error()
		} else {
			r.ResolvedAddress = resolvedAddress
		}

		result = append(result, r)
	}

	return &dymnstypes.QueryResolveDymNameAddressesResponse{
		ResolvedAddresses: result,
	}, nil
}

func (q queryServer) DymNamesOwnedByAccount(goCtx context.Context, req *dymnstypes.QueryDymNamesOwnedByAccountRequest) (*dymnstypes.QueryDymNamesOwnedByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	dymNames, err := q.GetDymNamesOwnedBy(ctx, req.Owner, getEpochFromContextOrNow(ctx))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dymnstypes.QueryDymNamesOwnedByAccountResponse{
		DymNames: dymNames,
	}, nil
}

func (q queryServer) SellOrder(goCtx context.Context, req *dymnstypes.QuerySellOrderRequest) (*dymnstypes.QuerySellOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if !dymnsutils.IsValidDymName(req.DymName) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid dym name: %s", req.DymName)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	so := q.GetSellOrder(ctx, req.DymName)
	if so == nil {
		return nil, status.Errorf(codes.NotFound, "no active Sell Order for '%s' at this moment", req.DymName)
	}

	return &dymnstypes.QuerySellOrderResponse{
		Result: *so,
	}, nil
}

func (q queryServer) HistoricalSellOrder(goCtx context.Context, req *dymnstypes.QueryHistoricalSellOrderRequest) (*dymnstypes.QueryHistoricalSellOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if !dymnsutils.IsValidDymName(req.DymName) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid dym name: %s", req.DymName)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hso := q.GetHistoricalSellOrders(ctx, req.DymName)

	return &dymnstypes.QueryHistoricalSellOrderResponse{
		Result: hso,
	}, nil
}

func (q queryServer) EstimateRegisterName(goCtx context.Context, req *dymnstypes.QueryEstimateRegisterNameRequest) (*dymnstypes.QueryEstimateRegisterNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if !dymnsutils.IsValidDymName(req.Name) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid dym name: %s", req.Name)
	}

	if req.Duration < 1 {
		return nil, status.Error(codes.InvalidArgument, "duration must be at least 1 year")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)
	existingDymNameRecord := q.GetDymName(ctx, req.Name) // can be nil if not registered before

	if existingDymNameRecord != nil && existingDymNameRecord.Owner != req.Owner {
		// check take-over permission
		nowEpoch := getEpochFromContextOrNow(ctx)
		if !existingDymNameRecord.IsExpiredAtEpoch(nowEpoch) {
			return nil, status.Errorf(
				codes.PermissionDenied,
				"you are not the owner of '%s'", req.Name,
			)
		}

		// we ignore the grace period since this is just an estimation
	}

	estimation := EstimateRegisterName(
		params,
		req.Name,
		existingDymNameRecord,
		req.Owner,
		req.Duration,
	)
	return &estimation, nil
}

func getEpochFromContextOrNow(ctx sdk.Context) int64 {
	if !ctx.BlockTime().IsZero() {
		return ctx.BlockTime().Unix()
	}

	return time.Now().Unix()
}
