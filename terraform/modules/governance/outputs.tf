locals {
  manage_security_and_analysis           = true
  vulnerability_alerts                   = true
  secret_scanning_status                 = "enabled"
  secret_scanning_push_protection_status = "enabled"

  has_wiki                    = true
  has_issues                  = true
  allow_merge_commit          = false
  allow_squash_merge          = true
  squash_merge_commit_title   = "PR_TITLE"
  squash_merge_commit_message = "PR_BODY"
  allow_rebase_merge          = false
  delete_branch_on_merge      = false
}

output "manage_security_and_analysis" {
  description = "Whether Terraform should manage vulnerability alerts and security_and_analysis settings"
  value       = local.manage_security_and_analysis
}

output "vulnerability_alerts" {
  description = "Dependabot alerts setting applied when manage_security_and_analysis is true"
  value       = local.vulnerability_alerts
}

output "secret_scanning_status" {
  description = "Secret scanning status applied when manage_security_and_analysis is true"
  value       = local.secret_scanning_status
}

output "secret_scanning_push_protection_status" {
  description = "Push protection status applied when manage_security_and_analysis is true"
  value       = local.secret_scanning_push_protection_status
}

output "has_wiki" {
  description = "Whether the repository wiki is enabled"
  value       = local.has_wiki
}

output "has_issues" {
  description = "Whether issues are enabled for the repository"
  value       = local.has_issues
}

output "allow_merge_commit" {
  description = "Whether merge commits are allowed"
  value       = local.allow_merge_commit
}

output "allow_squash_merge" {
  description = "Whether squash merges are allowed"
  value       = local.allow_squash_merge
}

output "squash_merge_commit_title" {
  description = "Squash merge commit title setting"
  value       = local.squash_merge_commit_title
}

output "squash_merge_commit_message" {
  description = "Squash merge commit message setting"
  value       = local.squash_merge_commit_message
}

output "allow_rebase_merge" {
  description = "Whether rebase merges are allowed"
  value       = local.allow_rebase_merge
}

output "delete_branch_on_merge" {
  description = "Delete the branch on merge setting for the repository"
  value       = local.delete_branch_on_merge
}

output "main_default_ruleset_id" {
  description = "ID of the main-default repository ruleset"
  value       = var.repository_name == null ? null : github_repository_ruleset.main_default[0].id
}
