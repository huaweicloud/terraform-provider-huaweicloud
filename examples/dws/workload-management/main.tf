data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                  = var.security_group_name
  delete_default_rules  = var.security_group_delete_default_rules
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

# In order to allow DWS clients can connect to the cluster.
resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
  ports             = var.security_group_rule_ports
  protocol          = "tcp"
}

data "huaweicloud_dws_flavors" "test" {
  count = var.cluster_node_type == "" || var.cluster_version == "" ? 1 : 0

  availability_zone = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  vcpus             = var.cluster_vcpus
  memory            = var.cluster_memory
  datastore_type    = var.cluster_datastore_type
}

resource "huaweicloud_dws_cluster" "test" {
  name                  = var.cluster_name
  node_type             = var.cluster_node_type != "" ? var.cluster_node_type : try(data.huaweicloud_dws_flavors.test[0].flavors[0].flavor_id, null)
  number_of_node        = var.cluster_number_of_node
  number_of_cn          = var.cluster_number_of_cn
  version               = var.cluster_version != "" ? var.cluster_version : try(data.huaweicloud_dws_flavors.test[0].flavors[0].datastore_version, null)
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = var.availability_zone != "" ? var.availability_zone : try(data.huaweicloud_availability_zones.test[0].names[0], null)
  user_name             = var.cluster_admin_user_name
  user_pwd              = var.cluster_admin_user_pwd
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  volume {
    type     = var.cluster_volume_type
    capacity = var.cluster_volume_capacity
  }
}

resource "huaweicloud_dws_workload_queue" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  name       = var.workload_queue_name

  dynamic "configuration" {
    for_each = var.workload_queue_configurations

    content {
      resource_name  = configuration.value.resource_name
      resource_value = configuration.value.resource_value
    }
  }
}

resource "huaweicloud_dws_cluster_user" "test" {
  cluster_id   = huaweicloud_dws_cluster.test.id
  type         = "user"
  name         = var.user_name
  password     = var.user_password
  description  = var.user_description != "" ? var.user_description : null
  cascade      = var.user_cascade
  login        = var.user_login
  create_role  = var.user_create_role
  create_db    = var.user_create_db
  system_admin = var.user_system_admin
  audit_admin  = var.user_audit_admin
  inherit      = var.user_inherit
  use_ft       = var.user_use_ft
  conn_limit   = var.user_conn_limit
  replication  = var.user_replication
  valid_begin  = var.user_valid_begin
  valid_until  = var.user_valid_until

  dynamic "grant_list" {
    for_each = var.user_grant_list

    content {
      type                 = grant_list.value.type
      database             = grant_list.value.database
      schema_name          = grant_list.value.schema_name
      object_name          = grant_list.value.object_name
      all_object           = grant_list.value.all_object
      future               = grant_list.value.future
      future_object_owners = grant_list.value.future_object_owners
      column_names         = grant_list.value.column_names

      dynamic "privileges" {
        for_each = grant_list.value.privileges

        content {
          permission = privileges.value.permission
          grant_with = privileges.value.grant_with
        }
      }
    }
  }
}

resource "huaweicloud_dws_workload_queue_user_associate" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  queue_name = huaweicloud_dws_workload_queue.test.name
  user_names = [huaweicloud_dws_cluster_user.test.name]
}

resource "huaweicloud_dws_workload_plan" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  name       = var.workload_plan_name
}

resource "huaweicloud_dws_workload_plan_stage" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  plan_id    = huaweicloud_dws_workload_plan.test.id
  name       = var.workload_plan_stage_name
  month      = var.workload_plan_stage_month
  day        = var.workload_plan_stage_day
  start_time = var.workload_plan_stage_start_time
  end_time   = var.workload_plan_stage_end_time

  queues {
    name = huaweicloud_dws_workload_queue.test.name

    dynamic "configuration" {
      for_each = var.workload_plan_stage_configurations

      content {
        resource_name        = configuration.value.resource_name
        resource_value       = configuration.value.resource_value
        value_unit           = configuration.value.value_unit
        resource_description = configuration.value.resource_description
      }
    }
  }
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
  name       = var.exception_rule_name

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}
