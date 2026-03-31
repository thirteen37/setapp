package cmd

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/thirteen37/setapp/internal/db"
	"github.com/thirteen37/setapp/internal/model"
)

type mockStore struct {
	allApps        []model.App
	allAppsErr     error
	searchApps     []model.App
	searchAppsErr  error
	findApp        *model.App
	findAppErr     error
	loadCatsErr    error
	appCategories  []string
	appCatsErr     error
	allCategories  []model.Category
	allCatsErr     error
	appsByCategory []model.App
	appsByCatErr   error
}

func (m *mockStore) AllApps() ([]model.App, error)                { return m.allApps, m.allAppsErr }
func (m *mockStore) SearchApps(q string) ([]model.App, error)     { return m.searchApps, m.searchAppsErr }
func (m *mockStore) FindApp(name string) (*model.App, error)      { return m.findApp, m.findAppErr }
func (m *mockStore) LoadCategories(apps []model.App) error        { return m.loadCatsErr }
func (m *mockStore) AppCategories(pk int) ([]string, error)       { return m.appCategories, m.appCatsErr }
func (m *mockStore) AllCategories() ([]model.Category, error)     { return m.allCategories, m.allCatsErr }
func (m *mockStore) AppsByCategory(n string) ([]model.App, error) { return m.appsByCategory, m.appsByCatErr }
func (m *mockStore) Close() error                                 { return nil }

// saveDeps saves current dependency vars and returns a restore function.
func saveDeps(t *testing.T) {
	t.Helper()
	origOpenDB := openDB
	origInstalled := installedAppNames
	origExec := execCommand
	origRemove := removeAll
	origJSON := jsonOutput
	t.Cleanup(func() {
		openDB = origOpenDB
		installedAppNames = origInstalled
		execCommand = origExec
		removeAll = origRemove
		jsonOutput = origJSON
	})
}

// mockOpenDB sets openDB to return the given store.
func mockOpenDB(store db.Store) {
	openDB = func() (db.Store, error) { return store, nil }
}

// mockInstalled sets installedAppNames to return the given map.
func mockInstalled(names map[string]bool) {
	installedAppNames = func() map[string]bool { return names }
}

// noopExec sets execCommand to return a no-op command (echo).
func noopExec() {
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}
}

// executeCommand runs a cobra command with the given args and captures stdout.
func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}
