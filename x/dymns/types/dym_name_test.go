package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDymName_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*DymName)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection ALL
	tests := []struct {
		name            string
		Name            string
		Owner           string
		Controller      string
		ExpireAt        int64
		Configs         []DymNameConfig
		wantErr         bool
		wantErrContains string
	}{
		{
			name:       "valid dym name",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
		},
		{
			name:       "empty name",
			Name:       "",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "name is empty",
		},
		{
			name:       "bad name",
			Name:       "-a",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:       "empty owner",
			Name:       "bonded-pool",
			Owner:      "",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "owner is empty",
		},
		{
			name:       "bad owner",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "owner is not a valid bech32 account address",
		},
		{
			name:       "empty controller",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "controller is empty",
		},
		{
			name:       "bad controller",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "controller is not a valid bech32 account address",
		},
		{
			name:       "empty expire at",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   0,
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
			},
			wantErr:         true,
			wantErrContains: "expire at is empty",
		},
		{
			name:       "valid dym name without config",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
		},
		{
			name:       "bad config",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "",
					Value: "dym1fl48vsnmsdzcv85q5d2",
				},
			},
			wantErr:         true,
			wantErrContains: "dym name config value must be a valid bech32 account address",
		},
		{
			name:       "duplicate config",
			Name:       "bonded-pool",
			Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			ExpireAt:   time.Now().Unix(),
			Configs: []DymNameConfig{
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
				},
				{
					Type:  DymNameConfigType_NAME,
					Path:  "www",
					Value: "dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
				},
			},
			wantErr:         true,
			wantErrContains: "dym name config is not unique",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DymName{
				Name:       tt.Name,
				Owner:      tt.Owner,
				Controller: tt.Controller,
				ExpireAt:   tt.ExpireAt,
				Configs:    tt.Configs,
			}
			err := m.Validate()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

//goland:noinspection SpellCheckingInspection
func TestDymNameConfig_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*DymNameConfig)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		Type            DymNameConfigType
		ChainId         string
		Path            string
		Value           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "valid name config",
			Type:    DymNameConfigType_NAME,
			ChainId: "dymension_1100-1",
			Path:    "abc",
			Value:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:  "valid name config with multi-level path",
			Type:  DymNameConfigType_NAME,
			Path:  "abc.def",
			Value: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:    "valid name config with empty path",
			Type:    DymNameConfigType_NAME,
			ChainId: "dymension_1100-1",
			Path:    "",
			Value:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:    "valid name config with empty chain-id",
			Type:    DymNameConfigType_NAME,
			ChainId: "",
			Path:    "abc",
			Value:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:    "valid name config with empty chain-id and path",
			Type:    DymNameConfigType_NAME,
			ChainId: "",
			Path:    "",
			Value:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "not accept unknown type",
			Type:            DymNameConfigType_UNKNOWN,
			ChainId:         "",
			Path:            "",
			Value:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config type is not",
		},
		{
			name:            "bad chain-id",
			Type:            DymNameConfigType_NAME,
			ChainId:         "dymension_",
			Path:            "abc",
			Value:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config chain id must be a valid chain id format",
		},
		{
			name:            "bad path",
			Type:            DymNameConfigType_NAME,
			ChainId:         "",
			Path:            "-a",
			Value:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config path must be a valid dym name",
		},
		{
			name:            "bad multi-level path",
			Type:            DymNameConfigType_NAME,
			ChainId:         "",
			Path:            "a.b.",
			Value:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "dym name config path must be a valid dym name",
		},
		{
			name:    "value can be empty",
			Type:    DymNameConfigType_NAME,
			ChainId: "",
			Path:    "a",
			Value:   "",
		},
		{
			name:    "value can be empty",
			Type:    DymNameConfigType_NAME,
			ChainId: "",
			Path:    "",
			Value:   "",
		},
		{
			name:            "bad value",
			Type:            DymNameConfigType_NAME,
			ChainId:         "",
			Path:            "a",
			Value:           "0x01",
			wantErr:         true,
			wantErrContains: "dym name config value must be a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DymNameConfig{
				Type:    tt.Type,
				ChainId: tt.ChainId,
				Path:    tt.Path,
				Value:   tt.Value,
			}

			err := m.Validate()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOwnedDymNames_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*OwnedDymNames)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		DymNames        []string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "valid owned dym name",
			DymNames: []string{"bonded-pool", "not-bonded-pool"},
		},
		{
			name:     "allow empty",
			DymNames: []string{},
		},
		{
			name:            "bad dym name",
			DymNames:        []string{"bonded-pool", "-not-bonded-pool"},
			wantErr:         true,
			wantErrContains: "owned dym name is not a valid dym name:",
		},
		{
			name:            "bad dym name",
			DymNames:        []string{"-a"},
			wantErr:         true,
			wantErrContains: "owned dym name is not a valid dym name:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &OwnedDymNames{
				DymNames: tt.DymNames,
			}

			err := m.Validate()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDymName_IsExpiredAt(t *testing.T) {
	require.False(t, DymName{
		ExpireAt: time.Now().Unix(),
	}.IsExpiredAt(time.Now().UTC().Add(-time.Second)))
	require.True(t, DymName{
		ExpireAt: time.Now().Unix(),
	}.IsExpiredAt(time.Now().UTC().Add(time.Second)))
}

func TestDymName_GetSdkEvent(t *testing.T) {
	event := DymName{
		Name:       "a",
		Owner:      "b",
		Controller: "c",
		ExpireAt:   time.Date(2024, 01, 02, 03, 04, 05, 0, time.UTC).Unix(),
		Configs:    []DymNameConfig{{}, {}},
	}.GetSdkEvent()
	require.NotNil(t, event)
	require.Equal(t, EventTypeSetDymName, event.Type)
	require.Len(t, event.Attributes, 5)
	require.Equal(t, AttributeKeyDymName, event.Attributes[0].Key)
	require.Equal(t, "a", event.Attributes[0].Value)
	require.Equal(t, AttributeKeyDymNameOwner, event.Attributes[1].Key)
	require.Equal(t, "b", event.Attributes[1].Value)
	require.Equal(t, AttributeKeyDymNameController, event.Attributes[2].Key)
	require.Equal(t, "c", event.Attributes[2].Value)
	require.Equal(t, AttributeKeyDymNameExpiryEpoch, event.Attributes[3].Key)
	require.Equal(t, "1704164645", event.Attributes[3].Value)
	require.Equal(t, AttributeKeyDymNameConfigCount, event.Attributes[4].Key)
	require.Equal(t, "2", event.Attributes[4].Value)
}

func TestDymNameConfig_IsDelete(t *testing.T) {
	require.True(t, DymNameConfig{
		Value: "",
	}.IsDelete(), "if value is empty then it's delete")
	require.False(t, DymNameConfig{
		Value: "1",
	}.IsDelete(), "if value is not empty then it's not delete")
}
