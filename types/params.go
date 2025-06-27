package types

import "github.com/W3Tools/gosui/cryptography"

// GetCoinsParams defines the parameters for getting coins owned by a specific address.
type GetCoinsParams struct {
	Owner    string  `json:"owner"`
	CoinType *string `json:"coinType,omitempty"`
	Cursor   *string `json:"cursor,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
}

// GetAllCoinsParams defines the parameters for getting all coins owned by a specific address.
type GetAllCoinsParams struct {
	Owner  string  `json:"owner"`
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

// GetBalanceParams defines the parameters for getting the balance of a specific coin type owned by an address.
type GetBalanceParams struct {
	Owner    string  `json:"owner"`
	CoinType *string `json:"coinType,omitempty"`
}

// GetAllBalancesParams defines the parameters for getting all balances owned by a specific address.
type GetAllBalancesParams struct {
	Owner string `json:"owner"`
}

// GetCoinMetadataParams defines the parameters for getting metadata of a specific coin type.
type GetCoinMetadataParams struct {
	CoinType string `json:"coinType"`
}

// GetTotalSupplyParams defines the parameters for getting the total supply of a specific coin type.
type GetTotalSupplyParams struct {
	CoinType string `json:"coinType"`
}

// GetObjectParams defines the parameters for getting a specific object by its ID.
type GetObjectParams struct {
	ID      string                `json:"id"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

// MultiGetObjectsParams defines the parameters for getting multiple objects by their IDs.
type MultiGetObjectsParams struct {
	IDs     []string              `json:"ids"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

// GetOwnedObjectsParams defines the parameters for getting objects owned by a specific address.
type GetOwnedObjectsParams struct {
	Owner                  string                 `json:"owner"`
	Cursor                 *string                `json:"cursor,omitempty"`
	Limit                  *int                   `json:"limit,omitempty"`
	SuiObjectResponseQuery SuiObjectResponseQuery `json:",inline"`
}

// TryGetPastObjectParams defines the parameters for trying to get a past object by its ID and version.
type TryGetPastObjectParams struct {
	ID      string                `json:"id"`
	Version int                   `json:"version"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

// GetDynamicFieldsParams defines the parameters for getting dynamic fields of a specific object.
type GetDynamicFieldsParams struct {
	ParentID string  `json:"parentId"`
	Cursor   *string `json:"cursor,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
}

// GetDynamicFieldObjectParams defines the parameters for getting a specific dynamic field object by its parent ID and name.
type GetDynamicFieldObjectParams struct {
	ParentID string           `json:"parentId"`
	Name     DynamicFieldName `json:"name"`
}

// GetTransactionBlockParams defines the parameters for getting a specific transaction block by its digest.
type GetTransactionBlockParams struct {
	Digest  string                              `json:"digest"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

// MultiGetTransactionBlocksParams defines the parameters for getting multiple transaction blocks by their digests.
type MultiGetTransactionBlocksParams struct {
	Digests []string                            `json:"digests"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

// QueryTransactionBlocksParams defines the parameters for querying transaction blocks.
type QueryTransactionBlocksParams struct {
	Cursor                           *string                            `json:"cursor,omitempty"`
	Limit                            *int                               `json:"limit,omitempty"`
	Order                            *QueryTransactionBlocksParamsOrder `json:"order,omitempty"`
	SuiTransactionBlockResponseQuery SuiTransactionBlockResponseQuery   `json:",inline"`
}

// GetEventsParams defines the parameters for getting events by their digests.
type GetEventsParams struct {
	Digest string `json:"digests"`
}

// QueryEventsParams defines the parameters for querying events based on a filter.
type QueryEventsParams struct {
	Query           SuiEventFilter `json:"query"`
	Cursor          *EventID       `json:"cursor,omitempty"`
	Limit           *int           `json:"limit,omitempty"`
	DescendingOrder *bool          `json:"order,omitempty"` // default false(ascending order),
}

// GetProtocolConfigParams defines the parameters for getting the protocol configuration.
type GetProtocolConfigParams struct {
	Version *string `json:"version,omitempty"`
}

// GetCheckpointParams defines the parameters for getting a specific checkpoint by its ID.
type GetCheckpointParams struct {
	ID CheckpointID `json:"id"`
}

// GetCheckpointsParams defines the parameters for getting checkpoints with optional pagination and sorting.
type GetCheckpointsParams struct {
	Cursor          *string `json:"cursor,omitempty"`
	Limit           *int    `json:"limit,omitempty"`
	DescendingOrder bool    `json:"descendingOrder"`
}

// GetCommitteeInfoParams defines the parameters for getting committee information for a specific epoch.
type GetCommitteeInfoParams struct {
	Epoch *string `json:"epoch,omitempty"`
}

// SubscribeEventParams defines the parameters for subscribing to Sui events based on a filter.
type SubscribeEventParams struct {
	Filter SuiEventFilter `json:"filter"`
}

// SubscribeTransactionParams defines the parameters for subscribing to transaction events based on a filter.
type SubscribeTransactionParams struct {
	Filter TransactionFilter `json:"filter"`
}

// GetStakesParams defines the parameters for getting stakes owned by a specific address.
type GetStakesParams struct {
	Owner string `json:"owner"`
}

// GetStakesByIdsParams defines the parameters for getting stakes by their IDs.
type GetStakesByIdsParams struct {
	StakedSuiIds []string `json:"stakedSuiIds"`
}

// ResolveNameServiceNamesParams defines the parameters for resolving names in the Sui Name Service to addresses.
type ResolveNameServiceNamesParams struct {
	Address string  `json:"address"`
	Cursor  *string `json:"cursor,omitempty"`
	Limit   *int    `json:"limit,omitempty"`
}

// GetMoveFunctionArgTypesParams defines the parameters for getting argument types of a Move function.
type GetMoveFunctionArgTypesParams struct {
	Package  string `json:"package"`
	Module   string `json:"module"`
	Function string `json:"function"`
}

// GetNormalizedMoveModulesByPackageParams defines the parameters for getting normalized Move modules by package.
type GetNormalizedMoveModulesByPackageParams struct {
	Package string `json:"package"`
}

// GetNormalizedMoveModuleParams defines the parameters for getting a normalized Move module by package and module name.
type GetNormalizedMoveModuleParams struct {
	Package string `json:"package"`
	Module  string `json:"module"`
}

// GetNormalizedMoveFunctionParams defines the parameters for getting a normalized Move function by package, module, and function name.
type GetNormalizedMoveFunctionParams struct {
	Package  string `json:"package"`
	Module   string `json:"module"`
	Function string `json:"function"`
}

// GetNormalizedMoveStructParams defines the parameters for getting a normalized Move struct by package, module, and struct name.
type GetNormalizedMoveStructParams struct {
	Package string `json:"package"`
	Module  string `json:"module"`
	Struct  string `json:"struct"`
}

// ResolveNameServiceAddressParams defines the parameters for resolving a name in the Sui Name Service to an address.
type ResolveNameServiceAddressParams struct {
	Name string `json:"name"`
}

// DryRunTransactionBlockParams defines the parameters for dry-running a transaction block.
type DryRunTransactionBlockParams struct {
	TransactionBlock []byte `json:"transactionBlock"`
}

// DevInspectTransactionBlockParams defines the parameters for inspecting a transaction block in development mode.
type DevInspectTransactionBlockParams struct {
	Sender           string      `json:"sender"`
	TransactionBlock interface{} `json:"transactionBlock"`
	GasPrice         *uint64     `json:"gasPrice,omitempty"`
	Epoch            *string     `json:"epoch,omitempty"`
}

// ExecuteTransactionBlockParams defines the parameters for executing a transaction block.
type ExecuteTransactionBlockParams struct {
	TransactionBlock []byte                              `json:"transactionBlock"`
	Signature        []string                            `json:"signature"`
	Options          *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
	RequestType      *ExecuteTransactionRequestType      `json:"requestType,omitempty"`
}

// SignAndExecuteTransactionBlockParams defines the parameters for signing and executing a transaction block.
type SignAndExecuteTransactionBlockParams struct {
	TransactionBlock []byte                              `json:"transactionBlock"`
	Signer           cryptography.Signer                 `json:"signer"`
	Options          *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
	RequestType      *ExecuteTransactionRequestType      `json:"requestType,omitempty"`
}

// SuiObjectDataFilterMatchAll defines a filter that matches all specified SuiObjectDataFilters.
type SuiObjectDataFilterMatchAll struct {
	MatchAll []SuiObjectDataFilter `json:"MatchAll"`
}

// SuiObjectDataFilterMatchAny defines a filter that matches any of the specified SuiObjectDataFilters.
type SuiObjectDataFilterMatchAny struct {
	MatchAny []SuiObjectDataFilter `json:"MatchAny"`
}

// SuiObjectDataFilterMatchNone defines a filter that matches none of the specified SuiObjectDataFilters.
type SuiObjectDataFilterMatchNone struct {
	MatchNone []SuiObjectDataFilter `json:"MatchNone"`
}

// SuiObjectDataFilterPackage defines a filter for objects belonging to a specific package.
type SuiObjectDataFilterPackage struct {
	Package string `json:"Package"`
}

// SuiObjectDataFilterMoveModule defines a filter for objects belonging to a specific Move module.
type SuiObjectDataFilterMoveModule struct {
	MoveModule SuiObjectDataFilterMoveModuleStruct `json:"MoveModule"`
}

// SuiObjectDataFilterMoveModuleStruct defines the structure for a Move module filter.
type SuiObjectDataFilterMoveModuleStruct struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

// SuiObjectDataFilterStructType defines a filter for objects of a specific struct type.
type SuiObjectDataFilterStructType struct {
	StructType string `json:"StructType"`
}

// SuiObjectDataFilterAddressOwner defines a filter for objects owned by a specific address.
type SuiObjectDataFilterAddressOwner struct {
	AddressOwner string `json:"AddressOwner"`
}

// SuiObjectDataFilterObjectOwner defines a filter for objects owned by a specific object owner.
type SuiObjectDataFilterObjectOwner struct {
	ObjectOwner string `json:"ObjectOwner"`
}

// SuiObjectDataFilterObjectID defines a filter for objects with a specific object ID.
type SuiObjectDataFilterObjectID struct {
	ObjectID string `json:"ObjectId"`
}

// SuiObjectDataFilterObjectIds defines a filter for objects with specific object IDs.
type SuiObjectDataFilterObjectIds struct {
	ObjectIds []string `json:"ObjectIds"`
}

// SuiObjectDataFilterVersion defines a filter for objects with a specific version.
type SuiObjectDataFilterVersion struct {
	Version string `json:"Version"`
}

// QueryTransactionBlocksParamsOrder defines the order in which transaction blocks are queried.
type QueryTransactionBlocksParamsOrder string

var (
	// Ascending is the order for querying transaction blocks in ascending order.
	Ascending QueryTransactionBlocksParamsOrder = "ascending"
	// Descending is the order for querying transaction blocks in descending order.
	Descending QueryTransactionBlocksParamsOrder = "descending"
)

// TransactionFilterMoveFunction defines a filter for transactions that call a specific Move function.
type TransactionFilterMoveFunction struct {
	Function *string `json:"function,omitempty"`
	Module   *string `json:"module,omitempty"`
	Package  string  `json:"package"`
}

// TransactionFilterFromAndToAddress defines a filter for transactions that involve both a "from" and "to" address.
type TransactionFilterFromAndToAddress struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// TransactionFilterFromOrToAddress defines a filter for transactions that involve either a "from" or "to" address.
type TransactionFilterFromOrToAddress struct {
	Addr string `json:"addr"`
}

// SuiEventFilterMoveModule defines a filter for events related to a specific Move module.
type SuiEventFilterMoveModule struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

// SuiEventFilterMoveEventModule defines a filter for events related to a specific Move event module.
type SuiEventFilterMoveEventModule struct {
	Module  string `json:"module"`
	Package string `json:"package"`
}

// SuiEventFilterMoveEventField defines a filter for events related to a specific Move event field.
type SuiEventFilterMoveEventField struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// SuiEventFilterTimeRange defines a filter for events that occurred within a specific time range.
type SuiEventFilterTimeRange struct {
	EndTime   string `json:"endTime"`
	StartTime string `json:"startTime"`
}

// SuiEventFilters defines a slice of SuiEventFilter, allowing for multiple filters to be applied together.
type SuiEventFilters []SuiEventFilter
