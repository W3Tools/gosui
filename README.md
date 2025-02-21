# gosui

Golang SDK for Sui blockchain, migrated from https://github.com/W3Tools/go-modules/tree/main/gmsui

## Development Example

### Installation

```
go get github.com/W3Tools/gosui
```

### Import it into your program

```
import "github.com/W3Tools/gosui"
```

### Create a client to use JSON-RPC requests

```
package main

import (
	"context"
	"fmt"

	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

func main() {
	// Create a new SuiClient
	suiClient, err := client.NewSuiClient(context.Background(), client.GetFullNodeURL("mainnet"))
	if err != nil {
		panic(err)
	}

	// Get the balance of the owner (0x0)
	balance, err := suiClient.GetBalance(types.GetBalanceParams{Owner: "0x0", CoinType: &utils.SUI_TYPE_ARG})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Balance: %s\n", balance.TotalBalance)

	// Here you can use suiClient to call all RPC methods
}
```

### Create or import mnemonics

```
package main

import (
	"fmt"

	"github.com/W3Tools/gosui/cryptography"
	"github.com/W3Tools/gosui/keypairs/ed25519"
	"github.com/W3Tools/gosui/keypairs/secp256k1"
	"github.com/W3Tools/gosui/keypairs/secp256r1"
)

func main() {
	// Create mnemonic list using the BIP-39 standard
	mnemonics, err := cryptography.GenerateMnemonic()
	if err != nil {
		panic(err)
	}
	fmt.Printf("New mnemonic: %s\n", mnemonics)

	// Or use existing mnemonics
	// mnemonics = "abandon ..."

	// Derive keypair from mnemonics using the Ed25519 algorithm
	ed25519Keypair, err := ed25519.DeriveKeypair(mnemonics, "") // Default derivation path: m/44'/784'/0'/0'/0'
	if err != nil {
		panic(err)
	}
	fmt.Printf("ed25519 address: %s\n", ed25519Keypair.ToSuiAddress())

	// Derive keypair from mnemonics using the Secp256k1 algorithm
	secp256k1Keypair, err := secp256k1.DeriveKeypair(mnemonics, "") // Default derivation path: m/54'/784'/0'/0/0
	if err != nil {
		panic(err)
	}
	fmt.Printf("secp256k1 address: %s\n", secp256k1Keypair.ToSuiAddress())

	// Derive keypair from mnemonics using the Secp256r1 algorithm
	secp256r1Keypair, err := secp256r1.DeriveKeypair(mnemonics, "") // Default derivation path: m/74'/784'/0'/0/0
	if err != nil {
		panic(err)
	}
	fmt.Printf("secp256r1 address: %s\n", secp256r1Keypair.ToSuiAddress())
}
```

### Import or export an address using Sui private key (private key with the suiprivkey prefix)

```
package main

import (
	"fmt"

	"github.com/W3Tools/gosui/cryptography"
	"github.com/W3Tools/gosui/keypairs/ed25519"
	"github.com/W3Tools/gosui/keypairs/secp256k1"
	"github.com/W3Tools/gosui/keypairs/secp256r1"
)

func main() {
	privkey := "suiprivkey1..."

	// Decode the private key
	parsed, err := cryptography.DecodeSuiPrivateKey(privkey)
	if err != nil {
		panic(err)
	}

	// Create a keypair from the private key
	var keypair cryptography.Keypair
	switch parsed.Scheme {
	case cryptography.Ed25519Scheme:
		keypair, err = ed25519.FromSecretKey(parsed.SecretKey, false)
		if err != nil {
			panic(err)
		}
	case cryptography.Secp256k1Scheme:
		keypair, err = secp256k1.FromSecretKey(parsed.SecretKey, false)
		if err != nil {
			panic(err)
		}
	case cryptography.Secp256r1Scheme:
		keypair, err = secp256r1.FromSecretKey(parsed.SecretKey, false)
		if err != nil {
			panic(err)
		}
	}

	// Sign a personal message
	message := []byte("Hello, world!")
	signature, err := keypair.SignPersonalMessage(message)
	if err != nil {
		panic(err)
	}
	fmt.Printf("message bytes: %v, signature: %v\n", signature.Bytes, signature.Signature)

	// Export Sui private key (private key with the suiprivkey prefix)
	expected, err := keypair.GetSecretKey()
	if err != nil {
		panic(err)
	}
	fmt.Printf("matching: %v\n", expected == privkey)
}
```

### Create a Transaction

