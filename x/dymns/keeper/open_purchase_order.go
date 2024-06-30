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

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixOpenPurchaseOrderByDymName)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var opo dymnstypes.OpenPurchaseOrder
		k.cdc.MustUnmarshal(iterator.Value(), &opo)
		list = append(list, opo)
	}

	return list
}
