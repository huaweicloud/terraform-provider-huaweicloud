---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_instance"
description: |-
  Manages a TaurusDB HTAP StarRocks instance resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_instance

Manages a TaurusDB HTAP StarRocks instance resource within HuaweiCloud.

## Example Usage

### Create a postPaid Cluster StarRocks instance

```hcl
variable "taurusdb_instance_id" {}
variable "name" {}
variable "availability_zone" {}
variable "fe_flavor_id" {}
variable "be_flavor_id" {}
variable "db_root_pwd" {}
variable "enterprise_project_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = var.taurusdb_instance_id
  name              = var.name
  fe_flavor_id      = var.fe_flavor_id
  be_flavor_id      = var.be_flavor_id
  az_code           = var.availability_zone
  db_root_pwd       = var.db_root_pwd
  fe_count          = 3
  be_count          = 3
  az_mode           = "single"
  time_zone         = "UTC+08:00"
  enable_users_sync = "true"

  engine {
    type    = "star-rocks"
    version = "3.1.15.80"
  }

  ha {
    mode = "Cluster"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = var.enterprise_project_id
    }
  }
}
```

### Create a postPaid Single-Node StarRocks instance

```hcl
variable "taurusdb_instance_id" {}
variable "name" {}
variable "availability_zone" {}
variable "fe_flavor_id" {}
variable "be_flavor_id" {}
variable "db_root_pwd" {}
variable "enterprise_project_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = var.name
  fe_flavor_id      = var.fe_flavor_id
  be_flavor_id      = var.be_flavor_id
  az_code           = var.availability_zone
  db_root_pwd       = var.db_root_pwd
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  enable_users_sync = "true"

  engine {
    type    = "star-rocks"
    version = "3.1.15.80"
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = var.enterprise_project_id
    }
  }
}
```

### Create a prePaid Single-Node StarRocks instance

