package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k Keeper) MigrateChainIds(ctx sdk.Context, replacement []dymnstypes.MigrateChainId) error {
	previousChainIdsToNewChainId := make(map[string]string)

	for _, r := range replacement {
		previousChainIdsToNewChainId[r.PreviousChainId] = r.NewChainId
	}

	if err := k.migrateChainIdsInParams(ctx, previousChainIdsToNewChainId); err != nil {
		return err
	}

	if err := k.migrateChainIdsInDymNames(ctx, previousChainIdsToNewChainId); err != nil {
		return err
	}

	return nil
}

func (k Keeper) migrateChainIdsInParams(ctx sdk.Context, previousChainIdsToNewChainId map[string]string) error {
	params := k.GetParams(ctx)

	if len(params.Chains.AliasesByChainId) > 0 {
		newAliasesByChainId := make(map[string]dymnstypes.AliasesOfChainId)
		for chainId, aliases := range params.Chains.AliasesByChainId {
			if newChainId, isPreviousChainId := previousChainIdsToNewChainId[chainId]; isPreviousChainId {
				if _, foundDeclared := params.Chains.AliasesByChainId[newChainId]; foundDeclared {
					// we don't override, we keep the aliases of the new chain id

					// ignore and remove the aliases of the previous chain id
				} else {
					newAliasesByChainId[newChainId] = aliases
				}
			} else {
				newAliasesByChainId[chainId] = aliases
			}
		}
		params.Chains.AliasesByChainId = newAliasesByChainId
	}

	if len(params.Chains.CoinType60ChainIds) > 0 {
		newCoinType60UniqueChainIds := make(map[string]bool)

		for _, chainId := range params.Chains.CoinType60ChainIds {
			if newChainId, isPreviousChainId := previousChainIdsToNewChainId[chainId]; isPreviousChainId {
				newCoinType60UniqueChainIds[newChainId] = true
			} else {
				newCoinType60UniqueChainIds[chainId] = true
			}
		}

		newCoinType60ChainIds := make([]string, 0, len(newCoinType60UniqueChainIds))
		for chainId := range newCoinType60UniqueChainIds {
			newCoinType60ChainIds = append(newCoinType60ChainIds, chainId)
		}

		params.Chains.CoinType60ChainIds = newCoinType60ChainIds
	}

	if err := params.Validate(); err != nil {
		// TODO DymNS: write test-case that cover this error then delete this statement
		//  because SetParams already validate the params
		return errors.Wrap(err, "failed to update params")
	}

	if err := k.SetParams(ctx, params); err != nil {
		k.Logger(ctx).Error(
			"failed to update params",
			"error", err,
			"migration-state", "aborted",
		)
		return err
	}

	return nil
}

func (k Keeper) migrateChainIdsInDymNames(ctx sdk.Context, previousChainIdsToNewChainId map[string]string) error {
	// We only migrate for Dym-Names that not expired to reduce IO needed.

	nonExpiredDymNames := k.GetAllNonExpiredDymNames(ctx, ctx.BlockTime().Unix())
	if len(nonExpiredDymNames) < 1 {
		return nil
	}

	for _, dymName := range nonExpiredDymNames {
		newConfigs := make([]dymnstypes.DymNameConfig, len(dymName.Configs))
		var anyConfigUpdated bool
		for i, config := range dymName.Configs {
			if config.ChainId != "" {
				if newChainId, isPreviousChainId := previousChainIdsToNewChainId[config.ChainId]; isPreviousChainId {
					config.ChainId = newChainId
					anyConfigUpdated = true
				}
			}

			newConfigs[i] = config
		}

		if !anyConfigUpdated {
			// Skip migration for this Dym-Name if nothing updated to reduce IO.
			continue
		}

		dymName.Configs = newConfigs

		if err := dymName.Validate(); err != nil {
			k.Logger(ctx).Error(
				"failed to migrate chain ids for Dym-Name",
				"dymName", dymName.Name,
				"step", "Validate",
				"error", err,
				"migration-state", "continue",
			)
			// Skip migration for this Dym-Name.
			// We don't want to break the migration process for other Dym-Names.
			// The replacement should be done later by the owner.
			continue
		}

		// from here, any step can procedures dirty state, so we need to abort the migration

		if err := k.BeforeDymNameConfigChanged(ctx, dymName.Name); err != nil {
			k.Logger(ctx).Error(
				"failed to migrate chain ids for Dym-Name",
				"dymName", dymName.Name,
				"step", "BeforeDymNameConfigChanged",
				"error", err,
				"migration-state", "aborted",
			)
			return err
		}

		if err := k.SetDymName(ctx, dymName); err != nil {
			k.Logger(ctx).Error(
				"failed to migrate chain ids for Dym-Name",
				"dymName", dymName.Name,
				"step", "SetDymName",
				"error", err,
				"migration-state", "aborted",
			)
			return err
		}

		if err := k.AfterDymNameConfigChanged(ctx, dymName.Name); err != nil {
			k.Logger(ctx).Error(
				"failed to migrate chain ids for Dym-Name",
				"dymName", dymName.Name,
				"step", "AfterDymNameConfigChanged",
				"error", err,
				"migration-state", "aborted",
			)
			return err
		}

		k.Logger(ctx).Info("migrated chain ids for Dym-Name", "dymName", dymName.Name)
	}

	return nil
}
