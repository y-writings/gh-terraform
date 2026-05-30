package main

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/huh"
)

var validRunTarget = []byte(`locals {
  repositories = {
    repo_aaaaaaaa = {
      name = "existing-repo"
    }
  }
}`)

func TestParseArgsRequiresExactlyOnePath(t *testing.T) {
	t.Parallel()

	if _, err := parseArgs(nil); err == nil {
		t.Fatal("parseArgs(nil) returned nil error")
	}

	if _, err := parseArgs([]string{"one", "two"}); err == nil {
		t.Fatal("parseArgs with two args returned nil error")
	}

	path, err := parseArgs([]string{"../../terraform/work-repositories/locals.tf"})
	if err != nil {
		t.Fatalf("parseArgs returned unexpected error: %v", err)
	}
	if path != "../../terraform/work-repositories/locals.tf" {
		t.Fatalf("path = %q", path)
	}
}

func TestRunRejectsMissingTargetFile(t *testing.T) {
	t.Parallel()

	err := run([]string{"/path/that/does/not/exist/locals.tf"})
	if err == nil {
		t.Fatal("run returned nil error for missing file")
	}
	if !strings.Contains(err.Error(), "read target file") {
		t.Fatalf("error = %q, want read target file", err.Error())
	}
}

func TestRunWithDepsTreatsPromptAbortAsCancellation(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	writeCalled := false
	fmtCalled := false

	err := runWithDeps([]string{"locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return validRunTarget, nil
		},
		writeFile: func(path string, data []byte, perm os.FileMode) error {
			writeCalled = true
			return nil
		},
		prompt: func(existingNames map[string]struct{}) (repositoryInput, error) {
			return repositoryInput{}, huh.ErrUserAborted
		},
		generateKey: func(existingKeys map[string]struct{}) (string, error) {
			return "repo_bbbbbbbb", nil
		},
		confirm: func(targetPath string, key string, input repositoryInput) (bool, error) {
			return true, nil
		},
		terraformFmt: func(targetPath string) error {
			fmtCalled = true
			return nil
		},
		stdout: &stdout,
	})
	if err != nil {
		t.Fatalf("runWithDeps returned error: %v", err)
	}
	if got := stdout.String(); got != "No changes made.\n" {
		t.Fatalf("stdout = %q, want cancellation message", got)
	}
	if writeCalled {
		t.Fatal("writeFile was called after prompt abort")
	}
	if fmtCalled {
		t.Fatal("terraformFmt was called after prompt abort")
	}
}

func TestRunWithDepsTreatsRejectedConfirmationAsCancellation(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	writeCalled := false
	fmtCalled := false

	err := runWithDeps([]string{"locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return validRunTarget, nil
		},
		writeFile: func(path string, data []byte, perm os.FileMode) error {
			writeCalled = true
			return nil
		},
		prompt: func(existingNames map[string]struct{}) (repositoryInput, error) {
			return repositoryInput{Name: "new-repo"}, nil
		},
		generateKey: func(existingKeys map[string]struct{}) (string, error) {
			return "repo_bbbbbbbb", nil
		},
		confirm: func(targetPath string, key string, input repositoryInput) (bool, error) {
			return false, nil
		},
		terraformFmt: func(targetPath string) error {
			fmtCalled = true
			return nil
		},
		stdout: &stdout,
	})
	if err != nil {
		t.Fatalf("runWithDeps returned error: %v", err)
	}
	if got := stdout.String(); got != "No changes made.\n" {
		t.Fatalf("stdout = %q, want cancellation message", got)
	}
	if writeCalled {
		t.Fatal("writeFile was called after rejected confirmation")
	}
	if fmtCalled {
		t.Fatal("terraformFmt was called after rejected confirmation")
	}
}

func TestRunWithDepsReturnsFmtErrorAfterWrite(t *testing.T) {
	t.Parallel()

	fmtErr := errors.New("terraform fmt failed")
	writeCalled := false

	err := runWithDeps([]string{"locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return validRunTarget, nil
		},
		writeFile: func(path string, data []byte, perm os.FileMode) error {
			writeCalled = true
			return nil
		},
		prompt: func(existingNames map[string]struct{}) (repositoryInput, error) {
			return repositoryInput{Name: "new-repo"}, nil
		},
		generateKey: func(existingKeys map[string]struct{}) (string, error) {
			return "repo_bbbbbbbb", nil
		},
		confirm: func(targetPath string, key string, input repositoryInput) (bool, error) {
			return true, nil
		},
		terraformFmt: func(targetPath string) error {
			return fmtErr
		},
		stdout: &bytes.Buffer{},
	})
	if !errors.Is(err, fmtErr) {
		t.Fatalf("error = %v, want fmt error", err)
	}
	if !writeCalled {
		t.Fatal("writeFile was not called before terraformFmt error")
	}
}

