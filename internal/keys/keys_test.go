package keys

import (
	"crypto/ed25519"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateEd25519Key tests the generation of an Ed25519 private key.
func TestGenerateEd25519Key(t *testing.T) {
	privKey, err := GenerateEd25519Key()
	require.NoError(t, err, "Expected no error when generating Ed25519 key")
	assert.NotNil(t, privKey, "Expected a non-nil private key")
	assert.IsType(t, ed25519.PrivateKey{}, privKey, "Expected an Ed25519 private key")
}

// TestGenerateRSAKey tests the generation of RSA private keys with valid and invalid bit sizes.
func TestGenerateRSAKey(t *testing.T) {
	// Test with valid bit size
	privKey, err := GenerateRSAKey(2048)
	require.NoError(t, err, "Expected no error when generating RSA key with 2048 bits")
	assert.NotNil(t, privKey, "Expected a non-nil private key")
	assert.IsType(t, &rsa.PrivateKey{}, privKey, "Expected an RSA private key")

	// Test with invalid bit size
	_, err = GenerateRSAKey(1024)
	assert.Error(t, err, "Expected an error when generating RSA key with less than 2048 bits")
	assert.EqualError(t, err, "RSA key size must be at least 2048 bits", "Expected specific error message")
}

// TestSerializePrivateKey tests the serialization of private keys, both unencrypted and encrypted.
func TestSerializePrivateKey(t *testing.T) {
	// Generate an Ed25519 key for testing
	privKey, err := GenerateEd25519Key()
	require.NoError(t, err, "Failed to generate Ed25519 key for testing")

	// Test serialization without password
	pemData, err := SerializePrivateKey(privKey, "")
	require.NoError(t, err, "Expected no error when serializing without password")
	assert.Contains(t, string(pemData), "-----BEGIN PRIVATE KEY-----", "Expected PEM header for unencrypted key")

	// Test serialization with password
	pemData, err = SerializePrivateKey(privKey, "secret")
	require.NoError(t, err, "Expected no error when serializing with password")
	assert.Contains(t, string(pemData), "-----BEGIN ENCRYPTED PRIVATE KEY-----", "Expected PEM header for encrypted key")
}

// TestParsePrivateKey tests parsing of PEM-encoded private keys, both unencrypted and encrypted.
func TestParsePrivateKey(t *testing.T) {
	// Generate an Ed25519 key for testing
	privKey, err := GenerateEd25519Key()
	require.NoError(t, err, "Failed to generate Ed25519 key for testing")

	// Serialize the key without password
	pemData, err := SerializePrivateKey(privKey, "")
	require.NoError(t, err, "Failed to serialize private key without password")

	// Parse the unencrypted key
	parsedKey, err := ParsePrivateKey(pemData, "")
	require.NoError(t, err, "Expected no error when parsing unencrypted key")
	assert.Equal(t, privKey, parsedKey, "Parsed key should match original key")

	// Serialize the key with password
	pemData, err = SerializePrivateKey(privKey, "secret")
	require.NoError(t, err, "Failed to serialize private key with password")

	// Parse the encrypted key with correct password
	parsedKey, err = ParsePrivateKey(pemData, "secret")
	require.NoError(t, err, "Expected no error when parsing encrypted key with correct password")
	assert.Equal(t, privKey, parsedKey, "Parsed key should match original key")

	// Attempt to parse the encrypted key with incorrect password
	_, err = ParsePrivateKey(pemData, "wrong")
	assert.Error(t, err, "Expected an error when parsing with incorrect password")
}
