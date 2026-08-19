# Create VPC and subnet for the TaurusDB instance
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

# Query available flavors to get the instance flavor and availability zones
data "huaweicloud_taurusdb_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = var.availability_zone_mode
}

locals {
  # Get available AZs from the flavor's az_status (TaurusDB AZs must come from flavors, not huaweicloud_availability_zones)
  available_azs = try([for k, v in data.huaweicloud_taurusdb_flavors.test.flavors[0].az_status : k if v == "normal"], [])
  master_az     = var.master_availability_zone != "" ? var.master_availability_zone : try(local.available_azs[0], "")
}

# Create security group for the TaurusDB instance
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

# Create TaurusDB instance with SQL filter enabled for SQL control rules and auto throttling
resource "huaweicloud_taurusdb_instance" "test" {
  name                             = var.instance_name
  password                         = var.instance_password != "" ? var.instance_password : try(random_password.test[0].result, null)
  flavor                           = var.instance_flavor_ref != "" ? var.instance_flavor_ref : try(data.huaweicloud_taurusdb_flavors.test.flavors[0].name, "")
  vpc_id                           = huaweicloud_vpc.test.id
  subnet_id                        = huaweicloud_vpc_subnet.test.id
  security_group_id                = huaweicloud_networking_secgroup.test.id
  mode                             = var.instance_mode
  availability_zone_mode           = var.availability_zone_mode
  master_availability_zone         = local.master_az
  read_replicas                    = var.read_replicas
  enterprise_project_id            = var.enterprise_project_id
  volume_type                      = var.volume_type
  time_zone                        = var.time_zone
  port                             = var.instance_db_port
  ssl_option                       = var.ssl_option
  sql_filter_enabled               = var.sql_filter_enabled
  slow_log_show_original_switch    = var.slow_log_show_original_switch
  table_name_case_sensitivity      = var.table_name_case_sensitivity
  multi_tenant_switch              = var.multi_tenant_switch
  maintain_begin                   = var.maintain_begin
  maintain_end                     = var.maintain_end
  description                      = var.description
  seconds_level_monitoring_enabled = var.seconds_level_monitoring_enabled
  seconds_level_monitoring_period  = var.seconds_level_monitoring_enabled ? var.seconds_level_monitoring_period : null
  audit_log_enabled                = var.audit_log_enabled
  audit_log_keep_days              = var.audit_log_keep_days
  reserve_audit_logs               = var.reserve_audit_logs

  datastore {
    engine  = "gaussdb-mysql"
    version = "8.0"
  }

  backup_strategy {
    start_time = var.instance_backup_time_window
    keep_days  = tostring(var.instance_backup_keep_days)
  }

  tags = var.tags

  lifecycle {
    ignore_changes = [
      password, ssl_option, reserve_audit_logs, datastore[0].version,
    ]
  }
}

locals {
  # Get the master node ID for SQL control rule and auto throttling
  # The auto throttling feature is only available for primary nodes
  master_node_id = try([for node in huaweicloud_taurusdb_instance.test.nodes : node.id if node.type == "master"][0], "")
}

# Create TaurusDB SQL concurrency control rule
# Requires sql_filter_enabled = true on the instance
resource "huaweicloud_taurusdb_sql_control_rule" "test" {
  instance_id     = huaweicloud_taurusdb_instance.test.id
  node_id         = local.master_node_id
  sql_type        = var.sql_control_rule_sql_type
  pattern         = var.sql_control_rule_pattern
  max_concurrency = var.sql_control_rule_max_concurrency
}

# Create TaurusDB SQL auto throttling
# This feature is only available for primary (master) nodes
resource "huaweicloud_taurusdb_sql_auto_throttling" "test" {
  instance_id     = huaweicloud_taurusdb_instance.test.id
  node_id         = local.master_node_id
  start_time      = var.sql_auto_throttling_start_time
  end_time        = var.sql_auto_throttling_end_time
  condition       = var.sql_auto_throttling_condition
  cpu_usage       = var.sql_auto_throttling_cpu_usage
  active_sessions = var.sql_auto_throttling_active_sessions
  clear_time      = var.sql_auto_throttling_clear_time
  duration        = var.sql_auto_throttling_duration
  max_concurrency = var.sql_auto_throttling_max_concurrency
  retain_sql_rule = var.sql_auto_throttling_retain_sql_rule

  lifecycle {
    ignore_changes = [
      # instance_id and node_id are NoneUpdatable; ignore to prevent destroy failure
      # when the referenced instance resource is being destroyed
      instance_id, node_id, retain_sql_rule,
    ]
  }

  depends_on = [huaweicloud_taurusdb_sql_control_rule.test]
}
