package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgTransferOwnership_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		DymName         string
		NewOwner        string
		Owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "valid",
			DymName:  "a",
			NewOwner: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:    "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "new owner and owner can not be the same",
			DymName:         "a",
			NewOwner:        "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "new owner must be different from the current owner",
		},
		{
			name:            "missing name",
			DymName:         "",
			NewOwner:        "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "missing new owner",
			DymName:         "a",
			NewOwner:        "",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "new owner is not a valid bech32 account address",
		},
		{
			name:            "invalid new owner",
			DymName:         "a",
			NewOwner:        "dym1tygms3xhhs3yv487phx",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "new owner is not a valid bech32 account address",
		},
		{
			name:            "missing owner",
			DymName:         "a",
			NewOwner:        "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "invalid owner",
			DymName:         "a",
			NewOwner:        "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "dym1fl48vsnmsdzcv85q5d2",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "new owner must be dym1",
			DymName:         "a",
			NewOwner:        "nim1tygms3xhhs3yv487phx3dw4a95jn7t7l4kreyj",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "new owner is not a valid bech32 account address",
		},
		{
			name:            "owner must be dym1",
			DymName:         "a",
			NewOwner:        "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
			Owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgTransferOwnership{
				Name:     tt.DymName,
				NewOwner: tt.NewOwner,
				Owner:    tt.Owner,
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
