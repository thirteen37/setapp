package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
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
	d, err := openDB()
	if err != nil {
		return err
	}
	defer d.Close()

	cats, err := d.AllCategories()
	if err != nil {
		return err
	}

	if jsonOutput {
		printJSON(cmd, cats)
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION")
	for _, c := range cats {
		fmt.Fprintf(w, "%s\t%s\n", c.Name, c.Description)
	}
	return w.Flush()
}
