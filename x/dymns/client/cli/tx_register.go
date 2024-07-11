package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/spf13/cobra"
)

const (
	flagYears = "years"
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

			years, _ := cmd.Flags().GetUint16(flagYears)
			if years < 1 {
				return fmt.Errorf("years must be greater than 0, specify by flag --%s", flagYears)
			}

			msg := &dymnstypes.MsgRegisterName{
				Name:     dymName,
				Duration: int32(years),
				Owner:    clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Uint16(flagYears, 0, "number of years to register the Dym-Name for")

	return cmd
}