func TestValidateRepositoryName(t *testing.T) {
	t.Parallel()

	existing := map[string]struct{}{"existing-repo": {}}

	tests := []struct {
		name    string
		wantErr string
	}{
		{name: "new-repo"},
		{name: "repo.name_1"},
		{name: "", wantErr: "repository name is required"},
		{name: "   ", wantErr: "repository name is required"},
		{name: "existing-repo", wantErr: "repository already exists"},
		{name: "contains space", wantErr: "repository name must be 1-100 characters and contain only letters, numbers, dots, underscores, or hyphens"},
		{name: strings.Repeat("a", 101), wantErr: "repository name must be 1-100 characters and contain only letters, numbers, dots, underscores, or hyphens"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRepositoryName(tt.name, existing)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateRepositoryName returned unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("validateRepositoryName returned nil error, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestGenerateRepositoryKey(t *testing.T) {
	t.Parallel()

	existing := map[string]struct{}{}
	key, err := generateRepositoryKey(existing)
	if err != nil {
		t.Fatalf("generateRepositoryKey returned error: %v", err)
	}

	if !regexp.MustCompile(`^repo_[0-9a-f]{8}$`).MatchString(key) {
		t.Fatalf("key %q does not match repo_ + 8 lowercase hex digits", key)
	}
}

func TestRenderRepositoryEntryMinimal(t *testing.T) {
	t.Parallel()

	got := renderRepositoryEntry("repo_ab12cd34", repositoryInput{Name: "new-repo"})
	want := strings.TrimSpace(`
repo_ab12cd34 = {
  name = "new-repo"
}`)

	if got != want {
		t.Fatalf("entry mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestRenderRepositoryEntryWithOptions(t *testing.T) {
	t.Parallel()

	got := renderRepositoryEntry("repo_ab12cd34", repositoryInput{
		Name:            "new-repo",
		EnableCodeQL:    true,
		GitHubAppTokens: []string{"pr_approver", "pr_creator"},
		PATTokens:       []string{"metrics"},
	})
	want := strings.TrimSpace(`
repo_ab12cd34 = {
  name          = "new-repo"
  enable_codeql = true
  github_app_tokens = {
    pr_creator  = local.github_app_token_presets.pr_creator
    pr_approver = local.github_app_token_presets.pr_approver
  }
  pat_tokens = {
    metrics = local.pat_token_presets.metrics
  }
}`)

	if got != want {
		t.Fatalf("entry mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestReadExistingRepositories(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_aaaaaaaa = {
      name = "existing-repo"
    }
  }
}`)

	existing, err := readExistingRepositories(src, "locals.tf")
	if err != nil {
		t.Fatalf("readExistingRepositories returned error: %v", err)
	}

	if _, ok := existing.Keys["repo_aaaaaaaa"]; !ok {
		t.Fatalf("missing parsed key")
	}
	if _, ok := existing.Names["existing-repo"]; !ok {
		t.Fatalf("missing parsed repository name")
	}
}

func TestAppendRepository(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  github_app_token_presets = {
    pr_creator = {}
  }

  repositories = {
    repo_aaaaaaaa = {
      name = "existing-repo"
    }
  }
}`)

	updated, err := appendRepository(src, "locals.tf", "repo_bbbbbbbb", repositoryInput{Name: "new-repo"})
	if err != nil {
		t.Fatalf("appendRepository returned error: %v", err)
	}

	formatted := string(updated)
	if !strings.Contains(formatted, `repo_aaaaaaaa = {`) {
		t.Fatalf("existing repository entry was not preserved:\n%s", formatted)
	}
	if !strings.Contains(formatted, `repo_bbbbbbbb = {`) {
		t.Fatalf("new repository key was not appended:\n%s", formatted)
	}
	if !strings.Contains(formatted, `name = "new-repo"`) {
		t.Fatalf("new repository name was not appended:\n%s", formatted)
	}

	existing, err := readExistingRepositories(updated, "locals.tf")
	if err != nil {
		t.Fatalf("readExistingRepositories returned error for updated HCL: %v", err)
	}
	if _, ok := existing.Keys["repo_bbbbbbbb"]; !ok {
		t.Fatalf("new repository key was not parsed from locals.repositories:\n%s", formatted)
	}
	if _, ok := existing.Names["new-repo"]; !ok {
		t.Fatalf("new repository name was not parsed from locals.repositories:\n%s", formatted)
	}
}

func TestAppendRepositoryRejectsNonObjectRepositories(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = merge({}, {})
}`)

	_, err := appendRepository(src, "locals.tf", "repo_bbbbbbbb", repositoryInput{Name: "new-repo"})
	if err == nil {
		t.Fatal("appendRepository returned nil error")
	}
	if !strings.Contains(err.Error(), "locals.repositories must be an object") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestAppendRepositoryUsesLocalsBlockWithRepositories(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  unrelated = true
}

locals {
  repositories = {
    repo_aaaaaaaa = {
      name = "existing-repo"
    }
  }
}`)

	updated, err := appendRepository(src, "locals.tf", "repo_bbbbbbbb", repositoryInput{Name: "new-repo"})
	if err != nil {
		t.Fatalf("appendRepository returned error: %v", err)
	}

	existing, err := readExistingRepositories(updated, "locals.tf")
	if err != nil {
		t.Fatalf("readExistingRepositories returned error: %v", err)
	}
	if _, ok := existing.Keys["repo_bbbbbbbb"]; !ok {
		t.Fatalf("missing appended key")
	}
	if _, ok := existing.Names["new-repo"]; !ok {
		t.Fatalf("missing appended repository name")
	}
	if !strings.Contains(string(updated), "unrelated = true") {
		t.Fatalf("first locals block was not preserved:\n%s", string(updated))
	}
}
