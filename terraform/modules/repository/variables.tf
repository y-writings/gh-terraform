variable "repository_name" {
  description = "Target repository name"
  type        = string
}

variable "repository_visibility" {
  description = "Repository visibility. Set to null to preserve imported visibility or provider defaults"
  type        = string
  default     = null

  validation {
    condition     = var.repository_visibility == null || contains(["public", "private"], var.repository_visibility)
    error_message = "repository_visibility must be null, public, or private for personal accounts."
  }
}

variable "manage_security_and_analysis" {
  description = "Whether Terraform should manage vulnerability alerts and security_and_analysis settings"
  type        = bool
  default     = false
}

variable "vulnerability_alerts" {
  description = "Dependabot alerts setting applied when manage_security_and_analysis is true"
  type        = bool
  default     = true
}

variable "secret_scanning_status" {
  description = "Secret scanning status applied when manage_security_and_analysis is true"
  type        = string
  default     = "enabled"

  validation {
    condition     = contains(["enabled", "disabled"], var.secret_scanning_status)
    error_message = "secret_scanning_status must be enabled or disabled."
  }
}

variable "secret_scanning_push_protection_status" {
  description = "Push protection status applied when manage_security_and_analysis is true"
  type        = string
  default     = "enabled"

  validation {
    condition     = contains(["enabled", "disabled"], var.secret_scanning_push_protection_status)
    error_message = "secret_scanning_push_protection_status must be enabled or disabled."
  }
}

variable "ruleset_enforcement" {
  description = "Repository ruleset enforcement mode"
  type        = string
  default     = "active"

  validation {
    condition     = contains(["active", "disabled"], var.ruleset_enforcement)
    error_message = "ruleset_enforcement must be active or disabled for personal accounts."
  }
}

variable "required_approving_review_count" {
  description = "Required approving review count for the main-default ruleset"
  type        = number
  default     = 1

  validation {
    condition     = var.required_approving_review_count >= 0
    error_message = "required_approving_review_count must be zero or greater."
  }

  validation {
    condition     = floor(var.required_approving_review_count) == var.required_approving_review_count
    error_message = "required_approving_review_count must be an integer."
  }
}

variable "allowed_merge_methods" {
  description = "Allowed merge methods for the main-default ruleset pull_request rule"
  type        = list(string)
  default     = ["squash"]

  validation {
    condition = alltrue([
      for method in var.allowed_merge_methods : contains(["merge", "squash", "rebase"], method)
    ])
    error_message = "allowed_merge_methods must contain only merge, squash, or rebase."
  }
}

variable "dismiss_stale_reviews_on_push" {
  description = "Dismiss stale reviews when new commits are pushed"
  type        = bool
  default     = true
}

variable "require_code_owner_review" {
  description = "Require code owner review for the main-default ruleset"
  type        = bool
  default     = false
}

variable "require_last_push_approval" {
  description = "Require last push approval for the main-default ruleset"
  type        = bool
  default     = false
}

variable "required_review_thread_resolution" {
  description = "Require all review threads to be resolved before merge"
  type        = bool
  default     = true
}

variable "required_code_scanning" {
  description = "Optional code scanning requirement for the main-default ruleset"
  type = object({
    tool                      = string
    alerts_threshold          = string
    security_alerts_threshold = string
  })
  default = null

  validation {
    condition     = var.required_code_scanning == null || trimspace(var.required_code_scanning.tool) != ""
    error_message = "required_code_scanning.tool must be a non-empty string."
  }

  validation {
    condition     = var.required_code_scanning == null || contains(["none", "errors", "errors_and_warnings", "all"], var.required_code_scanning.alerts_threshold)
    error_message = "required_code_scanning.alerts_threshold must be none, errors, errors_and_warnings, or all."
  }

  validation {
    condition     = var.required_code_scanning == null || contains(["none", "critical", "high_or_higher", "medium_or_higher", "all"], var.required_code_scanning.security_alerts_threshold)
    error_message = "required_code_scanning.security_alerts_threshold must be none, critical, high_or_higher, medium_or_higher, or all."
  }
}

variable "delete_branch_on_merge" {
  description = "Delete the branch on merge setting for the repository"
  type        = bool
  default     = false
}

variable "has_wiki" {
  description = "Whether the repository wiki is enabled"
  type        = bool
  default     = true
}

variable "bypass_repository_role_actor_id" {
  description = "RepositoryRole actor_id used for the bypass actor. The current baseline uses 5."
  type        = number
  default     = 5
}
