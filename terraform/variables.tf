variable "github_owner" {
  description = "GitHub personal account name"
  type        = string
  default     = "y-writings"
}

variable "repositories" {
  description = "Repositories to manage under the personal account"
  type = map(object({
    import_existing_repository = optional(bool)
    main_default_ruleset_id    = optional(string)
  }))
}
