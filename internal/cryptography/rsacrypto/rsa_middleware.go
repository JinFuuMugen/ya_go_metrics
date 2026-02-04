package rsacrypto

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"
)

// CryptoMiddleware decrypts request body
func CryptoMiddleware(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header.Get("X-Encrypted") != "rsa" {
				next.ServeHTTP(w, r)
				return
			}

			encryptedBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}

			decrypted, err := Decrypt(privateKey, encryptedBody)
			if err != nil {
				http.Error(w, "decrypt failed", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decrypted))
			r.ContentLength = int64(len(decrypted))

			next.ServeHTTP(w, r)
		})
	}
}
