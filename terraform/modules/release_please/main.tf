data "onepassword_vault" "github_app_tokens" {
  for_each = var.github_app_tokens

  name = each.value.vault_name
}

data "onepassword_vault" "pat_tokens" {
  for_each = var.pat_tokens

  name = each.value.vault_name
}

data "onepassword_item" "github_app_tokens" {
  for_each = var.github_app_tokens

  vault = data.onepassword_vault.github_app_tokens[each.key].uuid
  title = each.value.item_title
}

data "onepassword_item" "pat_tokens" {
  for_each = var.pat_tokens

  vault = data.onepassword_vault.pat_tokens[each.key].uuid
  title = each.value.item_title
}

resource "github_actions_secret" "github_app_private_key" {
  for_each = var.github_app_tokens

  repository      = var.repository_name
  secret_name     = each.value.private_key_secret_name
  plaintext_value = data.onepassword_item.github_app_tokens[each.key].private_key
}

resource "github_actions_variable" "github_app_id" {
  for_each = var.github_app_tokens

  repository    = var.repository_name
  variable_name = each.value.app_id_variable_name
  value         = data.onepassword_item.github_app_tokens[each.key].section_map[each.value.app_id_section].field_map[each.value.app_id_field].value
}

resource "github_actions_secret" "pat_token" {
  for_each = var.pat_tokens

  repository      = var.repository_name
  secret_name     = each.value.secret_name
  plaintext_value = data.onepassword_item.pat_tokens[each.key].section_map[each.value.section].field_map[each.value.field].value
}
