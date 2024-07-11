package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k msgServer) PurchaseName(goCtx context.Context, msg *dymnstypes.MsgPurchaseName) (*dymnstypes.MsgPurchaseNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName, opo, err := k.validatePurchase(ctx, msg)
	if err != nil {
		return nil, err
	}

	if opo.HighestBid != nil {
		// refund previous bidder
		if err := k.RefundBid(ctx, *opo.HighestBid); err != nil {
			return nil, err
		}
	}

	// deduct offer price from buyer's account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		sdk.MustAccAddressFromBech32(msg.Buyer),
		dymnstypes.ModuleName,
		sdk.Coins{msg.Offer},
	); err != nil {
		return nil, err
	}

	// record new highest bid
	opo.HighestBid = &dymnstypes.OpenPurchaseOrderBid{
		Bidder: msg.Buyer,
		Price:  msg.Offer,
	}

	// store OPO after highest bid updated
	if err := k.SetOpenPurchaseOrder(ctx, *opo); err != nil {
		return nil, err
	}

	// try to complete the purchase

	if opo.HasFinishedAtCtx(ctx) {
		if err := k.CompletePurchaseOrder(ctx, dymName.Name); err != nil {
			return nil, err
		}
	}

	miscParams := k.MiscParams(ctx)
	minimumTxGas := sdk.Gas(miscParams.GasCrudOpenPurchaseOrder)
	if consumedGas := ctx.GasMeter().GasConsumed(); consumedGas < minimumTxGas {
		ctx.GasMeter().ConsumeGas(minimumTxGas-consumedGas, "PurchaseName")
	}

	return &dymnstypes.MsgPurchaseNameResponse{}, nil
}

func (k msgServer) validatePurchase(ctx sdk.Context, msg *dymnstypes.MsgPurchaseName) (*dymnstypes.DymName, *dymnstypes.OpenPurchaseOrder, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return nil, nil, dymnstypes.ErrDymNameNotFound.Wrap(msg.Name)
	}

	if dymName.Owner == msg.Buyer {
		return nil, nil, sdkerrors.ErrLogic.Wrap("cannot purchase your own dym name")
	}

	opo := k.GetOpenPurchaseOrder(ctx, msg.Name)
	if opo == nil {
		return nil, nil, dymnstypes.ErrOpenPurchaseOrderNotFound.Wrap(msg.Name)
	}

	if opo.HasExpiredAtCtx(ctx) {
		return nil, nil, dymnstypes.ErrInvalidState.Wrap("cannot purchase an expired order")
	}

	if opo.HasFinishedAtCtx(ctx) {
		return nil, nil, dymnstypes.ErrInvalidState.Wrap("cannot purchase a completed order")
	}

	if msg.Offer.Denom != opo.MinPrice.Denom {
		return nil, nil, sdkerrors.ErrUnknownRequest.Wrapf(
			"offer denom does not match the order denom: %s != %s",
			msg.Offer.Denom, opo.MinPrice.Denom,
		)
	}

	if msg.Offer.IsLT(opo.MinPrice) {
		return nil, nil, sdkerrors.ErrInsufficientFunds.Wrap("offer is lower than minimum price")
	}

	if opo.HasSetSellPrice() {
		if !msg.Offer.IsLTE(*opo.SellPrice) { // overpaid protection
			return nil, nil, sdkerrors.ErrInsufficientFunds.Wrap("offer is higher than sell price")
		}
	}

	if opo.HighestBid != nil {
		if msg.Offer.IsLTE(opo.HighestBid.Price) {
			return nil, nil, sdkerrors.ErrInsufficientFunds.Wrap("new offer must be higher than current highest bid")
		}
	}

	return dymName, opo, nil
}
