package utils

import "fmt"

var (
	SUI_DECIMALS                      = 9
	MIST_PER_SUI                      = 1000000000
	GAS_SAFE_OVERHEAD          uint64 = 1000
	MAX_GAS                    uint64 = 50000000000
	MOVE_STDLIB_ADDRESS               = NormalizeSuiObjectId("0x1") // 0x0000000000000000000000000000000000000000000000000000000000000001
	SUI_FRAMEWORK_ADDRESS             = NormalizeSuiObjectId("0x2") // 0x0000000000000000000000000000000000000000000000000000000000000002
	SUI_SYSTEM_ADDRESS                = NormalizeSuiObjectId("0x3") // 0x0000000000000000000000000000000000000000000000000000000000000003
	SUI_CLOCK_OBJECT_ID               = NormalizeSuiObjectId("0x6") // 0x0000000000000000000000000000000000000000000000000000000000000006
	SUI_SYSTEM_MODULE_NAME            = "sui_system"
	SUI_TYPE_ARG                      = fmt.Sprintf("%s::sui::SUI", SUI_FRAMEWORK_ADDRESS)
	SUI_SYSTEM_STATE_OBJECT_ID        = NormalizeSuiObjectId("0x5")
)

// Endpoint: represents a SUI network endpoint
type Endpoint string

const (
	MainnetRPC  Endpoint = "https://fullnode.mainnet.sui.io"
	TestnetRPC  Endpoint = "https://fullnode.testnet.sui.io"
	DevnetRPC   Endpoint = "https://fullnode.devnet.sui.io"
	LocalnetRPC Endpoint = "http://127.0.0.1:9000"

	MainnetWSS Endpoint = "wss://fullnode.mainnet.sui.io"
	TestnetWSS Endpoint = "wss://fullnode.testnet.sui.io"

	TestnetFaucet  Endpoint = "https://faucet.testnet.sui.io/gas"
	DevnetFaucet   Endpoint = "https://faucet.devnet.sui.io"
	LocalnetFaucet Endpoint = "http://127.0.0.1:9123"
)

func (e Endpoint) String() string {
	return string(e)
}

// Network: represents the SUI network type
type Network string

const (
	Mainnet  Network = "mainnet"
	Testnet  Network = "testnet"
	Devnet   Network = "devnet"
	Localnet Network = "localnet"
)

func (n Network) String() string {
	return string(n)
}
