removed {
  from = github_repository.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = github_repository_ruleset.main_default

  lifecycle {
    destroy = false
  }
}
