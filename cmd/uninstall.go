package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/model"
)

var forceUninstall bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <app>",
	Short: "Uninstall a Setapp app",
	Long:  "Removes the app from /Applications/Setapp/. Setapp will detect the removal (like brew uninstall).",
	Args:  cobra.ExactArgs(1),
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "skip confirmation prompt")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	if !installedAppNames()[app.Name] {
		return fmt.Errorf("%s is not installed", app.Name)
	}

	appPath := model.AppPath(app.Name)
	out := cmd.OutOrStdout()

	if !forceUninstall {
		fmt.Fprintf(out, "Uninstall %s? This will remove %s [y/N] ", app.Name, appPath)
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Fprintln(out, "Cancelled.")
			return nil
		}
	}

	if err := removeAll(appPath); err != nil {
		return fmt.Errorf("failed to remove %s: %w", appPath, err)
	}

	fmt.Fprintf(out, "Uninstalled %s.\n", app.Name)
	return nil
}
