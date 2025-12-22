package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
)

// GetHMACSHA256 calculates HMACSHA256 hash for the given value using provided key.
func GetHMACSHA256(value []byte, key string) []byte {

	h := hmac.New(sha256.New, []byte(key))

	h.Write(value)
	hash := h.Sum(nil)

	return hash
}
