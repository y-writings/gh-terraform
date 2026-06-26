package tflocals

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadRepositoriesReturnsStaticKeysAndOptionalNames(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_named = {
      name = "named-repo"
      codeql = {}
    }
    repo_without_name = {}
  }
}`)

	got, err := ReadRepositories(src, "locals.tf")
	if err != nil {
		t.Fatalf("ReadRepositories returned error: %v", err)
	}

	want := []Repository{
		{Key: "repo_named", Name: "named-repo"},
		{Key: "repo_without_name"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("repositories mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestReadCodeQLTargetsReturnsRepositoriesWithLanguages(t *testing.T) {
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

	got, err := ReadCodeQLTargets(src, "locals.tf")
	if err != nil {
		t.Fatalf("ReadCodeQLTargets returned error: %v", err)
	}

	want := []CodeQLTarget{{Name: "enabled-repo", Languages: []string{"go", "javascript-typescript"}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("targets mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

func TestReadCodeQLTargetsRejectsLegacyEnableCodeQLWithoutLanguages(t *testing.T) {
	t.Parallel()

	src := []byte(`locals {
  repositories = {
    repo_legacy = {
      name          = "legacy-repo"
      enable_codeql = true
    }
  }
}`)

	_, err := ReadCodeQLTargets(src, "locals.tf")
	if err == nil {
		t.Fatal("ReadCodeQLTargets returned nil error")
	}
	if !strings.Contains(err.Error(), "legacy-repo") || !strings.Contains(err.Error(), "codeql.languages") {
		t.Fatalf("error = %q, want repository name and codeql.languages guidance", err.Error())
	}
}
