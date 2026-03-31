package cmd

import (
	"fmt"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunOpen(t *testing.T) {
	saveDeps(t)

	app := &model.App{Name: "Bartender"}

	t.Run("installed app", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{"Bartender": true})
		noopExec()

		_, err := executeCommand("open", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("not installed", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{})

		_, err := executeCommand("open", "Bartender")
		if err == nil {
			t.Fatal("expected error for uninstalled app")
		}
	})

	t.Run("app not found", func(t *testing.T) {
		mockOpenDB(&mockStore{findAppErr: fmt.Errorf("no app found")})

		_, err := executeCommand("open", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
