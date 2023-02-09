echo "You may need to install mockery by running \"brew install mockery\""
cd ~/go/src/github.com/Switcheo/ethermint/x/evm/types
mockery --name=QueryClient --filename=query_client.go --output=../../../rpc/backend/mocks
