locals {
  manage_security_and_analysis           = true
  vulnerability_alerts                   = true
  secret_scanning_status                 = "enabled"
  secret_scanning_push_protection_status = "enabled"

  visibility                  = "public"
  has_wiki                    = true
  has_issues                  = true
  allow_merge_commit          = false
  allow_squash_merge          = true
  squash_merge_commit_title   = "PR_TITLE"
  squash_merge_commit_message = "PR_BODY"
  allow_rebase_merge          = false
  delete_branch_on_merge      = true
}

resource "github_repository" "this" {
  name       = var.name
  visibility = local.visibility

  lifecycle {
    prevent_destroy = true
  }

  vulnerability_alerts = local.manage_security_and_analysis ? local.vulnerability_alerts : null

  dynamic "security_and_analysis" {
    for_each = local.manage_security_and_analysis ? [1] : []
    content {
      secret_scanning {
        status = local.secret_scanning_status
      }

      secret_scanning_push_protection {
        status = local.secret_scanning_push_protection_status
      }
    }
  }

  has_issues = local.has_issues
  has_wiki   = local.has_wiki

  delete_branch_on_merge = local.delete_branch_on_merge

  allow_merge_commit          = local.allow_merge_commit
  allow_squash_merge          = local.allow_squash_merge
  squash_merge_commit_title   = local.squash_merge_commit_title
  squash_merge_commit_message = local.squash_merge_commit_message
  allow_rebase_merge          = local.allow_rebase_merge
}
