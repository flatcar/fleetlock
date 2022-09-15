package main

import (
	"os"

	"github.com/flatcar/fleetlock/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		os.Exit(1)
	}
}
