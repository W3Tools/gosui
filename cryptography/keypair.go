package cryptography

import (
	"fmt"

	"github.com/W3Tools/gosui/b64"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

const (
	// PrivateKeySize is the size of a Sui private key in bytes.
	PrivateKeySize = 32
	// SuiPrivateKeyPrefix is the Bech32 prefix for Sui private keys.
	SuiPrivateKeyPrefix = "suiprivkey"
)

// ParsedKeypair defines a structure that holds the signature scheme and the secret key in bytes.
type ParsedKeypair struct {
	Scheme    SignatureScheme
	SecretKey []byte
}

// SignatureWithBytes defines a structure that holds the base64 encoded bytes of the signed message and the serialized signature.
type SignatureWithBytes struct {
	Bytes     string
	Signature SerializedSignature
}

// Signer is an interface that defines the methods for signing messages and retrieving public keys.
type Signer interface {
	Sign(bs []byte) ([]byte, error)
	SignWithIntent(bs []byte, intent IntentScope) (*SignatureWithBytes, error)
	SignTransactionBlock(bs []byte) (*SignatureWithBytes, error)
	SignPersonalMessage(bs []byte) (*SignatureWithBytes, error)
	ToSuiAddress() string
	SignData(data []byte) ([]byte, error)
	GetKeyScheme() SignatureScheme
	GetPublicKey() (PublicKey, error)
}

// Keypair is an interface that extends the Signer interface and adds a method to retrieve the secret key.
type Keypair interface {
	Signer
	GetSecretKey() (string, error)
}

// BaseSigner defines a struct that implements the Signer interface.
type BaseSigner struct {
	self Signer
}

// SignWithIntent signs a message with a specific intent scope. By combining the message bytes with the intent before hashing and signing,
// it ensures that a signed message is tied to a specific purpose and domain separator is provided
func (signer *BaseSigner) SignWithIntent(bs []byte, intent IntentScope) (*SignatureWithBytes, error) {
	intentMessage := MessageWithIntent(intent, bs)
	digest := blake2b.Sum256(intentMessage)

	publicKey, err := signer.self.GetPublicKey()
	if err != nil {
		return nil, err
	}

	sign, err := signer.self.Sign(digest[:])
	if err != nil {
		return nil, err
	}

	signature, err := ToSerializedSignature(SerializeSignatureInput{
		SignatureScheme: signer.self.GetKeyScheme(),
		Signature:       sign,
		PublicKey:       publicKey,
	})
	if err != nil {
		return nil, err
	}

	return &SignatureWithBytes{Bytes: b64.ToBase64(bs), Signature: signature}, nil
}

// SignTransactionBlock signs provided transaction block by calling `signWithIntent()` with a `TransactionData` provided as intent scope
func (signer *BaseSigner) SignTransactionBlock(bs []byte) (*SignatureWithBytes, error) {
	return signer.SignWithIntent(bs, TransactionData)
}

// SignPersonalMessage signs provided personal message by calling `signWithIntent()` with a `PersonalMessage` provided as intent scope
func (signer *BaseSigner) SignPersonalMessage(bs []byte) (*SignatureWithBytes, error) {
	msg, err := bcs.Marshal(bs)
	if err != nil {
		return nil, err
	}

	return signer.SignWithIntent(msg, PersonalMessage)
}

// ToSuiAddress converts the public key of the signer to a Sui address.
func (signer *BaseSigner) ToSuiAddress() string {
	publicKey, _ := signer.self.GetPublicKey()
	return publicKey.ToSuiAddress()
}

// abstract
// func (signer *BaseSigner) Sign(bs []byte) []byte         { return nil }
// func (signer *BaseSigner) SignData(data []byte) []byte   { return nil }
// func (signer *BaseSigner) GetKeyScheme() SignatureScheme { return "" }
// func (signer *BaseSigner) GetPublicKey() PublicKey       { return nil }

// BaseKeypair defines a struct that implements the Keypair interface.
type BaseKeypair struct {
	BaseSigner
}

// SetSelf sets the signer for the BaseKeypair.
func (base *BaseKeypair) SetSelf(signer Signer) {
	base.self = signer
}

// DecodeSuiPrivateKey decodes a Bech32 encoded Sui private key string into a ParsedKeypair object.
// This returns a ParsedKeypair object by validating the 33-byte Bech32 encoded string starting with `suiprivkey`,
// and parses out the signature scheme and the private key in bytes.
func DecodeSuiPrivateKey(value string) (*ParsedKeypair, error) {
	prefix, words, err := bech32.Decode(value)
	if err != nil {
		return nil, err
	}
	if prefix != SuiPrivateKeyPrefix {
		return nil, fmt.Errorf("invalid private key prefix")
	}
	extendedSecretKey, err := bech32.ConvertBits(words, 5, 8, false)
	if err != nil {
		return nil, err
	}
	secretKey := extendedSecretKey[1:]
	signatureScheme := SignatureFlagToScheme[extendedSecretKey[0]]
	return &ParsedKeypair{Scheme: signatureScheme, SecretKey: secretKey}, nil
}

// EncodeSuiPrivateKey returns a Bech32 encoded string starting with `suiprivkey`,
// encoding 33-byte `flag || bytes` for the given 32-byte private key and its signature scheme.
func EncodeSuiPrivateKey(bytes []byte, scheme string) (string, error) {
	if len(bytes) != PrivateKeySize {
		return "", fmt.Errorf("invalid bytes length")
	}
	flag := SignatureSchemeToFlag[scheme]
	privKeyBytes := append([]byte{flag}, bytes...)
	words, err := bech32.ConvertBits(privKeyBytes, 8, 5, true)
	if err != nil {
		return "", err
	}
	encoded, err := bech32.Encode(SuiPrivateKeyPrefix, words)
	if err != nil {
		return "", err
	}
	return encoded, nil
}
