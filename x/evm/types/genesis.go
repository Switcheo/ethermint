package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/evmos/ethermint/types"
)

// Validate performs a basic validation of a GenesisAccount fields.
func (ga GenesisAccount) Validate() error {
	if err := ethermint.ValidateAddress(ga.Address); err != nil {
		return err
	}
	return ga.Storage.Validate()
}

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Accounts:              []GenesisAccount{},
		Params:                DefaultParams(),
		EthToCosmosAddressMap: make(map[string]string),
		CosmosToEthAddressMap: make(map[string]string),
	}
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, accounts []GenesisAccount, ethToCosmosAddressMap, cosmosToEthAddressMap map[string]string) *GenesisState {
	return &GenesisState{
		Accounts:              accounts,
		Params:                params,
		EthToCosmosAddressMap: ethToCosmosAddressMap,
		CosmosToEthAddressMap: cosmosToEthAddressMap,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	seenAccounts := make(map[string]bool)
	for _, acc := range gs.Accounts {
		if seenAccounts[acc.Address] {
			return fmt.Errorf("duplicated genesis account %s", acc.Address)
		}
		if err := acc.Validate(); err != nil {
			return fmt.Errorf("invalid genesis account %s: %w", acc.Address, err)
		}
		seenAccounts[acc.Address] = true
	}
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	return gs.validateAddressMapping()
}

func (gs GenesisState) validateAddressMapping() error {
	if err := validateAddressMap(gs.CosmosToEthAddressMap); err != nil {
		return err
	}
	return validateAddressMap(gs.EthToCosmosAddressMap)
}

func validateAddressMap(addressMap map[string]string) error {
	seenAddressValue := make(map[string]bool)

	for key, value := range addressMap {
		if seenAddressValue[value] {
			return fmt.Errorf("duplicated address value: %s", value)
		}
		if key == value {
			return fmt.Errorf("found same address for key and value in address mapping. key: %s, value: %s", key, value)
		}
		_, err := sdk.AccAddressFromBech32(key)
		if err != nil {
			return fmt.Errorf("unable to convert address key to bech32, key:%s : %w", key, err)

		}
		_, err = sdk.AccAddressFromBech32(value)
		if err != nil {
			return fmt.Errorf("unable to convert address value to bech32, key:%s : %w", value, err)

		}
		seenAddressValue[value] = true

	}

	return nil
}
