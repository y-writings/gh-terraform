provider "github" {
  owner = "y-writings"
}

locals {
  fork_repository_names = [
    for repository in local.repositories : coalesce(try(repository.fork_name, null), "${repository.source_repo}-fork")
  ]

  canonical_fork_repository_names = [
    for name in local.fork_repository_names : lower(name)
  ]

  duplicate_fork_repository_names = distinct([
    for name in local.fork_repository_names : name
    if length([for candidate in local.canonical_fork_repository_names : candidate if candidate == lower(name)]) > 1
  ])

  fork_repositories_by_name = {
    for name in distinct(local.fork_repository_names) : name => local.repositories[index(local.fork_repository_names, name)]
  }
}

resource "github_repository" "fork" {
  for_each = local.fork_repositories_by_name

  name         = each.key
  description  = "A fork of ${each.value.source_owner}/${each.value.source_repo}"
  fork         = true
  source_owner = each.value.source_owner
  source_repo  = each.value.source_repo

  lifecycle {
    prevent_destroy = true

    precondition {
      condition     = length(each.key) <= 100
      error_message = "Fork repository name '${each.key}' exceeds GitHub's 100-character repository name limit. Set fork_name in terraform/fork-repositories/locals.tf to a shorter value."
    }

    precondition {
      condition     = length(local.duplicate_fork_repository_names) == 0
      error_message = "Fork repository names must be unique. Duplicate name(s): ${join(", ", local.duplicate_fork_repository_names)}. Set fork_name in terraform/fork-repositories/locals.tf to give one of the forks a unique name."
    }
  }
}
