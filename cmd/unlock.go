package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/flatcar-linux/fleetlock/pkg/client"
)

func unlock(group, id, url *string) *cobra.Command {
	return &cobra.Command{
		Use:   "unlock-if-held",
		Short: "Try to release (unlock) a slot that it was previously holding",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == nil {
				var err error
				id, err = machineID()
				if err != nil {
					return fmt.Errorf("getting machine ID: %w", err)
				}
			}

			c, err := client.New(&client.Config{
				ID:    *id,
				Group: *group,
				URL:   *url,
			})
			if err != nil {
				return fmt.Errorf("building the client: %w", err)
			}

			if err := c.UnlockIfHeld(context.Background()); err != nil {
				return fmt.Errorf("unlocking: %w", err)
			}

			return nil
		},
	}
}
