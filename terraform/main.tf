provider "github" {
  owner = var.github_owner
}

resource "github_repository_ruleset" "main_default" {
  name        = "main-default"
  repository  = var.repository_name
  target      = "branch"
  enforcement = "active"

  bypass_actors {
    actor_id    = 5
    actor_type  = "RepositoryRole"
    bypass_mode = "pull_request"
  }

  conditions {
    ref_name {
      include = ["~DEFAULT_BRANCH"]
      exclude = []
    }
  }

  rules {
    creation                = true
    deletion                = true
    update                  = true
    required_linear_history = true
    non_fast_forward        = true

    pull_request {
      required_approving_review_count = 0
    }
  }
}
