package main

import (
	"github.com/dstout-devops/hephaestus/internal/command"
	// _ "github.com/dstout-devops/hephaestus/internal/builtins"
)

func main() {
	app := command.NewApp()
	_ = app
	/*
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintf(os.Stderr, "certgopher: %s\n", err)
			os.Exit(1)
		}
	*/
}
