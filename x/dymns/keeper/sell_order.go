package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// SetSellOrder stores a Sell-Order into the KVStore.
func (k Keeper) SetSellOrder(ctx sdk.Context, so dymnstypes.SellOrder) error {
	if err := so.Validate(); err != nil {
		return err
	}

	if !so.HasSetSellPrice() {
		so.SellPrice = nil
	}

	// persist record
	store := ctx.KVStore(k.storeKey)
	soKey := dymnstypes.SellOrderKey(so.Name)
	bz := k.cdc.MustMarshal(&so)
	store.Set(soKey, bz)

	ctx.EventManager().EmitEvent(so.GetSdkEvent(dymnstypes.AttributeKeyDymNameSoActionNameSet))

	return nil
}

// GetSellOrder retrieves active Sell-Order of the corresponding Dym-Name from the KVStore.
// If the Sell-Order does not exist, nil is returned.
func (k Keeper) GetSellOrder(ctx sdk.Context, dymName string) *dymnstypes.SellOrder {
	store := ctx.KVStore(k.storeKey)
	soKey := dymnstypes.SellOrderKey(dymName)

	bz := store.Get(soKey)
	if bz == nil {
		return nil
	}

	var so dymnstypes.SellOrder
	k.cdc.MustUnmarshal(bz, &so)

	return &so
}

// DeleteSellOrder deletes the Sell-Order from the KVStore.
func (k Keeper) DeleteSellOrder(ctx sdk.Context, dymName string) {
	so := k.GetSellOrder(ctx, dymName)
	if so == nil {
		return
	}

	store := ctx.KVStore(k.storeKey)
	soKey := dymnstypes.SellOrderKey(dymName)
	store.Delete(soKey)

	ctx.EventManager().EmitEvent(so.GetSdkEvent(dymnstypes.AttributeKeyDymNameSoActionNameDelete))
}

// GetAllSellOrders returns all active Sell-Orders from the KVStore.
// No filter is applied.
func (k Keeper) GetAllSellOrders(ctx sdk.Context) (list []dymnstypes.SellOrder) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixSellOrder)
	defer func() {
		_ = iterator.Close() // nolint: errcheck
	}()

	for ; iterator.Valid(); iterator.Next() {
		var so dymnstypes.SellOrder
		k.cdc.MustUnmarshal(iterator.Value(), &so)
		list = append(list, so)
	}

	return list
}

// MoveSellOrderToHistorical moves the active Sell-Order record of the Dym-Name
// into historical, and deletes the original record from KVStore.
func (k Keeper) MoveSellOrderToHistorical(ctx sdk.Context, dymName string) error {
	// find active record
	so := k.GetSellOrder(ctx, dymName)
	if so == nil {
		return dymnstypes.ErrSellOrderNotFound.Wrap(dymName)
	}

	if so.HighestBid == nil {
		// in-case of no bid, check if the order has expired
		if !so.HasExpiredAtCtx(ctx) {
			return dymnstypes.ErrInvalidState.Wrapf(
				"Sell-Order of '%s' has not expired yet",
				dymName,
			)
		}
	}

	// remove the active record
	k.DeleteSellOrder(ctx, so.Name)

	// set historical records
	store := ctx.KVStore(k.storeKey)
	hSoKey := dymnstypes.HistoricalSellOrdersKey(dymName)
	bz := store.Get(hSoKey)

	var hSo dymnstypes.HistoricalSellOrders
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &hSo)
	}
	hSo.SellOrders = append(hSo.SellOrders, *so)

	var persist bool

	if ignorableErr := hSo.Validate(); ignorableErr != nil {
		k.Logger(ctx).Error(
			"historical sell order validation failed, skip persist this historical record",
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
		k.SetHistoricalSellOrders(ctx, dymName, hSo)

		var minExpiry int64 = -1
		for _, hSo := range hSo.SellOrders {
			if minExpiry < 0 || hSo.ExpireAt < minExpiry {
				minExpiry = hSo.ExpireAt
			}
		}
		if minExpiry > 0 {
			k.SetMinExpiryHistoricalSellOrder(ctx, dymName, minExpiry)
		}
	}

	return nil
}

func (k Keeper) SetHistoricalSellOrders(ctx sdk.Context, dymName string, hSo dymnstypes.HistoricalSellOrders) {
	store := ctx.KVStore(k.storeKey)
	hSoKey := dymnstypes.HistoricalSellOrdersKey(dymName)
	bz := k.cdc.MustMarshal(&hSo)
	store.Set(hSoKey, bz)
}

// GetHistoricalSellOrders retrieves Historical Sell-Orders of the corresponding Dym-Name from the KVStore.
func (k Keeper) GetHistoricalSellOrders(ctx sdk.Context, dymName string) []dymnstypes.SellOrder {
	store := ctx.KVStore(k.storeKey)
	hSoKey := dymnstypes.HistoricalSellOrdersKey(dymName)

	bz := store.Get(hSoKey)
	if bz == nil {
		return nil
	}

	var hSo dymnstypes.HistoricalSellOrders
	k.cdc.MustUnmarshal(bz, &hSo)

	return hSo.SellOrders
}

