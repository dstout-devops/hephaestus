package main

import (
	"fmt"
	"os"

	"github.com/dstout-devops/hephaestus/internal/command"
	// _ "github.com/dstout-devops/hephaestus/internal/builtins"
)

func main() {
	app := command.NewApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "hephaestus: %s\n", err)
		os.Exit(1)
	}
}
