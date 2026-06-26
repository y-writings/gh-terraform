package main

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

var multiTargetConfig = []byte(`locals {
  repositories = {
    repo_one = {
      name = "one-repo"
      codeql = {
        languages = ["go"]
      }
    }
    repo_two = {
      name = "two-repo"
      codeql = {
        languages = ["go", "python"]
      }
    }
  }
}`)

func TestReadTargetsReturnsRepositoriesWithCodeQLLanguages(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_enabled = {
      name = "enabled-repo"
      codeql = {
        languages = ["go", "javascript-typescript"]
      }
    }
    repo_disabled = {
      name = "disabled-repo"
    }
  }
}`)

	got, err := readTargets(src, "locals.tf")
	if err != nil {
		t.Fatalf("readTargets returned error: %v", err)
	}

	want := []codeqlTarget{{Name: "enabled-repo", Languages: []string{"go", "javascript-typescript"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("targets mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestReadTargetsRejectsLegacyEnableCodeQLWithoutLanguages(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_legacy = {
      name          = "legacy-repo"
      enable_codeql = true
    }
  }
}`)

	_, err := readTargets(src, "locals.tf")
	if err == nil {
		t.Fatal("readTargets returned nil error")
	}
	if !strings.Contains(err.Error(), "legacy-repo") || !strings.Contains(err.Error(), "codeql.languages") {
		t.Fatalf("error = %q, want repository name and codeql.languages guidance", err.Error())
	}
}

func TestReadTargetsRejectsCodeQLWithoutLanguages(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_missing_languages = {
      name = "missing-languages-repo"
      codeql = {}
    }
  }
}`)

	_, err := readTargets(src, "locals.tf")
	if err == nil {
		t.Fatal("readTargets returned nil error")
	}
	if !strings.Contains(err.Error(), "missing-languages-repo") || !strings.Contains(err.Error(), "codeql.languages") {
		t.Fatalf("error = %q, want repository name and codeql.languages guidance", err.Error())
	}
}

func TestBuildInvocationUsesGhUsecaseCodeQLDefaultSetup(t *testing.T) {
	t.Parallel()

	got := buildInvocation("y-writings", codeqlTarget{Name: "repo", Languages: []string{"go", "python"}})
	want := commandInvocation{
		Name: "gh-usecase",
		Args: []string{
			"codeql-default-setup",
			"--owner", "y-writings",
			"--repo", "repo",
			"--languages", "go,python",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invocation mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRunWithDepsDryRunPrintsFilteredInvocationWithoutExecuting(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	runCalled := false

	err := runWithDeps([]string{"--owner", "y-writings", "--repo", "two-repo", "--dry-run", "locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			if path != "locals.tf" {
				t.Fatalf("path = %q", path)
			}
			return multiTargetConfig, nil
		},
		runCommand: func(invocation commandInvocation) error {
			runCalled = true
			return nil
		},
		stdout: &stdout,
	})
	if err != nil {
		t.Fatalf("runWithDeps returned error: %v", err)
	}
	if runCalled {
		t.Fatal("runCommand was called during dry-run")
	}

	want := "gh-usecase codeql-default-setup --owner y-writings --repo two-repo --languages go,python\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunWithDepsAttemptsAllTargetsAndReportsFailures(t *testing.T) {
	t.Parallel()

	var invocations []commandInvocation
	err := runWithDeps([]string{"--owner", "y-writings", "locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return multiTargetConfig, nil
		},
		runCommand: func(invocation commandInvocation) error {
			invocations = append(invocations, invocation)
			if invocation.Args[4] == "one-repo" {
				return errors.New("boom")
			}
			return nil
		},
		stdout: &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("runWithDeps returned nil error")
	}
	if len(invocations) != 2 {
		t.Fatalf("invocations = %d, want 2", len(invocations))
	}
	if !strings.Contains(err.Error(), "one-repo") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %q, want failed repository name and cause", err.Error())
	}
}

func TestRunWithDepsRejectsRepoFilterWithNoMatchingTarget(t *testing.T) {
	t.Parallel()

	err := runWithDeps([]string{"--owner", "y-writings", "--repo", "missing-repo", "locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return multiTargetConfig, nil
		},
		runCommand: func(invocation commandInvocation) error {
			t.Fatal("runCommand must not be called for a missing repo filter")
			return nil
		},
		stdout: &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("runWithDeps returned nil error")
	}
	if !strings.Contains(err.Error(), "missing-repo") || !strings.Contains(err.Error(), "no CodeQL target") {
		t.Fatalf("error = %q, want missing repo guidance", err.Error())
	}
}
