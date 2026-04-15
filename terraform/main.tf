provider "github" {
  owner = "y-writings"
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
