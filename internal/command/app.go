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
	log logger.Logger
	cfg config.Config
}

// NewApp creates a new App instance with a logger.
func NewApp() *App {
	log := logger.NewLogger()
	return &App{log: log}
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
	a.log.Info("Loading configuration...")
	var err error
	a.cfg, err = config.LoadConfig()
	if err != nil {
		a.log.Error("Failed to load config", "error", err)
		return errors.New("configuration loading failed")
	}
	a.log.Info("Configuration loaded successfully")
	return nil
}

// GenerateKey generates and saves the private key.
func (a *App) GenerateKey() error {
	a.log.Info("Generating private key...")
	privKey, err := keys.GenerateEd25519Key() // Or GenerateRSAKey if needed
	if err != nil {
		a.log.Error("Failed to generate private key", "error", err)
		return errors.New("private key generation failed")
	}

	pemKey, err := keys.SerializePrivateKey(privKey, "")
	if err != nil {
		a.log.Error("Failed to serialize private key", "error", err)
		return errors.New("private key serialization failed")
	}

	err = os.WriteFile(a.cfg.Key.Output, pemKey, 0600)
	if err != nil {
		a.log.Error("Failed to save private key", "error", err, "path", a.cfg.Key.Output)
		return errors.New("private key saving failed")
	}
	a.log.Info("Private key saved successfully", "path", a.cfg.Key.Output)
	return nil
}

// GenerateCSR generates and saves the CSR.
func (a *App) GenerateCSR() error {
	a.log.Info("Generating CSR...")
	csrPem, err := csr.GenerateCSR(a.cfg.Key, a.cfg.CSR) // Assuming PrivateKey is stored or accessible
	if err != nil {
		a.log.Error("Failed to generate CSR", "error", err)
		return errors.New("CSR generation failed")
	}

	csrPath := "host.csr" // Could be configurable via a.cfg
	err = os.WriteFile(csrPath, csrPem, 0644)
	if err != nil {
		a.log.Error("Failed to save CSR", "error", err, "path", csrPath)
		return errors.New("CSR saving failed")
	}
	a.log.Info("CSR saved successfully", "path", csrPath)
	return nil
}
