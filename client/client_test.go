package client_test

import (
	"context"
	"testing"

	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/utils"
)

func TestSuiClientConnectionReuse(t *testing.T) {
	c, err := client.NewSuiClient(client.GetFullNodeURL(utils.Mainnet))
	if err != nil {
		t.Fatalf("Failed to create Sui client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	version, err := c.GetRPCAPIVersion(ctx)
	if err != nil {
		t.Fatalf("Failed to get RPC API version: %v", err)
	}
	t.Logf("RPC API Version: %s", version)

	version1, err := c.GetRPCAPIVersion(ctx)
	if err != nil {
		t.Fatalf("Failed to get RPC API version: %v", err)
	}
	t.Logf("RPC API Version: %s", version1)
}
