resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = (var.subnet_gateway_ip != "" ? var.subnet_gateway_ip :
  cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1))
}

data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_dcs_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  cache_mode     = "single"
  capacity       = var.instance_capacity
  engine_version = var.instance_engine_version
}

# Generate random password if not provided
resource "random_password" "test" {
  count = var.instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
}

resource "huaweicloud_dcs_instance" "test" {
  name                  = var.instance_name
  engine                = "Redis"
  enterprise_project_id = var.enterprise_project_id
  engine_version        = var.instance_engine_version
  capacity              = var.instance_capacity
  flavor                = (var.instance_flavor_id != "" ? var.instance_flavor_id :
  try(data.huaweicloud_dcs_flavors.test[0].flavors[0].name, null))
  availability_zones    = (var.availability_zone != "" ? [var.availability_zone] :
  try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), null))
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  password              = (var.instance_password != "" ? var.instance_password :
  try(random_password.test[0].result, null))
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew

  # If you want to change the `flavor` or `availability_zones`, you need to delete it from the `lifecycle.ignore_changes`.
  lifecycle {
    ignore_changes = [
      flavor,
      availability_zones,
    ]
  }
}

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}

locals {
  replication_list = try(data.huaweicloud_dcs_instance_shards.test.group_list[0].replication_list, [])
}

resource "huaweicloud_dcs_sessions_kill" "test" {
  instance_id  = huaweicloud_dcs_instance.test.id
  node_id      = try(local.replication_list[0].node_id, "")
  client_addrs = var.client_addrs
}
