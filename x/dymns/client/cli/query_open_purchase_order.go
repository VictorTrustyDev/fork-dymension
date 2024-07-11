package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/spf13/cobra"

	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
)

func CmdQueryOpenPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open-purchase-order [Dym-Name]",
		Aliases: []string{"opo"},
		Short:   "Get current active Open Purchase Order of a Dym-Name.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dymName := args[0]

			if !dymnsutils.IsValidDymName(dymName) {
				return fmt.Errorf("input Dym-Name '%s' is not a valid Dym-Name", dymName)
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := dymnstypes.NewQueryClient(clientCtx)

			res, err := queryClient.OpenPurchaseOrder(cmd.Context(), &dymnstypes.QueryOpenPurchaseOrderRequest{
				DymName: dymName,
			})
			if err != nil {
				return fmt.Errorf("failed to fetch Open Purchase Order of '%s': %w", dymName, err)
			}

			if res == nil {
				return fmt.Errorf("no active Open Purchase Order of '%s'", dymName)
			}

			return clientCtx.PrintProto(&res.Result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
