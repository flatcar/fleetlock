package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/flatcar-linux/fleetlock/pkg/client"
)

var lock = &cobra.Command{
	Use: "recursive-lock",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.DefaultClient

		c, err := client.New(url, group, id, httpClient)
		if err != nil {
			return fmt.Errorf("building the client: %w", err)
		}

		if err := c.RecursiveLock(); err != nil {
			return fmt.Errorf("locking: %w", err)
		}

		return nil
	},
}
