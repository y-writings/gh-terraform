provider "github" {
  owner = var.github_owner
}

import {
  for_each = local.repositories_to_import
  to       = github_repository.this[each.key]
  id       = each.key
}

import {
  for_each = local.rulesets_to_import
  to       = module.rulesets[each.key].github_repository_ruleset.main_default
  id       = "${each.key}:${each.value.main_default_ruleset_id}"
}

module "general" {
  source = "./modules/general"
}

module "advanced_security" {
  source = "./modules/advanced_security"
}

resource "github_repository" "this" {
  for_each = local.repositories

  name       = each.key
  visibility = each.value.visibility

  lifecycle {
    prevent_destroy = true
  }

  vulnerability_alerts = module.advanced_security.manage_security_and_analysis ? module.advanced_security.vulnerability_alerts : null

  dynamic "security_and_analysis" {
    for_each = module.advanced_security.manage_security_and_analysis ? [1] : []
    content {
      secret_scanning {
        status = module.advanced_security.secret_scanning_status
      }

      secret_scanning_push_protection {
        status = module.advanced_security.secret_scanning_push_protection_status
      }
    }
  }

  has_issues = module.general.has_issues
  has_wiki   = module.general.has_wiki

  delete_branch_on_merge = module.general.delete_branch_on_merge

  allow_merge_commit          = module.general.allow_merge_commit
  allow_squash_merge          = module.general.allow_squash_merge
  squash_merge_commit_title   = module.general.squash_merge_commit_title
  squash_merge_commit_message = module.general.squash_merge_commit_message
  allow_rebase_merge          = module.general.allow_rebase_merge
}

module "rulesets" {
  source   = "./modules/rulesets"
  for_each = local.repositories

  repository_name                   = github_repository.this[each.key].name
  ruleset_enforcement               = local.repository_governance.ruleset_enforcement
  required_approving_review_count   = local.repository_governance.required_approving_review_count
  dismiss_stale_reviews_on_push     = local.repository_governance.dismiss_stale_reviews_on_push
  require_code_owner_review         = local.repository_governance.require_code_owner_review
  require_last_push_approval        = local.repository_governance.require_last_push_approval
  required_review_thread_resolution = local.repository_governance.required_review_thread_resolution
  required_code_scanning            = local.repository_governance.required_code_scanning
}
