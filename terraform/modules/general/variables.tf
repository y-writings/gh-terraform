variable "has_wiki" {
  description = "Whether the repository wiki is enabled"
  type        = bool
  default     = true
}

variable "has_issues" {
  description = "Whether issues are enabled for the repository"
  type        = bool
  default     = true
}

variable "allow_merge_commit" {
  description = "Whether merge commits are allowed"
  type        = bool
  default     = false
}

variable "allow_squash_merge" {
  description = "Whether squash merges are allowed"
  type        = bool
  default     = true
}

variable "squash_merge_commit_title" {
  description = "Squash merge commit title setting"
  type        = string
  default     = "PR_TITLE"
}

variable "squash_merge_commit_message" {
  description = "Squash merge commit message setting"
  type        = string
  default     = "PR_BODY"
}

variable "allow_rebase_merge" {
  description = "Whether rebase merges are allowed"
  type        = bool
  default     = false
}

variable "delete_branch_on_merge" {
  description = "Delete the branch on merge setting for the repository"
  type        = bool
  default     = false
}
