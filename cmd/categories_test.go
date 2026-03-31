package cmd

import (
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunCategories(t *testing.T) {
	saveDeps(t)

	cats := []model.Category{
		{Name: "Productivity", Description: "Productivity apps", Position: 1},
		{Name: "Utilities", Description: "Utility apps", Position: 2},
	}

	t.Run("table output", func(t *testing.T) {
		mockOpenDB(&mockStore{allCategories: cats})
		jsonOutput = false

		out, err := executeCommand("categories")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Productivity") || !strings.Contains(out, "Utilities") {
			t.Error("expected both categories in output")
		}
	})

	t.Run("json output", func(t *testing.T) {
		mockOpenDB(&mockStore{allCategories: cats})
		jsonOutput = true

		out, err := executeCommand("categories")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"name": "Productivity"`) {
			t.Error("expected JSON output")
		}
	})
}
