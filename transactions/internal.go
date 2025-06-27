package transactions

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/W3Tools/go-sui-sdk/v2/lib"
	"github.com/W3Tools/go-sui-sdk/v2/move_types"
	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
	"github.com/fardream/go-bcs/bcs"
)

// Create Transaction Result With NormalizedMoveFunction Return Count
func (txb *Transaction) createTransactionResult(count int) []*sui_types.Argument {
	nestedResult1 := uint16(len(txb.builder.Commands) - 1)
	returnArguments := make([]*sui_types.Argument, count)
	for i := 0; i < count; i++ {
		returnArguments[i] = &sui_types.Argument{
			NestedResult: &struct {
				Result1 uint16
				Result2 uint16
			}{
				Result1: nestedResult1,
				Result2: uint16(i),
			},
		}
	}

	return returnArguments
}

func setGasPrice(ctx context.Context, txb *Transaction) error {
	if txb.GasConfig.Price == 0 {
		referenceGasPrice, err := txb.client.GetReferenceGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to get reference gas price, err: %v", err)
		}

		txb.GasConfig.Price = referenceGasPrice.Uint64() + 1
	}

	return nil
}

func setGasBudget(ctx context.Context, txb *Transaction) error {
	if txb.GasConfig.Budget == 0 {
		tx := *txb

		dryRunResult, err := tx.DryRunTransactionBlock(ctx)
		if err != nil {
			return fmt.Errorf("failed to dry run transaction block, err: %v", err)
		}

		if dryRunResult.Effects.Status.Status != "success" {
			return fmt.Errorf("dry run failed, could not automatically determine a budget: %v", dryRunResult.Effects.Status.Error)
		}

		safeOverhead := utils.GasSafeOverhead * tx.GasConfig.Price

		computationCost, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.ComputationCost, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse computation cost, err: %v", err)
		}
		baseComputationCostWithOverhead := computationCost + safeOverhead

		storageCost, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.StorageCost, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse storage cost, err: %v", err)
		}
		storageRebate, err := strconv.ParseUint(dryRunResult.Effects.GasUsed.StorageRebate, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse storage rebate, err: %v", err)
		}

		cost := baseComputationCostWithOverhead + storageCost
		if storageRebate > cost {
			txb.GasConfig.Budget = baseComputationCostWithOverhead
		} else {
			gasBudget := baseComputationCostWithOverhead + storageCost - storageRebate
			if gasBudget > baseComputationCostWithOverhead {
				txb.GasConfig.Budget = gasBudget
			} else {
				txb.GasConfig.Budget = baseComputationCostWithOverhead
			}
		}
	}

	return nil
}

func setGasPayment(ctx context.Context, txb *Transaction) error {
	if len(txb.GasConfig.Payment) == 0 {
		owner := txb.GasConfig.Owner
		if owner == "" {
			owner = txb.Sender.String()
		}
		coins, err := txb.client.GetCoins(ctx, types.GetCoinsParams{Owner: owner, CoinType: &utils.SuiTypeArg})
		if err != nil {
			return fmt.Errorf("failed to get coins, err: %v", err)
		}

		paymentCoins := make([]*sui_types.ObjectRef, 0)
		for _, coin := range coins.Data {
			objectRef, err := coinStructToObjectRef(coin)
			if err != nil {
				return fmt.Errorf("failed to create object reference, err: %v", err)
			}

			paymentCoins = append(paymentCoins, objectRef)
		}

		if len(paymentCoins) == 0 {
			return fmt.Errorf("no valid gas coins found for the transaction")
		}

		txb.GasConfig.Payment = paymentCoins
	}

	return nil
}

// Check if the param is tx_context.TxContext
func isTxContext(param types.SuiMoveNormalizedType) bool {
	structType := extractStructTag(param)
	if structType == nil {
		return false
	}

	return structType.Struct.Address == "0x2" && structType.Struct.Module == "tx_context" && structType.Struct.Name == "TxContext"
}

// Extract NormalizedMoveFunction Type
func extractStructTag(normalizedType types.SuiMoveNormalizedType) *types.SuiMoveNormalizedTypeStruct {
	_struct, ok := normalizedType.(types.SuiMoveNormalizedTypeStruct)
	if ok {
		return &_struct
	}

	ref := extractReference(normalizedType)
	mutRef := extractMutableReference(normalizedType)

	if ref != nil {
		return extractStructTag(ref)
	}

	if mutRef != nil {
		return extractStructTag(mutRef)
	}

	return nil
}

func extractReference(normalizedType types.SuiMoveNormalizedType) types.SuiMoveNormalizedType {
	reference, ok := normalizedType.(types.SuiMoveNormalizedTypeReference)
	if ok {
		return reference.Reference.SuiMoveNormalizedType
	}
	return nil
}

