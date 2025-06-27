package types

// Balance defines the structure for a balance in SUI.
type Balance struct {
	CoinObjectCount int               `json:"coinObjectCount"`
	CoinType        string            `json:"coinType"`
	LockedBalance   map[string]string `json:"lockedBalance"`
	TotalBalance    string            `json:"totalBalance"`
}

// BalanceChange defines the structure for a balance change in SUI.
type BalanceChange struct {
	Owner    ObjectOwnerWrapper `json:"owner"`
	CoinType string             `json:"coinType"`
	Amount   string             `json:"amount"`
}

// Checkpoint defines the structure for a checkpoint in SUI.
type Checkpoint struct {
	Epoch                      string                 `json:"epoch"`
	SequenceNumber             string                 `json:"sequenceNumber"`
	Digest                     string                 `json:"digest"`
	NetworkTotalTransactions   string                 `json:"networkTotalTransactions"`
	PreviousDigest             string                 `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary GasCostSummary         `json:"epochRollingGasCostSummary"`
	TimestampMs                string                 `json:"timestampMs"`
	EndOfEpochData             *EndOfEpochData        `json:"endOfEpochData,omitempty"`
	Transactions               []string               `json:"transactions"`
	CheckpointCommitments      []CheckpointCommitment `json:"checkpointCommitments"`
	ValidatorSignature         string                 `json:"validatorSignature"`
}

// GasCostSummary defines the structure for gas cost summary in SUI.
type GasCostSummary struct {
	ComputationCost         string `json:"computationCost"`
	StorageCost             string `json:"storageCost"`
	StorageRebate           string `json:"storageRebate"`
	NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
}

// EndOfEpochData defines the structure for end of epoch data in SUI.
type EndOfEpochData struct {
	EpochCommitments         []CheckpointCommitment `json:"epochCommitments"`
	NextEpochCommittee       [][2]string            `json:"nextEpochCommittee"`
	NextEpochProtocolVersion string                 `json:"nextEpochProtocolVersion"`
}

// CheckpointCommitment defines the structure for a checkpoint commitment in SUI.
type CheckpointCommitment struct {
	ECMHLiveObjectSetDigest ECMHLiveObjectSetDigest `json:"ecmhLiveObjectSetDigest"`
}

// ECMHLiveObjectSetDigest defines the structure for ECMH live object set digest in SUI.
type ECMHLiveObjectSetDigest struct {
	Digest []int `json:"digest"` // TODO: bytes?
}

// CheckpointID defines the type for checkpoint IDs in SUI.
type CheckpointID string

// Claim defines the structure for a claim in SUI.
type Claim struct {
	IndexMod4 uint64 `json:"indexMod4"`
	Value     string `json:"value"`
}

// CoinStruct defines the structure for a coin object in SUI.
type CoinStruct struct {
	Balance             string `json:"balance"`
	CoinObjectID        string `json:"coinObjectId"`
	CoinType            string `json:"coinType"`
	Digest              string `json:"digest"`
	PreviousTransaction string `json:"previousTransaction"`
	Version             string `json:"version"`
}

// CommitteeInfo defines the structure for committee information in SUI.
type CommitteeInfo struct {
	Epoch      string      `json:"epoch"`
	Validators [][2]string `json:"validators"`
}

// DelegatedStake defines the structure for delegated stake in SUI.
type DelegatedStake struct {
	ValidatorAddress string               `json:"validatorAddress"`
	StakingPool      string               `json:"stakingPool"`
	Stakes           []StakeObjectWrapper `json:"stakes"`
}

// DevInspectResults defines the structure for the results of a dev inspect operation in SUI.
type DevInspectResults struct {
	Effects TransactionEffects   `json:"effects"`
	Error   string               `json:"error,omitempty"`
	Events  []SuiEvent           `json:"events"`
	Results []SuiExecutionResult `json:"results,omitempty"`
}

// TransactionEffects defines the structure for transaction effects in SUI.
type TransactionEffects struct {
	MessageVersion       string                                      `json:"messageVersion"`
	Status               ExecutionStatus                             `json:"status"`
	ExecutedEpoch        string                                      `json:"executedEpoch"`
	GasUsed              GasCostSummary                              `json:"gasUsed"`
	ModifiedAtVersions   []TransactionBlockEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	SharedObjects        []SuiObjectRef                              `json:"sharedObjects,omitempty"`
	TransactionDigest    string                                      `json:"transactionDigest"`
	Created              []OwnedObjectRef                            `json:"created,omitempty"`
	Mutated              []OwnedObjectRef                            `json:"mutated,omitempty"`
	Deleted              []SuiObjectRef                              `json:"deleted,omitempty"`
	GasObject            OwnedObjectRef                              `json:"gasObject"`
	EventsDigest         *string                                     `json:"eventsDigest,omitempty"`
	Dependencies         []string                                    `json:"dependencies,omitempty"`
	Unwrapped            []OwnedObjectRef                            `json:"unwrapped,omitempty"`
	UnwrappedThenDeleted []SuiObjectRef                              `json:"unwrappedThenDeleted,omitempty"`
	Wrapped              []SuiObjectRef                              `json:"wrapped,omitempty"`
}

// ExecutionStatus defines the structure for the execution status of a transaction in SUI.
type ExecutionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// TransactionBlockEffectsModifiedAtVersions defines the structure for modified at versions in transaction block effects.
type TransactionBlockEffectsModifiedAtVersions struct {
	ObjectID       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

// SuiObjectRef defines the structure for a reference to an object in SUI.
type SuiObjectRef struct {
	ObjectID string `json:"objectId"`
	Version  uint64 `json:"version"`
	Digest   string `json:"digest"`
}

// OwnedObjectRef defines the structure for an owned object reference in SUI.
type OwnedObjectRef struct {
	Owner     ObjectOwnerWrapper `json:"owner"`
	Reference SuiObjectRef       `json:"reference"`
}

// SuiEventBase defines the base structure for an event in SUI.
type SuiEventBase struct {
	ID                EventID     `json:"id"`
	PackageID         string      `json:"packageId"`
	TransactionModule string      `json:"transactionModule"`
	Sender            string      `json:"sender"`
	Type              string      `json:"type"`
	ParsedJSON        interface{} `json:"parsedJson"`
	Bcs               string      `json:"bcs"`
}

// SuiEvent defines the structure for an event in SUI, including a timestamp.
type SuiEvent struct {
	SuiEventBase
	TimestampMs string `json:"timestampMs,omitempty"`
}

// EventID defines the structure for an event ID in SUI.
type EventID struct {
	TxDigest string `json:"txDigest"`
	EventSeq string `json:"eventSeq"`
}

// SuiExecutionResult defines the structure for the result of executing a transaction in SUI.
type SuiExecutionResult struct {
	MutableReferenceOutputs [][3]interface{} `json:"mutableReferenceOutputs,omitempty"` // TODO: [SuiArgument, bytes, string][]
	ReturnValues            [][2]interface{} `json:"returnValues,omitempty"`            // TODO: interface -> [bytes, string][]
}

// DisplayFieldsResponse defines the structure for the response of display fields in SUI.
type DisplayFieldsResponse struct {
	Data  map[string]*string          `json:"data"`
	Error *ObjectResponseErrorWrapper `json:"error"`
}

// DryRunTransactionBlockResponse defines the structure for the response of a dry run transaction block in SUI.
type DryRunTransactionBlockResponse struct {
	Effects        TransactionEffects       `json:"effects"`
	Events         []SuiEvent               `json:"events"`
	ObjectChanges  []SuiObjectChangeWrapper `json:"objectChanges"`
	BalanceChanges []BalanceChange          `json:"balanceChanges"`
	Input          TransactionBlockData     `json:"input"`
}

// TransactionBlockData defines the structure for transaction block data in SUI.
type TransactionBlockData struct {
	MessageVersion string                         `json:"messageVersion"`
	Transaction    SuiTransactionBlockKindWrapper `json:"transaction"`
	Sender         string                         `json:"sender"`
	GasData        SuiGasData                     `json:"gasData"`
}

// SuiGasData defines the structure for gas data in SUI transactions.
type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	Owner   string         `json:"owner"`
	Price   string         `json:"price"`
	Budget  string         `json:"budget"`
}

// DynamicFieldInfo defines the structure for dynamic field information in SUI.
type DynamicFieldInfo struct {
	Name       DynamicFieldName `json:"name"`
	BcsName    string           `json:"bcsName"`
	Type       DynamicFieldType `json:"type"`
	ObjectType string           `json:"objectType"`
	ObjectID   string           `json:"objectId"`
	Version    int64            `json:"version"`
	Digest     string           `json:"digest"`
}

// DynamicFieldName defines the structure for a dynamic field name in SUI.
type DynamicFieldName struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// DynamicFieldType defines the type of a dynamic field in SUI.
type DynamicFieldType string

var (
	// DynamicField is a type for dynamic fields in SUI.
	DynamicField DynamicFieldType = "DynamicField"
	// DynamicObject is a type for dynamic objects in SUI.
	DynamicObject DynamicFieldType = "DynamicObject"
)

// ExecuteTransactionRequestType defines the type of transaction execution request in SUI.
type ExecuteTransactionRequestType string

var (
	// WaitForEffectsCert indicates that the transaction should wait for an effects certificate.
	WaitForEffectsCert ExecuteTransactionRequestType = "WaitForEffectsCert"
	// WaitForLocalExecution is the default type for executing a transaction in SUI.
	WaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

// SuiObjectData defines the structure for object data in SUI.
type SuiObjectData struct {
	ObjectID            string                 `json:"objectId"`
	Version             string                 `json:"version"`
	Digest              string                 `json:"digest"`
	Type                *string                `json:"type,omitempty"`
	Owner               *ObjectOwnerWrapper    `json:"owner,omitempty"`
	PreviousTransaction *string                `json:"previousTransaction,omitempty"`
	StorageRebate       *string                `json:"storageRebate,omitempty"`
	Display             *DisplayFieldsResponse `json:"display,omitempty"`
	Content             *SuiParsedDataWrapper  `json:"content,omitempty"`
	Bcs                 *RawDataWrapper        `json:"bcs,omitempty"`
}

// SuiObjectDataOptions defines the options for displaying SUI object data.
type SuiObjectDataOptions struct {
	ShowBcs                 bool `json:"showBcs,omitempty"`
	ShowContent             bool `json:"showContent,omitempty"`
	ShowDisplay             bool `json:"showDisplay,omitempty"`
	ShowOwner               bool `json:"showOwner,omitempty"`
	ShowPreviousTransaction bool `json:"showPreviousTransaction,omitempty"`
	ShowStorageRebate       bool `json:"showStorageRebate,omitempty"`
	ShowType                bool `json:"showType,omitempty"`
}

// SuiObjectResponseQuery defines the structure for querying SUI object responses.
type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

// PaginatedCoins defines a paginated response for coins in SUI.
type PaginatedCoins struct {
	Data        []CoinStruct `json:"data"`
	HasNextPage bool         `json:"hasNextPage"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
}

