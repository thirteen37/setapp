package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
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
	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	apps, err := d.AllApps()
	if err != nil {
		return err
	}

	installed := model.InstalledAppNames()
	var result []model.App
	for i := range apps {
		if installed[apps[i].Name] {
			apps[i].Installed = true
			result = append(result, apps[i])
		}
	}

	if jsonOutput {
		printJSON(result)
		return nil
	}

	if len(result) == 0 {
		fmt.Println("No Setapp apps installed.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVENDOR\tVERSION")
	for _, a := range result {
		fmt.Fprintf(w, "%s\t%s\t%s\n", a.Name, a.Vendor, a.Version)
	}
	return w.Flush()
}
