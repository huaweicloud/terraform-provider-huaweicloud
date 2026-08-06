resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = 2

  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  remote_ip_prefix  = "192.168.0.0/16"
  protocol          = "tcp"
  direction         = count.index == 0 ? "ingress" : "egress"
  ports             = count.index == 0 ? "3306" : null
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = var.rds_db_type
  db_version    = var.rds_db_version
  instance_mode = var.rds_instance_mode
}

# RDS MySQL instance
resource "huaweicloud_rds_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.test,
  ]

  name              = var.rds_name
  flavor            = var.rds_flavor != "" ? var.rds_flavor : try(data.huaweicloud_rds_flavors.test.flavors[0].name, null)
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  fixed_ip          = var.rds_fixed_ip

  availability_zone = [
    try(data.huaweicloud_availability_zones.test.names[0], ""),
  ]

  db {
    password = var.db_password
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

# DRS connection for MySQL
resource "huaweicloud_drs_connection" "test" {
  name        = var.connection_name
  db_type     = "mysql"
  description = var.description

  endpoint {
    endpoint_name = "cloud_mysql"
    instance_id   = huaweicloud_rds_instance.test.id
    db_port       = var.db_port
    db_user       = var.db_user
    db_password   = var.db_password
  }

  vpc {
    vpc_id            = huaweicloud_rds_instance.test.vpc_id
    subnet_id         = huaweicloud_rds_instance.test.subnet_id
    security_group_id = huaweicloud_networking_secgroup.test.id
  }

  ssl {
    ssl_link = false
  }

  config {
    driver_name = var.driver_name
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
    ]
  }
}
