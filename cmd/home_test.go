package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thirteen37/setapp/internal/model"
)

func TestRunHome(t *testing.T) {
	saveDeps(t)

	t.Run("no args opens setapp", func(t *testing.T) {
		noopExec()

		_, err := executeCommand("home")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with app opens website", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp: &model.App{Name: "Bartender", MarketingURL: "https://bartender.app"},
		})
		noopExec()

		out, err := executeCommand("home", "Bartender")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Opening Bartender website") {
			t.Error("expected opening message")
		}
	})

	t.Run("no marketing URL", func(t *testing.T) {
		mockOpenDB(&mockStore{
			findApp: &model.App{Name: "NoURL"},
		})

		_, err := executeCommand("home", "NoURL")
		if err == nil {
			t.Fatal("expected error for missing URL")
		}
	})

	t.Run("app not found", func(t *testing.T) {
		mockOpenDB(&mockStore{findAppErr: fmt.Errorf("no app found")})

		_, err := executeCommand("home", "nope")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
