package rsa_crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	return priv, &priv.PublicKey
}

func TestCryptoMiddleware_NoEncryption(t *testing.T) {
	priv, _ := generateTestKeys(t)

	handler := CryptoMiddleware(priv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("plain-text"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rec.Code)
	}

	if rec.Body.String() != "plain-text" {
		t.Fatalf("body was modified: %q", rec.Body.String())
	}
}

func TestCryptoMiddleware_DecryptSuccess(t *testing.T) {
	priv, pub := generateTestKeys(t)

	original := []byte("secret payload")

	encrypted, err := Encrypt(pub, original)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	handler := CryptoMiddleware(priv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(encrypted))
	req.Header.Set("X-Encrypted", "rsa")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", rec.Code)
	}

	if !bytes.Equal(rec.Body.Bytes(), original) {
		t.Fatalf("decrypted body mismatch: %q", rec.Body.Bytes())
	}
}

func TestCryptoMiddleware_InvalidCiphertext(t *testing.T) {
	priv, _ := generateTestKeys(t)

	handler := CryptoMiddleware(priv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("garbage")))
	req.Header.Set("X-Encrypted", "rsa")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCryptoMiddleware_EmptyBody(t *testing.T) {
	priv, _ := generateTestKeys(t)

	handler := CryptoMiddleware(priv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Encrypted", "rsa")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
