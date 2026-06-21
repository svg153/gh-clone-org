# GitHub CLI Extension: gh-clone-org

A GitHub CLI extension to clone all repositories in an organization as quickly as possible.

**v2.0.0** — Complete rewrite in Go!

## Installation

```sh
gh extension install svg153/gh-clone-org
```

## Usage

```sh
gh clone-org <org>
```

For help, run:

```sh
gh clone-org --help
```

## Features

- **Parallel cloning** — Uses all available CPU cores
- **Org validation** — Checks if the target is an organization (not a user)
- **Update mode** — Update existing repos and clone new ones (`--update-org-folder`)
- **Skip clone protection** — Disable `GIT_CLONE_PROTECTION_ACTIVE` (`--disable-clone-protection`)
- **Multi-account** — Custom SSH host for multiple GitHub accounts (`-s github.com-company`)
- **Pagination** — Handles orgs with 100+ repos

## Roadmap (v2+)

See [Issue #2](https://github.com/svg153/gh-clone-org/issues/2) for the full roadmap. Upcoming features:

- [ ] Filter archived/forks repos
- [ ] Rate limiting with exponential backoff
- [ ] Dry-run mode
- [ ] Verbose logging (`-v`, `-vv`, `-vvv`)
- [ ] User mode (`--user`)
- [ ] Configuration profiles

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
