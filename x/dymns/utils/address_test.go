package utils

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
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
			name:                            "valid bech32 account address, Interchain Account",
			address:                         "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul",
			matchAccountAddressBech32Prefix: true,
			want:                            true,
		},
		{
			name: "reject bech32 account address which neither 20 nor 32 bytes, case 19 bytes",
			address: func() string {
				bz := make([]byte, 19)
				addr, err := bech32.ConvertAndEncode("dym", bz)
				require.NoError(t, err)
				return addr
			}(),
			want: false,
		},
		{
			name: "reject bech32 account address which neither 20 nor 32 bytes, case 21 bytes",
			address: func() string {
				bz := make([]byte, 21)
				addr, err := bech32.ConvertAndEncode("dym", bz)
				require.NoError(t, err)
				return addr
			}(),
			want: false,
		},
		{
			name: "reject bech32 account address which neither 20 nor 32 bytes, case 31 bytes",
			address: func() string {
				bz := make([]byte, 31)
				addr, err := bech32.ConvertAndEncode("dym", bz)
				require.NoError(t, err)
				return addr
			}(),
			want: false,
		},
		{
			name: "reject bech32 account address which neither 20 nor 32 bytes, case 33 bytes",
			address: func() string {
				bz := make([]byte, 33)
				addr, err := bech32.ConvertAndEncode("dym", bz)
				require.NoError(t, err)
				return addr
			}(),
			want: false,
		},
		{
			name: "reject bech32 account address which neither 20 nor 32 bytes, case 128 bytes",
			address: func() string {
				bz := make([]byte, 128)
				addr, err := bech32.ConvertAndEncode("dym", bz)
				require.NoError(t, err)
				return addr
			}(),
			want: false,
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
			gotValid := IsValidBech32AccountAddress(tt.address, tt.matchAccountAddressBech32Prefix)
			require.Equal(t, tt.want, gotValid)
		})
	}
}

//goland:noinspection SpellCheckingInspection
func TestIsValid0xAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "allow 0x hex address with 20 bytes",
			address: "0x1234567890123456789012345678901234567890",
			want:    true,
		},
		{
			name:    "allow 0x hex address with 32 bytes, Interchain Account",
			address: "0x1234567890123456789012345678901234567890123456789012345678901234",
			want:    true,
		},
		{
			name:    "disallow 0x hex address with 19 bytes",
			address: "0x123456789012345678901234567890123456789",
			want:    false,
		},
		{
			name:    "disallow 0x hex address with 21 bytes",
			address: "0x12345678901234567890123456789012345678901",
			want:    false,
		},
		{
			name:    "disallow 0x hex address with 31 bytes",
			address: "0x123456789012345678901234567890123456789012345678901234567890123",
			want:    false,
		},
		{
			name:    "disallow 0x hex address with 33 bytes",
			address: "0x12345678901234567890123456789012345678901234567890123456789012345",
			want:    false,
		},
		{
			name:    "disallow valid bech32 address",
			address: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			want:    false,
		},
		{
			name:    "disallow valid bech32 address, Interchain Account",
			address: "dym1zg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg69v7yszg6qrz80ul",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsValid0xAddress(tt.address))
		})
	}
}