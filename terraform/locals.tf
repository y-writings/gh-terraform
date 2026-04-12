locals {
  repository_defaults = {
    visibility                             = null
    import_existing_repository             = false
    manage_security_and_analysis           = false
    vulnerability_alerts                   = true
    secret_scanning_status                 = "enabled"
    secret_scanning_push_protection_status = "enabled"
    delete_branch_on_merge                 = false
    has_wiki                               = true
    ruleset_enforcement                    = "active"
    required_approving_review_count        = 1
    allowed_merge_methods                  = ["squash"]
    dismiss_stale_reviews_on_push          = true
    require_code_owner_review              = false
    require_last_push_approval             = false
    required_review_thread_resolution      = true
    required_code_scanning                 = null
    main_default_ruleset_id                = null
  }

  repositories = {
    for name, config in var.repositories :
    name => {
      visibility                             = config.visibility != null ? config.visibility : local.repository_defaults.visibility
      import_existing_repository             = coalesce(config.import_existing_repository, local.repository_defaults.import_existing_repository)
      manage_security_and_analysis           = coalesce(config.manage_security_and_analysis, local.repository_defaults.manage_security_and_analysis)
      vulnerability_alerts                   = coalesce(config.vulnerability_alerts, local.repository_defaults.vulnerability_alerts)
      secret_scanning_status                 = coalesce(config.secret_scanning_status, local.repository_defaults.secret_scanning_status)
      secret_scanning_push_protection_status = coalesce(config.secret_scanning_push_protection_status, local.repository_defaults.secret_scanning_push_protection_status)
      delete_branch_on_merge                 = coalesce(config.delete_branch_on_merge, local.repository_defaults.delete_branch_on_merge)
      has_wiki                               = coalesce(config.has_wiki, local.repository_defaults.has_wiki)
      ruleset_enforcement                    = coalesce(config.ruleset_enforcement, local.repository_defaults.ruleset_enforcement)
      required_approving_review_count        = coalesce(config.required_approving_review_count, local.repository_defaults.required_approving_review_count)
      allowed_merge_methods                  = coalesce(config.allowed_merge_methods, local.repository_defaults.allowed_merge_methods)
      dismiss_stale_reviews_on_push          = coalesce(config.dismiss_stale_reviews_on_push, local.repository_defaults.dismiss_stale_reviews_on_push)
      require_code_owner_review              = coalesce(config.require_code_owner_review, local.repository_defaults.require_code_owner_review)
      require_last_push_approval             = coalesce(config.require_last_push_approval, local.repository_defaults.require_last_push_approval)
      required_review_thread_resolution      = coalesce(config.required_review_thread_resolution, local.repository_defaults.required_review_thread_resolution)
      required_code_scanning                 = config.required_code_scanning != null ? config.required_code_scanning : local.repository_defaults.required_code_scanning
      main_default_ruleset_id                = config.main_default_ruleset_id != null ? config.main_default_ruleset_id : local.repository_defaults.main_default_ruleset_id
    }
  }

  repositories_to_import = {
    for name, config in local.repositories :
    name => config
    if config.import_existing_repository
  }

  rulesets_to_import = {
    for name, config in local.repositories :
    name => config
    if config.main_default_ruleset_id != null
  }
}
