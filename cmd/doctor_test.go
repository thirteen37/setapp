package cmd

import (
	"os/exec"
	"strings"
	"testing"
)

func TestParsePlistOutput(t *testing.T) {
	input := `account=user@example.com
subscriptionState=active
subscriptionStartDate=2024-01-15
subscriptionExpirationDate=2025-01-15
gracePeriodExpirationDate=2025-02-15`

	info := &doctorInfo{}
	parsePlistOutput(input, info)

	if info.Account != "user@example.com" {
		t.Errorf("Account = %q, want %q", info.Account, "user@example.com")
	}
	if info.Subscription != "active" {
		t.Errorf("Subscription = %q, want %q", info.Subscription, "active")
	}
	if info.Since != "2024-01-15" {
		t.Errorf("Since = %q, want %q", info.Since, "2024-01-15")
	}
	if info.Expires != "2025-01-15" {
		t.Errorf("Expires = %q, want %q", info.Expires, "2025-01-15")
	}
	if info.GracePeriod != "2025-02-15" {
		t.Errorf("GracePeriod = %q, want %q", info.GracePeriod, "2025-02-15")
	}
}

func TestParsePlistOutputEmpty(t *testing.T) {
	info := &doctorInfo{}
	parsePlistOutput("", info)

	if info.Account != "" || info.Subscription != "" {
		t.Error("expected empty fields for empty input")
	}
}

func TestRunDoctor(t *testing.T) {
	saveDeps(t)

	t.Run("text output", func(t *testing.T) {
		mockInstalled(map[string]bool{"App1": true, "App2": true})
		// Mock exec to fail (no Setapp defaults available)
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		}
		jsonOutput = false

		out, err := executeCommand("doctor")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "Installed apps: 2") {
			t.Errorf("expected installed count, got: %s", out)
		}
	})

	t.Run("json output", func(t *testing.T) {
		mockInstalled(map[string]bool{"App1": true})
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		}
		jsonOutput = true

		out, err := executeCommand("doctor")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"installed_count": 1`) {
			t.Errorf("expected JSON output, got: %s", out)
		}
	})
}
