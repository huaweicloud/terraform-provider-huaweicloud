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
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
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
    uuid = huaweicloud_vpc_subnet.test.id
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
  source_vpc_id            = huaweicloud_vpc.test.id
}

resource "huaweicloud_sdrs_protected_instance" "test" {
  name                 = var.protected_instance_name
  group_id             = huaweicloud_sdrs_protection_group.test.id
  server_id            = huaweicloud_compute_instance.test.id
  cluster_id           = var.cluster_id
  primary_subnet_id    = huaweicloud_vpc_subnet.test.id
  primary_ip_address   = var.primary_ip_address
  delete_target_server = var.delete_target_server
  delete_target_eip    = var.delete_target_eip
  description          = var.protected_instance_description
  tags                 = var.protected_instance_tags
}
