package app

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/evmos/ethermint/encoding"
)

func BenchmarkEthermintApp_ExportAppStateAndValidators(b *testing.B) {
	db := dbm.NewMemDB()
	app := NewEthermintApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		encoding.MakeConfig(ModuleBasics),
		simtestutil.NewAppOptionsWithFlagHome(DefaultNodeHome),
		baseapp.SetChainID(ChainID),
	)

	genesisState := NewTestGenesisState(app.AppCodec())
	stateBytes, err := json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		b.Fatal(err)
	}

	// Initialize the chain
	app.InitChain(
		&abci.RequestInitChain{
			ChainId:       "ethermint_9000-1",
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Making a new app object with the db, so that initchain hasn't been called
		app2 := NewEthermintApp(
			log.NewNopLogger(),
			db,
			nil,
			true,
			map[int64]bool{},
			DefaultNodeHome,
			0,
			encoding.MakeConfig(ModuleBasics),
			simtestutil.NewAppOptionsWithFlagHome(DefaultNodeHome),
			baseapp.SetChainID(ChainID),
		)
		if _, err := app2.ExportAppStateAndValidators(false, []string{}, []string{}); err != nil {
			b.Fatal(err)
		}
	}
}
