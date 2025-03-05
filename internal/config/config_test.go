package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestViperConfigLoader_LoadConfig tests loading a complete config file.
func TestViperConfigLoader_LoadConfig(t *testing.T) { // Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Write test config data
	configData := []byte(`
key:
  bits: 2048
csr:
  common_name: "example.com"
  organization: "Example Org"
`)
	err := os.WriteFile(configPath, configData, 0644)
	assert.NoError(t, err)

	// Set CONFIG_PATH
	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")

	// Create ViperConfigLoader
	loader := NewViperConfigLoader()

	// Load config
	cfg, err := loader.LoadConfig()
	assert.NoError(t, err)

	// Check values from file
	assert.Equal(t, 2048, cfg.Key.Size)
	assert.Equal(t, "example.com", cfg.CSR.CommonName)
	assert.Equal(t, "Example Org", cfg.CSR.Organization)
	// Check defaults
	assert.Equal(t, "private.key", cfg.Key.Output)
	assert.Equal(t, "certificate.pem", cfg.Certificate.Output)
}

// TestViperConfigLoader_LoadConfig_Defaults tests defaults when fields are missing.
func TestViperConfigLoader_LoadConfig_Defaults(t *testing.T) { // Create a temporary directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Write a minimal YAML config file
	configData := []byte(`
csr:
  common_name: "example.com"
`)
	err := os.WriteFile(configPath, configData, 0644)
	assert.NoError(t, err)

	os.Setenv("CONFIG_PATH", configPath)
	defer os.Unsetenv("CONFIG_PATH")

	loader := NewViperConfigLoader()
	cfg, err := loader.LoadConfig()
	assert.NoError(t, err)

	// Check zero values for unspecified fields
	assert.Equal(t, 0, cfg.Key.Size)
	assert.Equal(t, "example.com", cfg.CSR.CommonName)
	assert.Equal(t, "", cfg.CSR.Organization)
	// Check defaults set by Viper
	assert.Equal(t, "private.key", cfg.Key.Output)
	assert.Equal(t, "certificate.pem", cfg.Certificate.Output)
}

// TestViperConfigLoader_LoadConfig_Error tests error when file is missing.
func TestViperConfigLoader_LoadConfig_Error(t *testing.T) {
	// Set CONFIG_PATH to a non-existent file
	os.Setenv("CONFIG_PATH", "non_existent.yml")
	defer os.Unsetenv("CONFIG_PATH")

	loader := NewViperConfigLoader()
	_, err := loader.LoadConfig()
	assert.Error(t, err)
}

// TestViperConfigLoader_LoadConfig_Malformed tests error with malformed YAML.
func TestViperConfigLoader_LoadConfig_Malformed(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config.yml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write malformed YAML
	err = os.WriteFile(tempFile.Name(), []byte(`key: [}`), 0644)
	assert.NoError(t, err)

	os.Setenv("CONFIG_PATH", tempFile.Name())
	defer os.Unsetenv("CONFIG_PATH")

	loader := NewViperConfigLoader()
	_, err = loader.LoadConfig()
	assert.Error(t, err)
}
