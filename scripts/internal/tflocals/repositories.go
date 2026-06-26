package tflocals

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Repository struct {
	Key  string
	Name string
}

type CodeQLTarget struct {
	Name      string
	Languages []string
}

func ReadRepositories(src []byte, filename string) ([]Repository, error) {
	repositories, err := repositoriesObjectExpression(src, filename)
	if err != nil {
		return nil, err
	}

	result := make([]Repository, 0, len(repositories.Items))
	for _, item := range repositories.Items {
		key, ok := staticObjectKey(item.KeyExpr)
		if !ok {
			continue
		}

		repository := Repository{Key: key}
		value, ok := item.ValueExpr.(*hclsyntax.ObjectConsExpr)
		if ok {
			repository.Name, _ = repositoryNameValue(value)
		}
		result = append(result, repository)
	}

	return result, nil
}

func ReadCodeQLTargets(src []byte, filename string) ([]CodeQLTarget, error) {
	repositories, err := repositoriesObjectExpression(src, filename)
	if err != nil {
		return nil, err
	}

	var targets []CodeQLTarget
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
			targets = append(targets, CodeQLTarget{Name: name, Languages: codeql})
			continue
		}

		if enabled, ok := repositoryEnableCodeQL(repository); ok && enabled {
			return nil, fmt.Errorf("repository %s uses legacy enable_codeql without codeql.languages", name)
		}
	}

	return targets, nil
}

func repositoriesObjectExpression(src []byte, filename string) (*hclsyntax.ObjectConsExpr, error) {
	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", filename, diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("parse %s: unexpected HCL body", filename)
	}

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
