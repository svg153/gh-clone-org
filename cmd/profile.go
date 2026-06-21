package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ConfigProfile represents a saved configuration profile
type ConfigProfile struct {
	SkipArchived     bool     `yaml:"skip_archived"`
	SkipForks        bool     `yaml:"skip_forks"`
	Path             string   `yaml:"path"`
	RateLimitDelay   int      `yaml:"rate_limit_delay"`
	Verbose          int      `yaml:"verbose"`
	IncludePatterns  []string `yaml:"include_patterns"`
	ExcludePatterns  []string `yaml:"exclude_patterns"`
	Limit            int      `yaml:"limit"`
	DryRun           bool     `yaml:"dry_run"`
	UpdateOrgFolder  bool     `yaml:"update_org_folder"`
	DisableProtection bool    `yaml:"disable_clone_protection"`
	ServerHostSSH    string   `yaml:"server_host_ssh"`
}

// DefaultProfiles returns built-in profile configurations
func DefaultProfiles() map[string]ConfigProfile {
	return map[string]ConfigProfile{
		"full": {
			SkipArchived: false,
			SkipForks:    false,
			Verbose:      1,
		},
		"minimal": {
			SkipArchived: true,
			SkipForks:    true,
			Verbose:      0,
		},
		"no-forks": {
			SkipArchived: false,
			SkipForks:    true,
			Verbose:      1,
		},
	}
}

// LoadProfile loads a configuration profile from file
func LoadProfile(name string) (*ConfigProfile, error) {
	profiles := DefaultProfiles()
	
	// Check built-in profiles first
	if p, ok := profiles[strings.ToLower(name)]; ok {
		return &p, nil
	}

	// Try to load from file
	paths := []string{
		".gh-clone-org.yaml",
		filepath.Join(os.Getenv("HOME"), ".gh-clone-org.yaml"),
		filepath.Join(os.Getenv("HOME"), ".config", "gh-clone-org", "profiles.yaml"),
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var allProfiles map[string]ConfigProfile
		if err := yaml.Unmarshal(data, &allProfiles); err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
		}

		if p, ok := allProfiles[name]; ok {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("profile %s not found (built-in: full, minimal, no-forks)", name)
}

// ApplyProfile applies a profile's settings to the config, only overriding defaults
func ApplyProfile(cfg *config, profile *ConfigProfile) {
	// Only apply if CLI flag was not explicitly set
	// This is handled by checking if the value differs from zero/false defaults
	// and the CLI flag was not provided.
	
	// For simplicity, profile values override when explicitly set
	if profile.SkipArchived {
		cfg.skipArchived = true
	}
	if profile.SkipForks {
		cfg.skipForks = true
	}
	if profile.Path != "" {
		cfg.path = profile.Path
	}
	if profile.Verbose > 0 {
		cfg.verbose = profile.Verbose
	}
	if len(profile.IncludePatterns) > 0 {
		cfg.includePatterns = profile.IncludePatterns
	}
	if len(profile.ExcludePatterns) > 0 {
		cfg.excludePatterns = profile.ExcludePatterns
	}
	if profile.Limit > 0 {
		cfg.limit = profile.Limit
	}
	if profile.DryRun {
		cfg.dryRun = true
	}
	if profile.UpdateOrgFolder {
		cfg.updateOrgFolder = true
	}
	if profile.DisableProtection {
		cfg.disableCloneProtection = true
	}
	if profile.ServerHostSSH != "" {
		cfg.serverHostSSH = profile.ServerHostSSH
	}
}
