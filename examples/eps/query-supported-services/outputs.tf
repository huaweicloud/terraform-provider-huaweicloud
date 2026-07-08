output "matched_services" {
  description = "The services that match the filter criteria"
  value       = local.matched_services
}

output "matched_service_by_only_exact_service_name" {
  description = "The service that match the specified exact service name"
  value       = local.matched_service_by_only_exact_service_name
}
