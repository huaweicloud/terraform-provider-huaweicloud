data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = var.volume_availability_zone != "" ? var.volume_availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
  name              = var.volume_name
  volume_type       = var.volume_type
  size              = var.volume_size
  description       = var.voluem_description
  multiattach       = var.vouleme_multiattach
}

resource "huaweicloud_evs_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = var.snapshot_name
  description = var.snapshot_description
  metadata    = var.snapshot_metadata
  force       = var.snapshot_force
}
