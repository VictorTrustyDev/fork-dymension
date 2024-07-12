package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDefaultGenesis(t *testing.T) {
	defaultGenesis := DefaultGenesis()
	require.NotNil(t, defaultGenesis)
	require.NoError(t, defaultGenesis.Validate())
}

//goland:noinspection SpellCheckingInspection
func TestGenesisState_Validate(t *testing.T) {
	defaultGenesis := DefaultGenesis()
	require.NoError(t, defaultGenesis.Validate())

	t.Run("valid genesis", func(t *testing.T) {
		require.NoError(t, (GenesisState{
			Params: DefaultParams(),
			DymNames: []DymName{
				{
					Name:       "bonded-pool",
					Owner:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					ExpireAt:   time.Now().Unix(),
				},
			},
			OpenPurchaseOrderBids: []OpenPurchaseOrderBid{
				{
					Bidder: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
					Price: sdk.Coin{
						Denom:  "adym",
						Amount: sdk.OneInt(),
					},
				},
			},
		}).Validate())
	})

	t.Run("invalid params", func(t *testing.T) {
		require.Error(t, (GenesisState{
			Params: Params{
				Price: DefaultPriceParams(),
				Misc: MiscParams{
					BeginEpochHookIdentifier: "invalid",
				},
			},
		}).Validate())

		require.Error(t, (GenesisState{
			Params: Params{
				Price: PriceParams{},
				Misc: MiscParams{
					BeginEpochHookIdentifier: "invalid",
				},
			},
		}).Validate())
	})

	t.Run("invalid dym names", func(t *testing.T) {
		require.Error(t, (GenesisState{
			Params: DefaultParams(),
			DymNames: []DymName{
				{
					Name: "",
				},
			},
		}).Validate())
	})

	t.Run("invalid bid", func(t *testing.T) {
		require.Error(t, (GenesisState{
			Params: DefaultParams(),
			OpenPurchaseOrderBids: []OpenPurchaseOrderBid{
				{
					Bidder: "",
				},
			},
		}).Validate())
	})
}
