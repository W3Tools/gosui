package verify

import (
	"fmt"

	"github.com/W3Tools/gosui/cryptography"
	"github.com/W3Tools/gosui/keypairs/ed25519"
	"github.com/W3Tools/gosui/keypairs/secp256k1"
	"github.com/W3Tools/gosui/keypairs/secp256r1"
)

// PublicKeyFromRawBytes converts raw bytes to a PublicKey based on the provided signature scheme.
func PublicKeyFromRawBytes(signatureScheme cryptography.SignatureScheme, bs []byte) (cryptography.PublicKey, error) {
	switch signatureScheme {
	case cryptography.Ed25519Scheme:
		return ed25519.NewPublicKey(bs)
	case cryptography.Secp256k1Scheme:
		return secp256k1.NewPublicKey(bs)
	case cryptography.Secp256r1Scheme:
		return secp256r1.NewPublicKey(bs)
	case cryptography.MultiSigScheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	case cryptography.ZkLoginScheme:
		return nil, fmt.Errorf("unimplemented %v", signatureScheme)
	default:
		return nil, fmt.Errorf("unsupported signature scheme %v", signatureScheme)
	}
}
