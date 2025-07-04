package transactions

import (
	"context"
	"fmt"

	"github.com/W3Tools/go-sui-sdk/v2/sui_types"
	"github.com/W3Tools/gosui/client"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

// UnresolvedParameter defines unresolved parameters for a transaction command.
type UnresolvedParameter struct {
	Arguments UnresolvedArguments      `json:"argument"`
	Objects   map[int]UnresolvedObject `json:"object"` // map key is index of argument
}

// UnresolvedArgument defines an unresolved argument for a transaction command.
type UnresolvedArgument struct {
	Pure     any
	Object   *sui_types.ObjectArg
	Argument *sui_types.Argument
	resolved bool
}

// UnresolvedArguments defines a slice of unresolved arguments.
type UnresolvedArguments []*UnresolvedArgument

// UnresolvedObject defines an unresolved object for a transaction command.
type UnresolvedObject struct {
	ObjectID string
	Mutable  bool
}

// NewUnresolvedParameter creates and returns a new UnresolvedParameter instance.
func NewUnresolvedParameter(count int) *UnresolvedParameter {
	return &UnresolvedParameter{
		Arguments: make(UnresolvedArguments, count),
		Objects:   make(map[int]UnresolvedObject),
	}
}

// merge merges another UnresolvedParameter into the current one.
func (up *UnresolvedParameter) merge(dest *UnresolvedParameter) {
	if dest == nil {
		return
	}

	count := len(up.Arguments)

	up.Arguments = append(up.Arguments, dest.Arguments...)
	for idx, obj := range dest.Objects {
		up.Objects[idx+count] = obj
	}
}

// resolveAndParseToArguments resolves objects and parses all arguments to Sui types.
func (up *UnresolvedParameter) resolveAndParseToArguments(ctx context.Context, suiClient *client.SuiClient, txb *Transaction) ([]sui_types.Argument, error) {
	err := up.resolveObjects(ctx, suiClient)
	if err != nil {
		return nil, fmt.Errorf("can not resolve objects: %v", err)
	}

	return up.toArguments(txb)
}

// resolveObjects resolves all unresolved objects using the Sui client.
func (up *UnresolvedParameter) resolveObjects(ctx context.Context, suiClient *client.SuiClient) error {
	var ids []string
	for idx, resolve := range up.Objects {
		if entry := cache.GetSharedObject(resolve.ObjectID); entry != nil {
			up.Arguments[idx] = &UnresolvedArgument{Object: entry.ToObjectArg(resolve.Mutable), resolved: true}
		} else {
			ids = append(ids, resolve.ObjectID)
		}
	}

	if len(ids) > 0 {
		var objects []*types.SuiObjectResponse
		objects, err := suiClient.MultiGetObjects(ctx, types.MultiGetObjectsParams{IDs: ids, Options: &types.SuiObjectDataOptions{ShowOwner: true}})
		if err != nil {
			return fmt.Errorf("can not call jsonrpc to multi get objects: %v", err)
		}

		objectMap := utils.SliceToMap(objects, func(v *types.SuiObjectResponse) string {
			if v.Data != nil {
				return v.Data.ObjectID
			}
			return ""
		})

		for idx, resolveObject := range up.Objects {
			if idx >= len(up.Arguments) {
				return fmt.Errorf("can not resolve object at index %d, out of range", idx)
			}

			if argument := up.Arguments[idx]; argument != nil {
				if argument.resolved {
					continue
				}
			}

			object := objectMap[utils.NormalizeSuiObjectID(resolveObject.ObjectID)]
			if object == nil {
				return fmt.Errorf("can not fetch object with id [%s] at index %d", resolveObject.ObjectID, idx)
			}

			objectArg, err := objectResponseToObjectArg(object, resolveObject.Mutable)
			if err != nil {
				return fmt.Errorf("can not convert object response to object arg at index %d: %v", idx, err)
			}

			up.Arguments[idx] = &UnresolvedArgument{Object: objectArg}

			if objectArg.SharedObject != nil {
				cache.AddSharedObject(
					&SharedObjectCacheEntry{
						ObjectID:             &objectArg.SharedObject.Id,
						InitialSharedVersion: &objectArg.SharedObject.InitialSharedVersion,
					},
				)
			}
		}
	}

	return nil
}

// toArguments parses all unresolved arguments to Sui types.
func (up *UnresolvedParameter) toArguments(txb *Transaction) ([]sui_types.Argument, error) {
	arguments := make([]sui_types.Argument, len(up.Arguments))
	for idx, input := range up.Arguments {
		if input.Pure != nil {
			value, err := txb.builder.Pure(input.Pure)
			if err != nil {
				return nil, fmt.Errorf("can not create pure argument at index %d: %v", idx, err)
			}
			arguments[idx] = value
		} else if input.Object != nil {
			value, err := txb.builder.Obj(*input.Object)
			if err != nil {
				return nil, fmt.Errorf("can not create object argument at index %d: %v", idx, err)
			}
			arguments[idx] = value
		} else if input.Argument != nil {
			arguments[idx] = *input.Argument
		} else {
			return nil, fmt.Errorf("invalid input argument at index: %v, pure: %v, object: %v, argument: %v", idx, input.Pure, input.Object, input.Argument)
		}
	}

	return arguments, nil
}
