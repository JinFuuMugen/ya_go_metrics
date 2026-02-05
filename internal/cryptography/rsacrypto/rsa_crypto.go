package rsacrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// func LoadPublicKey(path string) (*rsa.PublicKey, error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot read publickey file: %w", err)
// 	}

// 	block, _ := pem.Decode(data)
// 	if block == nil {
// 		return nil, fmt.Errorf("invalid PEM public key")
// 	}

// 	cert, err := x509.ParseCertificate(block.Bytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot parse certificate: %w", err)
// 	}

// 	pub, ok := cert.PublicKey.(*rsa.PublicKey)
// 	if !ok {
// 		return nil, fmt.Errorf("not RSA public key")
// 	}

// 	return pub, nil
// }

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %w", err)
	}

	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return pub, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read privatekey file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func Encrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func Decrypt(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
