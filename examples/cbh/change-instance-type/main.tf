resource "huaweicloud_cbh_change_instance_type" "test" {
  server_id         = var.server_id
  availability_zone = var.availability_zone
}
