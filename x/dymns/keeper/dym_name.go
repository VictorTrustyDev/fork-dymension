package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// SetDymName stores a Dym-Name into the KVStore.
func (k Keeper) SetDymName(ctx sdk.Context, dymName dymnstypes.DymName) error {
	if err := dymName.Validate(); err != nil {
		return err
	}

	// persist record
	store := ctx.KVStore(k.storeKey)
	dymNameKey := dymnstypes.DymNameKey(dymName.Name)
	bz := k.cdc.MustMarshal(&dymName)
	store.Set(dymNameKey, bz)
	ctx.EventManager().EmitEvent(dymName.GetSdkEvent())

	// finally persist reverse lookup by owner account address
	return k.SetReverseMappingOwnerToOwnedDymName(ctx, dymName.Owner, dymName.Name)
}

// GetDymName returns a Dym-Name from the KVStore.
func (k Keeper) GetDymName(ctx sdk.Context, name string) *dymnstypes.DymName {
	store := ctx.KVStore(k.storeKey)
	dymNameKey := dymnstypes.DymNameKey(name)

	bz := store.Get(dymNameKey)
	if bz == nil {
		return nil
	}

	var dymName dymnstypes.DymName
	k.cdc.MustUnmarshal(bz, &dymName)

	return &dymName
}

func (k Keeper) GetAllNonExpiredDymNames(ctx sdk.Context, anchorEpoch int64) (list []dymnstypes.DymName) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixDymName)
	defer iterator.Close() // nolint: errcheck

	for ; iterator.Valid(); iterator.Next() {
		var dymName dymnstypes.DymName
		k.cdc.MustUnmarshal(iterator.Value(), &dymName)

		if dymName.ExpireAt < anchorEpoch {
			continue
		}
		list = append(list, dymName)
	}

	return list
}

func (k Keeper) SetReverseMappingOwnerToOwnedDymName(ctx sdk.Context, owner, name string) error {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountKey(bzAccAddr)

	var existingOwnedDymNames dymnstypes.OwnedDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(dymNamesOwnedByAccountKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &existingOwnedDymNames)
		for _, owned := range existingOwnedDymNames.DymNames {
			if owned == name {
				return nil
			}
		}

		existingOwnedDymNames.DymNames = append(existingOwnedDymNames.DymNames, name)
	} else {
		existingOwnedDymNames = dymnstypes.OwnedDymNames{
			DymNames: []string{
				name,
			},
		}
	}

	bz = k.cdc.MustMarshal(&existingOwnedDymNames)
	store.Set(dymNamesOwnedByAccountKey, bz)

	return nil
}

func (k Keeper) GetDymNamesOwnedBy(ctx sdk.Context, owner string) ([]dymnstypes.DymName, error) {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return nil, dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountKey(bzAccAddr)

	var existingOwnedDymNames dymnstypes.OwnedDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(dymNamesOwnedByAccountKey)
	if bz == nil {
		return nil, nil
	}

	k.cdc.MustUnmarshal(bz, &existingOwnedDymNames)

	var dymNames []dymnstypes.DymName
	for _, owned := range existingOwnedDymNames.DymNames {
		dymName := k.GetDymName(ctx, owned)
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
	// TODO DymNS: call this function when a dym-name is expired
	// TODO DymNS: call this function when a dym-name ownership is transferred

	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountKey(bzAccAddr)

	var existingOwnedDymNames dymnstypes.OwnedDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(dymNamesOwnedByAccountKey)
	if bz == nil {
		// no mapping to remove
		return nil
	}

	k.cdc.MustUnmarshal(bz, &existingOwnedDymNames)

	var newOwnedDymNames []string
	for _, owned := range existingOwnedDymNames.DymNames {
		if owned == name {
			continue
		}
		newOwnedDymNames = append(newOwnedDymNames, owned)
	}
	if len(newOwnedDymNames) == len(existingOwnedDymNames.DymNames) {
		// no mapping to remove
		return nil
	}

	if len(newOwnedDymNames) == 0 {
		// no more owned dym-names, remove the mapping
		store.Delete(dymNamesOwnedByAccountKey)
		return nil
	}

	existingOwnedDymNames.DymNames = newOwnedDymNames
	bz = k.cdc.MustMarshal(&existingOwnedDymNames)
	store.Set(dymNamesOwnedByAccountKey, bz)

	return nil
}
