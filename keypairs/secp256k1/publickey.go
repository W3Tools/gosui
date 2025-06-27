package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/cryptography"
	"github.com/btcsuite/btcd/btcec/v2"
)

var (
	_ cryptography.PublicKey = (*PublicKey)(nil)
)

// PublicKey defines a Secp256k1 public key used for signing transactions.
type PublicKey struct {
	data []byte
	cryptography.BasePublicKey
}

// NewPublicKey creates a new Secp256k1 public key from a base64 encoded string or byte slice.
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
	if len(publicKey.data) != cryptography.Secp256k1PublicKeySize {
		return nil, fmt.Errorf("invalid public key input, expected %d bytes, got %d", cryptography.Secp256k1PublicKeySize, len(publicKey.data))
	}
	publicKey.SetSelf(publicKey)
	return
}

// Equals checks if the provided public key is equal to this Secp256k1 public key.
func (k *PublicKey) Equals(publicKey cryptography.PublicKey) bool {
	return k.BasePublicKey.Equals(publicKey)
}

// ToRawBytes returns the raw bytes of the Secp256k1 public key.
func (k *PublicKey) ToRawBytes() []byte {
	return k.data
}

// Flag returns the flag for the Secp256k1 public key.
func (k *PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.Secp256k1Scheme]
}

// Verify checks whether the signature is valid for the provided message using the Secp256k1 public key.
func (k *PublicKey) Verify(message []byte, signature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(signature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.Secp256k1Scheme {
		return false, fmt.Errorf("invalid signature scheme")
	}

	if !bytes.Equal(k.ToRawBytes(), parsed.PubKey) {
		return false, fmt.Errorf("signature does not match public key")
	}

	// Parse the signature into r and s components
	r := new(big.Int).SetBytes(parsed.Signature[:32])
	s := new(big.Int).SetBytes(parsed.Signature[32:])

	pubKey, err := btcec.ParsePubKey(k.ToRawBytes())
	if err != nil {
		return false, err
	}

	// Verify the signature
	hash := sha256.Sum256(message)
	verified := ecdsa.Verify(pubKey.ToECDSA(), hash[:], r, s)

	return verified, nil
}
