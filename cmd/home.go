package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
)

var homeCmd = &cobra.Command{
	Use:   "home [app]",
	Short: "Open Setapp or an app's website",
	Long:  "Opens the Setapp UI, or an app's vendor website if specified (like brew home).",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHome,
}

func init() {
	rootCmd.AddCommand(homeCmd)
}

func runHome(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return exec.Command("open", "setappDiscovery://").Run()
	}

	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	if app.MarketingURL == "" {
		return fmt.Errorf("no website URL available for %s", app.Name)
	}

	fmt.Printf("Opening %s website...\n", app.Name)
	return exec.Command("open", app.MarketingURL).Run()
}
