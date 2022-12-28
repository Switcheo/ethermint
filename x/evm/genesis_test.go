package evm_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/testutil/test"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/evmos/ethermint/x/evm/types"
)

func (suite *EvmTestSuite) TestInitGenesis() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := common.HexToAddress(privkey.PubKey().Address().String())

	var vmdb *statedb.StateDB

	testCases := []struct {
		name     string
		malleate func()
		genState *types.GenesisState
		expPanic bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesisState(),
			false,
		},
		{
			"valid account",
			func() {
				vmdb.AddBalance(address, big.NewInt(1))
			},
			&types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
						Storage: types.Storage{
							{Key: common.BytesToHash([]byte("key")).String(), Value: common.BytesToHash([]byte("value")).String()},
						},
					},
				},
			},
			false,
		},
		{
			"account not found",
			func() {},
			&types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
		{
			"invalid account type",
			func() {
				acc := authtypes.NewBaseAccountWithAddress(address.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			&types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
		{
			"invalid code hash",
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			&types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
						Code:    "ffffffff",
					},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values
			vmdb = suite.StateDB()

			tc.malleate()
			vmdb.Commit()

			if tc.expPanic {
				suite.Require().Panics(
					func() {
						_ = evm.InitGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper, *tc.genState)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						_ = evm.InitGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper, *tc.genState)
					},
				)
			}
		})
	}
}

func (suite *EvmTestSuite) TestInitGenesisAddressMapping() {

	testCases := []struct {
		name     string
		genState *types.GenesisState
	}{
		{
			"eth-cosmos and cosmos-eth mapping should be populated correctly in the accountKeeper",
			&types.GenesisState{
				Params:   types.DefaultParams(),
				Accounts: []types.GenesisAccount{},
				EthToCosmosAddressMap: map[string]string{
					test.Maker1.String(): test.Maker2.String(),
					test.Maker3.String(): test.Maker4.String(),
					test.Maker5.String(): test.Maker6.String(),
				},
				CosmosToEthAddressMap: map[string]string{
					test.Maker2.String(): test.Maker1.String(),
					test.Maker4.String(): test.Maker3.String(),
					test.Maker6.String(): test.Maker5.String(),
				},
			},
		},
		{
			"eth-cosmos and cosmos-eth mapping should panic if address is not convertible to bech32",
			&types.GenesisState{
				Params:                types.DefaultParams(),
				Accounts:              []types.GenesisAccount{},
				EthToCosmosAddressMap: map[string]string{"tswth1": "tswth2"},
				CosmosToEthAddressMap: map[string]string{"tswth3": "tswth4"},
			},
		},
	}

	suite.Run(testCases[0].name, func() {
		_ = evm.InitGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper, *testCases[0].genState)
		ethCosmosMap := suite.app.AccountKeeper.Store(suite.ctx, authtypes.EthAddressToCosmosAddressKey)
		cosmosEthMap := suite.app.AccountKeeper.Store(suite.ctx, authtypes.CosmosAddressToEthAddressKey)

		itr := ethCosmosMap.Iterator(nil, nil)
		defer itr.Close()
		var ethCosmosMapSize int
		for ; itr.Valid(); itr.Next() {
			ethCosmosMapSize++
		}

		itr2 := cosmosEthMap.Iterator(nil, nil)
		defer itr2.Close()
		var cosmosEthMapSize int
		for ; itr2.Valid(); itr2.Next() {
			cosmosEthMapSize++
		}

		suite.Require().Equal(cosmosEthMapSize, len(testCases[0].genState.CosmosToEthAddressMap))
		suite.Require().Equal(ethCosmosMapSize, len(testCases[0].genState.EthToCosmosAddressMap))

		suite.Require().Equal(sdk.AccAddress(ethCosmosMap.Get(test.Maker1)).String(), testCases[0].genState.EthToCosmosAddressMap[test.Maker1.String()])
		suite.Require().Equal(sdk.AccAddress(ethCosmosMap.Get(test.Maker3)).String(), testCases[0].genState.EthToCosmosAddressMap[test.Maker3.String()])
		suite.Require().Equal(sdk.AccAddress(ethCosmosMap.Get(test.Maker5)).String(), testCases[0].genState.EthToCosmosAddressMap[test.Maker5.String()])

		suite.Require().Equal(sdk.AccAddress(cosmosEthMap.Get(test.Maker2)).String(), testCases[0].genState.CosmosToEthAddressMap[test.Maker2.String()])
		suite.Require().Equal(sdk.AccAddress(cosmosEthMap.Get(test.Maker4)).String(), testCases[0].genState.CosmosToEthAddressMap[test.Maker4.String()])
		suite.Require().Equal(sdk.AccAddress(cosmosEthMap.Get(test.Maker6)).String(), testCases[0].genState.CosmosToEthAddressMap[test.Maker6.String()])

	})

	suite.Run(testCases[1].name, func() {
		suite.Require().Panics(func() {
			_ = evm.InitGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper, *testCases[1].genState)
		})
	})

}

func (suite *EvmTestSuite) TestExportGenesisAddressMapping() {

	testCases := []struct {
		name string
	}{
		{
			"export genesis should export all address mapping correctly",
		},
	}

	suite.Run(testCases[0].name, func() {
		suite.app.AccountKeeper.SetCorrespondingAddresses(suite.ctx, test.Maker1, test.Maker2)
		suite.app.AccountKeeper.SetCorrespondingAddresses(suite.ctx, test.Maker3, test.Maker4)
		genesisState := evm.ExportGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper)

		ethCosmosMap := suite.app.AccountKeeper.Store(suite.ctx, authtypes.EthAddressToCosmosAddressKey)
		cosmosEthMap := suite.app.AccountKeeper.Store(suite.ctx, authtypes.CosmosAddressToEthAddressKey)

		itr := ethCosmosMap.Iterator(nil, nil)
		defer itr.Close()
		var ethCosmosMapSize int
		for ; itr.Valid(); itr.Next() {
			ethCosmosMapSize++
		}

		itr2 := cosmosEthMap.Iterator(nil, nil)
		defer itr2.Close()
		var cosmosEthMapSize int
		for ; itr2.Valid(); itr2.Next() {
			cosmosEthMapSize++
		}

		suite.Require().Equal(cosmosEthMapSize, len(genesisState.CosmosToEthAddressMap))
		suite.Require().Equal(ethCosmosMapSize, len(genesisState.EthToCosmosAddressMap))

		suite.Require().Equal(sdk.AccAddress(ethCosmosMap.Get(test.Maker2)).String(), genesisState.EthToCosmosAddressMap[test.Maker2.String()])
		suite.Require().Equal(sdk.AccAddress(ethCosmosMap.Get(test.Maker4)).String(), genesisState.EthToCosmosAddressMap[test.Maker4.String()])

		suite.Require().Equal(sdk.AccAddress(cosmosEthMap.Get(test.Maker1)).String(), genesisState.CosmosToEthAddressMap[test.Maker1.String()])
		suite.Require().Equal(sdk.AccAddress(cosmosEthMap.Get(test.Maker3)).String(), genesisState.CosmosToEthAddressMap[test.Maker3.String()])

	})

}
