package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func (k *Keeper) MergeUserAccount(ctx sdk.Context, msg *types.MsgMergeAccount) (*types.MsgMergeAccountResponse, error) {
	panic("WIP")

	//addr, err := sdk.AccAddressFromBech32(msg.Creator)
	//if err != nil {
	//	return nil, sdkerrors.Wrapf(types.ErrInvalidCreatorAddress, "address %s convert from bech 32 to accAddress fail", msg.Creator)
	//}
	//
	//if k.AccountHasAlreadyBeenMerged(ctx, addr, msg.IsEthAddress) {
	//	return nil, sdkerrors.Wrapf(types.ErrAccountMerged, "merging not required. Mapping already exists and no eth account found to merge. Msg: %s", msg)
	//}
	//pubKeyBz, err := hex.DecodeString(msg.PubKey)
	//if err != nil {
	//	return nil, sdkerrors.Wrapf(types.ErrInvalidPubKey, "unable to decode public key hex : %s", msg.PubKey)
	//}
	//
	//var cosmosAcc authtypes.AccountI
	//var ethAcc authtypes.AccountI
	//var newCosmosAccCreated bool
	//
	//if msg.IsEthAddress {
	//	ethAcc = k.accountKeeper.GetAccount(ctx, addr)
	//	cosmosPubkey := &secp256k1.PubKey{Key: pubKeyBz}
	//	cosmosAddr := sdk.AccAddress(cosmosPubkey.Address())
	//	cosmosAcc = k.accountKeeper.GetAccount(ctx, cosmosAddr)
	//	if cosmosAcc == nil {
	//		cosmosAcc = k.addNewCosmosAccount(ctx, cosmosAddr)
	//		newCosmosAccCreated = true
	//	}
	//} else {
	//	ethPubkey := &ethsecp256k1.PubKey{Key: pubKeyBz}
	//	ethAddr := sdk.AccAddress(ethPubkey.Address())
	//	ethAcc = k.accountKeeper.GetAccount(ctx, ethAddr)
	//	cosmosAcc = k.accountKeeper.GetAccount(ctx, addr)
	//	if ethAcc == nil {
	//		err = k.moveEthBankBalanceToCosmosAddress(ctx, ethAddr, addr)
	//		if err != nil {
	//			return nil, err
	//		}
	//		k.accountKeeper.SetCorrespondingAddresses(ctx, cosmosAcc.GetAddress(), ethAddr)
	//		emitMergeAccountEvent(ctx, ethAddr, cosmosAcc.GetAddress(), newCosmosAccCreated)
	//		return &types.MsgMergeAccountResponse{
	//			MergedCosmosAddress: cosmosAcc.GetAddress().String(),
	//		}, nil
	//	}
	//}
	//
	//err = k.mergeEthAndCosmosAccounts(ctx, ethAcc, cosmosAcc)
	//if err != nil {
	//	return nil, err
	//}
	//
	//k.accountKeeper.SetCorrespondingAddresses(ctx, cosmosAcc.GetAddress(), ethAcc.GetAddress())
	//
	//emitMergeAccountEvent(ctx, ethAcc.GetAddress(), cosmosAcc.GetAddress(), newCosmosAccCreated)
	//
	//return &types.MsgMergeAccountResponse{
	//	MergedCosmosAddress: cosmosAcc.GetAddress().String(),
	//}, nil

}

func (k *Keeper) addNewCosmosAccount(ctx sdk.Context, cosmosAddress sdk.AccAddress) authtypes.AccountI {
	newAccount := k.accountKeeper.NewAccountWithAddress(ctx, cosmosAddress)

	//To standardise all accounts created by evm module to have an empty string code hash for external accounts.
	if acct, ok := newAccount.(ethermint.EthAccountI); ok {
		emptyCodeHash := common.BytesToHash(crypto.Keccak256(nil))
		acct.SetCodeHash(emptyCodeHash)
	}
	k.accountKeeper.SetAccount(ctx, newAccount)
	return newAccount
}

func (k *Keeper) mergeEthAndCosmosAccounts(ctx sdk.Context, ethAcc authtypes.AccountI, cosmosAcc authtypes.AccountI) error {
	if err := k.moveEthBankBalanceToCosmosAddress(ctx, ethAcc.GetAddress(), cosmosAcc.GetAddress()); err != nil {
		return err
	}
	k.setLargerNonce(ctx, ethAcc, cosmosAcc)
	k.accountKeeper.RemoveAccount(ctx, ethAcc)
	return nil
}

// Compares the nonce between eth and cosmos acc and set the larger number to the cosmos acc.
// Prevents future replay attacks
func (k *Keeper) setLargerNonce(ctx sdk.Context, ethAcc authtypes.AccountI, cosmosAcc authtypes.AccountI) {
	ethNonce := ethAcc.GetSequence()
	cosmosNonce := cosmosAcc.GetSequence()
	if ethNonce > cosmosNonce {
		_ = cosmosAcc.SetSequence(ethNonce)
		k.accountKeeper.SetAccount(ctx, cosmosAcc)
	}
}

func (k *Keeper) moveEthBankBalanceToCosmosAddress(ctx sdk.Context, ethAddress sdk.AccAddress, cosmosAddress sdk.AccAddress) error {
	panic("WIP")
	//ethCoins := k.bankKeeper.GetAllBalances(ctx, ethAddress)
	//if len(ethCoins) > 0 {
	//	for _, ethCoin := range ethCoins {
	//		err := k.bankKeeper.SendCoins(ctx, ethAddress, cosmosAddress, sdk.Coins{sdk.Coin{Denom: ethCoin.GetDenom(), Amount: ethCoin.Amount}})
	//		if err != nil {
	//			return sdkerrors.Wrapf(err, "move balance from eth account: %s to cosmos account: %s for denom %s failed", ethAddress, cosmosAddress, ethCoin.GetDenom())
	//		}
	//	}
	//}
	//return nil
}

func (k *Keeper) AccountHasAlreadyBeenMerged(ctx sdk.Context, address sdk.AccAddress, isEthAddress bool) bool {
	if isEthAddress {
		return k.accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, address) != nil && !k.accountKeeper.HasExactAccount(ctx, address)
	}
	ethAddress := k.accountKeeper.GetCorrespondingEthAddressIfExists(ctx, address)
	return ethAddress != nil && !k.accountKeeper.HasExactAccount(ctx, ethAddress)

}

func emitMergeAccountEvent(ctx sdk.Context, ethAddress sdk.AccAddress, cosmosAddress sdk.AccAddress, newComsosAccCreated bool) {

	ctx.EventManager().EmitTypedEvents(&types.MergeAccountEvent{
		CosmosAddress:       cosmosAddress.String(),
		EthAddress:          ethAddress.String(),
		NewCosmosAccCreated: newComsosAccCreated,
	})

}
