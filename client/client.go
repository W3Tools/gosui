package client

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/types"
	"github.com/W3Tools/gosui/utils"
)

// SuiClient is a client for interacting with the Sui blockchain via its RPC API.
type SuiClient struct {
	rpc        string
	requestID  int
	httpClient *http.Client
}

// SuiTransportRequestOptions defines the options for a Sui transport request.
type SuiTransportRequestOptions struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

// NewSuiClient creates a new Sui client with the given RPC URL.
func NewSuiClient(rpc string) (*SuiClient, error) {
	_, err := url.ParseRequestURI(rpc)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    5,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	return &SuiClient{rpc: rpc, requestID: 1, httpClient: httpClient}, nil
}

// RPC returns the RPC URL of the Sui client.
func (client SuiClient) RPC() string {
	return client.rpc
}

// Close the HTTP client connections.
func (client *SuiClient) Close() {
	client.httpClient.CloseIdleConnections()
}

// Call any RPC method
func (client *SuiClient) Call(ctx context.Context, method string, params []any, response any) error {
	return client.request(ctx, SuiTransportRequestOptions{Method: method, Params: params}, &response)
}

// GetRPCAPIVersion returns the version of the Sui RPC API.
func (client *SuiClient) GetRPCAPIVersion(ctx context.Context) (string, error) {
	var response struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	err := client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "rpc.discover",
			Params: []any{},
		},
		&response,
	)
	if err != nil {
		return "", err
	}

	return response.Info.Version, nil
}

// GetCoins fetches the Coin objects owned by an address
func (client *SuiClient) GetCoins(ctx context.Context, input types.GetCoinsParams) (response *types.PaginatedCoins, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	if input.CoinType != nil || *input.CoinType != "" {
		normalized := utils.NormalizeSuiCoinType(*input.CoinType)
		input.CoinType = &normalized
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getCoins",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.CoinType, input.Cursor, input.Limit},
		},
		&response,
	)
}

// GetAllCoins fetches all Coin objects owned by an address
func (client *SuiClient) GetAllCoins(ctx context.Context, input types.GetAllCoinsParams) (response *types.PaginatedCoins, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getAllCoins",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.Cursor, input.Limit},
		},
		&response,
	)
}

// GetBalance fetches the balance of a specific coin type owned by an address.
func (client *SuiClient) GetBalance(ctx context.Context, input types.GetBalanceParams) (response *types.Balance, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	if input.CoinType != nil || *input.CoinType != "" {
		normalized := utils.NormalizeSuiCoinType(*input.CoinType)
		input.CoinType = &normalized
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getBalance",
			Params: []any{utils.NormalizeSuiAddress(input.Owner), input.CoinType},
		},
		&response,
	)
}

// GetAllBalances fetches all balances of all coin types owned by an address.
func (client *SuiClient) GetAllBalances(ctx context.Context, input types.GetAllBalancesParams) (response []*types.Balance, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getAllBalances",
			Params: []any{utils.NormalizeSuiAddress(input.Owner)},
		},
		&response,
	)
}

// GetCoinMetadata fetch metadata for a coin type
func (client *SuiClient) GetCoinMetadata(ctx context.Context, input types.GetCoinMetadataParams) (response *types.CoinMetadata, err error) {
	if input.CoinType != "" {
		input.CoinType = utils.NormalizeSuiCoinType(input.CoinType)
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getCoinMetadata",
			Params: []any{input.CoinType},
		},
		&response,
	)
}

// GetTotalSupply fetch total supply for a coin
func (client *SuiClient) GetTotalSupply(ctx context.Context, input types.GetTotalSupplyParams) (response *types.CoinSupply, err error) {
	if input.CoinType != "" {
		input.CoinType = utils.NormalizeSuiCoinType(input.CoinType)
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getTotalSupply",
			Params: []any{input.CoinType},
		},
		&response,
	)
}

// GetObject returns the object for the given ID
func (client *SuiClient) GetObject(ctx context.Context, input types.GetObjectParams) (response *types.SuiObjectResponse, err error) {
	if input.ID == "" || !utils.IsValidSuiObjectID(utils.NormalizeSuiObjectID(input.ID)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getObject",
			Params: []any{utils.NormalizeSuiObjectID(input.ID), input.Options},
		},
		&response,
	)
}

