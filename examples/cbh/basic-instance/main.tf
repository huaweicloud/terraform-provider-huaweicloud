data "huaweicloud_cbh_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_cbh_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0
  type  = var.instance_flavor_type
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

resource "huaweicloud_cbh_instance" "test" {
  name              = var.instance_name
  flavor_id         = var.instance_flavor_id == "" ? try(data.huaweicloud_cbh_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_cbh_availability_zones.test[0].availability_zones[0].name, null) : var.availability_zone
  password          = var.instance_password
  charging_mode     = var.charging_mode
  period_unit       = var.period_unit
  period            = var.period
  auto_renew        = var.auto_renew

  # If you want to change some of the following parameters, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zone,
    ]
  }
}
