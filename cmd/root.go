package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// RootCmd is the root command for gh-clone-org
var RootCmd = &cobra.Command{
	Use:   "gh-clone-org",
	Short: "Clone all repositories in a GitHub organization",
	Long: `Clone all repositories from a GitHub organization to a local folder.
If a repository already exists, it will update it. Repositories are cloned in parallel.`,
}

// Execute executes the root command
func Execute() error {
	return RootCmd.Execute()
}

// SetVersion sets the version info
func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
	RootCmd.Version = fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date)
}
