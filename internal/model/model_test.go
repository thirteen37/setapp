package model

import (
	"testing"
	"time"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, ""},
		{-1, ""},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
	}
	for _, tt := range tests {
		got := FormatSize(tt.bytes)
		if got != tt.want {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestFirstReleaseTime(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a := App{}
		if !a.FirstReleaseTime().IsZero() {
			t.Error("expected zero time for nil FirstRelease")
		}
	})
	t.Run("valid", func(t *testing.T) {
		// Core Data epoch: 2001-01-01 00:00:00 UTC
		// 0.0 should map to 2001-01-01
		v := 0.0
		a := App{FirstRelease: &v}
		got := a.FirstReleaseTime()
		want := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("FirstReleaseTime() = %v, want %v", got, want)
		}
	})
}

func TestLastReleaseTime(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a := App{}
		if !a.LastReleaseTime().IsZero() {
			t.Error("expected zero time for nil LastRelease")
		}
	})
	t.Run("valid", func(t *testing.T) {
		v := 0.0
		a := App{LastRelease: &v}
		got := a.LastReleaseTime()
		want := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("LastReleaseTime() = %v, want %v", got, want)
		}
	})
}

func TestStatusString(t *testing.T) {
	if got := (App{Installed: true}).StatusString(); got != "installed" {
		t.Errorf("StatusString() = %q, want %q", got, "installed")
	}
	if got := (App{Installed: false}).StatusString(); got != "available" {
		t.Errorf("StatusString() = %q, want %q", got, "available")
	}
}

func TestAppPath(t *testing.T) {
	got := AppPath("CleanMyMac")
	want := "/Applications/Setapp/CleanMyMac.app"
	if got != want {
		t.Errorf("AppPath(%q) = %q, want %q", "CleanMyMac", got, want)
	}
}
