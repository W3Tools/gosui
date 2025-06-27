package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/W3Tools/gosui/cryptography"
)

var (
	_ cryptography.Keypair = (*Keypair)(nil)
)

// DefaultDerivationPath is the default derivation path for Ed25519 keypairs.
const DefaultDerivationPath = "m/44'/784'/0'/0'/0'"

// KeypairData defines the public key (32 bytes) and secret key (64 bytes) for an Ed25519 keypair.
type KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// Keypair defines an Ed25519 keypair used for signing transactions.
type Keypair struct {
	keypair *KeypairData
	cryptography.BaseKeypair
}

// NewKeypair creates or generates a new Ed25519 keypair instance.
func NewKeypair(keypair *KeypairData) (*Keypair, error) {
	if keypair == nil {
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		keypair = &KeypairData{PublicKey: publicKey, SecretKey: privateKey}
	}

	kp := &Keypair{keypair: keypair}
	kp.SetSelf(kp)
	return kp, nil
}

// GetKeyScheme returns the signature scheme of the Ed25519 keypair.
func (k *Keypair) GetKeyScheme() cryptography.SignatureScheme {
	return cryptography.Ed25519Scheme
}

// GetPublicKey returns the public key of the Ed25519 keypair.
func (k *Keypair) GetPublicKey() (cryptography.PublicKey, error) {
	return NewPublicKey(k.keypair.PublicKey)
}

// GetSecretKey returns the secret key of the Ed25519 keypair as a base64 encoded string.
func (k *Keypair) GetSecretKey() (string, error) {
	return cryptography.EncodeSuiPrivateKey(k.keypair.SecretKey[:cryptography.PrivateKeySize], k.GetKeyScheme())
}

// Sign returns the signature for the provided data using Ed25519.
func (k *Keypair) Sign(data []byte) ([]byte, error) {
	return k.SignData(data)
}

// SignData signs the provided data using the Ed25519 keypair's secret key.
func (k *Keypair) SignData(data []byte) ([]byte, error) {
	return ed25519.Sign(k.keypair.SecretKey, data), nil
}

// GenerateKeypair creates a new Ed25519 keypair with a random secret key.
func GenerateKeypair() (*Keypair, error) {
	return NewKeypair(nil)
}

// FromSecretKey creates an Ed25519 keypair from a raw secret key byte array (seed).
// This is NOT the private scalar which is result of hashing and bit clamping of the raw secret key.
func FromSecretKey(secretKey []byte, skipValidation bool) (*Keypair, error) {
	secretKeyLength := len(secretKey)
	if secretKeyLength != cryptography.PrivateKeySize {
		return nil, fmt.Errorf("wrong secretKey size. expected %d bytes, got %d", cryptography.PrivateKeySize, secretKeyLength)
	}

	privateKey := ed25519.NewKeyFromSeed(secretKey)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	if !skipValidation {
		signData := []byte("sui validation")
		signature := ed25519.Sign(privateKey, signData)
		if !ed25519.Verify(publicKey, signData, signature) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}

	return NewKeypair(&KeypairData{PublicKey: publicKey, SecretKey: privateKey})
}

// DeriveKeypair derives an Ed25519 keypair from a mnemonic phrase and a derivation path.
// The mnemonics must be normalized and validated against the english wordlist.
// If path is none, it will default to m/44'/784'/0'/0'/0'
// Otherwise the path must be compliant to SLIP-0010 in form m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
func DeriveKeypair(mnemonics, path string) (*Keypair, error) {
	if path == "" {
		path = DefaultDerivationPath
	}
	if !cryptography.IsValidHardenedPath(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}

	seedHex, err := cryptography.MnemonicToSeedHex(mnemonics)
	if err != nil {
		return nil, err
	}

	key, err := cryptography.DerivePath(path, seedHex)
	if err != nil {
		return nil, err
	}

	return FromSecretKey(key.Key, false)
}

// DeriveKeypairFromSeed derives an Ed25519 keypair from a seed hex string and a derivation path.
// If path is none, it will default to m/44'/784'/0'/0'/0'
// Otherwise the path must be compliant to SLIP-0010 in form m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
func DeriveKeypairFromSeed(seedHex, path string) (*Keypair, error) {
	if path == "" {
		path = DefaultDerivationPath
	}

	if !cryptography.IsValidHardenedPath(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}

	key, err := cryptography.DerivePath(path, seedHex)
	if err != nil {
		return nil, err
	}
	return FromSecretKey(key.Key, false)
}

// PublicKey returns the public key of the Ed25519 keypair.
func (k *Keypair) PublicKey() []byte {
	return k.keypair.PublicKey
}

// SecretKey returns the secret key of the Ed25519 keypair.
func (k *Keypair) SecretKey() []byte {
	return k.keypair.SecretKey
}