// PaginatedDynamicFieldInfos defines a paginated response for dynamic field infos in SUI.
type PaginatedDynamicFieldInfos struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  *string            `json:"nextCursor,omitempty"`
	HasNextPage bool               `json:"hasNextPage"`
}

// PaginatedEvents defines a paginated response for events in SUI.
type PaginatedEvents struct {
	Data        []SuiEvent `json:"data"`
	NextCursor  *EventID   `json:"nextCursor,omitempty"`
	HasNextPage bool       `json:"hasNextPage"`
}

// PaginatedStrings defines a paginated response for strings in SUI.
type PaginatedStrings struct {
	Data        []string `json:"data"`
	NextCursor  *string  `json:"nextCursor,omitempty"`
	HasNextPage bool     `json:"hasNextPage"`
}

// PaginatedObjectsResponse defines a paginated response for SuiObjectResponse in SUI.
type PaginatedObjectsResponse struct {
	Data        []SuiObjectResponse `json:"data"`
	NextCursor  *string             `json:"nextCursor,omitempty"`
	HasNextPage bool                `json:"hasNextPage"`
}

// SuiObjectResponse defines the response for a SUI object.
type SuiObjectResponse struct {
	Data  *SuiObjectData              `json:"data,omitempty"`
	Error *ObjectResponseErrorWrapper `json:"error,omitempty"`
}

