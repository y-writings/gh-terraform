# gh-terraform

personal account 配下の複数 GitHub リポジトリに、同一の統制ルールを Terraform で適用するための構成です。

## ディレクトリ構成

```text
.
├── terraform/
│   ├── main.tf                    # provider / module 呼び出し
│   ├── locals.tf                  # 管理対象リポジトリ一覧
│   ├── versions.tf                # Terraform / GitHub provider のバージョン制約
│   ├── modules/
│   │   ├── governance/
│   │   │   ├── main.tf            # github_repository_ruleset と共通 ruleset baseline
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── versions.tf
│   │   └── repository/
│   │       ├── main.tf            # fixed github_repository baseline
│   │       ├── variables.tf
│   │       ├── outputs.tf
│   │       └── versions.tf
├── .github/
├── .commitlintrc.yaml
├── lefthook.yaml
└── release-please-config.json
```

## この構成で管理するもの

- personal account (`y-writings`) 配下の複数リポジトリ
- 各 repo 共通の `main-default` repository ruleset
  - `~DEFAULT_BRANCH` に対する creation / update / deletion 制限
  - linear history 必須
  - force push 防止
  - pull request 必須
- `github_repository` による repo 設定
  - fixed public visibility
  - issues / wiki / merge method
  - module-owned governance baseline

### 現在採用している統一 baseline

- Repository ruleset (`main-default`)
  - `~DEFAULT_BRANCH`
  - Pull request 必須
  - 承認 1 件必須
  - squash merge のみ許可
  - stale review dismissal 有効
  - review thread resolution 必須
  - creation / update / deletion 制限
  - linear history 必須
  - force push 防止
  - CodeQL 必須 (`errors_and_warnings` / `high_or_higher`)
- Security baseline
  - Dependabot alerts: `enabled`
  - Secret scanning: `enabled`
  - Push protection: `enabled`
  - `terraform/modules/repository` が保持し、root input では上書きしない
- General repository baseline
  - Wiki: `enabled`
  - Issues: `enabled`
  - Delete branch on merge: `enabled`
  - Merge commit: `disabled`
  - Squash merge: `enabled`
  - Squash merge commit title: `PR_TITLE`
  - Squash merge commit message: `PR_BODY`
  - Rebase merge: `disabled`
  - `terraform/modules/repository` が保持し、root input では上書きしない

## personal account 前提での設計方針

- `github_organization_ruleset` は使いません
- `modules/repository` に fixed repository baseline を、`modules/governance` に repo ごとの ruleset baseline を保持します
- 将来 repo を追加するときは `terraform/locals.tf` の `repositories` map に 1 エントリ追加するだけです
- shared repository baseline と advanced security baseline は `terraform/modules/repository` に固定で保持します
- governance のルール内容は module 内に固定で保持し、root input では変更できません
- 既存 repo や ruleset を state に載せる必要がある場合は Terraform CLI の `import` を使います

## 管理対象リポジトリ

```hcl
repositories = {
  dotfiles         = {}
  templates        = {}
  container        = {}
  karabiner-config = {}
  gh-terraform     = {}
}
```

`terraform/main.tf` では provider の owner を `y-writings` に固定し、`terraform/locals.tf` で管理対象 repo 一覧を保持します。

### repositories の意味

- shared repository baseline (`has_wiki` / `has_issues` / merge method defaults / `delete_branch_on_merge`) は `terraform/modules/repository` に固定で保持します
- advanced security baseline (`Dependabot alerts` / `Secret scanning` / `Push protection`) は `terraform/modules/repository` に固定で保持します
- `visibility`: `terraform/modules/repository` に固定で保持され、root では指定しません
- `repositories`: 管理対象 repo 名を key にした map です。value は空 object (`{}`) を指定します

## 使い方

1. `terraform/locals.tf` の `repositories` を必要に応じて調整します
2. `GITHUB_TOKEN` を設定します
3. 既存 repo や既存 `main-default` ruleset を state に載せる必要がある場合は Terraform CLI の `import` を使います
4. `terraform plan` で差分を確認し、問題なければ apply します

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform plan
terraform -chdir=terraform apply
```

## 既存 state の移行

- repository 本体や既存 `main-default` ruleset を state に載せる場合は Terraform CLI の `import` を使います
- baseline は fully module-owned のため、既存の外部入力や CI で baseline の個別設定を渡している場合は削除してください
- この構成では governance の repo ごと override は持たないため、plan では統一 baseline への収束差分が出ます

## provider 上の制約メモ

- `github_repository_ruleset` は repo 単位 resource のため、personal account では `for_each` で横展開します
- `github_repository_ruleset` の import は `<repository>:<ruleset_id>` 形式です
- `evaluate` enforcement は organization ruleset 向けであり、personal account では `active` / `disabled` を使います
- live baseline の取り込み時は、既存 ruleset を import してから共通 baseline に寄せます
- GitHub ruleset の `code_quality` ルールは現行 Terraform provider では表現できません。既存 repo にある場合は manual drift として扱い、apply 前後で GitHub 側確認が必要です
- 現在の module-owned governance baseline は public repository 前提です。personal account では managed repositories を public 前提で扱ってください
