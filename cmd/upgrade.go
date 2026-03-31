package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [app]",
	Short: "Check for updates",
	Long:  "Opens Setapp to check for updates. Optionally specify an app to open its page (like brew upgrade).",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUpgrade,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Opening Setapp to check for updates...")
		return execCommand("open", "setappDiscovery://").Run()
	}

	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	if app.SharingURL == "" {
		return fmt.Errorf("no Setapp page URL available for %s", app.Name)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Opening Setapp page for %s...\n", app.Name)
	return execCommand("open", app.SharingURL).Run()
}
