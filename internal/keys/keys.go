package keys

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"

	"github.com/youmark/pkcs8"
)

// GenerateEd25519Key generates a new ed25519 private key.
func GenerateEd25519Key() (ed25519.PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	return priv, err
}

// GenerateRSAKey generates a new RSA private key with the specified bit size.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, errors.New("RSA key size must be at least 2048 bits")
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// SerializePrivateKey serializes the private key to PEM format, optionally encrypted with a password.
func SerializePrivateKey(key crypto.PrivateKey, password string) ([]byte, error) {
	var pass []byte
	if password != "" {
		pass = []byte(password)
	}
	der, err := pkcs8.MarshalPrivateKey(key, pass, nil)
	if err != nil {
		return nil, err
	}
	pemType := "PRIVATE KEY"
	if pass != nil {
		pemType = "ENCRYPTED PRIVATE KEY"
	}
	pemBlock := &pem.Block{
		Type:  pemType,
		Bytes: der,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

// ParsePrivateKey parses a PEM-encoded private key, optionally decrypting it with a password.
func ParsePrivateKey(pemData []byte, password string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	var pass []byte
	if password != "" {
		pass = []byte(password)
	}
	key, _, err := pkcs8.ParsePrivateKey(block.Bytes, pass)
	if err != nil {
		return nil, err
	}
	privKey, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not a crypto.PrivateKey")
	}
	return privKey, nil
}
