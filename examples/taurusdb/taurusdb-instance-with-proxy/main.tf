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
  availability_zone_mode = "multi"
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

# Generate random password if not provided
resource "random_password" "test" {
  count = var.instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
}

# Create TaurusDB instance
resource "huaweicloud_taurusdb_instance" "test" {
  name                             = var.instance_name
  password                         = var.instance_password != "" ? var.instance_password : try(random_password.test[0].result, null)
  flavor                           = var.instance_flavor_ref != "" ? var.instance_flavor_ref : try(data.huaweicloud_taurusdb_flavors.test.flavors[0].name, "")
  vpc_id                           = huaweicloud_vpc.test.id
  subnet_id                        = huaweicloud_vpc_subnet.test.id
  security_group_id                = huaweicloud_networking_secgroup.test.id
  mode                             = "Cluster"
  availability_zone_mode           = "multi"
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
      password, ssl_option, reserve_audit_logs
    ]
  }
}

# Query proxy flavors after the instance is created
data "huaweicloud_taurusdb_proxy_flavors" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
}

locals {
  # Sort nodes by name for consistent weight assignment
  sort_nodes = tolist(values({ for node in huaweicloud_taurusdb_instance.test.nodes : node.name => node }))
}

# Create TaurusDB database proxy
resource "huaweicloud_taurusdb_proxy" "test" {
  instance_id              = huaweicloud_taurusdb_instance.test.id
  flavor                   = try(data.huaweicloud_taurusdb_proxy_flavors.test.flavor_groups[0].flavors[0].spec_code, "")
  node_num                 = var.proxy_node_num
  proxy_name               = var.proxy_name
  proxy_mode               = var.proxy_mode
  route_mode               = var.route_mode
  subnet_id                = huaweicloud_vpc_subnet.test.id
  new_node_auto_add_status = var.proxy_new_node_auto_add_status
  new_node_weight          = var.proxy_new_node_weight
  port                     = var.proxy_port
  transaction_split        = var.proxy_transaction_split
  consistence_mode         = var.proxy_consistence_mode
  connection_pool_type     = var.proxy_connection_pool_type
  open_access_control      = var.proxy_open_access_control
  access_control_type      = var.proxy_open_access_control ? var.access_control_type : null
  dns_name_prefix          = var.proxy_dns_name_prefix

  dynamic "access_control_ip_list" {
    for_each = var.proxy_access_control_ip_list

    content {
      ip          = access_control_ip_list.value.ip
      description = access_control_ip_list.value.description
    }
  }

  master_node_weight {
    id     = try(local.sort_nodes[0].id, "")
    weight = var.proxy_master_node_weight
  }

  dynamic "readonly_nodes_weight" {
    for_each = length(local.sort_nodes) > 1 ? slice(local.sort_nodes, 1, length(local.sort_nodes)) : []

    content {
      id     = readonly_nodes_weight.value.id
      weight = var.proxy_readonly_node_weight
    }
  }

  dynamic "parameters" {
    for_each = var.proxy_parameters

    content {
      name      = parameters.value.name
      value     = parameters.value.value
      elem_type = parameters.value.elem_type
    }
  }

  lifecycle {
    ignore_changes = [
      new_node_weight, proxy_mode, readonly_nodes_weight, parameters,
    ]
  }
}
