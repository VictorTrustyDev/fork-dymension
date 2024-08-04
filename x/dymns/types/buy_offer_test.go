package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
)

func TestBuyOffer_HasCounterpartyOfferPrice(t *testing.T) {
	require.False(t, (&BuyOffer{
		CounterpartyOfferPrice: nil,
	}).HasCounterpartyOfferPrice())
	require.False(t, (&BuyOffer{
		CounterpartyOfferPrice: &sdk.Coin{},
	}).HasCounterpartyOfferPrice())
	require.False(t, (&BuyOffer{
		CounterpartyOfferPrice: dymnsutils.TestCoinP(0),
	}).HasCounterpartyOfferPrice())
	require.True(t, (&BuyOffer{
		CounterpartyOfferPrice: dymnsutils.TestCoinP(1),
	}).HasCounterpartyOfferPrice())
}

func TestBuyOffer_Validate(t *testing.T) {
	t.Run("nil obj", func(t *testing.T) {
		m := (*BuyOffer)(nil)
		require.Error(t, m.Validate())
	})

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name                   string
		offerId                string
		dymName                string
		_type                  MarketOrderType
		buyer                  string
		offerPrice             sdk.Coin
		counterpartyOfferPrice *sdk.Coin
		wantErr                bool
		wantErrContains        string
	}{
		{
			name:                   "pass - valid offer",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: nil,
		},
		{
			name:                   "pass - valid offer with counterparty offer price",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: dymnsutils.TestCoinP(2),
		},
		{
			name:                   "pass - valid offer without counterparty offer price",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: nil,
		},
		{
			name:            "fail - empty offer ID",
			offerId:         "",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "ID of offer is empty",
		},
		{
			name:            "fail - type is Unknown",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_UNKNOWN,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Buy-Order type must be",
		},
		{
			name:            "fail - type is Alias (not yet supported_",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_ALIAS,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Buy-Order type must be",
		},
		{
			name:            "fail - bad offer ID",
			offerId:         "@",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "ID of offer is not a valid offer id",
		},
		{
			name:            "fail - empty name",
			offerId:         "1",
			dymName:         "",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Dym-Name of offer is empty",
		},
		{
			name:            "fail - bad name",
			offerId:         "1",
			dymName:         "@",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "Dym-Name of offer is not a valid dym name",
		},
		{
			name:            "fail - bad buyer",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "0x1",
			offerPrice:      dymnsutils.TestCoin(1),
			wantErr:         true,
			wantErrContains: "buyer is not a valid bech32 account address",
		},
		{
			name:            "fail - offer price is zero",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(0),
			wantErr:         true,
			wantErrContains: "offer price is zero",
		},
		{
			name:            "fail - offer price is empty",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      sdk.Coin{},
			wantErr:         true,
			wantErrContains: "offer price is zero",
		},
		{
			name:            "fail - offer price is negative",
			offerId:         "1",
			dymName:         "a",
			_type:           MarketOrderType_MOT_DYM_NAME,
			buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:      dymnsutils.TestCoin(-1),
			wantErr:         true,
			wantErrContains: "offer price is negative",
		},
		{
			name:    "fail - offer price is invalid",
			offerId: "1",
			dymName: "a",
			_type:   MarketOrderType_MOT_DYM_NAME,
			buyer:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice: sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "offer price is invalid",
		},
		{
			name:                   "pass - counter-party offer price is zero",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: dymnsutils.TestCoinP(0),
		},
		{
			name:                   "pass - counter-party offer price is empty",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: &sdk.Coin{},
		},
		{
			name:                   "fail - counter-party offer price is negative",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: dymnsutils.TestCoinP(-1),
			wantErr:                true,
			wantErrContains:        "counterparty offer price is negative",
		},
		{
			name:       "fail - counter-party offer price is invalid",
			offerId:    "1",
			dymName:    "a",
			_type:      MarketOrderType_MOT_DYM_NAME,
			buyer:      "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice: dymnsutils.TestCoin(1),
			counterpartyOfferPrice: &sdk.Coin{
				Denom:  "-",
				Amount: sdk.OneInt(),
			},
			wantErr:         true,
			wantErrContains: "counterparty offer price is invalid",
		},
		{
			name:                   "pass - counterparty offer price can be less than offer price",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(2),
			counterpartyOfferPrice: dymnsutils.TestCoinP(1),
			wantErr:                false,
		},
		{
			name:                   "pass - counterparty offer price can be equals to offer price",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(2),
			counterpartyOfferPrice: dymnsutils.TestCoinP(2),
			wantErr:                false,
		},
		{
			name:                   "pass - counterparty offer price can be greater than offer price",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(2),
			counterpartyOfferPrice: dymnsutils.TestCoinP(3),
			wantErr:                false,
		},
		{
			name:                   "fail - counterparty offer price denom must match offer price denom",
			offerId:                "1",
			dymName:                "a",
			_type:                  MarketOrderType_MOT_DYM_NAME,
			buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			offerPrice:             dymnsutils.TestCoin(1),
			counterpartyOfferPrice: dymnsutils.TestCoin2P(sdk.NewInt64Coin("u"+params.BaseDenom, 2)),
			wantErr:                true,
			wantErrContains:        "counterparty offer price denom is different from offer price denom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BuyOffer{
				Id:                     tt.offerId,
				Name:                   tt.dymName,
				Type:                   tt._type,
				Buyer:                  tt.buyer,
				OfferPrice:             tt.offerPrice,
				CounterpartyOfferPrice: tt.counterpartyOfferPrice,
			}

			err := m.Validate()
			if tt.wantErr {
				require.NotEmpty(t, tt.wantErrContains, "mis-configured test case")
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrContains)
				return
			}

			require.NoError(t, err)
		})
	}
}

