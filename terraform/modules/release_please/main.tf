locals {
  vault_name     = "dev"
  item_title     = "release-please-bot"
  app_id_section = "info"
  app_id_field   = "app_id"

  metrics_token_item_title   = "metrics-token"
  metrics_token_section_name = "info"
  metrics_token_field_name   = "token"
  metrics_token_secret_name  = "METRICS_TOKEN"

  app_private_key_secret_name = "RELEASE_PLEASE_APP_PRIVATE_KEY"
  app_id_variable_name        = "RELEASE_PLEASE_APP_ID"

  item    = data.onepassword_item.this
  section = local.item.section_map[local.app_id_section]

  metrics_token_item    = data.onepassword_item.metrics_token
  metrics_token_section = local.metrics_token_item.section_map[local.metrics_token_section_name]

  app_private_key = local.item.private_key
  app_id          = local.section.field_map[local.app_id_field].value
  metrics_token   = local.metrics_token_section.field_map[local.metrics_token_field_name].value
}

data "onepassword_vault" "dev" {
  name = local.vault_name
}

data "onepassword_item" "this" {
  vault = data.onepassword_vault.dev.uuid
  title = local.item_title
}

data "onepassword_item" "metrics_token" {
  vault = data.onepassword_vault.dev.uuid
  title = local.metrics_token_item_title
}

resource "github_actions_secret" "app_private_key" {
  repository      = var.repository_name
  secret_name     = local.app_private_key_secret_name
  plaintext_value = local.app_private_key
}

resource "github_actions_variable" "app_id" {
  repository    = var.repository_name
  variable_name = local.app_id_variable_name
  value         = local.app_id
}

resource "github_actions_secret" "metrics_token" {
  count = var.repository_name == "y-writings" ? 1 : 0

  repository      = var.repository_name
  secret_name     = local.metrics_token_secret_name
  plaintext_value = local.metrics_token
}
