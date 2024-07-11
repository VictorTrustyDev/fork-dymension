package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"strings"
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

// GetDymNameWithExpirationCheck returns a Dym-Name from the KVStore, if the Dym-Name is not expired.
// Returns nil if Dym-Name does not exist or is expired.
func (k Keeper) GetDymNameWithExpirationCheck(ctx sdk.Context, name string, nowEpoch int64) *dymnstypes.DymName {
	// Legacy TODO DymNS: always use this on queries
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return nil
	}

	if dymName.ExpireAt < nowEpoch {
		return nil
	}

	return dymName
}

func (k Keeper) DeleteDymName(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	dymNameKey := dymnstypes.DymNameKey(name)
	store.Delete(dymNameKey)
}

func (k Keeper) GetAllNonExpiredDymNames(ctx sdk.Context, nowEpoch int64) (list []dymnstypes.DymName) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, dymnstypes.KeyPrefixDymName)
	defer func() {
		_ = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var dymName dymnstypes.DymName
		k.cdc.MustUnmarshal(iterator.Value(), &dymName)

		if dymName.ExpireAt < nowEpoch {
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

func (k Keeper) GetDymNamesOwnedBy(
	ctx sdk.Context, owner string, nowEpoch int64,
) ([]dymnstypes.DymName, error) {
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

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountKey(accAddr)

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

// PruneDymName removes a Dym-Name from the KVStore, as well as all related records.
func (k Keeper) PruneDymName(ctx sdk.Context, name string) error {
	// remove OPO (force, ignore active OPO)
	k.DeleteOpenPurchaseOrder(ctx, name)

	// remove historical OPO
	k.DeleteHistoricalOpenPurchaseOrders(ctx, name)
	k.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, name, 0)

	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return nil
	}

	// remove reverse mapping
	if err := k.RemoveReverseMappingOwnerToOwnedDymName(ctx, dymName.Owner, dymName.Name); err != nil {
		return err
	}

	// remove config
	// This seems not necessary because we are going to remove the record anyway,
	// but just let it here to clear the business logic
	dymName.Configs = nil   // all configuration should be removed
	dymName.Owner = ""      // no one owns it anyone
	dymName.Controller = "" // no one controls it anyone

	// remove record
	k.DeleteDymName(ctx, name)

	return nil
}

func (k Keeper) ResolveByDymNameAddress(ctx sdk.Context, dymNameAddress string) (outputAddress string, err error) {
	subNameParts, name, chainIdOrAlias, parseErr := ParseDymNameAddress(dymNameAddress)
	if parseErr != nil {
		ctx.Logger().Debug("failed to parse Dym-Name address", "dym-name-address", dymNameAddress, "error", parseErr)
		err = parseErr
		return
	}

	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		err = dymnstypes.ErrDymNameNotFound.Wrap(name)
		return
	}
	if dymName.IsExpiredAt(ctx.BlockTime()) {
		err = dymnstypes.ErrDymNameNotFound.Wrap(name)
		return
	}

	defer func() {
		if outputAddress == "" {
			err = sdkerrors.ErrInvalidRequest.Wrap("no resolution found")
		}
	}()

	var resolvedChainId string
	if chainIdOrAlias == ctx.ChainID() {
		resolvedChainId = chainIdOrAlias
	} else {
		aliasParams := k.AliasParams(ctx)
		if len(aliasParams.ByChainId) > 0 {
		loopFindChainId:
			for chainId, aliases := range aliasParams.ByChainId {
				if chainIdOrAlias == chainId {
					resolvedChainId = chainId
					break
				}

				if len(aliases.Aliases) == 0 {
					continue
				}

				for _, alias := range aliases.Aliases {
					if alias == chainIdOrAlias {
						resolvedChainId = chainId
						break loopFindChainId
					}
				}
			}
		}
	}

	if resolvedChainId == "" {
		if !dymnsutils.IsValidChainIdFormat(chainIdOrAlias) {
			err = dymnstypes.ErrBadDymNameAddress.Wrap("chain-id is not well-formed")
			return
		}
		resolvedChainId = chainIdOrAlias
	}

	var lookupChainIdConfig, lookupSubName string
	if resolvedChainId == ctx.ChainID() {
		lookupChainIdConfig = ""
	} else {
		lookupChainIdConfig = resolvedChainId
	}
	if len(subNameParts) == 0 {
		lookupSubName = ""
	} else {
		lookupSubName = strings.Join(subNameParts, ".")
	}

	for _, config := range dymName.Configs {
		if config.Type != dymnstypes.DymNameConfigType_NAME {
			continue
		}

		if config.ChainId != lookupChainIdConfig {
			continue
		}

		if config.Path != lookupSubName {
			continue
		}

		outputAddress = config.Value
		return
	}

	// no resolution found

	if lookupSubName == "" {
		if lookupChainIdConfig == "" { // for host chain
			outputAddress = dymName.Owner
		} else {
			// TODO DymNS: implement fallback RollApp chain-id
		}
	}

	return
}

