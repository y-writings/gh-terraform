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
  delete_branch_on_merge                 = each.value.delete_branch_on_merge
  has_wiki                               = each.value.has_wiki
  manage_security_and_analysis           = local.repository_governance.manage_security_and_analysis
  vulnerability_alerts                   = local.repository_governance.vulnerability_alerts
  secret_scanning_status                 = local.repository_governance.secret_scanning_status
  secret_scanning_push_protection_status = local.repository_governance.secret_scanning_push_protection_status
  ruleset_enforcement                    = local.repository_governance.ruleset_enforcement
  required_approving_review_count        = local.repository_governance.required_approving_review_count
  dismiss_stale_reviews_on_push          = local.repository_governance.dismiss_stale_reviews_on_push
  require_code_owner_review              = local.repository_governance.require_code_owner_review
  require_last_push_approval             = local.repository_governance.require_last_push_approval
  required_review_thread_resolution      = local.repository_governance.required_review_thread_resolution
  required_code_scanning                 = local.repository_governance.required_code_scanning
}
