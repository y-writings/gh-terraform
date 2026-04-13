resource "github_repository_ruleset" "main_default" {
  name        = "main-default"
  repository  = var.repository_name
  target      = "branch"
  enforcement = var.ruleset_enforcement

  bypass_actors {
    actor_id    = var.bypass_repository_role_actor_id
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
      required_approving_review_count   = var.required_approving_review_count
      allowed_merge_methods             = ["squash"]
      dismiss_stale_reviews_on_push     = var.dismiss_stale_reviews_on_push
      require_code_owner_review         = var.require_code_owner_review
      require_last_push_approval        = var.require_last_push_approval
      required_review_thread_resolution = var.required_review_thread_resolution
    }

    dynamic "required_code_scanning" {
      for_each = var.required_code_scanning == null ? [] : [var.required_code_scanning]
      content {
        required_code_scanning_tool {
          tool                      = required_code_scanning.value.tool
          alerts_threshold          = required_code_scanning.value.alerts_threshold
          security_alerts_threshold = required_code_scanning.value.security_alerts_threshold
        }
      }
    }
  }
}