```hcl
variable "taurusdb_instance_id" {}
variable "name" {}
variable "availability_zone" {}
variable "fe_flavor_id" {}
variable "be_flavor_id" {}
variable "db_root_pwd" {}
variable "security_group_id" {}
variable "enterprise_project_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  instance_id       = huaweicloud_taurusdb_instance.test.id
  name              = var.name
  fe_flavor_id      = var.fe_flavor_id
  be_flavor_id      = var.be_flavor_id
  az_code           = var.availability_zone
  db_root_pwd       = var.db_root_pwd
  fe_count          = 1
  be_count          = 1
  az_mode           = "single"
  time_zone         = "UTC+07:00"
  security_group_id = var.security_group_id
  enable_users_sync = "true"
  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
  auto_renew        = "true"

  engine {
    type    = "star-rocks"
    version = "3.1.15.80"
  }

  ha {
    mode = "Single"
  }

  fe_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  be_volume {
    io_type        = "SSD"
    capacity_in_gb = 50
  }

  tags_info {
    sys_tags {
      key   = "_sys_enterprise_project_id"
      value = var.enterprise_project_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the TaurusDB instance ID.

* `name` - (Required, String, NoneUpdatable) Specifies the StarRocks instance name.
  The name must start with a letter and consist of `4` to `64` characters. Only letters (case-sensitive),
  digits, hyphens (-), and underscores (_) are allowed.

* `engine` - (Required, List, NoneUpdatable) Specifies the engine information.
  The [engine](#engine_block) structure is documented below.

* `ha` - (Required, List, NoneUpdatable) Specifies the deployment information.
  The [ha](#ha_block) structure is documented below.

* `fe_flavor_id` - (Required, String) Specifies the specification ID of the frontend node.
  The value can be obtained from the `id` field in the response to querying HTAP specifications.

* `be_flavor_id` - (Required, String) Specifies the specification ID of the backend node.
  The value can be obtained from the `id` field in the response to querying HTAP specifications.

* `db_root_pwd` - (Required, String) Specifies the database password.
  The password must consist of `8` to `32` characters and contain at least three types of the following:
  uppercase letters, lowercase letters, digits, and special characters (~!@#%^*-_=+?,()&).

* `fe_count` - (Required, Int, NoneUpdatable) Specifies the number of frontend nodes.
  For a single-node instance, the value is fixed to **1**. For a cluster instance, the value ranges from `3` to `10`.

* `be_count` - (Required, Int, NoneUpdatable) Specifies the number of backend nodes.
  For a single-node instance, the value is fixed to **1**. For a cluster instance, the value ranges from `3` to `10`.

* `az_mode` - (Required, String, NoneUpdatable) Specifies the AZ type. Currently, only **single** is supported.

* `fe_volume` - (Required, List, NoneUpdatable) Specifies the storage information of the frontend node.
  The [fe_volume](#volume_block) structure is documented below.

* `be_volume` - (Required, List, NoneUpdatable) Specifies the storage information of the backend node.
  The [be_volume](#volume_block) structure is documented below.

* `az_code` - (Required, String, NoneUpdatable) Specifies the AZ code.

* `time_zone` - (Optional, String) Specifies the time zone. The default time zone is **UTC+08:00**.

* `tags_info` - (Required, List, NoneUpdatable) Specifies the tag information.
  The [tags_info](#tags_info_block) structure is documented below.

* `security_group_id` - (Optional, String) Specifies the security group ID.
  By default, the value is the same as the ID of the security group associated with the TaurusDB instance.

* `enable_users_sync` - (Optional, String) Specifies whether users synchronization is enabled.
  Valid values are **true** and **false**.

* `open_slow_log_switch` - (Optional, String) Specifies whether to enable the slow query log original text switch.
  Valid values are **true** and **false**.

* `be_parameter_values` - (Optional, Map) Specifies a map contains mappings of parameter name and value to modify.
  The parameter name must be based on the default parameter template of the backend nodes.

* `fe_parameter_values` - (Optional, Map) Specifies a map contains mappings of parameter name and value to modify.
  The parameter name must be based on the default parameter template of the frontend nodes.

* `charging_mode` - (Optional, String, NoneUpdatable) Specifies the charging mode of the HTAP StarRocks instance.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.

* `period_unit` - (Optional, String, NoneUpdatable) Specifies the charging period unit of the HTAP StarRocks instance.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `period` - (Optional, Int, NoneUpdatable) Specifies the charging period of the HTAP StarRocks instance.
  If `period_unit` is set to **month**, the value ranges from `1` to `9`.
  If `period_unit` is set to **year**, the value ranges from `1` to `3`.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.

* `auto_renew` - (Optional, String) Specifies whether auto-renew is enabled. Valid values are **true** and **false**.

<a name="engine_block"></a>
The `engine` block supports:

* `type` - (Required, String, NoneUpdatable) Specifies the engine type.
  Only **star-rocks** is supported.

* `version` - (Required, String, NoneUpdatable) Specifies the major version number of the engine.

<a name="ha_block"></a>
The `ha` block supports:

* `mode` - (Required, String, NoneUpdatable) Specifies the deployment mode.
  Valid values are **Single** and **Cluster**.

<a name="volume_block"></a>
The `fe_volume` and `be_volume` blocks support:

* `io_type` - (Required, String, NoneUpdatable) Specifies the storage type. The valid values are **SSD** and **ESSD**.
  The value can be obtained from the response to querying HTAP engine resources.

* `capacity_in_gb` - (Required, Int, NoneUpdatable) Specifies the disk capacity in GB.
  The value ranges from `50` to `1,000`. The increment is `10` GB.

<a name="tags_info_block"></a>
The `tags_info` block supports:

* `sys_tags` - (Required, List, NoneUpdatable) Specifies the system tags.
  The [sys_tags](#sys_tags_block) structure is documented below.

<a name="sys_tags_block"></a>
The `sys_tags` block supports:

* `key` - (Required, String, NoneUpdatable) Specifies the tag key.
  The key **_sys_enterprise_project_id** must be set to specify the enterprise project ID.

* `value` - (Required, String, NoneUpdatable) Specifies the tag value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of the StarRocks instance (htap_instance_id).

* `status` - The instance status. Valid values are **creating**, **normal**, **abnormal**, and **createfail**.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `db_user` - The database user. The default value is **root**.

* `port` - The database port number.

* `enable_ssl` - Whether SSL is enabled.

* `data_vip` - The data IP address of the frontend node.

* `fe_node_volume_code` - The storage type of the frontend node.

* `be_node_volume_code` - The storage type of the backend node.

* `fe_node_volume_size` - The storage space of the frontend node.

* `be_node_volume_size` - The storage space of the backend node.

* `pay_model` - The payment method. Valid values are **0** (pay-per-use) and **1** (yearly/monthly).

* `groups` - The instance groups. The [groups](#groups_block) structure is documented below.

* `actions` - The actions performed on the instance. The [actions](#actions_block) structure is documented below.

* `can_enable_public_access` - Whether public access can be enabled.

* `public_ip` - The public IP address.

* `current_backup_node_id` - The current backup node ID.

* `cluster_mode` - The cluster mode.

* `is_frozen` - Whether the instance is frozen.

* `frozen_time` - The frozen time.

* `bak_period` - The backup period.

* `bak_keep_day` - The number of days to keep backups.

* `bak_expected_start_time` - The expected backup start time.

* `bak_store_version_id` - The backup store version ID.

* `bak_store_version` - The backup store version.

* `bak_store_type` - The backup store type.

* `create_at` - The creation time (Unix timestamp).

* `update_at` - The update time (Unix timestamp).

* `delete_at` - The deletion time (Unix timestamp).

* `db_port` - The database port number.

* `param_group` - The parameter group information.

* `create_fail_error_code` - The error code when creation fails.

* `ops_window` - The operation window. The [ops_window](#ops_window_block) structure is documented below.

* `tags_info` - The tag information. The [tags_info](#tags_info_block) structure is documented below.

* `backup_used_space` - The backup used space.

* `port_info` - The port information. The [port_info](#port_info_block) structure is documented below.

* `support_data_replication` - Whether data replication is supported.

* `ssl_option` - Whether SSL is enabled.

* `dedicated_resource_id` - The dedicated resource ID.

* `users_sync_switch_on` - Whether users sync is enabled.

* `data_store_version_id` - The DB version ID.

* `data_store_version` - The DB version.

* `data_store_type` - The DB engine.

* `enterprise_project_id` - The enterprise project ID.

* `new_version_available` - Whether there is a new DB version available.

* `project_id` - The project ID of a tenant in a region.

* `be_configurations` - The parameter configuration information of the backend nodes.
  The [configurations](#htap_starrocks_parameters_configurations_attr) structure is documented below.

* `be_parameters` - The list of parameter objects of the backend nodes.
  The [parameter_values](#htap_starrocks_parameters_parameter_values_attr) structure is documented below.

* `fe_configurations` - The parameter configuration information of the frontend nodes.
  The [configurations](#htap_starrocks_parameters_configurations_attr) structure is documented below.

* `fe_parameters` - The list of parameter objects of the frontend nodes.
  The [parameter_values](#htap_starrocks_parameters_parameter_values_attr) structure is documented below.

<a name="groups_block"></a>
The `groups` block supports:

* `id` - The group ID.

* `name` - The group name.

* `group_type_name` - The instance group type name.

* `group_node_type` - The type of nodes in the instance group. Valid values are **be** and **fe**.

* `nodes` - The instance node information. The [nodes](#nodes_block) structure is documented below.

<a name="nodes_block"></a>
The `nodes` block supports:

* `id` - The instance node ID.

* `name` - The instance node name.

* `type` - The instance node type. Valid values are **be**, **fe-leader**, **fe-follower**, and **fe-observer**.

* `status` - The node status.

* `flavor_id` - The node specification ID.

* `flavor_ref` - The node specification code.

* `iass_flavor_ref` - The IaaS specification code.

* `cpu` - The vCPUs of the instance node.

* `mem` - The memory size (GB) of the instance node.

* `db_port` - The database port number.

* `az_code` - The AZ code.

* `az_description` - The AZ description.

* `az_type` - The AZ type.

* `region_code` - The region where the instance is deployed.

* `period` - The instance node subscription period.

* `volume` - The instance node storage information.
  The [volume](#node_volume_block) structure is documented below.

* `datastore` - The database information.
  The [datastore](#node_datastore_block) structure is documented below.

* `priority` - The node priority.

* `frozen_flag` - The frozen flag.

* `pay_model` - The billing mode.

* `order_id` - The order ID.

* `traffic_ip` - The data IP address.

* `traffic_ipv6` - The data IPv6 address.

* `create_at` - The time when the node was created (Unix timestamp).

* `update_at` - The time when the node was updated (Unix timestamp).

* `max_connections` - The maximum number of public network connections.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `need_restart` - Whether a reboot is required for parameter update.

* `sg_id` - The security group ID.

* `param_group` - The parameter template information.
  The [param_group](#node_param_group_block) structure is documented below.

<a name="node_volume_block"></a>
The `volume` block supports:

* `type` - The storage type of the instance node.

* `size` - The storage space of the instance node.

<a name="node_datastore_block"></a>
The `datastore` block supports:

* `id` - The database ID.

* `type` - The database type.

* `version` - The database version.

<a name="node_param_group_block"></a>
The `param_group` block supports:

* `id` - The parameter template ID.

* `name` - The parameter template name.

<a name="actions_block"></a>
The `actions` block supports:

* `id` - The action ID.

* `action` - The action name.

* `object_id` - The object ID.

* `type` - The action type.

* `job_id` - The job ID.

* `status` - The action status.

* `created_at` - The creation time.

* `updated_at` - The update time.

<a name="ops_window_block"></a>
The `ops_window` block supports:

* `period` - The operation period.

* `start_time` - The start time.

* `end_time` - The end time.

<a name="tags_info_block"></a>
The `tags_info` block supports:

* `tags` - The user-defined tags. The structure is documented below.

* `sys_tags` - The system tags. The structure is documented below.

<a name="port_info_block"></a>
The `port_info` block supports:

* `mysql_port` - The MySQL port number.

<a name="htap_starrocks_parameters_configurations_attr"></a>
The `configurations` block supports:

* `configuration_id` - The parameter template ID.

* `datastore_version_name` - The DB version name.

* `datastore_name` - The DB engine name in the parameter template.

* `created` - The time when the parameter template was created.

* `updated` - The time when the parameter template was updated.

<a name="htap_starrocks_parameters_parameter_values_attr"></a>
The `parameter_values` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `restart_required` - Whether a reboot is required.

* `readonly` - Whether the parameter is read-only.

* `value_range` - The value range of the parameter.

* `type` - The parameter type.

* `description` - The parameter description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

The StarRocks instance can be imported using the `instance_id` and the `htap_instance_id`, separated by a slash, e.g.

```bash
terraform import huaweicloud_taurusdb_htap_starrocks_instance.test <taurusdb_instance_id>/<htap_instance_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `db_root_pwd`, `be_parameter_values`,
`fe_parameter_values`, `period_unit`, `period`, `auto_renew`. It is generally recommended running `terraform plan` after
importing a TaurusDB HTAP StarRocks instance. You can then decide if changes should be applied to the instance,
or the resource definition should be updated to align with the instance. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_taurusdb_htap_starrocks_instance" "test" {
  ...

lifecycle {
  ignore_changes = [
    "db_root_pwd", "be_parameter_values", "fe_parameter_values", "period_unit", "period", "auto_renew"
  ]
}
}
```
