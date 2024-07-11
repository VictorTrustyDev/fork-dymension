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

	if !opo.HasSetSellPrice() {
		opo.SellPrice = nil
	}

	// persist record
	store := ctx.KVStore(k.storeKey)
	opoKey := dymnstypes.OpenPurchaseOrderKey(opo.Name)
	bz := k.cdc.MustMarshal(&opo)
	store.Set(opoKey, bz)

	ctx.EventManager().EmitEvent(opo.GetSdkEvent(dymnstypes.AttributeKeyDymNameOpoActionNameSet))

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
func (k Keeper) DeleteOpenPurchaseOrder(ctx sdk.Context, dymName string) {
	opo := k.GetOpenPurchaseOrder(ctx, dymName)
	if opo == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	opoKey := dymnstypes.OpenPurchaseOrderKey(dymName)
	store.Delete(opoKey)

	ctx.EventManager().EmitEvent(opo.GetSdkEvent(dymnstypes.AttributeKeyDymNameOpoActionNameDelete))
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
	k.DeleteOpenPurchaseOrder(ctx, opo.Name)

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
		- This method will be called in an epoch hooks.
		- By skipping persisting the invalid historical record, we can prevent the chain from being halted.
		*/
	} else {
		// only persist if passed validation
		persist = true
	}

	if persist {
		k.SetHistoricalOpenPurchaseOrders(ctx, dymName, hOpo)

		var minExpiry int64 = -1
		for _, hOpo := range hOpo.OpenPurchaseOrders {
			if minExpiry < 0 || hOpo.ExpireAt < minExpiry {
				minExpiry = hOpo.ExpireAt
			}
		}
		if minExpiry > 0 {
			k.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, minExpiry)
		}
	}

	return nil
}

func (k Keeper) SetHistoricalOpenPurchaseOrders(ctx sdk.Context, dymName string, hOpo dymnstypes.HistoricalOpenPurchaseOrders) {
	store := ctx.KVStore(k.storeKey)
	hOpoKey := dymnstypes.HistoricalOpenPurchaseOrdersKey(dymName)
	bz := k.cdc.MustMarshal(&hOpo)
	store.Set(hOpoKey, bz)
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

// DeleteHistoricalOpenPurchaseOrders deletes the Historical Open-Purchase-Orders of specific Dym-Name from the KVStore.
func (k Keeper) DeleteHistoricalOpenPurchaseOrders(ctx sdk.Context, dymName string) {
	store := ctx.KVStore(k.storeKey)
	hOpoKey := dymnstypes.HistoricalOpenPurchaseOrdersKey(dymName)
	store.Delete(hOpoKey)
}

// CompletePurchaseOrder completes the purchase of the Dym-Name, give value to the previous owner, and transfer ownership.
func (k Keeper) CompletePurchaseOrder(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return dymnstypes.ErrDymNameNotFound.Wrap(name)
	}

	// here we don't check Dym-Name expiration, because it can not happen,
	// and there is a grace period for the owner to renew the Dym-Name in case bad things happen

	opo := k.GetOpenPurchaseOrder(ctx, name)
	if opo == nil {
		return dymnstypes.ErrOpenPurchaseOrderNotFound.Wrap(name)
	}

	if !opo.HasFinishedAtCtx(ctx) {
		return dymnstypes.ErrInvalidState.Wrap("Open-Purchase-Order has not finished yet")
	}

	// the OPO can be expired at this point,
	// in case the highest bid is lower than sell price or no sell price is set,
	// so the order is expired, but no logic to complete the purchase, then will be completed via hooks

	if opo.HighestBid == nil {
		return dymnstypes.ErrInvalidState.Wrap("no bid placed")
	}

	newOwner := opo.HighestBid.Bidder

	// complete the purchase

	previousOwner := dymName.Owner

	// give value to the previous owner
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		dymnstypes.ModuleName,
		sdk.MustAccAddressFromBech32(previousOwner),
		sdk.Coins{opo.HighestBid.Price},
	); err != nil {
		return err
	}

	// move the OPO to history
	if err := k.MoveOpenPurchaseOrderToHistorical(ctx, dymName.Name); err != nil {
		return err
	}

	// transfer ownership

	// remove reverse mapping
	if err := k.RemoveReverseMappingOwnerToOwnedDymName(ctx, previousOwner, dymName.Name); err != nil {
		return err
	}

	// update Dym records to prevent any potential mistake
	dymName.Owner = newOwner
	dymName.Controller = newOwner
	dymName.Configs = nil

	// persist updated DymName
	k.SetDymName(ctx, *dymName)

	return nil
}

func (k Keeper) SetActiveOpenPurchaseOrdersExpiration(ctx sdk.Context, opo dymnstypes.ActiveOpenPurchaseOrdersExpiration) error {
	// persist record
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&opo)
	store.Set(dymnstypes.KeyActiveOpenPurchaseOrdersExpiration, bz)
	return nil
}

func (k Keeper) GetActiveOpenPurchaseOrdersExpiration(ctx sdk.Context) dymnstypes.ActiveOpenPurchaseOrdersExpiration {
	store := ctx.KVStore(k.storeKey)

	var record dymnstypes.ActiveOpenPurchaseOrdersExpiration

	bz := store.Get(dymnstypes.KeyActiveOpenPurchaseOrdersExpiration)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &record)
	}

	if record.ExpiryByName == nil {
		record.ExpiryByName = make(map[string]int64)
	}

	return record
}

func (k Keeper) SetMinExpiryHistoricalOpenPurchaseOrder(ctx sdk.Context, dymName string, minExpiry int64) {
	store := ctx.KVStore(k.storeKey)
	key := dymnstypes.MinExpiryHistoricalOpenPurchaseOrdersKey(dymName)
	if minExpiry < 1 {
		store.Delete(key)
	} else {
		store.Set(key, sdk.Uint64ToBigEndian(uint64(minExpiry)))
	}
}

func (k Keeper) GetMinExpiryHistoricalOpenPurchaseOrder(ctx sdk.Context, dymName string) (minExpiry int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := dymnstypes.MinExpiryHistoricalOpenPurchaseOrdersKey(dymName)
	bz := store.Get(key)
	if bz != nil {
		minExpiry = int64(sdk.BigEndianToUint64(bz))
		found = true
	}
	return
}

func (k Keeper) GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx sdk.Context) (nameToMinExpiry map[string]int64) {
	store := ctx.KVStore(k.storeKey)

	nameToMinExpiry = make(map[string]int64)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixMinExpiryHistoricalOpenPurchaseOrders)
	defer func() {
		_ = iterator.Close() // nolint: errcheck
	}()

	for ; iterator.Valid(); iterator.Next() {
		dymName := string(iterator.Key()[len(dymnstypes.KeyPrefixMinExpiryHistoricalOpenPurchaseOrders):])
		minExpiry := int64(sdk.BigEndianToUint64(iterator.Value()))

		nameToMinExpiry[dymName] = minExpiry
	}

	return
}
