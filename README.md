# gh-terraform

personal account 配下の複数 GitHub リポジトリに、同一の統制ルールを Terraform で適用するための構成です。

## ディレクトリ構成

```text
.
├── terraform/
│   ├── work-repositories/         # 自分の作業リポジトリ用 root module / state
│   │   ├── main.tf                # provider / module 呼び出し
│   │   ├── locals.tf              # 管理対象リポジトリ一覧
│   │   ├── versions.tf            # Terraform / provider のバージョン制約
│   │   └── terraform.tfvars.example
│   ├── fork-repositories/         # fork リポジトリ用 root module / state
│   │   ├── main.tf                # provider / fork repository resource
│   │   ├── locals.tf              # fork 元リポジトリ一覧
│   │   └── versions.tf            # Terraform / provider のバージョン制約
│   ├── modules/
│   │   ├── governance/
│   │   │   ├── main.tf            # github_repository_ruleset と共通 ruleset baseline
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── versions.tf
│   │   ├── release_please/
│   │   │   ├── main.tf            # release-please 用の Actions secret / variable baseline
│   │   │   ├── variables.tf
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
  - GitHub Actions workflow permissions
  - GitHub Actions permissions
  - release-please 用の共通 Actions secret / variable
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
  - Default workflow permissions: `write`
  - Allow GitHub Actions to create and approve pull requests: `enabled`
  - Actions permissions: `enabled`
  - Allowed actions and reusable workflows: `all`
  - Require actions to be pinned to a full-length commit SHA: `enabled`
  - `terraform/modules/repository` が保持し、root input では上書きしない

## root module 分割方針

- `terraform/work-repositories` は自分の作業リポジトリ用の root module です
- `terraform/fork-repositories` は fork リポジトリ用の root module です
- 共有 module は `terraform/modules` に置き、root module から相対パスで参照します
- root module ごとに Terraform の実行ディレクトリと state を分けます
- fork リポジトリ管理では、fork 作成自体を `github_repository` の `fork = true` / `source_owner` / `source_repo` で管理します

## personal account 前提での設計方針

- `github_organization_ruleset` は使いません
- `modules/repository` に fixed repository baseline を、`modules/governance` に repo ごとの ruleset baseline を保持します
- `modules/release_please` に release-please 用の共通 Actions secret / variable baseline を保持します
- 将来 repo を追加するときは、原則として `terraform/work-repositories/locals.tf` の `repositories` map に 1 エントリ追加します
- ただし `sha_pinning_required = true` を baseline として適用するため、既存 workflow が tag / branch 参照の Action を使っている repo は事前に full-length commit SHA 参照へ寄せる必要があります
- shared repository baseline と advanced security baseline は `terraform/modules/repository` に固定で保持します
- GitHub Actions workflow permissions baseline も `terraform/modules/repository` に固定で保持します
- GitHub Actions permissions baseline も `terraform/modules/repository` に固定で保持します
- release-please 用の `RELEASE_PLEASE_APP_PRIVATE_KEY` / `RELEASE_PLEASE_APP_ID` は `terraform/modules/release_please` が 1Password から取得して各 repo に適用します
- 1Password 側では固定で `op://dev/release-please-bot/private key` と `op://dev/release-please-bot/info/app_id` に対応する値を参照します
- governance のルール内容は module 内に固定で保持し、root input では変更できません
- fork リポジトリ管理は最小構成とし、通常リポジトリ用の governance / release-please baseline は適用しません
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

`terraform/work-repositories/main.tf` では provider の owner を `y-writings` に固定し、`terraform/work-repositories/locals.tf` で管理対象 repo 一覧を保持します。

### repositories の意味

