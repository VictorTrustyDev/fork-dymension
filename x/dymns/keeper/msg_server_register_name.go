package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"time"
)

func (k msgServer) RegisterName(goCtx context.Context, msg *dymnstypes.MsgRegisterName) (*dymnstypes.MsgRegisterNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName, params, err := k.validateRegisterName(ctx, msg)
	if err != nil {
		return nil, err
	}

	addDurationInSeconds := 86400 * 365 * int64(msg.Duration)

	var firstYearPrice sdkmath.Int
	nameLength := len(msg.Name)
	if nameLength == 1 {
		firstYearPrice = params.Price.Price_1Letter
	} else if nameLength == 2 {
		firstYearPrice = params.Price.Price_2Letters
	} else if nameLength == 3 {
		firstYearPrice = params.Price.Price_3Letters
	} else if nameLength == 4 {
		firstYearPrice = params.Price.Price_4Letters
	} else {
		firstYearPrice = params.Price.Price_5PlusLetters
	}

	var pruneAnyHistoricalData bool
	var totalCost sdk.Coin
	if dymName == nil {
		// register new
		pruneAnyHistoricalData = true

		dymName = &dymnstypes.DymName{
			Name:       msg.Name,
			Owner:      msg.Owner,
			Controller: msg.Owner,
			ExpireAt:   ctx.BlockTime().Unix() + addDurationInSeconds,
			Configs:    nil,
		}

		totalCost = sdk.NewCoin(
			params.Price.PriceDenom,
			firstYearPrice.Add( // first year has different price
				params.Price.PriceExtends.Mul(
					sdkmath.NewInt(int64(
						msg.Duration-1, // subtract first year
					)),
				),
			),
		)
	} else if dymName.Owner == msg.Owner {
		if dymName.IsExpiredAt(ctx.BlockTime()) {
			// renew
			pruneAnyHistoricalData = true

			dymName = &dymnstypes.DymName{
				Name:       msg.Name,
				Owner:      msg.Owner,
				Controller: msg.Owner,
				ExpireAt:   ctx.BlockTime().Unix() + addDurationInSeconds,
				Configs:    nil,
			}
		} else {
			// extends
			pruneAnyHistoricalData = false

			// just add duration, no need to change any existing configuration
			dymName.ExpireAt += addDurationInSeconds
		}

		totalCost = sdk.NewCoin(
			params.Price.PriceDenom,
			params.Price.PriceExtends.Mul(
				sdkmath.NewInt(msg.Duration),
			),
		)
	} else {
		// take over
		pruneAnyHistoricalData = true

		dymName = &dymnstypes.DymName{
			Name:       msg.Name,
			Owner:      msg.Owner,
			Controller: msg.Owner,
			ExpireAt:   ctx.BlockTime().Unix() + addDurationInSeconds,
			Configs:    nil,
		}

		totalCost = sdk.NewCoin(
			params.Price.PriceDenom,
			firstYearPrice.Add( // first year has different price
				params.Price.PriceExtends.Mul(
					sdkmath.NewInt(int64(
						msg.Duration-1, // subtract first year
					)),
				),
			),
		)
	}

	if !totalCost.IsPositive() {
		panic(sdkerrors.ErrLogic.Wrapf("total cost is not positive: %s", totalCost.String()))
	}

	if !totalCost.Equal(msg.ConfirmPayment) {
		return nil, dymnstypes.ErrUnAcknowledgedPayment.Wrapf("Actual payment is %s != %s provided by user", totalCost.String(), msg.ConfirmPayment)
	}

	// At this place we don't do compare actual payment with estimated payment calculated by EstimateRegisterName
	// because in-case there is different between them, it would prevent user to registration/renew.

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx,
		sdk.MustAccAddressFromBech32(msg.Owner),
		dymnstypes.ModuleName,
		sdk.NewCoins(totalCost),
	); err != nil {
		return nil, err
	}

	if err := k.bankKeeper.BurnCoins(ctx, dymnstypes.ModuleName, sdk.NewCoins(totalCost)); err != nil {
		return nil, err
	}

	if pruneAnyHistoricalData {
		if err := k.PruneDymName(ctx, msg.Name); err != nil {
			return nil, err
		}
	}

	if err := k.SetDymName(ctx, *dymName); err != nil {
		return nil, err
	}

	if err := k.AfterDymNameOwnerChanged(ctx, dymName.Name); err != nil {
		return nil, err
	}

	if err := k.AfterDymNameConfigChanged(ctx, dymName.Name); err != nil {
		return nil, err
	}

	return &dymnstypes.MsgRegisterNameResponse{}, nil
}

func (k msgServer) validateRegisterName(ctx sdk.Context, msg *dymnstypes.MsgRegisterName) (*dymnstypes.DymName, *dymnstypes.Params, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	params := k.GetParams(ctx)

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName != nil {
		if dymName.Owner == msg.Owner {
			// just renew or extends
		} else {
			if !dymName.IsExpiredAt(ctx.BlockTime()) {
				return nil, nil, sdkerrors.ErrUnauthorized
			}

			// take over

			// check grace period.
			// Grace period is the time period after the Dym-Name expired
			// that the previous owner can re-purchase the Dym-Name and no one else can take over.
			// This follow domain specification to prevent user mistake.
			gracePeriodInSeconds := int64(params.Misc.DaysGracePeriod) * 86400
			dymNameCanBeTakeOverAfterEpoch := dymName.ExpireAt + gracePeriodInSeconds

			if ctx.BlockTime().Unix() < dymNameCanBeTakeOverAfterEpoch {
				// still in grace period
				return nil, nil, dymnstypes.ErrGracePeriod.Wrapf(
					"can be taken over after %s", time.Unix(dymNameCanBeTakeOverAfterEpoch, 0).UTC().Format(time.DateTime),
				)
			}

			// allowed to take over
		}
	}

	return dymName, &params, nil
}

func EstimateRegisterName(
	params dymnstypes.Params,
	name string,
	existingDymName *dymnstypes.DymName,
	newOwner string,
	duration int64,
) dymnstypes.QueryEstimateRegisterNameResponse {
	getLengthBasedPrice := func() sdkmath.Int {
		switch len(name) {
		case 1:
			return params.Price.Price_1Letter
		case 2:
			return params.Price.Price_2Letters
		case 3:
			return params.Price.Price_3Letters
		case 4:
			return params.Price.Price_4Letters
		default:
			return params.Price.Price_5PlusLetters
		}
	}

	var newFirstYearPrice, extendsPrice sdkmath.Int

	if existingDymName != nil && existingDymName.Owner == newOwner {
		// Dym-Name exists and just renew or extends by the same owner

		newFirstYearPrice = sdk.ZeroInt() // regardless of expired or not, we don't charge this
		extendsPrice = params.Price.PriceExtends.Mul(
			sdkmath.NewInt(duration),
		)
	} else {
		// new registration or take over
		newFirstYearPrice = getLengthBasedPrice() // charge based on name length for the first year
		if duration > 1 {
			extendsPrice = params.Price.PriceExtends.Mul(
				sdkmath.NewInt(duration - 1), // subtract first year, which has different price
			)
		} else {
			extendsPrice = sdk.ZeroInt()
		}
	}

	return dymnstypes.QueryEstimateRegisterNameResponse{
		FirstYearPrice: sdk.NewCoin(params.Price.PriceDenom, newFirstYearPrice),
		ExtendPrice:    sdk.NewCoin(params.Price.PriceDenom, extendsPrice),
		TotalPrice:     sdk.NewCoin(params.Price.PriceDenom, newFirstYearPrice.Add(extendsPrice)),
	}
}