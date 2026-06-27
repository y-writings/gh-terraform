package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/y-writings/gh-terraform/scripts/internal/tflocals"
	"github.com/y-writings/gh-usecase/codeqldefaultsetup"
)

type codeqlTarget = tflocals.CodeQLTarget
type reconcileOutput = codeqldefaultsetup.Output
type restClient = codeqldefaultsetup.Client

type restClientFactory func() (restClient, error)
type reconcileTargetFunc func(context.Context, string, codeqlTarget) (reconcileOutput, error)

type runDeps struct {
	readFile        func(string) ([]byte, error)
	reconcileTarget reconcileTargetFunc
	stdout          io.Writer
}

type runOptions struct {
	Owner      string
	Repo       string
	DryRun     bool
	TargetPath string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	return runWithDeps(args, runDeps{
		readFile:        os.ReadFile,
		reconcileTarget: newReconcileTarget(newRESTClient),
		stdout:          os.Stdout,
	})
}

func runWithDeps(args []string, deps runDeps) error {
	options, err := parseArgs(args)
	if err != nil {
		return err
	}

	src, err := deps.readFile(options.TargetPath)
	if err != nil {
		return fmt.Errorf("read target file: %w", err)
	}

	targets, err := readTargets(src, options.TargetPath)
	if err != nil {
		return err
	}

	var failures []string
	matched := 0
	ctx := context.Background()
	for _, target := range targets {
		if options.Repo != "" && target.Name != options.Repo {
			continue
		}
		matched++

		if options.DryRun {
			fmt.Fprintln(deps.stdout, dryRunLine(options.Owner, target))
			continue
		}
		output, err := deps.reconcileTarget(ctx, options.Owner, target)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", target.Name, err))
			continue
		}
		if err := writeReconcileOutput(deps.stdout, output); err != nil {
			return err
		}
	}
	if options.Repo != "" && matched == 0 {
		return fmt.Errorf("no CodeQL target found for repository %s", options.Repo)
	}
	if len(failures) > 0 {
		return fmt.Errorf("codeql reconcile failed: %s", strings.Join(failures, "; "))
	}

	return nil
}

func parseArgs(args []string) (runOptions, error) {
	var options runOptions
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--owner":
			i++
			if i >= len(args) || args[i] == "" {
				return runOptions{}, fmt.Errorf("--owner requires a value")
			}
			options.Owner = args[i]
		case "--repo":
			i++
			if i >= len(args) || args[i] == "" {
				return runOptions{}, fmt.Errorf("--repo requires a value")
			}
			options.Repo = args[i]
		case "--dry-run":
			options.DryRun = true
		default:
			if strings.HasPrefix(args[i], "-") {
				return runOptions{}, fmt.Errorf("unknown option %s", args[i])
			}
			if options.TargetPath != "" {
				return runOptions{}, fmt.Errorf("usage: reconcile-codeql --owner <owner> [--repo <repo>] [--dry-run] <path-to-locals.tf>")
			}
			options.TargetPath = args[i]
		}
	}

	if options.Owner == "" {
		return runOptions{}, fmt.Errorf("--owner is required")
	}
	if options.TargetPath == "" {
		return runOptions{}, fmt.Errorf("usage: reconcile-codeql --owner <owner> [--repo <repo>] [--dry-run] <path-to-locals.tf>")
	}
	return options, nil
}

func dryRunLine(owner string, target codeqlTarget) string {
	return fmt.Sprintf(
		"reconcile CodeQL default setup --owner %s --repo %s --languages %s",
		owner,
		target.Name,
		strings.Join(target.Languages, ","),
	)
}

func writeReconcileOutput(stdout io.Writer, output reconcileOutput) error {
	encoded, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("encode reconcile output: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, string(encoded)); err != nil {
		return fmt.Errorf("write reconcile output: %w", err)
	}
	return nil
}

func newReconcileTarget(newClient restClientFactory) reconcileTargetFunc {
	return func(ctx context.Context, owner string, target codeqlTarget) (reconcileOutput, error) {
		input := codeqldefaultsetup.Input{
			Owner:     owner,
			Repo:      target.Name,
			Languages: target.Languages,
		}
		if _, err := codeqldefaultsetup.Validate(input); err != nil {
			return reconcileOutput{}, err
		}

		client, err := newClient()
		if err != nil {
			return reconcileOutput{}, err
		}

		return codeqldefaultsetup.Reconcile(ctx, client, input)
	}
}

func newRESTClient() (restClient, error) {
	return api.NewRESTClient(api.ClientOptions{
		Headers: map[string]string{
			"Accept":               "application/vnd.github+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
	})
}

func readTargets(src []byte, filename string) ([]codeqlTarget, error) {
	return tflocals.ReadCodeQLTargets(src, filename)
}
