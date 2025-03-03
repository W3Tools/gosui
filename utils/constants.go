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

const (
	BvTestnetEndpoint   = "https://sui-testnet-endpoint.blockvision.org"
	BvMainnetEndpoint   = "https://sui-mainnet-endpoint.blockvision.org"
	SuiTestnetEndpoint  = "https://fullnode.testnet.sui.io"
	SuiMainnetEndpoint  = "https://fullnode.mainnet.sui.io"
	SuiDevnetEndpoint   = "https://fullnode.devnet.sui.io"
	SuiLocalnetEndpoint = "http://127.0.0.1:9000"

	WssBvTestnetEndpoint  = "wss://sui-testnet-endpoint.blockvision.org/websocket"
	WssBvMainnetEndpoint  = "wss://sui-mainnet-endpoint.blockvision.org/websocket"
	WssSuiTestnetEndpoint = "wss://fullnode.testnet.sui.io"
	WssSuiMainnetEndpoint = "wss://fullnode.mainnet.sui.io"

	FaucetTestnetEndpoint  = "https://faucet.testnet.sui.io/gas"
	FaucetDevnetEndpoint   = "https://faucet.devnet.sui.io"
	FaucetLocalnetEndpoint = "http://127.0.0.1:9123"
)

const (
	SuiMainnet  = "mainnet"
	SuiTestnet  = "testnet"
	SuiDevnet   = "devnet"
	SuiLocalnet = "localnet"
)
