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
  count = length(var.instance_flavors_filter)

  db_type           = lookup(var.instance_flavors_filter[count.index], "db_type")
  db_version        = lookup(var.instance_flavors_filter[count.index], "db_version")
  instance_mode     = lookup(var.instance_flavors_filter[count.index], "instance_mode")
  group_type        = lookup(var.instance_flavors_filter[count.index], "group_type")
  vcpus             = lookup(var.instance_flavors_filter[count.index], "vcpus")
  memory            = lookup(var.instance_flavors_filter[count.index], "memory")
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
  flavor              = var.instance_flavor_id != "" ? var.instance_flavor_id : try([for o in data.huaweicloud_rds_flavors.test : o.flavors[0].name if o.instance_mode == "ha"][0], null)
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = length(var.availability_zones) > 0 ? var.availability_zones : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 2), [])
  ha_replication_mode = var.ha_replication_mode

  db {
    type     = try([for o in var.instance_flavors_filter : lookup(o, "db_type", "")][0], "MySQL")
    version  = try([for o in var.instance_flavors_filter : lookup(o, "db_version", "")][0], "MySQL")
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

resource "huaweicloud_rds_read_replica_instance" "test" {
  primary_instance_id = huaweicloud_rds_instance.test.id
  name                = var.replica_instance_name
  flavor              = var.replica_instance_flavor_id != "" ? var.replica_instance_flavor_id : try([for o in data.huaweicloud_rds_flavors.test : o.flavors[0].name if o.instance_mode == "replica"][0], null)
  availability_zone   = length(var.availability_zones) > 0 ? element(var.availability_zones, 1) : try(data.huaweicloud_availability_zones.test[0].names[0], null)

  volume {
    type = var.replica_instance_volume_type
    size = var.replica_instance_volume_size
  }

  lifecycle {
    ignore_changes = [
      flavor,
      availability_zone,
    ]
  }
}
