package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/W3Tools/gosui/cryptography"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/blake2b"
)

var (
	_ cryptography.Keypair = (*Keypair)(nil)
)

// DefaultDerivationPath is the default derivation path for Secp256r1 keypairs.
const DefaultDerivationPath = "m/74'/784'/0'/0/0"

// KeypairData defines the public and secret key data for a Secp256r1 keypair.
type KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// Keypair defines a Secp256r1 keypair used for signing transactions.
type Keypair struct {
	keypair *KeypairData
	cryptography.BaseKeypair
}

// NewKeypair creates a new Secp256r1 keypair.
func NewKeypair(keypair *KeypairData) (*Keypair, error) {
	if keypair == nil {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		publicKey := elliptic.MarshalCompressed(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
		keypair = &KeypairData{PublicKey: publicKey, SecretKey: privateKey.D.Bytes()}
	}

	kp := &Keypair{keypair: keypair}
	kp.SetSelf(kp)
	return kp, nil
}

// GetKeyScheme returns the signature scheme of the keypair.
func (k *Keypair) GetKeyScheme() cryptography.SignatureScheme {
	return cryptography.Secp256r1Scheme
}

// GetPublicKey returns the public key of the keypair.
func (k *Keypair) GetPublicKey() (cryptography.PublicKey, error) {
	return NewPublicKey(k.keypair.PublicKey)
}

// GetSecretKey returns the secret key of the keypair as a hex-encoded string.
func (k *Keypair) GetSecretKey() (string, error) {
	return cryptography.EncodeSuiPrivateKey(k.keypair.SecretKey, k.GetKeyScheme())
}

// Sign signs the provided data and returns the signature.
func (k *Keypair) Sign(data []byte) ([]byte, error) {
	return k.SignData(data)
}

// SignData signs the provided data using the keypair and returns the signature.
func (k *Keypair) SignData(data []byte) ([]byte, error) {
	hexMessage := sha256.Sum256(data)

	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = new(big.Int).SetBytes(k.keypair.SecretKey)
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKey.D.Bytes())
	r, s, err := deterministicSign(privKey, hexMessage[:])
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

// GenerateKeypair generates a new Secp256r1 keypair with a random secret key.
func GenerateKeypair() (*Keypair, error) {
	return NewKeypair(nil)
}

// FromSecretKey creates a keypair from a raw secret key byte array.
// This method should only be used to recreate a keypair from a previously generated secret key.
// Generating keypairs from a random seed should be done with the {@link Keypair.fromSeed} method.
func FromSecretKey(secretKey []byte, skipValidation bool) (*Keypair, error) {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(secretKey)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())
	publicKey := elliptic.MarshalCompressed(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)

	if !skipValidation {
		signData := []byte("sui validation")
		hash := blake2b.Sum256(signData)
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
		if err != nil {
			return nil, err
		}

		if !ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}
	return NewKeypair(&KeypairData{PublicKey: publicKey, SecretKey: secretKey})
}

// FromSeed generates a Secp256r1 keypair from a 32-byte seed.
func FromSeed(seed []byte) (*Keypair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), strings.NewReader(string(seed)))
	if err != nil {
		return nil, err
	}
	publicKey := elliptic.MarshalCompressed(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	secretKey := privateKey.D.Bytes()

	return NewKeypair(&KeypairData{PublicKey: publicKey, SecretKey: secretKey})
}

// DeriveKeypair derives a Secp256r1 keypair from mnemonics and a specified derivation path.
// If path is none, it will default to m/74'/784'/0'/0/0
// Otherwise the path must be compliant to BIP-32 in form m/74'/784'/{account_index}'/{change_index}/{address_index}.
func DeriveKeypair(mnemonics string, path string) (*Keypair, error) {
	if path == "" {
		path = DefaultDerivationPath
	}

	if !cryptography.IsValidBIP32Path(path) {
		return nil, fmt.Errorf("invalid derivation path")
	}
	seed, err := cryptography.MnemonicToSeed(mnemonics)
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	key, err := deriveChildKeyFromPath(masterKey, path)
	if err != nil {
		return nil, err
	}
	return FromSecretKey(key.Key, false)
}

func deriveChildKeyFromPath(masterKey *bip32.Key, path string) (*bip32.Key, error) {
	segments := strings.Split(path, "/")[1:]
	key := masterKey

	for _, segment := range segments {
		if strings.HasSuffix(segment, "'") {
			index, err := strconv.Atoi(strings.TrimSuffix(segment, "'"))
			if err != nil {
				return nil, err
			}
			key, err = key.NewChildKey(bip32.FirstHardenedChild + uint32(index))
			if err != nil {
				return nil, err
			}
		} else {
			index, err := strconv.Atoi(segment)
			if err != nil {
				return nil, err
			}
			key, err = key.NewChildKey(uint32(index))
			if err != nil {
				return nil, err
			}
		}
	}

	return key, nil
}

// PublicKey returns the public key bytes of the keypair.
func (k *Keypair) PublicKey() []byte {
	return k.keypair.PublicKey
}

// SecretKey returns the secret key bytes of the keypair.
func (k *Keypair) SecretKey() []byte {
	return k.keypair.SecretKey
}
