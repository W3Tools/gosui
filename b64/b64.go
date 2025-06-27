package b64

import "encoding/base64"

// FromBase64 decodes a base64 encoded string into a byte slice.
func FromBase64(base64String string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64String)
}

// ToBase64 encodes a byte slice into a base64 encoded string.
func ToBase64(bs []byte) string {
	return base64.StdEncoding.EncodeToString(bs)
}
