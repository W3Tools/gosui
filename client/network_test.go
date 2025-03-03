package client_test

import (
	"testing"

	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/utils"
)

func TestGetFullNodeURL(t *testing.T) {
	tests := []struct {
		network utils.Network
		want    utils.Endpoint
	}{
		{"mainnet", utils.MainnetRPC},
		{"testnet", utils.TestnetRPC},
		{"devnet", utils.DevnetRPC},
		{"localnet", utils.LocalnetRPC},
		{"unknown", utils.DevnetRPC},
	}

	for _, tt := range tests {
		t.Run(tt.network.String(), func(t *testing.T) {
			got := client.GetFullNodeURL(tt.network)
			if got != tt.want.String() {
				t.Errorf("GetFullNodeURL(%q) = %v; want %v", tt.network, got, tt.want)
			}
		})
	}
}
