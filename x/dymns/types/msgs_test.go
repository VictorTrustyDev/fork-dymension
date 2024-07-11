package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgs_Signers(t *testing.T) {
	t.Run("get signers", func(t *testing.T) {
		//goland:noinspection GoDeprecation,SpellCheckingInspection
		msgs := []sdk.Msg{
			&MsgRegisterName{
				Owner: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgTransferOwnership{
				Owner: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgSetController{
				Owner: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgUpdateResolveAddress{
				Controller: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgPutAdsSellName{
				Owner: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgCancelAdsSellName{
				Owner: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
			&MsgPurchaseName{
				Buyer: "dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
			},
		}

		for _, msg := range msgs {
			require.Len(t, msg.GetSigners(), 1)
		}
	})

	t.Run("bad signers should panic", func(t *testing.T) {
		msgs := []sdk.Msg{
			&MsgRegisterName{},
			&MsgTransferOwnership{},
			&MsgSetController{},
			&MsgUpdateResolveAddress{},
			&MsgPutAdsSellName{},
			&MsgCancelAdsSellName{},
			&MsgPurchaseName{},
		}

		for _, msg := range msgs {
			require.Panics(t, func() {
				_ = msg.GetSigners()
			})
		}
	})
}

func TestMsgs_ImplementLegacyMsg(t *testing.T) {
	//goland:noinspection GoDeprecation
	msgs := []legacytx.LegacyMsg{
		&MsgRegisterName{},
		&MsgTransferOwnership{},
		&MsgSetController{},
		&MsgUpdateResolveAddress{},
		&MsgPutAdsSellName{},
		&MsgCancelAdsSellName{},
		&MsgPurchaseName{},
	}

	for _, msg := range msgs {
		require.Equal(t, RouterKey, msg.Route())
		require.NotEmpty(t, msg.Type())
		require.NotEmpty(t, msg.GetSignBytes())
	}
}
