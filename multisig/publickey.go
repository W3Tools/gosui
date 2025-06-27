package multisig

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/W3Tools/gosui/b64"
	"github.com/W3Tools/gosui/cryptography"
	"github.com/W3Tools/gosui/utils"
	"github.com/W3Tools/gosui/verify"
	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

// MaxSignerInMultisig is the maximum number of signers allowed in a multisig.
// MinSignerInMultisig is the minimum number of signers required in a multisig.
const (
	MaxSignerInMultisig = 10
	MinSignerInMultisig = 1
)

var (
	_ cryptography.PublicKey = (*PublicKey)(nil)
)

// StringPubKeyEnumWeightPair defines a public key and weight pair for multisig.
type StringPubKeyEnumWeightPair struct {
	PubKey string `json:"pubKey"`
	Weight uint8  `json:"weight"`
}

// ParsedPartialMultiSigSignature defines a parsed partial multisig signature.
type ParsedPartialMultiSigSignature struct {
	SignatureScheme cryptography.SignatureScheme
	Signature       []byte
	PublicKey       cryptography.PublicKey
	Weight          uint8
}

// PublicKey defines a multisig public key.
type PublicKey struct {
	rawBytes          []byte
	multisigPublicKey cryptography.MultiSigPublicKeyStruct
	publicKeys        []cryptography.MultiSigPublicKeyPair
	cryptography.BasePublicKey
}

// PublicKeyWeightPair defines a public key and weight pair for multisig construction.
type PublicKeyWeightPair struct {
	PublicKey cryptography.PublicKey
	Weight    uint8
}

// NewPublicKey creates a new multisig PublicKey from various input types.
func NewPublicKey[T string | []byte | cryptography.MultiSigPublicKeyStruct](value T) (multisig *PublicKey, err error) {
	multisig = new(PublicKey)
	switch v := any(value).(type) {
	case string:
		rawBytes, err := b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
		multisig.rawBytes = rawBytes

		multisigPublicKeyStruct := new(cryptography.MultiSigPublicKeyStruct)
		_, err = bcs.Unmarshal(multisig.rawBytes, &multisigPublicKeyStruct)
		if err != nil {
			return nil, err
		}
		multisig.multisigPublicKey = *multisigPublicKeyStruct
	case []byte:
		multisig.rawBytes = v

		multisigPublicKeyStruct := new(cryptography.MultiSigPublicKeyStruct)
		_, err = bcs.Unmarshal(multisig.rawBytes, &multisigPublicKeyStruct)
		if err != nil {
			return nil, err
		}
		multisig.multisigPublicKey = *multisigPublicKeyStruct
	case cryptography.MultiSigPublicKeyStruct:
		multisig.multisigPublicKey = v

		rawBytes, err := bcs.Marshal(multisig.multisigPublicKey)
		if err != nil {
			return nil, err
		}
		multisig.rawBytes = rawBytes
	}
	if multisig.multisigPublicKey.Threshold < 1 {
		return nil, fmt.Errorf("invalid threshold")
	}

	seenPublicKeys := make(map[string]bool)

	for _, v := range multisig.multisigPublicKey.PubKeyMap {
		publicKeyString := string(v.PubKey[:])
		if ok := seenPublicKeys[publicKeyString]; ok {
			return nil, fmt.Errorf("multisig does not support duplicate public keys")
		}
		seenPublicKeys[publicKeyString] = true

		if v.Weight < 1 {
			return nil, fmt.Errorf("invalid weight")
		}

		pubKey, err := verify.PublicKeyFromRawBytes(cryptography.SignatureFlagToScheme[v.PubKey[0]], v.PubKey[1:])
		if err != nil {
			return nil, err
		}

		multisig.publicKeys = append(multisig.publicKeys, cryptography.MultiSigPublicKeyPair{PublicKey: pubKey, Weight: v.Weight})
	}

	var totalWeight uint16 = 0
	for _, pubkey := range multisig.publicKeys {
		totalWeight = totalWeight + uint16(pubkey.Weight)
	}

	if multisig.multisigPublicKey.Threshold > totalWeight {
		return nil, fmt.Errorf("unreachable threshold")
	}

	if len(multisig.publicKeys) > MaxSignerInMultisig {
		return nil, fmt.Errorf("max number of signers in a multisig is %d", MaxSignerInMultisig)
	}

	if len(multisig.publicKeys) < MinSignerInMultisig {
		return nil, fmt.Errorf("min number of signers in a multisig is %d", MinSignerInMultisig)
	}

	multisig.SetSelf(multisig)
	return
}

