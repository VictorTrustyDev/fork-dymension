package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsValidBech32AccountAddress(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                            string
		address                         string
		matchAccountAddressBech32Prefix bool
		want                            bool
	}{
		{
			name:                            "valid bech32 account address",
			address:                         "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			matchAccountAddressBech32Prefix: true,
			want:                            true,
		},
		{
			name:                            "bad checksum bech32 account address",
			address:                         "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9feu",
			matchAccountAddressBech32Prefix: true,
			want:                            false,
		},
		{
			name:                            "bad bech32 account address",
			address:                         "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3",
			matchAccountAddressBech32Prefix: true,
			want:                            false,
		},
		{
			name:                            "not bech32 address",
			address:                         "0x4fea76427b8345861e80a3540a8a9d936fd39391",
			matchAccountAddressBech32Prefix: true,
			want:                            false,
		},
		{
			name:                            "not bech32 address",
			address:                         "0x4fea76427b8345861e80a3540a8a9d936fd39391",
			matchAccountAddressBech32Prefix: false,
			want:                            false,
		},
		{
			name:                            "valid bech32 account address but mis-match HRP",
			address:                         "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			matchAccountAddressBech32Prefix: true,
			want:                            false,
		},
		{
			name:                            "valid bech32 account address ignore mis-match HRP",
			address:                         "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			matchAccountAddressBech32Prefix: false,
			want:                            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsValidBech32AccountAddress(tt.address, tt.matchAccountAddressBech32Prefix))
		})
	}
}
