locals {
  repository_defaults = {
    import_existing_repository = false
    main_default_ruleset_id    = null
  }

  repositories = {
    for name, config in var.repositories :
    name => {
      import_existing_repository = coalesce(config.import_existing_repository, local.repository_defaults.import_existing_repository)
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
