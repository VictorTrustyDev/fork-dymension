package types

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestDymName_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*DymName)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection SpellCheckingInspection
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

func TestDymNameConfig_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*DymNameConfig)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection SpellCheckingInspection
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
			name:            "not accept hex address value",
			Type:            DymNameConfigType_NAME,
			ChainId:         "",
			Path:            "",
			Value:           "0x1234567890123456789012345678901234567890",
			wantErr:         true,
			wantErrContains: "must be a valid bech32 account address",
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
		{
			name:            "reject value not normalized",
			Type:            DymNameConfigType_NAME,
			ChainId:         "",
			Path:            "",
			Value:           "Dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "must be lowercase",
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

func TestReverseLookupDymNames_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*ReverseLookupDymNames)(nil)
		require.Error(t, m.Validate())
	})

	tests := []struct {
		name            string
		DymNames        []string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "valid reverse lookup record",
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
			wantErrContains: "invalid dym name:",
		},
		{
			name:            "bad dym name",
			DymNames:        []string{"-a"},
			wantErr:         true,
			wantErrContains: "invalid dym name:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ReverseLookupDymNames{
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

func TestDymName_IsExpiredAtEpoch(t *testing.T) {
	require.False(t, DymName{
		ExpireAt: time.Now().Unix(),
	}.IsExpiredAtEpoch(time.Now().UTC().Add(-time.Second).Unix()))
	require.True(t, DymName{
		ExpireAt: time.Now().Unix(),
	}.IsExpiredAtEpoch(time.Now().UTC().Add(time.Second).Unix()))
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

func TestDymNameConfig_GetIdentity(t *testing.T) {
	tests := []struct {
		name    string
		_type   DymNameConfigType
		chainId string
		path    string
		value   string
		want    string
	}{
		{
			name:    "combination of Type & Chain Id & Path, exclude Value",
			_type:   DymNameConfigType_NAME,
			chainId: "1",
			path:    "2",
			value:   "3",
			want:    "name|1|2",
		},
		{
			name:    "combination of Type & Chain Id & Path, exclude Value",
			_type:   DymNameConfigType_NAME,
			chainId: "1",
			path:    "2",
			value:   "",
			want:    "name|1|2",
		},
		{
			name:    "normalize material fields",
			_type:   DymNameConfigType_NAME,
			chainId: "AaA",
			path:    "bBb",
			value:   "",
			want:    "name|aaa|bbb",
		},
		{
			name:    "use String() of type",
			_type:   DymNameConfigType_UNKNOWN,
			chainId: "1",
			path:    "2",
			want:    "unknown|1|2",
		},
		{
			name:    "use String() of type",
			_type:   DymNameConfigType_NAME,
			chainId: "1",
			path:    "2",
			want:    "name|1|2",
		},
		{
			name:    "respect empty chain-id",
			_type:   DymNameConfigType_NAME,
			chainId: "",
			path:    "2",
			want:    "name||2",
		},
		{
			name:    "respect empty path",
			_type:   DymNameConfigType_NAME,
			chainId: "1",
			path:    "",
			want:    "name|1|",
		},
		{
			name:    "respect empty chain-id and path",
			_type:   DymNameConfigType_NAME,
			chainId: "",
			path:    "",
			want:    "name||",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DymNameConfig{
				Type:    tt._type,
				ChainId: tt.chainId,
				Path:    tt.path,
				Value:   tt.value,
			}
			require.Equal(t, tt.want, m.GetIdentity())
		})
	}

	t.Run("normalize material fields", func(t *testing.T) {
		require.Equal(t, DymNameConfig{
			ChainId: "AaA",
			Path:    "bBb",
			Value:   "123",
		}.GetIdentity(), DymNameConfig{
			ChainId: "aAa",
			Path:    "BbB",
			Value:   "456",
		}.GetIdentity())
	})
}

func TestDymNameConfig_IsDelete(t *testing.T) {
	require.True(t, DymNameConfig{
		Value: "",
	}.IsDelete(), "if value is empty then it's delete")
	require.False(t, DymNameConfig{
		Value: "1",
	}.IsDelete(), "if value is not empty then it's not delete")
}

//goland:noinspection SpellCheckingInspection
func TestDymName_GetAddressReverseMappingRecords(t *testing.T) {
	const dymName = "a"
	const ownerBech32 = "dym1zg69v7yszg69v7yszg69v7yszg69v7ys8xdv96"
	const ownerBech32AtNim = "nim1zg69v7yszg69v7yszg69v7yszg69v7yspkhdt9"
	const ownerHex = "0x1234567890123456789012345678901234567890"
	const bondedPoolBech32 = "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue"
	const bondedPoolHex = "0x4fea76427b8345861e80a3540a8a9d936fd39391"

	const icaBech32 = "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul"
	const icaBech32AtNim = "nim1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qe9zz9m"
	const icaHex = "0x1234567890123456789012345678901234567890123456789012345678901234"

	tests := []struct {
		name                                 string
		configs                              []DymNameConfig
		customFuncCheckChainIdIsCoinType60   func(string) bool
		wantPanic                            bool
		wantConfiguredAddressesToDymNames    map[string]ReverseLookupDymNames
		wantCoinType60HexAddressesToDymNames map[string]ReverseLookupDymNames
	}{
		{
			name: "pass",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - hex address is parsed correctly",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - configured bech32 address is kept as is",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim, // not dym1, it's nim1
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - hex address is distinct",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim, // not dym1, it's nim1, but still the owner
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: { // only one
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - able to detect default config address when not configured",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
				// not include default config
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: { // default config resolved to owner
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: { // default config resolved to owner
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - respect default config when it is not owner",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   bondedPoolBech32, // not the owner
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				bondedPoolBech32: { // respect
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				bondedPoolHex: { // respect
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - respect default config when it is not owner",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   bondedPoolBech32, // not the owner
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "a",
					Value:   bondedPoolBech32,
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				bondedPoolBech32: { // respect
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				bondedPoolHex: { // respect
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - respect default config when it is not owner",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   bondedPoolBech32, // not owner
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim, // but this is owner, in different bech32 prefix
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - chains not coin-type-60 will not have hex records",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "cosmoshub-4",
					Path:    "",
					Value:   "cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r",
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
			},
			customFuncCheckChainIdIsCoinType60: func(_ string) bool {
				return false // no chain is coin-type-60
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				"cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r": {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - chains not coin-type-60 will not have hex records, mixed with chain that is coin-type-60",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "cosmoshub-4",
					Path:    "",
					Value:   "cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r",
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   "dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
			},
			customFuncCheckChainIdIsCoinType60: func(chaiId string) bool {
				return chaiId == "nim_1122-1" // only NIM is coin-type-60
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				"cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r": {
					DymNames: []string{dymName},
				},
				"dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4": {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
				"0x42f1c98751ec7a02046aa3f10e8eadcb2674c820": {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - host chain is coin-type-60, regardless of func check chain id is coin-type-60 or not",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "cosmoshub-4",
					Path:    "",
					Value:   "cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r",
				},
			},
			customFuncCheckChainIdIsCoinType60: func(_ string) bool {
				return false // no chain is coin-type-60
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				"cosmos1tygms3xhhs3yv487phx3dw4a95jn7t7lpm470r": {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name:      "fail - not accept malformed config",
			configs:   []DymNameConfig{{}},
			wantPanic: true,
		},
		{
			name: "fail - not accept malformed config, not bech32 address value",
			configs: []DymNameConfig{{
				Type:    DymNameConfigType_NAME,
				ChainId: "",
				Path:    "a",
				Value:   "0x1234567890123456789012345678901234567890",
			}},
			wantPanic: true,
		},
		{
			name: "fail - not accept malformed config, default config is not bech32 address of host",
			configs: []DymNameConfig{{
				Type:    DymNameConfigType_NAME,
				ChainId: "",
				Path:    "",
				Value:   ownerBech32AtNim,
			}},
			wantPanic: true,
		},
		{
			name: "fail - not accept malformed config, not valid bech32 address",
			configs: []DymNameConfig{{
				Type:    DymNameConfigType_NAME,
				ChainId: "",
				Path:    "a",
				Value:   ownerBech32 + "a",
			}},
			wantPanic: true,
		},
		{
			name: "fail - not accept malformed config, default config is not bech32 address of host",
			configs: []DymNameConfig{{
				Type:    DymNameConfigType_NAME,
				ChainId: "",
				Path:    "",
				Value:   ownerBech32 + "a",
			}},
			wantPanic: true,
		},
		{
			name: "pass - ignore empty value config",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   ownerBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   "", // empty value
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32: {
					DymNames: []string{dymName},
				},
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
			},
		},
		{
			name: "pass - allow Interchain Account",
			configs: []DymNameConfig{
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "",
					Value:   icaBech32,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "ica",
					Value:   icaBech32AtNim,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "nim_1122-1",
					Path:    "",
					Value:   ownerBech32AtNim,
				},
				{
					Type:    DymNameConfigType_NAME,
					ChainId: "",
					Path:    "bonded-pool",
					Value:   bondedPoolBech32,
				},
			},
			wantConfiguredAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerBech32AtNim: {
					DymNames: []string{dymName},
				},
				bondedPoolBech32: {
					DymNames: []string{dymName},
				},
				icaBech32: {
					DymNames: []string{dymName},
				},
				icaBech32AtNim: {
					DymNames: []string{dymName},
				},
			},
			wantCoinType60HexAddressesToDymNames: map[string]ReverseLookupDymNames{
				ownerHex: {
					DymNames: []string{dymName},
				},
				bondedPoolHex: {
					DymNames: []string{dymName},
				},
				icaHex: {
					DymNames: []string{dymName},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DymName{
				Name:       dymName,
				Owner:      ownerBech32,
				Controller: ownerBech32,
				ExpireAt:   1,
				Configs:    tt.configs,
			}

			if tt.wantPanic {
				require.Panics(t, func() {
					_, _ = m.GetAddressReverseMappingRecords(func(_ string) bool {
						return true
					})
				})

				return
			}

			useFuncCheckChainIdIsCoinType60 := func(_ string) bool {
				return true
			}
			if tt.customFuncCheckChainIdIsCoinType60 != nil {
				useFuncCheckChainIdIsCoinType60 = tt.customFuncCheckChainIdIsCoinType60
			}

			gotConfiguredAddressesToDymNames, gotCoinType60HexAddressesToDymNames := m.GetAddressReverseMappingRecords(useFuncCheckChainIdIsCoinType60)
			if !reflect.DeepEqual(gotConfiguredAddressesToDymNames, tt.wantConfiguredAddressesToDymNames) {
				t.Errorf("gotConfiguredAddressesToDymNames = %v, want %v", gotConfiguredAddressesToDymNames, tt.wantConfiguredAddressesToDymNames)
			}
			if !reflect.DeepEqual(gotCoinType60HexAddressesToDymNames, tt.wantCoinType60HexAddressesToDymNames) {
				t.Errorf("gotCoinType60HexAddressesToDymNames = %v, want %v", gotCoinType60HexAddressesToDymNames, tt.wantCoinType60HexAddressesToDymNames)
			}
		})
	}

	t.Run("func check chain-id is coin-type 60 chain is required", func(t *testing.T) {
		dymName := &DymName{
			Name:       "a",
			Owner:      ownerBech32,
			Controller: ownerBech32,
			ExpireAt:   1,
		}

		require.NoError(t, dymName.Validate())

		require.Panics(t, func() {
			_, _ = dymName.GetAddressReverseMappingRecords(nil)
		})
	})
}

func TestReverseLookupDymNames_Distinct(t *testing.T) {
	tests := []struct {
		name             string
		providedDymNames []string
		wantDistinct     []string
	}{
		{
			name:             "distinct",
			providedDymNames: []string{"a", "b", "b", "a", "c", "d"},
			wantDistinct:     []string{"a", "b", "c", "d"},
		},
		{
			name:             "distinct of single",
			providedDymNames: []string{"a"},
			wantDistinct:     []string{"a"},
		},
		{
			name:             "empty",
			providedDymNames: []string{},
			wantDistinct:     []string{},
		},
		{
			name:             "nil",
			providedDymNames: nil,
			wantDistinct:     []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ReverseLookupDymNames{
				DymNames: tt.providedDymNames,
			}
			distinct := m.Distinct().DymNames
			want := tt.wantDistinct

			sort.Strings(distinct)
			sort.Strings(want)

			require.Equal(t, want, distinct)
		})
	}
}

func TestReverseLookupDymNames_Combine(t *testing.T) {
	tests := []struct {
		name             string
		providedDymNames []string
		otherDymNames    []string
		wantCombined     []string
	}{
		{
			name:             "combined",
			providedDymNames: []string{"a", "b"},
			otherDymNames:    []string{"c", "d"},
			wantCombined:     []string{"a", "b", "c", "d"},
		},
		{
			name:             "combined, distinct",
			providedDymNames: []string{"a", "b"},
			otherDymNames:    []string{"b", "c", "d"},
			wantCombined:     []string{"a", "b", "c", "d"},
		},
		{
			name:             "combined, distinct",
			providedDymNames: []string{"a"},
			otherDymNames:    []string{"a"},
			wantCombined:     []string{"a"},
		},
		{
			name:             "combine empty with other",
			providedDymNames: nil,
			otherDymNames:    []string{"a"},
			wantCombined:     []string{"a"},
		},
		{
			name:             "combine empty with other",
			providedDymNames: []string{"a"},
			otherDymNames:    nil,
			wantCombined:     []string{"a"},
		},
		{
			name:             "combine empty with other",
			providedDymNames: nil,
			otherDymNames:    []string{"a", "b"},
			wantCombined:     []string{"a", "b"},
		},
		{
			name:             "combine with other empty",
			providedDymNames: []string{"a", "b"},
			otherDymNames:    nil,
			wantCombined:     []string{"a", "b"},
		},
		{
			name:             "distinct source",
			providedDymNames: []string{"a", "b", "a"},
			otherDymNames:    []string{"c", "c", "d"},
			wantCombined:     []string{"a", "b", "c", "d"},
		},
		{
			name:             "both empty",
			providedDymNames: []string{},
			otherDymNames:    []string{},
			wantCombined:     []string{},
		},
		{
			name:             "both nil",
			providedDymNames: nil,
			otherDymNames:    nil,
			wantCombined:     []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ReverseLookupDymNames{
				DymNames: tt.providedDymNames,
			}
			other := ReverseLookupDymNames{
				DymNames: tt.otherDymNames,
			}
			combined := m.Combine(other).DymNames
			want := tt.wantCombined

			sort.Strings(combined)
			sort.Strings(want)

			require.Equal(t, want, combined)
		})
	}
}

func TestReverseLookupDymNames_Exclude(t *testing.T) {
	tests := []struct {
		name                 string
		providedDymNames     []string
		toBeExcludedDymNames []string
		want                 []string
	}{
		{
			name:                 "exclude",
			providedDymNames:     []string{"a", "b", "c", "d"},
			toBeExcludedDymNames: []string{"b", "d"},
			want:                 []string{"a", "c"},
		},
		{
			name:                 "exclude all",
			providedDymNames:     []string{"a", "b", "c", "d"},
			toBeExcludedDymNames: []string{"d", "c", "b", "a"},
			want:                 []string{},
		},
		{
			name:                 "exclude none",
			providedDymNames:     []string{"a", "b", "c", "d"},
			toBeExcludedDymNames: []string{},
			want:                 []string{"a", "b", "c", "d"},
		},
		{
			name:                 "exclude nil",
			providedDymNames:     []string{"a", "b", "c", "d"},
			toBeExcludedDymNames: []string{},
			want:                 []string{"a", "b", "c", "d"},
		},
		{
			name:                 "none exclude",
			providedDymNames:     []string{},
			toBeExcludedDymNames: []string{"a", "b", "c", "d"},
			want:                 []string{},
		},
		{
			name:                 "nil exclude",
			providedDymNames:     nil,
			toBeExcludedDymNames: []string{"a", "b", "c", "d"},
			want:                 []string{},
		},
		{
			name:                 "distinct after exclude",
			providedDymNames:     []string{"a", "a", "b"},
			toBeExcludedDymNames: []string{"b", "d"},
			want:                 []string{"a"},
		},
		{
			name:                 "exclude partial",
			providedDymNames:     []string{"a", "b", "c"},
			toBeExcludedDymNames: []string{"b", "c", "d"},
			want:                 []string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ReverseLookupDymNames{
				DymNames: tt.providedDymNames,
			}
			other := ReverseLookupDymNames{
				DymNames: tt.toBeExcludedDymNames,
			}
			combined := m.Exclude(other).DymNames
			want := tt.want

			sort.Strings(combined)
			sort.Strings(want)

			require.Equal(t, want, combined)
		})
	}
}