func extractMutableReference(normalizedType types.SuiMoveNormalizedType) types.SuiMoveNormalizedType {
	mutableReference, ok := normalizedType.(types.SuiMoveNormalizedTypeMutableReference)
	if ok {
		return mutableReference.MutableReference.SuiMoveNormalizedType
	}
	return nil
}

// Resolve Parameter

// Allowed types are sui_types.Argument, string -> object id, *TransactionInputGasCoin
func (txb *Transaction) resolveSplitCoinsCoin(coin any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(1)

	reflectValue := reflect.ValueOf(coin)
	switch reflectValue.Type() {
	case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
		unresolvedParameter.Arguments[0] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
	case reflect.TypeOf((*TransactionInputGasCoin)(nil)): // gas coin
		unresolvedParameter.Arguments[0] = &UnresolvedArgument{Argument: &sui_types.Argument{GasCoin: &lib.EmptyEnum{}}}
	case reflect.TypeOf(""): // object id
		unresolvedParameter.Objects[0] = UnresolvedObject{Mutable: false, ObjectID: reflectValue.String()}
	default:
		return nil, fmt.Errorf("input coin should one of address(string), sui_types.Argument or *TransactionInputGasCoin, got %v", reflectValue.Type().String())
	}

	return unresolvedParameter, nil
}

// Allowed types are uint, uint8, uint16, uint32, uint64, sui_types.Argument
func (txb *Transaction) resolveSplitCoinsAmounts(amounts []any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(len(amounts))

	for idx, amount := range amounts {
		reflectValue := reflect.ValueOf(amount)
		switch reflectValue.Type() {
		case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
		case reflect.TypeOf(uint(0)), reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)):
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Pure: amount}
		default:
			return nil, fmt.Errorf("input amount should be uint or sui_types.Argument at index %d, got %v", idx, reflectValue.Type().String())
		}
	}

	return unresolvedParameter, nil
}

// Parse TransferObjects Params

// Allowed types are sui_types.Argument, string -> object id
func (txb *Transaction) resolveTransferObjectsObjects(objects []any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(len(objects))

	for idx, object := range objects {
		reflectValue := reflect.ValueOf(object)
		switch reflectValue.Type() {
		case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
		case reflect.TypeOf((*TransactionInputGasCoin)(nil)): // gas coin
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: &sui_types.Argument{GasCoin: &lib.EmptyEnum{}}}
		case reflect.TypeOf(""): // object id
			unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: false, ObjectID: reflectValue.String()}
		default:
			return nil, fmt.Errorf("input object should one of address(string), sui_types.Argument or *TransactionInputGasCoin at index %d, got %v", idx, reflectValue.Type().String())
		}
	}

	return unresolvedParameter, nil
}

// Allowed types are sui_types.Argument, string -> address
func (txb *Transaction) resolveTransferObjectsAddress(address any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(1)

	reflectValue := reflect.ValueOf(address)
	switch reflectValue.Type() {
	case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
		unresolvedParameter.Arguments[0] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
	case reflect.TypeOf(""): // address
		suiAddress, err := sui_types.NewAddressFromHex(reflectValue.String())
		if err != nil {
			return nil, fmt.Errorf("input address must conform to the address(string), got %v", reflectValue.String())
		}

		unresolvedParameter.Arguments[0] = &UnresolvedArgument{Pure: suiAddress}
	default:
		return nil, fmt.Errorf("input address should be address(string) or sui_types.Argument, got %v", reflectValue.Type().String())
	}

	return unresolvedParameter, nil
}

// Parse MergeCoins Params

// Allowed types are sui_types.Argument, string -> object id
func (txb *Transaction) resolveMergeCoinsDestination(destination any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(1)

	reflectValue := reflect.ValueOf(destination)
	switch reflectValue.Type() {
	case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
		unresolvedParameter.Arguments[0] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
	case reflect.TypeOf(""): // address
		unresolvedParameter.Objects[0] = UnresolvedObject{Mutable: false, ObjectID: reflectValue.String()}
	default:
		return nil, fmt.Errorf("input destination should be address(string) or sui_types.Argument, got %v", reflectValue.Type().String())
	}

	return unresolvedParameter, nil
}

// Allowed types are sui_types.Argument, string -> object id
func (txb *Transaction) resolveMergeCoinsSources(sources []any) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(len(sources))

	for idx, source := range sources {
		reflectValue := reflect.ValueOf(source)
		switch reflectValue.Type() {
		case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
		case reflect.TypeOf(""): // object id
			unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: false, ObjectID: reflectValue.String()}
		default:
			return nil, fmt.Errorf("input source should be address(string) or sui_types.Argument at index %d, got %v", idx, reflectValue.Type().String())
		}
	}

	return unresolvedParameter, nil
}

