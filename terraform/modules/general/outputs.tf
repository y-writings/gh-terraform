output "has_wiki" {
  description = "Whether the repository wiki is enabled"
  value       = var.has_wiki
}

output "has_issues" {
  description = "Whether issues are enabled for the repository"
  value       = var.has_issues
}

output "allow_merge_commit" {
  description = "Whether merge commits are allowed"
  value       = var.allow_merge_commit
}

output "allow_squash_merge" {
  description = "Whether squash merges are allowed"
  value       = var.allow_squash_merge
}

output "squash_merge_commit_title" {
  description = "Squash merge commit title setting"
  value       = var.squash_merge_commit_title
}

output "squash_merge_commit_message" {
  description = "Squash merge commit message setting"
  value       = var.squash_merge_commit_message
}

output "allow_rebase_merge" {
  description = "Whether rebase merges are allowed"
  value       = var.allow_rebase_merge
}

output "delete_branch_on_merge" {
  description = "Delete the branch on merge setting for the repository"
  value       = var.delete_branch_on_merge
}