// PaginatedTransactionResponse defines a paginated response for transaction blocks in SUI.
type PaginatedTransactionResponse struct {
	Data        []SuiTransactionBlockResponse `json:"data"`
	NextCursor  *string                       `json:"nextCursor,omitempty"`
	HasNextPage bool                          `json:"hasNextPage"`
}

// SuiTransactionBlockResponse defines the response for a transaction block in SUI.
type SuiTransactionBlockResponse struct {
	Digest                  string                    `json:"digest"`
	Transaction             *SuiTransactionBlock      `json:"transaction,omitempty"`
	RawTransaction          string                    `json:"rawTransaction,omitempty"`
	Effects                 *TransactionEffects       `json:"effects,omitempty"`
	Events                  []*SuiEvent               `json:"events,omitempty"`
	ObjectChanges           []*SuiObjectChangeWrapper `json:"objectChanges,omitempty"`
	BalanceChanges          []*BalanceChange          `json:"balanceChanges,omitempty"`
	TimestampMs             *string                   `json:"timestampMs,omitempty"`
	Checkpoint              *string                   `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                     `json:"confirmedLocalExecution,omitempty"`
	Errors                  []string                  `json:"errors,omitempty"`
}

// SuiTransactionBlock defines a transaction block in SUI.
type SuiTransactionBlock struct {
	Data         TransactionBlockData `json:"data"`
	TxSignatures []string             `json:"txSignatures"`
}

// ProtocolConfig defines the protocol configuration in SUI.
type ProtocolConfig struct {
	MinSupportedProtocolVersion string                         `json:"minSupportedProtocolVersion"`
	MaxSupportedProtocolVersion string                         `json:"maxSupportedProtocolVersion"`
	ProtocolVersion             string                         `json:"protocolVersion"`
	FeatureFlags                map[string]bool                `json:"featureFlags"`
	Attributes                  map[string]ProtocolConfigValue `json:"attributes"`
}

// SuiActiveJwk defines an active JWK in SUI.
type SuiActiveJwk struct {
	Epoch string   `json:"epoch"`
	Jwk   SuiJWK   `json:"jwk"`
	JwkID SuiJwkID `json:"jwk_id"`
}

// SuiJWK defines a JWK in SUI.
type SuiJWK struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kty string `json:"kty"`
	N   string `json:"n"`
}

// SuiJwkID defines a JWK ID in SUI.
type SuiJwkID struct {
	Iss string `json:"iss"`
	Kid string `json:"kid"`
}

// SuiAuthenticatorStateExpire defines the expiration state of an authenticator in SUI.
type SuiAuthenticatorStateExpire struct {
	MinEpoch string `json:"min_epoch"`
}

// SuiChangeEpoch defines the change of an epoch in SUI.
type SuiChangeEpoch struct {
	ComputationCharge     string `json:"computation_charge"`
	Epoch                 string `json:"epoch"`
	EpochStartTimestampMs string `json:"epoch_start_timestamp_ms"`
	StorageCharge         string `json:"storage_charge"`
	StorageRebate         string `json:"storage_rebate"`
}

// CoinMetadata defines the metadata of a coin in SUI.
type CoinMetadata struct {
	Decimals    uint8  `json:"decimals"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
}

