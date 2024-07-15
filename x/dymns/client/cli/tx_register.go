package cli

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	flagYears          = "years"
	flagConfirmPayment = "confirm-payment"
)

func NewRegisterDymNameTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [Dym-Name]",
		Short: "Register a new Dym-Name or Extends the duration of an owned Dym-Name.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dymName := args[0]
			if !dymnsutils.IsValidDymName(dymName) {
				return fmt.Errorf("input Dym-Name '%s' is not a valid Dym-Name", dymName)
			}

			years, _ := cmd.Flags().GetInt64(flagYears)
			if years < 1 {
				return fmt.Errorf("years must be greater than 0, specify by flag --%s", flagYears)
			}

			buyer := clientCtx.GetFromAddress().String()

			confirmPaymentStr, _ := cmd.Flags().GetString(flagConfirmPayment)
			if confirmPaymentStr == "" {
				// mode query to get the estimated payment amount
				queryClient := dymnstypes.NewQueryClient(clientCtx)

				resEst, err := queryClient.EstimateRegisterName(cmd.Context(), &dymnstypes.QueryEstimateRegisterNameRequest{
					Name:     dymName,
					Duration: years,
					Owner:    buyer,
				})
				if err != nil {
					return fmt.Errorf("failed to estimate registration/renew fee for '%s': %w", dymName, err)
				}

				toEstimatedAmount := func(amount sdkmath.Int) string {
					return fmt.Sprintf("%s %s", amount.QuoRaw(1e18), strings.ToUpper(params.DisplayDenom))
				}

				fmt.Println("Estimated payment amount:")
				if resEst.FirstYearPrice.IsNil() || resEst.FirstYearPrice.IsZero() {
					fmt.Println("- Registration fee: None")
				} else {
					fmt.Println("- Registration fee + first year fee: ", resEst.FirstYearPrice)
					fmt.Printf("  (~ %s)\n", toEstimatedAmount(resEst.FirstYearPrice.Amount))
				}
				fmt.Print("- Extends duration fee: ")
				if resEst.ExtendPrice.IsNil() || resEst.ExtendPrice.IsZero() {
					fmt.Println("None")
				} else {
					fmt.Println(resEst.ExtendPrice)
					fmt.Printf("  (~ %s)\n", toEstimatedAmount(resEst.ExtendPrice.Amount))
				}
				fmt.Println("- Total fee: ", resEst.TotalPrice)
				fmt.Printf("  (~ %s)\n", toEstimatedAmount(resEst.TotalPrice.Amount))

				fmt.Printf("Supplying flag '--%s=%s to submit the transaction'\n", flagConfirmPayment, resEst.TotalPrice.String())

				return nil
			}

			confirmPayment, err := sdk.ParseCoinNormalized(confirmPaymentStr)
			if err != nil {
				return fmt.Errorf("invalid confirm payment: %v", err)
			}

			return submitRegistration(clientCtx, &dymnstypes.MsgRegisterName{
				Name:           dymName,
				Duration:       years,
				Owner:          buyer,
				ConfirmPayment: confirmPayment,
			}, cmd)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Int64(flagYears, 0, "number of years to register the Dym-Name for")
	cmd.Flags().String(flagConfirmPayment, "", "confirm payment for the Dym-Name registration, without this flag, the command will query the estimated payment amount")

	return cmd
}

func submitRegistration(clientCtx client.Context, msg *dymnstypes.MsgRegisterName, cmd *cobra.Command) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
