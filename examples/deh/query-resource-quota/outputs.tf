output "quotas_with_usage" {
  description = "The quotas that have been used"
  value       = local.quotas_with_usage
}

output "quotas_available" {
  description = "The quotas that have available capacity"
  value       = local.quotas_available
}

output "quotas_exhausted" {
  description = "The quotas that are fully used"
  value       = local.quotas_exhausted
}
