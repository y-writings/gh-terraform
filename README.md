# gh-terraform

GitHub リポジトリの初期設定を Terraform で管理するためのテンプレートです。

## ディレクトリ構成

```text
.
├── terraform/                  # Terraformコード一式
│   ├── main.tf                 # GitHub Provider設定と repository / repository ruleset の定義
│   ├── variables.tf            # github_owner / repository_name の入力変数
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
  - GitHub Advanced Security: `enabled`
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
3. 既存リポジトリを管理対象にする場合は import してください。

```bash
export GITHUB_TOKEN="<your-token>"
terraform -chdir=terraform init
terraform -chdir=terraform import github_repository.main <owner>/<repository>
terraform -chdir=terraform plan
```

デフォルト設定では `y-writings/gh-terraform` を対象にします。
