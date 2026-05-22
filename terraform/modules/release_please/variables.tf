variable "repository_name" {
  description = "Repository name"
  type        = string

  validation {
    condition     = trimspace(var.repository_name) != ""
    error_message = "repository_name must be a non-empty string."
  }
}

variable "github_app_tokens" {
  description = "GitHub App token configurations keyed by stable token name."
  type = map(object({
    vault_name              = string
    item_title              = string
    app_id_section          = string
    app_id_field            = string
    private_key_secret_name = string
    app_id_variable_name    = string
  }))
  default = {}
}

variable "pat_tokens" {
  description = "PAT token configurations keyed by stable token name."
  type = map(object({
    vault_name  = string
    item_title  = string
    section     = string
    field       = string
    secret_name = string
  }))
  default = {}
}
