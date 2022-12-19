package ante

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

// SignatureVerificationGasConsumer is the type of function that is used to both
// consume gas when verifying signatures and also to accept or reject different types of pubkeys
// This is where apps can define their own PubKey
type SignatureVerificationGasConsumer = func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error

// Eip712SigGasConsumeDecorator Consume parameter-defined amount of gas for each signature according to the passed-in SignatureVerificationGasConsumer function
// before calling the next AnteHandler
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
// Customisation:
// 1.  To convert pubKey to ethsecp256k1 type for eth signature + cosmos address to pass pubkey type test because only ethsecp256k1 pubkey is allowed to pass
type Eip712SigGasConsumeDecorator struct {
	ak             AccountKeeper
	sigGasConsumer SignatureVerificationGasConsumer
}

func NewEip712SigGasConsumeDecorator(ak AccountKeeper, sigGasConsumer SignatureVerificationGasConsumer) Eip712SigGasConsumeDecorator {
	return Eip712SigGasConsumeDecorator{
		ak:             ak,
		sigGasConsumer: sigGasConsumer,
	}
}

func (sgcd Eip712SigGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	params := sgcd.ak.GetParams(ctx)
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	signerAddrs := sigTx.GetSigners()

	for i, sig := range sigs {
		signerAcc, err := authante.GetSignerAcc(ctx, sgcd.ak, signerAddrs[i])
		if err != nil {
			return ctx, err
		}

		pubKey := signerAcc.GetPubKey()

		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.
		if simulate && pubKey == nil {
			pubKey = simSecp256k1Pubkey
		}
		// To allow eth signature + cosmos address to pass pubkey type test because only ethsecp256k1 pubkey is allowed to pass
		if ethsecp256k1.KeyType == sig.PubKey.Type() && bytes.Equal(signerAddrs[i].Bytes(), pubKey.Address().Bytes()) {
			pubKey = &ethsecp256k1.PubKey{Key: pubKey.Bytes()}
		}
		// make a SignatureV2 with PubKey filled in from above
		sig = signing.SignatureV2{
			PubKey:   pubKey,
			Data:     sig.Data,
			Sequence: sig.Sequence,
		}

		err = sgcd.sigGasConsumer(ctx.GasMeter(), sig, params)
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
