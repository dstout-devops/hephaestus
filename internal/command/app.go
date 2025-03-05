package command

import (
	"flag"

	"github.com/dstout-devops/hephaestus/internal/logger"
)

type App struct {
	configPath string
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(args []string) error {
	log := logger.NewLogger()
	_ = log
	fs := flag.NewFlagSet("hephaestus", flag.ExitOnError)
	fs.StringVar(&a.configPath, "config", "config.yaml", "Path to the configuration file")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return nil
}