// DeleteHistoricalSellOrders deletes the Historical Sell-Orders of specific Dym-Name from the KVStore.
func (k Keeper) DeleteHistoricalSellOrders(ctx sdk.Context, dymName string) {
	store := ctx.KVStore(k.storeKey)
	hSoKey := dymnstypes.HistoricalSellOrdersKey(dymName)
	store.Delete(hSoKey)
}

// CompleteSellOrder completes the active sell order of the Dym-Name,
// give value to the previous owner, and transfer ownership to new owner.
func (k Keeper) CompleteSellOrder(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return dymnstypes.ErrDymNameNotFound.Wrap(name)
	}

	// here we don't check Dym-Name expiration, because it can not happen,
	// and there is a grace period for the owner to renew the Dym-Name in case bad things happen

	so := k.GetSellOrder(ctx, name)
	if so == nil {
		return dymnstypes.ErrSellOrderNotFound.Wrap(name)
	}

	if !so.HasFinishedAtCtx(ctx) {
		return dymnstypes.ErrInvalidState.Wrap("Sell-Order has not finished yet")
	}

	// the SO can be expired at this point,
	// in case the highest bid is lower than sell price or no sell price is set,
	// so the order is expired, but no logic to complete the SO, then will be completed via hooks

	if so.HighestBid == nil {
		return dymnstypes.ErrInvalidState.Wrap("no bid placed")
	}

	newOwner := so.HighestBid.Bidder

	// complete the Sell

	previousOwner := dymName.Owner

	// give value to the previous owner
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		dymnstypes.ModuleName,
		sdk.MustAccAddressFromBech32(previousOwner),
		sdk.Coins{so.HighestBid.Price},
	); err != nil {
		return err
	}

	// move the SO to history
	if err := k.MoveSellOrderToHistorical(ctx, dymName.Name); err != nil {
		return err
	}

	// transfer ownership

	if err := k.BeforeDymNameOwnerChanged(ctx, dymName.Name); err != nil {
		return err
	}

	if err := k.BeforeDymNameConfigChanged(ctx, dymName.Name); err != nil {
		return err
	}

	// update Dym records to prevent any potential mistake
	dymName.Owner = newOwner
	dymName.Controller = newOwner
	dymName.Configs = nil

	// persist updated DymName
	if err := k.SetDymName(ctx, *dymName); err != nil {
		return err
	}

	if err := k.AfterDymNameOwnerChanged(ctx, dymName.Name); err != nil {
		return err
	}

	if err := k.AfterDymNameConfigChanged(ctx, dymName.Name); err != nil {
		return err
	}

	return nil
}

func (k Keeper) SetActiveSellOrdersExpiration(ctx sdk.Context, so dymnstypes.ActiveSellOrdersExpiration) error {
	// persist record
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&so)
	store.Set(dymnstypes.KeyActiveSellOrdersExpiration, bz)
	return nil
}

func (k Keeper) GetActiveSellOrdersExpiration(ctx sdk.Context) dymnstypes.ActiveSellOrdersExpiration {
	store := ctx.KVStore(k.storeKey)

	var record dymnstypes.ActiveSellOrdersExpiration

	bz := store.Get(dymnstypes.KeyActiveSellOrdersExpiration)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &record)
	}

	if record.ExpiryByName == nil {
		record.ExpiryByName = make(map[string]int64)
	}

	return record
}

func (k Keeper) SetMinExpiryHistoricalSellOrder(ctx sdk.Context, dymName string, minExpiry int64) {
	store := ctx.KVStore(k.storeKey)
	key := dymnstypes.MinExpiryHistoricalSellOrdersKey(dymName)
	if minExpiry < 1 {
		store.Delete(key)
	} else {
		store.Set(key, sdk.Uint64ToBigEndian(uint64(minExpiry)))
	}
}

func (k Keeper) GetMinExpiryHistoricalSellOrder(ctx sdk.Context, dymName string) (minExpiry int64, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := dymnstypes.MinExpiryHistoricalSellOrdersKey(dymName)
	bz := store.Get(key)
	if bz != nil {
		minExpiry = int64(sdk.BigEndianToUint64(bz))
		found = true
	}
	return
}

func (k Keeper) GetMinExpiryOfAllHistoricalSellOrders(ctx sdk.Context) (nameToMinExpiry map[string]int64) {
	store := ctx.KVStore(k.storeKey)

	nameToMinExpiry = make(map[string]int64)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixMinExpiryHistoricalSellOrders)
	defer func() {
		_ = iterator.Close() // nolint: errcheck
	}()

	for ; iterator.Valid(); iterator.Next() {
		dymName := string(iterator.Key()[len(dymnstypes.KeyPrefixMinExpiryHistoricalSellOrders):])
		minExpiry := int64(sdk.BigEndianToUint64(iterator.Value()))

		nameToMinExpiry[dymName] = minExpiry
	}

	return
}