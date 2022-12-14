package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) Params
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}