func ParseDymNameAddress(
	dymNameAddress string,
) (
	subNameParts []string, dymName string, chainIdOrAlias string, err error,
) {
	dymNameAddress = strings.ToLower(strings.TrimSpace(dymNameAddress))

	lastDotIndex := strings.LastIndex(dymNameAddress, ".")
	lastAtIndex := strings.LastIndex(dymNameAddress, "@")

	if lastAtIndex > -1 && lastDotIndex > -1 {
		if lastDotIndex > lastAtIndex {
			// do not accept '.' at chain-id/alias part
			err = dymnstypes.ErrBadDymNameAddress.Wrap("misplaced '.'")
			return
		}
	}

	firstDotIndex := strings.IndexRune(dymNameAddress, '.')
	firstAtIndex := strings.IndexRune(dymNameAddress, '@')
	if firstAtIndex > -1 {
		if firstAtIndex != lastAtIndex {
			err = dymnstypes.ErrBadDymNameAddress.Wrap("multiple '@' found")
			return
		}
	}

	if firstDotIndex == 0 || firstAtIndex == 0 {
		err = dymnstypes.ErrBadDymNameAddress
		return
	}

	if lastCharIdx := len(dymNameAddress) - 1; firstDotIndex == lastCharIdx ||
		firstAtIndex == lastCharIdx ||
		lastDotIndex == lastCharIdx ||
		lastAtIndex == lastCharIdx {
		err = dymnstypes.ErrBadDymNameAddress
		return
	}

	if strings.Contains(strings.ReplaceAll(strings.ReplaceAll(dymNameAddress, ".", "|"), "@", "|"), "||") {
		err = dymnstypes.ErrBadDymNameAddress
		return
	}

	chunks := strings.FieldsFunc(dymNameAddress, func(r rune) bool {
		return r == '.' || r == '@'
	})
	for i, chunk := range chunks {
		normalizedChunk := strings.TrimSpace(chunk)
		if normalizedChunk != chunk {
			err = dymnstypes.ErrBadDymNameAddress
			return
		}
		chunks[i] = normalizedChunk
	}

	if len(chunks) == 1 {
		// only Dym-Name, without chain-id/alias,... That is not accepted
		err = dymnstypes.ErrBadDymNameAddress
		return
	}

	chainIdOrAlias = chunks[len(chunks)-1]
	dymName = chunks[len(chunks)-2]
	if len(chunks) > 2 {
		subNameParts = chunks[:len(chunks)-2]
	}

	if !dymnsutils.IsValidChainIdFormat(chainIdOrAlias) &&
		!dymnsutils.IsValidAlias(chainIdOrAlias) {
		err = dymnstypes.ErrBadDymNameAddress.Wrap("chain-id/alias is not well-formed")
		return
	}

	if !dymnsutils.IsValidDymName(dymName) {
		err = dymnstypes.ErrBadDymNameAddress.Wrap("Dym-Name is not well-formed")
		return
	}

	if len(subNameParts) > 0 {
		for _, subNamePart := range subNameParts {
			if !dymnsutils.IsValidDymName(subNamePart) {
				err = dymnstypes.ErrBadDymNameAddress.Wrap("sub-Dym-Name is not well-formed")
				return
			}
		}
	}

	return
}
