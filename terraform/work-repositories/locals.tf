locals {
  github_app_token_presets = {
    pr_creator = {
      vault_name              = "dev"
      item_title              = "pr-creator-bot"
      app_id_section          = "info"
      app_id_field            = "app_id"
      private_key_secret_name = "PR_CREATOR_APP_PRIVATE_KEY"
      app_id_variable_name    = "PR_CREATOR_APP_ID"
    }

    pr_approver = {
      vault_name              = "dev"
      item_title              = "pr-approver-bot"
      app_id_section          = "INFO"
      app_id_field            = "app_id"
      private_key_secret_name = "PR_APPROVER_APP_PRIVATE_KEY"
      app_id_variable_name    = "PR_APPROVER_APP_ID"
    }
  }

  pat_token_presets = {
    metrics = {
      vault_name  = "dev"
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
        pr_approver = local.github_app_token_presets.pr_approver
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
        pr_approver = local.github_app_token_presets.pr_approver
      }
    }
    repo_fe83b6f2 = {
      name = "y-writings"
    }
    repo_6e7bb53d = {
      name          = "calver-beacon-action"
      enable_codeql = true
      github_app_tokens = {
        pr_creator  = local.github_app_token_presets.pr_creator
        pr_approver = local.github_app_token_presets.pr_approver
      }
    }
    repo_cf0c042d = {
      name = "oc-logger"
    }
    repo_b004ad62 = {
      name = "opencode-keyflow"
    }
    repo_d1f71e9e = {
      name          = "driftline"
      enable_codeql = true
      github_app_tokens = {
        pr_creator  = local.github_app_token_presets.pr_creator
        pr_approver = local.github_app_token_presets.pr_approver
      }
    }
    repo_5fc995a0 = {
      name          = "pr-seal-action"
      enable_codeql = true
      github_app_tokens = {
        pr_creator  = local.github_app_token_presets.pr_creator
        pr_approver = local.github_app_token_presets.pr_approver
      }
    }
    repo_ea2e4c1b = {
      name = "gh-usecase"
      github_app_tokens = {
        pr_creator  = local.github_app_token_presets.pr_creator
        pr_approver = local.github_app_token_presets.pr_approver
      }
    }
  }
}
