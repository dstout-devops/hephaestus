package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTempConfig creates a temporary YAML file with the given content and returns its path.
func createTempConfig(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config*.yml")
	assert.NoError(t, err, "Failed to create temporary file")
	_, err = tmpfile.Write([]byte(content))
	assert.NoError(t, err, "Failed to write to temporary file")
	err = tmpfile.Close()
	assert.NoError(t, err, "Failed to close temporary file")
	return tmpfile.Name()
}

// TestLoadConfigAllFields tests loading a configuration with all fields present.
func TestLoadConfigAllFields(t *testing.T) {
	configContent := `
key:
  bits: 2048
  output: "private.key"
csr:
  common_name: "example.com"
  organization: "Example Org"
  organizational_unit: "IT"
  country: "US"
  state: "California"
  locality: "San Francisco"
  ip_address: "192.168.1.1"
endpoint: "https://example.com"
esf:
  program_id: "prog123"
  service_id: "svc456"
  application_id: "app789"
certificate:
  output: "certificate.pem"
`
	configPath := createTempConfig(t, configContent)
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Remove(configPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg, err := LoadConfig()
	assert.NoError(t, err, "Expected no error when loading complete config")

	expected := Config{
		Key: KeyConfig{Bits: 2048, Output: "private.key"},
		CSR: CSRConfig{
			CommonName:         "example.com",
			Organization:       "Example Org",
			OrganizationalUnit: "IT",
			Country:            "US",
			State:              "California",
			Locality:           "San Francisco",
			IPAddress:          "192.168.1.1",
		},
		Endpoint: "https://example.com",
		ESF: ESFConfig{
			ProgramID:     "prog123",
			ServiceID:     "svc456",
			ApplicationID: "app789",
		},
		Certificate: CertificateConfig{Output: "certificate.pem"},
	}

	assert.Equal(t, expected, cfg, "Loaded config does not match expected")
}

// TestLoadConfigPartial tests loading a configuration with some fields missing.
func TestLoadConfigPartial(t *testing.T) {
	configContent := `
key:
  bits: 2048
csr:
  common_name: "example.com"
endpoint: "https://example.com"
`
	configPath := createTempConfig(t, configContent)
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Remove(configPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg, err := LoadConfig()
	assert.NoError(t, err, "Expected no error when loading partial config")

	expected := Config{
		Key: KeyConfig{Bits: 2048, Output: "private.key"}, // default output
		CSR: CSRConfig{
			CommonName: "example.com",
			// Other fields are zero values
		},
		Endpoint: "https://example.com",
		ESF:      ESFConfig{
			// All fields are zero values
		},
		Certificate: CertificateConfig{Output: "certificate.pem"}, // default output
	}

	assert.Equal(t, expected, cfg, "Loaded config does not match expected with defaults")
}

// TestLoadConfigInvalidYAML tests loading a configuration with invalid YAML.
func TestLoadConfigInvalidYAML(t *testing.T) {
	configContent := `key: [invalid`
	configPath := createTempConfig(t, configContent)
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Remove(configPath)
	defer os.Unsetenv("CONFIG_PATH")

	_, err := LoadConfig()
	assert.Error(t, err, "Expected error when loading invalid YAML")
}

// TestLoadConfigMissingFile tests loading a configuration when the file is missing.
func TestLoadConfigMissingFile(t *testing.T) {
	os.Setenv("CONFIG_PATH", "/nonexistent/config.yml")
	defer os.Unsetenv("CONFIG_PATH")

	_, err := LoadConfig()
	assert.Error(t, err, "Expected error when config file is missing")
}
