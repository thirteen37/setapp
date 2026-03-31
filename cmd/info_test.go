package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunInfo(t *testing.T) {
	saveDeps(t)

	app := &model.App{
		PK:      1,
		Name:    "Bartender",
		Vendor:  "Surtees Studios",
		Version: "4.0",
		Tagline: "Organize menu bar",
		Size:    52428800, // 50 MB
	}

	t.Run("text output", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp:       app,
			appCategories: []string{"Utilities"},
		})
		mockInstalled(map[string]bool{"Bartender": true})
		jsonOutput = false

		out, err := executeCommand("info", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Name:         Bartender") {
			t.Error("expected name in output")
		}
		if !strings.Contains(out, "installed") {
			t.Error("expected installed status")
		}
		if !strings.Contains(out, "Utilities") {
			t.Error("expected categories")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findAppErr: fmt.Errorf("no app found matching %q", "nope"),
		})
		jsonOutput = false

		_, err := executeCommand("info", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("json output", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp:       app,
			appCategories: []string{"Utilities"},
		})
		mockInstalled(map[string]bool{})
		jsonOutput = true

		out, err := executeCommand("info", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"name": "Bartender"`) {
			t.Error("expected JSON output")
		}
	})
}
