resource "huaweicloud_secmaster_workspace" "test" {
  name                  = var.workspace_name
  project_name          = var.workspace_project_name
  description           = var.workspace_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.workspace_tags
}
