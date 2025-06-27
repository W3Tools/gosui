package types

import (
	"encoding/json"
	"errors"
)

// ObjectRead is an interface that represents different types of object read results.
type ObjectRead interface {
	isObjectRead()
}

// ObjectReadVersionFound defines a structure for a successful object read with version found.
type ObjectReadVersionFound struct {
	Details SuiObjectData `json:"detail"`
	Status  string        `json:"status"`
}

// ObjectReadObjectNotExists defines a structure for an object read where the object does not exist.
type ObjectReadObjectNotExists struct {
	Details string `json:"detail"`
	Status  string `json:"status"`
}

// ObjectReadObjectDeleted defines a structure for an object read where the object has been deleted.
type ObjectReadObjectDeleted struct {
	Details SuiObjectRef `json:"detail"`
	Status  string       `json:"status"`
}

// ObjectReadVersionNotFound defines a structure for an object read where the requested version was not found.
type ObjectReadVersionNotFound struct {
	Details [2]string `json:"detail"`
	Status  string    `json:"status"`
}

// ObjectReadVersionTooHigh defines a structure for an object read where the requested version is too high.
type ObjectReadVersionTooHigh struct {
	Details VersionTooHighDetails `json:"detail"`
	Status  string                `json:"status"`
}

// VersionTooHighDetails defines the details for a version too high error.
type VersionTooHighDetails struct {
	AskedVersion  string `json:"asked_version"`
	LatestVersion string `json:"latest_version"`
	ObjectID      string `json:"object_id"`
}

func (ObjectReadVersionFound) isObjectRead()    {}
func (ObjectReadObjectNotExists) isObjectRead() {}
func (ObjectReadObjectDeleted) isObjectRead()   {}
func (ObjectReadVersionNotFound) isObjectRead() {}
func (ObjectReadVersionTooHigh) isObjectRead()  {}

// ObjectReadWrapper is a wrapper for ObjectRead that allows unmarshalling from JSON.
type ObjectReadWrapper struct {
	ObjectRead
}

// UnmarshalJSON custom unmarshaller for ObjectReadWrapper
func (w *ObjectReadWrapper) UnmarshalJSON(data []byte) error {
	type Status struct {
		Status string `json:"status"`
	}

	var status Status
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	switch status.Status {
	case "VersionFound":
		var or ObjectReadVersionFound
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "ObjectNotExists":
		var or ObjectReadObjectNotExists
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "ObjectDeleted":
		var or ObjectReadObjectDeleted
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "VersionNotFound":
		var or ObjectReadVersionNotFound
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	case "VersionTooHigh":
		var or ObjectReadVersionTooHigh
		if err := json.Unmarshal(data, &or); err != nil {
			return err
		}
		w.ObjectRead = or
	default:
		return errors.New("unknown ObjectRead type")
	}

	return nil
}

// MarshalJSON custom marshaller for ObjectReadWrapper
func (w ObjectReadWrapper) MarshalJSON() ([]byte, error) {
	switch obj := w.ObjectRead.(type) {
	case ObjectReadVersionFound:
		return json.Marshal(ObjectReadVersionFound{Details: obj.Details, Status: obj.Status})
	case ObjectReadObjectNotExists:
		return json.Marshal(ObjectReadObjectNotExists{Details: obj.Details, Status: obj.Status})
	case ObjectReadObjectDeleted:
		return json.Marshal(ObjectReadObjectDeleted{Details: obj.Details, Status: obj.Status})
	case ObjectReadVersionNotFound:
		return json.Marshal(ObjectReadVersionNotFound{Details: obj.Details, Status: obj.Status})
	case ObjectReadVersionTooHigh:
		return json.Marshal(ObjectReadVersionTooHigh{Details: obj.Details, Status: obj.Status})
	default:
		return nil, errors.New("unknown ObjectRead type")
	}
}
