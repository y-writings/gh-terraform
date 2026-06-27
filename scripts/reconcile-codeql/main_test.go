package main

import (
	"bytes"
	"context"
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

type failingWriter struct{}

func (failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

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

func TestDryRunLineDescribesPackageReconcileOperation(t *testing.T) {
	t.Parallel()

	got := dryRunLine("y-writings", codeqlTarget{Name: "repo", Languages: []string{"go", "python"}})
	want := "reconcile CodeQL default setup --owner y-writings --repo repo --languages go,python"

	if got != want {
		t.Fatalf("dryRunLine = %q, want %q", got, want)
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
		reconcileTarget: func(ctx context.Context, owner string, target codeqlTarget) (reconcileOutput, error) {
			runCalled = true
			return reconcileOutput{}, nil
		},
		stdout: &stdout,
	})
	if err != nil {
		t.Fatalf("runWithDeps returned error: %v", err)
	}
	if runCalled {
		t.Fatal("runCommand was called during dry-run")
	}

	want := "reconcile CodeQL default setup --owner y-writings --repo two-repo --languages go,python\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunWithDepsPrintsReconcileOutput(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	err := runWithDeps([]string{"--owner", "y-writings", "--repo", "two-repo", "locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return multiTargetConfig, nil
		},
		reconcileTarget: func(ctx context.Context, owner string, target codeqlTarget) (reconcileOutput, error) {
			return reconcileOutput{Owner: owner, Repo: target.Name, Changed: true}, nil
		},
		stdout: &stdout,
	})
	if err != nil {
		t.Fatalf("runWithDeps returned error: %v", err)
	}

	for _, want := range []string{`"owner": "y-writings"`, `"repo": "two-repo"`, `"changed": true`} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q, want to contain %q", stdout.String(), want)
		}
	}
}

func TestWriteReconcileOutputReturnsWriterError(t *testing.T) {
	t.Parallel()

	err := writeReconcileOutput(failingWriter{}, reconcileOutput{Owner: "y-writings", Repo: "repo"})
	if err == nil {
		t.Fatal("writeReconcileOutput returned nil error")
	}
	if !strings.Contains(err.Error(), "write failed") {
		t.Fatalf("error = %q, want writer error", err.Error())
	}
}

func TestRunWithDepsAttemptsAllTargetsAndReportsFailures(t *testing.T) {
	t.Parallel()

	var targets []codeqlTarget
	err := runWithDeps([]string{"--owner", "y-writings", "locals.tf"}, runDeps{
		readFile: func(path string) ([]byte, error) {
			return multiTargetConfig, nil
		},
		reconcileTarget: func(ctx context.Context, owner string, target codeqlTarget) (reconcileOutput, error) {
			if owner != "y-writings" {
				t.Fatalf("owner = %q", owner)
			}
			targets = append(targets, target)
			if target.Name == "one-repo" {
				return reconcileOutput{}, errors.New("boom")
			}
			return reconcileOutput{}, nil
		},
		stdout: &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("runWithDeps returned nil error")
	}
	if len(targets) != 2 {
		t.Fatalf("targets = %d, want 2", len(targets))
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
		reconcileTarget: func(ctx context.Context, owner string, target codeqlTarget) (reconcileOutput, error) {
			t.Fatal("reconcileTarget must not be called for a missing repo filter")
			return reconcileOutput{}, nil
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

func TestReconcileTargetRejectsInvalidLanguagesBeforeCreatingClient(t *testing.T) {
	t.Parallel()

	clientCalled := false
	reconcile := newReconcileTarget(func() (restClient, error) {
		clientCalled = true
		return nil, errors.New("client should not be created")
	})

	_, err := reconcile(context.Background(), "y-writings", codeqlTarget{Name: "repo", Languages: []string{"invalid"}})
	if err == nil {
		t.Fatal("reconcile returned nil error")
	}
	if !strings.Contains(err.Error(), "languages must contain only") {
		t.Fatalf("error = %q, want language validation error", err.Error())
	}
	if clientCalled {
		t.Fatal("client factory was called before input validation")
	}
}
