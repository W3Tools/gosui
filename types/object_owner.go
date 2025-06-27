package types

import (
	"encoding/json"
	"errors"
)

// ObjectOwner is an interface that represents different types of object owners.
type ObjectOwner interface {
	isObjectOwner()
}

// ObjectOwnerAddressOwner defines an object owner that is an address owner.
type ObjectOwnerAddressOwner struct {
	AddressOwner string `json:"AddressOwner"`
}

// ObjectOwnerObjectOwner defines an object owner that is another object owner.
type ObjectOwnerObjectOwner struct {
	ObjectOwner string `json:"ObjectOwner"`
}

// ObjectOwnerShared defines an object owner that is shared.
type ObjectOwnerShared struct {
	Shared ObjectOwnerSharedData `json:"Shared"`
}

// ObjectOwnerSharedData defines the data structure for shared object owners.
type ObjectOwnerSharedData struct {
	InitialSharedVersion uint64 `json:"initial_shared_version"`
}

// ObjectOwnerImmutable defines an immutable object owner, represented as a string.
type ObjectOwnerImmutable string

func (ObjectOwnerAddressOwner) isObjectOwner() {}
func (ObjectOwnerObjectOwner) isObjectOwner()  {}
func (ObjectOwnerShared) isObjectOwner()       {}
func (ObjectOwnerImmutable) isObjectOwner()    {}

// ObjectOwnerWrapper is a wrapper for ObjectOwner that allows unmarshalling from JSON.
type ObjectOwnerWrapper struct {
	ObjectOwner
}

// UnmarshalJSON custom unmarshaller for ObjectOwnerWrapper
func (w *ObjectOwnerWrapper) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		w.ObjectOwner = ObjectOwnerImmutable(s)
		return nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if addressOwner, ok := obj["AddressOwner"]; ok {
		var o ObjectOwnerAddressOwner
		if err := json.Unmarshal(addressOwner, &o.AddressOwner); err != nil {
			return err
		}
		w.ObjectOwner = o
		return nil
	}

	if objectOwner, ok := obj["ObjectOwner"]; ok {
		var o ObjectOwnerObjectOwner
		if err := json.Unmarshal(objectOwner, &o.ObjectOwner); err != nil {
			return err
		}
		w.ObjectOwner = o
		return nil
	}

	if shared, ok := obj["Shared"]; ok {
		var o ObjectOwnerShared
		if err := json.Unmarshal(shared, &o.Shared); err != nil {
			return err
		}

		w.ObjectOwner = o
		return nil
	}

	return errors.New("unknown ObjectOwner type")
}

// MarshalJSON custom marshaller for ObjectOwnerWrapper
func (w ObjectOwnerWrapper) MarshalJSON() ([]byte, error) {
	switch o := w.ObjectOwner.(type) {
	case ObjectOwnerAddressOwner:
		return json.Marshal(ObjectOwnerAddressOwner{AddressOwner: o.AddressOwner})
	case ObjectOwnerObjectOwner:
		return json.Marshal(ObjectOwnerObjectOwner{ObjectOwner: o.ObjectOwner})
	case ObjectOwnerShared:
		return json.Marshal(ObjectOwnerShared{Shared: o.Shared})
	case ObjectOwnerImmutable:
		return json.Marshal(string(o))
	default:
		return nil, errors.New("unknown ObjectOwner type")
	}
}
