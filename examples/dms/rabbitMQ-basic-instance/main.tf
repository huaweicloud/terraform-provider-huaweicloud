data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
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

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  type               = var.instance_flavor_type
  flavor_id          = var.instance_flavor_id
  storage_spec_code  = var.instance_storage_spec_code
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, var.availability_zone_number), null) : var.availability_zones
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name                  = var.instance_name
  engine_version        = var.instance_engine_version
  flavor_id             = var.instance_flavor_id == "" ? try(data.huaweicloud_dms_rabbitmq_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, var.availability_zone_number), null) : var.availability_zones
  broker_num            = var.instance_broker_num
  storage_space         = var.instance_storage_space
  storage_spec_code     = var.instance_storage_spec_code
  ssl_enable            = var.instance_ssl_enable
  access_user           = var.instance_access_user_name
  password              = var.instance_password
  description           = var.instance_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.instance_tags
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew

  # If you want to change some of the following parameters, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zones,
    ]
  }
}
