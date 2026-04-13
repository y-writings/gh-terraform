variable "manage_security_and_analysis" {
  description = "Whether Terraform should manage vulnerability alerts and security_and_analysis settings"
  type        = bool
  default     = false
}

variable "vulnerability_alerts" {
  description = "Dependabot alerts setting applied when manage_security_and_analysis is true"
  type        = bool
  default     = true
}

variable "secret_scanning_status" {
  description = "Secret scanning status applied when manage_security_and_analysis is true"
  type        = string
  default     = "enabled"

  validation {
    condition     = contains(["enabled", "disabled"], var.secret_scanning_status)
    error_message = "secret_scanning_status must be enabled or disabled."
  }
}

variable "secret_scanning_push_protection_status" {
  description = "Push protection status applied when manage_security_and_analysis is true"
  type        = string
  default     = "enabled"

  validation {
    condition     = contains(["enabled", "disabled"], var.secret_scanning_push_protection_status)
    error_message = "secret_scanning_push_protection_status must be enabled or disabled."
  }
}
