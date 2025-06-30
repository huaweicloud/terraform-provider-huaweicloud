data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id  = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  visibility = "public"
  os         = "Ubuntu"
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name       = var.kps_key_pair_name
  public_key = var.kps_public_key
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = var.as_configuration_name
  instance_config {
    image              = try(data.huaweicloud_images_images.test.images[0].id, "")
    flavor             = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
    key_name           = huaweicloud_kps_keypair.test.id
    security_group_ids = [huaweicloud_networking_secgroup.test.id]

    metadata  = var.as_metadata
    user_data = var.as_user_data

    dynamic "disk" {
      for_each = var.as_disks

      content {
        size        = disk.value["size"]
        volume_type = disk.value["volume_type"]
        disk_type   = disk.value["disk_type"]
      }
    }

    dynamic "public_ip" {
      for_each = var.as_public_ip

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
