resource "huaweicloud_gaussdb_agency_permission_policy" "test" {
  agency_name       = var.agency_name
  bind_role_names   = var.bind_role_names
  unbind_role_names = var.unbind_role_names
}
