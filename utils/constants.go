package utils

import "fmt"

var (
	// SuiDecimals is the number of decimals used in SUI tokens
	SuiDecimals = 9
	// MistPerSui is the number of mist per SUI token
	MistPerSui = 1000000000
	// GasSafeOverhead is the overhead for gas in SUI transactions
	GasSafeOverhead uint64 = 1000
	// MaxGas is the maximum gas limit for SUI transactions
	MaxGas uint64 = 50000000000
	// MoveStdlibAddress is the address of the Move standard library in SUI, 0x0000000000000000000000000000000000000000000000000000000000000001
	MoveStdlibAddress = NormalizeSuiObjectID("0x1")
	// SuiFrameworkAddress is the address of the SUI framework, 0x0000000000000000000000000000000000000000000000000000000000000002
	SuiFrameworkAddress = NormalizeSuiObjectID("0x2")
	// SuiSystemAddress is the address of the SUI system, 0x0000000000000000000000000000000000000000000000000000000000000003
	SuiSystemAddress = NormalizeSuiObjectID("0x3")
	// SuiClockObjectID is the object ID for SUI clock, 0x0000000000000000000000000000000000000000000000000000000000000006
	SuiClockObjectID = NormalizeSuiObjectID("0x6")
	// SuiSystemModuleName is the name of the SUI system module
	SuiSystemModuleName = "sui_system"
	// SuiTypeArg is the type argument for SUI, used in Move types, 0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI
	SuiTypeArg = fmt.Sprintf("%s::sui::SUI", SuiFrameworkAddress)
	// SuiSystemStateObjectID is the object ID for the SUI system state, 0x0000000000000000000000000000000000000000000000000000000000000005
	SuiSystemStateObjectID = NormalizeSuiObjectID("0x5")
)

// Endpoint represents a SUI network endpoint
type Endpoint string

const (
	// MainnetRPC is the RPC endpoint for the SUI mainnet
	MainnetRPC Endpoint = "https://fullnode.mainnet.sui.io"
	// TestnetRPC is the RPC endpoint for the SUI testnet
	TestnetRPC Endpoint = "https://fullnode.testnet.sui.io"
	// DevnetRPC is the RPC endpoint for the SUI devnet
	DevnetRPC Endpoint = "https://fullnode.devnet.sui.io"
	// LocalnetRPC is the RPC endpoint for the SUI localnet
	LocalnetRPC Endpoint = "http://127.0.0.1:9000"

	// MainnetWSS is the WebSocket endpoint for the SUI mainnet
	MainnetWSS Endpoint = "wss://fullnode.mainnet.sui.io"
	// TestnetWSS is the WebSocket endpoint for the SUI testnet
	TestnetWSS Endpoint = "wss://fullnode.testnet.sui.io"

	// TestnetFaucet is the faucet endpoint for the SUI testnet
	TestnetFaucet Endpoint = "https://faucet.testnet.sui.io/gas"
	// DevnetFaucet is the faucet endpoint for the SUI devnet
	DevnetFaucet Endpoint = "https://faucet.devnet.sui.io"
	// LocalnetFaucet is the faucet endpoint for the SUI localnet
	LocalnetFaucet Endpoint = "http://127.0.0.1:9123"
)

func (e Endpoint) String() string {
	return string(e)
}

// Network represents the SUI network type
type Network string

const (
	// Mainnet is the main SUI network
	Mainnet Network = "mainnet"
	// Testnet is the test SUI network
	Testnet Network = "testnet"
	// Devnet is the development SUI network
	Devnet Network = "devnet"
	// Localnet is the local SUI network
	Localnet Network = "localnet"
)

func (n Network) String() string {
	return string(n)
}
