variable "repository_name" {
  description = "Repository name"
  type        = string

  validation {
    condition     = trimspace(var.repository_name) != ""
    error_message = "repository_name must be a non-empty string."
  }
}

variable "metrics_token" {
  description = "Complete METRICS_TOKEN configuration. Set to null to skip creating the secret."
  type = object({
    vault_name  = string
    item_title  = string
    section     = string
    field       = string
    secret_name = string
  })
}

variable "release_please_token" {
  description = "Complete release-please GitHub App configuration. Set to null to skip creating the secret and variable."
  type = object({
    vault_name              = string
    item_title              = string
    app_id_section          = string
    app_id_field            = string
    private_key_secret_name = string
    app_id_variable_name    = string
  })
}

variable "changelog_approver_token" {
  description = "Complete changelog approver GitHub App configuration. Set to null to skip creating the secret and variable."
  type = object({
    vault_name              = string
    item_title              = string
    app_id_section          = string
    app_id_field            = string
    private_key_secret_name = string
    app_id_variable_name    = string
  })
}
