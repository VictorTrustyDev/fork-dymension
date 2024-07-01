package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// SetOpenPurchaseOrder stores an Open-Purchase-Order into the KVStore.
func (k Keeper) SetOpenPurchaseOrder(ctx sdk.Context, opo dymnstypes.OpenPurchaseOrder) error {
	if err := opo.Validate(); err != nil {
		return err
	}

	// persist record
	store := ctx.KVStore(k.storeKey)
	opoKey := dymnstypes.OpenPurchaseOrderKey(opo.Name)
	bz := k.cdc.MustMarshal(&opo)
	store.Set(opoKey, bz)
	ctx.EventManager().EmitEvent(opo.GetSdkEvent())

	return nil
}

// GetOpenPurchaseOrder retrieves Open-Purchase-Order of the corresponding Dym-Name from the KVStore.
// If the Open-Purchase-Order does not exist, nil is returned.
func (k Keeper) GetOpenPurchaseOrder(ctx sdk.Context, dymName string) *dymnstypes.OpenPurchaseOrder {
	store := ctx.KVStore(k.storeKey)
	opoKey := dymnstypes.OpenPurchaseOrderKey(dymName)

	bz := store.Get(opoKey)
	if bz == nil {
		return nil
	}

	var opo dymnstypes.OpenPurchaseOrder
	k.cdc.MustUnmarshal(bz, &opo)

	return &opo
}

// DeleteOpenPurchaseOrder deletes the Open-Purchase-Order from the KVStore.
func (k Keeper) DeleteOpenPurchaseOrder(ctx sdk.Context, opo dymnstypes.OpenPurchaseOrder) {
	store := ctx.KVStore(k.storeKey)
	opoKey := dymnstypes.OpenPurchaseOrderKey(opo.Name)
	store.Delete(opoKey)
}

// GetAllOpenPurchaseOrders returns all Open-Purchase-Orders from the KVStore.
// No filter is applied.
func (k Keeper) GetAllOpenPurchaseOrders(ctx sdk.Context) (list []dymnstypes.OpenPurchaseOrder) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixOpenPurchaseOrder)
	defer func() {
		_ = iterator.Close() // nolint: errcheck
	}()

	for ; iterator.Valid(); iterator.Next() {
		var opo dymnstypes.OpenPurchaseOrder
		k.cdc.MustUnmarshal(iterator.Value(), &opo)
		list = append(list, opo)
	}

	return list
}

// MoveOpenPurchaseOrderToHistorical moves the active Open-Purchase-Order record of the Dym-Name
// into historical, and deletes the original record from KVStore.
func (k Keeper) MoveOpenPurchaseOrderToHistorical(ctx sdk.Context, dymName string) error {
	// find active record
	opo := k.GetOpenPurchaseOrder(ctx, dymName)
	if opo == nil {
		return dymnstypes.ErrOpenPurchaseOrderNotFound.Wrap(dymName)
	}

	if opo.HighestBid == nil {
		// in-case of no bid, check if the order has expired
		if !opo.HasExpiredAtCtx(ctx) {
			return dymnstypes.ErrInvalidState.Wrapf(
				"Open-Purchase-Order of '%s' has not expired yet",
				dymName,
			)
		}
	}

	// remove the active record
	k.DeleteOpenPurchaseOrder(ctx, *opo)

	// set historical records
	store := ctx.KVStore(k.storeKey)
	hOpoKey := dymnstypes.HistoricalOpenPurchaseOrdersKey(dymName)
	bz := store.Get(hOpoKey)

	var hOpo dymnstypes.HistoricalOpenPurchaseOrders
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &hOpo)
	}
	hOpo.OpenPurchaseOrders = append(hOpo.OpenPurchaseOrders, *opo)

	var persist bool

	if ignorableErr := hOpo.Validate(); ignorableErr != nil {
		k.Logger(ctx).Error(
			"historical open purchase order validation failed, skip persist this historical record",
			"error", ignorableErr,
		)

		// skip persisting historical record

		/**
		Why do we skip persisting the historical record when it fails validation?
		- The historical record is not an important data for the chain to function.
		- By skipping persisting the invalid historical record, we can prevent the chain from being halted.
		*/
	} else {
		// only persist if passed validation
		persist = true
	}

	if persist {
		bz = k.cdc.MustMarshal(&hOpo)
		store.Set(hOpoKey, bz)
	}

	return nil
}

// GetHistoricalOpenPurchaseOrders retrieves Historical Open-Purchase-Orders of the corresponding Dym-Name from the KVStore.
func (k Keeper) GetHistoricalOpenPurchaseOrders(ctx sdk.Context, dymName string) []dymnstypes.OpenPurchaseOrder {
	store := ctx.KVStore(k.storeKey)
	hOpoKey := dymnstypes.HistoricalOpenPurchaseOrdersKey(dymName)

	bz := store.Get(hOpoKey)
	if bz == nil {
		return nil
	}

	var hOpo dymnstypes.HistoricalOpenPurchaseOrders
	k.cdc.MustUnmarshal(bz, &hOpo)

	return hOpo.OpenPurchaseOrders
}
