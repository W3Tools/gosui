package types

import (
	"encoding/json"
	"fmt"
)

// SuiMoveID defines a structure for a Move ID in the Sui blockchain.
type SuiMoveID struct {
	ID string `json:"id"`
}

// SuiMoveTable defines a structure for a Move table in the Sui blockchain.
type SuiMoveTable struct {
	Fields SuiMoveTableFields `json:"fields"`
	Type   string             `json:"type"`
}

// SuiMoveTableFields defines the fields of a Move table in the Sui blockchain.
type SuiMoveTableFields struct {
	ID   SuiMoveID `json:"id"`
	Size string    `json:"size"`
}

// SuiMoveString defines a structure for a Move string in the Sui blockchain.
type SuiMoveString struct {
	Type   string              `json:"type"`
	Fields SuiMoveStringFields `json:"fields"`
}

// SuiMoveStringFields defines the fields of a Move string in the Sui blockchain.
type SuiMoveStringFields struct {
	Name string `json:"name"`
}

// SuiMoveDynamicField defines a structure for a Move dynamic field in the Sui blockchain.
type SuiMoveDynamicField[TypeFields any, TypeName any] struct {
	ID    SuiMoveID                            `json:"id"`
	Name  TypeName                             `json:"name"`
	Value SuiMoveDynamicFieldValue[TypeFields] `json:"value"`
}

// SuiMoveDynamicFieldValue defines the value of a Move dynamic field in the Sui blockchain.
type SuiMoveDynamicFieldValue[TypeFields any] struct {
	Fields TypeFields `json:"fields"`
	Type   string     `json:"type"`
}

// MoveEventModuleConfig define the type as MoveModule/MoveEventModule. Events emitted, defined on the specified Move module.
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#filtering-event-queries
type MoveEventModuleConfig struct {
	Package string `toml:"Package,omitempty"`
	Module  string `toml:"Module,omitempty"`
}

// Join returns the full module path in the format "Package::Module".
func (ec *MoveEventModuleConfig) Join() string {
	return fmt.Sprintf("%s::%s", ec.Package, ec.Module)
}

// JoinEventName returns the full event name in the format "Package::Module::EventName".
func (ec *MoveEventModuleConfig) JoinEventName(name string) string {
	return fmt.Sprintf("%s::%s::%s", ec.Package, ec.Module, name)
}

// ParseEvent parses a SuiEvent's ParsedJSON field into a specified type T.
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#move-event-structure
func ParseEvent[T any](event SuiEvent) (*T, error) {
	jsonBytes, err := json.Marshal(event.ParsedJSON)
	if err != nil {
		return nil, err
	}

	parsedJSON := new(T)
	if err := json.Unmarshal(jsonBytes, &parsedJSON); err != nil {
		return nil, err
	}
	return parsedJSON, nil
}
