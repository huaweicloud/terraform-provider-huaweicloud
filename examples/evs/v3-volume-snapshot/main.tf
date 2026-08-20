data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  count = var.volume_image_id == "" ? 1 : 0

  visibility = var.volume_image_visibility
  os         = var.volume_image_os
}

resource "huaweicloud_evsv3_volume" "test" {
  region            = var.region_name
  volume_type       = var.volume_type
  availability_zone = var.volume_availability_zone != "" ? var.volume_availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
  description       = var.volume_description
  image_id          = var.volume_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.volume_image_id
  metadata          = var.volume_metadata
  multiattach       = var.volume_multiattach
  name              = var.volume_name
  size              = var.volume_size
  tags              = var.volume_tags
}

resource "huaweicloud_evsv3_snapshot" "test" {
  region      = var.region_name
  volume_id   = huaweicloud_evsv3_volume.test.id
  name        = var.snapshot_name
  metadata    = var.snapshot_metadata
  description = var.snapshot_description
}
