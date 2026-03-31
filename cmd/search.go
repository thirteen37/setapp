package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
	"github.com/thirteen37/setapp/internal/model"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search available apps",
	Long:  "Search all Setapp apps by name, keyword, tagline, or vendor (like brew search).",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	apps, err := d.SearchApps(args[0])
	if err != nil {
		return err
	}

	installed := model.InstalledAppNames()
	for i := range apps {
		apps[i].Installed = installed[apps[i].Name]
	}

	if jsonOutput {
		printJSON(apps)
		return nil
	}

	if len(apps) == 0 {
		fmt.Printf("No apps found matching %q.\n", args[0])
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVENDOR\tTAGLINE\tSTATUS")
	for _, a := range apps {
		status := ""
		if a.Installed {
			status = "*"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.Name, a.Vendor, a.Tagline, status)
	}
	return w.Flush()
}
