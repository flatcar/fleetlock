package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/flatcar-linux/fleetlock/pkg/client"
)

func lock(group, id, url *string) *cobra.Command {
	return &cobra.Command{
		Use:   "recursive-lock",
		Short: "Try to reserve (lock) a slot for rebooting",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := http.DefaultClient

			c, err := client.New(&client.Config{
				URL:   *url,
				Group: *group,
				ID:    *id,
				HTTP:  httpClient,
			})
			if err != nil {
				return fmt.Errorf("building the client: %w", err)
			}

			if err := c.RecursiveLock(context.Background()); err != nil {
				return fmt.Errorf("locking: %w", err)
			}

			return nil
		},
	}
}
