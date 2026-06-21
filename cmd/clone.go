package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone-org [org]",
	Short: "Clone all repositories in a GitHub organization",
	Long: `Clone all repositories from a GitHub organization to a local folder.
If a repository already exists, it will update it. Repositories are cloned in parallel.

Examples:
  gh clone-org github
  gh clone-org github -p ~/github
  gh clone-org github -s github.com-company`,
	Args: cobra.MaximumNArgs(1),
	RunE: runClone,
}

func init() {
	// Replace RootCmd with cloneCmd so gh-clone-org IS the clone command
	RootCmd = cloneCmd
	RootCmd.Short = "Clone all repositories in a GitHub organization"
	RootCmd.Long = `Clone all repositories from a GitHub organization to a local folder.
If a repository already exists, it will update it. Repositories are cloned in parallel.`

	RootCmd.Flags().StringP("org", "o", "", "GitHub organization (positional arg alias)")
	RootCmd.Flags().StringP("path", "p", "", "Path to clone repositories (default: current directory)")
	RootCmd.Flags().Bool("update-org-folder", false, "Update existing repositories and clone new ones")
	RootCmd.Flags().Bool("disable-clone-protection", false, "Disable GIT_CLONE_PROTECTION_ACTIVE")
	RootCmd.Flags().StringP("server-host-ssh", "s", "github.com", "SSH server host for multi-account setups")
	
	// Issue #3: Filters
	RootCmd.Flags().Bool("skip-archived", false, "Skip archived repositories")
	RootCmd.Flags().Bool("skip-forks", false, "Skip forked repositories")
	RootCmd.Flags().StringArray("include-pattern", []string{}, "Glob pattern to include repos (can be repeated)")
	RootCmd.Flags().StringArray("exclude-pattern", []string{}, "Glob pattern to exclude repos (can be repeated)")
	RootCmd.Flags().Int("limit", 0, "Maximum number of repositories to clone (0 = unlimited)")
	RootCmd.Flags().Bool("dry-run", false, "List repos that would be cloned without actually cloning")
}

type config struct {
	organization           string
	path                   string
	updateOrgFolder        bool
	disableCloneProtection bool
	serverHostSSH          string
	// Issue #3: Filters
	skipArchived         bool
	skipForks            bool
	includePatterns      []string
	excludePatterns      []string
	limit                int
	// Issue #5: Dry-run
	dryRun bool
}

func runClone(cmd *cobra.Command, args []string) error {
	cfg := &config{}

	// Parse organization from positional arg or flag
	if len(args) > 0 {
		cfg.organization = args[0]
	}
	orgFlag, _ := cmd.Flags().GetString("org")
	if orgFlag != "" {
		cfg.organization = orgFlag
	}

	if cfg.organization == "" {
		return fmt.Errorf("organization is required")
	}

	// Validate organization type
	if err := validateOrg(cfg.organization); err != nil {
		return err
	}

	// Parse path
	pathFlag, _ := cmd.Flags().GetString("path")
	if pathFlag != "" {
		cfg.path = pathFlag
	}
	if cfg.path == "" {
		var err error
		cfg.path, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Parse other flags
	cfg.updateOrgFolder, _ = cmd.Flags().GetBool("update-org-folder")
	cfg.disableCloneProtection, _ = cmd.Flags().GetBool("disable-clone-protection")
	cfg.serverHostSSH, _ = cmd.Flags().GetString("server-host-ssh")
	
	// Issue #3: Parse filters
	cfg.skipArchived, _ = cmd.Flags().GetBool("skip-archived")
	cfg.skipForks, _ = cmd.Flags().GetBool("skip-forks")
	cfg.includePatterns, _ = cmd.Flags().GetStringArray("include-pattern")
	cfg.excludePatterns, _ = cmd.Flags().GetStringArray("exclude-pattern")
	cfg.limit, _ = cmd.Flags().GetInt("limit")
	cfg.dryRun, _ = cmd.Flags().GetBool("dry-run")

	return cloneOrg(cfg)
}

// validateOrg checks that the given name is an organization, not a user
func validateOrg(name string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Check if it's a user
	var user struct {
		Type string `json:"type"`
	}
	if err := client.Get(fmt.Sprintf("users/%s", name), &user); err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}
	if user.Type == "User" {
		return fmt.Errorf("this extension only works with organizations, not users")
	}

	// Check org exists
	var org struct {
		Login string `json:"login"`
	}
	if err := client.Get(fmt.Sprintf("orgs/%s", name), &org); err != nil {
		return fmt.Errorf("organization %s does not exist", name)
	}

	return nil
}