func (txb *Transaction) resolveMakeMoveVecType(vecType string) *move_types.TypeTag {
	switch strings.ToLower(vecType) {
	case "bool":
		return &move_types.TypeTag{Bool: &lib.EmptyEnum{}}
	case "u8":
		return &move_types.TypeTag{U8: &lib.EmptyEnum{}}
	case "u16":
		return &move_types.TypeTag{U16: &lib.EmptyEnum{}}
	case "u32":
		return &move_types.TypeTag{U32: &lib.EmptyEnum{}}
	case "u64":
		return &move_types.TypeTag{U64: &lib.EmptyEnum{}}
	case "u128":
		return &move_types.TypeTag{U128: &lib.EmptyEnum{}}
	case "u256":
		return &move_types.TypeTag{U256: &lib.EmptyEnum{}}
	case "address":
		return &move_types.TypeTag{Address: &lib.EmptyEnum{}}
	case "signer":
		return &move_types.TypeTag{Signer: &lib.EmptyEnum{}}
	default:
		inner := strings.TrimPrefix(vecType, "vector<")
		inner = strings.TrimSuffix(inner, ">")

		if inner != vecType {
			return &move_types.TypeTag{Vector: txb.resolveMakeMoveVecType(inner)}
		}
		return nil
	}
}

func (txb *Transaction) resolveMakeMoveElement(eles []interface{}) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(len(eles))

	for idx, element := range eles {
		reflectValue := reflect.ValueOf(element)
		switch reflectValue.Type() {
		case reflect.TypeOf((*sui_types.Argument)(nil)): // nest result
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: reflectValue.Interface().(*sui_types.Argument)}
		case reflect.TypeOf(uint(0)), reflect.TypeOf(uint8(0)), reflect.TypeOf(uint16(0)), reflect.TypeOf(uint32(0)), reflect.TypeOf(uint64(0)), reflect.TypeOf(&bcs.Uint128{}), reflect.TypeOf(&bcs.Uint256{}):
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Pure: element}
		case reflect.TypeOf(""): // object id
			unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: false, ObjectID: reflectValue.String()}
		default:
			return nil, fmt.Errorf("input amount should be uint or sui_types.Argument at index %d, got %v", idx, reflectValue.Type().String())
		}
	}

	return unresolvedParameter, nil
}

// Resolve Function
func (txb *Transaction) resolveMoveFunction(ctx context.Context, pkg, mod, fn string, arguments []interface{}, typeArguments []string) (inputArguments []sui_types.Argument, inputTypeArguments []move_types.TypeTag, returnsCount int, err error) {
	normalized, err := getNormalizedMoveFunctionFromCache(ctx, txb.client, pkg, mod, fn)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("can not get normalized move function in command %d: %v", len(txb.builder.Commands), err)
	}

	if len(normalized.Parameters) > 0 && isTxContext(normalized.Parameters[len(normalized.Parameters)-1].SuiMoveNormalizedType) {
		normalized.Parameters = normalized.Parameters[:len(arguments)]
	}

	if len(arguments) != len(normalized.Parameters) || len(typeArguments) != len(normalized.TypeParameters) {
		return nil, nil, 0, fmt.Errorf("incorrect number of arguments or type arguments in command %d, required arguments: %d, type arguments: %d", len(txb.builder.Commands), len(normalized.Parameters), len(normalized.TypeParameters))
	}

	inputTypeArguments, err = txb.resolveFunctionTypeArguments(typeArguments)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("can not resolve function type arguments in command %d: %v", len(txb.builder.Commands), err)
	}

	unresolvedParameter, err := txb.resolveFunctionArguments(arguments, normalized.Parameters)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("can not resolve function arguments in command %d: %v", len(txb.builder.Commands), err)
	}

	inputArguments, err = unresolvedParameter.resolveAndParseToArguments(ctx, txb.client, txb)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("can not parse unresolved parameter to arguments in command %d: %v", len(txb.builder.Commands), err)
	}
	return inputArguments, inputTypeArguments, len(normalized.Return), nil
}

