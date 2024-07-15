package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgPurchaseName_ValidateBasic(t *testing.T) {
	validOffer := dymnsutils.TestCoin(100)

	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name            string
		DymName         string
		Offer           sdk.Coin
		Buyer           string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "valid",
			DymName: "abc",
			Offer:   validOffer,
			Buyer:   "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		},
		{
			name:            "not allow missing name",
			DymName:         "",
			Offer:           validOffer,
			Buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "bad name",
			DymName:         "-a",
			Offer:           validOffer,
			Buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "name is not a valid dym name",
		},
		{
			name:            "missing offer",
			DymName:         "abc",
			Buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "invalid offer",
		},
		{
			name:            "offer can not be zero",
			DymName:         "abc",
			Offer:           dymnsutils.TestCoin(0),
			Buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "offer must be positive",
		},
		{
			name:    "offer can not be negative",
			DymName: "abc",
			Offer: sdk.Coin{
				Denom:  params.BaseDenom,
				Amount: sdk.NewInt(-1),
			},
			Buyer:           "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			wantErr:         true,
			wantErrContains: "invalid offer",
		},
		{
			name:            "missing buyer",
			DymName:         "abc",
			Offer:           validOffer,
			Buyer:           "",
			wantErr:         true,
			wantErrContains: "buyer is not a valid bech32 account address",
		},
		{
			name:            "invalid buyer",
			DymName:         "abc",
			Offer:           validOffer,
			Buyer:           "dym1fl48vsnmsdzcv",
			wantErr:         true,
			wantErrContains: "buyer is not a valid bech32 account address",
		},
		{
			name:            "buyer must be dym1",
			DymName:         "abc",
			Offer:           validOffer,
			Buyer:           "nim1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3pklgjx",
			wantErr:         true,
			wantErrContains: "buyer is not a valid bech32 account address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MsgPurchaseName{
				Name:  tt.DymName,
				Offer: tt.Offer,
				Buyer: tt.Buyer,
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
