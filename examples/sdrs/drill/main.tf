data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].ids[0], "") : var.instance_flavor_id
  os         = var.instance_image_os_type
  visibility = var.instance_image_visibility
}

resource "huaweicloud_vpc" "test" {
  count = 2

  name = count.index == 0 ? var.vpc_name : "${var.vpc_name}-drill"
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 2

  vpc_id            = huaweicloud_vpc.test[count.index].id
  name              = count.index == 0 ? var.subnet_name : "${var.subnet_name}-drill"
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 8, 0), 1) : var.subnet_gateway_ip
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.ecs_instance_name
  availability_zone  = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, "") : var.instance_flavor_id
  image_id           = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.instance_image_id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test[0].id
  }
}

data "huaweicloud_sdrs_domain" "test" {
  name = var.sdrs_domain_name
}

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = var.protection_group_name
  source_availability_zone = var.source_availability_zone != "" ? var.source_availability_zone : try(data.huaweicloud_availability_zones.test.names[0], null)
  target_availability_zone = var.target_availability_zone != "" ? var.target_availability_zone : try(data.huaweicloud_availability_zones.test.names[1], null)
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test[0].id
}

resource "huaweicloud_sdrs_protected_instance" "test" {
  name                 = var.protected_instance_name
  group_id             = huaweicloud_sdrs_protection_group.test.id
  server_id            = huaweicloud_compute_instance.test.id
  delete_target_server = true
  delete_target_eip    = true
}

resource "huaweicloud_sdrs_drill" "test" {
  name         = var.drill_name
  group_id     = huaweicloud_sdrs_protection_group.test.id
  drill_vpc_id = huaweicloud_vpc.test[1].id

  depends_on = [
    huaweicloud_sdrs_protected_instance.test,
    huaweicloud_vpc_subnet.test[1],
  ]
}
