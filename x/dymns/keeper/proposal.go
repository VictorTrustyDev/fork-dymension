package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func (k Keeper) MigrateChainIds(ctx sdk.Context, replacement []dymnstypes.MigrateChainId) error {
	previousChainIds := make(map[string]dymnstypes.MigrateChainId)

	for _, r := range replacement {
		previousChainIds[r.PreviousChainId] = r
	}

	// migrate params
	params := k.GetParams(ctx)

	if len(params.Chains.AliasesByChainId) > 0 {
		newAliasesByChainId := make(map[string]dymnstypes.AliasesOfChainId)
		for chainId, aliases := range params.Chains.AliasesByChainId {
			if r, isPreviousChainId := previousChainIds[chainId]; isPreviousChainId {
				if _, foundDeclared := params.Chains.AliasesByChainId[r.NewChainId]; foundDeclared {
					// we don't override, we keep the aliases of the new chain id

					// ignore and remove the aliases of the previous chain id
				} else {
					newAliasesByChainId[r.NewChainId] = aliases
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
			if r, isPreviousChainId := previousChainIds[chainId]; isPreviousChainId {
				newCoinType60UniqueChainIds[r.NewChainId] = true
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
		k.Logger(ctx).Error("failed to update params", "error", err)
		return err
	}

	// end of params migration

	// Migrate chain ids within configs of Dym-Name.
	// We only migrate for Dym-Names that not expired to reduce IO needed.

	nonExpiredDymNames := k.GetAllNonExpiredDymNames(ctx, ctx.BlockTime().Unix())
	if len(nonExpiredDymNames) > 0 {
		for _, dymName := range nonExpiredDymNames {
			newConfigs := make([]dymnstypes.DymNameConfig, len(dymName.Configs))
			var anyConfigUpdated bool
			for i, config := range dymName.Configs {
				if config.ChainId != "" {
					if r, isPreviousChainId := previousChainIds[config.ChainId]; isPreviousChainId {
						config.ChainId = r.NewChainId
						anyConfigUpdated = true
					}
				}

				newConfigs[i] = config
			}

			if !anyConfigUpdated {
				continue
			}

			dymName.Configs = newConfigs

			if err := dymName.Validate(); err != nil {
				k.Logger(ctx).Error(
					"failed to migrate chain ids for Dym-Name",
					"dymName", dymName.Name,
					"error", err,
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
				)
				return err
			}

			if err := k.SetDymName(ctx, dymName); err != nil {
				k.Logger(ctx).Error(
					"failed to migrate chain ids for Dym-Name",
					"dymName", dymName.Name,
					"step", "SetDymName",
					"error", err,
				)
				return err
			}

			if err := k.AfterDymNameConfigChanged(ctx, dymName.Name); err != nil {
				k.Logger(ctx).Error(
					"failed to migrate chain ids for Dym-Name",
					"dymName", dymName.Name,
					"step", "AfterDymNameConfigChanged",
					"error", err,
				)
				return err
			}

			k.Logger(ctx).Info("migrated chain ids for Dym-Name", "dymName", dymName.Name)
		}
	}

	// end of Dym-Name config migration

	return nil
}
