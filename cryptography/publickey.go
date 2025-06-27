package cryptography

import (
	"bytes"
	"encoding/hex"
	"reflect"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/utils"
	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

// PublicKey is an interface that defines the methods for a public key in the Sui cryptography system.
type PublicKey interface {
	Equals(publicKey PublicKey) bool
	ToBase64() string
	ToSuiPublicKey() string
	VerifyWithIntent(bs []byte, signature SerializedSignature, intent IntentScope) (bool, error)
	VerifyPersonalMessage(message []byte, signature SerializedSignature) (bool, error)
	VerifyTransactionBlock(transactionBlock []byte, signature SerializedSignature) (bool, error)
	ToSuiBytes() []byte
	ToSuiAddress() string
	ToRawBytes() []byte
	Flag() uint8
	Verify(message []byte, signature SerializedSignature) (bool, error)
}

// BasePublicKey defines a struct that implements the PublicKey interface.
type BasePublicKey struct {
	self PublicKey
}

// Equals checks if the provided public key is equal to the base public key.
func (base *BasePublicKey) Equals(publicKey PublicKey) bool {
	return reflect.DeepEqual(base.self.ToRawBytes(), publicKey.ToRawBytes())
}

// ToBase64 returns the base-64 representation of the public key.
func (base *BasePublicKey) ToBase64() string {
	return b64.ToBase64(base.self.ToRawBytes())
}

// ToSuiPublicKey returns the Sui representation of the public key encoded in base-64.
// A Sui public key is formed by the concatenation of the scheme flag with the raw bytes of the public key
func (base *BasePublicKey) ToSuiPublicKey() string {
	suiBytes := base.ToSuiBytes()
	return b64.ToBase64(suiBytes)
}

// VerifyWithIntent verifies that the signature is valid for the provided message bytes and intent.
func (base *BasePublicKey) VerifyWithIntent(bs []byte, signature SerializedSignature, intent IntentScope) (bool, error) {
	intentMessage := MessageWithIntent(intent, bs)
	digest := blake2b.Sum256(intentMessage)
	return base.self.Verify(digest[:], signature)
}

// VerifyPersonalMessage verifies that the signature is valid for the provided personal message.
func (base *BasePublicKey) VerifyPersonalMessage(message []byte, signature SerializedSignature) (bool, error) {
	msg, err := bcs.Marshal(message)
	if err != nil {
		return false, err
	}
	return base.VerifyWithIntent(msg, signature, PersonalMessage)
}

// VerifyTransactionBlock verifies that the signature is valid for the provided transaction block.
func (base *BasePublicKey) VerifyTransactionBlock(transactionBlock []byte, signature SerializedSignature) (bool, error) {
	return base.VerifyWithIntent(transactionBlock, signature, TransactionData)
}

// ToSuiBytes returns the bytes representation of the public key prefixed with the signature scheme flag.
func (base *BasePublicKey) ToSuiBytes() []byte {
	rawBytes := base.self.ToRawBytes()
	suiBytes := new(bytes.Buffer)
	suiBytes.Write([]byte{base.self.Flag()})
	suiBytes.Write(rawBytes)

	return suiBytes.Bytes()
}

// ToSuiAddress returns the Sui address associated with this public key.
func (base *BasePublicKey) ToSuiAddress() string {
	digest := blake2b.Sum256(base.ToSuiBytes())
	return utils.NormalizeSuiAddress(hex.EncodeToString(digest[:])[:utils.SuiAddressLength*2])
}

// abstract: Return the byte array representation of the public key
// func (base *BasePublicKey) ToRawBytes() []byte

// abstract: Return signature scheme flag of the public key
// func (base *BasePublicKey) Flag() uint8

// // abstract: Verifies that the signature is valid for for the provided message
//
//	func (base *BasePublicKey) Verify(message []byte, signature SerializedSignature) (bool, error) {
//		return false, nil
//	}

// SetSelf sets the self reference for the BasePublicKey.
func (base *BasePublicKey) SetSelf(pubkey PublicKey) {
	base.self = pubkey
}
