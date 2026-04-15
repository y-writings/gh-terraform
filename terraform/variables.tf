variable "github_owner" {
  description = "GitHub personal account name"
  type        = string
  default     = "y-writings"
}

variable "repositories" {
  description = "Repositories to manage under the personal account"
  type = map(object({
    visibility                 = optional(string)
    import_existing_repository = optional(bool)
    main_default_ruleset_id    = optional(string)
  }))

  validation {
    condition = alltrue([
      for repository in values(var.repositories) : repository.visibility == null || repository.visibility == "public"
    ])
    error_message = "repository visibility must be omitted or set to public for managed personal-account repositories."
  }

}
