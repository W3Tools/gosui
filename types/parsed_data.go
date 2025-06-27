package types

import (
	"encoding/json"
	"errors"
)

// SuiParsedData is an interface that defines a Sui parsed data type.
type SuiParsedData interface {
	isSuiParsedData()
}

// SuiParsedMoveObjectData defines a Sui Move object parsed data.
type SuiParsedMoveObjectData struct {
	DataType          string            `json:"dataType"`
	Type              string            `json:"type"`
	HasPublicTransfer bool              `json:"hasPublicTransfer"`
	Fields            MoveStructWrapper `json:"fields"`
}

// SuiParsedPackageData defines a Sui package parsed data.
type SuiParsedPackageData struct {
	DataType     string                  `json:"dataType"`
	Disassembled *map[string]interface{} `json:"disassembled,omitempty"`
}

// isSuiParsedData implements the SuiParsedData interface for SuiParsedMoveObjectData.
func (SuiParsedMoveObjectData) isSuiParsedData() {}

// isSuiParsedData implements the SuiParsedData interface for SuiParsedPackageData.
func (SuiParsedPackageData) isSuiParsedData() {}

// SuiParsedDataWrapper defines a wrapper for SuiParsedData to support custom JSON marshaling and unmarshaling.
type SuiParsedDataWrapper struct {
	SuiParsedData
}

// UnmarshalJSON decodes JSON data into a SuiParsedDataWrapper.
func (w *SuiParsedDataWrapper) UnmarshalJSON(data []byte) error {
	type DataType struct {
		DataType string `json:"dataType"`
	}

	var dataType DataType
	if err := json.Unmarshal(data, &dataType); err != nil {
		return err
	}

	switch dataType.DataType {
	case "moveObject":
		var p SuiParsedMoveObjectData
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		w.SuiParsedData = p
		return nil
	case "package":
		var p SuiParsedPackageData
		if err := json.Unmarshal(data, &p); err != nil {
			return err
		}
		w.SuiParsedData = p
		return nil
	default:
		return errors.New("unknown SuiParsedData type")
	}
}

// MarshalJSON encodes a SuiParsedDataWrapper into JSON.
func (w *SuiParsedDataWrapper) MarshalJSON() ([]byte, error) {
	switch data := w.SuiParsedData.(type) {
	case SuiParsedMoveObjectData:
		return json.Marshal(SuiParsedMoveObjectData{DataType: data.DataType, Type: data.Type, HasPublicTransfer: data.HasPublicTransfer, Fields: data.Fields})
	case SuiParsedPackageData:
		return json.Marshal(SuiParsedPackageData{DataType: data.DataType, Disassembled: data.Disassembled})
	default:
		return nil, errors.New("unknown SuiParsedData type")
	}
}
