package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunInstall(t *testing.T) {
	saveDeps(t)

	app := &model.App{Name: "Bartender", SharingURL: "https://setapp.com/bartender"}

	t.Run("already installed", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{"Bartender": true})

		out, err := executeCommand("install", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "already installed") {
			t.Error("expected already installed message")
		}
	})

	t.Run("opens setapp page", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: app})
		mockInstalled(map[string]bool{})
		noopExec()

		out, err := executeCommand("install", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Opening Setapp page") {
			t.Error("expected opening message")
		}
	})

	t.Run("no sharing URL", func(t *testing.T) {
		mockOpenDB(&mockStore{findApp: &model.App{Name: "NoURL"}})
		mockInstalled(map[string]bool{})

		_, err := executeCommand("install", "NoURL")
		if err == nil {
			t.Fatal("expected error for missing URL")
		}
	})

	t.Run("app not found", func(t *testing.T) {
		mockOpenDB(&mockStore{findAppErr: fmt.Errorf("no app found")})

		_, err := executeCommand("install", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
