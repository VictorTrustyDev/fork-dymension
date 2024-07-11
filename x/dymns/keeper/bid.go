package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// GenesisRefundBid refunds the bid in genesis initialization.
// This action will mint coins to the module account and send coins to the bidder.
func (k Keeper) GenesisRefundBid(ctx sdk.Context, opoBid dymnstypes.OpenPurchaseOrderBid) error {
	return k.refundBid(ctx, opoBid, true)
}

// RefundBid refunds the bid.
// This action will send coins from module account to the bidder.
func (k Keeper) RefundBid(ctx sdk.Context, opoBid dymnstypes.OpenPurchaseOrderBid) error {
	return k.refundBid(ctx, opoBid, false)
}

func (k Keeper) refundBid(ctx sdk.Context, opoBid dymnstypes.OpenPurchaseOrderBid, genesis bool) error {
	if err := opoBid.Validate(); err != nil {
		return err
	}

	if genesis {
		// During genesis initialization progress, the module account has no balance, so we mint coins.
		// Otherwise, the module account should have enough balance to refund the bid.
		if err := k.bankKeeper.MintCoins(ctx, dymnstypes.ModuleName, sdk.Coins{opoBid.Price}); err != nil {
			return err
		}
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		dymnstypes.ModuleName,
		sdk.MustAccAddressFromBech32(opoBid.Bidder),
		sdk.Coins{opoBid.Price},
	); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			dymnstypes.EventTypeDymNameRefundBid,
			sdk.NewAttribute(dymnstypes.AttributeKeyDymNameRefundBidder, opoBid.Bidder),
			sdk.NewAttribute(dymnstypes.AttributeKeyDymNameRefundAmount, opoBid.Price.String()),
		),
	)

	return nil
}
