package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// AddReverseMappingOwnerToOwnedDymName stores a reverse mapping from owner to owned Dym-Name into the KVStore.
func (k Keeper) AddReverseMappingOwnerToOwnedDymName(ctx sdk.Context, owner, name string) error {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(bzAccAddr)

	return k.GenericAddReverseLookupDymNamesRecord(ctx, dymNamesOwnedByAccountKey, name)
}

// GetDymNamesOwnedBy returns all Dym-Names owned by the account address.
func (k Keeper) GetDymNamesOwnedBy(
	ctx sdk.Context, owner string, nowEpoch int64,
) ([]dymnstypes.DymName, error) {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return nil, dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(bzAccAddr)

	existingOwnedDymNames := k.GenericGetReverseLookupDymNamesRecord(ctx, dymNamesOwnedByAccountKey)

	var dymNames []dymnstypes.DymName
	for _, owned := range existingOwnedDymNames.DymNames {
		dymName := k.GetDymNameWithExpirationCheck(ctx, owned, nowEpoch)
		if dymName == nil {
			// dym-name not found, skip
			continue
		}
		if dymName.Owner != owner {
			// dym-name owner mismatch, skip
			continue
		}
		dymNames = append(dymNames, *dymName)
	}

	return dymNames, nil
}

func (k Keeper) RemoveReverseMappingOwnerToOwnedDymName(ctx sdk.Context, owner, name string) error {
	accAddr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return dymnstypes.ErrInvalidOwner.Wrapf("owner `%s` is not a valid bech32 account address: %v", owner, err)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(accAddr)

	return k.GenericRemoveReverseLookupDymNamesRecord(ctx, dymNamesOwnedByAccountKey, name)
}

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
