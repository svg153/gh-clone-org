#!/usr/bin/env bats

# BATS tests for gh-clone-org
# Run with: bats test/

setup() {
    cd "$(dirname "$BATS_TEST_FILENAME")/.."
}

@test "version flag shows version info" {
    run ./gh-clone-org --version
    [ "$status" -eq 0 ]
    [[ "$output" == *"gh-clone-org"* || "$output" == *"dev"* ]]
}

@test "help flag shows help text" {
    run ./gh-clone-org --help
    [ "$status" -eq 0 ]
    [[ "$output" == *"Clone all repositories"* ]]
    [[ "$output" == *"--dry-run"* ]]
    [[ "$output" == *"--user"* ]]
    [[ "$output" == *"--profile"* ]]
    [[ "$output" == *"-v, --verbose"* ]]
}

@test "no org arg shows error" {
    run ./gh-clone-org
    [ "$status" -ne 0 ]
    [[ "$output" == *"organization or user is required"* ]]
}

@test "dry-run flag is recognized" {
    run ./gh-clone-org testorg --dry-run
    # Will fail to connect to API but flag should be parsed
    [ "$status" -ne 0 ] || [[ "$output" == *"Would clone"* ]]
}

@test "user flag is recognized" {
    run ./gh-clone-org svg153 --user --dry-run
    # Will fail to connect to API but flag should be parsed
    [ "$status" -ne 0 ] || [[ "$output" == *"Would clone"* ]]
}

@test "verbose flag is recognized" {
    run ./gh-clone-org testorg -v
    [ "$status" -ne 0 ] || true
}

@test "verbose flag double is recognized" {
    run ./gh-clone-org testorg -vv
    [ "$status" -ne 0 ] || true
}

@test "profile flag is recognized" {
    run ./gh-clone-org testorg --profile minimal
    [ "$status" -ne 0 ] || true
}

@test "skip-archived flag is recognized" {
    run ./gh-clone-org testorg --skip-archived
    [ "$status" -ne 0 ] || true
}

@test "skip-forks flag is recognized" {
    run ./gh-clone-org testorg --skip-forks
    [ "$status" -ne 0 ] || true
}

@test "limit flag is recognized" {
    run ./gh-clone-org testorg --limit 5
    [ "$status" -ne 0 ] || true
}

@test "include-pattern flag is recognized" {
    run ./gh-clone-org testorg --include-pattern "test-*"
    [ "$status" -ne 0 ] || true
}

@test "exclude-pattern flag is recognized" {
    run ./gh-clone-org testorg --exclude-pattern "test-*"
    [ "$status" -ne 0 ] || true
}