// SuiMoveAbilitySet defines a set of Move abilities in SUI.
type SuiMoveAbilitySet struct {
	Abilities []SuiMoveAbility `json:"abilities"`
}

// SuiMoveAbility defines a Move ability in SUI.
type SuiMoveAbility string

var (
	// Copy is a Move ability in SUI.
	Copy SuiMoveAbility = "Copy"
	// Drop is a Move ability in SUI.
	Drop SuiMoveAbility = "Drop"
	// Store is a Move ability in SUI.
	Store SuiMoveAbility = "Store"
	// Key is a Move ability in SUI.
	Key SuiMoveAbility = "Key"
)

// SuiMoveModuleID defines a Move module ID in SUI.
type SuiMoveModuleID struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// SuiMoveNormalizedField defines a normalized Move field in SUI.
type SuiMoveNormalizedField struct {
	Name string                       `json:"name"`
	Type SuiMoveNormalizedTypeWrapper `json:"type"`
}

// SuiMoveNormalizedFunction defines a normalized Move function in SUI.
type SuiMoveNormalizedFunction struct {
	Visibility     SuiMoveVisibility               `json:"visibility"`
	IsEntry        bool                            `json:"isEntry"`
	TypeParameters []SuiMoveAbilitySet             `json:"typeParameters"`
	Parameters     []*SuiMoveNormalizedTypeWrapper `json:"parameters"`
	Return         []*SuiMoveNormalizedTypeWrapper `json:"return"`
}

