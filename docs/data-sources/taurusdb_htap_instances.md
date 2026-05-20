---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_instances"
description: |-
  Use this data source to query the list of TaurusDB HTAP instances within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_instances

Use this data source to query the list of TaurusDB HTAP instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_htap_instances" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP instances.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the TaurusDB instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of HTAP instances.
  The [instances](#taurusdb_htap_instances_attr) structure is documented below.

* `max_htap_instance_num_of_taurus` - The maximum number of HTAP instances for TaurusDB.

<a name="taurusdb_htap_instances_attr"></a>
The `instances` block supports:

* `id` - The ID of the HTAP instance.

* `name` - The name of the HTAP instance.

* `engine_name` - The HTAP DB engine name.

* `engine_version` - The HTAP DB engine version.

* `project_id` - The project ID of a tenant in a region.

* `instance_state` - The instance state information.
  The [instance_state](#taurusdb_htap_instances_instance_state) structure is documented below.

* `create_at` - The creation time of the HTAP instance.

* `is_frozen` - Whether the HTAP instance is frozen.

* `ha_mode` - The deployment mode of the HTAP instance.

* `pay_model` - The billing mode.
  The valid values are as follows:
  + **0**: pay-per-use
  + **1**: yearly/monthly

* `order_id` - The order ID for the yearly/monthly subscription.

* `alter_order_id` - The alternative order ID for the yearly/monthly subscription.

* `data_vip` - The private IP address.

* `readable_node_infos` - The readable node information.
  The [readable_node_infos](#taurusdb_htap_instances_readable_node_infos) structure is documented below.

* `proxy_ips` - The proxy IP addresses.

* `data_vip_v6` - The private IPv6 address.

* `port` - The database port.

* `available_zones` - The availability zone information.
  The [available_zones](#taurusdb_htap_instances_available_zones) structure is documented below.

* `current_actions` - The instance actions.
  The [current_actions](#taurusdb_htap_instances_current_actions) structure is documented below.

* `volume_type` - The storage type.
  The valid values are as follows:
  + **SSD**: ultra-high I/O
  + **ESSD**: extreme SSD

* `server_type` - The server type.

* `enterprise_project_id` - The enterprise project ID.

* `dedicated_resource_id` - The dedicated resource pool ID.

* `network` - The network information.
  The [network](#taurusdb_htap_instances_network) structure is documented below.

* `ch_master_node_id` - The ClickHouse primary node ID.

* `node_num` - The number of nodes.

<a name="taurusdb_htap_instances_instance_state"></a>
The `instance_state` block supports:

* `instance_status` - The status of the HTAP instance.
  The valid values are as follows:
  + **creating**: An instance is being created.
  + **normal**: An instance is available.
  + **abnormal**: An instance is abnormal.
  + **createfail**: An instance failed to be created.

* `create_fail_error_code` - The error code for HTAP instance creation failures.

* `fail_message` - The error message for HTAP instance creation failures.

* `wait_restart_for_params` - Whether a reboot is required for parameter updates.

<a name="taurusdb_htap_instances_readable_node_infos"></a>
The `readable_node_infos` block supports:

* `data_ip` - The IP address of the readable node.

* `node_id` - The ID of the readable node.

* `node_name` - The name of the readable node.

<a name="taurusdb_htap_instances_available_zones"></a>
The `available_zones` block supports:

* `code` - The AZ code.

* `description` - The AZ description.

* `az_type` - The AZ type.

<a name="taurusdb_htap_instances_current_actions"></a>
The `current_actions` block supports:

* `id` - The instance or node action ID.

* `action` - The instance or node action name.

* `object_id` - The object ID of the instance or node action.

* `type` - The instance or node action type.

* `job_id` - The task ID of the instance or node action.

* `status` - The instance or node action status.

* `created_at` - The creation time of the instance or node action.

* `updated_at` - The update time of the instance or node action.

<a name="taurusdb_htap_instances_network"></a>
The `network` block supports:

* `vpc_id` - The VPC ID.

* `sub_net_id` - The subnet ID.

* `security_group_id` - The security group ID.
