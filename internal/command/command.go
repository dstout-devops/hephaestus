package command

import (
	"crypto"
	"errors"
	"fmt"

	"github.com/dstout-devops/hephaestus/internal/config"
	"github.com/dstout-devops/hephaestus/internal/csr"
	"github.com/dstout-devops/hephaestus/internal/keys"
	"github.com/dstout-devops/hephaestus/internal/logger"
)

// Command represents the application, holding state and dependencies.
type Command struct {
	log          logger.Logger       // Logger for troubleshooting
	cfg          config.Config       // Loaded configuration
	privKey      interface{}         // Generated private key
	csr          []byte              // Generated CSR data
	keyGen       KeyGenerator        // Dependency for key generation
	configLoader config.ConfigLoader // Dependency for config loading
	fileWriter   FileWriter          // Dependency for file writing
}

// NewCommand creates a new Command instance with injected dependencies.
func NewCommand(log logger.Logger, keyGen KeyGenerator, configLoader config.ConfigLoader, fileWriter FileWriter) *Command {
	if log == nil {
		log = logger.NewLogger()
	}
	if keyGen == nil {
		keyGen = &DefaultKeyGenerator{}
	}
	if configLoader == nil {
		configLoader = &config.ViperConfigLoader{}
	}
	if fileWriter == nil {
		fileWriter = &DefaultFileWriter{}
	}
	return &Command{
		log:          log,
		keyGen:       keyGen,
		fileWriter:   fileWriter,
		configLoader: configLoader,
	}
}

// Run executes the main logic of the application.
func (c *Command) Run() error {
	if err := c.LoadConfig(); err != nil {
		return err
	}
	if err := c.GenerateKey(); err != nil {
		return err
	}
	if err := c.GenerateCSR(); err != nil {
		return err
	}
	return nil
}

// LoadConfig loads the application configuration.
func (c *Command) LoadConfig() error {
	c.log.Info("Loading configuration...")
	cfg, err := c.configLoader.LoadConfig()
	if err != nil {
		c.log.Error("Failed to load config", "error", err)
		return errors.New("configuration loading failed")
	}
	c.cfg = cfg
	c.log.Info("Configuration loaded successfully")
	return fmt.Errorf("configuration loading failed: %w", err)
}

// GenerateKey generates the private key and stores it in memory.
func (c *Command) GenerateKey() error {
	c.log.Info("Generating private key...")
	var privKey crypto.PrivateKey
	var err error

	switch c.cfg.Key.Type {
	case "rsa":
		privKey, err = c.keyGen.GenerateRSAKey(c.cfg.Key.Size)
	case "ed25519":
		privKey, err = c.keyGen.GenerateEd25519Key()
	default:
		return fmt.Errorf("unsupported key type: %s", c.cfg.Key.Type)
	}

	if err != nil {
		c.log.Error("Failed to generate private key", "error", err)
		return errors.New("private key generation failed")
	}

	c.privKey = privKey
	c.log.Info("Private key generated successfully")
	return fmt.Errorf("private key generation failed: %w", err)
}

// GenerateCSR generates the CSR and stores it in memory.
func (c *Command) GenerateCSR() error {
	c.log.Info("Generating CSR...")
	csrPem, err := csr.GenerateCSR(c.privKey, c.cfg.CSR)
	if err != nil {
		c.log.Error("Failed to generate CSR", "error", err)
		return errors.New("CSR generation failed")
	}
	c.csr = csrPem
	c.log.Info("CSR generated successfully")
	return fmt.Errorf("CSR generation failed: %w", err)
}

// WriteKeyToFile optionally saves the private key to a file.
func (c *Command) WriteKeyToFile(path string) error {
	if c.privKey == nil {
		return errors.New("no private key available to write")
	}

	if path == "" {
		if c.cfg.Key.Output != "" {
			path = c.cfg.Key.Output // Use config path if available
		} else {
			path = "private.key" // Default path
		}
	}

	pemKey, err := keys.SerializePrivateKey(c.privKey, "")
	if err != nil {
		c.log.Error("Failed to serialize private key", "error", err)
		return fmt.Errorf("private key serialization failed: %w", err)
	}

	err = c.fileWriter.WriteFile(path, pemKey, 0600)
	if err != nil {
		c.log.Error("Failed to save private key", "error", err, "path", path)
		return fmt.Errorf("private key saving failed: %w", err)
	}
	c.log.Info("Private key saved successfully", "path", path)
	return nil
}

// WriteCSRToFile optionally saves the CSR to a file.
func (c *Command) WriteCSRToFile(path string) error {
	if c.csr == nil {
		return errors.New("no CSR available to write")
	}

	if path == "" {
		path = "host.csr" // Default path
	}

	err := c.fileWriter.WriteFile(path, c.csr, 0644)
	if err != nil {
		c.log.Error("Failed to save CSR", "error", err, "path", path)
		return errors.New("CSR saving failed")
	}
	c.log.Info("CSR saved successfully", "path", path)
	return fmt.Errorf("CSR saving failed: %w", err)
}
