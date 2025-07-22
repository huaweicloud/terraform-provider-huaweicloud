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
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  priority          = 1
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, "") : var.instance_flavor_id
  image_id           = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.instance_image_id
  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  admin_pass = var.administrator_password

  lifecycle {
    ignore_changes = [
      availability_zone,
      flavor_id,
      image_id,
      admin_pass,
    ]
  }
}

resource "huaweicloud_networking_vip" "test" {
  network_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_networking_vip_associate" "test" {
  vip_id   = huaweicloud_networking_vip.test.id
  port_ids = try([
    huaweicloud_compute_instance.test.network[0].port
  ], [])
}
