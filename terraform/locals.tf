locals {
  repository_defaults = {
    visibility                 = null
    import_existing_repository = false
    delete_branch_on_merge     = false
    has_wiki                   = true
    main_default_ruleset_id    = null
  }

  repository_governance_defaults = {
    enable_required_code_scanning     = true
    ruleset_enforcement               = "active"
    required_approving_review_count   = 1
    dismiss_stale_reviews_on_push     = true
    require_code_owner_review         = false
    require_last_push_approval        = false
    required_review_thread_resolution = true
    required_code_scanning = {
      tool                      = "CodeQL"
      alerts_threshold          = "errors_and_warnings"
      security_alerts_threshold = "high_or_higher"
    }
  }

  repository_governance_enable_required_code_scanning = coalesce(
    var.repository_governance.enable_required_code_scanning,
    local.repository_governance_defaults.enable_required_code_scanning,
  )

  repository_governance = {
    enable_required_code_scanning     = local.repository_governance_enable_required_code_scanning
    ruleset_enforcement               = coalesce(var.repository_governance.ruleset_enforcement, local.repository_governance_defaults.ruleset_enforcement)
    required_approving_review_count   = coalesce(var.repository_governance.required_approving_review_count, local.repository_governance_defaults.required_approving_review_count)
    dismiss_stale_reviews_on_push     = coalesce(var.repository_governance.dismiss_stale_reviews_on_push, local.repository_governance_defaults.dismiss_stale_reviews_on_push)
    require_code_owner_review         = coalesce(var.repository_governance.require_code_owner_review, local.repository_governance_defaults.require_code_owner_review)
    require_last_push_approval        = coalesce(var.repository_governance.require_last_push_approval, local.repository_governance_defaults.require_last_push_approval)
    required_review_thread_resolution = coalesce(var.repository_governance.required_review_thread_resolution, local.repository_governance_defaults.required_review_thread_resolution)
    required_code_scanning            = local.repository_governance_enable_required_code_scanning ? (var.repository_governance.required_code_scanning != null ? var.repository_governance.required_code_scanning : local.repository_governance_defaults.required_code_scanning) : null
  }

  repositories = {
    for name, config in var.repositories :
    name => {
      visibility                 = config.visibility != null ? config.visibility : local.repository_defaults.visibility
      import_existing_repository = coalesce(config.import_existing_repository, local.repository_defaults.import_existing_repository)
      delete_branch_on_merge     = coalesce(config.delete_branch_on_merge, local.repository_defaults.delete_branch_on_merge)
      has_wiki                   = coalesce(config.has_wiki, local.repository_defaults.has_wiki)
      main_default_ruleset_id    = config.main_default_ruleset_id != null ? config.main_default_ruleset_id : local.repository_defaults.main_default_ruleset_id
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
