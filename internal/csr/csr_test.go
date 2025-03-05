package csr

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"testing"

	"github.com/dstout-devops/hephaestus/internal/config" // Adjust to your actual module path
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateCSR_Success tests the GenerateCSR function with a fully populated configuration.
func TestGenerateCSR_Success(t *testing.T) {
	// Generate an RSA private key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate RSA private key")

	// Define a fully populated CSR configuration
	cfg := config.CSRConfig{
		CommonName:         "test.com",
		Organization:       "Test Org",
		OrganizationalUnit: "IT",
		Country:            "US",
		State:              "California",
		Locality:           "San Francisco",
		IPAddress:          "192.168.1.1",
	}

	// Generate the CSR
	csrPem, err := GenerateCSR(privKey, cfg)
	require.NoError(t, err, "GenerateCSR should not return an error")
	assert.NotEmpty(t, csrPem, "CSR PEM should not be empty")

	// Decode the PEM to verify its contents
	block, _ := pem.Decode(csrPem)
	require.NotNil(t, block, "PEM decoding should return a non-nil block")
	assert.Equal(t, "CERTIFICATE REQUEST", block.Type, "PEM block type should be CERTIFICATE REQUEST")

	// Parse the CSR to verify subject and IP fields
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "failed to parse CSR")

	// Verify subject fields
	assert.Equal(t, cfg.CommonName, csr.Subject.CommonName, "CommonName should match")
	assert.Equal(t, cfg.Organization, csr.Subject.Organization[0], "Organization should match")
	assert.Equal(t, cfg.OrganizationalUnit, csr.Subject.OrganizationalUnit[0], "OrganizationalUnit should match")
	assert.Equal(t, cfg.Country, csr.Subject.Country[0], "Country should match")
	assert.Equal(t, cfg.State, csr.Subject.Province[0], "State should match")
	assert.Equal(t, cfg.Locality, csr.Subject.Locality[0], "Locality should match")

	// Verify IP addresst := s.T()
	assert.Len(t, csr.IPAddresses, 1, "CSR should have exactly one IP address")
	expectedIP := net.ParseIP("192.168.1.1")
	actualIP := csr.IPAddresses[0]
	assert.True(t, expectedIP.Equal(actualIP), "IP addresses should match: expected %v, got %v", expectedIP, actualIP)
}

// TestGenerateCSR_InvalidIP tests the GenerateCSR function with an invalid IP address.
func TestGenerateCSR_InvalidIP(t *testing.T) {
	// Generate an RSA private key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate RSA private key")

	// Define a configuration with an invalid IP address
	cfg := config.CSRConfig{
		IPAddress: "invalid_ip",
	}

	// Attempt to generate the CSR and expect an error
	_, err = GenerateCSR(privKey, cfg)
	require.Error(t, err, "GenerateCSR should return an error for invalid IP")
	assert.Equal(t, "invalid IP address in configuration", err.Error(), "error message should match expected")
}

// TestGenerateCSR_MinimalConfig tests the GenerateCSR function with a minimal configuration.
func TestGenerateCSR_MinimalConfig(t *testing.T) {
	// Generate an RSA private key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate RSA private key")

	// Define a minimal CSR configuration
	cfg := config.CSRConfig{
		CommonName: "test.com",
	}

	// Generate the CSR
	csrPem, err := GenerateCSR(privKey, cfg)
	require.NoError(t, err, "GenerateCSR should not return an error")
	assert.NotEmpty(t, csrPem, "CSR PEM should not be empty")

	// Decode the PEM to verify its contents
	block, _ := pem.Decode(csrPem)
	require.NotNil(t, block, "PEM decoding should return a non-nil block")
	assert.Equal(t, "CERTIFICATE REQUEST", block.Type, "PEM block type should be CERTIFICATE REQUEST")

	// Parse the CSR to verify subject and IP fields
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "failed to parse CSR")

	// Verify subject fields
	assert.Equal(t, cfg.CommonName, csr.Subject.CommonName, "CommonName should match")
	assert.Empty(t, csr.Subject.Organization, "Organization should be empty")
	assert.Empty(t, csr.Subject.OrganizationalUnit, "OrganizationalUnit should be empty")
	assert.Empty(t, csr.Subject.Country, "Country should be empty")
	assert.Empty(t, csr.Subject.Province, "State should be empty")
	assert.Empty(t, csr.Subject.Locality, "Locality should be empty")

	// Verify no IP addresses are set
	assert.Empty(t, csr.IPAddresses, "IPAddresses should be empty")
}
