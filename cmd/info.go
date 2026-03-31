package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
	"github.com/thirteen37/setapp/internal/model"
)

var infoCmd = &cobra.Command{
	Use:   "info <app>",
	Short: "Show detailed app info",
	Long:  "Show detailed information about a Setapp app (like brew info).",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	app, err := d.FindApp(args[0])
	if err != nil {
		return err
	}

	cats, err := d.AppCategories(app.PK)
	if err != nil {
		return err
	}
	app.Categories = cats
	app.Installed = model.InstalledAppNames()[app.Name]

	if jsonOutput {
		printJSON(app)
		return nil
	}

	fmt.Printf("Name:         %s\n", app.Name)
	fmt.Printf("Vendor:       %s\n", app.Vendor)
	if app.Version != "" {
		fmt.Printf("Version:      %s\n", app.Version)
	}
	fmt.Printf("Status:       %s\n", app.StatusString())
	if app.Tagline != "" {
		fmt.Printf("Tagline:      %s\n", app.Tagline)
	}
	if len(app.Categories) > 0 {
		fmt.Printf("Categories:   %s\n", strings.Join(app.Categories, ", "))
	}
	if app.Size > 0 {
		fmt.Printf("Size:         %s\n", model.FormatSize(app.Size))
	}
	if app.MinOS != "" {
		fmt.Printf("Min macOS:    %s\n", app.MinOS)
	}
	if app.MarketingURL != "" {
		fmt.Printf("Website:      %s\n", app.MarketingURL)
	}
	if app.SharingURL != "" {
		fmt.Printf("Setapp page:  %s\n", app.SharingURL)
	}
	if app.Keywords != "" {
		fmt.Printf("Keywords:     %s\n", app.Keywords)
	}
	if !app.FirstReleaseTime().IsZero() {
		fmt.Printf("First release: %s\n", app.FirstReleaseTime().Format("2006-01-02"))
	}
	if !app.LastReleaseTime().IsZero() {
		fmt.Printf("Last release:  %s\n", app.LastReleaseTime().Format("2006-01-02"))
	}
	if app.Description != "" {
		fmt.Printf("\n%s\n", app.Description)
	}

	return nil
}
