package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k Keeper) GenericAddReverseLookupDymNamesRecord(ctx sdk.Context, key []byte, name string) error {
	var modifiedRecord dymnstypes.ReverseLookupDymNames

	modifiedRecord = dymnstypes.ReverseLookupDymNames{
		DymNames: []string{
			name,
		},
	}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz != nil {
		var existingRecord dymnstypes.ReverseLookupDymNames
		k.cdc.MustUnmarshal(bz, &existingRecord)

		modifiedRecord = existingRecord.Combine(
			modifiedRecord,
		)

		if len(modifiedRecord.DymNames) == len(existingRecord.DymNames) {
			// no new mapping to add
			return nil
		}
	}

	bz = k.cdc.MustMarshal(&modifiedRecord)
	store.Set(key, bz)

	return nil
}

func (k Keeper) GenericGetReverseLookupDymNamesRecord(
	ctx sdk.Context, key []byte,
) dymnstypes.ReverseLookupDymNames {
	var existingRecord dymnstypes.ReverseLookupDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &existingRecord)
	}

	return existingRecord
}

func (k Keeper) GenericRemoveReverseLookupDymNamesRecord(ctx sdk.Context, key []byte, name string) error {
	var existingRecord dymnstypes.ReverseLookupDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		// no mapping to remove
		return nil
	}

	k.cdc.MustUnmarshal(bz, &existingRecord)

	modifiedRecord := existingRecord.Exclude(dymnstypes.ReverseLookupDymNames{
		DymNames: []string{name},
	})

	if len(modifiedRecord.DymNames) == 0 {
		// no more, remove record
		store.Delete(key)
		return nil
	}

	bz = k.cdc.MustMarshal(&modifiedRecord)
	store.Set(key, bz)

	return nil
}
