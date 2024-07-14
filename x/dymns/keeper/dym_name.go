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

// SetReverseMappingOwnerToOwnedDymName stores a reverse mapping from owner to owned Dym-Name into the KVStore.
func (k Keeper) SetReverseMappingOwnerToOwnedDymName(ctx sdk.Context, owner, name string) error {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(bzAccAddr)

	var existingOwnedDymNames dymnstypes.OwnedDymNames

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(dymNamesOwnedByAccountKey)
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &existingOwnedDymNames)
		for _, owned := range existingOwnedDymNames.DymNames {
			if owned == name {
				// reverse lookup already exists
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

// GetDymNamesOwnedBy returns all Dym-Names owned by the account address.
func (k Keeper) GetDymNamesOwnedBy(
	ctx sdk.Context, owner string, nowEpoch int64,
) ([]dymnstypes.DymName, error) {
	_, bzAccAddr, err := bech32.DecodeAndConvert(owner)
	if err != nil {
		return nil, dymnstypes.ErrInvalidOwner.Wrap(owner)
	}

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(bzAccAddr)

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

	dymNamesOwnedByAccountKey := dymnstypes.DymNamesOwnedByAccountRvlKey(accAddr)

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
	subName, name, chainIdOrAlias, parseErr := ParseDymNameAddress(dymNameAddress)
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

	tryResolveFromConfig := func(lookupChainIdConfig string) (value string, found bool) {
		if lookupChainIdConfig == ctx.ChainID() {
			lookupChainIdConfig = ""
		}

		for _, config := range dymName.Configs {
			if config.Type != dymnstypes.DymNameConfigType_NAME {
				continue
			}

			if config.ChainId != lookupChainIdConfig {
				continue
			}

			if config.Path != subName {
				continue
			}

			return config.Value, true
		}

		return "", false
	}

	// first attempt to resolve, to see if the chain-id/alias is user-configured
	var found bool
	outputAddress, found = tryResolveFromConfig(chainIdOrAlias)
	if found {
		return
	}

	// end of first attempt

	var resolvedToChainId string

	if chainIdOrAlias == ctx.ChainID() {
		resolvedToChainId = chainIdOrAlias
	} else if chainId, success := k.tryResolveChainIdOrAliasToChainId(ctx, chainIdOrAlias); success {
		resolvedToChainId = chainId
	} else {
		// treat it as chain-id
		resolvedToChainId = chainIdOrAlias
	}

	// second attempt to resolve
	outputAddress, found = tryResolveFromConfig(resolvedToChainId)
	if found {
		return
	}
	// end of second attempt

	// no more attempts to resolve from config

	// try fallback

	if subName != "" {
		// no need to fallback
		return
	}

	if resolvedToChainId == ctx.ChainID() {
		outputAddress = dymName.Owner
		return
	}

	// not the host chain, try to resolve if is a RollApp

	isRollAppId := k.IsRollAppId(ctx, resolvedToChainId)
	if !isRollAppId {
		// fallback does not apply for non-RollApp
		return
	}

	rollAppBech32Prefix, found := k.GetRollAppBech32Prefix(ctx, resolvedToChainId)
	if !found {
		// fallback does not apply for RollApp does not have this metadata
		return
	}

	accAddr := sdk.MustAccAddressFromBech32(dymName.Owner)
	rollAppBasedBech32Addr, convertErr := bech32.ConvertAndEncode(rollAppBech32Prefix, accAddr)
	if convertErr != nil {
		err = sdkerrors.ErrInvalidAddress.Wrapf(
			"failed to convert '%s' to RollApp-based address: %v", dymName.Owner, convertErr,
		)
		return
	}

	outputAddress = rollAppBasedBech32Addr
	return
}

func (k Keeper) tryResolveChainIdOrAliasToChainId(ctx sdk.Context, chainIdOrAlias string) (resolvedToChainId string, success bool) {
	aliasParams := k.AliasParams(ctx)
	if len(aliasParams.ByChainId) > 0 {
		for chainId, aliases := range aliasParams.ByChainId {
			if chainIdOrAlias == chainId {
				return chainId, true
			}

			for _, alias := range aliases.Aliases {
				if alias == chainIdOrAlias {
					return chainId, true
				}
			}
		}
	}

	if isRollAppId := k.IsRollAppId(ctx, chainIdOrAlias); isRollAppId {
		return chainIdOrAlias, true
	}

	if rollAppId, found := k.GetRollAppIdByAlias(ctx, chainIdOrAlias); found {
		// TODO DymNS: require RollApp alias to use dymnsutils.IsValidAlias to validate the alias
		return rollAppId, true
	}

	return
}

func ParseDymNameAddress(
	dymNameAddress string,
) (
	subName string, dymName string, chainIdOrAlias string, err error,
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
		subNameParts := chunks[:len(chunks)-2]
		for _, subNamePart := range subNameParts {
			if !dymnsutils.IsValidDymName(subNamePart) {
				err = dymnstypes.ErrBadDymNameAddress.Wrap("sub-Dym-Name is not well-formed")
				return
			}
		}
		subName = strings.Join(subNameParts, ".")
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

	return
}
