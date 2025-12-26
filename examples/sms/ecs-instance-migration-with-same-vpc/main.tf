data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "22"
  remote_ip_prefix  = huaweicloud_vpc_subnet.test.cidr
  priority          = 1
  action            = "allow"
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].ids[0], null)
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  image_id           = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  flavor_id          = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  availability_zone  = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  admin_pass         = var.instance_admin_password
  system_disk_size   = var.instance_system_disk_size
  system_disk_type   = var.instance_system_disk_type

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

# ST.001 Disable
resource "huaweicloud_compute_instance" "destination_server" {
  # ST.001 Enable
  name               = var.destination_instance_name
  image_id           = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
  flavor_id          = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
  availability_zone  = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  admin_pass         = var.instance_admin_password
  system_disk_size   = var.destination_instance_system_disk_size
  system_disk_type   = var.destination_instance_system_disk_type

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_sms_source_server" "test" {
  ip                 = huaweicloud_compute_instance.test.access_ip_v4
  name               = var.source_server_name
  os_type            = "LINUX"
  os_version         = var.source_server_os_version
  firmware           = var.source_server_firmware
  boot_loader        = var.source_server_boot_loader
  has_rsync          = var.source_server_has_rsync
  paravirtualization = var.source_server_paravirtualization
  cpu_quantity       = var.source_server_cpu_quantity
  memory             = var.source_server_memory
  agent_version      = var.source_server_agent_version

  dynamic "disks" {
    for_each = var.source_server_disks != null ? [var.source_server_disks] : []

    content {
      name            = disks.value.name
      partition_style = disks.value.partition_style
      device_use      = disks.value.device_use
      size            = disks.value.size
      used_size       = disks.value.used_size

      dynamic "physical_volumes" {
        for_each = disks.value.physical_volumes
        content {
          device_use  = physical_volumes.value.device_use
          file_system = physical_volumes.value.file_system
          mount_point = physical_volumes.value.mount_point
          name        = physical_volumes.value.name
        }
      }
    }
  }

  networks {
    name    = "eth0"
    ip      = huaweicloud_compute_instance.test.access_ip_v4
    netmask = cidrnetmask(huaweicloud_vpc_subnet.test.cidr)
    gateway = huaweicloud_vpc_subnet.test.gateway_ip
    mtu     = 1
    mac     = huaweicloud_compute_instance.test.network[0].mac
  }
}

resource "huaweicloud_sms_task" "test" {
  type             = var.migrate_task_type
  os_type          = "LINUX"
  source_server_id = huaweicloud_sms_source_server.test.id
  auto_start       = var.task_auto_start
  target_server_id = huaweicloud_compute_instance.destination_server.id
  migration_ip     = huaweicloud_compute_instance.destination_server.access_ip_v4
  action           = var.task_action

  dynamic "target_server_disks" {
    for_each = var.task_target_server_disks != null ? [var.task_target_server_disks] : []

    content {
      name        = target_server_disks.value.name
      size        = target_server_disks.value.size
      device_type = target_server_disks.value.device_type
      disk_id     = huaweicloud_compute_instance.destination_server.system_disk_id

      dynamic "physical_volumes" {
        for_each = target_server_disks.value.physical_volumes

        content {
          name        = physical_volumes.value.name
          size        = physical_volumes.value.size
          device_type = physical_volumes.value.device_type
          index       = physical_volumes.value.volume_index
          file_system = physical_volumes.value.file_system
          mount_point = physical_volumes.value.mount_point
        }
      }
    }
  }

  dynamic "configurations" {
    for_each = var.task_configurations

    content {
      config_key   = configurations.value.config_key
      config_value = configurations.value.config_value
    }
  }
}
