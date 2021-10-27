// Package cmd implements fleetlockctl CLI.
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// Command returns CLI command to be executed.
func Command() *cobra.Command {
	cli := &cobra.Command{Use: "fleetlockctl"}

	var group, id, url string

	cli.PersistentFlags().StringVarP(&group, "group", "g", "default", "FleetLock group")
	cli.PersistentFlags().StringVarP(&id, "id", "i", "", "FleetLock instance ID (e.g. content of /etc/machine-id file)")
	cli.PersistentFlags().StringVarP(&url, "url", "u", "", "FleetLock endpoint URL")

	cli.AddCommand(lock(&group, &id, &url))
	cli.AddCommand(unlock(&group, &id, &url))

	return cli
}

// machineID is a helper to return unique ID
// of the machine.
func machineID() (*string, error) {
	id, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		return nil, fmt.Errorf("reading machine ID from file: %w", err)
	}

	sid := string(id)

	return &sid, nil
}
