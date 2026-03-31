package cmd

import (
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunSearch(t *testing.T) {
	saveDeps(t)

	t.Run("with results", func(t *testing.T) {
		mockOpenDB(&mockStore{
			searchApps: []model.App{
				{Name: "Bartender", Vendor: "Surtees", Tagline: "Organize menu bar"},
			},
		})
		mockInstalled(map[string]bool{"Bartender": true})
		jsonOutput = false

		out, err := executeCommand("search", "bar")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Bartender") {
			t.Error("expected Bartender in output")
		}
		if !strings.Contains(out, "*") {
			t.Error("expected installed marker")
		}
	})

	t.Run("no results", func(t *testing.T) {
		mockOpenDB(&mockStore{})
		mockInstalled(map[string]bool{})
		jsonOutput = false

		out, err := executeCommand("search", "nonexistent")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "No apps found") {
			t.Error("expected no apps message")
		}
	})

	t.Run("json output", func(t *testing.T) {
		mockOpenDB(&mockStore{
			searchApps: []model.App{
				{Name: "Bartender", Vendor: "Surtees"},
			},
		})
		mockInstalled(map[string]bool{})
		jsonOutput = true

		out, err := executeCommand("search", "bar")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"name": "Bartender"`) {
			t.Error("expected JSON output")
		}
	})
}
