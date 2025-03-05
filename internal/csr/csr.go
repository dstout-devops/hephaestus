package csr

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"net"

	"github.com/dstout-devops/hephaestus/internal/config" // Replace with your actual module path
)

// GenerateCSR creates a Certificate Signing Request (CSR) using the provided private key and CSR configuration.
func GenerateCSR(privKey interface{}, cfg config.CSRConfig) ([]byte, error) {
	// Create the CSR subject using pkix.Name with default zero values
	subject := pkix.Name{
		CommonName: cfg.CommonName, // CommonName is a string, so it can be set directly
	}

	// Conditionally set Organization if not empty
	if cfg.Organization != "" {
		subject.Organization = []string{cfg.Organization}
	}

	// Conditionally set OrganizationalUnit if not empty
	if cfg.OrganizationalUnit != "" {
		subject.OrganizationalUnit = []string{cfg.OrganizationalUnit}
	}

	// Conditionally set Country if not empty
	if cfg.Country != "" {
		subject.Country = []string{cfg.Country}
	}

	// Conditionally set Province (State) if not empty
	if cfg.State != "" {
		subject.Province = []string{cfg.State}
	}

	// Conditionally set Locality if not empty
	if cfg.Locality != "" {
		subject.Locality = []string{cfg.Locality}
	}

	// Handle IP address if provided
	var ipAddresses []net.IP
	if cfg.IPAddress != "" {
		ip := net.ParseIP(cfg.IPAddress)
		if ip == nil {
			return nil, errors.New("invalid IP address in configuration")
		}
		ipAddresses = append(ipAddresses, ip)
	}

	// Create the CSR template
	csrTemplate := &x509.CertificateRequest{
		Subject:     subject,
		IPAddresses: ipAddresses,
		// Add DNSNames or EmailAddresses here if your config supports them
	}

	// Generate the CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privKey)
	if err != nil {
		return nil, err
	}

	// PEM-encode the CSR
	csrPemBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	}
	csrPem := pem.EncodeToMemory(csrPemBlock)
	if csrPem == nil {
		return nil, errors.New("failed to PEM-encode CSR")
	}

	return csrPem, nil
}
