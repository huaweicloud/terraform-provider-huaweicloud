data "huaweicloud_identity_projects" "test" {
  # ST.002 Disable
  name = var.region_name
  # ST.002 Enable
}

locals {
  exact_project_id = try([for v in data.huaweicloud_identity_projects.test.projects : v.id if v.name == var.region_name][0], null)
}

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = local.exact_project_id

  dynamic "resources" {
    for_each = var.associated_resources_configuration

    content {
      resource_type = resources.value.type
      resource_id   = resources.value.id
    }
  }

  tags = var.resource_tags
}
