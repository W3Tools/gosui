package types

import (
	"encoding/json"
	"errors"
)

// SuiMoveNormalizedType is an interface for normalized Move types in Sui.
type SuiMoveNormalizedType interface {
	isSuiMoveNormalizedType()
}

// SuiMoveNormalizedTypeString defines a string type for normalized Move types in Sui.
type SuiMoveNormalizedTypeString string

// SuiMoveNormalizedTypeStruct defines a struct type for normalized Move types in Sui.
type SuiMoveNormalizedTypeStruct struct {
	Struct SuiMoveNormalizedTypeStructStruct `json:"Struct"`
}

// SuiMoveNormalizedTypeStructStruct defines the struct details for a normalized Move struct type in Sui.
type SuiMoveNormalizedTypeStructStruct struct {
	Address       string                         `json:"address"`
	Module        string                         `json:"module"`
	Name          string                         `json:"name"`
	TypeArguments []SuiMoveNormalizedTypeWrapper `json:"typeArguments"`
}

// SuiMoveNormalizedTypeVector defines a vector type for normalized Move types in Sui.
type SuiMoveNormalizedTypeVector struct {
	Vector SuiMoveNormalizedTypeWrapper `json:"Vector"`
}

// SuiMoveNormalizedTypeTypeParameter defines a type parameter for normalized Move types in Sui.
type SuiMoveNormalizedTypeTypeParameter struct {
	TypeParameter uint64 `json:"TypeParameter"`
}

// SuiMoveNormalizedTypeReference defines a reference type for normalized Move types in Sui.
type SuiMoveNormalizedTypeReference struct {
	Reference SuiMoveNormalizedTypeWrapper `json:"Reference"`
}

// SuiMoveNormalizedTypeMutableReference defines a mutable reference type for normalized Move types in Sui.
type SuiMoveNormalizedTypeMutableReference struct {
	MutableReference SuiMoveNormalizedTypeWrapper `json:"MutableReference"`
}

func (SuiMoveNormalizedTypeString) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedTypeStruct) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedTypeVector) isSuiMoveNormalizedType()           {}
func (SuiMoveNormalizedTypeTypeParameter) isSuiMoveNormalizedType()    {}
func (SuiMoveNormalizedTypeReference) isSuiMoveNormalizedType()        {}
func (SuiMoveNormalizedTypeMutableReference) isSuiMoveNormalizedType() {}

// SuiMoveNormalizedTypeWrapper defines a wrapper for SuiMoveNormalizedType.
type SuiMoveNormalizedTypeWrapper struct {
	SuiMoveNormalizedType
}

// UnmarshalJSON decodes a SuiMoveNormalizedTypeWrapper from JSON.
func (w *SuiMoveNormalizedTypeWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.SuiMoveNormalizedType = SuiMoveNormalizedTypeString(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if _struct, ok := obj["Struct"]; ok {
		var s SuiMoveNormalizedTypeStruct
		if err := json.Unmarshal(_struct, &s.Struct); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if vector, ok := obj["Vector"]; ok {
		var s SuiMoveNormalizedTypeVector
		if err := json.Unmarshal(vector, &s.Vector); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if typeParameter, ok := obj["TypeParameter"]; ok {
		var s SuiMoveNormalizedTypeTypeParameter
		if err := json.Unmarshal(typeParameter, &s.TypeParameter); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if reference, ok := obj["Reference"]; ok {
		var s SuiMoveNormalizedTypeReference
		if err := json.Unmarshal(reference, &s.Reference); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	if mutableReference, ok := obj["MutableReference"]; ok {
		var s SuiMoveNormalizedTypeMutableReference
		if err := json.Unmarshal(mutableReference, &s.MutableReference); err != nil {
			return err
		}
		w.SuiMoveNormalizedType = s
		return nil
	}

	return errors.New("unknown SuiMoveNormalizedType type")
}

// MarshalJSON encodes a SuiMoveNormalizedTypeWrapper to JSON.
func (w SuiMoveNormalizedTypeWrapper) MarshalJSON() ([]byte, error) {
	switch t := w.SuiMoveNormalizedType.(type) {
	case SuiMoveNormalizedTypeString:
		return json.Marshal(string(t))
	case SuiMoveNormalizedTypeStruct:
		return json.Marshal(SuiMoveNormalizedTypeStruct{Struct: t.Struct})
	case SuiMoveNormalizedTypeVector:
		return json.Marshal(SuiMoveNormalizedTypeVector{Vector: t.Vector})
	case SuiMoveNormalizedTypeTypeParameter:
		return json.Marshal(SuiMoveNormalizedTypeTypeParameter{TypeParameter: t.TypeParameter})
	case SuiMoveNormalizedTypeReference:
		return json.Marshal(SuiMoveNormalizedTypeReference{Reference: t.Reference})
	case SuiMoveNormalizedTypeMutableReference:
		return json.Marshal(SuiMoveNormalizedTypeMutableReference{MutableReference: t.MutableReference})
	default:
		return nil, errors.New("unknown SuiMoveNormalizedType")
	}
}
