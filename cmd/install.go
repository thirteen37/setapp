package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <app>",
	Short: "Install an app via Setapp",
	Long:  "Opens the app's Setapp page where you can click Install (like brew install).",
	Args:  cobra.ExactArgs(1),
	RunE:  runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	if installedAppNames()[app.Name] {
		fmt.Fprintf(cmd.OutOrStdout(), "%s is already installed.\n", app.Name)
		return nil
	}

	if app.SharingURL == "" {
		return fmt.Errorf("no Setapp page URL available for %s", app.Name)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Opening Setapp page for %s...\n", app.Name)
	return execCommand("open", app.SharingURL).Run()
}
