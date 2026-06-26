# CodeQL Reconcile Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `work:codeql` task that reads managed work repository configuration and reconciles GitHub CodeQL default setup through `gh-usecase`.

**Architecture:** Keep Terraform responsible for repository resources and CodeQL-required rulesets. Add a separate Go CLI under `scripts/reconcile-codeql` that parses `terraform/work-repositories/locals.tf`, selects repositories with `codeql.languages`, and shells out to `gh-usecase codeql-default-setup`. The reconcile command is independent of Terraform diff, so it can fix GitHub-side drift even when Terraform has no changes.

**Tech Stack:** Go 1.26.3, HashiCorp HCL parser, existing `gh-usecase` CLI, Terraform HCL locals, mise tasks.

---

## Commit Policy For This Workspace

Do not run `git commit` unless the user explicitly asks for a commit.

## File Structure

- Create `scripts/reconcile-codeql/go.mod` and `go.sum`: standalone Go module for the reconcile CLI.
- Create `scripts/reconcile-codeql/main.go`: argument parsing, HCL reading, `codeql.languages` extraction, command execution, result summary.
- Create `scripts/reconcile-codeql/main_test.go`: offline tests for argument validation, HCL parsing, command construction, target filtering, failure aggregation, and dry-run behavior.
- Modify `scripts/add-work-repository/main.go`: generate `codeql.languages` instead of legacy `enable_codeql`.
- Modify `scripts/add-work-repository/main_test.go`: cover the new CodeQL HCL shape.
- Modify `terraform/work-repositories/locals.tf`: replace `enable_codeql = true` with `codeql = { languages = ["go"] }`; add CodeQL config for `gh-usecase`.
- Modify `terraform/work-repositories/main.tf`: pass CodeQL-required ruleset state from `try(each.value.codeql, null) != null`.
- Modify `.mise/config.toml`: add `work:codeql` task.
- Modify `README.md`: document `work:codeql` and the fact that it is independent of Terraform diff.

## Tasks

### Task 1: Add failing reconcile parser tests

**Files:**
- Create: `scripts/reconcile-codeql/main_test.go`

- [x] Add tests that define `locals.repositories` with one repository containing `codeql.languages`, one without CodeQL, and one existing legacy `enable_codeql` entry.
- [x] Verify the parser returns only repositories with `codeql.languages`.
- [x] Verify missing languages or legacy-only `enable_codeql` returns a validation error instead of guessing languages.
- [x] Run `go -C scripts/reconcile-codeql test ./...` and confirm it fails because the module and functions do not exist yet.

### Task 2: Implement HCL parsing and target extraction

**Files:**
- Create: `scripts/reconcile-codeql/go.mod`
- Create: `scripts/reconcile-codeql/main.go`

- [x] Implement `readTargets(src []byte, filename string) ([]codeqlTarget, error)` using `hclsyntax.ParseConfig`.
- [x] Extract `name` and static `codeql.languages` string list from repository object entries.
- [x] Reject entries that set `enable_codeql = true` without `codeql.languages` with a clear error.
- [x] Run `go -C scripts/reconcile-codeql test ./...` and confirm parser tests pass.

### Task 3: Add command construction and dry-run tests

**Files:**
- Modify: `scripts/reconcile-codeql/main_test.go`
- Modify: `scripts/reconcile-codeql/main.go`

- [x] Add tests for `buildInvocation(owner, target)` returning `gh-usecase codeql-default-setup --owner <owner> --repo <name> --languages <csv>`.
- [x] Add tests for `--repo <name>` filtering.
- [x] Add tests for `--dry-run` printing commands without executing them.
- [x] Implement the smallest CLI plumbing to pass the tests.
- [x] Run `go -C scripts/reconcile-codeql test ./...`.

### Task 4: Add execution and failure aggregation

**Files:**
- Modify: `scripts/reconcile-codeql/main_test.go`
- Modify: `scripts/reconcile-codeql/main.go`

- [x] Add tests that inject a fake command runner and verify all targets are attempted even if one fails.
- [x] Verify the final error contains the failed repository names.
- [x] Implement command execution with `os/exec` in production dependencies.
- [x] Run `go -C scripts/reconcile-codeql test ./...`.

### Task 5: Wire configuration and tasks

**Files:**
- Modify: `terraform/work-repositories/locals.tf`
- Modify: `terraform/work-repositories/main.tf`
- Modify: `.mise/config.toml`

- [x] Change CodeQL-enabled repositories from `enable_codeql = true` to `codeql = { languages = ["go"] }`.
- [x] Add `codeql = { languages = ["go"] }` to `gh-usecase`.
- [x] Change the governance module call to `enable_codeql = try(each.value.codeql, null) != null`.
- [x] Add `work:codeql` mise task running `go -C scripts/reconcile-codeql run . --owner y-writings ../../terraform/work-repositories/locals.tf`.
- [x] Run `terraform fmt terraform/work-repositories/locals.tf terraform/work-repositories/main.tf`.

### Task 5a: Keep add-work-repository on the current config shape

**Files:**
- Modify: `scripts/add-work-repository/main_test.go`
- Modify: `scripts/add-work-repository/main.go`

- [x] Update the render test to expect `codeql = { languages = ["go"] }` instead of `enable_codeql = true`.
- [x] Replace the CodeQL boolean prompt with a language multi-select.
- [x] Generate `codeql.languages` when at least one language is selected.
- [x] Run `go -C scripts/add-work-repository test ./...`.

### Task 6: Documentation and verification

**Files:**
- Modify: `README.md`

- [x] Document `mise run work:codeql`.
- [x] State that Terraform diff does not control whether the reconcile task runs; the task runs only when invoked, and `gh-usecase` only PATCHes when GitHub differs.
- [x] Run `go -C scripts/add-work-repository test ./...`.
- [x] Run `go -C scripts/reconcile-codeql test ./...`.
- [x] Run `terraform fmt -check terraform/work-repositories/locals.tf terraform/work-repositories/main.tf`.
