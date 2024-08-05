package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ dymnstypes.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) dymnstypes.QueryServer {
	return &queryServer{Keeper: keeper}
}

// Params queries the parameters of the module.
func (q queryServer) Params(goCtx context.Context, _ *dymnstypes.QueryParamsRequest) (*dymnstypes.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := q.GetParams(ctx)

	return &dymnstypes.QueryParamsResponse{Params: params}, nil
}

// DymName queries a Dym-Name by its name.
func (q queryServer) DymName(goCtx context.Context, req *dymnstypes.QueryDymNameRequest) (*dymnstypes.QueryDymNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName := q.GetDymNameWithExpirationCheck(ctx, req.DymName)

	return &dymnstypes.QueryDymNameResponse{DymName: dymName}, nil
}

// ResolveDymNameAddresses resolves multiple Dym-Name Addresses to account address of each pointing to.
//
// For example:
//   - "my-name@dym" => "dym1a..."
//   - "another.my-name@dym" => "dym1b..."
//   - "my-name@nim" => "nim1..."
//   - (extra format) "0x1234...6789@nim" => "nim1..."
//   - (extra format) "dym1a...@nim" => "nim1..."
func (q queryServer) ResolveDymNameAddresses(goCtx context.Context, req *dymnstypes.QueryResolveDymNameAddressesRequest) (*dymnstypes.QueryResolveDymNameAddressesResponse, error) {
	if req == nil || len(req.Addresses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// There is a phishing attack vector like this: dym1.....@dym
	// With the current implementation, it is limited to 20 characters per name/sub-name
	// so, it is easier to recognize: dym1234.5678@dym

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

// DymNamesOwnedByAccount queries the Dym-Names owned by an account.
func (q queryServer) DymNamesOwnedByAccount(goCtx context.Context, req *dymnstypes.QueryDymNamesOwnedByAccountRequest) (*dymnstypes.QueryDymNamesOwnedByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	dymNames, err := q.GetDymNamesOwnedBy(ctx, req.Owner)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dymnstypes.QueryDymNamesOwnedByAccountResponse{
		DymNames: dymNames,
	}, nil
}

// SellOrder queries the active SO of a Dym-Name.
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

// HistoricalSellOrder queries the historical SOs of a Dym-Name.
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

// EstimateRegisterName estimates the cost to register a Dym-Name.
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
		if !existingDymNameRecord.IsExpiredAtCtx(ctx) {
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

// ReverseResolveAddress resolves multiple account addresses to Dym-Name Addresses which point to each.
// This function may return multiple possible Dym-Name-Addresses those point to each of the input address.
//
// For example: when we have "my-name@dym" resolves to "dym1a..."
// so reverse resolve will return "my-name@dym" when input is "dym1a..."
func (q queryServer) ReverseResolveAddress(goCtx context.Context, req *dymnstypes.QueryReverseResolveAddressRequest) (*dymnstypes.QueryReverseResolveAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.Addresses) < 1 {
		return nil, status.Error(codes.InvalidArgument, "no addresses provided")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	workingChainId := req.WorkingChainId
	if workingChainId == "" {
		workingChainId = ctx.ChainID()
	}

	result := make(map[string]dymnstypes.ReverseResolveAddressResult)
	// Describe usage of Go Map: non-consensus state, for querying purpose only.

	addErrorResult := func(address string, err error) {
		result[address] = dymnstypes.ReverseResolveAddressResult{
			Error: err.Error(),
		}
	}

	addResult := func(address string, candidates []dymnstypes.ReverseResolvedDymNameAddress) {
		dymNameAddress := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			dymNameAddress = append(dymNameAddress, candidate.String())
		}

		result[address] = dymnstypes.ReverseResolveAddressResult{
			Candidates: dymNameAddress,
		}
	}

	for _, address := range req.Addresses {
		if !dymnsutils.PossibleAccountRegardlessChain(address) {
			// Simply ignore invalid address.
			// Invalid address is not included here to prevent wasting resources due to bad requests.
			continue
		}

		candidates, err := q.ReverseResolveDymNameAddress(ctx, address, workingChainId)
		if err != nil {
			addErrorResult(address, err)
			continue
		}

		addResult(address, candidates)
	}

	return &dymnstypes.QueryReverseResolveAddressResponse{
		Result:         result,
		WorkingChainId: workingChainId,
	}, nil
}

// TranslateAliasOrChainIdToChainId tries to translate an alias/handle to a chain id.
// If an alias/handle can not be translated to chain-id, it is treated as a chain-id and returns.
func (q queryServer) TranslateAliasOrChainIdToChainId(goCtx context.Context, req *dymnstypes.QueryTranslateAliasOrChainIdToChainIdRequest) (*dymnstypes.QueryTranslateAliasOrChainIdToChainIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.AliasOrChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "empty alias or chain id")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	resolvedToChainId, success := q.tryResolveChainIdOrAliasToChainId(ctx, req.AliasOrChainId)
	if !success {
		resolvedToChainId = req.AliasOrChainId
	}

	return &dymnstypes.QueryTranslateAliasOrChainIdToChainIdResponse{
		ChainId: resolvedToChainId,
	}, nil
}

// BuyOfferById queries a buy offer by its id.
func (q queryServer) BuyOfferById(goCtx context.Context, req *dymnstypes.QueryBuyOfferByIdRequest) (*dymnstypes.QueryBuyOfferByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if !dymnstypes.IsValidBuyOfferId(req.Id) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid offer id: %s", req.Id)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	offer := q.GetBuyOffer(ctx, req.Id)
	if offer == nil {
		return nil, status.Error(codes.NotFound, "offer not found")
	}

	return &dymnstypes.QueryBuyOfferByIdResponse{
		Offer: *offer,
	}, nil
}

// BuyOffersPlacedByAccount queries the all the buy offers placed by an account.
func (q queryServer) BuyOffersPlacedByAccount(goCtx context.Context, req *dymnstypes.QueryBuyOffersPlacedByAccountRequest) (*dymnstypes.QueryBuyOffersPlacedByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	offers, err := q.GetBuyOffersByBuyer(ctx, req.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dymnstypes.QueryBuyOffersPlacedByAccountResponse{
		Offers: offers,
	}, nil
}

// BuyOffersByDymName queries all the buy offers of a Dym-Name.
func (q queryServer) BuyOffersByDymName(goCtx context.Context, req *dymnstypes.QueryBuyOffersByDymNameRequest) (*dymnstypes.QueryBuyOffersByDymNameResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	offers, err := q.GetBuyOffersOfDymName(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dymnstypes.QueryBuyOffersByDymNameResponse{
		Offers: offers,
	}, nil
}

// BuyOffersOfDymNamesOwnedByAccount queries all the buy offers of all Dym-Names owned by an account.
func (q queryServer) BuyOffersOfDymNamesOwnedByAccount(goCtx context.Context, req *dymnstypes.QueryBuyOffersOfDymNamesOwnedByAccountRequest) (*dymnstypes.QueryBuyOffersOfDymNamesOwnedByAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	ownedDymNames, err := q.GetDymNamesOwnedBy(ctx, req.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	offers := make([]dymnstypes.BuyOffer, 0)
	for _, dymName := range ownedDymNames {
		offersOfDymName, err := q.GetBuyOffersOfDymName(ctx, dymName.Name)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		offers = append(offers, offersOfDymName...)
	}

	return &dymnstypes.QueryBuyOffersOfDymNamesOwnedByAccountResponse{
		Offers: offers,
	}, nil
}
