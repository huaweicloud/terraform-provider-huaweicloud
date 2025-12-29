resource "huaweicloud_bms_instance_password_reset" "test" {
  server_id    = var.bms_instance_id
  new_password = var.bms_instance_new_password
}
