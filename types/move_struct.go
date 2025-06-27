package types

import (
	"encoding/json"
	"errors"
)

// MoveStruct is an interface that defines the structure of a Move struct in the Move programming language.
type MoveStruct interface {
	isMoveStruct()
	isMoveValue()
}

// MoveStructMoveValue defines a structure for a Move struct that contains a slice of MoveValueWrapper.
type MoveStructMoveValue []MoveValueWrapper

// MoveStructFieldsType defines a structure for a Move struct that contains a type and a map of fields.
type MoveStructFieldsType struct {
	Type   string                      `json:"type"`
	Fields map[string]MoveValueWrapper `json:"fields"`
}

// MoveStructMap defines a structure for a Move struct that contains a map of string keys to MoveValueWrapper values.
type MoveStructMap map[string]MoveValueWrapper

func (MoveStructMoveValue) isMoveStruct()  {}
func (MoveStructFieldsType) isMoveStruct() {}
func (MoveStructMap) isMoveStruct()        {}

func (MoveStructMoveValue) isMoveValue()  {}
func (MoveStructFieldsType) isMoveValue() {}
func (MoveStructMap) isMoveValue()        {}

// MoveStructWrapper defines a wrapper for the MoveStruct interface to handle different kinds of Move structs.
type MoveStructWrapper struct {
	MoveStruct
}

// UnmarshalJSON implements the json.Unmarshaler interface for MoveStructWrapper.
func (w *MoveStructWrapper) UnmarshalJSON(data []byte) error {
	var mvs MoveStructMoveValue
	if err := json.Unmarshal(data, &mvs); err == nil {
		w.MoveStruct = mvs
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["fields"]; ok {
		var ms MoveStructFieldsType
		if err := json.Unmarshal(data, &ms); err != nil {
			return err
		}
		w.MoveStruct = ms
		return nil
	}
	var ms MoveStructMap
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}
	w.MoveStruct = ms
	return nil
}

// MarshalJSON implements the json.Marshaler interface for MoveStructWrapper.
func (w MoveStructWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.MoveStruct.(type) {
	case MoveStructMoveValue:
		return json.Marshal([]MoveValueWrapper(v))
	case MoveStructFieldsType:
		return json.Marshal(MoveStructFieldsType{Fields: v.Fields, Type: v.Type})
	case MoveStructMap:
		return json.Marshal(v)
	default:
		return nil, errors.New("unknown MoveStruct type")
	}
}

// MoveValue is an interface that defines the structure of a Move value in the Move programming language.
type MoveValue interface {
	isMoveValue()
}

// MoveNumberValue defines a numeric value in the Move programming language, represented as a uint64.
type MoveNumberValue uint64

// MoveBooleanValue defines a boolean value in the Move programming language.
type MoveBooleanValue bool

// MoveStringValue defines a string value in the Move programming language.
type MoveStringValue string

// MoveValueMoveValues defines a slice of MoveValueWrapper, which can contain multiple Move values.
type MoveValueMoveValues []MoveValueWrapper

// MoveIDValue defines a structure for a Move value that contains an ID, typically used to reference objects in the Move programming language.
type MoveIDValue struct {
	ID string `json:"id"`
}

// MoveStructValue defines a structure for a Move value that contains a MoveStruct, which can be a Move struct or a map of fields.
type MoveStructValue MoveStruct

func (MoveNumberValue) isMoveValue()     {}
func (MoveBooleanValue) isMoveValue()    {}
func (MoveStringValue) isMoveValue()     {}
func (MoveValueMoveValues) isMoveValue() {}
func (MoveIDValue) isMoveValue()         {}

// MoveValueWrapper defines a wrapper for the MoveValue interface to handle different kinds of Move values.
type MoveValueWrapper struct {
	MoveValue
}

// UnmarshalJSON implements the json.Unmarshaler interface for MoveValueWrapper.
func (w *MoveValueWrapper) UnmarshalJSON(data []byte) error {
	var num uint64
	if err := json.Unmarshal(data, &num); err == nil {
		w.MoveValue = MoveNumberValue(num)
		return nil
	}
	var bol bool
	if err := json.Unmarshal(data, &bol); err == nil {
		w.MoveValue = MoveBooleanValue(bol)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		w.MoveValue = MoveStringValue(str)
		return nil
	}

	var mvs []MoveValueWrapper
	if err := json.Unmarshal(data, &mvs); err == nil {
		w.MoveValue = MoveValueMoveValues(mvs)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if _, ok := obj["id"]; ok {
		var mid MoveIDValue
		if err := json.Unmarshal(data, &mid); err == nil {
			w.MoveValue = mid
			return nil
		}
	} else {
		var ms MoveStructWrapper
		if err := json.Unmarshal(data, &ms); err == nil {
			w.MoveValue = MoveStructValue(ms)
			return nil
		}
	}

	return errors.New("unknown MoveValue type")
}

// MarshalJSON implements the json.Marshaler interface for MoveValueWrapper.
func (w MoveValueWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.MoveValue.(type) {
	case MoveNumberValue:
		return json.Marshal(uint64(v))
	case MoveBooleanValue:
		return json.Marshal(bool(v))
	case MoveStringValue:
		return json.Marshal(string(v))
	case MoveIDValue:
		return json.Marshal(MoveIDValue{ID: v.ID})
	case MoveStructValue:
		return json.Marshal(MoveStruct(v))
	case MoveValueMoveValues:
		return json.Marshal(v)
	default:
		return nil, errors.New("unknown MoveValue type")
	}
}

// SuiMoveFunctionArgType is an interface that defines the structure of a Move function argument type in the Sui blockchain.
type SuiMoveFunctionArgType interface {
	isSuiMoveFunctionArgType()
}

// SuiMoveFunctionArgStringType defines a string type for Move function arguments in the Sui blockchain.
type SuiMoveFunctionArgStringType string

// SuiMoveFunctionArgObjectType defines an object type for Move function arguments in the Sui blockchain.
type SuiMoveFunctionArgObjectType struct {
	Object ObjectValueKind `json:"Object"`
}

func (SuiMoveFunctionArgStringType) isSuiMoveFunctionArgType() {}
func (SuiMoveFunctionArgObjectType) isSuiMoveFunctionArgType() {}

// SuiMoveFunctionArgTypeWrapper defines a wrapper for the SuiMoveFunctionArgType interface to handle different kinds of Move function argument types.
type SuiMoveFunctionArgTypeWrapper struct {
	SuiMoveFunctionArgType
}

// UnmarshalJSON implements the json.Unmarshaler interface for SuiMoveFunctionArgTypeWrapper.
func (w *SuiMoveFunctionArgTypeWrapper) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		w.SuiMoveFunctionArgType = SuiMoveFunctionArgStringType(str)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _, ok := obj["Object"]; ok {
		var at SuiMoveFunctionArgObjectType
		if err := json.Unmarshal(data, &at); err != nil {
			return err
		}

		w.SuiMoveFunctionArgType = at
		return nil
	}

	return errors.New("unknown SuiMoveFunctionArgType type")
}

// MarshalJSON implements the json.Marshaler interface for SuiMoveFunctionArgTypeWrapper.
func (w SuiMoveFunctionArgTypeWrapper) MarshalJSON() ([]byte, error) {
	switch v := w.SuiMoveFunctionArgType.(type) {
	case SuiMoveFunctionArgStringType:
		return json.Marshal(string(v))
	case SuiMoveFunctionArgObjectType:
		return json.Marshal(SuiMoveFunctionArgObjectType{Object: v.Object})
	default:
		return nil, errors.New("unknown SuiMoveFunctionArgType type")
	}
}
