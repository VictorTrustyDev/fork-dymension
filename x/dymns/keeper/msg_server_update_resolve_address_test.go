package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testkeeper "github.com/dymensionxyz/dymension/v3/testutil/keeper"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//goland:noinspection SpellCheckingInspection
func Test_msgServer_UpdateResolveAddress(t *testing.T) {
	now := time.Now().UTC()

	const chainId = "dymension_1100-1"

	setupTest := func() (dymnskeeper.Keeper, sdk.Context) {
		dk, _, _, ctx := testkeeper.DymNSKeeper(t)
		ctx = ctx.WithBlockHeader(tmproto.Header{
			Time: now,
		})
		ctx = ctx.WithChainID(chainId)

		return dk, ctx
	}

	const owner = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const controller = "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4"
	const recordName = "bonded-pool"

	tests := []struct {
		name            string
		dymName         *dymnstypes.DymName
		msg             *dymnstypes.MsgUpdateResolveAddress
		wantErr         bool
		wantErrContains string
		wantDymName     *dymnstypes.DymName
	}{
		{
			name: "fail - reject if message not pass validate basic",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg:             &dymnstypes.MsgUpdateResolveAddress{},
			wantErr:         true,
			wantErrContains: dymnstypes.ErrValidationFailed.Error(),
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
		},
		{
			name:    "fail - Dym-Name does not exists",
			dymName: nil,
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr:         true,
			wantErrContains: dymnstypes.ErrDymNameNotFound.Error(),
			wantDymName:     nil,
		},
		{
			name: "fail - reject if Dym-Name expired",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() - 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr:         true,
			wantErrContains: "Dym-Name is already expired",
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() - 1,
			},
		},
		{
			name: "fail - reject if sender is neither owner nor controller",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			},
			wantErr:         true,
			wantErrContains: sdkerrors.ErrUnauthorized.Error(),
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
		},
		{
			name: "fail - reject if sender is owner but not controller",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: owner,
			},
			wantErr:         true,
			wantErrContains: "please use controller account",
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
		},
		{
			name: "fail - reject if config is not valid",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "0x1",
				Controller: controller,
			},
			wantErr:         true,
			wantErrContains: dymnstypes.ErrValidationFailed.Error(),
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
		},
		{
			name: "success - can update",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "",
						Path:    "",
						Value:   owner,
					},
				},
			},
		},
		{
			name: "success - add new record if not exists",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: owner,
					},
				},
			},
		},
		{
			name: "success - override record if exists",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: owner,
					},
				},
			},
		},
		{
			name: "success - remove record if new resolve to empty, single-config",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "",
				SubName:    "a",
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs:    nil,
			},
		},
		{
			name: "success - remove record if new resolve to empty, single-config, not match any",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "",
				SubName:    "non-exists",
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
				},
			},
		},
		{
			name: "success - remove record if new resolve to empty, multi-config, first",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "",
				SubName:    "a",
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
		},
		{
			name: "success - remove record if new resolve to empty, multi-configs, last",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "",
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
				},
			},
		},
		{
			name: "success - remove record if new resolve to empty, multi-config, not any of existing",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  "",
				SubName:    "non-exists",
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "a",
						Value: owner,
					},
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Path:  "",
						Value: controller,
					},
				},
			},
		},
		{
			name: "success - expiry not changed",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 99,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 99,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:  dymnstypes.DymNameConfigType_NAME,
						Value: owner,
					},
				},
			},
		},
		{
			name: "success - chain-id automatically removed from record if is host chain-id",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ChainId:    chainId,
				SubName:    "a",
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "", // empty
						Path:    "a",
						Value:   owner,
					},
				},
			},
		},
		{
			name: "success - chain-id automatically removed from record if is host chain-id",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "", // originally empty
						Path:    "a",
						Value:   controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ChainId:    chainId,
				SubName:    "a",
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "", // empty
						Path:    "a",
						Value:   owner,
					},
				},
			},
		},
		{
			name: "success - chain-id recorded if is NOT host chain-id",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ChainId:    "blumbus_100-1",
				SubName:    "a",
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_100-1",
						Path:    "a",
						Value:   owner,
					},
				},
			},
		},
		{
			name: "success - do not override record with different chain-id",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "",
						Path:    "a",
						Value:   owner,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ChainId:    "blumbus_100-1",
				SubName:    "a",
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "",
						Path:    "a",
						Value:   owner,
					},
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_100-1",
						Path:    "a",
						Value:   owner,
					},
				},
			},
		},
		{
			name: "success - do not override record with different chain-id",
			dymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "",
						Path:    "a",
						Value:   controller,
					},
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_100-1",
						Path:    "a",
						Value:   controller,
					},
				},
			},
			msg: &dymnstypes.MsgUpdateResolveAddress{
				ChainId:    "blumbus_100-1",
				SubName:    "a",
				ResolveTo:  owner,
				Controller: controller,
			},
			wantErr: false,
			wantDymName: &dymnstypes.DymName{
				Owner:      owner,
				Controller: controller,
				ExpireAt:   now.Unix() + 1,
				Configs: []dymnstypes.DymNameConfig{
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "",
						Path:    "a",
						Value:   controller,
					},
					{
						Type:    dymnstypes.DymNameConfigType_NAME,
						ChainId: "blumbus_100-1",
						Path:    "a",
						Value:   owner,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, ctx := setupTest()

			if tt.dymName != nil {
				if tt.dymName.Name == "" {
					tt.dymName.Name = recordName
				}
				err := dk.SetDymName(ctx, *tt.dymName)
				require.NoError(t, err)
			}
			if tt.wantDymName != nil && tt.wantDymName.Name == "" {
				tt.wantDymName.Name = recordName
			}

			tt.msg.Name = recordName
			resp, err := dymnskeeper.NewMsgServerImpl(dk).UpdateResolveAddress(ctx, tt.msg)
			laterDymName := dk.GetDymName(ctx, tt.msg.Name)

			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				require.Nil(t, resp)

				if tt.wantDymName != nil {
					require.Equal(t, *tt.wantDymName, *laterDymName)
				} else {
					require.Nil(t, laterDymName)
				}

				require.Less(t,
					ctx.GasMeter().GasConsumed(), dymnstypes.OpGasConfig,
					"should not consume params gas on failed operation",
				)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, laterDymName)
			require.Equal(t, *tt.wantDymName, *laterDymName)

			require.GreaterOrEqual(t,
				ctx.GasMeter().GasConsumed(), dymnstypes.OpGasConfig,
				"should consume params gas",
			)
		})
	}
}
