data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].ids[0], null)
  visibility = var.image_visibility
  os         = var.image_os
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
    image    = var.instance_image_id != "" ? var.instance_image_id : try(data.huaweicloud_images_images.test[0].images[0].id, null)
    flavor   = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null)
    key_name = huaweicloud_kps_keypair.test.id

    dynamic "disk" {
      for_each = var.disk_configurations

      content {
        disk_type   = disk.value.disk_type
        volume_type = disk.value.volume_type
        size        = disk.value.volume_size
      }
    }
  }
}

resource "huaweicloud_as_group" "test" {
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = huaweicloud_vpc.test.id
  scaling_group_name       = var.group_name
  desire_instance_number   = var.desire_instance_number
  min_instance_number      = var.min_instance_number
  max_instance_number      = var.max_instance_number
  delete_publicip          = var.is_delete_publicip
  delete_instances         = var.is_delete_instances ? "yes" : "no"
  cool_down_time           = var.group_cool_down_time
  description              = var.group_description
  enterprise_project_id    = var.enterprise_project_id
  tags                     = var.group_tags

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = huaweicloud_as_group.test.id
  name             = var.task_name

  scheduled_policy {
    launch_time      = var.task_execute_time
    recurrence_type  = var.task_type
    recurrence_value = var.task_frequency
    start_time       = var.task_start_time
    end_time         = var.task_end_time
  }

  instance_number {
    min    = var.min_instance_number
    max    = var.max_instance_number
    desire = var.desire_instance_number
  }
}
