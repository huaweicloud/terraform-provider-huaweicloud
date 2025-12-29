resource "huaweicloud_bms_volume_attach" "test" {
  server_id = var.server_id
  volume_id = var.volume_id
  device    = var.device
}