// MultiGetObjects returns the list of objects for the given IDs
func (client *SuiClient) MultiGetObjects(ctx context.Context, input types.MultiGetObjectsParams) (response []*types.SuiObjectResponse, err error) {
	idmap, ids := make(map[string]struct{}, 0), make([]string, 0)
	for _, id := range input.IDs {
		normalized := utils.NormalizeSuiObjectID(id)
		if id == "" || !utils.IsValidSuiObjectID(normalized) {
			return nil, fmt.Errorf("invalid sui object id %s", id)
		}

		if _, ok := idmap[normalized]; !ok {
			idmap[normalized] = struct{}{}
			ids = append(ids, normalized)
		}
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_multiGetObjects",
			Params: []any{ids, input.Options},
		},
		&response,
	)
}

// GetOwnedObjects returns the list of objects owned by an address
func (client *SuiClient) GetOwnedObjects(ctx context.Context, input types.GetOwnedObjectsParams) (response *types.PaginatedObjectsResponse, err error) {
	if input.Owner == "" || !utils.IsValidSuiAddress(utils.NormalizeSuiAddress(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getOwnedObjects",
			Params: []any{
				utils.NormalizeSuiAddress(input.Owner),
				types.SuiObjectResponseQuery{
					Filter:  input.SuiObjectResponseQuery.Filter,
					Options: input.SuiObjectResponseQuery.Options,
				},
				input.Cursor,
				input.Limit,
			},
		},
		&response,
	)
}

// TryGetPastObject attempts to get an object at a specific version, returning an error if the object does not exist.
func (client *SuiClient) TryGetPastObject(ctx context.Context, input types.TryGetPastObjectParams) (response *types.ObjectReadWrapper, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_tryGetPastObject",
			Params: []any{utils.NormalizeSuiObjectID(input.ID), input.Version, input.Options},
		},
		&response,
	)
}

// GetDynamicFields returns the dynamic fields for a given object ID, paginated.
func (client *SuiClient) GetDynamicFields(ctx context.Context, input types.GetDynamicFieldsParams) (response *types.DynamicFieldPage, err error) {
	if input.ParentID == "" || !utils.IsValidSuiObjectID(utils.NormalizeSuiObjectID(input.ParentID)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFields",
			Params: []any{utils.NormalizeSuiObjectID(input.ParentID), input.Cursor, input.Limit},
		},
		&response,
	)
}

// GetDynamicFieldObject returns the dynamic field object for a given parent ID and field name.
func (client *SuiClient) GetDynamicFieldObject(ctx context.Context, input types.GetDynamicFieldObjectParams) (response *types.SuiObjectResponse, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFieldObject",
			Params: []any{input.ParentID, input.Name},
		},
		&response,
	)
}

// GetTransactionBlock returns the transaction block for a given digest.
func (client *SuiClient) GetTransactionBlock(ctx context.Context, input types.GetTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	if !utils.IsValidTransactionDigest(input.Digest) {
		return nil, fmt.Errorf("invalid transaction digest")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getTransactionBlock",
			Params: []any{input.Digest, input.Options},
		},
		&response,
	)
}

// MultiGetTransactionBlocks returns the transaction blocks for a list of digests.
func (client *SuiClient) MultiGetTransactionBlocks(ctx context.Context, input types.MultiGetTransactionBlocksParams) (response []*types.SuiTransactionBlockResponse, err error) {
	digestmap, digests := make(map[string]struct{}, 0), make([]string, 0)
	for _, digest := range input.Digests {
		if digest == "" || !utils.IsValidTransactionDigest(digest) {
			return nil, fmt.Errorf("invalid transaction digest %s", digest)
		}

		if _, ok := digestmap[digest]; !ok {
			digestmap[digest] = struct{}{}
			digests = append(digests, digest)
		}
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_multiGetTransactionBlocks",
			Params: []any{digests, input.Options},
		},
		&response,
	)
}

// QueryTransactionBlocks returns transaction blocks based on the provided query parameters.
func (client *SuiClient) QueryTransactionBlocks(ctx context.Context, input types.QueryTransactionBlocksParams) (response *types.PaginatedTransactionResponse, err error) {
	var order bool = true
	if input.Order != nil {
		order = (input.Order == &types.Descending)
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_queryTransactionBlocks",
			Params: []any{input.SuiTransactionBlockResponseQuery, input.Cursor, input.Limit, &order},
		},
		&response,
	)
}

// GetTotalTransactionBlocks returns the total number of transaction blocks in the Sui network.
func (client *SuiClient) GetTotalTransactionBlocks(ctx context.Context) (response *big.Int, err error) {
	var data string
	err = client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getTotalTransactionBlocks",
			Params: []any{},
		},
		&data,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(data, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", data)
	}

	return response, nil
}

