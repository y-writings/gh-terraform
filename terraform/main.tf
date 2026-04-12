provider "github" {
  owner = var.github_owner
}

import {
  for_each = local.repositories_to_import
  to       = module.repository[each.key].github_repository.this
  id       = each.key
}

import {
  for_each = local.rulesets_to_import
  to       = module.repository[each.key].github_repository_ruleset.main_default
  id       = "${each.key}:${each.value.main_default_ruleset_id}"
}

module "repository" {
  source   = "./modules/repository"
  for_each = local.repositories

  repository_name                        = each.key
  repository_visibility                  = each.value.visibility
  manage_security_and_analysis           = each.value.manage_security_and_analysis
  vulnerability_alerts                   = each.value.vulnerability_alerts
  secret_scanning_status                 = each.value.secret_scanning_status
  secret_scanning_push_protection_status = each.value.secret_scanning_push_protection_status
  delete_branch_on_merge                 = each.value.delete_branch_on_merge
  has_wiki                               = each.value.has_wiki
  ruleset_enforcement                    = each.value.ruleset_enforcement
  required_approving_review_count        = each.value.required_approving_review_count
  allowed_merge_methods                  = each.value.allowed_merge_methods
  dismiss_stale_reviews_on_push          = each.value.dismiss_stale_reviews_on_push
  require_code_owner_review              = each.value.require_code_owner_review
  require_last_push_approval             = each.value.require_last_push_approval
  required_review_thread_resolution      = each.value.required_review_thread_resolution
  required_code_scanning                 = each.value.required_code_scanning
}