// SuiMoveVisibility defines the visibility of a Move entity in SUI.
type SuiMoveVisibility string

var (
	// Private is a Move visibility in SUI.
	Private SuiMoveVisibility = "Private"
	// Public is a Move visibility in SUI.
	Public SuiMoveVisibility = "Public"
	// Friend is a Move visibility in SUI.
	Friend SuiMoveVisibility = "Friend"
)

// SuiMoveNormalizedModule defines a normalized Move module in SUI.
type SuiMoveNormalizedModule struct {
	FileFormatVersion int                                  `json:"fileFormatVersion"`
	Address           string                               `json:"address"`
	Name              string                               `json:"name"`
	Friends           []SuiMoveModuleID                    `json:"friends"`
	Structs           map[string]SuiMoveNormalizedStruct   `json:"structs"`
	ExposedFunctions  map[string]SuiMoveNormalizedFunction `json:"exposedFunctions"`
}

// SuiMoveNormalizedStruct defines a normalized Move struct in SUI.
type SuiMoveNormalizedStruct struct {
	Abilities      SuiMoveAbilitySet            `json:"abilities"`
	TypeParameters []SuiMoveStructTypeParameter `json:"typeParameters"`
	Fields         []SuiMoveNormalizedField     `json:"fields"`
}

// SuiMoveStructTypeParameter defines a struct type parameter in SUI Move.
type SuiMoveStructTypeParameter struct {
	Constraints SuiMoveAbilitySet `json:"constraints"`
	IsPhantom   bool              `json:"isPhantom"`
}

// MoveCallSuiTransaction defines a Move call transaction in SUI.
type MoveCallSuiTransaction struct {
	Package       string               `json:"package"`
	Module        string               `json:"module"`
	Function      string               `json:"function"`
	TypeArguments []*string            `json:"type_arguments,omitempty"`
	Arguments     []SuiArgumentWrapper `json:"arguments"`
}

