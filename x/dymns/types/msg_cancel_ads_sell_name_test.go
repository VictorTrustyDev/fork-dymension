package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgCancelAdsSellName_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		DymName         string
		Owner           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "valid",
			DymName: "abc",
			Owner:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "not allow empty name",
			DymName:         "",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "not allow invalid name",
			DymName:         "-a",
			Owner:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "invalid owner",
			DymName:         "a",
			Owner:           "dym1fl48vsnmsdzcv85q5",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "missing owner",
			DymName:         "a",
			Owner:           "",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:            "owner must be dym1",
			DymName:         "a",
			Owner:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgCancelAdsSellName{
				Name:  tt.DymName,
				Owner: tt.Owner,
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