// FromPublicKeys creates a multisig PublicKey from a list of public keys and a threshold.
func (multisig *PublicKey) FromPublicKeys(publicKeys []PublicKeyWeightPair, threshold uint16) (*PublicKey, error) {
	pubkeys := make([]*cryptography.PubKeyEnumWeightPair, len(publicKeys))
	for i, v := range publicKeys {
		pubkeys[i] = &cryptography.PubKeyEnumWeightPair{PubKey: v.PublicKey.ToSuiBytes(), Weight: v.Weight}
	}

	return NewPublicKey(cryptography.MultiSigPublicKeyStruct{PubKeyMap: pubkeys, Threshold: threshold})
}

// Equals checks if the multisig public key equals another public key.
func (multisig *PublicKey) Equals(publicKey cryptography.PublicKey) bool {
	return multisig.BasePublicKey.Equals(publicKey)
}

// ToSuiAddress returns the Sui address associated with the multisig public key.
func (multisig *PublicKey) ToSuiAddress() string {
	tmp := new(bytes.Buffer)
	tmp.WriteByte(cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme])

	threshold, _ := bcs.Marshal(multisig.multisigPublicKey.Threshold)
	tmp.Write(threshold)

	for _, publicKey := range multisig.publicKeys {
		tmp.Write(publicKey.PublicKey.ToSuiBytes())
		tmp.WriteByte(publicKey.Weight)
	}

	sum256 := blake2b.Sum256(tmp.Bytes())
	return utils.NormalizeShortSuiAddress(hex.EncodeToString(sum256[:])[:64])
}

// ToRawBytes returns the byte array representation of the multisig public key.
func (multisig *PublicKey) ToRawBytes() []byte {
	return multisig.rawBytes
}

// GetPublicKeys returns the public keys in the multisig public key.
func (multisig *PublicKey) GetPublicKeys() []cryptography.MultiSigPublicKeyPair {
	return multisig.publicKeys
}

// GetThreshold returns the threshold of the multisig public key.
func (multisig *PublicKey) GetThreshold() uint16 {
	return multisig.multisigPublicKey.Threshold
}

// Flag returns the signature scheme flag for the multisig public key.
func (multisig *PublicKey) Flag() uint8 {
	return cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme]
}

// Verify checks if the signature is valid for the provided message.
func (multisig *PublicKey) Verify(message []byte, multisigSignature cryptography.SerializedSignature) (bool, error) {
	parsed, err := cryptography.ParseSerializedSignature(multisigSignature)
	if err != nil {
		return false, err
	}

	if parsed.SignatureScheme != cryptography.MultiSigScheme {
		return false, err
	}

	thisMultisig := parsed.Multisig

	bs1, err := bcs.Marshal(multisig.multisigPublicKey)
	if err != nil {
		return false, err
	}
	bs2, err := bcs.Marshal(thisMultisig.MultisigPubKey)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(bs1, bs2) {
		return false, err
	}

	var signatureWeight uint16 = 0
	partialParsedData, err := ParsePartialSignatures(thisMultisig)
	if err != nil {
		return false, err
	}

	for _, data := range partialParsedData {
		signature, err := cryptography.ToSerializedSignature(cryptography.SerializeSignatureInput{
			SignatureScheme: cryptography.SignatureFlagToScheme[data.PublicKey.Flag()],
			Signature:       data.Signature,
			PublicKey:       data.PublicKey,
		})
		if err != nil {
			return false, err
		}

		pass, err := data.PublicKey.Verify(message, signature)
		if err != nil {
			return false, err
		}
		if !pass {
			return false, nil
		}

		signatureWeight += uint16(data.Weight)
	}
	return (signatureWeight >= multisig.multisigPublicKey.Threshold), nil
}

