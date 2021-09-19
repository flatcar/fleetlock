package cmd

import (
	"github.com/spf13/cobra"
)

var (
	group, id, url string

	Root = &cobra.Command{Use: "fleetlockctl"}
)

func init() {
	Root.PersistentFlags().StringVarP(&group, "group", "g", "default", "Fleet-Lock group")
	Root.PersistentFlags().StringVarP(&id, "id", "i", "", "Fleet-Lock instance ID (/etc/machine-id for example)")
	Root.PersistentFlags().StringVarP(&url, "url", "u", "", "Fleet-Lock endpoint URL")

	Root.AddCommand(lock)
	Root.AddCommand(unlock)
}
