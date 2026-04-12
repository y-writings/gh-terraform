variable "github_owner" {
  description = "GitHub personal account name"
  type        = string
  default     = "y-writings"
}

variable "repositories" {
  description = "Repositories to manage under the personal account"
  type = map(object({
    visibility                             = optional(string)
    import_existing_repository             = optional(bool)
    manage_security_and_analysis           = optional(bool)
    vulnerability_alerts                   = optional(bool)
    secret_scanning_status                 = optional(string)
    secret_scanning_push_protection_status = optional(string)
    delete_branch_on_merge                 = optional(bool)
    has_wiki                               = optional(bool)
    ruleset_enforcement                    = optional(string)
    required_approving_review_count        = optional(number)
    allowed_merge_methods                  = optional(list(string))
    dismiss_stale_reviews_on_push          = optional(bool)
    require_code_owner_review              = optional(bool)
    require_last_push_approval             = optional(bool)
    required_review_thread_resolution      = optional(bool)
    required_code_scanning = optional(object({
      tool                      = string
      alerts_threshold          = string
      security_alerts_threshold = string
    }))
    main_default_ruleset_id = optional(string)
  }))

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.visibility == null || contains(["public", "private"], repository.visibility)
    ])
    error_message = "repository visibility must be null, public, or private for personal accounts."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.secret_scanning_status == null || contains(["enabled", "disabled"], repository.secret_scanning_status)
    ])
    error_message = "secret_scanning_status must be enabled or disabled."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.secret_scanning_push_protection_status == null || contains(["enabled", "disabled"], repository.secret_scanning_push_protection_status)
    ])
    error_message = "secret_scanning_push_protection_status must be enabled or disabled."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.ruleset_enforcement == null || contains(["active", "disabled"], repository.ruleset_enforcement)
    ])
    error_message = "ruleset_enforcement must be active or disabled for personal-account repository rulesets."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.required_approving_review_count == null || repository.required_approving_review_count >= 0
    ])
    error_message = "required_approving_review_count must be zero or greater."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.required_approving_review_count == null || floor(repository.required_approving_review_count) == repository.required_approving_review_count
    ])
    error_message = "required_approving_review_count must be an integer."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.allowed_merge_methods == null || alltrue([
        for method in repository.allowed_merge_methods : contains(["merge", "squash", "rebase"], method)
      ])
    ])
    error_message = "allowed_merge_methods must contain only merge, squash, or rebase."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.required_code_scanning == null || contains(["CodeQL"], repository.required_code_scanning.tool)
    ])
    error_message = "required_code_scanning.tool must currently be CodeQL."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.required_code_scanning == null || contains(["none", "errors", "errors_and_warnings", "all"], repository.required_code_scanning.alerts_threshold)
    ])
    error_message = "required_code_scanning.alerts_threshold must be none, errors, errors_and_warnings, or all."
  }

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.required_code_scanning == null || contains(["none", "critical", "high_or_higher", "medium_or_higher", "all"], repository.required_code_scanning.security_alerts_threshold)
    ])
    error_message = "required_code_scanning.security_alerts_threshold must be none, critical, high_or_higher, medium_or_higher, or all."
  }
}
