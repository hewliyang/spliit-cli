package main

import (
	"os"

	"github.com/m1a1/spliit-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
