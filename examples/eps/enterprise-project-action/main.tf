resource "huaweicloud_enterprise_project" "test" {
  count = var.enterprise_project_id == "" ? 1 : 0

  name        = var.enterprise_project_name
  description = var.enterprise_project_description
  type        = var.enterprise_project_type
  enable      = var.enterprise_project_enable
  delete_flag = var.delete_flag

  lifecycle {
    ignore_changes = [
      enable
    ]
  }
}

resource "huaweicloud_enterprise_project_action" "test" {
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : try(huaweicloud_enterprise_project.test[0].id, "")
  action                = var.enterprise_project_action
}
