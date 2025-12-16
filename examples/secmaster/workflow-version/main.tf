data "huaweicloud_secmaster_workspaces" "test" {
  count = var.workspace_id == "" ? 1 : 0

  name = var.workspace_name
}

data "huaweicloud_secmaster_workflows" "test" {
  count = var.workflow_id == "" ? 1 : 0

  workspace_id = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  # ST.002 Disable
  name         = var.workflow_name
  # ST.002 Enable
}

resource "huaweicloud_secmaster_workflow_version" "test" {
  workspace_id  = var.workspace_id != "" ? var.workspace_id : try(data.huaweicloud_secmaster_workspaces.test[0].workspaces[0].id, null)
  workflow_id   = var.workflow_id != "" ? var.workflow_id : try(data.huaweicloud_secmaster_workflows.test[0].workflows[0].id, null)
  name          = var.workflow_name
  taskflow      = var.workflow_version_taskflow
  taskconfig    = var.workflow_version_taskconfig
  taskflow_type = var.workflow_version_taskflow_type
  aop_type      = var.workflow_version_aop_type
  description   = var.workflow_version_description
}
