package client

import "github.com/W3Tools/gosui/utils"

func GetFullNodeURL(network string) string {
	switch network {
	case utils.SuiMainnet:
		return utils.SuiMainnetEndpoint
	case utils.SuiTestnet:
		return utils.SuiTestnetEndpoint
	case utils.SuiDevnet:
		return utils.SuiDevnetEndpoint
	case utils.SuiLocalnet:
		return utils.SuiLocalnet
	default:
		return utils.SuiDevnetEndpoint
	}
}
