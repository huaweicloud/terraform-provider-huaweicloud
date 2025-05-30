# List the availability zones in the current region.
data "huaweicloud_availability_zones" "test" {}

# List the ecs flavors in the specified availability zone.
data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

# Query the specified system image.
data "huaweicloud_images_images" "test" {
  flavor_id  = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  visibility = "public"
  os         = "Ubuntu"
}

# Create a security group.
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# Create a KPS keypair.
resource "huaweicloud_kps_keypair" "test" {
  name       = var.key_pair_name
  public_key = var.public_key
}

# Create a AS configuration.
resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = var.configuration_name

  instance_config {
    image              = try(data.huaweicloud_images_images.test.images[0].id, "")
    flavor             = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
    key_name           = huaweicloud_kps_keypair.test.id
    security_group_ids = [huaweicloud_networking_secgroup.test.id]

    metadata = {
      some_key = "some_value"
    }

    user_data = <<EOT
#!/bin/sh
echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
EOT

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    public_ip {
      eip {
        ip_type = "5_bgp"
        bandwidth {
          size          = 10
          share_type    = "PER"
          charging_mode = "traffic"
        }
      }
    }
  }
}
