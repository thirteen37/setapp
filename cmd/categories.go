package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/thirteen37/setapp/internal/db"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List app categories",
	RunE:  runCategories,
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}

func runCategories(cmd *cobra.Command, args []string) error {
	d, err := db.Open()
	if err != nil {
		return err
	}
	defer d.Close()

	cats, err := d.AllCategories()
	if err != nil {
		return err
	}

	if jsonOutput {
		printJSON(cats)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION")
	for _, c := range cats {
		fmt.Fprintf(w, "%s\t%s\n", c.Name, c.Description)
	}
	return w.Flush()
}
