# GitHub CLI Extension: gh-clone-org

A GitHub CLI extension to clone all repositories in an organization or user account as quickly as possible.

**v2.0.0** — Complete rewrite in Go with Cobra!

## Installation

```sh
gh extension install svg153/gh-clone-org
```

Or build from source:

```sh
git clone https://github.com/svg153/gh-clone-org.git
cd gh-clone-org
go build -o gh-clone-org .
sudo cp gh-clone-org /usr/local/bin/
```

## Quick Start

```sh
# Clone all repos in an organization
gh clone-org github

# Clone to a specific directory
gh clone-org github -p ~/github

# Clone a user's personal repos
gh clone-org svg153 --user

# Preview what would be cloned (dry-run)
gh clone-org github --dry-run
```

## Usage

```
gh clone-org [org|user] [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `-o, --org` | GitHub organization (positional arg alias) |
| `-p, --path` | Path to clone repositories (default: current directory) |
| `-s, --server-host-ssh` | SSH server host for multi-account setups (default: github.com) |
| `--update-org-folder` | Update existing repos and clone new ones |
| `--disable-clone-protection` | Disable `GIT_CLONE_PROTECTION_ACTIVE` |
| `--skip-archived` | Skip archived repositories |
| `--skip-forks` | Skip forked repositories |
| `--include-pattern` | Glob pattern to include repos (can be repeated) |
| `--exclude-pattern` | Glob pattern to exclude repos (can be repeated) |
| `--limit` | Maximum number of repositories to clone (0 = unlimited) |
| `--dry-run` | List repos that would be cloned without actually cloning |
| `-v, --verbose` | Increase verbosity: `-v` = clone status, `-vv` = git output, `-vvv` = debug info |
| `--user` | Clone user's personal repos instead of organization repos |
| `--profile` | Configuration profile (built-in: full, minimal, no-forks, or custom) |
| `--version` | Show version info |
| `-h, --help` | Help for this command |

## Examples

### Basic usage

```sh
# Clone all repos in an organization
gh clone-org kubernetes

# Clone to a specific directory
gh clone-org kubernetes -p ~/k8s

# Clone with verbose output
gh clone-org kubernetes -v
```

### Filters

```sh
# Skip archived and forked repos
gh clone-org github --skip-archived --skip-forks

# Clone only specific repos
gh clone-org github --include-pattern "cli-*"

# Exclude certain repos
gh clone-org github --exclude-pattern "test-*" --exclude-pattern "docs-*"

# Clone only first 10 repos
gh clone-org github --limit 10
```

### User mode

```sh
# Clone a user's personal repos
gh clone-org svg153 --user

# User mode with filters
gh clone-org svg153 --user --skip-forks --skip-archived
```

### Dry-run

```sh
# See what would be cloned without actually cloning
gh clone-org github --dry-run
gh clone-org github --dry-run --skip-forks --limit 5
```

### Configuration profiles

Built-in profiles:

- **full**: No filters, verbose mode
- **minimal**: Skip archived + forks, quiet mode
- **no-forks**: Skip forks only

```sh
# Use a built-in profile
gh clone-org github --profile minimal

# Use a custom profile from .gh-clone-org.yaml
gh clone-org github --profile production
```

Custom profiles in `.gh-clone-org.yaml`:

```yaml
production:
  skip_archived: true
  skip_forks: true
  verbose: 1
  limit: 50
  exclude_patterns:
    - "test-*"
    - "docs-*"
```

### Multi-account

```sh
# Clone from a GitHub Enterprise instance
gh clone-org mycompany -s github.mycompany.com

# Clone from a custom SSH host
gh clone-org mycompany -s github.company-ssh.com
```

## Architecture

- **Parallel cloning**: Uses `runtime.NumCPU()` workers with a semaphore channel
- **Rate limiting**: Exponential backoff with jitter (1s base, 60s max, 25% jitter)
- **Pagination**: Handles orgs with 100+ repos via GitHub API pagination
- **SSH URLs**: Clones via SSH for authenticated access
- **Error handling**: Collects all errors and reports them at the end

## Contributing

1. Fork the repository
2. Create a feature branch (`feat/your-feature`)
3. Make your changes
4. Run tests: `bats test/`
5. Commit with [Conventional Commits](https://www.conventionalcommits.org/) format
6. Push and create a Pull Request

### Commit format

```
type(scope): description

- feat: add new feature
- fix: fix a bug
- docs: update documentation
- test: add or update tests
- refactor: code refactoring
- ci: CI/CD changes
- chore: maintenance tasks
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