- shared repository baseline (`has_wiki` / `has_issues` / merge method defaults / `delete_branch_on_merge`) は `terraform/modules/repository` に固定で保持します
- advanced security baseline (`Dependabot alerts` / `Secret scanning` / `Push protection`) は `terraform/modules/repository` に固定で保持します
- GitHub Actions workflow permissions baseline (`default_workflow_permissions` / `can_approve_pull_request_reviews`) は `terraform/modules/repository` に固定で保持します
- GitHub Actions permissions baseline (`enabled` / `allowed_actions` / `sha_pinning_required`) は `terraform/modules/repository` に固定で保持します
- release-please 用の repository secret / variable baseline は `terraform/modules/release_please` に固定で保持します
- `visibility`: `terraform/modules/repository` に固定で保持され、root では指定しません
- `repositories`: 管理対象 repo 名を key にした map です。value は空 object (`{}`) を指定します
- `sha_pinning_required = true` により、managed repo の workflow で使う Action は full-length commit SHA で pin されている前提になります（reusable workflow の参照は GitHub の仕様上 tag 利用が残る場合があります）
- Artifact and log retention は GitHub API では設定できますが、現行の Terraform GitHub provider では repository scope の設定項目として未対応のため、この repo では baseline 管理していません

## 使い方

1. `terraform/work-repositories/locals.tf` の `repositories` を必要に応じて調整します
2. `GITHUB_TOKEN` を設定します
3. 1Password desktop app 認証を使う場合は `terraform/work-repositories/terraform.tfvars` に `onepassword_account` を設定し、アプリ側で **Integrate with other apps** を有効にします
4. 既存 repo や既存 `main-default` ruleset を state に載せる必要がある場合は Terraform CLI の `import` を使います
5. `mise run plan` で差分を確認し、問題なければ `mise run apply` を実行します

```bash
export GITHUB_TOKEN="<your-token>"
cp terraform/work-repositories/terraform.tfvars.example terraform/work-repositories/terraform.tfvars
terraform -chdir=terraform/work-repositories init \
  -backend-config="bucket=<gcs-bucket>" \
  -backend-config="prefix=work-repositories"
mise run plan
mise run apply
```

backend は `terraform/work-repositories/backend.tf` で `backend "gcs" {}` の空ブロックだけを宣言し、bucket / prefix などの値は `terraform init -backend-config=...` で外から注入します。ローカルに `terraform/work-repositories/config.gcs.tfbackend` を作成して `terraform -chdir=terraform/work-repositories init -backend-config=config.gcs.tfbackend` のように渡すこともできます。

`terraform/work-repositories/terraform.tfvars` には次のように設定します。

```hcl
onepassword_account = "my.1password.com"
```

`.mise/config.toml` の task は内部で `terraform -chdir=terraform/work-repositories plan` / `terraform -chdir=terraform/work-repositories apply` を実行します。

## fork リポジトリ管理

`terraform/fork-repositories` は fork 作成用の root module です。`terraform/work-repositories` とは別の実行ディレクトリとして扱うため、state も分離されます。

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform/fork-repositories init \
  -backend-config="bucket=<gcs-bucket>" \
  -backend-config="prefix=fork-repositories"
mise run fork:plan
mise run fork:apply
```

fork リポジトリ側も `terraform/fork-repositories/backend.tf` で `backend "gcs" {}` の空ブロックだけを宣言し、GCS backend の値は init 時の `-backend-config` で注入します。ローカルに `terraform/fork-repositories/config.gcs.tfbackend` を作成する場合も、git 管理外として扱います。

`terraform/fork-repositories/locals.tf` には次のように fork 元を配列で指定します。fork 後の repo 名は、デフォルトでは `${source_repo}-fork` にします。repo 名が重複する場合や GitHub の repository name limit（100 文字）を超える場合は、`fork_name` で短く一意な repo 名を明示します。

```hcl
locals {
  repositories = [
    {
      source_owner = "zed-industries"
      source_repo  = "zed"
    },
    {
      source_owner = "another-owner"
      source_repo  = "zed"
      fork_name    = "another-zed-fork"
    },
  ]
}
```

`terraform/fork-repositories/main.tf` では personal account (`y-writings`) 前提で provider の owner を固定しています。`github_repository.fork` には `prevent_destroy = true` を設定し、誤って fork repo を destroy しないようにしています。

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
- GitHub Actions の Artifact and log retention も GitHub API では設定可能ですが、現行 Terraform GitHub provider では repository scope の resource / attribute がないため manual drift として扱います
- 現在の module-owned governance baseline は public repository 前提です。personal account では managed repositories を public 前提で扱ってください
