locals {
  # Filter the standard editions
  se_version_names = [for n in data.huaweicloud_rds_engine_versions.default.versions[*].name : n if length(regexall("^\\d+_SE$", n)) > 0]
  # Filter the web editions
  web_version_names = [for n in data.huaweicloud_rds_engine_versions.default.versions[*].name : n if length(regexall("^\\d+_WEB$", n)) > 0]
}

data "huaweicloud_availability_zones" "default" {}

# Query an available engine version list of the SQL server database.
data "huaweicloud_rds_engine_versions" "default" {
  type = "SQLServer"
}

# Query an available flavor of web edition using specified vcpus and memory numbers.
data "huaweicloud_rds_flavors" "default" {
  db_type = "SQLServer"
  # Use the latest version of the RDS engine version
  db_version    = element(local.web_version_names, length(local.web_version_names) - 1)
  instance_mode = "single"
}

# Query an available flavor of standard edition using specified vcpus and memory numbers.
data "huaweicloud_rds_flavors" "ha_database" {
  db_type = "SQLServer"
  # Use latest version of the RDS engine versions
  db_version    = element(local.se_version_names, length(local.se_version_names) - 1)
  instance_mode = "ha"
}

resource "huaweicloud_vpc" "default" {
  name = var.vpc_name
  cidr = "192.168.128.0/20"
}

resource "huaweicloud_vpc_subnet" "default" {
  name       = var.subnet_name
  vpc_id     = huaweicloud_vpc.default.id
  cidr       = "192.168.128.0/24"
  gateway_ip = "192.168.128.1"
}

resource "huaweicloud_networking_secgroup" "default" {
  name = var.security_group_name
}

resource "random_password" "password" {
  length           = 32
  special          = true
  min_upper        = 1
  min_special      = 1
  min_numeric      = 1
  override_special = "_%@"
}

# Create a single-type database
resource "huaweicloud_rds_instance" "default" {
  name              = var.single_instance_name
  flavor            = data.huaweicloud_rds_flavors.default.flavors[0].name
  vpc_id            = huaweicloud_vpc.default.id
  subnet_id         = huaweicloud_vpc_subnet.default.id
  security_group_id = huaweicloud_networking_secgroup.default.id

  availability_zone = [
    data.huaweicloud_availability_zones.default.names[0],
  ]

  db {
    type     = "SQLServer"
    version  = element(local.web_version_names, length(local.web_version_names) - 1)
    password = random_password.password.result
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}

# Create an HA-type database
resource "huaweicloud_rds_instance" "ha_database" {
  name                = var.ha_instance_name
  flavor              = data.huaweicloud_rds_flavors.ha_database.flavors[0].name
  vpc_id              = huaweicloud_vpc.default.id
  subnet_id           = huaweicloud_vpc_subnet.default.id
  security_group_id   = huaweicloud_networking_secgroup.default.id
  ha_replication_mode = "sync"

  availability_zone = [
    data.huaweicloud_availability_zones.default.names[0],
    data.huaweicloud_availability_zones.default.names[1],
  ]

  db {
    type     = "SQLServer"
    version  = element(local.se_version_names, length(local.se_version_names) - 1)
    password = random_password.password.result
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 7
  }
}
