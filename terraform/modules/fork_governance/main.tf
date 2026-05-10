locals {
  actions_enabled = false
}

resource "github_actions_repository_permissions" "this" {
  repository = var.repository_name

  enabled = local.actions_enabled
}
