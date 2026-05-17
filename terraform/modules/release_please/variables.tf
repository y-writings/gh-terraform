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
