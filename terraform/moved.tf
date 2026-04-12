moved {
  from = github_repository.this
  to   = module.repository["gh-terraform"].github_repository.this
}

moved {
  from = github_repository_ruleset.main_default
  to   = module.repository["gh-terraform"].github_repository_ruleset.main_default
}
