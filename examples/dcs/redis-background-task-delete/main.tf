resource "huaweicloud_dcs_background_task_delete" "test" {
  instance_id = var.dcs_instance_id
  task_id     = var.background_task_id
}
