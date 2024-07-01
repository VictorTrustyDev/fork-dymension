package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorePrefixes(t *testing.T) {
	t.Run("ensure key prefixes are not mistakenly modified", func(t *testing.T) {
		require.Equal(t, []byte{0x01}, KeyPrefixDymName)
		require.Equal(t, []byte{0x02}, KeyPrefixDymNameOwnedByAccount)
		require.Equal(t, []byte{0x03}, KeyPrefixOpenPurchaseOrder)
		require.Equal(t, []byte{0x04}, KeyPrefixHistoricalOpenPurchaseOrders)
	})
}

func TestDymNameKey(t *testing.T) {
	//goland:noinspection ALL
	tests := []struct {
		dymName string
		want    string
	}{
		{
			dymName: "a",
			want:    "013ac225168df54212a25c1c01fd35bebfea408fdac2e31ddd6f80a4bbf9a5f1cb",
		},
		{
			dymName: "b",
			want:    "01b5553de315e0edf504d9150af82dafa5c4667fa618ed0a6f19c69b41166c5510",
		},
		{
			dymName: "bonded-pool",
			want:    "01fc965ffded6bec70c93eb105ac8eb1678697e1d4bb50e05d5f21f3c02bea4993",
		},
	}
	for _, tt := range tests {
		t.Run(tt.dymName, func(t *testing.T) {
			require.Equal(t, tt.want, hex.EncodeToString(DymNameKey(tt.dymName)))
		})
	}
}
