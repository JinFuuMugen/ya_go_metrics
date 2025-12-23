package cryptography

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
)

// ValidateHashMiddleware returns an HTTP middleware that validates the HMAC-SHA256 hash of incoming requests.
// When a secret key is configured, the middleware also adds a
// "HashSHA256" header to the response.
func ValidateHashMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			hashString := r.Header.Get("HashSHA256")
			if hashString != "" {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "cannot read body", http.StatusBadRequest)
				}
				defer r.Body.Close()

				newBody := io.NopCloser(bytes.NewReader(body))

				hash := GetHMACSHA256(body, cfg.Key)
				calculatedHashString := hex.EncodeToString(hash)

				if hashString != calculatedHashString {
					http.Error(w, "hash differs", http.StatusBadRequest)
					return
				}

				r.Body = newBody
			}

			next.ServeHTTP(w, r)

			if cfg.Key != "" {
				responseHash := GetHMACSHA256([]byte(""), cfg.Key)
				responseHashString := hex.EncodeToString(responseHash)
				w.Header().Set("HashSHA256", responseHashString)
			}
		})
	}
}
