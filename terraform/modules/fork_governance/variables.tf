variable "repository_name" {
  description = "Target fork repository name"
  type        = string

  validation {
    condition     = trimspace(var.repository_name) != ""
    error_message = "repository_name must be a non-empty string."
  }
}
