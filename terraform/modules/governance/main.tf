resource "github_repository_ruleset" "main_default" {
  count = var.repository_name == null ? 0 : 1

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
      required_approving_review_count   = 1
      allowed_merge_methods             = ["squash"]
      dismiss_stale_reviews_on_push     = true
      require_code_owner_review         = false
      require_last_push_approval        = false
      required_review_thread_resolution = true
    }

    required_code_scanning {
      required_code_scanning_tool {
        tool                      = "CodeQL"
        alerts_threshold          = "errors_and_warnings"
        security_alerts_threshold = "high_or_higher"
      }
    }
  }
}
