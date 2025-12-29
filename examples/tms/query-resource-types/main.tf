data "huaweicloud_tms_resource_types" "test" {
  service_name = var.exact_service_name != "" ? var.exact_service_name : null
}

locals {
  # All service names registered on TMS service (using fuzzy matching based on user-specified service name)
  regex_matched_service_names = distinct(var.fuzzy_service_name != "" ? [
    for v in data.huaweicloud_tms_resource_types.test.types[*].service_name : v if length(regexall(var.fuzzy_service_name, v)) > 0
  ] : data.huaweicloud_tms_resource_types.test.types[*].service_name)

  # All resource types (object including the resource type name and the service name to which the resource type
  # belongs) registered on TMS service (using fuzzy matching based on user-specified service name)
  regex_matched_resource_types_by_only_fuzzy_service_name = var.fuzzy_service_name != "" ? [
    for v in data.huaweicloud_tms_resource_types.test.types : v if length(regexall(var.fuzzy_service_name, v.service_name)) > 0
  ] : data.huaweicloud_tms_resource_types.test.types

  # All resource types (object including the resource type name and the service name to which the resource type
  # belongs) registered on TMS service (using fuzzy matching based on user-specified service name or resource type name)
  regex_matched_resource_types = var.fuzzy_resource_type_name != "" ? [
    for v in local.regex_matched_resource_types_by_only_fuzzy_service_name : v if length(regexall(var.fuzzy_resource_type_name, v.name)) > 0
  ] : local.regex_matched_resource_types_by_only_fuzzy_service_name
}
