package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type BankKeeper struct {
	bankkeeper.Keeper
}

// TransferEswthDust Placeholder function definition as the actual function is defined in EvmBankKeeper in carbon
// Needed for compile to succeed
func (k BankKeeper) TransferEswthDust(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress) error {
	return nil
}
