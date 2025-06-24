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

type SuiClientOptions struct {
	URL string
}

type SuiClient struct {
	rpc        string
	requestId  int
	httpClient *http.Client
}

type SuiTransportRequestOptions struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

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

	return &SuiClient{rpc: rpc, requestId: 1, httpClient: httpClient}, nil
}

func (client SuiClient) RPC() string {
	return client.rpc
}

func (client *SuiClient) Close() {
	client.httpClient.CloseIdleConnections()
}

// Invoke any RPC method
func (client *SuiClient) Call(ctx context.Context, method string, params []any, response any) error {
	return client.request(ctx, SuiTransportRequestOptions{Method: method, Params: params}, &response)
}

func (client *SuiClient) GetRpcApiVersion(ctx context.Context) (string, error) {
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

// Get all Coin<`coin_type`> objects owned by an address.
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

// Get all Coin objects owned by an address.
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

// Get the total coin balance for one coin type, owned by the address owner.
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

// Get the total coin balance for all coin types, owned by the address owner.
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

// Fetch CoinMetadata for a given coin type
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

// Fetch total supply for a coin
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

// Get details about an object
func (client *SuiClient) GetObject(ctx context.Context, input types.GetObjectParams) (response *types.SuiObjectResponse, err error) {
	if input.ID == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.ID)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getObject",
			Params: []any{utils.NormalizeSuiObjectId(input.ID), input.Options},
		},
		&response,
	)
}

// Batch get details about a list of objects. If any of the object ids are duplicates the call will fail
func (client *SuiClient) MultiGetObjects(ctx context.Context, input types.MultiGetObjectsParams) (response []*types.SuiObjectResponse, err error) {
	idmap, ids := make(map[string]struct{}, 0), make([]string, 0)
	for _, id := range input.IDs {
		normalized := utils.NormalizeSuiObjectId(id)
		if id == "" || !utils.IsValidSuiObjectId(normalized) {
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

// Get all objects owned by an address
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

// Return the object information for a specified version
func (client *SuiClient) TryGetPastObject(ctx context.Context, input types.TryGetPastObjectParams) (response *types.ObjectReadWrapper, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_tryGetPastObject",
			Params: []any{utils.NormalizeSuiObjectId(input.ID), input.Version, input.Options},
		},
		&response,
	)
}

// Return the list of dynamic field objects owned by an object
func (client *SuiClient) GetDynamicFields(ctx context.Context, input types.GetDynamicFieldsParams) (response *types.DynamicFieldPage, err error) {
	if input.ParentId == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.ParentId)) {
		return nil, fmt.Errorf("invalid sui object id")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFields",
			Params: []any{utils.NormalizeSuiObjectId(input.ParentId), input.Cursor, input.Limit},
		},
		&response,
	)
}

// Return the dynamic field object information for a specified object
func (client *SuiClient) GetDynamicFieldObject(ctx context.Context, input types.GetDynamicFieldObjectParams) (response *types.SuiObjectResponse, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getDynamicFieldObject",
			Params: []any{input.ParentId, input.Name},
		},
		&response,
	)
}

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

// Get transaction blocks for a given query criteria
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

// Get total number of transactions
func (client *SuiClient) GetTotalTransactionBlocks(ctx context.Context) (response *big.Int, err error) {
	var response_ string
	err = client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getTotalTransactionBlocks",
			Params: []any{},
		},
		&response_,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(response_, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", response_)
	}

	return response, nil
}

func (client *SuiClient) SubscribeTransaction(ctx context.Context, input types.SubscribeTransactionParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// SuiGetEvents implements the method `sui_getEvents`, gets transaction events.
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

// Get events for a given query criteria
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

// Subscribe to get notifications whenever an event matching the filter occurs
func (client *SuiClient) SubscribeEvent(ctx context.Context, input types.SubscribeEventParams) (response any, err error) {
	return nil, fmt.Errorf("unimplemented")
}

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

// Get the sequence number of the latest checkpoint that has been executed
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

// Returns information about a given checkpoint
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

// Returns historical checkpoints paginated
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

// Getting the reference gas price for the network
func (client *SuiClient) GetReferenceGasPrice(ctx context.Context) (response *big.Int, err error) {
	var response_ string
	err = client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getReferenceGasPrice",
			Params: []any{},
		},
		&response_,
	)
	if err != nil {
		return nil, err
	}

	response, ok := new(big.Int).SetString(response_, 10)
	if !ok {
		return nil, fmt.Errorf("got invalid string number %s", response_)
	}

	return response, nil
}

// Return the latest system state content.
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

// Return the committee information for the asked epoch
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

// Return the Validators APYs
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

// Return the delegated stakes for an address
func (client *SuiClient) GetStakes(ctx context.Context, input types.GetStakesParams) (response []*types.DelegatedStake, err error) {
	if input.Owner == "" || !utils.IsValidSuiObjectId(utils.NormalizeSuiObjectId(input.Owner)) {
		return nil, fmt.Errorf("invalid sui address")
	}

	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "suix_getStakes",
			Params: []any{utils.NormalizeSuiObjectId(input.Owner)},
		},
		&response,
	)
}

// Return the delegated stakes queried by id.
func (client *SuiClient) GetStakesByIds(ctx context.Context, input types.GetStakesByIdsParams) (response []*types.DelegatedStake, err error) {
	idmap, ids := make(map[string]struct{}, 0), make([]string, 0)
	for _, id := range input.StakedSuiIds {
		normalized := utils.NormalizeSuiObjectId(id)
		if id == "" || !utils.IsValidSuiObjectId(normalized) {
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

// Get Move function argument types like read, write and full access
func (client *SuiClient) GetMoveFunctionArgTypes(ctx context.Context, input types.GetMoveFunctionArgTypesParams) (response []types.SuiMoveFunctionArgTypeWrapper, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getMoveFunctionArgTypes",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// Get a map from module name to structured representations of Move modules
func (client *SuiClient) GetNormalizedMoveModulesByPackage(ctx context.Context, input types.GetNormalizedMoveModulesByPackageParams) (response *types.SuiMoveNormalizedModules, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModulesByPackage",
			Params: []any{utils.NormalizeSuiObjectId(input.Package)},
		},
		&response,
	)
}

// Get a structured representation of Move module
func (client *SuiClient) GetNormalizedMoveModule(ctx context.Context, input types.GetNormalizedMoveModuleParams) (response *types.SuiMoveNormalizedModule, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveModule",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module},
		},
		&response,
	)
}

// Get a structured representation of Move function
func (client *SuiClient) GetNormalizedMoveFunction(ctx context.Context, input types.GetNormalizedMoveFunctionParams) (response *types.SuiMoveNormalizedFunction, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveFunction",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Function},
		},
		&response,
	)
}

// Get a structured representation of Move struct
func (client *SuiClient) GetNormalizedMoveStruct(ctx context.Context, input types.GetNormalizedMoveStructParams) (response *types.SuiMoveNormalizedStruct, err error) {
	return response, client.request(
		ctx,
		SuiTransportRequestOptions{
			Method: "sui_getNormalizedMoveStruct",
			Params: []any{utils.NormalizeSuiObjectId(input.Package), input.Module, input.Struct},
		},
		&response,
	)
}

// Dry run a transaction block and return the result.
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

// Runs the transaction block in dev-inspect mode.
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
