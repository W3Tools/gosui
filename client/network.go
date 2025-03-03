package client

import "github.com/W3Tools/gosui/utils"

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
