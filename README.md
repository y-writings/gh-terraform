# gh-terraform

GitHub リポジトリの初期設定を Terraform で管理するためのテンプレートです。

## ディレクトリ構成

```text
.
├── terraform/                  # Terraformコード一式
│   ├── main.tf                 # GitHub Provider設定と repository ruleset の定義
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
- `main`（デフォルトブランチ）向けの `repository ruleset` (`main-default`) 作成
  - Restrict creations / updates / deletions
  - Require linear history
  - Require pull request before merging
  - Block force pushes
  - Repository admin ロールを Pull Request のみバイパス可能

## 使い方

1. 必要に応じて `terraform/terraform.tfvars.example` を `terraform/terraform.tfvars` にコピーして値を調整します。
2. GitHub トークンを環境変数として設定します。

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform plan
```

デフォルト設定では `y-writings/gh-terraform` を対象にし、既存リポジトリを初回 apply 時に import します（`import_existing_repository = true`）。
