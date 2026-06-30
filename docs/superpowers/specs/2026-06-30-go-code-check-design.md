# Go Code Check Design

## Goal

Add a small, reliable code-checking layer for this Go toolkit repository. The first version should make local checks and GitHub CI run the same core commands, so contributors can catch formatting, vet, and test failures before code is merged.

## Scope

The first version includes:

- A `Makefile` with local check targets.
- A GitHub Actions workflow for push and pull request checks.
- Format verification with `gofmt`.
- Static correctness checks with `go vet ./...`.
- Test execution with `go test ./...`.

The first version does not include:

- `staticcheck`.
- `golangci-lint`.
- Security scanning.
- Coverage thresholds.
- AI-based code review.

Those can be added later after the baseline check flow is stable.

## Architecture

The check flow has two layers:

1. Local developer entrypoint: `make check`.
2. Remote CI entrypoint: `.github/workflows/ci.yml`.

The CI workflow delegates repository-specific validation to `make check` instead of duplicating all commands in YAML. This keeps local and CI behavior aligned.

## Make Targets

The `Makefile` should expose these targets:

- `fmt`: verify all Go files are formatted with `gofmt`.
- `vet`: run `go vet ./...`.
- `test`: run `go test ./...`.
- `check`: run `fmt`, `vet`, and `test` in order.

The `fmt` target should fail if any Go file would be changed by `gofmt`. It should not silently rewrite files during CI checks.

## CI Flow

The GitHub Actions workflow should run on:

- `push`
- `pull_request`

The workflow should:

1. Check out the repository.
2. Install the Go version from `go.mod`.
3. Download modules.
4. Run `make check`.

If any check fails, the workflow should fail.

## Error Handling

All targets should exit non-zero when their underlying command fails.

The format target should print the unformatted file paths before failing, so the developer knows what to fix.

## Testing Strategy

After implementation, verify locally with:

```bash
make check
```

Also verify that the workflow YAML is syntactically valid enough for GitHub Actions by keeping it minimal and using standard actions.

## Future Extensions

Possible follow-up improvements:

- Add `staticcheck` as `make staticcheck`.
- Add `golangci-lint` after choosing a rule profile.
- Add `go test -race ./...` if runtime allows.
- Add coverage reporting without enforcing thresholds at first.
- Add dependency and secret scanning.
