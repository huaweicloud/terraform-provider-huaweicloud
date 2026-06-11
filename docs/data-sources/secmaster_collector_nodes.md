---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_nodes"
description: |-
  Use this data source to get the list of collector nodes.
---

# huaweicloud_secmaster_collector_nodes

Use this data source to get the list of collector nodes.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_nodes" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `health_status` - (Optional, String) Specifies the health status of the node.  
  The valid values are as follows:
  + **NORMAL**
  + **ANOMALIES**
  + **FAULTS**
  + **LOST_CONTACT**

* `node_id` - (Optional, String) Specifies the node ID.

* `node_name` - (Optional, String) Specifies the node name.

* `ip_address` - (Optional, String) Specifies the IP address.

* `sort_key` - (Optional, String) Specifies the sort field.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The value can be **asc** or **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The result data list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `channel_instance_refer_count` - The channel instance refer count.

* `create_by` - The IAM user ID.

* `custom_label` - The custom label.

* `description` - The description.

* `device_type` - The device type.

* `ip_address` - The IP address.

* `monitor` - The monitor information.

  The [monitor](#monitor_struct) structure is documented below.

* `node_expansion` - The node expansion information.

  The [node_expansion](#node_expansion_struct) structure is documented below.

* `node_id` - The node ID.

* `node_name` - The node name.

* `os_type` - The OS type.

* `private_ip_address` - The private IP address.

* `project_id` - The project ID.

* `region` - The region.

* `specification` - The specification.

* `update_time` - The update time, in milliseconds.

* `vpc_endpoint_address` - The VPC endpoint address.

* `vpc_endpoint_id` - The VPC endpoint ID.

* `vpc_id` - The VPC ID.

* `workspace_id` - The workspace ID.

<a name="monitor_struct"></a>
The `monitor` block supports:

* `cpu_idle` - The CPU idle percentage.

* `cpu_usage` - The CPU usage.

* `disk_count` - The disk count.

* `disk_usage` - The disk usage.

* `down_pps` - The download packets per second.

* `health_status` - The health status.

* `heart_beat` - The heartbeat status.

* `heart_beat_time` - The last heartbeat time.

* `memory_cache` - The memory cache size.

* `memory_count` - The memory count.

* `memory_free` - The free memory size.

* `memory_shared` - The shared memory size.

* `memory_usage` - The memory usage.

* `mini_on_online` - Whether the node is online.

* `read_rate` - The disk read rate.

* `up_pps` - The upload packets per second.

* `write_rate` - The disk write rate.

<a name="node_expansion_struct"></a>
The `node_expansion` block supports:

* `custom_label` - The custom label.

* `data_center` - The data center.

* `description` - The description.

* `maintainer` - The maintainer.

* `network_plane` - The network plane.

* `node_id` - The node ID.
