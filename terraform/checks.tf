check "shared_governance_requires_public_repositories" {
  assert {
    condition = (
      !local.repository_governance.manage_security_and_analysis &&
      local.repository_governance.required_code_scanning == null
      ) || alltrue([
        for repository in values(local.repositories) : repository.visibility == "public"
    ])

    error_message = "The current shared repository_governance baseline assumes public repositories when security_and_analysis or required_code_scanning is enabled. For personal accounts, set every managed repository visibility to public or change the global repository_governance baseline before applying."
  }
}
