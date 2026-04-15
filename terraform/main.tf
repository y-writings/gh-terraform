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
  to       = module.governance_rulesets[each.key].github_repository_ruleset.main_default[0]
  id       = "${each.key}:${each.value.main_default_ruleset_id}"
}

module "governance" {
  source = "./modules/governance"
}

module "governance_rulesets" {
  source   = "./modules/governance"
  for_each = local.repositories

  repository_name = github_repository.this[each.key].name
}

resource "github_repository" "this" {
  for_each = local.repositories

  name       = each.key
  visibility = each.value.visibility

  lifecycle {
    prevent_destroy = true
  }

  vulnerability_alerts = module.governance.manage_security_and_analysis ? module.governance.vulnerability_alerts : null

  dynamic "security_and_analysis" {
    for_each = module.governance.manage_security_and_analysis ? [1] : []
    content {
      secret_scanning {
        status = module.governance.secret_scanning_status
      }

      secret_scanning_push_protection {
        status = module.governance.secret_scanning_push_protection_status
      }
    }
  }

  has_issues = module.governance.has_issues
  has_wiki   = module.governance.has_wiki

  delete_branch_on_merge = module.governance.delete_branch_on_merge

  allow_merge_commit          = module.governance.allow_merge_commit
  allow_squash_merge          = module.governance.allow_squash_merge
  squash_merge_commit_title   = module.governance.squash_merge_commit_title
  squash_merge_commit_message = module.governance.squash_merge_commit_message
  allow_rebase_merge          = module.governance.allow_rebase_merge
}