// SubscribeTransaction unimplement
func (client *SuiClient) SubscribeTransaction(ctx context.Context, input types.SubscribeTransactionParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// GetEvents returns the events for a given event digest.
func (client *SuiClient) GetEvents(ctx context.Context, input types.GetEventsParams) (response []*types.SuiEventBase, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getEvents",
			Params: []any{input.Digest},
		},
		&response,
	)
}

// QueryEvents returns events based on the provided query parameters.
func (client *SuiClient) QueryEvents(ctx context.Context, input types.QueryEventsParams) (response *types.PaginatedEvents, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_queryEvents",
			Params: []any{input.Query, input.Cursor, input.Limit, input.DescendingOrder},
		},
		&response,
	)
}

// SubscribeEvent unimplemented
func (client *SuiClient) SubscribeEvent(ctx context.Context, input types.SubscribeEventParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// GetProtocolConfig returns the protocol configuration for a specific version.
func (client *SuiClient) GetProtocolConfig(ctx context.Context, input types.GetProtocolConfigParams) (response *types.ProtocolConfig, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getProtocolConfig",
			Params: []any{input.Version},
		},
		&response,
	)
}

// GetLatestCheckpointSequenceNumber returns the latest checkpoint sequence number.
func (client *SuiClient) GetLatestCheckpointSequenceNumber(ctx context.Context) (response string, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getLatestCheckpointSequenceNumber",
			Params: []any{},
		},
		&response,
	)
}

// GetCheckpoint returns the checkpoint for a given ID.
func (client *SuiClient) GetCheckpoint(ctx context.Context, input types.GetCheckpointParams) (response *types.Checkpoint, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getCheckpoint",
			Params: []any{input.ID},
		},
		&response,
	)
}

// GetCheckpoints returns a paginated list of checkpoints, starting from the specified cursor.
func (client *SuiClient) GetCheckpoints(ctx context.Context, input types.GetCheckpointsParams) (response *types.CheckpointPage, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getCheckpoints",
			Params: []any{input.Cursor, input.Limit, input.DescendingOrder},
		},
		&response,
	)
}

// GetReferenceGasPrice returns the reference gas price for the Sui network.
func (client *SuiClient) GetReferenceGasPrice(ctx context.Context) (response *big.Int, err error) {
	var data string
	err = client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getReferenceGasPrice",
			Params: []any{},
		},
		&data,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(data, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", data)
	}

	return response, nil
}

// GetLatestSuiSystemState returns the latest Sui system state summary.
func (client *SuiClient) GetLatestSuiSystemState(ctx context.Context) (response *types.SuiSystemStateSummary, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getLatestSuiSystemState",
			Params: []any{},
		},
		&response,
	)
}

// GetCommitteeInfo returns the committee information for a specific epoch.
func (client *SuiClient) GetCommitteeInfo(ctx context.Context, input types.GetCommitteeInfoParams) (response *types.CommitteeInfo, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getCommitteeInfo",
			Params: []any{input.Epoch},
		},
		&response,
	)
}

// GetValidatorsApy returns the list of validators apy
func (client *SuiClient) GetValidatorsApy(ctx context.Context) (response *types.ValidatorsApy, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getValidatorsApy",
			Params: []any{},
		},
		&response,
	)
}

// GetChainIdentifier returns the chain identifier for the Sui network.
func (client *SuiClient) GetChainIdentifier(ctx context.Context) (response string, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getChainIdentifier",
			Params: []any{},
		},
		&response,
	)
}

// GetStakes returns the list of delegated stakes for a given owner address.
func (client *SuiClient) GetStakes(ctx context.Context, input types.GetStakesParams) (response []*types.DelegatedStake, err error) {
	if input.Owner == "" || !utils.IsValidSuiObjectID(utils.NormalizeSuiObjectID(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getStakes",
			Params: []any{utils.NormalizeSuiObjectID(input.Owner)},
		},
		&response,
	)
}

// GetStakesByIds returns the list of delegated stakes for a given list of staked SUI IDs.
func (client *SuiClient) GetStakesByIds(ctx context.Context, input types.GetStakesByIdsParams) (response []*types.DelegatedStake, err error) {
	idmap, ids := make(map[string]struct{}, 0), make([]string, 0)
	for _, id := range input.StakedSuiIds {
		normalized := utils.NormalizeSuiObjectID(id)
		if id == "" || !utils.IsValidSuiObjectID(normalized) {
			return nil, fmt.Errorf("invalid sui object id %s", id)
		}

		if _, ok := idmap[id]; !ok {
			idmap[id] = struct{}{}
			ids = append(ids, id)
		}
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getStakesByIds",
			Params: []any{ids},
		},
		&response,
	)
}

