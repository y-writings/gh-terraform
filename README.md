# gh-terraform

GitHub リポジトリの初期設定を Terraform で管理するためのテンプレートです。

## 何を設定するか

- GitHub Provider (`integrations/github`) の有効化
- `main`（デフォルトブランチ）向けの `repository ruleset` (`main-default`) 作成
  - Restrict creations / updates / deletions
  - Require linear history
  - Require pull request before merging
  - Block force pushes
  - Repository admin ロールを Pull Request のみバイパス可能

## 使い方

1. 必要に応じて `terraform.tfvars.example` を `terraform.tfvars` にコピーして値を調整します。
2. GitHub トークンを環境変数として設定します。

```bash
export GITHUB_TOKEN="<your-token>"
terraform init
terraform plan
```

デフォルト設定では `y-writings/gh-terraform` を対象にします。
