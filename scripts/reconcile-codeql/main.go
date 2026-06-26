package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type codeqlTarget struct {
	Name      string
	Languages []string
}

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
	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", filename, diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("parse %s: unexpected HCL body", filename)
	}

	repositories, err := repositoriesObjectExpression(body)
	if err != nil {
		return nil, err
	}

	var targets []codeqlTarget
	for _, item := range repositories.Items {
		repositoryKey, _ := staticObjectKey(item.KeyExpr)
		repository, ok := item.ValueExpr.(*hclsyntax.ObjectConsExpr)
		if !ok {
			continue
		}

		name, hasName := repositoryNameValue(repository)
		if !hasName {
			name = repositoryKey
		}

		codeql, hasCodeQL, err := repositoryCodeQL(repository)
		if err != nil {
			return nil, err
		}
		if hasCodeQL {
			if !hasName {
				return nil, fmt.Errorf("repository %s has codeql but is missing name", repositoryKey)
			}
			targets = append(targets, codeqlTarget{Name: name, Languages: codeql})
			continue
		}

		if enabled, ok := repositoryEnableCodeQL(repository); ok && enabled {
			return nil, fmt.Errorf("repository %s uses legacy enable_codeql without codeql.languages", name)
		}
	}

	return targets, nil
}

func repositoriesObjectExpression(body *hclsyntax.Body) (*hclsyntax.ObjectConsExpr, error) {
	for _, block := range body.Blocks {
		if block.Type != "locals" || len(block.Labels) != 0 {
			continue
		}

		attr, ok := block.Body.Attributes["repositories"]
		if !ok {
			continue
		}
		repositories, ok := attr.Expr.(*hclsyntax.ObjectConsExpr)
		if !ok {
			return nil, fmt.Errorf("locals.repositories must be an object")
		}
		return repositories, nil
	}

	return nil, fmt.Errorf("locals.repositories is missing")
}

func staticObjectKey(expr hclsyntax.Expression) (string, bool) {
	switch expr := expr.(type) {
	case *hclsyntax.ObjectConsKeyExpr:
		return staticObjectKey(expr.Wrapped)
	case *hclsyntax.ScopeTraversalExpr:
		if len(expr.Traversal) != 1 {
			return "", false
		}
		root, ok := expr.Traversal[0].(hcl.TraverseRoot)
		return root.Name, ok
	case *hclsyntax.LiteralValueExpr:
		if expr.Val.Type() != cty.String {
			return "", false
		}
		return expr.Val.AsString(), true
	default:
		return "", false
	}
}

func repositoryNameValue(repository *hclsyntax.ObjectConsExpr) (string, bool) {
	for _, item := range repository.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok || key != "name" {
			continue
		}

		return staticStringValue(item.ValueExpr)
	}

	return "", false
}

func repositoryCodeQL(repository *hclsyntax.ObjectConsExpr) ([]string, bool, error) {
	for _, item := range repository.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok || key != "codeql" {
			continue
		}

		codeql, ok := item.ValueExpr.(*hclsyntax.ObjectConsExpr)
		if !ok {
			return nil, true, fmt.Errorf("codeql must be an object")
		}

		languages, ok, err := codeqlLanguages(codeql)
		if err != nil {
			return nil, true, err
		}
		if !ok {
			name, _ := repositoryNameValue(repository)
			return nil, true, fmt.Errorf("repository %s codeql.languages is required", name)
		}
		return languages, true, nil
	}

	return nil, false, nil
}

func codeqlLanguages(codeql *hclsyntax.ObjectConsExpr) ([]string, bool, error) {
	for _, item := range codeql.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok || key != "languages" {
			continue
		}

		languages, ok := staticStringList(item.ValueExpr)
		if !ok {
			return nil, true, fmt.Errorf("codeql.languages must be a static string list")
		}
		if len(languages) == 0 {
			return nil, true, fmt.Errorf("codeql.languages must not be empty")
		}
		for _, language := range languages {
			if strings.TrimSpace(language) == "" {
				return nil, true, fmt.Errorf("codeql.languages must not contain empty values")
			}
		}
		return languages, true, nil
	}

	return nil, false, nil
}

func repositoryEnableCodeQL(repository *hclsyntax.ObjectConsExpr) (bool, bool) {
	for _, item := range repository.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok || key != "enable_codeql" {
			continue
		}

		value, ok := staticBoolValue(item.ValueExpr)
		return value, ok
	}

	return false, false
}

func staticStringValue(expr hclsyntax.Expression) (string, bool) {
	switch expr := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if expr.Val.Type() != cty.String {
			return "", false
		}
		return expr.Val.AsString(), true
	case *hclsyntax.TemplateExpr:
		if len(expr.Parts) != 1 {
			return "", false
		}
		return staticStringValue(expr.Parts[0])
	default:
		return "", false
	}
}

func staticBoolValue(expr hclsyntax.Expression) (bool, bool) {
	literal, ok := expr.(*hclsyntax.LiteralValueExpr)
	if !ok || literal.Val.Type() != cty.Bool {
		return false, false
	}
	return literal.Val.True(), true
}

func staticStringList(expr hclsyntax.Expression) ([]string, bool) {
	tuple, ok := expr.(*hclsyntax.TupleConsExpr)
	if !ok {
		return nil, false
	}

	languages := make([]string, 0, len(tuple.Exprs))
	for _, expr := range tuple.Exprs {
		language, ok := staticStringValue(expr)
		if !ok {
			return nil, false
		}
		languages = append(languages, language)
	}
	return languages, true
}