func (txb *Transaction) resolveFunctionArguments(inputArguments []interface{}, requiredArguments []*types.SuiMoveNormalizedTypeWrapper) (*UnresolvedParameter, error) {
	unresolvedParameter := NewUnresolvedParameter(len(requiredArguments))
	for idx, parameter := range requiredArguments {
		reflecetInput := reflect.ValueOf(inputArguments[idx])
		if reflecetInput.Type() == reflect.TypeOf((*sui_types.Argument)(nil)) {
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Argument: reflecetInput.Interface().(*sui_types.Argument)}
			continue
		}

		reflectParameter := reflect.ValueOf(parameter.SuiMoveNormalizedType)

		switch reflectParameter.Type() {
		case reflect.TypeOf(types.SuiMoveNormalizedTypeVector{}):
			unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Pure: inputArguments[idx]}
		case reflect.TypeOf(types.SuiMoveNormalizedTypeString("")):
			// Here we are only supporting pure types
			switch reflectParameter.String() {
			case "Bool", "U8", "U16", "U32", "U64", "U128", "U256":
				if reflecetInput.Kind() == reflect.String {
					return nil, fmt.Errorf("input parameter must be bool or unsigned integer at index %d, got %v", idx, reflecetInput.Type())
				}
				unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Pure: inputArguments[idx]}
			case "Address":
				address, err := sui_types.NewAddressFromHex(utils.NormalizeSuiAddress(inputArguments[idx].(string)))
				if err != nil {
					return nil, fmt.Errorf("input parameter must conform to the address(string) at index %d, got %v", idx, inputArguments[idx])
				}
				unresolvedParameter.Arguments[idx] = &UnresolvedArgument{Pure: address}
			default:
				return nil, fmt.Errorf("function string parameter [%v] is not supported at index %d", reflectParameter.String(), idx)
			}
		default:
			if reflecetInput.Type().Kind() != reflect.String {
				return nil, fmt.Errorf("input parameter must be address(string) at index %d, got %v", idx, reflecetInput.Type())
			}
			switch reflectParameter.Type() {
			case reflect.TypeOf(types.SuiMoveNormalizedTypeReference{}):
				unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: false, ObjectID: inputArguments[idx].(string)}
			case reflect.TypeOf(types.SuiMoveNormalizedTypeMutableReference{}):
				unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: true, ObjectID: inputArguments[idx].(string)}
			case reflect.TypeOf(types.SuiMoveNormalizedTypeStruct{}):
				unresolvedParameter.Objects[idx] = UnresolvedObject{Mutable: false, ObjectID: inputArguments[idx].(string)}
			default:
				return nil, fmt.Errorf("function parameter [%v] is not supported at index %d", reflectParameter.Type(), idx)
			}
		}
	}
	return unresolvedParameter, nil
}

func (txb *Transaction) resolveFunctionTypeArguments(typeArguments []string) (inputTypeArguments []move_types.TypeTag, err error) {
	inputTypeArguments = []move_types.TypeTag{}

	for idx, arg := range typeArguments {
		entry := strings.Split(arg, "::")
		if len(entry) != 3 {
			return nil, fmt.Errorf("input type arguments at index %d must be in the format 'address::module::name', got [%v]", idx, arg)
		}

		object, err := sui_types.NewObjectIdFromHex(entry[0])
		if err != nil {
			return nil, fmt.Errorf("input type arguments at index %d must be in the format 'address::module::name', got [%v]: %v", idx, entry[0], err)
		}

		typeTag := move_types.TypeTag{
			Struct: &move_types.StructTag{
				Address: *object,
				Module:  move_types.Identifier(entry[1]),
				Name:    move_types.Identifier(entry[2]),
			},
		}
		inputTypeArguments = append(inputTypeArguments, typeTag)
	}
	return
}

// Convert ObjectResponse to ObjectArg
func objectResponseToObjectArg(data *types.SuiObjectResponse, mutable bool) (*sui_types.ObjectArg, error) {
	if data == nil || data.Data == nil {
		return nil, fmt.Errorf("invalid object response")
	}

	objectArg := new(sui_types.ObjectArg)

	id, err := sui_types.NewObjectIdFromHex(data.Data.ObjectID)
	if err != nil {
		return nil, fmt.Errorf("object id [%s] must be sui address: %v", data.Data.ObjectID, err)
	}

	owner := data.Data.Owner.ObjectOwner
	switch owner := owner.(type) {
	case types.ObjectOwnerShared:
		// Shared object: just set the initial shared version
		objectArg.SharedObject = &struct {
			Id                   move_types.AccountAddress
			InitialSharedVersion uint64
			Mutable              bool
		}{
			Id:                   *id,
			InitialSharedVersion: owner.Shared.InitialSharedVersion,
			Mutable:              mutable,
		}
	default:
		// Other object: set the version and digest
		objectRef, err := ObjectStringRef{ObjectID: data.Data.ObjectID, Version: data.Data.Version, Digest: data.Data.Digest}.ToObjectRef()
		if err != nil {
			return nil, fmt.Errorf("can not convert object ref: %v", err)
		}

		objectArg.ImmOrOwnedObject = objectRef
	}

	return objectArg, nil
}
