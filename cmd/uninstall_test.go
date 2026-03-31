package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunUninstall(t *testing.T) {
	saveDeps(t)

	app := &model.App{Name: "Bartender"}

	t.Run("not installed", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{})

		_, err := executeCommand("uninstall", "Bartender")
		if err == nil {
			t.Fatal("expected error for uninstalled app")
		}
	})

	t.Run("force uninstall", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{"Bartender": true})
		var removedPath string
		removeAll = func(path string) error {
			removedPath = path
			return nil
		}

		out, err := executeCommand("uninstall", "--force", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Uninstalled Bartender") {
			t.Error("expected uninstalled message")
		}
		if removedPath != "/Applications/Setapp/Bartender.app" {
			t.Errorf("removed wrong path: %s", removedPath)
		}
	})

	t.Run("force uninstall remove error", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{"Bartender": true})
		removeAll = func(path string) error {
			return fmt.Errorf("permission denied")
		}

		_, err := executeCommand("uninstall", "--force", "Bartender")
		if err == nil {
			t.Fatal("expected error from removeAll")
		}
	})

	t.Run("app not found", func(t *testing.T) {
		mockOpenDB(&mockStore{findAppErr: fmt.Errorf("no app found")})

		_, err := executeCommand("uninstall", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
