variable "has_wiki" {
  description = "Whether the repository wiki is enabled"
  type        = bool
  default     = true
}

variable "delete_branch_on_merge" {
  description = "Delete the branch on merge setting for the repository"
  type        = bool
  default     = false
}
