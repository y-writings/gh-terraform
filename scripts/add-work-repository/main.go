package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type repositoryInput struct {
	Name            string
	EnableCodeQL    bool
	GitHubAppTokens []string
	PATTokens       []string
}

type existingRepositories struct {
	Keys  map[string]struct{}
	Names map[string]struct{}
}

type runDeps struct {
	readFile     func(string) ([]byte, error)
	writeFile    func(string, []byte, os.FileMode) error
	prompt       func(map[string]struct{}) (repositoryInput, error)
	generateKey  func(map[string]struct{}) (string, error)
	confirm      func(string, string, repositoryInput) (bool, error)
	terraformFmt func(string) error
	stdout       io.Writer
}

var repositoryNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,100}$`)

var githubAppTokenOrder = []string{"pr_creator", "pr_approver"}

var patTokenOrder = []string{"metrics"}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	return runWithDeps(args, runDeps{
		readFile:     os.ReadFile,
		writeFile:    os.WriteFile,
		prompt:       promptRepositoryInput,
		generateKey:  generateRepositoryKey,
		confirm:      confirmAddition,
		terraformFmt: runTerraformFmt,
		stdout:       os.Stdout,
	})
}

func runWithDeps(args []string, deps runDeps) error {
	targetPath, err := parseArgs(args)
	if err != nil {
		return err
	}

	src, err := deps.readFile(targetPath)
	if err != nil {
		return fmt.Errorf("read target file: %w", err)
	}

	existing, err := readExistingRepositories(src, targetPath)
	if err != nil {
		return err
	}

	input, err := deps.prompt(existing.Names)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(deps.stdout, "No changes made.")
			return nil
		}
		return err
	}

	key, err := deps.generateKey(existing.Keys)
	if err != nil {
		return err
	}

	confirmed, err := deps.confirm(targetPath, key, input)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(deps.stdout, "No changes made.")
			return nil
		}
		return err
	}
	if !confirmed {
		fmt.Fprintln(deps.stdout, "No changes made.")
		return nil
	}

	updated, err := appendRepository(src, targetPath, key, input)
	if err != nil {
		return err
	}
	if err := deps.writeFile(targetPath, updated, 0o644); err != nil {
		return fmt.Errorf("write target file: %w", err)
	}
	if err := deps.terraformFmt(targetPath); err != nil {
		return err
	}

	fmt.Fprintf(deps.stdout, "Added %s for %q to %s\n", key, input.Name, targetPath)
	fmt.Fprintln(deps.stdout, "Next: mise run work:plan")
	return nil
}

func parseArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: add-work-repository <path-to-locals.tf>")
	}
	return args[0], nil
}

func validateRepositoryName(name string, existing map[string]struct{}) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("repository name is required")
	}
	if !repositoryNamePattern.MatchString(name) {
		return fmt.Errorf("repository name must be 1-100 characters and contain only letters, numbers, dots, underscores, or hyphens")
	}
	if _, ok := existing[name]; ok {
		return fmt.Errorf("repository already exists: %s", name)
	}
	return nil
}

func promptRepositoryInput(existingNames map[string]struct{}) (repositoryInput, error) {
	var input repositoryInput
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository name").
				CharLimit(100).
				Value(&input.Name).
				Validate(func(value string) error {
					return validateRepositoryName(strings.TrimSpace(value), existingNames)
				}),
			huh.NewConfirm().
				Title("Enable CodeQL").
				Value(&input.EnableCodeQL),
			huh.NewMultiSelect[string]().
				Title("GitHub App tokens").
				Options(
					huh.NewOption("pr_creator", "pr_creator"),
					huh.NewOption("pr_approver", "pr_approver"),
				).
				Value(&input.GitHubAppTokens),
			huh.NewMultiSelect[string]().
				Title("PAT tokens").
				Options(huh.NewOption("metrics", "metrics")).
				Value(&input.PATTokens),
		),
	)
	if err := form.Run(); err != nil {
		return repositoryInput{}, err
	}

	input.Name = strings.TrimSpace(input.Name)
	return input, nil
}

func confirmAddition(targetPath string, key string, input repositoryInput) (bool, error) {
	confirmed := false
	err := huh.NewConfirm().
		Title("Add repository?").
		Description(fmt.Sprintf("Target: %s\n\n%s", targetPath, renderRepositoryEntry(key, input))).
		Value(&confirmed).
		Run()
	return confirmed, err
}

func runTerraformFmt(targetPath string) error {
	cmd := exec.Command("terraform", "fmt", targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform fmt %s: %w\n%s", targetPath, err, strings.TrimSpace(string(output)))
	}
	return nil
}

func generateRepositoryKey(existing map[string]struct{}) (string, error) {
	for range 256 {
		buf := make([]byte, 4)
		if _, err := rand.Read(buf); err != nil {
			return "", fmt.Errorf("generate repository key: %w", err)
		}

		key := "repo_" + hex.EncodeToString(buf)
		if _, ok := existing[key]; !ok {
			return key, nil
		}
	}

	return "", fmt.Errorf("could not generate unique repository key after 256 attempts")
}

func readExistingRepositories(src []byte, filename string) (existingRepositories, error) {
	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return existingRepositories{}, fmt.Errorf("parse %s: %s", filename, diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return existingRepositories{}, fmt.Errorf("parse %s: unexpected HCL body", filename)
	}

	repositories, err := repositoriesObjectExpression(body)
	if err != nil {
		return existingRepositories{}, err
	}

	existing := existingRepositories{
		Keys:  map[string]struct{}{},
		Names: map[string]struct{}{},
	}
	for _, item := range repositories.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok {
			continue
		}
		existing.Keys[key] = struct{}{}

		value, ok := item.ValueExpr.(*hclsyntax.ObjectConsExpr)
		if !ok {
			continue
		}
		if name, ok := repositoryNameValue(value); ok {
			existing.Names[name] = struct{}{}
		}
	}

	return existing, nil
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

func appendRepository(src []byte, filename string, key string, input repositoryInput) ([]byte, error) {
	if _, err := readExistingRepositories(src, filename); err != nil {
		return nil, err
	}

	file, diags := hclwrite.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", filename, diags.Error())
	}

	block := repositoriesLocalsBlock(file.Body())
	if block == nil {
		return nil, fmt.Errorf("locals.repositories is missing")
	}

	attr := block.Body().GetAttribute("repositories")
	if attr == nil {
		return nil, fmt.Errorf("locals.repositories is missing")
	}

	tokens, err := appendObjectAttributeTokens(attr.Expr().BuildTokens(nil), buildRepositoryEntryTokens(key, input))
	if err != nil {
		return nil, fmt.Errorf("locals.repositories must be an object")
	}
	block.Body().SetAttributeRaw("repositories", tokens)

	return file.Bytes(), nil
}

func repositoriesLocalsBlock(body *hclwrite.Body) *hclwrite.Block {
	for _, block := range body.Blocks() {
		if block.Type() == "locals" && len(block.Labels()) == 0 && block.Body().GetAttribute("repositories") != nil {
			return block
		}
	}

	return nil
}

func appendObjectAttributeTokens(object hclwrite.Tokens, entry hclwrite.Tokens) (hclwrite.Tokens, error) {
	for i := len(object) - 1; i >= 0; i-- {
		if object[i].Type != hclsyntax.TokenCBrace {
			continue
		}

		updated := make(hclwrite.Tokens, 0, len(object)+len(entry))
		updated = append(updated, object[:i]...)
		updated = append(updated, entry...)
		updated = append(updated, object[i:]...)
		return updated, nil
	}

	return nil, fmt.Errorf("object closing brace not found")
}

func renderRepositoryEntry(key string, input repositoryInput) string {
	return strings.TrimSpace(string(hclwrite.Format(buildRepositoryEntryTokens(key, input).Bytes())))
}

func buildRepositoryEntryTokens(key string, input repositoryInput) hclwrite.Tokens {
	tokens := hclwrite.TokensForObject([]hclwrite.ObjectAttrTokens{
		{
			Name:  hclwrite.TokensForIdentifier(key),
			Value: buildRepositoryValueTokens(input),
		},
	})

	return tokens[1 : len(tokens)-1]
}

func buildRepositoryValueTokens(input repositoryInput) hclwrite.Tokens {
	attrs := []hclwrite.ObjectAttrTokens{
		{
			Name:  hclwrite.TokensForIdentifier("name"),
			Value: hclwrite.TokensForValue(cty.StringVal(input.Name)),
		},
	}

	if input.EnableCodeQL {
		attrs = append(attrs, hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier("enable_codeql"),
			Value: hclwrite.TokensForValue(cty.BoolVal(true)),
		})
	}

	if tokens := buildPresetTokenMap(input.GitHubAppTokens, githubAppTokenOrder, "github_app_token_presets"); len(tokens) > 0 {
		attrs = append(attrs, hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier("github_app_tokens"),
			Value: tokens,
		})
	}

	if tokens := buildPresetTokenMap(input.PATTokens, patTokenOrder, "pat_token_presets"); len(tokens) > 0 {
		attrs = append(attrs, hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier("pat_tokens"),
			Value: tokens,
		})
	}

	return hclwrite.TokensForObject(attrs)
}

func buildPresetTokenMap(selected []string, allowed []string, presetName string) hclwrite.Tokens {
	selectedSet := make(map[string]struct{}, len(selected))
	for _, name := range selected {
		selectedSet[name] = struct{}{}
	}

	attrs := make([]hclwrite.ObjectAttrTokens, 0, len(allowed))
	for _, name := range allowed {
		if _, ok := selectedSet[name]; !ok {
			continue
		}

		attrs = append(attrs, hclwrite.ObjectAttrTokens{
			Name: hclwrite.TokensForIdentifier(name),
			Value: hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{Name: "local"},
				hcl.TraverseAttr{Name: presetName},
				hcl.TraverseAttr{Name: name},
			}),
		})
	}

	if len(attrs) == 0 {
		return nil
	}

	return hclwrite.TokensForObject(attrs)
}
