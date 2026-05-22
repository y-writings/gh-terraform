locals {
  github_app_token_presets = {
    release_please = {
      item_title              = "release-please-bot"
      app_id_section          = "info"
      app_id_field            = "app_id"
      private_key_secret_name = "RELEASE_PLEASE_APP_PRIVATE_KEY"
      app_id_variable_name    = "RELEASE_PLEASE_APP_ID"
    }

    changelog_approver = {
      item_title              = "changelog-approver-bot"
      app_id_section          = "INFO"
      app_id_field            = "app_id"
      private_key_secret_name = "CHANGELOG_APPROVER_APP_PRIVATE_KEY"
      app_id_variable_name    = "CHANGELOG_APPROVER_APP_ID"
    }
  }

  pat_token_presets = {
    metrics = {
      item_title  = "metrics-token"
      section     = "info"
      field       = "token"
      secret_name = "METRICS_TOKEN"
    }
  }

  repositories = {
    repo_93e3a3b5 = {
      name = "dotfiles"
    }
    repo_6a83d2cc = {
      name = "templates"
      github_app_tokens = {
        changelog_approver = merge(local.github_app_token_presets.changelog_approver, {
          vault_name = "dev"
        })
      }
    }
    repo_247d31ce = {
      name = "container"
    }
    repo_7e3c6bbd = {
      name = "karabiner-config"
    }
    repo_5e6c65a5 = {
      name = "gh-terraform"
      github_app_tokens = {
        changelog_approver = merge(local.github_app_token_presets.changelog_approver, {
          vault_name = "dev"
        })
      }
    }
    repo_fe83b6f2 = {
      name = "y-writings"
      pat_tokens = {
        metrics = merge(local.pat_token_presets.metrics, {
          vault_name = "dev"
        })
      }
    }
    repo_6e7bb53d = {
      name = "calver-beacon-action"
      github_app_tokens = {
        release_please = merge(local.github_app_token_presets.release_please, {
          vault_name = "dev"
        })
        changelog_approver = merge(local.github_app_token_presets.changelog_approver, {
          vault_name = "dev"
        })
      }
    }
    repo_cf0c042d = {
      name = "oc-logger"
    }
    repo_b004ad62 = {
      name = "opencode-keyflow"
    }
  }
}
