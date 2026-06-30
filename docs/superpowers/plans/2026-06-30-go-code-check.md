# Go Code Check Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a baseline Go code-checking workflow that runs formatting verification, vet, and tests locally and in GitHub Actions.

**Architecture:** The repository gets one local entrypoint, `make check`, that owns the actual validation commands. GitHub Actions delegates to `make check` so local and CI behavior stay aligned.

**Tech Stack:** Go, Make, GitHub Actions, `gofmt`, `go vet`, `go test`.

## Global Constraints

- The first version includes a `Makefile` with local check targets.
- The first version includes a GitHub Actions workflow for push and pull request checks.
- The first version verifies formatting with `gofmt`.
- The first version runs static correctness checks with `go vet ./...`.
- The first version runs tests with `go test ./...`.
- The first version does not include `staticcheck`, `golangci-lint`, security scanning, coverage thresholds, or AI-based code review.
- The `fmt` target must fail if any Go file would be changed by `gofmt`.
- The `fmt` target must print unformatted file paths before failing.
- CI must run on `push` and `pull_request`.
- CI must run `make check`.

---

## File Structure

- Create `Makefile`: local developer entrypoint for formatting verification, vet, tests, and the aggregate check.
- Create `.github/workflows/ci.yml`: GitHub Actions workflow that checks out the repo, sets up Go from `go.mod`, downloads modules, and runs `make check`.

### Task 1: Local Make Check Entrypoint

**Files:**
- Create: `Makefile`

**Interfaces:**
- Consumes: Go packages under `./...`.
- Produces: `make fmt`, `make vet`, `make test`, and `make check`.

- [ ] **Step 1: Confirm there is no existing Makefile**

Run:

```bash
test ! -f Makefile
```

Expected: command exits with status `0`.

- [ ] **Step 2: Create the Makefile**

Create `Makefile` with exactly this content:

```makefile
.PHONY: fmt vet test check

GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

fmt:
	@unformatted="$$(gofmt -l $(GO_FILES))"; \
	if [ -n "$$unformatted" ]; then \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet test
```

- [ ] **Step 3: Verify the Makefile exposes the expected targets**

Run:

```bash
make -n check
```

Expected: output includes commands for `gofmt -l`, `go vet ./...`, and `go test ./...`.

- [ ] **Step 4: Run the local check**

Run:

```bash
make check
```

Expected: command exits with status `0`.

- [ ] **Step 5: Commit the local check entrypoint**

Run:

```bash
git add Makefile
git commit -m "chore: add go code check targets"
```

Expected: commit succeeds and includes only `Makefile`.

### Task 2: GitHub Actions CI

**Files:**
- Create: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `make check` from Task 1.
- Produces: a CI workflow named `CI` that runs on `push` and `pull_request`.

- [ ] **Step 1: Confirm there is no existing CI workflow**

Run:

```bash
test ! -f .github/workflows/ci.yml
```

Expected: command exits with status `0`.

- [ ] **Step 2: Create the GitHub Actions workflow**

Create `.github/workflows/ci.yml` with exactly this content:

```yaml
name: CI

on:
  push:
  pull_request:

jobs:
  check:
    runs-on: ubuntu-latest

    steps:
      - name: Check out repository
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Download modules
        run: go mod download

      - name: Run checks
        run: make check
```

- [ ] **Step 3: Verify the workflow references the local check entrypoint**

Run:

```bash
rg -n "make check|go-version-file: go.mod|pull_request|push" .github/workflows/ci.yml
```

Expected: output includes all four patterns.

- [ ] **Step 4: Run the same check CI will run**

Run:

```bash
make check
```

Expected: command exits with status `0`.

- [ ] **Step 5: Commit the CI workflow**

Run:

```bash
git add .github/workflows/ci.yml
git commit -m "ci: run go code checks"
```

Expected: commit succeeds and includes only `.github/workflows/ci.yml`.

## Final Verification

- [ ] **Step 1: Run all local checks**

Run:

```bash
make check
```

Expected: command exits with status `0`.

- [ ] **Step 2: Inspect the final diff against the design commit**

Run:

```bash
git status --short
git log --oneline -3
```

Expected: working tree is clean, and recent commits include:

```text
ci: run go code checks
chore: add go code check targets
docs: add go code check design
```