//goland:noinspection SpellCheckingInspection
func TestBuyOffer_GetSdkEvent(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		event := BuyOffer{
			Id:                     "1",
			Name:                   "a",
			Buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			OfferPrice:             dymnsutils.TestCoin(1),
			CounterpartyOfferPrice: dymnsutils.TestCoinP(2),
		}.GetSdkEvent("action-name")
		require.NotNil(t, event)
		require.Equal(t, EventTypeBuyOffer, event.Type)
		require.Len(t, event.Attributes, 5)
		require.Equal(t, AttributeKeyBoId, event.Attributes[0].Key)
		require.Equal(t, "1", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyBoName, event.Attributes[1].Key)
		require.Equal(t, "a", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyBoOfferPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyBoCounterpartyOfferPrice, event.Attributes[3].Key)
		require.Equal(t, "2"+params.BaseDenom, event.Attributes[3].Value)
		require.Equal(t, AttributeKeySoActionName, event.Attributes[4].Key)
		require.Equal(t, "action-name", event.Attributes[4].Value)
	})

	t.Run("no counterparty offer price", func(t *testing.T) {
		event := BuyOffer{
			Id:                     "1",
			Name:                   "a",
			Buyer:                  "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			OfferPrice:             dymnsutils.TestCoin(1),
			CounterpartyOfferPrice: nil,
		}.GetSdkEvent("action-name")
		require.NotNil(t, event)
		require.Equal(t, EventTypeBuyOffer, event.Type)
		require.Len(t, event.Attributes, 5)
		require.Equal(t, AttributeKeyBoId, event.Attributes[0].Key)
		require.Equal(t, "1", event.Attributes[0].Value)
		require.Equal(t, AttributeKeyBoName, event.Attributes[1].Key)
		require.Equal(t, "a", event.Attributes[1].Value)
		require.Equal(t, AttributeKeyBoOfferPrice, event.Attributes[2].Key)
		require.Equal(t, "1"+params.BaseDenom, event.Attributes[2].Value)
		require.Equal(t, AttributeKeyBoCounterpartyOfferPrice, event.Attributes[3].Key)
		require.Empty(t, event.Attributes[3].Value)
		require.Equal(t, AttributeKeySoActionName, event.Attributes[4].Key)
		require.Equal(t, "action-name", event.Attributes[4].Value)
	})
}
