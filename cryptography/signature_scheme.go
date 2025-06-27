package cryptography

// SignatureScheme defines the type for signature schemes used in Sui.
type SignatureScheme = string

// SignatureFlag defines the type for signature flags, which are used to identify the signature scheme in serialized signatures.
type SignatureFlag = uint8

var (
	// Ed25519Scheme is the signature scheme for Ed25519.
	Ed25519Scheme SignatureScheme = "ED25519"
	// Secp256k1Scheme is the signature scheme for Secp256k1.
	Secp256k1Scheme SignatureScheme = "Secp256k1"
	// Secp256r1Scheme is the signature scheme for Secp256r1.
	Secp256r1Scheme SignatureScheme = "Secp256r1"
	// MultiSigScheme is the signature scheme for MultiSig.
	MultiSigScheme SignatureScheme = "MultiSig"
	// ZkLoginScheme is the signature scheme for ZkLogin.
	ZkLoginScheme SignatureScheme = "ZkLogin"
)

// SignatureSchemeToFlag is a map that associates signature schemes with their corresponding flags.
var SignatureSchemeToFlag = map[SignatureScheme]SignatureFlag{
	Ed25519Scheme:   0x00,
	Secp256k1Scheme: 0x01,
	Secp256r1Scheme: 0x02,
	MultiSigScheme:  0x03,
	ZkLoginScheme:   0x05,
}

// SignatureSchemeToSize is a map that associates signature schemes with their corresponding signature sizes in bytes.
var SignatureSchemeToSize = map[SignatureScheme]int{
	Ed25519Scheme:   32,
	Secp256k1Scheme: 33,
	Secp256r1Scheme: 33,
}

// SignatureFlagToScheme is a map that associates signature flags with their corresponding signature schemes.
var SignatureFlagToScheme = map[SignatureFlag]SignatureScheme{
	0x00: Ed25519Scheme,
	0x01: Secp256k1Scheme,
	0x02: Secp256r1Scheme,
	0x03: MultiSigScheme,
	0x05: ZkLoginScheme,
}
