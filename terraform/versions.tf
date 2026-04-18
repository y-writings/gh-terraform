terraform {
  required_version = ">= 1.7.0"

  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }

    onepassword = {
      source  = "1Password/onepassword"
      version = "~> 3.3"
    }
  }
}
