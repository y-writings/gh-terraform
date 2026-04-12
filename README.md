# gh-terraform

GitHub リポジトリの初期設定を Terraform で管理するためのテンプレートです。

## ディレクトリ構成

```text
.
├── terraform/                  # Terraformコード一式
│   ├── main.tf                 # GitHub Provider設定と repository / repository ruleset の定義
│   ├── variables.tf            # github_owner / repository_name などの入力変数
│   ├── versions.tf             # Terraform本体と GitHub Provider のバージョン制約
│   └── terraform.tfvars.example
├── .github/                    # GitHub Actions などのリポジトリ運用設定
├── .commitlintrc.yaml          # commitlint 設定
├── lefthook.yaml               # lefthook 設定
└── release-please-config.json  # release-please 設定
```

## 何を設定するか

- GitHub Provider (`integrations/github`) の有効化
- `github_repository` によるセキュリティ関連設定
  - Repository visibility: `repository_visibility` を明示した場合だけ Terraform で管理し、未指定時は既存リポジトリの visibility と新規作成時の provider/GitHub デフォルトに任せる
  - `repository_visibility` を明示した場合のみ、Dependabot alerts (`vulnerability_alerts`) を `false`、Secret scanning / Push protection を `disabled` に設定
  - `repository_visibility` 未指定時は、既存リポジトリ import/new create ともにこれらの security 設定は Terraform から変更しない
- `main`（デフォルトブランチ）向けの `repository ruleset` (`main-default`) 作成
  - Restrict creations / updates / deletions
  - Require linear history
  - Require pull request before merging
  - Block force pushes
  - Repository admin ロールを Pull Request のみバイパス可能

## 添付画面との対応状況

Terraform provider (`integrations/github`) でこのテンプレートから直接設定している項目:

- Dependabot alerts
- Secret Protection / Push protection

Terraform provider 側で現時点このテンプレートでは未対応（対応 API/リソースが見当たらない、または専用設定が未提供）として扱っている項目:

- GitHub Advanced Security の enablement 管理（public は GitHub 側で常時有効、private/internal は org のライセンス構成によって provider 制約を受けるため、このテンプレートでは未設定）
- Private vulnerability reporting
- Dependency graph
- Automatic dependency submission
- Dependabot malware alerts
- Grouped security updates
- Dependabot version updates
- CodeQL analysis のセットアップ状態
- Copilot Autofix
- Protection rules (check runs failure threshold)

## 使い方

1. 必要に応じて `terraform/terraform.tfvars.example` を `terraform/terraform.tfvars` にコピーして値を調整します。
2. GitHub トークンを環境変数として設定します。
3. 既存リポジトリを管理対象にする場合だけ、`terraform/terraform.tfvars` で `import_existing_repository = true` を設定します。
4. `terraform plan` で差分を確認し、問題なければ apply します。

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform plan
```

デフォルト設定では `import_existing_repository = false` のため、既存リポジトリを誤って import してセキュリティ設定を変更しないようにしています。

既存リポジトリを Terraform 管理に取り込む場合は、`import_existing_repository = true` を明示的に設定してください。Terraform の `import` ブロックにより、初回 plan/apply 時に対象リポジトリ名を使って state に取り込みます。このとき `repository_visibility` を未指定にすると、既存 visibility と既存 security 設定を維持します。

新規リポジトリを作成する場合に public リポジトリとして運用したいなら、`repository_visibility = "public"` を明示してください。`repository_visibility` を未指定にした場合の visibility は provider/GitHub のデフォルトに従い、このテンプレートは visibility や security 設定を明示的には変更しません。

このリポジトリの Terraform は基本的に public リポジトリ向けに利用する想定です。そのため `terraform.tfvars.example` には import-safe な最小構成だけを置き、public リポジトリを新規作成する場合だけ `repository_visibility = "public"` を追加する形にしています。`import_existing_repository` は import ブロックを有効化するための一時スイッチであり、visibility の desired state そのものは `repository_visibility` で表現します。
