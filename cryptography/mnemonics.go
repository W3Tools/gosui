package cryptography

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/tyler-smith/go-bip39"
)

// IsValidHardenedPath checks if a given path is compliant with SLIP-0010 m/44'/784'/{account_index}'/{change_index}'/{address_index}'.
// e.g. `m/44'/784'/0'/0'/0'`
func IsValidHardenedPath(path string) bool {
	matched, _ := regexp.MatchString("^m\\/44'\\/784'\\/[0-9]+'\\/[0-9]+'\\/[0-9]+'+$", path)
	return matched
}

// IsValidBIP32Path checks if a given path is compliant with BIP-32 for Secp256k1 and Secp256r1.
// m/54'/784'/{account_index}'/{change_index}/{address_index}
// for Secp256k1 and m/74'/784'/{account_index}'/{change_index}/{address_index} for Secp256r1.
// Note that the purpose for Secp256k1 is registered as 54, to differentiate from Ed25519 with purpose 44.
// e.g. `m/54'/784'/0'/0/0`
func IsValidBIP32Path(path string) bool {
	matched, _ := regexp.MatchString("^m\\/(54|74)'\\/784'\\/[0-9]+'\\/[0-9]+\\/[0-9]+$", path)
	return matched
}

// MnemonicToSeed converts a mnemonic phrase to a seed using KDF to derive 64 bytes of key data from mnemonic with empty password.
// mnemonics 12 words string split by spaces.
func MnemonicToSeed(mnemonics string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonics) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	return bip39.NewSeed(mnemonics, ""), nil
}

// MnemonicToSeedHex converts a mnemonic phrase to a seed in hex format.
// Derive the seed in hex format from a 12-word mnemonic string.
// mnemonics 12 words string split by spaces.
func MnemonicToSeedHex(mnemonics string) (string, error) {
	seed, err := MnemonicToSeed(mnemonics)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(seed), nil
}

// GenerateMnemonic generates a new mnemonic phrase with 12 words.
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}