```
package main

import (
	"context"
	"fmt"

	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/transactions"
	"github.com/W3Tools/gosui/utils"
)

func main() {
	// Create a new Sui client
	suiClient, err := client.NewSuiClient(context.Background(), client.GetFullNodeURL("mainnet"))
	if err != nil {
		panic(err)
	}

	{
		// Example 1. Transfer an object to another address
		tx := transactions.NewTransaction(suiClient)
		if err := tx.TransferObjects([]interface{}{"${OBJECT_ID}"}, "${RECIPIENT_ADDRESS}"); err != nil {
			panic(err)
		}
		// Execute or dry-run the transaction
	}

	{
		// Example 2. Transfer coins to another address
		tx := transactions.NewTransaction(suiClient)

		// 2.1 Merge all coins into one coin object
		coinObjects := []string{"0x1..", "0x2..", "0x3..", "0x4.."} // list of coin objects, the coin type must be the same
		if err := tx.MergeCoins(coinObjects[0], []interface{}{coinObjects[1], coinObjects[2], coinObjects[3]}); err != nil {
			panic(err)
		}

		// 2.2 Split the coin object for transfer, inputCoins contains two coin objects because the two values (1*1e9 and 2*1e9) are cut
		inputCoins, err := tx.SplitCoins(coinObjects[0], []interface{}{uint64(1 * 1e9), uint64(2 * 1e9)})
		if err != nil {
			panic(err)
		}

		// 2.3 Transfer the coin object to another address, RECIPIENT_ADDRESS_ONE gets 1*1e9 and RECIPIENT_ADDRESS_TWO gets 2*1e9
		if err := tx.TransferObjects([]interface{}{inputCoins[0]}, "${RECIPIENT_ADDRESS_ONE}"); err != nil {
			panic(err)
		}
		if err := tx.TransferObjects([]interface{}{inputCoins[1]}, "${RECIPIENT_ADDRESS_TWO}"); err != nil {
			panic(err)
		}
		// Execute or dry-run the transaction
	}

	{
		// Example 3. Calling contract methods
		tx := transactions.NewTransaction(suiClient)

		// 3.1 Call `0x2::balance::zero` to generate an empty SUI balance
		output1, err := tx.MoveCall(
			fmt.Sprintf("%s::balance::zero", utils.SUI_FRAMEWORK_ADDRESS), // 0x2::balance::zero
			[]interface{}{},
			[]string{utils.SUI_TYPE_ARG},
		)
		if err != nil {
			panic(err)
		}
		suiBalance := output1[0]

		// 3.2 Call `0x2::coin::from_balance` to wrap suiBalance into coin
		output2, err := tx.MoveCall(
			fmt.Sprintf("%s::coin::from_balance", utils.SUI_FRAMEWORK_ADDRESS), // 0x2::coin::from_balance
			[]interface{}{suiBalance}, // TxContext will be automatically added
			[]string{utils.SUI_TYPE_ARG},
		)
		if err != nil {
			panic(err)
		}
		suiCoin := output2[0]

		// 3.3 Transfer the coin object to another address
		if err := tx.TransferObjects([]interface{}{suiCoin}, "${RECIPIENT_ADDRESS}"); err != nil {
			panic(err)
		}
		// Execute or dry-run the transaction
	}
}
```

### Execute or dry-run the transaction

```
package main

import (
	"context"
	"fmt"

	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/keypairs/ed25519"
	"github.com/W3Tools/gosui/transactions"
	"github.com/W3Tools/gosui/types"
)

func main() {
	// Create a new Sui client
	suiClient, err := client.NewSuiClient(context.Background(), client.GetFullNodeURL("mainnet"))
	if err != nil {
		panic(err)
	}

	tx := transactions.NewTransaction(suiClient)
	// TODO: implement your transaction here

	{
		// Dry run the transaction
		tx.SetSender("${EXECUTOR_ADDRESS}")
		result, err := tx.DryRunTransactionBlock()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Dry run result: %+v\n", result)
	}

	{
		// Sign and execute the transaction
		keypair, err := ed25519.DeriveKeypair("${EXECUTOR_MNEMONIC}", ed25519.DefaultEd25519DerivationPath)
		if err != nil {
			panic(err)
		}

		tx.SetGasBudget(100000000) // Set the gas budget to 0.1 SUI, will automatically add if not set
		tx.SetGasPrice(750)        // Set the gas price to 0.00000075 SUI, will automatically add if not set

		_, transactionBytes, err := tx.Build(keypair.ToSuiAddress())
		if err != nil {
			panic(err)
		}

		signature, err := keypair.SignTransactionBlock(transactionBytes)
		if err != nil {
			panic(err)
		}

		result, err := suiClient.ExecuteTransactionBlock(types.ExecuteTransactionBlockParams{
			TransactionBlock: transactionBytes,
			Signature:        []string{signature.Signature},
			RequestType:      &types.WaitForLocalExecution,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Execution result: %+v\n", result)

		// You can also use suiClient.SignAndExecuteTransactionBlock
		// result, err := suiClient.SignAndExecuteTransactionBlock(types.SignAndExecuteTransactionBlockParams{
		// 	TransactionBlock: transactionBytes,
		// 	Signer:           keypair,
		// 	Options: &types.SuiTransactionBlockResponseOptions{
		// 		ShowInput: true,
		// 	},
		// 	RequestType: &types.WaitForLocalExecution,
		// })
	}
}
```
