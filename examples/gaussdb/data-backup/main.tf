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
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
  ports             = var.security_group_rule_ports
  protocol          = "tcp"
}

data "huaweicloud_availability_zones" "test" {
  count = var.instance_availability_zones == "" ? 1 : 0
}

resource "random_password" "test" {
  length           = 16
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
  special          = true
  override_special = "~!@#%^*-_=+?"
}

# Create GaussDB instance
resource "huaweicloud_gaussdb_instance" "test" {
  name                  = var.instance_name
  flavor                = var.instance_flavor
  password              = var.instance_password != "" ? var.instance_password : random_password.test.result
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = var.instance_availability_zones != "" ? var.instance_availability_zones : join(",", slice(data.huaweicloud_availability_zones.test[0].names, 0, 3))
  port                  = var.instance_db_port
  enterprise_project_id = var.enterprise_project_id

  ha {
    mode             = var.instance_ha_mode
    replication_mode = var.instance_ha_replication_mode
    consistency      = var.instance_ha_consistency
  }

  replica_num = 3

  volume {
    type = var.instance_volume_type
    size = var.instance_volume_size
  }

  lifecycle {
    ignore_changes = [
      flavor,
    ]
  }
}

# Create a manual backup for the GaussDB instance
resource "huaweicloud_gaussdb_backup" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  name        = var.backup_name
  description = var.backup_description
}