// ResolveNameServiceNames resolves a list of names to their corresponding addresses.
func (client *SuiClient) ResolveNameServiceNames(ctx context.Context, input types.ResolveNameServiceNamesParams) (response *types.ResolvedNameServiceNames, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_resolveNameServiceNames",
			Params: []any{input.Address, input.Cursor, input.Limit},
		},
		&response,
	)
}

// ResolveNameServiceAddress resolves a name to its corresponding address.
func (client *SuiClient) ResolveNameServiceAddress(ctx context.Context, input types.ResolveNameServiceAddressParams) (response string, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_resolveNameServiceAddress",
			Params: []any{input.Name},
		},
		&response,
	)
}

// GetMoveFunctionArgTypes returns the argument types for a Move function.
func (client *SuiClient) GetMoveFunctionArgTypes(ctx context.Context, input types.GetMoveFunctionArgTypesParams) (response []types.SuiMoveFunctionArgTypeWrapper, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getMoveFunctionArgTypes",
			Params: []any{utils.NormalizeSuiObjectID(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// GetNormalizedMoveModulesByPackage returns a structured representation of Move modules in a package.
func (client *SuiClient) GetNormalizedMoveModulesByPackage(ctx context.Context, input types.GetNormalizedMoveModulesByPackageParams) (response *types.SuiMoveNormalizedModules, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModulesByPackage",
			Params: []any{utils.NormalizeSuiObjectID(input.Package)},
		},
		&response,
	)
}

// GetNormalizedMoveModule returns a structured representation of a Move module in a package.
func (client *SuiClient) GetNormalizedMoveModule(ctx context.Context, input types.GetNormalizedMoveModuleParams) (response *types.SuiMoveNormalizedModule, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModule",
			Params: []any{utils.NormalizeSuiObjectID(input.Package), input.Module},
		},
		&response,
	)
}

// GetNormalizedMoveFunction returns a structured representation of a Move function in a module.
func (client *SuiClient) GetNormalizedMoveFunction(ctx context.Context, input types.GetNormalizedMoveFunctionParams) (response *types.SuiMoveNormalizedFunction, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveFunction",
			Params: []any{utils.NormalizeSuiObjectID(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// GetNormalizedMoveStruct returns a structured representation of a Move struct in a module.
func (client *SuiClient) GetNormalizedMoveStruct(ctx context.Context, input types.GetNormalizedMoveStructParams) (response *types.SuiMoveNormalizedStruct, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveStruct",
			Params: []any{utils.NormalizeSuiObjectID(input.Package), input.Module, input.Struct},
		},
		&response,
	)
}

// DryRunTransactionBlock simulates the execution of a transaction block without committing it to the blockchain.
func (client *SuiClient) DryRunTransactionBlock(ctx context.Context, input types.DryRunTransactionBlockParams) (response *types.DryRunTransactionBlockResponse, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_dryRunTransactionBlock",
			Params: []any{b64.ToBase64(input.TransactionBlock)},
		},
		&response,
	)
}

// DevInspectTransactionBlock runs the transaction block in dev-inspect mode.
// Which allows for nearly any transaction (or Move call) with any arguments.
// Detailed results are provided, including both the transaction effects and any return values.
func (client *SuiClient) DevInspectTransactionBlock(ctx context.Context, input types.DevInspectTransactionBlockParams) (response *types.DevInspectResults, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_devInspectTransactionBlock",
			Params: []any{input.Sender, input.TransactionBlock, input.GasPrice, input.Epoch},
		},
		&response,
	)
}

// ExecuteTransactionBlock executes a transaction block on the Sui network.
func (client *SuiClient) ExecuteTransactionBlock(ctx context.Context, input types.ExecuteTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_executeTransactionBlock",
			Params: []any{b64.ToBase64(input.TransactionBlock), input.Signature, input.Options, input.RequestType},
		},
		&response,
	)
}

// SignAndExecuteTransactionBlock signs and executes a transaction block using the provided signer.
func (client *SuiClient) SignAndExecuteTransactionBlock(ctx context.Context, input types.SignAndExecuteTransactionBlockParams) (response *types.SuiTransactionBlockResponse, err error) {
	signatureData, err := input.Signer.SignTransactionBlock(input.TransactionBlock)
	if err != nil {
		return nil, err
	}

	return client.ExecuteTransactionBlock(
		ctx,
		types.ExecuteTransactionBlockParams{
			TransactionBlock: input.TransactionBlock,
			Signature:        []string{signatureData.Signature},
			Options:          input.Options,
			RequestType:      input.RequestType,
		},
	)
}
