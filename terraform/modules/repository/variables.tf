variable "name" {
  description = "Repository name"
  type        = string

  validation {
    condition     = trimspace(var.name) != ""
    error_message = "name must be a non-empty string."
  }
}
