package main

import (
	"fmt"
	"os"

	"github.com/flatcar-linux/fleetlock/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to execute command: %v", err)
	}
}