// SuiSystemStateSummary defines the summary of the system state in SUI.
type SuiSystemStateSummary struct {
	Epoch                                 string                `json:"epoch"`
	ProtocolVersion                       string                `json:"protocolVersion"`
	SystemStateVersion                    string                `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  string                `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       string                `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     string                `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                string                `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            string                `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                string                `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       string                `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 string                `json:"epochStartTimestampMs"`
	EpochDurationMs                       string                `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                string                `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     string                `json:"maxValidatorCount"`
	MinValidatorJoiningStake              string                `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            string                `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        string                `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          string                `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   string                `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       string                `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount string                `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              string                `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              int                   `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            string                `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsID             string                `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           string                `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []string              `json:"pendingRemovals"`
	StakingPoolMappingsID                 string                `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               string                `json:"stakingPoolMappingsSize"`
	InactivePoolsID                       string                `json:"inactivePoolsId"`
	InactivePoolsSize                     string                `json:"inactivePoolsSize"`
	ValidatorCandidatesID                 string                `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               string                `json:"validatorCandidatesSize"`
	AtRiskValidators                      [][2]string           `json:"atRiskValidators"`
	ValidatorReportRecords                [][2]interface{}      `json:"validatorReportRecords"`
}

// SuiValidatorSummary defines the summary of a validator in SUI.
type SuiValidatorSummary struct {
	SuiAddress                   string  `json:"suiAddress"`
	ProtocolPubkeyBytes          string  `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes           string  `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes            string  `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes       string  `json:"proofOfPossessionBytes"`
	Name                         string  `json:"name"`
	Description                  string  `json:"description"`
	ImageURL                     string  `json:"imageUrl"`
	ProjectURL                   string  `json:"projectUrl"`
	NetAddress                   string  `json:"netAddress"`
	P2pAddress                   string  `json:"p2pAddress"`
	PrimaryAddress               string  `json:"primaryAddress"`
	WorkerAddress                string  `json:"workerAddress"`
	NextEpochProtocolPubkeyBytes *string `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   *string `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  *string `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   *string `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          *string `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          *string `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      *string `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       *string `json:"nextEpochWorkerAddress"`
	VotingPower                  string  `json:"votingPower"`
	OperationCapID               string  `json:"operationCapId"`
	GasPrice                     string  `json:"gasPrice"`
	CommissionRate               string  `json:"commissionRate"`
	NextEpochStake               string  `json:"nextEpochStake"`
	NextEpochGasPrice            string  `json:"nextEpochGasPrice"`
	NextEpochCommissionRate      string  `json:"nextEpochCommissionRate"`
	StakingPoolID                string  `json:"stakingPoolId"`
	StakingPoolActivationEpoch   *string `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch *string `json:"stakingPoolDeactivationEpoch"`
	StakingPoolSuiBalance        string  `json:"stakingPoolSuiBalance"`
	RewardsPool                  string  `json:"rewardsPool"`
	PoolTokenBalance             string  `json:"poolTokenBalance"`
	PendingStake                 string  `json:"pendingStake"`
	PendingTotalSuiWithdraw      string  `json:"pendingTotalSuiWithdraw"`
	PendingPoolTokenWithdraw     string  `json:"pendingPoolTokenWithdraw"`
	ExchangeRatesID              string  `json:"exchangeRatesId"`
	ExchangeRatesSize            string  `json:"exchangeRatesSize"`
}

// SuiTransactionBlockBuilderMode defines the mode for building transaction blocks in SUI.
type SuiTransactionBlockBuilderMode string

const (
	// Commit is the default mode for building transaction blocks in SUI.
	Commit SuiTransactionBlockBuilderMode = "Commit"
	// DevInspect is used for simulating transaction execution without committing it to the blockchain.
	DevInspect SuiTransactionBlockBuilderMode = "DevInspect"
)

// CoinSupply defines the supply of a coin in SUI.
type CoinSupply struct {
	Value string `json:"value"`
}

// TransactionBlockBytes defines the bytes of a transaction block in SUI.
type TransactionBlockBytes struct {
	Gas          []SuiObjectRef           `json:"gas"`
	InputObjects []InputObjectKindWrapper `json:"inputObjects"`
	TxBytes      string                   `json:"txBytes"`
}

// SuiTransactionBlockResponseOptions defines options for transaction block responses in SUI.
type SuiTransactionBlockResponseOptions struct {
	ShowInput          bool `json:"showInput,omitempty"`
	ShowEffects        bool `json:"showEffects,omitempty"`
	ShowEvents         bool `json:"showEvents,omitempty"`
	ShowObjectChanges  bool `json:"showObjectChanges,omitempty"`
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
	ShowRawInput       bool `json:"showRawInput,omitempty"`
}

// SuiTransactionBlockResponseQuery defines a query for transaction block responses in SUI.
type SuiTransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

// TypeOrigin defines the origin of a type in SUI.
type TypeOrigin struct {
	ModuleName   string `json:"module_name"`
	DataTypeName string `json:"datatype_name"`
	Package      string `json:"package"`
}

// UpgradeInfo contains information about an upgrade in SUI.
type UpgradeInfo struct {
	UpgradedID      string `json:"upgraded_id"`
	UpgradedVersion int64  `json:"upgraded_version"`
}

// ValidatorApy defines the APY (Annual Percentage Yield) for a specific validator in SUI.
type ValidatorApy struct {
	Address string  `json:"address"`
	APY     float64 `json:"apy"`
}

// ValidatorsApy defines a collection of validator APYs for a specific epoch in SUI.
type ValidatorsApy struct {
	APYs  []ValidatorApy `json:"apys"`
	Epoch string         `json:"epoch"`
}

// ZkLoginAuthenticator defines the structure for a ZK login authenticator in SUI.
type ZkLoginAuthenticator struct {
	Inputs        ZkLoginInputs    `json:"inputs"`
	MaxEpoch      string           `json:"maxEpoch"`
	UserSignature SignatureWrapper `json:"userSignature"`
}

// ZkLoginInputs defines the inputs required for a ZK login.
type ZkLoginInputs struct {
	AddressSeed      string       `json:"addressSeed"`
	HeaderBase64     string       `json:"headerBase64"`
	IssBase64Details Claim        `json:"issBase64Details"`
	ProofPoints      ZkLoginProof `json:"proofPoints"`
}

// ZkLoginProof defines the proof points used in a ZK login.
type ZkLoginProof struct {
	A []string   `json:"a"`
	B [][]string `json:"b"`
	C []string   `json:"c"`
}

// PaginatedCheckpoints defines a paginated response for checkpoints in SUI.
type PaginatedCheckpoints struct {
	Data        []Checkpoint `json:"data"`
	HasNextPage bool         `json:"hasNextPage"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
}

