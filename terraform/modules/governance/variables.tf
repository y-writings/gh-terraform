variable "repository_name" {
  description = "Target repository name"
  type        = string
  default     = null

  validation {
    condition     = var.repository_name == null || trimspace(var.repository_name) != ""
    error_message = "repository_name must be null or a non-empty string."
  }
}
