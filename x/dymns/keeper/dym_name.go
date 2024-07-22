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
//
// Note:
//  1. Must call BeforeDymNameOwnerChanged and AfterDymNameOwnerChanged before and after calling this function when updating owner.
//  2. Must call BeforeDymNameConfigChanged and AfterDymNameConfigChanged before and after calling this function when updating configuration.
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

	return nil
}

// BeforeDymNameOwnerChanged must be called before updating the owner of a Dym-Name.
func (k Keeper) BeforeDymNameOwnerChanged(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return nil
	}

	if err := k.RemoveReverseMappingOwnerToOwnedDymName(ctx, dymName.Owner, dymName.Name); err != nil {
		return err
	}

	return nil
}

// AfterDymNameOwnerChanged must be called after the owner of a Dym-Name is changed.
func (k Keeper) AfterDymNameOwnerChanged(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return dymnstypes.ErrDymNameNotFound.Wrap(name)
	}

	if err := k.AddReverseMappingOwnerToOwnedDymName(ctx, dymName.Owner, name); err != nil {
		return err
	}

	return nil
}

// BeforeDymNameConfigChanged must be called before updating the configuration of a Dym-Name.
func (k Keeper) BeforeDymNameConfigChanged(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return nil
	}

	configuredAddresses, hexAddresses := dymName.GetAddressesForReverseMapping()
	for configuredAddress := range configuredAddresses {
		if err := k.RemoveReverseMappingConfiguredAddressToDymName(ctx, configuredAddress, name); err != nil {
			return err
		}
	}
	for hexAddress := range hexAddresses {
		var bzAddr []byte
		bzAddr = dymnsutils.GetBytesFromHexAddress(hexAddress)
		if err := k.RemoveReverseMappingHexAddressToDymName(ctx, bzAddr, name); err != nil {
			return err
		}
	}

	return nil
}

// AfterDymNameConfigChanged must be called after the configuration of a Dym-Name is changed.
func (k Keeper) AfterDymNameConfigChanged(ctx sdk.Context, name string) error {
	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return dymnstypes.ErrDymNameNotFound.Wrap(name)
	}

	configuredAddresses, hexAddresses := dymName.GetAddressesForReverseMapping()
	for configuredAddress := range configuredAddresses {
		if err := k.AddReverseMappingConfiguredAddressToDymName(ctx, configuredAddress, name); err != nil {
			return err
		}
	}
	for hexAddress := range hexAddresses {
		var bzAddr []byte
		bzAddr = dymnsutils.GetBytesFromHexAddress(hexAddress)
		if err := k.AddReverseMappingHexAddressToDymName(ctx, bzAddr, name); err != nil {
			return err
		}
	}

	return nil
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

