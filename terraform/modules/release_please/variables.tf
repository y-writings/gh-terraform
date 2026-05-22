variable "repository_name" {
  description = "Repository name"
  type        = string

  validation {
    condition     = trimspace(var.repository_name) != ""
    error_message = "repository_name must be a non-empty string."
  }
}

variable "enable_metrics_token" {
  description = "Whether to create the METRICS_TOKEN Actions secret for this repository."
  type        = bool
  default     = false
}

variable "enable_release_please_token" {
  description = "Whether to create the release-please GitHub App Actions secret and variable for this repository."
  type        = bool
  default     = false
}
