package transactions

import (
	"context"
	"fmt"
	"strings"

	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
	"github.com/fardream/go-bcs/bcs"
)

// Transaction defines a programmable transaction builder for Sui.
type Transaction struct {
	client  *client.SuiClient
	builder *sui_types.ProgrammableTransactionBuilder
	ctx     context.Context

	Sender    *sui_types.SuiAddress `json:"sender"`
	GasConfig *GasData              `json:"gasConfig"`
}

// GasData defines gas configuration for a transaction.
type GasData struct {
	Budget  uint64                 `json:"budget"`
	Price   uint64                 `json:"price"`
	Owner   string                 `json:"owner"`
	Payment []*sui_types.ObjectRef `json:"payment"`
}

// TransactionInputGasCoin defines a gas coin input for a transaction.
type TransactionInputGasCoin struct {
	GasCoin bool `json:"gasCoin"`
}

// NewTransaction creates a new Transaction instance with a SuiClient.
func NewTransaction(client *client.SuiClient) *Transaction {
	return &Transaction{
		client:  client,
		builder: sui_types.NewProgrammableTransactionBuilder(),

		GasConfig: new(GasData),
	}
}

// TransactionBuilder returns the underlying ProgrammableTransactionBuilder.
func (txb *Transaction) TransactionBuilder() *sui_types.ProgrammableTransactionBuilder {
	return txb.builder
}

// Client returns the SuiClient associated with the transaction.
func (txb *Transaction) Client() *client.SuiClient {
	return txb.client
}

// Context returns the context associated with the transaction.
func (txb *Transaction) Context() context.Context {
	return txb.ctx
}

// Gas returns a TransactionInputGasCoin indicating a gas coin input.
func (txb *Transaction) Gas() *TransactionInputGasCoin {
	return &TransactionInputGasCoin{GasCoin: true}
}

// SplitCoins encodes a split coins command in the transaction.
func (txb *Transaction) SplitCoins(ctx context.Context, coin interface{}, amounts []interface{}) (returnArguments []*sui_types.Argument, err error) {
	if len(amounts) == 0 {
		return nil, nil
	}

	unresolvedParameter, err := txb.resolveSplitCoinsCoin(coin)
	if err != nil {
		return nil, fmt.Errorf("can not resolve coin in command %d: %v", len(txb.builder.Commands), err)
	}

	amountArguments, err := txb.resolveSplitCoinsAmounts(amounts)
	if err != nil {
		return nil, fmt.Errorf("can not resolve amounts in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(amountArguments)
	arguments, err := unresolvedParameter.resolveAndParseToArguments(ctx, txb.client, txb)
	if err != nil {
		return nil, fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			SplitCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  arguments[0],
				Arguments: arguments[1:],
			},
		},
	)

	return txb.createTransactionResult(len(amounts)), nil
}

// TransferObjects encodes a transfer objects command in the transaction.
func (txb *Transaction) TransferObjects(ctx context.Context, objects []interface{}, address interface{}) error {
	if len(objects) == 0 {
		return nil
	}

	unresolvedParameter, err := txb.resolveTransferObjectsObjects(objects)
	if err != nil {
		return fmt.Errorf("can not resolve objects in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedAddressArgument, err := txb.resolveTransferObjectsAddress(address)
	if err != nil {
		return fmt.Errorf("can not resolve address in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(unresolvedAddressArgument)
	arguments, err := unresolvedParameter.resolveAndParseToArguments(ctx, txb.client, txb)
	if err != nil {
		return fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			TransferObjects: &struct {
				Arguments []sui_types.Argument
				Argument  sui_types.Argument
			}{
				Arguments: arguments[:len(arguments)-1],
				Argument:  arguments[len(arguments)-1],
			},
		},
	)

	return nil
}

// MergeCoins encodes a merge coins command in the transaction.
func (txb *Transaction) MergeCoins(ctx context.Context, destination interface{}, sources []interface{}) (err error) {
	if len(sources) == 0 {
		return nil
	}

	unresolvedParameter, err := txb.resolveMergeCoinsDestination(destination)
	if err != nil {
		return fmt.Errorf("failed to resolve destination in command %d, err: %v", len(txb.builder.Commands), err)
	}

	unresolvedSourceParameter, err := txb.resolveMergeCoinsSources(sources)
	if err != nil {
		return fmt.Errorf("failed to resolve sources in command %d, err: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter.merge(unresolvedSourceParameter)
	arguments, err := unresolvedParameter.resolveAndParseToArguments(ctx, txb.client, txb)
	if err != nil {
		return fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			MergeCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  arguments[0],
				Arguments: arguments[1:],
			},
		},
	)

	return nil
}

// MakeMoveVec encodes a make move vector command in the transaction.
func (txb *Transaction) MakeMoveVec(ctx context.Context, vecType string, arguments []interface{}) ([]*sui_types.Argument, error) {
	typeTag := txb.resolveMakeMoveVecType(vecType)
	unresolvedParameter, err := txb.resolveMakeMoveElement(arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve make move vec element in command %d, err: %v", len(txb.builder.Commands), err)
	}

	inputArguments, err := unresolvedParameter.resolveAndParseToArguments(ctx, txb.client, txb)
	if err != nil {
		return nil, fmt.Errorf("can not resolve and parse to arguments in command %d, err: %v", len(txb.builder.Commands), err)
	}

	txb.builder.Command(
		sui_types.Command{
			MakeMoveVec: &struct {
				TypeTag   *move_types.TypeTag `bcs:"optional"`
				Arguments []sui_types.Argument
			}{
				TypeTag:   typeTag,
				Arguments: inputArguments,
			},
		},
	)

	return txb.createTransactionResult(1), nil
}

// MoveCall encodes a programmable Move call transaction.
func (txb *Transaction) MoveCall(ctx context.Context, target string, arguments []interface{}, typeArguments []string) (returnArguments []*sui_types.Argument, err error) {
	entry := strings.Split(target, "::")
	if len(entry) != 3 {
		return nil, fmt.Errorf("invalid target [%s]", target)
	}
	var pkg, mod, fn = utils.NormalizeSuiObjectID(entry[0]), entry[1], entry[2]

	inputArguments, inputTypeArguments, returnsCount, err := txb.resolveMoveFunction(ctx, pkg, mod, fn, arguments, typeArguments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function arguments, err: %v", err)
	}

	packageID, err := sui_types.NewAddressFromHex(pkg)
	if err != nil {
		return nil, fmt.Errorf("invalid package address [%v]", err)
	}

	txb.builder.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       *packageID,
				Module:        move_types.Identifier(mod),
				Function:      move_types.Identifier(fn),
				Arguments:     inputArguments,
				TypeArguments: inputTypeArguments,
			},
		},
	)

	return txb.createTransactionResult(returnsCount), nil
}

