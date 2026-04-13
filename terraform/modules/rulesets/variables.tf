variable "repository_name" {
  description = "Target repository name"
  type        = string
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
  description = "Required approving review count for the shared main-default ruleset"
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

variable "dismiss_stale_reviews_on_push" {
  description = "Dismiss stale reviews when new commits are pushed in the shared ruleset"
  type        = bool
  default     = true
}

variable "require_code_owner_review" {
  description = "Require code owner review for the shared main-default ruleset"
  type        = bool
  default     = false
}

variable "require_last_push_approval" {
  description = "Require last push approval for the shared main-default ruleset"
  type        = bool
  default     = false
}

variable "required_review_thread_resolution" {
  description = "Require all review threads to be resolved before merge in the shared ruleset"
  type        = bool
  default     = true
}

variable "required_code_scanning" {
  description = "Optional code scanning requirement for the shared main-default ruleset"
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

variable "bypass_repository_role_actor_id" {
  description = "RepositoryRole actor_id used for the bypass actor. The current baseline uses 5."
  type        = number
  default     = 5
}
