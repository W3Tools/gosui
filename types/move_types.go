package types

import (
	"encoding/json"
	"fmt"
)

type SuiMoveId struct {
	Id string `json:"id"`
}

type SuiMoveTable struct {
	Fields SuiMoveTableFields `json:"fields"`
	Type   string             `json:"type"`
}

type SuiMoveTableFields struct {
	Id   SuiMoveId `json:"id"`
	Size string    `json:"size"`
}

type SuiMoveString struct {
	Type   string              `json:"type"`
	Fields SuiMoveStringFields `json:"fields"`
}

type SuiMoveStringFields struct {
	Name string `json:"name"`
}

type SuiMoveDynamicField[TypeFields any, TypeName any] struct {
	Id    SuiMoveId                            `json:"id"`
	Name  TypeName                             `json:"name"`
	Value SuiMoveDynamicFieldValue[TypeFields] `json:"value"`
}

type SuiMoveDynamicFieldValue[TypeFields any] struct {
	Fields TypeFields `json:"fields"`
	Type   string     `json:"type"`
}

// Define the type as MoveModule/MoveEventModule. Events emitted, defined on the specified Move module.
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#filtering-event-queries
type MoveEventModuleConfig struct {
	Package string `toml:"Package,omitempty"`
	Module  string `toml:"Module,omitempty"`
}

func (ec *MoveEventModuleConfig) Join() string {
	return fmt.Sprintf("%s::%s", ec.Package, ec.Module)
}

func (ec *MoveEventModuleConfig) JoinEventName(name string) string {
	return fmt.Sprintf("%s::%s::%s", ec.Package, ec.Module, name)
}

// Parsing custom event json
// Reference: https://docs.sui.io/guides/developer/sui-101/using-events#move-event-structure
func ParseEvent[T any](event SuiEvent) (*T, error) {
	jsonBytes, err := json.Marshal(event.ParsedJson)
	if err != nil {
		return nil, err
	}

	parsedJson := new(T)
	if err := json.Unmarshal(jsonBytes, &parsedJson); err != nil {
		return nil, err
	}
	return parsedJson, nil
}
