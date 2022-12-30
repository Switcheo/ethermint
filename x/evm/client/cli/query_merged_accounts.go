package cli

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/evmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
)

func GetMergedAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merged-accounts",
		Short: "Query all merged accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MergedAccountsAddressMappings(context.Background(), &types.QueryMergedAccountsMappingRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
