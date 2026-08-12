# SecMaster workspace
resource "huaweicloud_secmaster_workspace" "test" {
  name         = var.workspace_name
  project_name = var.region_name
  description  = var.workspace_description
}

# SecMaster dataspace
resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = huaweicloud_secmaster_workspace.test.id
  dataspace_name = var.dataspace_name
  description    = var.dataspace_description
}
