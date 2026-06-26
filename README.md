# gh-terraform

personal account (`y-writings`) 配下の GitHub リポジトリを Terraform で管理します。

## Requirements

- Terraform `>= 1.7.0`
- `mise`
- `GITHUB_TOKEN`
- GCS backend bucket
- 1Password CLI `op`

GCS backend bucket: `scripts/create-bucket/`

## Structure

```text
terraform/
├── work-repositories/       # 通常の作業リポジトリ用 root module
├── fork-repositories/       # fork リポジトリ用 root module
└── modules/
    ├── repository/          # repository 共通設定
    ├── governance/          # 通常 repo の ruleset
    ├── github_actions_credentials/
    └── fork_governance/     # fork repo 向けの最小 governance
```

## Initialize

```bash
cp terraform/work-repositories/terraform.tfvars.example terraform/work-repositories/terraform.tfvars
op signin
mise run tfinit
```

## Add Work Repo

```bash
mise run work:add-repo
```

When CodeQL languages are selected, the repository entry gets a `codeql.languages` block. That block enables the CodeQL-required ruleset in Terraform and is also used by the CodeQL reconcile task.

## Reconcile CodeQL

```bash
mise run work:codeql
```

`work:codeql` reads `terraform/work-repositories/locals.tf` and runs `gh-usecase codeql-default-setup` for repositories with `codeql.languages`.

This task is independent of Terraform diff. It runs only when invoked, even if `terraform plan` has no changes. The underlying `gh-usecase` command reads the current GitHub CodeQL default setup first and only PATCHes when GitHub differs from the desired languages.

## Check

```bash
mise run tfcheck
```

## Plan / Apply

work repo:

```bash
export GITHUB_TOKEN="<your-token>"
mise run work:plan
mise run work:apply
mise run work:codeql
```

fork repo:

```bash
export GITHUB_TOKEN="<your-token>"
mise run fork:plan
mise run fork:apply
```
