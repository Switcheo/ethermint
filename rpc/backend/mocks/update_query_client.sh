# ensure you have mockery installed first (brew install mockery)
cd ~/go/src/github.com/Switcheo/ethermint/x/evm/types
mockery --name=QueryClient --filename=query_client.go --output=../../../rpc/backend/mocks
