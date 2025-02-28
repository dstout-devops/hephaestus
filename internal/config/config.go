package config

import (
	"os"

	"github.com/spf13/viper"
)

// Config represents the top-level configuration structure.
type Config struct {
	Key         KeyConfig         `mapstructure:"key"`
	CSR         CSRConfig         `mapstructure:"csr"`
	Endpoint    string            `mapstructure:"endpoint"`
	ESF         ESFConfig         `mapstructure:"esf"`
	Certificate CertificateConfig `mapstructure:"certificate"`
}

// KeyConfig holds key-related settings.
type KeyConfig struct {
	Bits   int    `mapstructure:"bits"`
	Output string `mapstructure:"output"`
}

// CSRConfig holds CSR-related settings.
type CSRConfig struct {
	CommonName         string `mapstructure:"common_name"`
	Organization       string `mapstructure:"organization"`
	OrganizationalUnit string `mapstructure:"organizational_unit"`
	Country            string `mapstructure:"country"`
	State              string `mapstructure:"state"`
	Locality           string `mapstructure:"locality"`
	IPAddress          string `mapstructure:"ip_address"`
}

// ESFConfig holds ESF identifiers.
type ESFConfig struct {
	ProgramID     string `mapstructure:"program_id"`
	ServiceID     string `mapstructure:"service_id"`
	ApplicationID string `mapstructure:"application_id"`
}

// CertificateConfig holds certificate-related settings.
type CertificateConfig struct {
	Output string `mapstructure:"output"`
}

// LoadConfig loads the application configuration from a YAML file.
func LoadConfig() (Config, error) {
	// Set default values for optional fields
	viper.SetDefault("key.output", "private.key")
	viper.SetDefault("certificate.output", "certificate.pem")

	// Check if a custom config path is provided via environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Default to looking for config.yaml in the current directory
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// Unmarshal the configuration into the Config struct
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
