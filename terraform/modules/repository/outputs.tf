output "repository_name" {
  description = "Managed repository name"
  value       = github_repository.this.name
}

output "main_default_ruleset_id" {
  description = "ID of the main-default repository ruleset"
  value       = github_repository_ruleset.main_default.id
}
