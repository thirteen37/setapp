package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SetappDir = "/Applications/Setapp"
	// Core Data epoch: 2001-01-01 00:00:00 UTC
	coreDataEpochOffset = 978307200
)

type App struct {
	PK               int      `json:"-"`
	Identifier       int      `json:"identifier"`
	Name             string   `json:"name"`
	BundleIdentifier string   `json:"bundle_identifier,omitempty"`
	Vendor           string   `json:"vendor"`
	Tagline          string   `json:"tagline,omitempty"`
	Description      string   `json:"description,omitempty"`
	MarketingURL     string   `json:"marketing_url,omitempty"`
	SharingURL       string   `json:"sharing_url,omitempty"`
	Size             int64    `json:"size,omitempty"`
	Version          string   `json:"version,omitempty"`
	MinOS            string   `json:"min_os,omitempty"`
	URLScheme        string   `json:"url_scheme,omitempty"`
	Keywords         string   `json:"keywords,omitempty"`
	FirstRelease     *float64 `json:"-"`
	LastRelease      *float64 `json:"-"`
	Categories       []string `json:"categories,omitempty"`
	Installed        bool     `json:"installed"`
}

func (a App) FirstReleaseTime() time.Time {
	if a.FirstRelease == nil {
		return time.Time{}
	}
	return time.Unix(int64(*a.FirstRelease)+coreDataEpochOffset, 0)
}

func (a App) LastReleaseTime() time.Time {
	if a.LastRelease == nil {
		return time.Time{}
	}
	return time.Unix(int64(*a.LastRelease)+coreDataEpochOffset, 0)
}

func (a App) StatusString() string {
	if a.Installed {
		return "installed"
	}
	return "available"
}

type Category struct {
	PK          int    `json:"-"`
	Identifier  int    `json:"identifier"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

func FormatSize(bytes int64) string {
	if bytes <= 0 {
		return ""
	}
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func InstalledAppNames() map[string]bool {
	installed := make(map[string]bool)
	entries, err := os.ReadDir(SetappDir)
	if err != nil {
		return installed
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".app") {
			installed[strings.TrimSuffix(name, ".app")] = true
		}
	}
	return installed
}

func AppPath(name string) string {
	return filepath.Join(SetappDir, name+".app")
}
