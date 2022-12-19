package types

import (
	"bytes"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

var _ sdk.Msg = &MsgMergeAccount{}

func NewMsgMergeAccount(creator string, pubKey string, isEthAddress bool) *MsgMergeAccount {
	return &MsgMergeAccount{
		Creator:      creator,
		PubKey:       pubKey,
		IsEthAddress: isEthAddress,
	}
}

// Route ...
func (msg MsgMergeAccount) Route() string {
	return RouterKey
}

// Type ...
func (msg MsgMergeAccount) Type() string {
	return "merge_account"
}

// GetSigners ...
func (msg *MsgMergeAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes ...
func (msg *MsgMergeAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic ...
func (msg *MsgMergeAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "creator: %v", err)
	}
	if msg.PubKey == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "PubKey must be present")

	}
	err = msg.validateCreatorWithPubKey()
	if err != nil {
		return sdkerrors.Wrapf(err, "validateCreatorWithPubKey fail")
	}
	return nil
}

func (msg *MsgMergeAccount) validateCreatorWithPubKey() error {
	var pubKey cryptotypes.PubKey
	pubKeyBz, err := hex.DecodeString(msg.PubKey)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKey, "validation fail: unable to decode public key hex : %s", msg.PubKey)
	}

	if msg.IsEthAddress {
		pubKey = &ethsecp256k1.PubKey{Key: pubKeyBz}
	} else {
		pubKey = &secp256k1.PubKey{Key: pubKeyBz}
	}

	addressBz := pubKey.Address()
	creatorAddressBz, err := msg.getCreatorAddressBytes()
	if err != nil {
		return err
	}
	if !bytes.Equal(addressBz, creatorAddressBz) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "creator does not match expected address generated from pubKey. msg: %+v", msg)

	}
	return nil
}

func (msg *MsgMergeAccount) getCreatorAddressBytes() ([]byte, error) {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	return creator.Bytes(), nil
}
