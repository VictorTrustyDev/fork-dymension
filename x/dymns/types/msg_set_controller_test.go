package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgSetController_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		DymName         string
		Controller      string
		Owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:       "valid",
			DymName:    "a",
			Controller: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "controller and owner can be the same",
			DymName:    "a",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "controller and owner can be the different",
			DymName:    "a",
			Controller: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "missing name",
			DymName:         "",
			Controller:      "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "missing controller",
			DymName:         "a",
			Controller:      "",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "controller is not a valid bech32 account address",
		},
		{
			name:            "invalid controller",
			DymName:         "a",
			Controller:      "dym1tygms3xhhs3yv487phx",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "controller is not a valid bech32 account address",
		},
		{
			name:            "missing owner",
			DymName:         "a",
			Controller:      "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "invalid owner",
			DymName:         "a",
			Controller:      "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "dym1fl48vsnmsdzcv85q5d2",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "controller must be dym1",
			DymName:         "a",
			Controller:      "nim1tygms3xhhs3yv487phx3dw4a95jn7t7l4kreyj",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "controller is not a valid bech32 account address",
		},
		{
			name:            "owner must be dym1",
			DymName:         "a",
			Controller:      "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgSetController{
				Name:       tt.DymName,
				Controller: tt.Controller,
				Owner:      tt.Owner,
			}

			err := m.ValidateBasic()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
