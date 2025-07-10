locals {
  az          = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  subnet_cidr = var.subnet_cidr == "" ? cidrsubnet(try(huaweicloud_vpc.test[0].cidr, "192.168.0.0/16"), 8, 0) : var.subnet_cidr
  gateway     = var.gateway == "" ? cidrhost(cidrsubnet(try(huaweicloud_vpc.test[0].cidr, "192.168.0.0/16"), 8, 0), 1) : var.gateway
}

data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_rds_flavors" "test" {
  count = var.flavor_id == "" ? 1 : 0

  db_type           = var.db_type
  db_version        = var.db_version
  instance_mode     = var.instance_mode
  group_type        = var.group_type
  vcpus             = var.vcpus
  availability_zone = local.az
}

resource "huaweicloud_vpc" "test" {
  count = var.vpc_id == "" ? 1 : 0

  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = var.subnet_id == "" ? 1 : 0

  vpc_id            = try(huaweicloud_vpc.test[0].id, null)
  name              = var.subnet_name
  cidr              = local.subnet_cidr
  gateway_ip        = local.gateway
  availability_zone = local.az
}

resource "huaweicloud_networking_secgroup" "test" {
  count = var.secgroup_id == "" ? 1 : 0

  name = var.secgroup_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = var.secgroup_id == "" ? 1 : 0

  security_group_id = try(huaweicloud_networking_secgroup.test[0].id, null)
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = var.vpc_cidr
  port_range_max    = var.db_port
  port_range_min    = var.db_port
  protocol          = "tcp"
}

resource "random_password" "test" {
  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
}

resource "huaweicloud_rds_instance" "test" {
  name              = var.instance_name
  flavor            = var.flavor_id != "" ? var.flavor_id : try(data.huaweicloud_rds_flavors.test[0].flavors[0].name, null)
  vpc_id            = var.vpc_id != "" ? var.vpc_id : huaweicloud_vpc.test[0].id
  subnet_id         = var.subnet_id != "" ? var.subnet_id : huaweicloud_vpc_subnet.test[0].id
  security_group_id = var.secgroup_id != "" ? var.secgroup_id : huaweicloud_networking_secgroup.test[0].id
  charging_mode     = var.charging_mode
  availability_zone = [local.az]

  db {
    type     = var.db_type
    version  = var.db_version
    password = random_password.test.result
    port     = var.db_port
  }

  volume {
    type = var.volume_type
    size = var.volume_size
  }

  backup_strategy {
    start_time = var.backup_time_window
    keep_days  = var.backup_keep_days
  }
}

resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = var.account_name
  password    = random_password.test.result
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = var.db_name
  character_set = var.character_set
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = var.db_name

  users {
    name     = huaweicloud_rds_mysql_account.test.name
    readonly = true
  }
}

resource "huaweicloud_rds_backup" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = var.backup_name

  depends_on = [ huaweicloud_rds_mysql_database_privilege.test]
}
