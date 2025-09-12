data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  performance_type  = var.instance_performance_type
  cpu_core_count    = var.instance_cpu_core_count
  memory_size       = var.instance_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  count = var.security_group_name != "" ? 1 : 0

  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "test" {
  name                    = var.instance_name
  image_id                = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, null) : var.instance_image_id
  flavor_id               = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  security_group_ids      = length(var.security_group_ids) == 0 ? huaweicloud_networking_secgroup.test[*].id : var.security_group_ids
  availability_zone       = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  admin_pass              = var.instance_admin_password
  description             = var.instance_description
  system_disk_type        = var.instance_system_disk_type
  system_disk_size        = var.instance_system_disk_size
  system_disk_iops        = var.instance_system_disk_iops
  system_disk_throughput  = var.instance_system_disk_throughput
  system_disk_dss_pool_id = var.instance_system_disk_dss_pool_id
  metadata                = var.instance_metadata
  tags                    = var.instance_tags
  enterprise_project_id   = var.enterprise_project_id
  eip_id                  = var.instance_eip_id
  eip_type                = var.instance_eip_type

  dynamic "bandwidth" {
    for_each = var.instance_bandwidth == null ? [] : [var.instance_bandwidth]

    content {
      share_type   = bandwidth.value.share_type
      id           = bandwidth.value.id
      size         = bandwidth.value.size
      charge_mode  = bandwidth.value.charge_mode
      extend_param = bandwidth.value.extend_param
    }
  }

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  charging_mode = "prePaid"
  period_unit   = var.period_unit
  period        = var.period
  auto_renew    = var.auto_renew

  lifecycle {
    precondition {
      condition     = (length(var.security_group_ids) != 0) != (var.security_group_name != "")
      error_message = "Exactly one of security_group_ids or security_group_name must be provided. Please set only one."
    }
  }
}
