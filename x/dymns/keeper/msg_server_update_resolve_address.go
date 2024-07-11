package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

// TODO DymNS: record reverse resolve

// TODO DymNS: fallback to main chain if chain-id is RollApp and not exists in configuration

func (k msgServer) UpdateResolveAddress(goCtx context.Context, msg *dymnstypes.MsgUpdateResolveAddress) (*dymnstypes.MsgUpdateResolveAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dymName, err := k.validateUpdateResolveAddress(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, newConfig := msg.GetDymNameConfig()
	if newConfig.ChainId == ctx.ChainID() {
		newConfig.ChainId = ""
	}
	newConfigIdentity := newConfig.GetIdentity()

	existingConfigCount := len(dymName.Configs)
	if newConfig.IsDelete() {
		var foundSameConfigIdAtIdx = -1
		for i, config := range dymName.Configs {
			if config.GetIdentity() == newConfigIdentity {
				foundSameConfigIdAtIdx = i
				break
			}
		}

		if foundSameConfigIdAtIdx < 0 {
			// no-config case also falls into this branch

			// do nothing
		} else {
			if existingConfigCount == 1 {
				dymName.Configs = nil
			} else {
				dymName.Configs[foundSameConfigIdAtIdx] = dymName.Configs[existingConfigCount-1]
				dymName.Configs = dymName.Configs[:existingConfigCount-1]
			}
		}
	} else {
		if existingConfigCount > 0 {
			var foundSameConfigId bool
			for i, config := range dymName.Configs {
				if config.GetIdentity() == newConfigIdentity {
					dymName.Configs[i] = newConfig
					foundSameConfigId = true
					break
				}
			}
			if !foundSameConfigId {
				dymName.Configs = append(dymName.Configs, newConfig)
			}
		} else {
			dymName.Configs = []dymnstypes.DymNameConfig{newConfig}
		}
	}

	if err := k.SetDymName(ctx, *dymName); err != nil {
		return nil, err
	}

	return &dymnstypes.MsgUpdateResolveAddressResponse{}, nil
}

func (k msgServer) validateUpdateResolveAddress(ctx sdk.Context, msg *dymnstypes.MsgUpdateResolveAddress) (*dymnstypes.DymName, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	dymName := k.GetDymName(ctx, msg.Name)
	if dymName == nil {
		return nil, dymnstypes.ErrDymNameNotFound.Wrap(msg.Name)
	}

	if dymName.IsExpiredAt(ctx.BlockTime()) {
		return nil, sdkerrors.ErrUnauthorized.Wrap("Dym-Name is already expired")
	}

	if dymName.Controller != msg.Controller {
		if dymName.Owner == msg.Controller {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf(
				"please use controller account '%s' to configure", dymName.Controller,
			)
		}

		return nil, sdkerrors.ErrUnauthorized
	}

	return dymName, nil
}
