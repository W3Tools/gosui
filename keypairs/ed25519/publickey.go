package ed25519

import (
	"crypto/ed25519"
	"fmt"
	"reflect"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/cryptography"
)

var (
	_ cryptography.PublicKey = (*PublicKey)(nil)
)

// PublicKey defines an Ed25519 public key.
type PublicKey struct {
	data []byte
	cryptography.BasePublicKey
}

// NewPublicKey creates a new Ed25519 public key from a base64 encoded string or byte slice.
func NewPublicKey[T string | []byte](value T) (publicKey *PublicKey, err error) {
	publicKey = new(PublicKey)
	switch v := any(value).(type) {
	case string:
		publicKey.data, err = b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		publicKey.data = v
	}

	if len(publicKey.data) != cryptography.Ed25519PublicKeySize {
		return nil, fmt.Errorf("invalid public key input. expected %v bytes, got %v", cryptography.Ed25519PublicKeySize, len(publicKey.data))
	}
	publicKey.SetSelf(publicKey)
	return
}

// ToRawBytes returns the raw bytes of the Ed25519 public key.
func (key *PublicKey) ToRawBytes() []byte {
	return key.data
}

// Flag returns the flag for the Ed25519 public key scheme.
func (key *PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.Ed25519Scheme]
}

// Verify checks if the signature is valid for the provided message using the Ed25519 public key.
func (key *PublicKey) Verify(message []byte, signature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(signature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.Ed25519Scheme {
		return false, fmt.Errorf("invalid signature scheme")
	}

	if !reflect.DeepEqual(key.ToRawBytes(), parsed.PubKey) {
		return false, fmt.Errorf("signature does not match public key")
	}

	return ed25519.Verify(parsed.PubKey, message, parsed.Signature), nil
}