// ObjectValueKind defines the kind of access to an object in SUI.
type ObjectValueKind string

const (
	// ByImmutableReference indicates that the object is accessed by an immutable reference.
	ByImmutableReference ObjectValueKind = "ByImmutableReference" // ByImmutableReference indicates that the object is accessed by an immutable reference.
	// ByMutableReference indicates that the object is accessed by a mutable reference.
	ByMutableReference ObjectValueKind = "ByMutableReference" // ByMutableReference indicates that the object is accessed by a mutable reference.
	// ByValue indicates that the object is accessed by value.
	ByValue ObjectValueKind = "ByValue" // ByValue indicates that the object is accessed by value.
)

// ResolvedNameServiceNames defines a paginated response for resolved name service names in SUI.
type ResolvedNameServiceNames struct {
	Data        []string `json:"data"`
	NextCursor  *string  `json:"nextCursor,omitempty"`
	HasNextPage bool     `json:"hasNextPage"`
}

// EpochInfo defines information about a specific epoch in SUI.
type EpochInfo struct {
	Epoch                  string                `json:"epoch"`
	Validators             []SuiValidatorSummary `json:"validators"`
	EpochTotalTransactions string                `json:"epochTotalTransactions"`
	FirstCheckpointID      string                `json:"firstCheckpointId"`
	EpochStartTimestamp    string                `json:"epochStartTimestamp"`
	EndOfEpochInfo         *EndOfEpochInfo       `json:"endOfEpochInfo"`
	ReferenceGasPrice      *uint64               `json:"referenceGasPrice"`
}

