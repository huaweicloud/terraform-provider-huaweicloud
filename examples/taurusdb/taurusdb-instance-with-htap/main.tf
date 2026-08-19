# Build the network for the TaurusDB instance and the HTAP StarRocks instance.
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

# Security group for the TaurusDB instance.
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# Generate a random password for the TaurusDB instance and the HTAP StarRocks instance.
resource "random_password" "test" {
  length           = 16
  special          = true
  override_special = "!@%^*-_=+"
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
}

# Query the availability zones, TaurusDB flavors and HTAP engine resources.
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_taurusdb_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = var.taurusdb_availability_zone_mode
}

data "huaweicloud_taurusdb_htap_flavors" "test" {
  engine_name            = "star-rocks"
  availability_zone_mode = "single"
}

data "huaweicloud_taurusdb_htap_datastores" "test" {
  engine_name = "star-rocks"
}

locals {
  # Get the first available AZ from the HTAP flavor's az_status.
  htap_az_status = try(data.huaweicloud_taurusdb_htap_flavors.test.flavors[0].az_status, {})
  all_az_normal  = [for k, v in local.htap_az_status : k if v == "normal"]
  az_code        = var.az_code != "" ? var.az_code : try(local.all_az_normal[0], "")

  # Filter the frontend and backend flavors by spec code and AZ status.
  be_flavors = [for f in data.huaweicloud_taurusdb_htap_flavors.test.flavors : f if length(regexall("sr-be", f.spec_code)) > 0 && contains(keys(f.az_status), local.az_code) && f.az_status[local.az_code] == "normal"]
  fe_flavors = [for f in data.huaweicloud_taurusdb_htap_flavors.test.flavors : f if length(regexall("sr-fe", f.spec_code)) > 0 && contains(keys(f.az_status), local.az_code) && f.az_status[local.az_code] == "normal"]

  fe_flavor_id   = var.fe_flavor_id != "" ? var.fe_flavor_id : try(local.fe_flavors[0].id, "")
  be_flavor_id   = var.be_flavor_id != "" ? var.be_flavor_id : try(local.be_flavors[0].id, "")
  engine_version = var.engine_version != "" ? var.engine_version : try(data.huaweicloud_taurusdb_htap_datastores.test.datastores[0].name, "")
}

# The TaurusDB instance that the HTAP StarRocks instance depends on.
resource "huaweicloud_taurusdb_instance" "test" {
  name                     = var.taurusdb_instance_name
  flavor                   = var.taurusdb_flavor_ref != "" ? var.taurusdb_flavor_ref : try(data.huaweicloud_taurusdb_flavors.test.flavors[0].name, "")
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  password                 = var.taurusdb_root_pwd != "" ? var.taurusdb_root_pwd : random_password.test.result
  enterprise_project_id    = var.enterprise_project_id
  availability_zone_mode   = var.taurusdb_availability_zone_mode
  master_availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "")
  read_replicas            = var.taurusdb_read_replicas

  datastore {
    engine  = "gaussdb-mysql"
    version = "8.0"
  }

  lifecycle {
    ignore_changes = [
      password, datastore[0].version,
    ]
  }
}

# The HTAP StarRocks instance.
resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id          = huaweicloud_taurusdb_instance.test.id
  name                 = var.htap_instance_name
  fe_flavor_id         = local.fe_flavor_id
  be_flavor_id         = local.be_flavor_id
  az_code              = local.az_code
  az_mode              = "single"
  fe_count             = var.fe_count
  be_count             = var.be_count
  db_root_pwd          = var.htap_db_root_pwd != "" ? var.htap_db_root_pwd : random_password.test.result
  time_zone            = var.time_zone
  security_group_id    = huaweicloud_networking_secgroup.test.id
  enable_users_sync    = var.enable_users_sync
  open_slow_log_switch = var.open_slow_log_switch

  engine {
    type    = "star-rocks"
    version = local.engine_version
  }

  ha {
    mode = var.ha_mode
  }

  fe_volume {
    io_type        = var.volume_io_type
    capacity_in_gb = var.fe_volume_capacity
  }

  be_volume {
    io_type        = var.volume_io_type
    capacity_in_gb = var.be_volume_capacity
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = var.enterprise_project_id
    }
  }

  be_parameter_values = var.be_parameter_values
  fe_parameter_values = var.fe_parameter_values

  lifecycle {
    ignore_changes = [
      db_root_pwd, be_parameter_values, fe_parameter_values,
    ]
  }
}
