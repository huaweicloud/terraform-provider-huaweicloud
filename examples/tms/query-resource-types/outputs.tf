output "regex_matched_service_names" {
  description = "The service names that match the regex pattern"
  value       = local.regex_matched_service_names
}

output "regex_matched_resource_types" {
  description = "The resource types that match the regex pattern"
  value       = local.regex_matched_resource_types
}
