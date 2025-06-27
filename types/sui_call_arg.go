package types

import (
	"encoding/json"
	"errors"
)

// SuiCallArg is an interface that defines a Sui call argument.
type SuiCallArg interface {
	isSuiCallArg()
}

// SuiCallArgImmOrOwnedObject defines an immutable or owned object argument for a Sui call.
type SuiCallArgImmOrOwnedObject struct {
	Type       string `json:"type"`
	Digest     string `json:"digest"`
	ObjectID   string `json:"objectId"`
	ObjectType string `json:"objectType"`
	Version    string `json:"version"`
}

// SuiCallArgSharedObject defines a shared object argument for a Sui call.
type SuiCallArgSharedObject struct {
	Type                 string `json:"type"`
	ObjectType           string `json:"objectType"`
	ObjectID             string `json:"objectId"`
	InitialSharedVersion string `json:"initialSharedVersion"`
	Mutable              bool   `json:"mutable"`
}

// SuiCallArgReceiving defines a receiving object argument for a Sui call.
type SuiCallArgReceiving struct {
	Type       string `json:"type"`
	Digest     string `json:"digest"`
	ObjectID   string `json:"objectId"`
	ObjectType string `json:"objectType"`
	Version    string `json:"version"`
}

// SuiCallArgPure defines a pure argument for a Sui call.
type SuiCallArgPure struct {
	Type      string      `json:"type"`
	ValueType *string     `json:"valueType,omitempty"`
	Value     interface{} `json:"value"`
}

func (SuiCallArgImmOrOwnedObject) isSuiCallArg() {}
func (SuiCallArgSharedObject) isSuiCallArg()     {}
func (SuiCallArgReceiving) isSuiCallArg()        {}
func (SuiCallArgPure) isSuiCallArg()             {}

// SuiCallArgWrapper defines a wrapper for SuiCallArg to support custom JSON marshaling and unmarshaling.
type SuiCallArgWrapper struct {
	SuiCallArg
}

// UnmarshalJSON decodes JSON data into a SuiCallArgWrapper.
func (w *SuiCallArgWrapper) UnmarshalJSON(data []byte) error {
	type Type struct {
		Type       string `json:"type"`
		ObjectType string `json:"objectType"`
	}
	var argType Type
	if err := json.Unmarshal(data, &argType); err != nil {
		return err
	}

	switch argType.Type {
	case "object":
		switch argType.ObjectType {
		case "immOrOwnedObject":
			var a SuiCallArgImmOrOwnedObject
			if err := json.Unmarshal(data, &a); err != nil {
				return err
			}
			w.SuiCallArg = a
		case "sharedObject":
			var a SuiCallArgSharedObject
			if err := json.Unmarshal(data, &a); err != nil {
				return err
			}
			w.SuiCallArg = a
		case "receiving":
			var a SuiCallArgReceiving
			if err := json.Unmarshal(data, &a); err != nil {
				return err
			}
			w.SuiCallArg = a
		default:
			return errors.New("unknown SuiCallArg object type")
		}
	case "pure":
		var a SuiCallArgPure
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		w.SuiCallArg = a
	default:
		return errors.New("unknown SuiCallArg type")
	}

	return nil
}

// MarshalJSON encodes a SuiCallArgWrapper into JSON.
func (w *SuiCallArgWrapper) MarshalJSON() ([]byte, error) {
	switch arg := w.SuiCallArg.(type) {
	case SuiCallArgImmOrOwnedObject:
		return json.Marshal(arg)
	case SuiCallArgSharedObject:
		return json.Marshal(arg)
	case SuiCallArgReceiving:
		return json.Marshal(arg)
	case SuiCallArgPure:
		return json.Marshal(arg)
	default:
		return nil, errors.New("unknown SuiCallArg type")
	}
}
