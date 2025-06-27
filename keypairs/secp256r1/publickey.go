package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/cryptography"
)

var (
	_ cryptography.PublicKey = (*PublicKey)(nil)
)

// PublicKey defines a Secp256r1 public key used for verifying signatures.
type PublicKey struct {
	data []byte
	cryptography.BasePublicKey
}

// NewPublicKey creates a new Secp256r1 public key from a base64 encoded string or byte slice.
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
	if len(publicKey.data) != cryptography.Secp256r1PublicKeySize {
		return nil, fmt.Errorf("invalid public key input, expected %d bytes, got %d", cryptography.Secp256r1PublicKeySize, len(publicKey.data))
	}
	publicKey.SetSelf(publicKey)
	return
}

// Equals checks if the provided public key is equal to this Secp256r1 public key.
func (k *PublicKey) Equals(publicKey cryptography.PublicKey) bool {
	return k.BasePublicKey.Equals(publicKey)
}

// ToRawBytes returns the byte array representation of the Secp256r1 public key.
func (k *PublicKey) ToRawBytes() []byte {
	return k.data
}

// Flag returns the Sui address flag associated with this Secp256r1 public key.
func (k *PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.Secp256r1Scheme]
}

// Verify checks whether the signature is valid for the provided message using the Secp256r1 public key.
func (k *PublicKey) Verify(message []byte, signature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(signature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.Secp256r1Scheme {
		return false, fmt.Errorf("invalid signature scheme")
	}

	if !reflect.DeepEqual(k.ToRawBytes(), parsed.PubKey) {
		return false, fmt.Errorf("signature does not match public key")
	}

	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), k.ToRawBytes())
	if x == nil || y == nil {
		return false, errors.New("error unmarshaling public key")
	}

	curve := elliptic.P256()
	pubKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}

	// Parse the signature into r and s components
	r := new(big.Int).SetBytes(parsed.Signature[:32])
	s := new(big.Int).SetBytes(parsed.Signature[32:])

	// Verify the signature
	hash := sha256.Sum256(message)
	valid := ecdsa.Verify(&pubKey, hash[:], r, s)
	return valid, nil
}
