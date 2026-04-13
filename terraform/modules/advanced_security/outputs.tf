output "manage_security_and_analysis" {
  description = "Whether Terraform should manage vulnerability alerts and security_and_analysis settings"
  value       = var.manage_security_and_analysis
}

output "vulnerability_alerts" {
  description = "Dependabot alerts setting applied when manage_security_and_analysis is true"
  value       = var.vulnerability_alerts
}

output "secret_scanning_status" {
  description = "Secret scanning status applied when manage_security_and_analysis is true"
  value       = var.secret_scanning_status
}

output "secret_scanning_push_protection_status" {
  description = "Push protection status applied when manage_security_and_analysis is true"
  value       = var.secret_scanning_push_protection_status
}
