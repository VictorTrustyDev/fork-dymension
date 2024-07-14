package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorePrefixes(t *testing.T) {
	t.Run("ensure key prefixes are not mistakenly modified", func(t *testing.T) {
		require.Equal(t, []byte{0x01}, KeyPrefixDymName, "do not change it, will break the app")
		require.Equal(t, []byte{0x02}, KeyPrefixRvlDymNameOwnedByAccount, "do not change it, will break the app")
		require.Equal(t, []byte{0x03}, KeyPrefixOpenPurchaseOrder, "do not change it, will break the app")
		require.Equal(t, []byte{0x05}, KeyPrefixHistoricalOpenPurchaseOrders, "do not change it, will break the app")
		require.Equal(t, []byte{0x06}, KeyPrefixMinExpiryHistoricalOpenPurchaseOrders, "do not change it, will break the app")
	})

	t.Run("ensure key are not mistakenly modified", func(t *testing.T) {
		require.Equal(t, []byte{0x04}, KeyActiveOpenPurchaseOrdersExpiration, "do not change it, will break the app")
	})
}

//goland:noinspection SpellCheckingInspection
func TestKeys(t *testing.T) {
	for _, dymName := range []string{"a", "b", "bonded-pool"} {
		t.Run(dymName, func(t *testing.T) {
			require.Equal(t, append(KeyPrefixDymName, []byte(dymName)...), DymNameKey(dymName))
			require.Equal(t, append(KeyPrefixOpenPurchaseOrder, []byte(dymName)...), OpenPurchaseOrderKey(dymName))
			require.Equal(t, append(KeyPrefixHistoricalOpenPurchaseOrders, []byte(dymName)...), HistoricalOpenPurchaseOrdersKey(dymName))
			require.Equal(t, append(KeyPrefixMinExpiryHistoricalOpenPurchaseOrders, []byte(dymName)...), MinExpiryHistoricalOpenPurchaseOrdersKey(dymName))
		})
	}

	for _, bech32Address := range []string{
		"dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		"dym1gtcunp63a3aqypr250csar4devn8fjpqulq8d4",
		"dym1tygms3xhhs3yv487phx3dw4a95jn7t7lnxec2d",
	} {
		t.Run(bech32Address, func(t *testing.T) {
			accAddr := sdk.MustAccAddressFromBech32(bech32Address)
			require.Equal(t, append(KeyPrefixRvlDymNameOwnedByAccount, accAddr.Bytes()...), DymNamesOwnedByAccountRvlKey(accAddr))
		})
	}
}
