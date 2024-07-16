package cli

import (
	"context"
	"cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	dymnskeeper "github.com/dymensionxyz/dymension/v3/x/dymns/keeper"
	dymnstypes "github.com/dymensionxyz/dymension/v3/x/dymns/types"
	dymnsutils "github.com/dymensionxyz/dymension/v3/x/dymns/utils"
	"github.com/spf13/cobra"
)

func NewUpdateResolveDymNameAddressTxCmd() *cobra.Command {
	//goland:noinspection SpellCheckingInspection
	cmd := &cobra.Command{
		Use:     "resolve [Dym-Name address] [?resolve to]",
		Short:   "Configure resolve Dym-Name address. 2nd arg if empty means to remove the configuration.",
		Example: "resolve bonded-pool.staking.dym dym1fl48vsnmsdzcv85q5d2q4z5ajdha8yu38x9fue",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var dymNameAddress, resolveTo string
			dymNameAddress = args[0]
			if len(args) > 1 {
				resolveTo = args[1]
			}

			queryClient := dymnstypes.NewQueryClient(clientCtx)

			resModuleParams, err := queryClient.Params(context.Background(), &dymnstypes.QueryParamsRequest{})
			if err != nil {
				return errors.Wrap(err, "failed to query module params")
			}
			moduleParams := resModuleParams.Params

			subName, dymName, chainIdOrAlias, err := dymnskeeper.ParseDymNameAddress(dymNameAddress)
			if err != nil {
				return errors.Wrap(err, "failed to parse input Dym-Name-Address")
			}

			chainId := func(chainIdOrAlias string) string {
				// translate to chain-id of is an alias
				for chainId, aliasesOfChainId := range moduleParams.Chains.AliasesByChainId {
					if chainId == chainIdOrAlias {
						return chainId
					}
					for _, alias := range aliasesOfChainId.Aliases {
						if alias == chainIdOrAlias {
							chainIdOrAlias = alias
							return chainId
						}
					}
				}
				return chainIdOrAlias
			}(chainIdOrAlias)

			if !dymnsutils.IsValidChainIdFormat(chainId) {
				return fmt.Errorf("input chain-id '%s' is not a valid chain-id", chainId)
			}

			msg := &dymnstypes.MsgUpdateResolveAddress{
				Name:       dymName,
				ChainId:    chainId,
				SubName:    subName,
				ResolveTo:  resolveTo,
				Controller: clientCtx.GetFromAddress().String(),
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
