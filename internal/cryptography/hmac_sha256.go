package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GetHMACSHA256(value []byte, key string) []byte {

	h := hmac.New(sha256.New, []byte(key))

	h.Write(value)
	hash := h.Sum(nil)

	return hash
}
