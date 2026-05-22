locals {
  changelog_approver_vault_name                  = "dev"
  changelog_approver_item_title                  = "changelog-approver-bot"
  changelog_approver_app_id_section              = "INFO"
  changelog_approver_app_id_field                = "app_id"
  changelog_approver_app_private_key_secret_name = "CHANGELOG_APPROVER_APP_PRIVATE_KEY"
  changelog_approver_app_id_variable_name        = "CHANGELOG_APPROVER_APP_ID"

  changelog_approver_item    = data.onepassword_item.changelog_approver
  changelog_approver_section = local.changelog_approver_item.section_map[local.changelog_approver_app_id_section]

  changelog_approver_app_private_key = local.changelog_approver_item.private_key
  changelog_approver_app_id          = local.changelog_approver_section.field_map[local.changelog_approver_app_id_field].value
}

data "onepassword_vault" "release_please_token" {
  count = var.release_please_token != null ? 1 : 0

  name = var.release_please_token.vault_name
}

data "onepassword_vault" "changelog_approver" {
  name = local.changelog_approver_vault_name
}

data "onepassword_vault" "metrics_token" {
  count = var.metrics_token != null ? 1 : 0

  name = var.metrics_token.vault_name
}

data "onepassword_item" "this" {
  count = var.release_please_token != null ? 1 : 0

  vault = data.onepassword_vault.release_please_token[0].uuid
  title = var.release_please_token.item_title
}

data "onepassword_item" "changelog_approver" {
  vault = data.onepassword_vault.changelog_approver.uuid
  title = local.changelog_approver_item_title
}

data "onepassword_item" "metrics_token" {
  count = var.metrics_token != null ? 1 : 0

  vault = data.onepassword_vault.metrics_token[0].uuid
  title = var.metrics_token.item_title
}

moved {
  from = github_actions_secret.app_private_key
  to   = github_actions_secret.app_private_key[0]
}

moved {
  from = github_actions_variable.app_id
  to   = github_actions_variable.app_id[0]
}

resource "github_actions_secret" "app_private_key" {
  count = var.release_please_token != null ? 1 : 0

  repository      = var.repository_name
  secret_name     = var.release_please_token.private_key_secret_name
  plaintext_value = data.onepassword_item.this[0].private_key
}

resource "github_actions_variable" "app_id" {
  count = var.release_please_token != null ? 1 : 0

  repository    = var.repository_name
  variable_name = var.release_please_token.app_id_variable_name
  value         = data.onepassword_item.this[0].section_map[var.release_please_token.app_id_section].field_map[var.release_please_token.app_id_field].value
}

resource "github_actions_secret" "changelog_approver_app_private_key" {
  repository      = var.repository_name
  secret_name     = local.changelog_approver_app_private_key_secret_name
  plaintext_value = local.changelog_approver_app_private_key
}

resource "github_actions_variable" "changelog_approver_app_id" {
  repository    = var.repository_name
  variable_name = local.changelog_approver_app_id_variable_name
  value         = local.changelog_approver_app_id
}

resource "github_actions_secret" "metrics_token" {
  count = var.metrics_token != null ? 1 : 0

  repository      = var.repository_name
  secret_name     = var.metrics_token.secret_name
  plaintext_value = data.onepassword_item.metrics_token[0].section_map[var.metrics_token.section].field_map[var.metrics_token.field].value
}
