data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  volume_type       = var.volume_type
  name              = var.volume_name
  size              = var.volume_size
  device_type       = var.volume_device_type
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = var.name
  type                  = var.type
  protection_type       = var.protection_type
  size                  = var.size
  enterprise_project_id = var.enterprise_project_id

  resources {
    includes = [huaweicloud_evs_volume.test.id]
  }
}