// CombinePartialSignatures combines multiple partial signatures into a multisig signature.
func (multisig *PublicKey) CombinePartialSignatures(signatures []cryptography.SerializedSignature) (cryptography.SerializedSignature, error) {
	if len(signatures) > MaxSignerInMultisig {
		return "", fmt.Errorf("max number of signatures in a multisig is %d", MaxSignerInMultisig)
	}

	var bitmap uint16 = 0
	compressedSignatures := make([]cryptography.CompressedSignature, len(signatures))
	for i := 0; i < len(signatures); i++ {
		parsed, err := cryptography.ParseSerializedSignature(signatures[i])
		if err != nil {
			return "", err
		}

		if parsed.SignatureScheme == cryptography.MultiSigScheme {
			return "", fmt.Errorf("multisig is not supported inside MultiSig")
		}
		if parsed.SignatureScheme == cryptography.ZkLoginScheme {
			return "", fmt.Errorf("unimplemented %v", parsed.SignatureScheme)
		}

		tmp := new(bytes.Buffer)
		tmp.Write([]byte{cryptography.SignatureSchemeToFlag[parsed.SignatureScheme]})
		tmp.Write(parsed.Signature)
		compressedSignatures[i] = cryptography.CompressedSignature{Signature: [65]byte(tmp.Bytes())}

		var publicKeyIndex *int
		for j := 0; j < len(multisig.publicKeys); j++ {
			if bytes.Equal(parsed.PubKey, multisig.publicKeys[j].PublicKey.ToRawBytes()) {
				if bitmap&(1<<j) > 0 {
					return "", fmt.Errorf("received multiple signatures from the same public key")
				}
				publicKeyIndex = &j
				break
			}
		}

		if publicKeyIndex == nil {
			return "", fmt.Errorf("received signature from unknown public key")
		}
		bitmap |= 1 << *publicKeyIndex
	}

	m := cryptography.MultiSigStruct{
		Sigs:           compressedSignatures,
		Bitmap:         bitmap,
		MultisigPubKey: multisig.multisigPublicKey,
	}

	bs, err := bcs.Marshal(&m)
	if err != nil {
		return "", err
	}

	tmp := new(bytes.Buffer)
	tmp.Write([]byte{cryptography.SignatureSchemeToFlag[cryptography.MultiSigScheme]})
	tmp.Write(bs)

	return b64.ToBase64(tmp.Bytes()), nil
}

// ParsePartialSignatures parses a multisig structure into an array of individual signatures.
func ParsePartialSignatures(multisig *cryptography.MultiSigStruct) ([]ParsedPartialMultiSigSignature, error) {
	res := make([]ParsedPartialMultiSigSignature, len(multisig.Sigs))

	for i := 0; i < len(multisig.Sigs); i++ {
		signatureScheme := cryptography.SignatureFlagToScheme[multisig.Sigs[i].Signature[0]]
		signature := multisig.Sigs[i].Signature[1:]

		bitmapIndices, err := AsIndices(multisig.Bitmap)
		if err != nil {
			return nil, err
		}
		pkIndex := bitmapIndices[i]
		pair := multisig.MultisigPubKey.PubKeyMap[pkIndex]
		pkBytes := pair.PubKey[1:]

		if signatureScheme == cryptography.MultiSigScheme {
			return nil, fmt.Errorf("multisig is not supported inside MultiSig")
		}

		publicKey, err := verify.PublicKeyFromRawBytes(signatureScheme, pkBytes)
		if err != nil {
			return nil, err
		}

		res[i] = ParsedPartialMultiSigSignature{
			SignatureScheme: signatureScheme,
			Signature:       signature,
			PublicKey:       publicKey,
			Weight:          pair.Weight,
		}

	}
	return res, nil
}

// AsIndices returns the indices of set bits in the bitmap.
func AsIndices(bitmap uint16) ([]byte, error) {
	if bitmap > 1024 {
		return nil, fmt.Errorf("invalid bitmap")
	}

	res := []byte{}
	for i := 0; i < 10; i++ {
		if (bitmap & (1 << i)) != 0 {
			res = append(res, byte(i))
		}
	}
	return res, nil
}

// PublicKeyFromSuiBytes creates a PublicKey from Sui bytes or base64 string.
func PublicKeyFromSuiBytes[T string | []byte](publicKey T) (pk cryptography.PublicKey, err error) {
	var bs []byte
	switch v := any(publicKey).(type) {
	case string:
		bs, err = b64.FromBase64(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		bs = v
	}

	signatureScheme := cryptography.SignatureFlagToScheme[bs[0]]

	if signatureScheme == cryptography.ZkLoginScheme {
		return nil, fmt.Errorf("zkLogin publicKey is not supported")
	}
	return verify.PublicKeyFromRawBytes(signatureScheme, bs[1:])
}
