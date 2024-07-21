package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
)

var _ epochstypes.EpochHooks = epochHooks{}

type epochHooks struct {
	Keeper
}

func (k Keeper) GetEpochHooks() epochstypes.EpochHooks {
	return epochHooks{
		Keeper: k,
	}
}

// BeforeEpochStart is the epoch start hook.
// Business logic is to prune historical sell orders.
func (e epochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := e.GetParams(ctx)

	if epochIdentifier != params.Misc.BeginEpochHookIdentifier {
		return nil
	}

	e.Keeper.Logger(ctx).Info("DymNS hook Before-Epoch-Start: triggered", "epoch-number", epochNumber, "epoch-identifier", epochIdentifier)

	return e.processCleanupHistoricalSellOrders(ctx, epochIdentifier, epochNumber, params)
}

func (e epochHooks) processCleanupHistoricalSellOrders(ctx sdk.Context, epochIdentifier string, epochNumber int64, params dymnstypes.Params) error {
	dk := e.Keeper

	/**
	We use this method instead of iterating through all historical sell orders.
	It helps reduce number of IO needed to read all historical sell orders.
	*/
	nameToMinExpiry := dk.GetMinExpiryOfAllHistoricalSellOrders(ctx)
	if len(nameToMinExpiry) < 1 {
		return nil
	}

	cleanBeforeEpochUTC := ctx.BlockTime().Unix() - 86400*int64(params.Misc.DaysPreservedClosedSellOrder)

	var cleanupHistoricalForDymNames []string
	for name, minExpiry := range nameToMinExpiry {
		if minExpiry > cleanBeforeEpochUTC {
			continue
		}

		cleanupHistoricalForDymNames = append(cleanupHistoricalForDymNames, name)
	}
	if len(cleanupHistoricalForDymNames) < 1 {
		return nil
	}

	e.Keeper.Logger(ctx).Info(
		"DymNS hook Before-Epoch-Start: processing cleanup historical sell orders",
		"count", len(cleanupHistoricalForDymNames),
		"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
	)

	for _, dymName := range cleanupHistoricalForDymNames {
		list := dk.GetHistoricalSellOrders(ctx, dymName)
		if len(list) < 1 {
			dk.SetMinExpiryHistoricalSellOrder(ctx, dymName, 0)
			continue
		}

		var keepList []dymnstypes.SellOrder
		for _, hso := range list {
			if hso.ExpireAt > cleanBeforeEpochUTC {
				keepList = append(keepList, hso)
			}
		}

		if len(keepList) == 0 {
			dk.DeleteHistoricalSellOrders(ctx, dymName)
			dk.SetMinExpiryHistoricalSellOrder(ctx, dymName, 0)
			continue
		}

		var newMinExpiry = keepList[0].ExpireAt
		for _, hso := range keepList {
			if hso.ExpireAt < newMinExpiry {
				newMinExpiry = hso.ExpireAt
			}
		}
		dk.SetMinExpiryHistoricalSellOrder(ctx, dymName, newMinExpiry)

		if len(keepList) != len(list) {
			hso := dymnstypes.HistoricalSellOrders{
				SellOrders: keepList,
			}
			dk.SetHistoricalSellOrders(ctx, dymName, hso)
		}
	}

	return nil
}

// AfterEpochEnd is the epoch end hook.
// Business logic is to move expired sell orders to historical
// and if SO has a winner, complete the sell order.
func (e epochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := e.GetParams(ctx)

	if epochIdentifier != params.Misc.EndEpochHookIdentifier {
		return nil
	}

	e.Keeper.Logger(ctx).Info("DymNS hook After-Epoch-End: triggered", "epoch-number", epochNumber, "epoch-identifier", epochIdentifier)

	return e.processActiveSellOrders(ctx, epochIdentifier, epochNumber)
}

func (e epochHooks) processActiveSellOrders(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	dk := e.Keeper

	aope := dk.GetActiveSellOrdersExpiration(ctx)
	nowEpochUTC := ctx.BlockTime().Unix()
	var finishedSOs []dymnstypes.SellOrder
	updateExpiryForSONames := make(map[string]int64)
	if len(aope.ExpiryByName) > 0 {
		for name, expiry := range aope.ExpiryByName {
			if expiry > nowEpochUTC {
				continue
			}

			so := dk.GetSellOrder(ctx, name)

			if so == nil {
				// remove the invalid entry
				delete(aope.ExpiryByName, name) // safe to delete during iteration
				continue
			}

			if !so.HasFinished(nowEpochUTC) {
				// invalid entry
				dk.Logger(ctx).Error(
					"DymNS hook After-Epoch-End: sell order has not finished",
					"name", name, "expiry", expiry, "now", nowEpochUTC,
					"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
				)

				// later update because it is Not safe to modify during iteration
				updateExpiryForSONames[name] = so.ExpireAt
				continue
			}

			finishedSOs = append(finishedSOs, *so)
		}
	}

	if len(finishedSOs) < 1 {
		// skip updating store
		return nil
	}

	dk.Logger(ctx).Info(
		"DymNS hook After-Epoch-End: processing finished SOs", "count", len(finishedSOs),
		"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
	)

	for name, newExpiry := range updateExpiryForSONames {
		aope.ExpiryByName[name] = newExpiry
	}
	updateExpiryForSONames = nil

	for _, so := range finishedSOs {
		delete(aope.ExpiryByName, so.Name) // delete the processed entry

		if so.HighestBid != nil {
			if err := dk.CompleteSellOrder(ctx, so.Name); err != nil {
				dk.Logger(ctx).Error(
					"DymNS hook After-Epoch-End: failed to complete sell order",
					"name", so.Name, "expiry", so.ExpireAt, "now", nowEpochUTC,
					"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
					"error", err,
				)
				return err
			}
			continue
		}

		// no bid placed, it just a normal expiry without winner,
		// in this case, just move to history
		if err := dk.MoveSellOrderToHistorical(ctx, so.Name); err != nil {
			dk.Logger(ctx).Error(
				"DymNS hook After-Epoch-End: failed to move expired sell order to historical",
				"name", so.Name, "expiry", so.ExpireAt, "now", nowEpochUTC,
				"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
				"error", err,
			)
			return err
		}
	}

	if err := dk.SetActiveSellOrdersExpiration(ctx, aope); err != nil {
		dk.Logger(ctx).Error(
			"DymNS hook After-Epoch-End: failed to update active SO expiry",
			"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
			"error", err,
		)
		return err
	}

	return nil
}
