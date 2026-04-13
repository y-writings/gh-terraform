check "shared_governance_requires_public_repositories" {
  assert {
    condition = (
      !module.advanced_security.manage_security_and_analysis &&
      local.repository_governance.required_code_scanning == null
      ) || alltrue([
        for repository in values(local.repositories) : repository.visibility == "public"
    ])

    error_message = "The current shared baseline assumes public repositories when the module-owned advanced_security baseline or repository_governance.required_code_scanning is enabled. For personal accounts, set every managed repository visibility to public before applying."
  }
}
