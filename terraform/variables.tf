variable "github_owner" {
  description = "GitHub personal account name"
  type        = string
  default     = "y-writings"
}

variable "repositories" {
  description = "Repositories to manage under the personal account"
  type = map(object({
    visibility                 = optional(string)
    import_existing_repository = optional(bool)
    delete_branch_on_merge     = optional(bool)
    has_wiki                   = optional(bool)
    main_default_ruleset_id    = optional(string)
  }))

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.visibility == null || contains(["public", "private"], repository.visibility)
    ])
    error_message = "repository visibility must be null, public, or private for personal accounts."
  }

}

variable "repository_governance" {
  description = "Shared governance controls applied uniformly to all managed repositories"
  type = object({
    manage_security_and_analysis           = optional(bool)
    enable_required_code_scanning          = optional(bool)
    vulnerability_alerts                   = optional(bool)
    secret_scanning_status                 = optional(string)
    secret_scanning_push_protection_status = optional(string)
    ruleset_enforcement                    = optional(string)
    required_approving_review_count        = optional(number)
    dismiss_stale_reviews_on_push          = optional(bool)
    require_code_owner_review              = optional(bool)
    require_last_push_approval             = optional(bool)
    required_review_thread_resolution      = optional(bool)
    required_code_scanning = optional(object({
      tool                      = string
      alerts_threshold          = string
      security_alerts_threshold = string
    }))
  })
  default = {}

  validation {
    condition     = var.repository_governance.secret_scanning_status == null || contains(["enabled", "disabled"], var.repository_governance.secret_scanning_status)
    error_message = "repository_governance.secret_scanning_status must be enabled or disabled."
  }

  validation {
    condition     = var.repository_governance.secret_scanning_push_protection_status == null || contains(["enabled", "disabled"], var.repository_governance.secret_scanning_push_protection_status)
    error_message = "repository_governance.secret_scanning_push_protection_status must be enabled or disabled."
  }

  validation {
    condition     = var.repository_governance.ruleset_enforcement == null || contains(["active", "disabled"], var.repository_governance.ruleset_enforcement)
    error_message = "repository_governance.ruleset_enforcement must be active or disabled for personal-account repository rulesets."
  }

  validation {
    condition     = var.repository_governance.required_approving_review_count == null || var.repository_governance.required_approving_review_count >= 0
    error_message = "repository_governance.required_approving_review_count must be zero or greater."
  }

  validation {
    condition     = var.repository_governance.required_approving_review_count == null || floor(var.repository_governance.required_approving_review_count) == var.repository_governance.required_approving_review_count
    error_message = "repository_governance.required_approving_review_count must be an integer."
  }

  validation {
    condition     = coalesce(var.repository_governance.enable_required_code_scanning, true) || var.repository_governance.required_code_scanning == null
    error_message = "repository_governance.required_code_scanning must be null or omitted when repository_governance.enable_required_code_scanning is false."
  }

  validation {
    condition     = var.repository_governance.required_code_scanning == null || trimspace(var.repository_governance.required_code_scanning.tool) != ""
    error_message = "repository_governance.required_code_scanning.tool must be a non-empty string."
  }

  validation {
    condition     = var.repository_governance.required_code_scanning == null || contains(["none", "errors", "errors_and_warnings", "all"], var.repository_governance.required_code_scanning.alerts_threshold)
    error_message = "repository_governance.required_code_scanning.alerts_threshold must be none, errors, errors_and_warnings, or all."
  }

  validation {
    condition     = var.repository_governance.required_code_scanning == null || contains(["none", "critical", "high_or_higher", "medium_or_higher", "all"], var.repository_governance.required_code_scanning.security_alerts_threshold)
    error_message = "repository_governance.required_code_scanning.security_alerts_threshold must be none, critical, high_or_higher, medium_or_higher, or all."
  }
}
