package types

import (
	"encoding/json"
	"errors"
)

// Signature is an interface that defines a Sui signature type.
type Signature interface {
	isSignature()
}

// SignatureEd25519 defines an Ed25519 Sui signature.
type SignatureEd25519 struct {
	Ed25519SuiSignature string `json:"Ed25519SuiSignature"`
}

// SignatureSecp256k1 defines a Secp256k1 Sui signature.
type SignatureSecp256k1 struct {
	Secp256k1SuiSignature string `json:"Secp256k1SuiSignature"`
}

// SignatureSecp256r1 defines a Secp256r1 Sui signature.
type SignatureSecp256r1 struct {
	Secp256r1SuiSignature string `json:"Secp256r1SuiSignature"`
}

// isSignature implements the Signature interface for SignatureEd25519.
func (SignatureEd25519) isSignature() {}

// isSignature implements the Signature interface for SignatureSecp256k1.
func (SignatureSecp256k1) isSignature() {}

// isSignature implements the Signature interface for SignatureSecp256r1.
func (SignatureSecp256r1) isSignature() {}

// SignatureWrapper defines a wrapper for Signature to support custom JSON marshaling and unmarshaling.
type SignatureWrapper struct {
	Signature
}

// UnmarshalJSON decodes JSON data into a SignatureWrapper.
func (w *SignatureWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["Ed25519SuiSignature"] != nil:
		var s SignatureEd25519
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	case obj["Secp256k1SuiSignature"] != nil:
		var s SignatureSecp256k1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	case obj["Secp256r1SuiSignature"] != nil:
		var s SignatureSecp256r1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.Signature = s
	default:
		return errors.New("unknown Signature type")
	}

	return nil
}

// MarshalJSON encodes a SignatureWrapper into JSON.
func (w *SignatureWrapper) MarshalJSON() ([]byte, error) {
	switch s := w.Signature.(type) {
	case SignatureEd25519:
		return json.Marshal(SignatureEd25519{
			Ed25519SuiSignature: s.Ed25519SuiSignature,
		})
	case SignatureSecp256k1:
		return json.Marshal(SignatureSecp256k1{
			Secp256k1SuiSignature: s.Secp256k1SuiSignature,
		})
	case SignatureSecp256r1:
		return json.Marshal(SignatureSecp256r1{
			Secp256r1SuiSignature: s.Secp256r1SuiSignature,
		})
	default:
		return nil, errors.New("unknown Signature type")
	}
}

// CompressedSignature is an interface that defines a compressed signature type.
type CompressedSignature interface {
	isCompressedSignature()
}

// CompressedSignatureEd25519 defines a compressed Ed25519 signature.
type CompressedSignatureEd25519 struct {
	Ed25519 string `json:"Ed25519"`
}

// CompressedSignatureSecp256k1 defines a compressed Secp256k1 signature.
type CompressedSignatureSecp256k1 struct {
	Secp256k1 string `json:"Secp256k1"`
}

// CompressedSignatureSecp256r1 defines a compressed Secp256r1 signature.
type CompressedSignatureSecp256r1 struct {
	Secp256r1 string `json:"Secp256r1"`
}

// isCompressedSignature implements the CompressedSignature interface for CompressedSignatureEd25519.
func (CompressedSignatureEd25519) isCompressedSignature() {}

// isCompressedSignature implements the CompressedSignature interface for CompressedSignatureSecp256k1.
func (CompressedSignatureSecp256k1) isCompressedSignature() {}

// isCompressedSignature implements the CompressedSignature interface for CompressedSignatureSecp256r1.
func (CompressedSignatureSecp256r1) isCompressedSignature() {}

// CompressedSignatureWrapper defines a wrapper for CompressedSignature to support custom JSON marshaling and unmarshaling.
type CompressedSignatureWrapper struct {
	CompressedSignature
}

// UnmarshalJSON decodes JSON data into a CompressedSignatureWrapper.
func (w *CompressedSignatureWrapper) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	switch {
	case obj["Ed25519"] != nil:
		var s CompressedSignatureEd25519
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	case obj["Secp256k1"] != nil:
		var s CompressedSignatureSecp256k1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	case obj["Secp256r1"] != nil:
		var s CompressedSignatureSecp256r1
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		w.CompressedSignature = s
	default:
		return errors.New("unknown CompressedSignature type")
	}

	return nil
}

// MarshalJSON encodes a CompressedSignatureWrapper into JSON.
func (w *CompressedSignatureWrapper) MarshalJSON() ([]byte, error) {
	switch s := w.CompressedSignature.(type) {
	case CompressedSignatureEd25519:
		return json.Marshal(CompressedSignatureEd25519{
			Ed25519: s.Ed25519,
		})
	case CompressedSignatureSecp256k1:
		return json.Marshal(CompressedSignatureSecp256k1{
			Secp256k1: s.Secp256k1,
		})
	case CompressedSignatureSecp256r1:
		return json.Marshal(CompressedSignatureSecp256r1{
			Secp256r1: s.Secp256r1,
		})
	default:
		return nil, errors.New("unknown CompressedSignature type")
	}
}
