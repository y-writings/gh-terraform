# gh-terraform

personal account 配下の複数 GitHub リポジトリに、同一の統制ルールを Terraform で適用するための構成です。

## ディレクトリ構成

```text
.
├── terraform/
│   ├── main.tf                    # provider / import / module 呼び出し
│   ├── locals.tf                  # リポジトリ metadata と共通 governance のデフォルト値を導出
│   ├── moved.tf                   # 旧 single-repo state を安全に state から外す定義
│   ├── variables.tf               # personal account 向けの入力変数
│   ├── versions.tf                # Terraform / GitHub provider のバージョン制約
│   ├── terraform.tfvars.example   # 5 リポジトリ管理のサンプル
│   ├── modules/
│   │   ├── advanced_security/
│   │   │   └── outputs.tf         # module-owned advanced security baseline
│   │   ├── general/
│   │   │   ├── variables.tf       # repo-specific metadata inputs
│   │   │   └── outputs.tf         # module-owned shared repository baseline
│   │   └── rulesets/
│   │       ├── main.tf            # github_repository_ruleset
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
  - module-owned general repository baseline
  - module-owned advanced security baseline

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
  - `terraform/modules/advanced_security` が保持し、root input では上書きしない
- General repository baseline
  - Issues: `enabled`
  - Merge commit: `disabled`
  - Squash merge: `enabled`
  - Squash merge commit title: `PR_TITLE`
  - Squash merge commit message: `PR_BODY`
  - Rebase merge: `disabled`
  - `terraform/modules/general` が保持し、repo metadata (`has_wiki`, `delete_branch_on_merge`) だけを input として受け取る

## personal account 前提での設計方針

- `github_organization_ruleset` は使いません
- 共通統制は `modules/general` / `modules/advanced_security` / `modules/rulesets` を `for_each` で展開して適用します
- 将来 repo を追加するときは `repositories` map に 1 エントリ追加するだけです
- shared repository baseline は `terraform/modules/general` に固定で保持します
- 共通統制のうち ruleset / code scanning は `repository_governance` で一元管理し、advanced security baseline は `terraform/modules/advanced_security` に固定で保持します
- 既存 repo の import は repo ごとの `import_existing_repository` で制御します
- 既存 `main-default` ruleset を import したいときだけ `main_default_ruleset_id` を設定します
- `delete_branch_on_merge` と `has_wiki` は governance ではなく repo metadata として repo ごとに保持します

## 入力例

```hcl
github_owner = "y-writings"

repository_governance = {
  enable_required_code_scanning          = true
  ruleset_enforcement                    = "active"
  required_approving_review_count        = 1
  dismiss_stale_reviews_on_push          = true
  require_code_owner_review              = false
  require_last_push_approval             = false
  required_review_thread_resolution      = true
  required_code_scanning = {
    tool                      = "CodeQL"
    alerts_threshold          = "errors_and_warnings"
    security_alerts_threshold = "high_or_higher"
  }
}

repositories = {
  dotfiles = {
    visibility                 = "public"
    import_existing_repository = true
    delete_branch_on_merge     = true
    has_wiki                   = false
    main_default_ruleset_id    = "14686753"
  }

  templates = {
    visibility                 = "public"
    import_existing_repository = true
    delete_branch_on_merge     = true
    has_wiki                   = true
    main_default_ruleset_id    = "14702457"
  }

  container = {
    visibility                 = "public"
    import_existing_repository = true
    has_wiki                   = true
  }

  karabiner-config = {
    visibility                 = "public"
    import_existing_repository = true
    has_wiki                   = false
    main_default_ruleset_id    = "14687410"
  }

  gh-terraform = {
    visibility                 = "public"
    import_existing_repository = true
    has_wiki                   = true
    main_default_ruleset_id    = "14959390"
  }
}
```

### 各属性の意味

- `repository_governance`: 全リポジトリに一律適用する統制ルール
- `repository_governance.enable_required_code_scanning`: shared baseline の code scanning requirement を global に有効/無効化するスイッチ
- shared repository baseline (`has_issues` / merge method defaults) は `terraform/modules/general` に固定で保持します
- advanced security baseline (`Dependabot alerts` / `Secret scanning` / `Push protection`) は `terraform/modules/advanced_security` に固定で保持します
- `visibility`: `public` / `private`。personal account のため `internal` は扱いません。ただし現在の shared governance baseline は public repository 前提です
- `import_existing_repository`: 既存 repo を初回 plan/apply 時に import するかどうか
- `delete_branch_on_merge`: merge 後にブランチを削除するかどうか
- `has_wiki`: repo ごとの wiki 有効/無効
- `main_default_ruleset_id`: 既存 `main-default` ruleset を import したい場合の ruleset ID

## 使い方

1. `terraform/terraform.tfvars.example` を `terraform/terraform.tfvars` にコピーして調整します
2. `GITHUB_TOKEN` を設定します
3. `repository_governance` で統一したい ruleset / code scanning ルールを定義します
4. 既存 repo を取り込む場合は対象 repo の `import_existing_repository = true` を設定します
5. 既存 `main-default` ruleset も state に載せたい場合は `main_default_ruleset_id` を設定します
6. `terraform plan` で差分を確認し、問題なければ apply します

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform plan
terraform -chdir=terraform apply
```

## 既存 state の移行

- 旧 single-repo 構成の `github_repository.this` / `github_repository_ruleset.main_default` は `moved.tf` の `removed { destroy = false }` で state から安全に外します
- repository 本体は import block で state へ取り込みます
- 既存 `main-default` ruleset を持つ repo は `main_default_ruleset_id` を設定して import します
- `general` baseline は module-owned に変更したため、shared repository baseline を root module call から上書きすることはありません。`terraform.tfvars` 側の cleanup は不要で、`has_wiki` と `delete_branch_on_merge` だけ引き続き repo metadata として設定します
- `advanced_security` baseline は module-owned に変更したため、既存の `terraform.tfvars` や CI の `repository_governance` object から `manage_security_and_analysis` / `vulnerability_alerts` / `secret_scanning_status` / `secret_scanning_push_protection_status` を削除してください
- つまり migration は「hard-coded moved」ではなく「old root state を remove → current module addresses へ import」で repository 名に依存せず行います
- この構成では governance の repo ごと override は持たないため、plan では統一 baseline への収束差分が出ます

## provider 上の制約メモ

- `github_repository_ruleset` は repo 単位 resource のため、personal account では `for_each` で横展開します
- `github_repository_ruleset` の import は `<repository>:<ruleset_id>` 形式です
- `evaluate` enforcement は organization ruleset 向けであり、personal account では `active` / `disabled` を使います
- live baseline の取り込み時は、既存 ruleset を import してから共通 baseline に寄せます
- GitHub ruleset の `code_quality` ルールは現行 Terraform provider では表現できません。既存 repo にある場合は manual drift として扱い、apply 前後で GitHub 側確認が必要です
- 現在の module-owned `advanced_security` baseline と `repository_governance.required_code_scanning` は public repository 前提です。personal account では managed repositories を public 前提で扱ってください
