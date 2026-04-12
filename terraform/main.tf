provider "github" {
  owner = var.github_owner
}

locals {
  effective_repository_visibility = var.repository_visibility != null ? var.repository_visibility : (var.import_existing_repository ? null : "private")
  configure_advanced_security     = local.effective_repository_visibility != null && local.effective_repository_visibility != "public"
}

import {
  for_each = var.import_existing_repository ? toset([var.repository_name]) : toset([])
  to       = github_repository.this
  id       = each.value
}

resource "github_repository" "this" {
  name       = var.repository_name
  visibility = local.effective_repository_visibility

  lifecycle {
    prevent_destroy = true
  }

  # Dependabot alerts
  vulnerability_alerts = false

  security_and_analysis {
    # GitHub Advanced Security
    dynamic "advanced_security" {
      for_each = local.configure_advanced_security ? [1] : []
      content {
        status = "enabled"
      }
    }

    # Secret Protection
    secret_scanning {
      status = "disabled"
    }

    # Push protection
    secret_scanning_push_protection {
      status = "disabled"
    }
  }

  has_issues = true
  has_wiki   = true

  allow_merge_commit          = false
  allow_squash_merge          = true
  squash_merge_commit_title   = "PR_TITLE"
  squash_merge_commit_message = "PR_BODY"
  allow_rebase_merge          = false
}

resource "github_repository_ruleset" "main_default" {
  name        = "main-default"
  repository  = github_repository.this.name
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
