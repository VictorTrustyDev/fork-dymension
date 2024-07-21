package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgUpdateResolveAddress_ValidateBasic(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		DymName         string
		ChainId         string
		SubName         string
		ResolveTo       string
		Controller      string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:       "valid",
			DymName:    "a",
			ChainId:    "dymension_1100-1",
			SubName:    "abc",
			ResolveTo:  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "missing dym-name",
			DymName:         "",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "bad dym-name",
			DymName:         "",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:       "valid config resolve with multi-level sub-name",
			DymName:    "a",
			SubName:    "abc.def",
			ResolveTo:  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "valid config resolve without sub-name",
			DymName:    "a",
			ChainId:    "dymension_1100-1",
			SubName:    "",
			ResolveTo:  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "valid config resolve with empty chain-id",
			DymName:    "a",
			ChainId:    "",
			SubName:    "abc",
			ResolveTo:  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "valid config resolve with empty chain-id and sub-name",
			DymName:    "a",
			ChainId:    "",
			SubName:    "",
			ResolveTo:  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "bad chain-id",
			DymName:         "a",
			ChainId:         "dymension_",
			SubName:         "abc",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config chain id must be a valid chain id format",
		},
		{
			name:            "bad sub-name",
			DymName:         "a",
			ChainId:         "",
			SubName:         "-a",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config path must be a valid dym name",
		},
		{
			name:            "bad sub-name, too long",
			DymName:         "a",
			ChainId:         "",
			SubName:         "123456789012345678901",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: ErrDymNameTooLong.Error(),
		},
		{
			name:            "bad multi-level sub-name",
			DymName:         "a",
			ChainId:         "",
			SubName:         "a.b.",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config path must be a valid dym name",
		},
		{
			name:       "resolve to can be empty to allow delete",
			DymName:    "a",
			ChainId:    "",
			SubName:    "a",
			ResolveTo:  "",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:       "resolve to can be empty to allow delete",
			DymName:    "a",
			ChainId:    "",
			SubName:    "",
			ResolveTo:  "",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "bad resolve to",
			DymName:         "a",
			ChainId:         "",
			SubName:         "a",
			ResolveTo:       "0x01",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config value must be a valid bech32 account address",
		},
		{
			name:            "resolve must be dym1 format if chain-id is empty",
			DymName:         "a",
			ChainId:         "",
			ResolveTo:       "nim1tygms3xhhs3yv487phx3dw4a95jn7t7l4kreyj",
			Controller:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "resolve address must be a valid bech32 account address on host chain",
		},
		{
			name:       "resolve to can be non-dym1 format if chain-id is not empty",
			DymName:    "a",
			ChainId:    "nim_1122-1",
			ResolveTo:  "nim1tygms3xhhs3yv487phx3dw4a95jn7t7l4kreyj",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "controller must be dym1",
			DymName:         "a",
			ResolveTo:       "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller:      "nim1tygms3xhhs3yv487phx3dw4a95jn7t7l4kreyj",
			wantErr:         true,
			wantErrContains: "controller is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgUpdateResolveAddress{
				Name:       tt.DymName,
				ChainId:    tt.ChainId,
				SubName:    tt.SubName,
				ResolveTo:  tt.ResolveTo,
				Controller: tt.Controller,
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

func TestMsgUpdateResolveAddress_GetDymNameConfig(t *testing.T) {
	tests := []struct {
		name       string
		DymName    string
		ChainId    string
		SubName    string
		ResolveTo  string
		Controller string
		wantName   string
		wantConfig DymNameConfig
	}{
		{
			name:       "assigned correctly",
			DymName:    "a",
			ChainId:    "dymension",
			SubName:    "sub",
			ResolveTo:  "r",
			Controller: "c",
			wantName:   "a",
			wantConfig: DymNameConfig{
				Type:    DymNameConfigType_NAME,
				ChainId: "dymension",
				Path:    "sub",
				Value:   "r",
			},
		},
		{
			name: "all empty",
			wantConfig: DymNameConfig{
				Type: DymNameConfigType_NAME,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgUpdateResolveAddress{
				Name:       tt.DymName,
				ChainId:    tt.ChainId,
				SubName:    tt.SubName,
				ResolveTo:  tt.ResolveTo,
				Controller: tt.Controller,
			}

			gotName, gotConfig := m.GetDymNameConfig()
			require.Equal(t, tt.wantName, gotName)
			require.Equal(t, tt.wantConfig, gotConfig)
		})
	}
}
