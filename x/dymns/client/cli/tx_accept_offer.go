package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/spf13/cobra"
)

func NewAcceptOfferToBuyDymNameTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("accept-offer [offer-id] [amount] %s", params.DisplayDenom),
		Short: "Accept an Offer-To-Buy your Dym-Name",
		Example: fmt.Sprintf(
			"$ %s tx %s accept-offer 1 50 %s --%s hub-user",
			version.AppName, dymnstypes.ModuleName,
			params.DisplayDenom,
			flags.FlagFrom,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			offerId := args[0]
			if !dymnsutils.IsValidBuyNameOfferId(offerId) {
				return fmt.Errorf("input Offer-ID '%s' is not a valid Offer-ID", offerId)
			}

			amount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || amount < 1 {
				return fmt.Errorf("amount must be a positive number")
			}

			denom := args[2]
			if !strings.EqualFold(denom, params.DisplayDenom) {
				return fmt.Errorf("denom must be %s", strings.ToUpper(params.DisplayDenom))
			}

			owner := clientCtx.GetFromAddress().String()
			if owner == "" {
				return fmt.Errorf("flag --%s is required", flags.FlagFrom)
			}

			queryClient := dymnstypes.NewQueryClient(clientCtx)

			resParams, err := queryClient.Params(context.Background(), &dymnstypes.QueryParamsRequest{})
			if err != nil {
				return err
			}

			msg := &dymnstypes.MsgAcceptOfferBuyName{
				OfferId: offerId,
				Owner:   owner,
				MinAccept: sdk.Coin{
					Denom:  resParams.Params.Price.PriceDenom,
					Amount: sdk.NewInt(int64(amount)).MulRaw(1e18),
				},
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
