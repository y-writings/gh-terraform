provider "github" {
  owner = "y-writings"
}

provider "onepassword" {
  account = var.onepassword_account
}

module "repository" {
  source   = "../modules/repository"
  for_each = local.repositories

  name = each.value.name
}

module "release_please" {
  source   = "../modules/release_please"
  for_each = local.repositories

  depends_on = [module.repository]

  repository_name          = each.value.name
  metrics_token            = try(each.value.metrics_token, null)
  release_please_token     = try(each.value.release_please_token, null)
  changelog_approver_token = try(each.value.changelog_approver_token, null)
}

module "governance_rulesets" {
  source   = "../modules/governance"
  for_each = local.repositories

  depends_on = [module.repository]

  repository_name = each.value.name
}
