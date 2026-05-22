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

  repository_name   = each.value.name
  github_app_tokens = try(each.value.github_app_tokens, {})
  pat_tokens        = try(each.value.pat_tokens, {})
}

module "governance_rulesets" {
  source   = "../modules/governance"
  for_each = local.repositories

  depends_on = [module.repository]

  repository_name = each.value.name
}
