package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/y-writings/gh-terraform/scripts/internal/tflocals"
)

type codeqlTarget = tflocals.CodeQLTarget

type commandInvocation struct {
	Name string
	Args []string
}

type runDeps struct {
	readFile   func(string) ([]byte, error)
	runCommand func(commandInvocation) error
	stdout     io.Writer
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
		readFile:   os.ReadFile,
		runCommand: runShellCommand,
		stdout:     os.Stdout,
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
	for _, target := range targets {
		if options.Repo != "" && target.Name != options.Repo {
			continue
		}
		matched++

		invocation := buildInvocation(options.Owner, target)
		if options.DryRun {
			fmt.Fprintln(deps.stdout, invocation.String())
			continue
		}
		if err := deps.runCommand(invocation); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", target.Name, err))
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

func buildInvocation(owner string, target codeqlTarget) commandInvocation {
	return commandInvocation{
		Name: "gh-usecase",
		Args: []string{
			"codeql-default-setup",
			"--owner", owner,
			"--repo", target.Name,
			"--languages", strings.Join(target.Languages, ","),
		},
	}
}

func runShellCommand(invocation commandInvocation) error {
	cmd := exec.Command(invocation.Name, invocation.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (i commandInvocation) String() string {
	return strings.Join(append([]string{i.Name}, i.Args...), " ")
}

func readTargets(src []byte, filename string) ([]codeqlTarget, error) {
	return tflocals.ReadCodeQLTargets(src, filename)
}
