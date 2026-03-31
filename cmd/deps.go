package cmd

import (
	"os"
	"os/exec"

	"github.com/thirteen37/setapp/internal/db"
	"github.com/thirteen37/setapp/internal/model"
)

// Replaceable dependencies for testing.
var (
	openDB            = func() (db.Store, error) { return db.Open() }
	installedAppNames = model.InstalledAppNames
	execCommand       = exec.Command
	removeAll         = os.RemoveAll
)
