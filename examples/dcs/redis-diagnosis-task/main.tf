resource "huaweicloud_dcs_diagnosis_task" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
}
