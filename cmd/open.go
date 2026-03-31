package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
	"github.com/thirteen37/setapp/internal/model"
)

var openCmd = &cobra.Command{
	Use:   "open <app>",
	Short: "Launch an installed app",
	Args:  cobra.ExactArgs(1),
	RunE:  runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	if !model.InstalledAppNames()[app.Name] {
		return fmt.Errorf("%s is not installed. Use 'setapp install %s' to install it", app.Name, app.Name)
	}

	appPath := model.AppPath(app.Name)
	return exec.Command("open", appPath).Run()
}
