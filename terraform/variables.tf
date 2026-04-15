variable "github_owner" {
  description = "GitHub personal account name"
  type        = string
  default     = "y-writings"
}

variable "repositories" {
  description = "Repositories to manage under the personal account"
  type        = map(object({}))
}
