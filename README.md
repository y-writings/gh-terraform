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
  - GitHub Advanced Security: visibility が `private` / `internal` のときのみ `enabled`（`public` や visibility 未指定時は provider 制約回避のため未設定）
  - Dependabot alerts (`vulnerability_alerts`): `false`
  - Secret scanning: `disabled`
  - Push protection: `disabled`
- `main`（デフォルトブランチ）向けの `repository ruleset` (`main-default`) 作成
  - Restrict creations / updates / deletions
  - Require linear history
  - Require pull request before merging
  - Block force pushes
  - Repository admin ロールを Pull Request のみバイパス可能

## 添付画面との対応状況

Terraform provider (`integrations/github`) でこのテンプレートから直接設定している項目:

- Dependabot alerts
- Advanced Security
- Secret Protection / Push protection

Terraform provider 側で現時点このテンプレートでは未対応（対応 API/リソースが見当たらない、または専用設定が未提供）として扱っている項目:

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

既存リポジトリを Terraform 管理に取り込む場合は、`import_existing_repository = true` を明示的に設定してください。Terraform の `import` ブロックにより、初回 plan/apply 時に対象リポジトリ名を使って state に取り込みます。
