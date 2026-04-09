# Migration to Nix Flakes — Proposal

**Project:** `go-composable-business-types`
**Date:** 2026-04-09
**Status:** Draft

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State Analysis](#2-current-state-analysis)
3. [Why Nix Flakes](#3-why-nix-flakes)
4. [Proposed Architecture](#4-proposed-architecture)
5. [Detailed Migration Plan](#5-detailed-migration-plan)
6. [Risk Assessment](#6-risk-assessment)
7. [Migration Checklist](#7-migration-checklist)
8. [Appendix: Reference flake.nix](#appendix-reference-flakenix)

---

## 1. Executive Summary

This proposal outlines the migration of `go-composable-business-types` from ad-hoc local tooling management to a **Nix Flakes**-based development environment. The goal is to achieve **fully reproducible, pinned, and declarative** developer toolchains — eliminating "works on my machine" problems and ensuring every contributor runs identical versions of Go, golangci-lint, go-enum, govulncheck, git-cliff, and just.

The migration is **non-breaking**: existing workflows (`just test`, `just lint`, CI pipelines) continue to work unchanged. Nix adds a new, superior entry point (`nix develop`) while the justfile remains the canonical task runner.

---

## 2. Current State Analysis

### 2.1 Project Characteristics

| Aspect              | Detail                                                        |
| ------------------- | ------------------------------------------------------------- |
| **Type**            | Go library (no `main` binary, no `cmd/` entrypoint)           |
| **Go version**      | 1.26.0 (`go.mod`), 1.26.1 (`.golangci.yml`)                   |
| **CI Go matrix**    | 1.26, 1.27, 1.28                                              |
| **Module path**     | `github.com/larsartmann/go-composable-business-types`         |
| **Key experiment**  | `GOEXPERIMENT=jsonv2` used everywhere (justfile, CI, linting) |
| **Code generation** | `go-enum` via `//go:generate go tool go-enum`                 |

### 2.2 Tool Inventory

Every tool listed below must be version-pinned in the Nix flake:

| Tool                          | Current Source            | Version                   | Used By                            |
| ----------------------------- | ------------------------- | ------------------------- | ---------------------------------- |
| **Go**                        | `actions/setup-go@v5`     | 1.26 / 1.27 / 1.28 matrix | build, test, generate              |
| **golangci-lint**             | `golangci-lint-action@v6` | v1.64                     | `just lint`, CI                    |
| **go-enum**                   | `go.mod` tool directive   | v0.9.2                    | `just generate`, CI                |
| **govulncheck**               | `go install @latest`      | unpinned                  | CI security job                    |
| **git-cliff**                 | `go install @latest`      | unpinned                  | `just changelog`, release workflow |
| **just**                      | Not managed               | system-dependent          | all justfile commands              |
| **gotools** (goimports, etc.) | implicit with Go          | Go 1.26 bundled           | `.golangci.yml` formatters         |
| **gofumpt**                   | via golangci-lint         | follows golangci-lint     | `.golangci.yml` formatters         |
| **gci**                       | via golangci-lint         | follows golangci-lint     | `.golangci.yml` formatters         |
| **golines**                   | via golangci-lint         | follows golangci-lint     | `.golangci.yml` formatters         |

### 2.3 Environment Variables

The `GOEXPERIMENT=jsonv2` flag is **critical** — it is set in:

- `justfile` (every recipe via `GOEXPERIMENT=jsonv2` prefix)
- `.github/workflows/ci.yml` (top-level `env:`)
- `.github/workflows/release.yml` (top-level `env:`)
- `.golangci.yml` (`run.build-tags: [goexperiment.jsonv2]`)

The Nix devShell **must** export this environment variable so that bare `go build ./...` works identically to `just build`.

### 2.4 Build Tags in `.golangci.yml`

```yaml
build-tags:
  - goexperiment.goroutineleakprofile
  - goexperiment.jsonv2
  - goexperiment.simd
```

These are handled within golangci-lint's own config and do not need Nix-level treatment, but the devShell should set `GOEXPERIMENT=jsonv2` for consistency.

### 2.5 CI Pipeline Summary

```
CI (push/PR to main|master)
├── test (Go 1.26, 1.27, 1.28 matrix + coverage + 80% threshold)
├── lint (golangci-lint v1.64, 5m timeout)
├── security (govulncheck)
├── generate (go generate + uncommitted-changes check)
└── benchmark

Release (tag push: YYYY-MM-DD.N)
├── test + lint + generate check
├── git-cliff → release_body.md
└── GitHub Release creation
```

---

## 3. Why Nix Flakes

### 3.1 Problems Solved

| Problem                     | Current Impact                              | Nix Solution                   |
| --------------------------- | ------------------------------------------- | ------------------------------ |
| **Unpinned govulncheck**    | `go install @latest` → non-reproducible CI  | Pinned in flake.lock           |
| **Unpinned git-cliff**      | `go install @latest` → drift between runs   | Pinned in flake.lock           |
| **No local tool isolation** | Tools pollute `$GOPATH/bin` globally        | Hermetic `nix develop` shell   |
| **Onboarding friction**     | "Install Go 1.26, golangci-lint, just, ..." | `nix develop` — done           |
| **Platform inconsistency**  | macOS vs Linux tool version mismatches      | Same flake.lock, same versions |
| **No direnv integration**   | Manual shell activation                     | `.envrc` → automatic           |

### 3.2 What Nix Flakes Will NOT Replace

- **GitHub Actions CI/CD** — CI continues using `actions/setup-go` and `golangci-lint-action` (these are well-maintained, fast, and CI-native)
- **justfile** — remains the canonical task runner; Nix provides the _environment_, just provides the _commands_
- **go.mod / go.sum** — Go's own dependency management is untouched
- **Codecov integration** — unchanged

### 3.3 Scope Decision: devShell-First

This proposal focuses on **devShell + checks** rather than package building because:

1. This is a **library**, not a binary — consumers `go get` it, they don't install it
2. `buildGoModule` would add maintenance burden (vendorHash updates) for zero consumer benefit
3. The primary value is reproducible **developer environments** and **CI-optional local checks**

Package outputs (`packages.default`) are included as a stretch goal (§5.4) but are not required for initial migration.

---

## 4. Proposed Architecture

### 4.1 File Structure

```
.
├── flake.nix          # Flake definition (new)
├── flake.lock         # Pinned dependency versions (new, auto-generated)
├── .envrc             # direnv integration (new)
├── justfile           # Unchanged — canonical task runner
├── go.mod             # Unchanged — Go dependency management
├── go.sum             # Unchanged
├── .golangci.yml      # Unchanged
├── cliff.toml         # Unchanged
└── ...
```

### 4.2 Flake Outputs

```
flake.nix
├── inputs
│   ├── nixpkgs (nixos-unstable)
│   └── flake-utils (multi-system support)
├── outputs
│   ├── devShells.<system>.default    # Primary: dev environment
│   ├── checks.<system>               # Optional: nix flake check
│   ├── packages.<system>.default     # Stretch: library metadata
│   └── formatter.<system>            # Optional: nix fmt for .nix files
```

### 4.3 Supported Systems

Matching the project's target platforms:

```nix
supportedSystems = [
  "x86_64-linux"     # CI (ubuntu-latest)
  "aarch64-linux"    # ARM CI / servers
  "x86_64-darwin"    # macOS Intel
  "aarch64-darwin"   # macOS Apple Silicon
];
```

### 4.4 Design Decision: flake-utils vs flake-parts

**Recommendation: `flake-utils`** (not `flake-parts`)

| Factor               | flake-utils                      | flake-parts                          |
| -------------------- | -------------------------------- | ------------------------------------ |
| Complexity           | Minimal                          | Module system, more to learn         |
| This project's needs | devShell + checks only           | Overkill for a Go library            |
| Maintenance          | Near-zero                        | Requires understanding module system |
| Ecosystem prevalence | Extremely common for Go projects | Better for multi-output repos        |

`flake-parts` adds value for repos with NixOS modules, Home Manager configs, or many packages. This project needs exactly one devShell. `flake-utils` is the right tool.

---

## 5. Detailed Migration Plan

### 5.1 Phase 1: Core devShell (Required)

**Goal:** `nix develop` gives you a complete, reproducible development environment.

**New files:**

- `flake.nix`
- `.envrc`

**devShell contents:**

```nix
devShells.default = pkgs.mkShell {
  packages = with pkgs; [
    # Go toolchain
    go_1_26

    # Linting & formatting
    golangci-lint

    # Go tooling
    gotools          # goimports, godoc, etc.
    gofumpt
    golines

    # Code generation (go-enum is a Go tool directive — available via `go tool go-enum`)
    # No separate Nix package needed; it's in go.mod

    # Security
    govulncheck

    # Changelog
    git-cliff

    # Task runner
    just

    # Misc
    gnumake          # for any potential Makefile interop
  ];

  GOEXPERIMENT = "jsonv2";

  shellHook = ''
    echo "go-composable-business-types dev shell"
    echo "Go: $(go version)"
    echo "GOEXPERIMENT=$GOEXPERIMENT"
  '';
};
```

**Key design decisions:**

1. **`go_1_26` not `go`** — Explicitly pin the minor version matching `go.mod`. Upgrading is a deliberate `flake.nix` edit, not an accident.
2. **`GOEXPERIMENT = "jsonv2"`** — Set as a Nix shell variable, so all tools and bare `go` commands pick it up automatically. This replaces the per-recipe prefix in the justfile (though the justfile remains unchanged for backwards compatibility).
3. **go-enum via `go tool`** — Since `go-enum` is declared in `go.mod` as a `tool` directive (`tool github.com/abice/go-enum`), it is already available via `go tool go-enum`. No separate Nix package needed.
4. **govulncheck pinned** — No more `go install @latest` drift.

**Verification steps:**

```bash
# Enter dev shell
nix develop

# Verify tools
go version                    # go1.26.x
echo $GOEXPERIMENT            # jsonv2
golangci-lint version         # pinned version
govulncheck --help            # available
git-cliff --version           # available
just --version                # available

# Verify workflows
just generate                 # go generate via go tool go-enum
just build                    # GOEXPERIMENT=jsonv2 go build
just test                     # all tests pass
just lint                     # golangci-lint passes
just check                    # full pipeline
```

### 5.2 Phase 2: direnv Integration (Recommended)

**Goal:** Automatic shell activation when entering the project directory.

**New file: `.envrc`**

```bash
# .envrc
if ! has nix_direnv_version || ! nix_direnv_version 3.0.0; then
  source_url "https://raw.githubusercontent.com/nix-community/nix-direnv/3.0.0/direnvrc" "sha256-21TMnI2xWX7HkSKj3Glrny5Ihoj3MLXJLRX/LJ6NNAQ="
fi

use flake
```

**Also consider: `nix-direnv`** for faster shell loading (caches the devShell profile).

**User prerequisites:**

```bash
# One-time setup (macOS with nix-darwin or home-manager)
nix profile install nixpkgs#direnv
nix profile install nixpkgs#nix-direnv

# Add to shell rc (~/.zshrc, ~/.bashrc)
eval "$(direnv hook zsh)"  # or bash

# Enable in project
direnv allow
```

### 5.3 Phase 3: Nix Checks (Optional Enhancement)

**Goal:** `nix flake check` runs the same quality gates as CI, locally.

```nix
checks = {
  build = pkgs.runCommandLocal "check-build" {
      nativeBuildInputs = [ pkgs.go_1_26 ];
      GOEXPERIMENT = "jsonv2";
    } ''
      cp -r ${./.}/* ./
      go build ./...
      touch $out
    '';

  test = pkgs.runCommandLocal "check-test" {
      nativeBuildInputs = [ pkgs.go_1_26 ];
      GOEXPERIMENT = "jsonv2";
    } ''
      cp -r ${./.}/* ./
      go test -race ./...
      touch $out
    '';

  lint = pkgs.runCommandLocal "check-lint" {
      nativeBuildInputs = [ pkgs.go_1_26 pkgs.golangci-lint ];
      GOEXPERIMENT = "jsonv2";
    } ''
      cp -r ${./.}/* ./
      golangci-lint run --timeout=5m
      touch $out
    '';

  vet = pkgs.runCommandLocal "check-vet" {
      nativeBuildInputs = [ pkgs.go_1_26 ];
      GOEXPERIMENT = "jsonv2";
    } ''
      cp -r ${./.}/* ./
      go vet ./...
      touch $out
    '';
};
```

**Note on checks:** Nix checks copy source into the Nix store, which means they don't see uncommitted changes. They verify the _committed_ state. For day-to-day development, `just check` remains the primary workflow. `nix flake check` is useful as a pre-push gate or CI alternative.

**Consideration:** Go's module system reads `go.sum` and may need network access. Nix builds are sandboxed by default. The `__noChroot` pattern or `buildGoModule` with vendorHash may be needed for fully sandboxed builds. This is why checks are **optional** — the devShell is the primary value.

### 5.4 Phase 4: Package Output (Stretch Goal)

**Goal:** Expose the library as a Nix package (primarily for downstream Nix consumers).

```nix
packages.default = pkgs.buildGoModule {
  pname = "go-composable-business-types";
  version = "0.0.0";  # overridden by --override-input or self.lastModifiedDate

  src = ./.;

  # vendorHash must be updated when go.mod/go.sum change
  # Use `pkgs.lib.fakeHash` during development, then replace with actual hash
  vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

  GOEXPERIMENT = "jsonv2";

  # Library-only: no subPackages, no installPhase binary
  # buildGoModule will build and verify all packages
};
```

**Maintenance burden:** Every `go.mod` change requires updating `vendorHash`. For a library consumed via `go get`, this provides minimal value. **Defer unless there's a downstream Nix consumer.**

### 5.5 Phase 5: Nix Formatter (Nice-to-Have)

```nix
formatter = pkgs.nixpkgs-fmt;
```

Enables `nix fmt` to format `flake.nix` consistently.

### 5.6 CI Considerations

**Option A (Recommended): Keep CI as-is**

GitHub Actions CI remains unchanged. It already uses well-maintained actions:

- `actions/setup-go@v5` — fast, cached Go installation
- `golangci-lint-action@v6` — optimized linting with caching
- `codecov/codecov-action@v4` — coverage reporting

These are **better suited for CI** than Nix-based builds because:

- GitHub-hosted runners have warm caches for popular actions
- `setup-go` caching is faster than Nix sandbox builds
- The CI matrix (Go 1.26/1.27/1.28) is trivially expressed with the matrix strategy

**Option B (Future): Nix-based CI**

If the project moves to a self-hosted CI or wants to test the Nix build:

```yaml
- name: Install Nix
  uses: DeterminateSystems/nix-installer-action@main

- name: Run checks
  run: nix flake check

- name: Run tests
  run: nix develop --command just test
```

This is **not recommended** for this project — the current CI is cleaner and faster.

---

## 6. Risk Assessment

### 6.1 Risks and Mitigations

| Risk                                     | Severity | Likelihood | Mitigation                                                                          |
| ---------------------------------------- | -------- | ---------- | ----------------------------------------------------------------------------------- |
| **nixpkgs Go 1.26 not yet available**    | High     | Low        | Use `go_1_26` from nixpkgs-unstable; fallback: overlay with custom Go build         |
| **golangci-lint version mismatch**       | Medium   | Medium     | Pin via `nixpkgs` revision or use `buildGoModule` for exact version                 |
| **govulncheck version drift in nixpkgs** | Low      | Medium     | govulncheck is simple; version differences rarely matter                            |
| **git-cliff not in nixpkgs**             | Medium   | Low        | `git-cliff` is packaged in nixpkgs; if missing, `buildGoModule` fallback            |
| **GOEXPERIMENT not respected by Nix Go** | High     | Low        | Go from nixpkgs respects `GOEXPERIMENT` env var like any Go install                 |
| **flake.lock grows stale**               | Low      | High       | Add `nix flake update` to monthly maintenance; Dependabot doesn't manage flake.lock |
| **Team unfamiliarity with Nix**          | Medium   | High       | Proposal includes `.envrc` for transparent activation; justfile unchanged           |
| **macOS aarch64-darwin build issues**    | Low      | Low        | nixpkgs has excellent Apple Silicon support since 2021                              |

### 6.2 Rollback Plan

The migration is **additive** — no existing files are modified. Rollback is:

```bash
rm flake.nix flake.lock .envrc
```

No changes to justfile, CI, go.mod, or any other files.

---

## 7. Migration Checklist

### Phase 1: Core devShell

- [ ] Install Nix (if not present): `sh <(curl -L https://nixos.org/nix/install) --daemon`
- [ ] Enable flakes: add to `/etc/nix/nix.conf` or `~/.config/nix/nix.conf`:
  ```
  experimental-features = nix-command flakes
  ```
- [ ] Create `flake.nix` (see Appendix)
- [ ] Run `nix develop` to verify all tools are available
- [ ] Run `just check` from within `nix develop` to verify full pipeline
- [ ] Commit `flake.nix` and generated `flake.lock`
- [ ] Update `.gitignore` to not ignore `flake.lock` (it should be committed)

### Phase 2: direnv Integration

- [ ] Install `direnv` and `nix-direnv`
- [ ] Create `.envrc` with `use flake`
- [ ] Run `direnv allow`
- [ ] Verify automatic shell activation on `cd`
- [ ] Commit `.envrc`

### Phase 3: Nix Checks (Optional)

- [ ] Add `checks` to `flake.nix`
- [ ] Run `nix flake check` and verify
- [ ] Consider adding to CI as an additional validation step

### Phase 4: Documentation Updates

- [ ] Update `AGENTS.md` with Nix commands
- [ ] Update `README.md` with "Getting Started" section mentioning Nix
- [ ] Update `SUPPORT.md` if it references tool installation

### Phase 5: Maintenance Cadence

- [ ] Monthly: `nix flake update` to refresh nixpkgs
- [ ] On Go version bump: update `go_1_XX` in `flake.nix`
- [ ] On golangci-lint version bump: verify nixpkgs has matching version

---

## Appendix: Reference flake.nix

This is the complete, production-ready `flake.nix` for Phase 1 + Phase 3:

```nix
{
  description = "go-composable-business-types — strongly typed, composable base values for business applications";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        # Development shell — the primary deliverable
        # Usage: nix develop
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Go toolchain (pinned to match go.mod)
            go_1_26

            # Linting & static analysis
            golangci-lint

            # Go extended tooling
            gotools

            # Formatting (also available via golangci-lint, but explicit for standalone use)
            gofumpt
            golines

            # Security scanning
            govulncheck

            # Changelog generation
            git-cliff

            # Task runner
            just
          ];

          # Critical: enable jsonv2 experiment globally
          GOEXPERIMENT = "jsonv2";

          shellHook = ''
            echo "go-composable-business-types"
            echo "  Go:        $(go version)"
            echo "  Lint:      $(golangci-lint version --format short 2>/dev/null || echo 'golangci-lint')"
            echo "  Govuln:    $(govulncheck version 2>/dev/null || echo 'govulncheck')"
            echo "  Cliff:     $(git-cliff --version 2>/dev/null || echo 'git-cliff')"
            echo "  GOEXP:     $GOEXPERIMENT"
            echo ""
            echo "Run 'just' to see available commands."
          '';
        };

        # Optional: nix flake check runs build, test, lint, vet
        checks =
          let
            src = ./.;
            goPkg = pkgs.go_1_26;
            goEnv = {
              GOEXPERIMENT = "jsonv2";
            };
          in
          {
            build = pkgs.runCommandLocal "check-build" (goEnv // {
              nativeBuildInputs = [ goPkg ];
            }) ''
              cp -r ${src}/* ./
              go build ./...
              touch $out
            '';

            test = pkgs.runCommandLocal "check-test" (goEnv // {
              nativeBuildInputs = [ goPkg ];
            }) ''
              cp -r ${src}/* ./
              go test -race ./...
              touch $out
            '';

            vet = pkgs.runCommandLocal "check-vet" (goEnv // {
              nativeBuildInputs = [ goPkg ];
            }) ''
              cp -r ${src}/* ./
              go vet ./...
              touch $out
            '';
          };

        # Optional: format .nix files
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}
```

---

## Summary of Decisions

| Decision            | Choice                           | Rationale                                                          |
| ------------------- | -------------------------------- | ------------------------------------------------------------------ |
| **Flake framework** | `flake-utils`                    | Minimal, sufficient for a Go library                               |
| **Primary output**  | `devShells.default`              | Library has no binary to ship                                      |
| **Package output**  | Deferred                         | No downstream Nix consumer; `go get` is the distribution mechanism |
| **Go version**      | `go_1_26` (explicit)             | Matches `go.mod`; prevents accidental upgrades                     |
| **GOEXPERIMENT**    | Shell-level env var              | Single source of truth; all tools inherit it                       |
| **go-enum**         | Via `go tool`                    | Already in `go.mod` tool directive; no Nix package needed          |
| **CI migration**    | None                             | Current GitHub Actions CI is optimal                               |
| **direnv**          | Recommended                      | Transparent activation; no manual `nix develop`                    |
| **Rollback**        | `rm flake.nix flake.lock .envrc` | Additive migration; zero risk                                      |
