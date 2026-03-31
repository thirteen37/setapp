package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/model"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed apps",
	Long:  "List apps installed via Setapp (like brew list).",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	apps, err := d.AllApps()
	if err != nil {
		return err
	}

	installed := installedAppNames()
	var result []model.App
	for i := range apps {
		if installed[apps[i].Name] {
			apps[i].Installed = true
			result = append(result, apps[i])
		}
	}

	if jsonOutput {
		printJSON(cmd, result)
		return nil
	}

	out := cmd.OutOrStdout()
	if len(result) == 0 {
		fmt.Fprintln(out, "No Setapp apps installed.")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVENDOR\tVERSION")
	for _, a := range result {
		fmt.Fprintf(w, "%s\t%s\t%s\n", a.Name, a.Vendor, a.Version)
	}
	return w.Flush()
}
