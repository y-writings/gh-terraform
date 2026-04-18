locals {
  vault_name     = "dev"
  item_title     = "release-please-bot"
  app_id_section = "info"
  app_id_field   = "app_id"

  app_private_key_secret_name = "RELEASE_PLEASE_APP_PRIVATE_KEY"
  app_id_variable_name        = "RELEASE_PLEASE_APP_ID"

  item    = data.onepassword_item.this
  section = local.item.section_map[local.app_id_section]

  app_private_key = local.item.private_key
  app_id          = local.section.field_map[local.app_id_field].value
}

data "onepassword_vault" "dev" {
  name = local.vault_name
}

data "onepassword_item" "this" {
  vault = data.onepassword_vault.dev.uuid
  title = local.item_title
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
