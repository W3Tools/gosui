package types

import (
	"encoding/json"
	"errors"
)

// ObjectResponseError is an interface that defines the methods for handling different types of object response errors in the Sui blockchain.
type ObjectResponseError interface {
	isObjectResponseError()
}

// ObjectResponseNotExistsError defines an error that occurs when an object does not exist in the Sui blockchain.
type ObjectResponseNotExistsError struct {
	Code     string `json:"code"`
	ObjectID string `json:"object_id"`
}

// ObjectResponseDynamicFieldNotFoundError defines an error that occurs when a dynamic field is not found in the Sui blockchain.
type ObjectResponseDynamicFieldNotFoundError struct {
	Code           string `json:"code"`
	ParentObjectID string `json:"parent_object_id"`
}

// ObjectResponseDeletedError defines an error that occurs when an object has been deleted in the Sui blockchain.
type ObjectResponseDeletedError struct {
	Code     string `json:"code"`
	Digest   string `json:"digest"`
	ObjectID string `json:"object_id"`
	Version  uint64 `json:"version"`
}

// ObjectResponseUnknownError defines an error that occurs when an unknown error happens in the Sui blockchain.
type ObjectResponseUnknownError struct {
	Code string `json:"code"`
}

// ObjectResponseDisplayErrorError defines an error that occurs when a display error happens in the Sui blockchain.
type ObjectResponseDisplayErrorError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (ObjectResponseNotExistsError) isObjectResponseError()            {}
func (ObjectResponseDynamicFieldNotFoundError) isObjectResponseError() {}
func (ObjectResponseDeletedError) isObjectResponseError()              {}
func (ObjectResponseUnknownError) isObjectResponseError()              {}
func (ObjectResponseDisplayErrorError) isObjectResponseError()         {}

// ObjectResponseErrorWrapper is a wrapper for ObjectResponseError that allows unmarshalling from JSON.
type ObjectResponseErrorWrapper struct {
	ObjectResponseError
}

// UnmarshalJSON custom unmarshaller for ObjectResponseErrorWrapper
func (w *ObjectResponseErrorWrapper) UnmarshalJSON(data []byte) error {
	type ErrorCode struct {
		Code string `json:"code"`
	}

	var errorCode ErrorCode
	if err := json.Unmarshal(data, &errorCode); err != nil {
		return err
	}

	switch errorCode.Code {
	case "notExists":
		var e ObjectResponseNotExistsError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "deleted":
		var e ObjectResponseDeletedError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "dynamicFieldNotFound":
		var e ObjectResponseDynamicFieldNotFoundError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "unknown":
		var e ObjectResponseUnknownError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	case "displayError":
		var e ObjectResponseDisplayErrorError
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}

		w.ObjectResponseError = e
		return nil
	default:
		return errors.New("unknown ObjectResponseError type")
	}
}

// MarshalJSON custom marshaller for ObjectResponseErrorWrapper
func (w *ObjectResponseErrorWrapper) MarshalJSON() ([]byte, error) {
	switch e := w.ObjectResponseError.(type) {
	case ObjectResponseNotExistsError:
		return json.Marshal(ObjectResponseNotExistsError{Code: e.Code, ObjectID: e.ObjectID})
	case ObjectResponseDynamicFieldNotFoundError:
		return json.Marshal(ObjectResponseDynamicFieldNotFoundError{Code: e.Code, ParentObjectID: e.ParentObjectID})
	case ObjectResponseDeletedError:
		return json.Marshal(ObjectResponseDeletedError{Code: e.Code, Digest: e.Digest, ObjectID: e.ObjectID, Version: e.Version})
	case ObjectResponseUnknownError:
		return json.Marshal(ObjectResponseUnknownError{Code: e.Code})
	case ObjectResponseDisplayErrorError:
		return json.Marshal(ObjectResponseDisplayErrorError{Code: e.Code, Error: e.Error})
	default:
		return nil, errors.New("unknown ObjectResponseError type")
	}
}
