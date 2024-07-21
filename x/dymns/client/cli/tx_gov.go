package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/dymensionxyz/dymension/v3/app/params"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// NewMigrateChainIdsCmd implements the command to submit an update proposal that update chain-id.
func NewMigrateChainIdsCmd() *cobra.Command {
	cmdCode := "migrate-chain-id"
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s PROPOSAL_FILE", cmdCode),
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal that update chain-id.",
		Long:  `Submit a proposal that update chain-id. The proposal details must be provided via a JSON file.`,
		Example: fmt.Sprintf(`$ %s tx gov submit-legacy-proposal %s proposal_file.json --from=<key_or_address>

Sample proposal file content:
// all fields are required
{
  "replacement": [{
      "previous_chain_id": "cosmoshub-3",
      "new_chain_id": "cosmoshub-4"
  },{
      "previous_chain_id": "columbus-4",
      "new_chain_id": "columbus-5"
  }]
}`,
			version.AppName,
			cmdCode,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			proposal, err := parseMigrateChainIdsProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := dymnstypes.NewMigrateChainIdsProposal(title, description, proposal.Replacement...)

			msg, err := govv1beta1.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagTitle); err != nil {
		panic(err)
	}

	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagDescription); err != nil {
		panic(err)
	}

	cmd.Flags().String(cli.FlagDeposit, "1000"+params.BaseDenom, "deposit of proposal")
	if err := cmd.MarkFlagRequired(cli.FlagDeposit); err != nil {
		panic(err)
	}

	return cmd
}

// parseMigrateChainIdsProposal reads and parses proposal from a JSON file.
func parseMigrateChainIdsProposal(cdc codec.JSONCodec, metadataFile string) (*dymnstypes.MigrateChainIdsProposal, error) {
	proposal := dymnstypes.MigrateChainIdsProposal{}

	contents, err := os.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return nil, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proposal: %w", err)
	}

	return &proposal, nil
}
