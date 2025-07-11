data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.configuration_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  performance_type  = var.configuration_flavor_performance_type
  cpu_core_count    = var.configuration_flavor_cpu_core_count
  memory_size       = var.configuration_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.configuration_image_id == "" ? 1 : 0

  flavor_id  = var.configuration_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.configuration_flavor_id
  visibility = var.configuration_image_visibility
  os         = var.configuration_image_os
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name       = var.keypair_name
  public_key = var.keypair_public_key != "" ? var.keypair_public_key : null
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = var.configuration_name

  instance_config {
    image              = var.configuration_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, null) : var.configuration_image_id
    flavor             = var.configuration_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.configuration_flavor_id
    key_name           = huaweicloud_kps_keypair.test.id
    security_group_ids = [huaweicloud_networking_secgroup.test.id]

    metadata  = var.configuration_metadata
    user_data = var.configuration_user_data

    dynamic "disk" {
      for_each = var.configuration_disks

      content {
        size        = disk.value["size"]
        volume_type = disk.value["volume_type"]
        disk_type   = disk.value["disk_type"]
      }
    }

    dynamic "public_ip" {
      for_each = var.configuration_public_ip

      content {
        eip {
          ip_type = public_ip.value.eip.ip_type
          bandwidth {
            size          = public_ip.value.eip.bandwidth.size
            share_type    = public_ip.value.eip.bandwidth.share_type
            charging_mode = public_ip.value.eip.bandwidth.charging_mode
          }
        }
      }
    }
  }
}
