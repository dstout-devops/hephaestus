package main

import (
	"fmt"
	"os"

	"github.com/dstout-devops/hephaestus/internal/command"
	// _ "github.com/dstout-devops/hephaestus/internal/builtins"
)

func main() {
	cmd := command.NewCommand(nil, nil, nil, nil)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "hephaestus: %s\n", err)
		os.Exit(1)
	}
}
