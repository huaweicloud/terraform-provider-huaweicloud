resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) < 1 ? 1 : 0
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = var.subnet_name
  cidr              = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip        = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
  availability_zone = length(var.availability_zones) > 0 ? element(var.availability_zones, 0) : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

data "huaweicloud_rds_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  db_type           = var.instance_db_type
  db_version        = var.instance_db_version
  instance_mode     = var.instance_mode
  group_type        = var.instance_flavor_group_type
  vcpus             = var.instance_flavor_vcpus
  memory            = var.instance_flavor_memory
  availability_zone = length(var.availability_zones) > 0 ? element(var.availability_zones, 0) : try(data.huaweicloud_availability_zones.test[0].names[0], null)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = var.vpc_cidr
  ports             = var.instance_db_port
  protocol          = "tcp"
}

resource "random_password" "test" {
  count = var.instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
}

resource "huaweicloud_rds_instance" "test" {
  name                = var.instance_name
  flavor              = var.instance_flavor_id != "" ? var.instance_flavor_id : try(data.huaweicloud_rds_flavors.test[0].flavors[0].name, null)
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = length(var.availability_zones) > 0 ? var.availability_zones : var.instance_mode == "ha" ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 2), []) : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 1), [])
  ha_replication_mode = var.instance_mode == "ha" ? var.ha_replication_mode : null

  db {
    type     = var.instance_db_type
    version  = var.instance_db_version
    port     = var.instance_db_port
    password = var.instance_password != "" ? var.instance_password : try(random_password.test[0].result, null)
  }

  volume {
    type = var.instance_volume_type
    size = var.instance_volume_size
  }

  backup_strategy {
    start_time = var.instance_backup_time_window
    keep_days  = var.instance_backup_keep_days
  }

  lifecycle {
    ignore_changes = [
      flavor,
      availability_zone,
    ]
  }
}

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = var.account_name
  password    = var.account_password != "" ? var.account_password : try(random_password.test[0].result, null)
}

resource "huaweicloud_rds_pg_account_privileges" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  user_name              = huaweicloud_rds_pg_account.test.name
  role_privileges        = ["CREATEROLE", "CREATEDB", "LOGIN", "REPLICATION"]
  system_role_privileges = ["pg_signal_backend"]
}

resource "huaweicloud_rds_pg_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = var.database_name
}

resource "huaweicloud_rds_pg_schema" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_pg_database.test.name
  owner       = huaweicloud_rds_pg_account.test.name
  schema_name = var.schema_name
}

resource "huaweicloud_rds_backup" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = var.backup_name

  depends_on = [huaweicloud_rds_pg_schema.test]
}
