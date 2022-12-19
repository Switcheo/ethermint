package ante

import (
	"bytes"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

// ValidateEip712MsgDecorator
//  1. Checks if Tx contains MsgMergeAccount.
//  2. Blocks invalid signer and signature combinations
type ValidateEip712MsgDecorator struct {
	evmKeeper EVMKeeper
}

func NewValidateEip712MsgDecorator(evmKeeper EVMKeeper) ValidateEip712MsgDecorator {
	return ValidateEip712MsgDecorator{
		evmKeeper: evmKeeper,
	}
}

func (vmmad ValidateEip712MsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	pubKeys, err := sigTx.GetPubKeys()
	signers := sigTx.GetSigners()
	msgs := tx.GetMsgs()

	msgMergeAccountExists, err := isMergeAccountTx(msgs)
	if err != nil {
		return ctx, err
	}
	for i, signer := range signers {
		if msgMergeAccountExists {
			if cosmosSigner(signer, pubKeys[i]) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "eth signature with cosmos signer is not allowed for merge account msg, only allowed for generic msgs")
			}

		} else {
			if ethSigner(signer, pubKeys[i]) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "eth signature with eth signer is not allowed for generic msgs, only for merge account msg")
			}
			if vmmad.cosmosSignerWithUnmergedAcc(ctx, signer, pubKeys[i]) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "tx not allowed for unmerged eth accounts")
			}

		}

	}
	return next(ctx, tx, simulate)
}

func (vmmad ValidateEip712MsgDecorator) cosmosSignerWithUnmergedAcc(ctx sdk.Context, signer sdk.AccAddress, pubKey cryptotypes.PubKey) bool {
	return cosmosSigner(signer, pubKey) && !vmmad.evmKeeper.AccountHasAlreadyBeenMerged(ctx, signer, false)
}

func cosmosSigner(signer sdk.AccAddress, pubKey cryptotypes.PubKey) bool {
	cosmosPubKey := &secp256k1.PubKey{Key: pubKey.Bytes()}
	return bytes.Equal(cosmosPubKey.Address(), signer)
}

func ethSigner(signer sdk.AccAddress, pubKey cryptotypes.PubKey) bool {
	ethPubKey := &ethsecp256k1.PubKey{Key: pubKey.Bytes()}
	return bytes.Equal(ethPubKey.Address(), signer)
}

func isMergeAccountTx(msgs []sdk.Msg) (bool, error) {
	var msgMergeAccountExists bool
	var msgMergeAccountIndex int
	//for i, msg := range msgs {
	//	if _, isMsgMergeAccount := msg.(*evmtypes.MsgMergeAccount); isMsgMergeAccount {
	//		msgMergeAccountExists = true
	//		msgMergeAccountIndex = i
	//		break
	//	}
	//}

	if msgMergeAccountExists && msgMergeAccountIndex != 0 {
		return true, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "a tx containing merge account msg should have it as the first msg")
	}
	return msgMergeAccountExists, nil

}
