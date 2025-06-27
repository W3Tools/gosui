package types

import (
	"encoding/json"
	"errors"
)

// SuiObjectChange is an interface for object changes in Sui.
type SuiObjectChange interface {
	isSuiObjectChange()
}

// SuiObjectChangePublished defines a published object change in Sui.
type SuiObjectChangePublished struct {
	Type      string   `json:"type"`
	PackageID string   `json:"packageId"`
	Version   string   `json:"version"`
	Digest    string   `json:"digest"`
	Modules   []string `json:"modules"`
}

// SuiObjectChangeTransferred defines a transferred object change in Sui.
type SuiObjectChangeTransferred struct {
	Type       string              `json:"type"`
	Sender     string              `json:"sender"`
	Recipient  *ObjectOwnerWrapper `json:"recipient,omitempty"`
	ObjectType string              `json:"objectType"`
	ObjectID   string              `json:"objectId"`
	Version    string              `json:"version"`
	Digest     string              `json:"digest"`
}

// SuiObjectChangeMutated defines a mutated object change in Sui.
type SuiObjectChangeMutated struct {
	Type            string              `json:"type"`
	Sender          string              `json:"sender"`
	Owner           *ObjectOwnerWrapper `json:"owner"`
	ObjectType      string              `json:"objectType"`
	ObjectID        string              `json:"objectId"`
	Version         string              `json:"version"`
	PreviousVersion string              `json:"previousVersion"`
	Digest          string              `json:"digest"`
}

// SuiObjectChangeDeleted defines a deleted object change in Sui.
type SuiObjectChangeDeleted struct {
	Type       string `json:"type"`
	Sender     string `json:"sender"`
	ObjectType string `json:"objectType"`
	ObjectID   string `json:"objectId"`
	Version    string `json:"version"`
}

// SuiObjectChangeWrapped defines a wrapped object change in Sui.
type SuiObjectChangeWrapped struct {
	Type       string `json:"type"`
	Sender     string `json:"sender"`
	ObjectType string `json:"objectType"`
	ObjectID   string `json:"objectId"`
	Version    string `json:"version"`
}

// SuiObjectChangeCreated defines a created object change in Sui.
type SuiObjectChangeCreated struct {
	Type       string              `json:"type"`
	Sender     string              `json:"sender"`
	Owner      *ObjectOwnerWrapper `json:"owner,omitempty"`
	ObjectType string              `json:"objectType"`
	ObjectID   string              `json:"objectId"`
	Version    string              `json:"version"`
	Digest     string              `json:"digest"`
}

func (SuiObjectChangePublished) isSuiObjectChange()   {}
func (SuiObjectChangeTransferred) isSuiObjectChange() {}
func (SuiObjectChangeMutated) isSuiObjectChange()     {}
func (SuiObjectChangeDeleted) isSuiObjectChange()     {}
func (SuiObjectChangeWrapped) isSuiObjectChange()     {}
func (SuiObjectChangeCreated) isSuiObjectChange()     {}

// SuiObjectChangeWrapper defines a wrapper for SuiObjectChange.
type SuiObjectChangeWrapper struct {
	SuiObjectChange
}

// UnmarshalJSON decodes a SuiObjectChangeWrapper from JSON.
func (w *SuiObjectChangeWrapper) UnmarshalJSON(data []byte) error {
	type Type struct {
		Type string `json:"type"`
	}

	var changeType Type
	if err := json.Unmarshal(data, &changeType); err != nil {
		return err
	}

	switch changeType.Type {
	case "published":
		var c SuiObjectChangePublished
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	case "transferred":
		var c SuiObjectChangeTransferred
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	case "mutated":
		var c SuiObjectChangeMutated
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	case "deleted":
		var c SuiObjectChangeDeleted
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	case "wrapped":
		var c SuiObjectChangeWrapped
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	case "created":
		var c SuiObjectChangeCreated
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		w.SuiObjectChange = c
	default:
		return errors.New("unknown SuiObjectChange type")
	}
	return nil
}

// MarshalJSON encodes a SuiObjectChangeWrapper to JSON.
func (w SuiObjectChangeWrapper) MarshalJSON() ([]byte, error) {
	switch change := w.SuiObjectChange.(type) {
	case SuiObjectChangePublished:
		return json.Marshal(change)
	case SuiObjectChangeTransferred:
		return json.Marshal(change)
	case SuiObjectChangeMutated:
		return json.Marshal(change)
	case SuiObjectChangeDeleted:
		return json.Marshal(change)
	case SuiObjectChangeWrapped:
		return json.Marshal(change)
	case SuiObjectChangeCreated:
		return json.Marshal(change)
	default:
		return nil, errors.New("unknown SuiObjectChange type")
	}
}
