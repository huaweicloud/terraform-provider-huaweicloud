data "huaweicloud_ram_resource_permissions" "test" {
  resource_type   = var.query_resource_type != "" ? var.query_resource_type : null
  permission_type = var.query_permission_type
  name            = var.query_permission_name != "" ? var.query_permission_name : null
}

resource "huaweicloud_ram_resource_share_permission" "test" {
  count = length(data.huaweicloud_ram_resource_permissions.test.permissions)

  resource_share_id = var.resource_share_id
  permission_id     = data.huaweicloud_ram_resource_permissions.test.permissions[count.index].id
  replace           = var.permission_replace
}
