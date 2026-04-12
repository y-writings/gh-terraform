variable "github_owner" {
  description = "GitHub organization or user name"
  type        = string
  default     = "y-writings"
}

variable "repository_name" {
  description = "Target repository name"
  type        = string
  default     = "gh-terraform"
}

variable "import_existing_repository" {
  description = "Import an existing repository into state before managing settings"
  type        = bool
  default     = true
}