// Build encodes and builds the transaction, returning the transaction data and its BCS-encoded bytes.
func (txb *Transaction) Build(ctx context.Context, sender string) (*sui_types.TransactionData, []byte, error) {
	txb.SetSenderIfNotSet(sender)

	if err := setGasPrice(ctx, txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas price when building transaction: %v", err)
	}
	if err := setGasBudget(ctx, txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas budget when building transaction: %v", err)
	}
	if err := setGasPayment(ctx, txb); err != nil {
		return nil, nil, fmt.Errorf("can not set gas payment when building transaction: %v", err)
	}

	tx := sui_types.NewProgrammable(
		*txb.Sender,
		txb.GasConfig.Payment,
		txb.builder.Finish(),
		txb.GasConfig.Budget,
		txb.GasConfig.Price,
	)
	bs, err := bcs.Marshal(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("can not marshal transaction: %v", err)
	}
	return &tx, bs, err
}

// DryRunTransactionBlock encodes and simulates the transaction block without executing it on-chain.
func (txb *Transaction) DryRunTransactionBlock(ctx context.Context) (*types.DryRunTransactionBlockResponse, error) {
	if txb.Sender == nil {
		return nil, fmt.Errorf("missing transaction sender")
	}

	if err := setGasPrice(ctx, txb); err != nil {
		return nil, fmt.Errorf("failed to set gas price, err: %v", err)
	}

	tx := sui_types.NewProgrammable(
		*txb.Sender,
		nil,
		txb.builder.Finish(),
		utils.MaxGas,
		txb.GasConfig.Price,
	)
	bs, err := bcs.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("can not marshal transaction, err: %v", err)
	}

	return txb.client.DryRunTransactionBlock(ctx, types.DryRunTransactionBlockParams{TransactionBlock: bs})
}

// DevInspectTransactionBlock encodes and simulates the transaction block for developer inspection.
func (txb *Transaction) DevInspectTransactionBlock(ctx context.Context) (*types.DevInspectResults, error) {
	if txb.Sender == nil {
		return nil, fmt.Errorf("missing transaction sender")
	}

	bs, err := bcs.Marshal(txb.builder.Finish())
	if err != nil {
		return nil, fmt.Errorf("can not marshal transaction: %v", err)
	}

	txBytes := append([]byte{0}, bs...)
	return txb.client.DevInspectTransactionBlock(ctx, types.DevInspectTransactionBlockParams{Sender: txb.Sender.String(), TransactionBlock: txBytes})
}

// SetSender sets the sender address for the transaction.
func (txb *Transaction) SetSender(sender string) {
	address, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		panic(fmt.Errorf("failed to create address from hex [%s], err: %v", sender, err))
	}

	txb.Sender = address
}

// SetSenderIfNotSet sets the sender address if it is not already set.
func (txb *Transaction) SetSenderIfNotSet(sender string) {
	if txb.Sender == nil {
		txb.SetSender(sender)
	}
}

// SetGasPrice sets the gas price for the transaction.
func (txb *Transaction) SetGasPrice(price uint64) {
	txb.GasConfig.Price = price
}

// SetGasBudget sets the gas budget for the transaction.
func (txb *Transaction) SetGasBudget(budget uint64) {
	txb.GasConfig.Budget = budget
}

// SetGasBudgetIfNotSet sets the gas budget if it is not already set.
func (txb *Transaction) SetGasBudgetIfNotSet(budget uint64) {
	if txb.GasConfig.Budget == 0 {
		txb.GasConfig.Budget = budget
	}
}

// SetGasOwner sets the gas owner for the transaction.
func (txb *Transaction) SetGasOwner(owner string) {
	txb.GasConfig.Owner = owner
}

// SetGasPayment sets the gas payment objects for the transaction.
func (txb *Transaction) SetGasPayment(payments []*sui_types.ObjectRef) {
	txb.GasConfig.Payment = payments
}
