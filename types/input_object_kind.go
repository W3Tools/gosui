package types

import (
	"encoding/json"
	"errors"
)

// InputObjectKind is an interface that defines the kind of input object in a transaction.
type InputObjectKind interface {
	isInputObjectKind()
}

// InputObjectKindMovePackage defines a kind of input object that represents a Move package.
type InputObjectKindMovePackage struct {
	MovePackage string `json:"MovePackage"`
}

// InputObjectKindImmOrOwnedMoveObject defines a kind of input object that represents an immutable or owned Move object.
type InputObjectKindImmOrOwnedMoveObject struct {
	ImmOrOwnedMoveObject SuiObjectRef `json:"ImmOrOwnedMoveObject"`
}

// InputObjectKindSharedMoveObject defines a kind of input object that represents a shared Move object.
type InputObjectKindSharedMoveObject struct {
	SharedMoveObject KindSharedMoveObject `json:"SharedMoveObject"`
}

// KindSharedMoveObject defines the structure for a shared Move object in a transaction.
type KindSharedMoveObject struct {
	ID                   string `json:"id"`
	InitialSharedVersion string `json:"initial_shared_version"`
	Mutable              bool   `json:"mutable,omitempty"`
}

func (InputObjectKindMovePackage) isInputObjectKind()          {}
func (InputObjectKindImmOrOwnedMoveObject) isInputObjectKind() {}
func (InputObjectKindSharedMoveObject) isInputObjectKind()     {}

// InputObjectKindWrapper defines a wrapper for the InputObjectKind interface to handle different kinds of input objects.
type InputObjectKindWrapper struct {
	InputObjectKind
}

// UnmarshalJSON implements the json.Unmarshaler interface for InputObjectKindWrapper.
func (w *InputObjectKindWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["MovePackage"] != nil:
		var s InputObjectKindMovePackage
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	case obj["ImmOrOwnedMoveObject"] != nil:
		var s InputObjectKindImmOrOwnedMoveObject
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	case obj["SharedMoveObject"] != nil:
		var s InputObjectKindSharedMoveObject
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.InputObjectKind = s
	default:
		return errors.New("unknown InputObjectKind type")
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface for InputObjectKindWrapper.
func (w *InputObjectKindWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.InputObjectKind.(type) {
	case InputObjectKindMovePackage:
		return json.Marshal(InputObjectKindMovePackage{
			MovePackage: t.MovePackage,
		})
	case InputObjectKindImmOrOwnedMoveObject:
		return json.Marshal(InputObjectKindImmOrOwnedMoveObject{
			ImmOrOwnedMoveObject: t.ImmOrOwnedMoveObject,
		})
	case InputObjectKindSharedMoveObject:
		return json.Marshal(InputObjectKindSharedMoveObject{
			SharedMoveObject: t.SharedMoveObject,
		})
	default:
		return nil, errors.New("unknown InputObjectKind type")
	}
}
