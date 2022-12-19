package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
	"strconv"
)

// CmdMergeAccount command build merge account transaction
func CmdMergeAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-account [public-key] [is-eth-address]",
		Short: "Merge Cosmos and Eth account transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsPublicKey := args[0]
			argsIsEthAddress, err := strconv.ParseBool(args[1])
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid isEthAddress: %s", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := evmtypes.NewMsgMergeAccount(
				clientCtx.GetFromAddress().String(),
				argsPublicKey,
				argsIsEthAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
