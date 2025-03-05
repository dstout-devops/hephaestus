package config

import (
	"os"

	"github.com/spf13/viper"
)

type ConfigLoader interface {
	LoadConfig() (Config, error)
}

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

// ViperConfigLoader implements the ConfigLoader interface using Viper.
type ViperConfigLoader struct {
	v *viper.Viper
}

// NewViperConfigLoader creates a new ViperConfigLoader with default settings.
func NewViperConfigLoader() *ViperConfigLoader {
	v := viper.New()
	v.SetDefault("key.output", "private.key")
	v.SetDefault("certificate.output", "certificate.pem")
	return &ViperConfigLoader{v: v}
}

// LoadConfig loads the configuration using the Viper instance.
func (l *ViperConfigLoader) LoadConfig() (Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		l.v.SetConfigFile(configPath)
	} else {
		l.v.SetConfigName("config")
		l.v.SetConfigType("yaml")
		l.v.AddConfigPath(".")
	}

	err := l.v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = l.v.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// LoadConfig is a convenience function for backward compatibility.
func LoadConfig() (Config, error) {
	loader := NewViperConfigLoader()
	return loader.LoadConfig()
}
