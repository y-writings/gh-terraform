variable "onepassword_account" {
  description = "1Password account sign-in address or account ID for the local onepassword provider"
  type        = string

  validation {
    condition     = trimspace(var.onepassword_account) != ""
    error_message = "onepassword_account must be a non-empty string."
  }
}
