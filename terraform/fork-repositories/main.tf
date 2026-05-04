provider "github" {
  owner = "y-writings"
}

resource "github_repository" "fork" {
  for_each = {
    for repository in local.repositories : "${repository.source_repo}-fork" => repository
  }

  name         = each.key
  description  = "A fork of ${each.value.source_owner}/${each.value.source_repo}"
  fork         = true
  source_owner = each.value.source_owner
  source_repo  = each.value.source_repo

  lifecycle {
    prevent_destroy = true
  }
}
