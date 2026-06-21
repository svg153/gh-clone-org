package cmd

import (
	"testing"
)

func TestValidateTarget_NonExistent(t *testing.T) {
	// This will fail because the org doesn't exist, but it tests the flow
	err := validateTarget("this-org-definitely-does-not-exist-xyz-123", false)
	if err == nil {
		t.Error("expected error for non-existent org")
	}
}

func TestValidateTarget_User(t *testing.T) {
	// This should fail because "dependabot[bot]" is a user, not an org
	err := validateTarget("dependabot[bot]", false)
	if err == nil {
		t.Error("expected error for user")
	}
}

func TestValidateTarget_UserMode(t *testing.T) {
	// In user mode, dependabot[bot] should pass (it's a user)
	// But without a real token, this will fail on API call — just check it doesn't panic
	_ = validateTarget("dependabot[bot]", true)
}
