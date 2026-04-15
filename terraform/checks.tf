check "shared_governance_requires_public_repositories" {
  assert {
    condition = alltrue([
      for repository in values(local.repositories) : repository.visibility == "public"
    ])

    error_message = "The fixed governance baseline requires every managed repository visibility to be public for personal accounts."
  }
}
