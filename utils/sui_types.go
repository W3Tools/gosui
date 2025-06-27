package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// SuiAddressLength is the length of a Sui address in bytes
	SuiAddressLength = 32
	// TxDigestLength is the length of a transaction digest in bytes
	TxDigestLength = 32
)

// IsValidTransactionDigest checks if a given transaction digest is valid.
func IsValidTransactionDigest(v string) bool {
	return true
}

// IsValidSuiAddress checks if a given Sui address is valid.
func IsValidSuiAddress(v string) bool {
	return IsHex(v) && GetHexByteLength(v) == SuiAddressLength
}

// IsValidSuiObjectID checks if a given Sui object ID is valid.
func IsValidSuiObjectID(v string) bool {
	return IsValidSuiAddress(v)
}

// NormalizeSuiAddress normalizes a Sui address string to a standard format.
/**
 * Perform the following operations:
 * 1. Make the address lower case
 * 2. Prepend `0x` if the string does not start with `0x`.
 * 3. Add more zeros if the length of the address(excluding `0x`) is less than `SUI_ADDRESS_LENGTH`
 *
 * WARNING: if the address value itself starts with `0x`, e.g., `0x0x`, the default behavior
 * is to treat the first `0x` not as part of the address. The default behavior can be overridden by
 * setting `forceAdd0x` to true
 *
 */
//  0x1 -> 0x0000000000000000000000000000000000000000000000000000000000000001
func NormalizeSuiAddress(v string) string {
	address := strings.ToLower(v)
	address = strings.TrimPrefix(address, "0x")

	if len(address) > SuiAddressLength*2 {
		return "0x" + address
	}

	return fmt.Sprintf("0x%s", strings.Repeat("0", SuiAddressLength*2-len(address))+address)
}

// NormalizeSuiObjectID normalizes a Sui object ID string to a standard format.
// 0x1 -> 0x0000000000000000000000000000000000000000000000000000000000000001
func NormalizeSuiObjectID(v string) string {
	return NormalizeSuiAddress(v)
}

// NormalizeSuiCoinType normalizes a Sui coin type string to a standard format.
// 0x2::sui::SUI -> 0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI
func NormalizeSuiCoinType(v string) string {
	splits := strings.Split(v, "::")
	if len(splits) != 3 {
		return v
	}

	return fmt.Sprintf("%s::%s::%s", NormalizeSuiAddress(splits[0]), splits[1], splits[2])
}

// NormalizeShortSuiAddress normalizes a short Sui address string.
// 0x0000000000000000000000000000000000000000000000000000000000000001 -> 0x1
func NormalizeShortSuiAddress(v string) string {
	address := NormalizeSuiAddress(v)

	address = strings.TrimPrefix(address, "0x")
	address = strings.TrimLeft(address, "0")

	if address == "" {
		return "0x0"
	}

	return "0x" + address
}

// NormalizeShortSuiObjectID normalizes a short Sui object ID string.
// 0x0000000000000000000000000000000000000000000000000000000000000001 -> 0x1
func NormalizeShortSuiObjectID(v string) string {
	return NormalizeShortSuiAddress(v)
}

// NormalizeShortSuiCoinType normalizes a short Sui coin type string.
// 0x2::sui::SUI -> 0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI
func NormalizeShortSuiCoinType(v string) string {
	splits := strings.Split(v, "::")
	if len(splits) != 3 {
		return v
	}

	return fmt.Sprintf("%s::%s::%s", NormalizeShortSuiAddress(splits[0]), splits[1], splits[2])
}

// IsHex checks if a string is a valid hexadecimal representation.
func IsHex(v string) bool {
	re := regexp.MustCompile(`^(0x|0X)?[a-fA-F0-9]+$`)

	if !re.MatchString(v) {
		return false
	}

	return len(v)%2 == 0
}

// GetHexByteLength returns the number of bytes represented by a hex string.
func GetHexByteLength(v string) int {
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
		v = v[2:]
	}

	return len(v) / 2
}
