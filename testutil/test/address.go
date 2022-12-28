package test

import sdk "github.com/cosmos/cosmos-sdk/types"

// Test maker's addresses
var (
	Maker1, Maker2, Maker3, Maker4, Maker5, Maker6 sdk.AccAddress
)

func init() {
	generateAddresses()
}

func generateAddresses() {
	//addresses use for testing
	Maker1, _ = sdk.GetFromBech32("ethm1d9wzmf532988dcehds2h2cw9ds93xuka8rd67k", "ethm")
	Maker2, _ = sdk.GetFromBech32("ethm12e3xkjythwaj8yd0fkm3hkdexhvkngw3fma6rk", "ethm")
	Maker3, _ = sdk.GetFromBech32("ethm18wysn7mm0w7ca0tkgvnh5hrff9ymysayjq38f2", "ethm")
	Maker4, _ = sdk.GetFromBech32("ethm1wx7sfduqz0utwuneg0gwvv4ptmzs7gqrzwtfup", "ethm")
	Maker5, _ = sdk.GetFromBech32("ethm195gwa7l29c0gycm9gk6qug3hqfhf7dsktvj8cy", "ethm")
	Maker6, _ = sdk.GetFromBech32("ethm1x2vk3d8m4lw7065ysryd5z5j4nvy5qc9ga8x2k", "ethm")
}
