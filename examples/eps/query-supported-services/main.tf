data "huaweicloud_enterprise_project_services" "test" {
  locale  = var.exact_locale != "" ? var.exact_locale : null
  service = var.fuzzy_service_name != "" ? var.fuzzy_service_name : null
}

locals {
  matched_services = data.huaweicloud_enterprise_project_services.test.services

  # All service names supported by EPS (using fuzzy matching based on user-specified service name)
  matched_service_by_only_exact_service_name = var.exact_service_name != "" ? [
    for v in data.huaweicloud_enterprise_project_services.test.services : v if v.service == var.exact_service_name
  ] : data.huaweicloud_enterprise_project_services.test.services
}
