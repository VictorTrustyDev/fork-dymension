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
// Business logic is to prune historical purchase orders.
func (e epochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := e.GetParams(ctx)

	if epochIdentifier != params.Misc.BeginEpochHookIdentifier {
		return nil
	}

	e.Keeper.Logger(ctx).Info("DymNS hook Before-Epoch-Start: triggered", "epoch-number", epochNumber, "epoch-identifier", epochIdentifier)

	return e.processCleanupHistoricalPurchaseOrder(ctx, epochIdentifier, epochNumber, params)
}

func (e epochHooks) processCleanupHistoricalPurchaseOrder(ctx sdk.Context, epochIdentifier string, epochNumber int64, params dymnstypes.Params) error {
	dk := e.Keeper

	/**
	We use this method instead of iterating through all historical purchase orders.
	It helps reduce number of IO needed to read all historical purchase orders.
	*/
	nameToMinExpiry := dk.GetMinExpiryOfAllHistoricalOpenPurchaseOrders(ctx)
	if len(nameToMinExpiry) < 1 {
		return nil
	}

	cleanBeforeEpochUTC := ctx.BlockTime().Unix() - 86400*int64(params.Misc.DaysPreservedClosedPurchaseOrder)

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
		"DymNS hook Before-Epoch-Start: processing cleanup historical purchase orders",
		"count", len(cleanupHistoricalForDymNames),
		"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
	)

	for _, dymName := range cleanupHistoricalForDymNames {
		list := dk.GetHistoricalOpenPurchaseOrders(ctx, dymName)
		if len(list) < 1 {
			dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, 0)
			continue
		}

		var keepList []dymnstypes.OpenPurchaseOrder
		for _, hOpo := range list {
			if hOpo.ExpireAt > cleanBeforeEpochUTC {
				keepList = append(keepList, hOpo)
			}
		}

		if len(keepList) == 0 {
			dk.DeleteHistoricalOpenPurchaseOrders(ctx, dymName)
			dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, 0)
			continue
		}

		var newMinExpiry = keepList[0].ExpireAt
		for _, hOpo := range keepList {
			if hOpo.ExpireAt < newMinExpiry {
				newMinExpiry = hOpo.ExpireAt
			}
		}
		dk.SetMinExpiryHistoricalOpenPurchaseOrder(ctx, dymName, newMinExpiry)

		if len(keepList) != len(list) {
			hOpo := dymnstypes.HistoricalOpenPurchaseOrders{
				OpenPurchaseOrders: keepList,
			}
			dk.SetHistoricalOpenPurchaseOrders(ctx, dymName, hOpo)
		}
	}

	return nil
}

// AfterEpochEnd is the epoch end hook.
// Business logic is to move expired open purchase orders to historical
// and if OPO has a winner, complete the purchase order.
func (e epochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := e.GetParams(ctx)

	if epochIdentifier != params.Misc.EndEpochHookIdentifier {
		return nil
	}

	e.Keeper.Logger(ctx).Info("DymNS hook After-Epoch-End: triggered", "epoch-number", epochNumber, "epoch-identifier", epochIdentifier)

	return e.processActiveOPO(ctx, epochIdentifier, epochNumber)
}

func (e epochHooks) processActiveOPO(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	dk := e.Keeper

	aope := dk.GetActiveOpenPurchaseOrdersExpiration(ctx)
	nowEpochUTC := ctx.BlockTime().Unix()
	var finishedOPOs []dymnstypes.OpenPurchaseOrder
	updateExpiryForOPONames := make(map[string]int64)
	if len(aope.ExpiryByName) > 0 {
		for name, expiry := range aope.ExpiryByName {
			if expiry > nowEpochUTC {
				continue
			}

			opo := dk.GetOpenPurchaseOrder(ctx, name)

			if opo == nil {
				// remove the invalid entry
				delete(aope.ExpiryByName, name) // safe to delete during iteration
				continue
			}

			if !opo.HasFinished(nowEpochUTC) {
				// invalid entry
				dk.Logger(ctx).Error(
					"DymNS hook After-Epoch-End: open purchase order has not finished",
					"name", name, "expiry", expiry, "now", nowEpochUTC,
					"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
				)

				// later update because it is Not safe to modify during iteration
				updateExpiryForOPONames[name] = opo.ExpireAt
				continue
			}

			finishedOPOs = append(finishedOPOs, *opo)
		}
	}

	if len(finishedOPOs) < 1 {
		// skip updating store
		return nil
	}

	dk.Logger(ctx).Info(
		"DymNS hook After-Epoch-End: processing finished OPOs", "count", len(finishedOPOs),
		"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
	)

	for name, newExpiry := range updateExpiryForOPONames {
		aope.ExpiryByName[name] = newExpiry
	}
	updateExpiryForOPONames = nil

	for _, opo := range finishedOPOs {
		delete(aope.ExpiryByName, opo.Name) // delete the processed entry

		if opo.HighestBid != nil {
			if err := dk.CompletePurchaseOrder(ctx, opo.Name); err != nil {
				dk.Logger(ctx).Error(
					"DymNS hook After-Epoch-End: failed to complete purchase order",
					"name", opo.Name, "expiry", opo.ExpireAt, "now", nowEpochUTC,
					"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
					"error", err,
				)
				return err
			}
			continue
		}

		// no bid placed, it just a normal expiry without winner,
		// in this case, just move to history
		if err := dk.MoveOpenPurchaseOrderToHistorical(ctx, opo.Name); err != nil {
			dk.Logger(ctx).Error(
				"DymNS hook After-Epoch-End: failed to move expired purchase order to historical",
				"name", opo.Name, "expiry", opo.ExpireAt, "now", nowEpochUTC,
				"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
				"error", err,
			)
			return err
		}
	}

	if err := dk.SetActiveOpenPurchaseOrdersExpiration(ctx, aope); err != nil {
		dk.Logger(ctx).Error(
			"DymNS hook After-Epoch-End: failed to update active OPO expiry",
			"epoch-number", epochNumber, "epoch-identifier", epochIdentifier,
			"error", err,
		)
		return err
	}

	return nil
}