// cloneOrg clones all repositories in the organization
func cloneOrg(cfg *config) error {
	// Get SSH URLs for all repos
	repos, err := getRepoSSHURLs(cfg.organization, cfg.serverHostSSH, cfg)
	if err != nil {
		return fmt.Errorf("failed to get repositories: %w", err)
	}

	fmt.Printf("Found %d repositories in %s\n", len(repos), cfg.organization)

	// Issue #5: Dry-run mode
	if cfg.dryRun {
		fmt.Printf("\nWould clone %d repositories:\n", len(repos))
		for i, url := range repos {
			name := strings.TrimSuffix(filepath.Base(url), ".git")
			fmt.Printf("  %d. %s (%s)\n", i+1, name, url)
		}
		fmt.Printf("\nTotal: %d repositories (none were cloned)\n", len(repos))
		return nil
	}

	// Ensure target directory exists
	if err := os.MkdirAll(cfg.path, 0755); err != nil {
		return fmt.Errorf("failed to create path: %w", err)
	}

	// Change to target directory
	if err := os.Chdir(cfg.path); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	// Clone each repo in parallel
	type job struct {
		url  string
		name string
	}

	jobs := make([]job, 0, len(repos))
	for _, url := range repos {
		name := strings.TrimSuffix(filepath.Base(url), ".git")
		jobs = append(jobs, job{url: url, name: name})
	}

	// Process jobs in parallel using nproc
	maxWorkers := runtime.NumCPU()
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	sem := make(chan struct{}, maxWorkers)
	errChan := make(chan error, len(jobs))

	for _, j := range jobs {
		go func(j job) {
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := cloneRepo(j.url, j.name, cfg); err != nil {
				errChan <- fmt.Errorf("failed to clone %s: %w", j.name, err)
				return
			}
		}(j)
	}

	// Wait for all goroutines
	var errs []string
	for i := 0; i < len(jobs); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errs))
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("%d errors occurred during clone", len(errs))
	}

	fmt.Printf("\nDone! Cloned %d repositories to %s\n", len(jobs), cfg.path)
	return nil
}

// getRepoSSHURLs fetches all SSH URLs for repos in an organization
func getRepoSSHURLs(org, serverHost string, cfg *config) ([]string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	var urls []string
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("orgs/%s/repos?per_page=%d&page=%d", org, perPage, page)
		var repos []struct {
			SSHURL   string `json:"ssh_url"`
			Name     string `json:"name"`
			Archived bool   `json:"archived"`
			Fork     bool   `json:"fork"`
		}
		if err := client.Get(url, &repos); err != nil {
			return nil, fmt.Errorf("failed to list repos: %w", err)
		}

		for _, r := range repos {
			// Issue #3: Apply filters
			if cfg.skipArchived && r.Archived {
				continue
			}
			if cfg.skipForks && r.Fork {
				continue
			}
			if len(cfg.includePatterns) > 0 {
				if !matchesAnyPattern(r.Name, cfg.includePatterns) {
					continue
				}
			}
			if len(cfg.excludePatterns) > 0 {
				if matchesAnyPattern(r.Name, cfg.excludePatterns) {
					continue
				}
			}
			if cfg.limit > 0 && len(urls) >= cfg.limit {
				break
			}

			u := r.SSHURL
			// Replace default host with custom SSH host if needed
			if serverHost != "github.com" {
				u = strings.Replace(u, "github.com", serverHost, 1)
			}
			urls = append(urls, u)
		}

		// Check if there are more pages
		if len(repos) < perPage {
			break
		}
		page++
		
		// Issue #3: Respect limit
		if cfg.limit > 0 && len(urls) >= cfg.limit {
			break
		}
	}

	return urls, nil
}

// matchesAnyPattern checks if a name matches any of the given glob patterns
func matchesAnyPattern(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, name)
		if matched {
			return true
		}
		// Also support regex patterns
		matched, _ = regexp.MatchString(pattern, name)
		if matched {
			return true
		}
	}
	return false
}

// cloneRepo clones or updates a single repository
func cloneRepo(url, name string, cfg *config) error {
	repoPath := filepath.Join(cfg.path, name)

	// Check if repo already exists
	if _, err := os.Stat(repoPath); err == nil {
		if cfg.updateOrgFolder {
			fmt.Printf("  Updating %s...\n", name)
			return updateRepo(repoPath)
		}
		fmt.Printf("  Skipping %s (already exists)\n", name)
		return nil
	}

	fmt.Printf("  Cloning %s...\n", name)

	// Build git clone command
	args := []string{"clone", "--quiet", url, name}
	cmd := exec.Command("git", args...)
	cmd.Dir = cfg.path

	// Set environment variables
	env := os.Environ()
	if cfg.disableCloneProtection {
		env = append(env, "GIT_CLONE_PROTECTION_ACTIVE=false")
	}
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w: %s", err, string(output))
	}

	return nil
}

// updateRepo runs git pull in an existing repository
func updateRepo(path string) error {
	cmd := exec.Command("git", "pull", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %w: %s", err, string(output))
	}
	return nil
}
