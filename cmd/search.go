package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
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
	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	apps, err := d.SearchApps(args[0])
	if err != nil {
		return err
	}

	installed := installedAppNames()
	for i := range apps {
		apps[i].Installed = installed[apps[i].Name]
	}

	if jsonOutput {
		printJSON(cmd, apps)
		return nil
	}

	out := cmd.OutOrStdout()
	if len(apps) == 0 {
		fmt.Fprintf(out, "No apps found matching %q.\n", args[0])
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
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
