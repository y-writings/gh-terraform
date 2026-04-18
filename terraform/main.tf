provider "github" {
  owner = "y-writings"
}

provider "onepassword" {
  account = var.onepassword_account
}

module "repository" {
  source   = "./modules/repository"
  for_each = local.repositories

  name = each.key
}

module "release_please" {
  source   = "./modules/release_please"
  for_each = local.repositories

  depends_on = [module.repository]

  repository_name = each.key
}

module "governance_rulesets" {
  source   = "./modules/governance"
  for_each = local.repositories

  depends_on = [module.repository]

  repository_name = each.key
}
