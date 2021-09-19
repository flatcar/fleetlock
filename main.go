package main

import (
	"fmt"
	"os"

	"github.com/flatcar-linux/fleetlock/cmd"
)

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to execute command: %v", err)
	}
}
