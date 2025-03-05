package command

import (
	"errors"
	"os"

	"github.com/dstout-devops/hephaestus/internal/config"
	"github.com/dstout-devops/hephaestus/internal/csr"
	"github.com/dstout-devops/hephaestus/internal/keys"
	"github.com/dstout-devops/hephaestus/internal/logger"
)

// App represents the application, holding state and methods.
type App struct {
	log     logger.Logger // Logger for troubleshooting
	cfg     config.Config // Loaded configuration
	privKey interface{}   // Generated private key
	csr     []byte        // Generated CSR data
	keyGen  KeyGenerator  // Dependency for key generation
}

// NewApp creates a new App instance with a logger and key generator.
func NewApp(keyGen KeyGenerator) *App {
	log := logger.NewLogger()
	if keyGen == nil {
		keyGen = &DefaultKeyGenerator{}
	}
	return &App{
		log:    log,
		keyGen: keyGen,
	}
}

// Run executes the main logic of the application.
func (a *App) Run() error {
	if err := a.LoadConfig(); err != nil {
		return err
	}
	if err := a.GenerateKey(); err != nil {
		return err
	}
	if err := a.GenerateCSR(); err != nil {
		return err
	}
	return nil
}

// LoadConfig loads the application configuration.
func (a *App) LoadConfig() error {
	// Assuming a.log is a logger field in App
	a.log.Info("Loading configuration...")
	cfg, err := config.ConfigLoader.LoadConfig(nil)
	if err != nil {
		a.log.Error("Failed to load config", "error", err)
		return errors.New("configuration loading failed")
	}
	a.cfg = cfg // Assuming a.cfg is a field in App to store the config
	a.log.Info("Configuration loaded successfully")
	return nil
}

// GenerateKey generates the private key and stores it in memory.
func (a *App) GenerateKey() error {
	a.log.Info("Generating private key...")
	privKey, err := a.keyGen.GenerateEd25519Key()
	if err != nil {
		a.log.Error("Failed to generate private key", "error", err)
		return errors.New("private key generation failed")
	}
	a.privKey = privKey
	a.log.Info("Private key generated successfully")
	return nil
}

// GenerateCSR generates the CSR and stores it in memory.
func (a *App) GenerateCSR() error {
	a.log.Info("Generating CSR...")
	csrPem, err := csr.GenerateCSR(a.privKey, a.cfg.CSR)
	if err != nil {
		a.log.Error("Failed to generate CSR", "error", err)
		return errors.New("CSR generation failed")
	}
	a.csr = csrPem
	a.log.Info("CSR generated successfully")
	return nil
}

// WriteKeyToFile optionally saves the private key to a file.
func (a *App) WriteKeyToFile(path string) error {
	if a.privKey == nil {
		return errors.New("no private key available to write")
	}

	if path == "" {
		if a.cfg.Key.Output != "" {
			path = a.cfg.Key.Output // Use config path if available
		} else {
			path = "private.key" // Default path
		}
	}

	pemKey, err := keys.SerializePrivateKey(a.privKey, "")
	if err != nil {
		a.log.Error("Failed to serialize private key", "error", err)
		return errors.New("private key serialization failed")
	}

	err = os.WriteFile(path, pemKey, 0600)
	if err != nil {
		a.log.Error("Failed to save private key", "error", err, "path", path)
		return errors.New("private key saving failed")
	}
	a.log.Info("Private key saved successfully", "path", path)
	return nil
}

// WriteCSRToFile optionally saves the CSR to a file.
func (a *App) WriteCSRToFile(path string) error {
	if a.csr == nil {
		return errors.New("no CSR available to write")
	}

	if path == "" {
		path = "host.csr" // Default path
	}

	err := os.WriteFile(path, a.csr, 0644)
	if err != nil {
		a.log.Error("Failed to save CSR", "error", err, "path", path)
		return errors.New("CSR saving failed")
	}
	a.log.Info("CSR saved successfully", "path", path)
	return nil
}