func (k Keeper) DeleteDymName(ctx sdk.Context, name string) error {
	if err := k.BeforeDymNameOwnerChanged(ctx, name); err != nil {
		return err
	}

	if err := k.BeforeDymNameConfigChanged(ctx, name); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	dymNameKey := dymnstypes.DymNameKey(name)
	store.Delete(dymNameKey)

	return nil
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

// PruneDymName removes a Dym-Name from the KVStore, as well as all related records.
func (k Keeper) PruneDymName(ctx sdk.Context, name string) error {
	// remove SO (force, ignore active SO)
	k.DeleteSellOrder(ctx, name)

	// remove historical SO
	k.DeleteHistoricalSellOrders(ctx, name)
	k.SetMinExpiryHistoricalSellOrder(ctx, name, 0)

	dymName := k.GetDymName(ctx, name)
	if dymName == nil {
		return nil
	}

	// remove config
	// This seems not necessary because we are going to remove the record anyway,
	// but just let it here to clear the business logic
	dymName.Configs = nil   // all configuration should be removed
	dymName.Owner = ""      // no one owns it anyone
	dymName.Controller = "" // no one controls it anyone

	return k.DeleteDymName(ctx, name)
}

func (k Keeper) ResolveByDymNameAddress(ctx sdk.Context, dymNameAddress string) (outputAddress string, err error) {
	subName, name, chainIdOrAlias, parseErr := ParseDymNameAddress(dymNameAddress)
	if parseErr != nil {
		ctx.Logger().Debug("failed to parse Dym-Name address", "dym-name-address", dymNameAddress, "error", parseErr)
		err = parseErr
		return
	}

	dymName := k.GetDymNameWithExpirationCheck(ctx, name, ctx.BlockTime().Unix())
	if dymName == nil {
		err = dymnstypes.ErrDymNameNotFound.Wrap(name)

		// Dym-Name not found, in this case, there are 3 possible reasons:
		// 1. Dym-Name does not exist
		// 2. Dym-Name was expired
		// 3. Resolve extra format 0x1234...6789@nim and dym1...@nim
		// First two cases, stop here.
		// If it is the third case, we need to resolve it.

		if subName != "" {
			return
		}

		var accAddr sdk.AccAddress
		if dymnsutils.IsValidHexAddress(name) {
			accAddr = dymnsutils.GetBytesFromHexAddress(name)
		} else if dymnsutils.IsValidBech32AccountAddress(name, false) {
			_, bz, errDecode := bech32.DecodeAndConvert(name)
			if errDecode != nil {
				return
			}

			accAddr = bz
		} else {
			// neither hex address nor bech32 account address
			return
		}

		chainId, success := k.tryResolveChainIdOrAliasToChainId(ctx, chainIdOrAlias)
		if !success {
			// not a known chain-id or alias
			return
		}

		// only accept resolve for host chain or RollApp

		if chainId == ctx.ChainID() {
			// is host chain
			outputAddress = accAddr.String()
			err = nil
			return
		}

		if !k.IsRollAppId(ctx, chainId) {
			return
		}

		bech32Prefix, found := k.GetRollAppBech32Prefix(ctx, chainId)
		if !found {
			// no bech32 prefix configured for this RollApp
			return
		}
		rollAppBasedBech32Addr, convertErr := bech32.ConvertAndEncode(bech32Prefix, accAddr)
		if convertErr != nil {
			return
		}

		outputAddress = rollAppBasedBech32Addr
		err = nil
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

	if chainId, success := k.tryResolveChainIdOrAliasToChainId(ctx, chainIdOrAlias); success {
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
	if chainIdOrAlias == ctx.ChainID() {
		return chainIdOrAlias, true
	}

	chainsParams := k.ChainsParams(ctx)
	if len(chainsParams.AliasesByChainId) > 0 {
		for chainId, aliases := range chainsParams.AliasesByChainId {
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

	if subName == "" {
		// when no sub-name, we have 2 valid formats that won't pass Dym-Name validation

		// 0x1234...6789@nim
		if dymnsutils.IsValidHexAddress(dymName) {
			return
		} else if dymnsutils.IsValidBech32AccountAddress(dymName, false) {
			return
		}
	}

	if !dymnsutils.IsValidDymName(dymName) {
		err = dymnstypes.ErrBadDymNameAddress.Wrap("Dym-Name is not well-formed")
		return
	}

	return
}

func (k Keeper) ReverseResolveDymNameAddress(ctx sdk.Context, inputAddress, workingChainId string) (outputDymNameAddresses dymnstypes.ReverseResolvedDymNameAddresses, err error) {
	inputAddress = strings.ToLower(inputAddress)
	nowEpoch := ctx.BlockTime().Unix()

	isBech32Addr := dymnsutils.IsValidBech32AccountAddress(inputAddress, false)
	is0xAddr := dymnsutils.IsValidHexAddress(inputAddress)

	if inputAddress == "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("input address cannot be blank")
	}

	if !isBech32Addr && !is0xAddr {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("not supported address format: %s", inputAddress)
	}

	if workingChainId == "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("working chain-id cannot be blank")
	}

	if !dymnsutils.IsValidChainIdFormat(workingChainId) {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("invalid chain-id format: %s", workingChainId)
	}

	workingChainIdIsHostChain := workingChainId == ctx.ChainID()
	isWorkingOnCoinType60Chain := workingChainIdIsHostChain || k.CheckChainIsCoinType60ByChainId(ctx, workingChainId)

	defer func() {
		if len(outputDymNameAddresses) < 1 {
			return
		}

		outputDymNameAddresses = outputDymNameAddresses.Distinct()
		outputDymNameAddresses = k.ReplaceChainIdWithAliasIfPossible(ctx, outputDymNameAddresses)
		outputDymNameAddresses.Sort()
	}()

	addConfigs := func(dymName dymnstypes.DymName, configs []dymnstypes.DymNameConfig) {
		for _, config := range configs {
			record := dymnstypes.ReverseResolvedDymNameAddress{
				SubName:        config.Path,
				Name:           dymName.Name,
				ChainIdOrAlias: config.ChainId,
			}
			if config.ChainId == "" {
				record.ChainIdOrAlias = ctx.ChainID()
			}

			if record.ChainIdOrAlias != workingChainId {
				// skip this config record
				continue
			}

			outputDymNameAddresses = append(outputDymNameAddresses, record)
		}
	}

	if is0xAddr {
		hexAddr := inputAddress
		bzAddr := dymnsutils.GetBytesFromHexAddress(hexAddr)

		// When the input address is a hex address,
		// we can assume the query comes from a coin-type-60 chain.
		// The working chain-id can be either host chain or another chain.
		//
		// Should we do a fallback lookup?
		// I think it should depend on the working chain-id is recognized or not.
		// - If the working chain-id is recognized, we should do a fallback lookup.
		// - If the working chain-id is NOT recognized, we should not do a fallback lookup.

		// But first, let try to convert it into a bech32 address to see if any result is available.
		var bech32Hrp string
		if workingChainIdIsHostChain {
			bech32Hrp = sdk.GetConfig().GetBech32AccountAddrPrefix()
		} else if rollappBech32Hrp, found := k.GetRollAppBech32Prefix(ctx, workingChainId); found {
			bech32Hrp = rollappBech32Hrp
		}

		if bech32Hrp != "" {
			bech32Addr := sdk.MustBech32ifyAddressBytes(bech32Hrp, bzAddr)
			dymNames, err1 := k.GetDymNamesContainsConfiguredAddress(ctx, bech32Addr, nowEpoch)
			if err1 != nil {
				return nil, err1
			}

			for _, dymName := range dymNames {
				configuredAddresses, _ := dymName.GetAddressesForReverseMapping()
				configs, found := configuredAddresses[bech32Addr]
				if !found {
					continue
				}

				addConfigs(dymName, configs)
			}

			if len(outputDymNameAddresses) > 0 {
				// we found at least one result, can stop here
				return
			}
		} else {
			// no bech32 prefix configured for this chain-id
		}

		// check if we should do a fallback lookup

		if !isWorkingOnCoinType60Chain {
			// we don't do fallback lookup for this case, just for safety purpose
			return
		}

		// we are going to do a fallback lookup
		dymNames, err2 := k.GetDymNamesContainsHexAddress(ctx, bzAddr, nowEpoch)
		if err2 != nil {
			return nil, err2
		}

		for _, dymName := range dymNames {
			_, hexAddresses := dymName.GetAddressesForReverseMapping()
			configs, found := hexAddresses[hexAddr]
			if !found {
				continue
			}

			for _, config := range configs {
				if config.IsDefaultNameConfig() {
					// only accept fallback for the case of default config
					outputDymNameAddresses = append(outputDymNameAddresses, dymnstypes.ReverseResolvedDymNameAddress{
						SubName:        "",
						Name:           dymName.Name,
						ChainIdOrAlias: workingChainId, // fallback
					})
					break
				}
			}
		}

		return
	}

	// try lookup using bech32 address

	bech32Addr := inputAddress

	dymNames, err1 := k.GetDymNamesContainsConfiguredAddress(ctx, bech32Addr, nowEpoch)
	if err1 != nil {
		return nil, err1
	}

	for _, dymName := range dymNames {
		configuredAddresses, _ := dymName.GetAddressesForReverseMapping()
		configs, found := configuredAddresses[bech32Addr]
		if !found {
			continue
		}

		addConfigs(dymName, configs)
	}

	if len(outputDymNameAddresses) > 0 {
		return
	}

	// There is no matching result from lookup by bech32 address,
	// we are going to give it one more try lookup by hex address
	// to see if any fallback is available.
	// If the working chain-id is a coin-type-60 chain-id.
	// If the working chain-id is NOT a coin-type-60 chain-id, does not satisfy the condition.

	if !isWorkingOnCoinType60Chain {
		// we don't do fallback lookup for this case, just for safety purpose
		return
	}

	_, bzHexAddr, err2 := bech32.DecodeAndConvert(bech32Addr)
	if err2 != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to decode bech32 address %s: %v", bech32Addr, err2)
	}

	dymNames, err3 := k.GetDymNamesContainsHexAddress(ctx, bzHexAddr, nowEpoch)
	if err3 != nil {
		return nil, err3
	}

	hexAddr := dymnsutils.GetHexAddressFromBytes(bzHexAddr)

	for _, dymName := range dymNames {
		_, hexAddresses := dymName.GetAddressesForReverseMapping()
		configs, found := hexAddresses[hexAddr]
		if !found {
			continue
		}

		for _, config := range configs {
			if config.IsDefaultNameConfig() {
				// only accept fallback for the case of default config
				outputDymNameAddresses = append(outputDymNameAddresses, dymnstypes.ReverseResolvedDymNameAddress{
					SubName:        "",
					Name:           dymName.Name,
					ChainIdOrAlias: workingChainId, // fallback
				})
				break
			}
		}
	}

	return
}

func (k Keeper) ReplaceChainIdWithAliasIfPossible(ctx sdk.Context, reverseResolvedRecords dymnstypes.ReverseResolvedDymNameAddresses) []dymnstypes.ReverseResolvedDymNameAddress {
	if len(reverseResolvedRecords) > 0 {
		for i, reverseResolvedRecord := range reverseResolvedRecords {
			if reverseResolvedRecord.ChainIdOrAlias == "" {
				reverseResolvedRecords[i].ChainIdOrAlias = ctx.ChainID()
			}
		}

		chainIdToAlias := make(map[string]string)

		chainsParams := k.ChainsParams(ctx)
		for chainId, aliases := range chainsParams.AliasesByChainId {
			if len(aliases.Aliases) > 0 {
				chainIdToAlias[chainId] = aliases.Aliases[0]
			}
		}

		for i, reverseResolvedRecord := range reverseResolvedRecords {
			chainId := reverseResolvedRecord.ChainIdOrAlias
			if alias, found := chainIdToAlias[chainId]; found {
				if len(alias) > 0 {
					reverseResolvedRecords[i].ChainIdOrAlias = alias
				}
				continue
			}

			isRollAppId := k.IsRollAppId(ctx, chainId)
			if !isRollAppId {
				chainIdToAlias[chainId] = chainId
				continue
			}

			alias, found := k.GetAliasByRollAppId(ctx, chainId)
			if !found {
				chainIdToAlias[chainId] = chainId
				continue
			}

			chainIdToAlias[chainId] = alias
			reverseResolvedRecords[i].ChainIdOrAlias = alias
		}
	}

	return reverseResolvedRecords
}
