package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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
	d, err := openDB()
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
	app.Installed = installedAppNames()[app.Name]

	if jsonOutput {
		printJSON(cmd, app)
		return nil
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Name:         %s\n", app.Name)
	fmt.Fprintf(out, "Vendor:       %s\n", app.Vendor)
	if app.Version != "" {
		fmt.Fprintf(out, "Version:      %s\n", app.Version)
	}
	fmt.Fprintf(out, "Status:       %s\n", app.StatusString())
	if app.Tagline != "" {
		fmt.Fprintf(out, "Tagline:      %s\n", app.Tagline)
	}
	if len(app.Categories) > 0 {
		fmt.Fprintf(out, "Categories:   %s\n", strings.Join(app.Categories, ", "))
	}
	if app.Size > 0 {
		fmt.Fprintf(out, "Size:         %s\n", model.FormatSize(app.Size))
	}
	if app.MinOS != "" {
		fmt.Fprintf(out, "Min macOS:    %s\n", app.MinOS)
	}
	if app.MarketingURL != "" {
		fmt.Fprintf(out, "Website:      %s\n", app.MarketingURL)
	}
	if app.SharingURL != "" {
		fmt.Fprintf(out, "Setapp page:  %s\n", app.SharingURL)
	}
	if app.Keywords != "" {
		fmt.Fprintf(out, "Keywords:     %s\n", app.Keywords)
	}
	if !app.FirstReleaseTime().IsZero() {
		fmt.Fprintf(out, "First release: %s\n", app.FirstReleaseTime().Format("2006-01-02"))
	}
	if !app.LastReleaseTime().IsZero() {
		fmt.Fprintf(out, "Last release:  %s\n", app.LastReleaseTime().Format("2006-01-02"))
	}
	if app.Description != "" {
		fmt.Fprintf(out, "\n%s\n", app.Description)
	}

	return nil
}
