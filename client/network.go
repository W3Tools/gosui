package client

import "github.com/W3Tools/gosui/utils"

// GetFullNodeURL returns the full node URL based on the provided network type.
func GetFullNodeURL(network utils.Network) string {
	switch network {
	case utils.Mainnet:
		return utils.MainnetRPC.String()
	case utils.Testnet:
		return utils.TestnetRPC.String()
	case utils.Devnet:
		return utils.DevnetRPC.String()
	case utils.Localnet:
		return utils.LocalnetRPC.String()
	default:
		return utils.DevnetRPC.String()
	}
}
