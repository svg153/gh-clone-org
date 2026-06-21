package cmd

import (
	"os"
	"testing"
)

func TestValidateOrg_NonExistent(t *testing.T) {
	// This will fail because the org doesn't exist, but it tests the flow
	err := validateOrg("this-org-definitely-does-not-exist-xyz-123")
	if err == nil {
		t.Error("expected error for non-existent org")
	}
}

func TestValidateOrg_User(t *testing.T) {
	// This should fail because "dependabot[bot]" is a user, not an org
	err := validateOrg("dependabot[bot]")
	if err == nil {
		t.Error("expected error for user")
	}
}

func TestMain(m *testing.M) {
	// Ensure GH_TOKEN is set for API calls
	if os.Getenv("GH_TOKEN") == "" {
		t := &testing.T{}
		t.Log("GH_TOKEN not set — tests will fail on API calls")
	}
	m.Run()
}
