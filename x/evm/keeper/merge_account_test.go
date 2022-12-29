package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/app"
	"github.com/evmos/ethermint/testutil"
	"github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"strconv"
)

var _ = Describe("Merge Account test", func() {

	var (
		ethermintApp  *app.EthermintApp
		ctx           sdk.Context
		k             keeper.Keeper
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
	)
	var (
		ethAccAddr, _    = sdk.AccAddressFromBech32("ethm1sdn3kaup7hw5pcmk5sf9j848rzcct7jp7jetft")
		cosmosAccAddr, _ = sdk.AccAddressFromBech32("ethm12w696c2ef6efgd69xdxt4jeamlc5kgqpdj9tt6")
		cosmosAddr       = cosmosAccAddr.String()
		ethAddr          = ethAccAddr.String()
		pubKey           = "039d10cebb893e15e2d192d116c9f9c07319fdfd1e7b9d178c50137aeb237979b6"
		evmDenom         = types.DefaultEVMDenom
	)

	BeforeEach(func() {
		ethermintApp = app.Setup(false, nil)
		ctx = ethermintApp.BaseApp.NewContext(false, tmproto.Header{})
		k = *ethermintApp.EvmKeeper
		accountKeeper = ethermintApp.AccountKeeper
		bankKeeper = ethermintApp.BankKeeper
	})

	Describe("Merge account from ETH account Scenario 1", func() {
		iusd := "iusd"
		initialEthAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		Context("when merging ETH account with an already present corresponding cosmos account", func() {
			BeforeEach(func() {
				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialEthAmount)))
				Expect(err).Should(BeNil())

				err = testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(iusd, initialEthAmount)))
				Expect(err).Should(BeNil())

				err = testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should be deleted and balances should be transferred to cosmos account. "+
				"Eth-cosmos address mapping should also be added into accountKeeper"+
				"Nonce in cosmos acc should be updated to the higher one(from ETH)", func() {
				//Artificially increase eth acc nonce for check later
				ethAcc := ethermintApp.AccountKeeper.GetAccount(ctx, ethAccAddr)
				var ethSequence uint64 = 1000
				_ = ethAcc.SetSequence(ethSequence)
				ethermintApp.AccountKeeper.SetAccount(ctx, ethAcc)

				msg := types.NewMsgMergeAccount(ethAddr, pubKey, true)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				originalEthAccountExists := accountKeeper.HasExactAccount(ctx, ethAccAddr)
				finalCosmosAccEvmDenomBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccEvmDenomBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedEvmDenomBalance := initialCosmosAmount.Add(initialEthAmount)
				finalCosmosAccIusdBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, iusd).Amount
				finalEthAccIusdBalance := bankKeeper.GetBalance(ctx, ethAccAddr, iusd).Amount
				mergedIusdBalance := initialEthAmount

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())
				Expect(corrCosmosAcc.GetSequence()).Should(Equal(ethSequence))

				Expect(originalEthAccountExists).Should(BeFalse())
				Expect(finalEthAccEvmDenomBalance).Should(Equal(finalCosmosAccEvmDenomBalance))
				Expect(finalCosmosAccEvmDenomBalance).Should(Equal(mergedEvmDenomBalance))
				//Other balances from eth account should also be moved to cosmos acc
				Expect(mergedIusdBalance).Should(Equal(finalCosmosAccIusdBalance))
				Expect(finalEthAccIusdBalance).Should(Equal(finalCosmosAccIusdBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from ETH account Scenario 2", func() {
		initialEthAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		initialCosmosAmount := sdk.NewIntWithDecimal(0, 8)
		Context("when merging ETH account with an absent cosmos account", func() {
			BeforeEach(func() {
				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialEthAmount)))
				Expect(err).Should(BeNil())

			})
			It("should be deleted and balance should be transferred to a newly created cosmos account. Eth-cosmos address mapping should also be added into accountKeeper", func() {

				msg := types.NewMsgMergeAccount(ethAddr, pubKey, true)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				newCosmosAccCreated := accountKeeper.HasExactAccount(ctx, cosmosAccAddr)
				originalEthAccountExists := accountKeeper.HasExactAccount(ctx, ethAccAddr)
				finalCosmosAccBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedBalance := initialCosmosAmount.Add(initialEthAmount)

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())

				Expect(originalEthAccountExists).Should(BeFalse())
				Expect(newCosmosAccCreated).Should(BeTrue())
				Expect(finalEthAccBalance).Should(Equal(finalCosmosAccBalance))
				Expect(finalCosmosAccBalance).Should(Equal(mergedBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeTrue())

			})
		})

	})

	Describe("Merge account from ETH account Scenario 3", func() {
		initialEthAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)

		Context("when merging ETH account with an already present merged cosmos account(mapping already exists), but original eth account still exists (ideally shouldn't happen)", func() {
			BeforeEach(func() {
				accountKeeper.SetCorrespondingAddresses(ctx, cosmosAccAddr, ethAccAddr)
				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialEthAmount)))
				Expect(err).Should(BeNil())
				err = testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should be deleted and balance should be transferred to cosmos account", func() {

				msg := types.NewMsgMergeAccount(ethAddr, pubKey, true)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				originalEthAccountExists := accountKeeper.HasExactAccount(ctx, ethAccAddr)
				finalCosmosAccBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedBalance := initialCosmosAmount.Add(initialEthAmount)

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())

				Expect(originalEthAccountExists).Should(BeFalse())
				Expect(finalEthAccBalance).Should(Equal(finalCosmosAccBalance))
				Expect(finalCosmosAccBalance).Should(Equal(mergedBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from ETH account Scenario 4", func() {
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)

		Context("when try to execute merge account again even though account already merged (ie address mapping already exists + eth account dont exist)", func() {
			BeforeEach(func() {
				accountKeeper.SetCorrespondingAddresses(ctx, cosmosAccAddr, ethAccAddr)

				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should throw an `account already merged` error", func() {

				msg := types.NewMsgMergeAccount(ethAddr, pubKey, true)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
					}
				}

				Expect(err).ShouldNot(BeNil())
				Expect(types.ErrAccountMerged.Is(err)).To(BeTrue())

				Expect(mergeAccountEventEmitted).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from Cosmos account Scenario 1", func() {
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		Context("when merging cosmos account without any eth account", func() {
			BeforeEach(func() {
				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should create eth-cosmos address mapping in accountKeeper", func() {

				msg := types.NewMsgMergeAccount(cosmosAddr, pubKey, false)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				finalCosmosAccBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedBalance := initialCosmosAmount

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())

				Expect(finalEthAccBalance).Should(Equal(finalCosmosAccBalance))
				Expect(finalCosmosAccBalance).Should(Equal(mergedBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from Cosmos account Scenario 2", func() {
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		initialEthAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		Context("when merging cosmos account with eth account present", func() {
			BeforeEach(func() {
				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialEthAmount)))
				Expect(err).Should(BeNil())

				err = testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should move balances from eth to cosmos account and eth account should be deleted. Eth-cosmos address mapping should also be added into accountKeeper", func() {

				msg := types.NewMsgMergeAccount(cosmosAddr, pubKey, false)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				originalEthAccountExists := accountKeeper.HasExactAccount(ctx, ethAccAddr)
				finalCosmosAccEvmDenomBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccEvmDenomBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedEvmDenomBalance := initialCosmosAmount.Add(initialEthAmount)

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())

				Expect(originalEthAccountExists).Should(BeFalse())
				Expect(finalEthAccEvmDenomBalance).Should(Equal(finalCosmosAccEvmDenomBalance))
				Expect(finalCosmosAccEvmDenomBalance).Should(Equal(mergedEvmDenomBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from Cosmos account Scenario 3", func() {
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		initialEthAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)
		Context("when merging ETH account with an already present merged cosmos account(mapping already exists), but original eth account still exists (ideally shouldn't happen)", func() {

			BeforeEach(func() {
				accountKeeper.SetCorrespondingAddresses(ctx, cosmosAccAddr, ethAccAddr)

				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, ethAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialEthAmount)))
				Expect(err).Should(BeNil())

				err = testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should move balances from eth to cosmos account and eth account should be deleted.", func() {

				msg := types.NewMsgMergeAccount(cosmosAddr, pubKey, false)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				var mergeAccountEvent sdk.Event
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
						mergeAccountEvent = event
					}
				}
				originalEthAccountExists := accountKeeper.HasExactAccount(ctx, ethAccAddr)
				finalCosmosAccEvmDenomBalance := bankKeeper.GetBalance(ctx, cosmosAccAddr, evmDenom).Amount
				finalEthAccEvmDenomBalance := bankKeeper.GetBalance(ctx, ethAccAddr, evmDenom).Amount
				mergedEvmDenomBalance := initialCosmosAmount.Add(initialEthAmount)

				corrCosmosAddress := accountKeeper.GetCorrespondingCosmosAddressIfExists(ctx, ethAccAddr)
				corrCosmosAcc := accountKeeper.GetAccount(ctx, corrCosmosAddress)
				corrEthAddress := accountKeeper.GetCorrespondingEthAddressIfExists(ctx, cosmosAccAddr)
				mergedEthAcc := accountKeeper.GetAccount(ctx, corrEthAddress)

				Expect(err).To(BeNil())

				Expect(originalEthAccountExists).Should(BeFalse())
				Expect(finalEthAccEvmDenomBalance).Should(Equal(finalCosmosAccEvmDenomBalance))
				Expect(finalCosmosAccEvmDenomBalance).Should(Equal(mergedEvmDenomBalance))

				Expect(corrCosmosAddress.String()).Should(Equal(cosmosAddr))
				Expect(corrEthAddress.String()).Should(Equal(ethAddr))

				Expect(corrCosmosAcc).ShouldNot(BeNil())
				Expect(corrCosmosAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(mergedEthAcc).ShouldNot(BeNil())
				Expect(mergedEthAcc.GetAddress().String()).Should(Equal(cosmosAddr))
				Expect(corrCosmosAcc).Should(Equal(mergedEthAcc))

				Expect(mergeAccountEventEmitted).Should(BeTrue())
				Expect(string(mergeAccountEvent.Attributes[2].Key)).Should(Equal(types.AttributeNewCosmosAccCreated))
				Expect(strconv.ParseBool(string(mergeAccountEvent.Attributes[2].Value))).Should(BeFalse())

			})
		})

	})

	Describe("Merge account from Cosmos account Scenario 4", func() {
		initialCosmosAmount := sdk.NewIntWithDecimal(1_000_000_000, 8)

		Context("when trying to execute merge cosmos account again even though account already merged (ie address mapping already exists + eth account dont exist)", func() {
			BeforeEach(func() {
				accountKeeper.SetCorrespondingAddresses(ctx, cosmosAccAddr, ethAccAddr)

				err := testutil.FundAccount(ethermintApp.BankKeeper, ctx, cosmosAccAddr, sdk.NewCoins(sdk.NewCoin(evmDenom, initialCosmosAmount)))
				Expect(err).Should(BeNil())

			})
			It("should throw an `account already merged` error", func() {

				msg := types.NewMsgMergeAccount(cosmosAddr, pubKey, false)
				err := k.MergeUserAccount(ctx, msg)
				events := ctx.EventManager().Events()
				var mergeAccountEventEmitted bool
				for _, event := range events {
					if event.Type == types.EventTypeMergeEthCosmosAccount {
						mergeAccountEventEmitted = true
					}
				}

				Expect(err).ShouldNot(BeNil())
				Expect(types.ErrAccountMerged.Is(err)).To(BeTrue())

				Expect(mergeAccountEventEmitted).Should(BeFalse())

			})
		})

	})

})
