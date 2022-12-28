package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Msg Merge Account", func() {

	Describe("ValidateBasic", func() {
		It("should pass validation when all attributes are valid (eth address)", func() {
			ethAddr := "ethm1sdn3kaup7hw5pcmk5sf9j848rzcct7jp7jetft"
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237979b6"

			msg := NewMsgMergeAccount(ethAddr, pubKey, true)
			err := msg.ValidateBasic()
			Expect(err).To(BeNil())
		})

		It("should pass validation when all attributes are valid (cosmos address)", func() {
			cosmosAddr := "ethm12w696c2ef6efgd69xdxt4jeamlc5kgqpdj9tt6"
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237979b6"

			msg := NewMsgMergeAccount(cosmosAddr, pubKey, false)
			err := msg.ValidateBasic()
			Expect(err).To(BeNil())

		})
		It("should fail when isEthAddress dont match the addr provided. isEthAddress is false but provided ethAddr", func() {
			ethAddr := "ethm1sdn3kaup7hw5pcmk5sf9j848rzcct7jp7jetft"
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237979b6"
			msg := NewMsgMergeAccount(ethAddr, pubKey, false)
			err := msg.ValidateBasic()
			Expect(err).ToNot(BeNil())
			Expect(sdkerrors.ErrInvalidRequest.Is(err)).To(BeTrue())
		})

		It("should fail when isEthAddress dont match the addr provided. isEthAddress is true but provided cosmosAddr", func() {
			cosmosAddr := "ethm1sdn3kaup7hw5pcmk5sf9j848rzcct7jp7jetft"
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237979b6"
			msg := NewMsgMergeAccount(cosmosAddr, pubKey, true)
			err := msg.ValidateBasic()
			Expect(err).ToNot(BeNil())
			Expect(sdkerrors.ErrInvalidRequest.Is(err)).To(BeTrue())
		})

		It("should fail when public key provided is invalid", func() {
			cosmosAddr := "ethm12w696c2ef6efgd69xdxt4jeamlc5kgqpdj9tt6"
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb23797z"
			msg := NewMsgMergeAccount(cosmosAddr, pubKey, true)
			err := msg.ValidateBasic()
			Expect(err).ToNot(BeNil())
			Expect(ErrInvalidPubKey.Is(err)).To(BeTrue())
		})
		It("should fail when public key is not provided", func() {
			cosmosAddr := "ethm12w696c2ef6efgd69xdxt4jeamlc5kgqpdj9tt6"
			msg := NewMsgMergeAccount(cosmosAddr, "", true)
			err := msg.ValidateBasic()
			Expect(err).ToNot(BeNil())
			Expect(sdkerrors.ErrInvalidRequest.Is(err)).To(BeTrue())
		})

		It("should fail when msg creator is not provided", func() {
			pubKey := "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237976"
			msg := NewMsgMergeAccount("", pubKey, true)
			err := msg.ValidateBasic()
			Expect(err).ToNot(BeNil())
			Expect(sdkerrors.ErrInvalidAddress.Is(err)).To(BeTrue())
		})

	})
})