// EndOfEpochInfo defines information about the end of an epoch in SUI.
type EndOfEpochInfo struct {
	LastCheckpointID             string `json:"lastCheckpointId"`
	EpochEndTimestamp            string `json:"epochEndTimestamp"`
	ProtocolVersion              string `json:"protocolVersion"`
	ReferenceGasPrice            string `json:"referenceGasPrice"`
	TotalStake                   string `json:"totalStake"`
	StorageFundReinvestment      string `json:"storageFundReinvestment"`
	StorageCharge                string `json:"storageCharge"`
	StorageRebate                string `json:"storageRebate"`
	StorageFundBalance           string `json:"storageFundBalance"`
	StakeSubsidyAmount           string `json:"stakeSubsidyAmount"`
	TotalGasFees                 string `json:"totalGasFees"`
	TotalStakeRewardsDistributed string `json:"totalStakeRewardsDistributed"`
	LeftoverStorageFundInflow    string `json:"leftoverStorageFundInflow"`
}

// EpochPage defines a paginated response for epochs in SUI.
type EpochPage struct {
	Data        []EpochInfo `json:"data"`
	NextCursor  string      `json:"nextCursor,omitempty"`
	HasNextPage bool        `json:"hasNextPage"`
}

// DynamicFieldPage defines a paginated response for dynamic fields in SUI.
type DynamicFieldPage struct {
	Data        []DynamicFieldInfo `json:"data"`
	NextCursor  *string            `json:"nextCursor,omitempty"`
	HasNextPage bool               `json:"hasNextPage"`
}

// CheckpointPage defines a paginated response for checkpoints in SUI.
type CheckpointPage struct {
	Data        []Checkpoint `json:"data"`
	NextCursor  *string      `json:"nextCursor,omitempty"`
	HasNextPage bool         `json:"hasNextPage"`
}

// SuiMoveNormalizedModules defines a collection of normalized Move modules in SUI.
type SuiMoveNormalizedModules map[string]SuiMoveNormalizedModule

// ProgrammableTransaction defines a transaction that can be executed on the SUI blockchain.
type ProgrammableTransaction struct {
	Transactions []SuiTransactionWrapper `json:"transactions"`
	Inputs       []SuiCallArgWrapper     `json:"inputs"`
}

// GetPastObjectRequest defines a request to get a past object by its ID and version.
type GetPastObjectRequest struct {
	ObjectID string `json:"objectId"`
	Version  string `json:"version"`
}

// LoadedChildObject defines a child object loaded from the SUI blockchain.
type LoadedChildObject struct {
	ObjectID       string `json:"objectId"`
	SequenceNumber string `json:"sequenceNumber"`
}

// LoadedChildObjectsResponse defines the response for loaded child objects.
type LoadedChildObjectsResponse struct {
	LoadedChildObjects []LoadedChildObject `json:"loadedChildObjects"`
}

// MoveCallParams defines the parameters for a Move call in SUI.
type MoveCallParams struct {
	Arguments       []interface{} `json:"arguments"`
	Function        string        `json:"function"`
	Module          string        `json:"module"`
	PackageObjectID string        `json:"packageObjectId"`
	TypeArguments   []string      `json:"typeArguments,omitempty"`
}

// TransferObjectParams defines the parameters for transferring an object in SUI.
type TransferObjectParams struct {
	ObjectID  string `json:"objectId"`
	Recipient string `json:"recipient"`
}
