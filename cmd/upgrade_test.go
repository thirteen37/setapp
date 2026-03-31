package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunUpgrade(t *testing.T) {
	saveDeps(t)

	t.Run("no args opens setapp", func(t *testing.T) {
		noopExec()

		out, err := executeCommand("upgrade")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Opening Setapp to check for updates") {
			t.Error("expected update message")
		}
	})

	t.Run("with app opens page", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp: &model.App{Name: "Bartender", SharingURL: "https://setapp.com/bartender"},
		})
		noopExec()

		out, err := executeCommand("upgrade", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Opening Setapp page for Bartender") {
			t.Error("expected opening message")
		}
	})

	t.Run("no sharing URL", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp: &model.App{Name: "NoURL"},
		})

		_, err := executeCommand("upgrade", "NoURL")
		if err == nil {
			t.Fatal("expected error for missing URL")
		}
	})

	t.Run("app not found", func(t *testing.T) {
		mockOpenDB(&mockStore{findAppErr: fmt.Errorf("no app found")})

		_, err := executeCommand("upgrade", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
