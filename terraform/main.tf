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
  to       = module.governance_rulesets[each.key].github_repository_ruleset.main_default[0]
  id       = "${each.key}:${each.value.main_default_ruleset_id}"
}

module "repository" {
  source   = "./modules/repository"
  for_each = local.repositories

  name = each.key
}

module "governance_rulesets" {
  source   = "./modules/governance"
  for_each = local.repositories

  repository_name = module.repository[each.key].name
}
