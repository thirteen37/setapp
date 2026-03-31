package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Show account and subscription status",
	Long:  "Display Setapp account info, subscription status, and app counts (like brew doctor).",
	RunE:  runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

type doctorInfo struct {
	Account        string `json:"account"`
	Subscription   string `json:"subscription,omitempty"`
	Since          string `json:"since,omitempty"`
	Expires        string `json:"expires,omitempty"`
	GracePeriod    string `json:"grace_period,omitempty"`
	InstalledCount int    `json:"installed_count"`
}

const plistScript = `
import plistlib, sys, os
tmpf = sys.argv[1]
with open(tmpf, 'rb') as f:
    d = plistlib.load(f)
print("account=" + d.get("CurrentUserAccount", ""))
kc = d.get("known_customers", [])
if kc:
    c = kc[0]
    for k in ["subscriptionState", "subscriptionStartDate", "subscriptionExpirationDate", "gracePeriodExpirationDate"]:
        v = c.get(k, "")
        if hasattr(v, "strftime"):
            v = v.strftime("%Y-%m-%d")
        print(k + "=" + str(v))
`

func runDoctor(cmd *cobra.Command, args []string) error {
	info := doctorInfo{
		InstalledCount: len(installedAppNames()),
	}

	// Export plist and parse with Python plistlib (handles nested binary data)
	plistData, err := execCommand("defaults", "export", "com.setapp.DesktopClient", "-").Output()
	if err == nil {
		tmpFile, tmpErr := os.CreateTemp("", "setapp-*.plist")
		if tmpErr == nil {
			tmpFile.Write(plistData)
			tmpFile.Close()
			out, pyErr := execCommand("python3", "-c", plistScript, tmpFile.Name()).Output()
			os.Remove(tmpFile.Name())
			if pyErr == nil {
				parsePlistOutput(string(out), &info)
			}
		}
	}

	// Fallback: read account directly
	if info.Account == "" {
		out, err := execCommand("defaults", "read", "com.setapp.DesktopClient", "CurrentUserAccount").Output()
		if err == nil {
			info.Account = strings.TrimSpace(string(out))
		}
	}

	if jsonOutput {
		printJSON(cmd, info)
		return nil
	}

	w := cmd.OutOrStdout()
	if info.Account != "" {
		fmt.Fprintf(w, "Account:        %s\n", info.Account)
	}
	if info.Subscription != "" {
		fmt.Fprintf(w, "Subscription:   %s\n", info.Subscription)
	}
	if info.Since != "" {
		fmt.Fprintf(w, "Since:          %s\n", info.Since)
	}
	if info.Expires != "" {
		fmt.Fprintf(w, "Expires:        %s\n", info.Expires)
	}
	if info.GracePeriod != "" {
		fmt.Fprintf(w, "Grace period:   %s\n", info.GracePeriod)
	}
	fmt.Fprintf(w, "Installed apps: %d\n", info.InstalledCount)

	return nil
}

func parsePlistOutput(output string, info *doctorInfo) {
	for _, line := range strings.Split(output, "\n") {
		k, v, ok := strings.Cut(line, "=")
		if !ok || v == "" {
			continue
		}
		switch k {
		case "account":
			info.Account = v
		case "subscriptionState":
			info.Subscription = v
		case "subscriptionStartDate":
			info.Since = v[:min(10, len(v))]
		case "subscriptionExpirationDate":
			info.Expires = v[:min(10, len(v))]
		case "gracePeriodExpirationDate":
			info.GracePeriod = v[:min(10, len(v))]
		}
	}
}
