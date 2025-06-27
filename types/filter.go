package types

// TransactionFilter defines a structure for filtering transactions based on various criteria.
type TransactionFilter struct {
	Checkpoint        *string                            `json:"Checkpoint,omitempty"`
	MoveFunction      *TransactionFilterMoveFunction     `json:"MoveFunction,omitempty"`
	InputObject       *string                            `json:"InputObject,omitempty"`
	ChangedObject     *string                            `json:"ChangedObject,omitempty"`
	FromAddress       *string                            `json:"FromAddress,omitempty"`
	ToAddress         *string                            `json:"ToAddress,omitempty"`
	FromAndToAddress  *TransactionFilterFromAndToAddress `json:"FromAndToAddress,omitempty"`
	FromOrToAddress   *TransactionFilterFromOrToAddress  `json:"FromOrToAddress,omitempty"`
	TransactionKind   *string                            `json:"TransactionKind,omitempty"`
	TransactionKindIn []*string                          `json:"TransactionKindIn,omitempty"`
}

// SuiObjectDataFilter defines a structure for filtering Sui objects based on various criteria.
type SuiObjectDataFilter struct {
	*SuiObjectDataFilterMatchAll
	*SuiObjectDataFilterMatchAny
	*SuiObjectDataFilterMatchNone
	*SuiObjectDataFilterPackage
	*SuiObjectDataFilterMoveModule
	*SuiObjectDataFilterStructType
	*SuiObjectDataFilterAddressOwner
	*SuiObjectDataFilterObjectOwner
	*SuiObjectDataFilterObjectID
	*SuiObjectDataFilterObjectIds
	*SuiObjectDataFilterVersion
}

// SuiEventFilter defines a structure for filtering Sui events based on various criteria.
type SuiEventFilter struct {
	Sender      *string `json:"Sender,omitempty"`
	Transaction *string `json:"Transaction,omitempty"`
	//Package         *string                         `json:"Package,omitempty"`
	MoveModule      *SuiEventFilterMoveModule      `json:"MoveModule,omitempty"`
	MoveEventType   *string                        `json:"MoveEventType,omitempty"`
	MoveEventModule *SuiEventFilterMoveEventModule `json:"MoveEventModule,omitempty"`
	//MoveEventField  *SuiEventFilter_MoveEventField  `json:"MoveEventField,omitempty"`
	TimeRange *SuiEventFilterTimeRange `json:"TimeRange,omitempty"`
	All       *SuiEventFilters         `json:"All,omitempty"`
	Any       *SuiEventFilters         `json:"Any,omitempty"`
	And       *SuiEventFilters         `json:"And,omitempty"`
	Or        *SuiEventFilters         `json:"Or,omitempty"`
}
