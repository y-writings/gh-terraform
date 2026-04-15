# gh-terraform

personal account 配下の複数 GitHub リポジトリに、同一の統制ルールを Terraform で適用するための構成です。

## ディレクトリ構成

```text
.
├── terraform/
│   ├── main.tf                    # provider / import / module 呼び出し
│   ├── locals.tf                  # リポジトリ metadata と import 対象を導出
│   ├── variables.tf               # personal account 向けの入力変数
│   ├── versions.tf                # Terraform / GitHub provider のバージョン制約
│   ├── terraform.tfvars.example   # 5 リポジトリ管理のサンプル
│   ├── modules/
│   │   └── governance/
│   │       ├── main.tf            # github_repository_ruleset と共通 baseline
│   │       ├── variables.tf
│   │       ├── outputs.tf
│   │       └── versions.tf
├── .github/
├── .commitlintrc.yaml
├── lefthook.yaml
└── release-please-config.json
```

## この構成で管理するもの

- personal account (`github_owner`) 配下の複数リポジトリ
- 各 repo 共通の `main-default` repository ruleset
  - `~DEFAULT_BRANCH` に対する creation / update / deletion 制限
  - linear history 必須
  - force push 防止
  - pull request 必須
- `github_repository` による repo 設定
  - visibility
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
  - `terraform/modules/governance` が保持し、root input では上書きしない
- General repository baseline
  - Wiki: `enabled`
  - Issues: `enabled`
  - Delete branch on merge: `disabled`
  - Merge commit: `disabled`
  - Squash merge: `enabled`
  - Squash merge commit title: `PR_TITLE`
  - Squash merge commit message: `PR_BODY`
  - Rebase merge: `disabled`
  - `terraform/modules/governance` が保持し、root input では上書きしない

## personal account 前提での設計方針

- `github_organization_ruleset` は使いません
- `modules/governance` に共通 baseline と repo ごとの ruleset を統合します
- 将来 repo を追加するときは `repositories` map に 1 エントリ追加するだけです
- shared repository baseline と advanced security baseline は `terraform/modules/governance` に固定で保持します
- governance のルール内容は module 内に固定で保持し、root input では変更できません
- 既存 repo の import は repo ごとの `import_existing_repository` で制御します
- 既存 `main-default` ruleset を import したいときだけ `main_default_ruleset_id` を設定します

## 入力例

```hcl
github_owner = "y-writings"

repositories = {
  dotfiles = {
    visibility                 = "public"
    import_existing_repository = true
    main_default_ruleset_id    = "14686753"
  }

  templates = {
    visibility                 = "public"
    import_existing_repository = true
    main_default_ruleset_id    = "14702457"
  }

  container = {
    visibility                 = "public"
    import_existing_repository = true
  }

  karabiner-config = {
    visibility                 = "public"
    import_existing_repository = true
    main_default_ruleset_id    = "14687410"
  }

  gh-terraform = {
    visibility                 = "public"
    import_existing_repository = true
    main_default_ruleset_id    = "14959390"
  }
}
```

### 各属性の意味

- shared repository baseline (`has_wiki` / `has_issues` / merge method defaults / `delete_branch_on_merge`) は `terraform/modules/governance` に固定で保持します
- advanced security baseline (`Dependabot alerts` / `Secret scanning` / `Push protection`) は `terraform/modules/governance` に固定で保持します
- `visibility`: managed repositories は `public` のみです。省略した場合も `public` として扱います
- `import_existing_repository`: 既存 repo を初回 plan/apply 時に import するかどうか
- `main_default_ruleset_id`: 既存 `main-default` ruleset を import したい場合の ruleset ID

## 使い方

1. `terraform/terraform.tfvars.example` を `terraform/terraform.tfvars` にコピーして調整します
2. `GITHUB_TOKEN` を設定します
3. 既存 repo を取り込む場合は対象 repo の `import_existing_repository = true` を設定します
4. 既存 `main-default` ruleset も state に載せたい場合は `main_default_ruleset_id` を設定します
5. `terraform plan` で差分を確認し、問題なければ apply します

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform plan
terraform -chdir=terraform apply
```

## 既存 state の移行

- repository 本体は import block で state へ取り込みます
- 既存 `main-default` ruleset を持つ repo は `main_default_ruleset_id` を設定して import します
- baseline は fully module-owned のため、既存の `terraform.tfvars` や CI で baseline の個別設定を渡している場合は削除してください
- この構成では governance の repo ごと override は持たないため、plan では統一 baseline への収束差分が出ます

## provider 上の制約メモ

- `github_repository_ruleset` は repo 単位 resource のため、personal account では `for_each` で横展開します
- `github_repository_ruleset` の import は `<repository>:<ruleset_id>` 形式です
- `evaluate` enforcement は organization ruleset 向けであり、personal account では `active` / `disabled` を使います
- live baseline の取り込み時は、既存 ruleset を import してから共通 baseline に寄せます
- GitHub ruleset の `code_quality` ルールは現行 Terraform provider では表現できません。既存 repo にある場合は manual drift として扱い、apply 前後で GitHub 側確認が必要です
- 現在の module-owned governance baseline は public repository 前提です。personal account では managed repositories を public 前提で扱ってください
