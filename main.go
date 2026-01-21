package main

import (
	"os"

	"github.com/hewliyang/spliit-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
