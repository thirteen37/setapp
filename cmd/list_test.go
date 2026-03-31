package cmd

import (
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunList(t *testing.T) {
	saveDeps(t)

	apps := []model.App{
		{Name: "Bartender", Vendor: "Surtees", Version: "4.0"},
		{Name: "CleanMyMac", Vendor: "MacPaw", Version: "5.0"},
	}

	t.Run("with installed apps", func(t *testing.T) {
		mockOpenDB(&mockStore{allApps: apps})
		mockInstalled(map[string]bool{"CleanMyMac": true})
		jsonOutput = false

		out, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "CleanMyMac") {
			t.Error("expected CleanMyMac in output")
		}
		if strings.Contains(out, "Bartender") {
			t.Error("Bartender should not appear (not installed)")
		}
	})

	t.Run("no installed apps", func(t *testing.T) {
		mockOpenDB(&mockStore{allApps: apps})
		mockInstalled(map[string]bool{})
		jsonOutput = false

		out, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "No Setapp apps installed") {
			t.Error("expected empty message")
		}
	})

	t.Run("json output", func(t *testing.T) {
		mockOpenDB(&mockStore{allApps: apps})
		mockInstalled(map[string]bool{"CleanMyMac": true})
		jsonOutput = true

		out, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"name": "CleanMyMac"`) {
			t.Error("expected JSON output with CleanMyMac")
		}
	})
}
