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

variable "repository_visibility" {
  description = "Repository visibility (public, private, or internal). Set to null to preserve imported visibility; new repositories default to private"
  type        = string
  default     = null
}

variable "import_existing_repository" {
  description = "Import an existing repository into state before managing settings. Enable this explicitly for existing repositories"
  type        = bool
  default     = false
}
