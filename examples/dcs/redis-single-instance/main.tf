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

data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_dcs_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  cache_mode     = "single"
  capacity       = var.instance_flavor_id == "" ? var.instance_capacity : null
  engine_version = var.instance_flavor_id == "" ? var.instance_engine_version : null
}

resource "huaweicloud_dcs_instance" "test" {
  name                  = var.instance_name
  engine                = "Redis"
  enterprise_project_id = var.enterprise_project_id
  engine_version        = var.instance_engine_version
  capacity              = var.instance_flavor_id == "" ? data.huaweicloud_dcs_flavors.test[0].flavors[0].capacity : var.instance_flavor_id
  flavor                = var.instance_flavor_id == "" ? data.huaweicloud_dcs_flavors.test[0].flavors[0].name : var.instance_flavor_id
  availability_zones    = var.availability_zone == "" ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), null) : [var.availability_zone]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  password              = var.instance_password
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew
}
