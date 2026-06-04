---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_all_instances_metrics"
description: |-
  Use this data source to query the metrics of all GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_all_instances_metrics

Use this data source to query the metrics of all GaussDB instances within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_gaussdb_all_instances_metrics" "test" {
}
```

### Filter by Instance ID

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_all_instances_metrics" "test" {
  instance_id = var.instance_id
}
```

### Filter by Instance Name

```hcl
variable "instance_name" {}

data "huaweicloud_gaussdb_all_instances_metrics" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instances metrics.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the instance name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instance metrics.
  The [instances](#gaussdb_all_instances_metrics_instances) structure is documented below.

<a name="gaussdb_all_instances_metrics_instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `name` - The instance name.

* `status` - The instance status. The valid values are:
  + **creating**: The instance is being created.
  + **normal**: The instance is normal.
  + **abnormal**: The instance is abnormal.
  + **createfail**: The instance creation failed.

* `mode` - The instance type.

* `engine_name` - The engine name.

* `engine_version` - The engine version.

* `solution` - The deployment mode.

* `disk_used_size` - The used size of the instance data disk.

* `disk_total_size` - The total size of the instance data disk.

* `disk_usage` - The usage percentage of the instance data disk.

* `p80` - The 80th percentile SQL response time.

* `p95` - The 95th percentile SQL response time.

* `deadlocks` - The number of deadlocks.

* `buffer_hit_ratio` - The buffer hit ratio.

* `nodes` - The list of instance node information.
  The [nodes](#gaussdb_all_instances_metrics_nodes) structure is documented below.

<a name="gaussdb_all_instances_metrics_nodes"></a>
The `nodes` block supports:

* `id` - The node ID.

* `name` - The node name.

* `role` - The node role. The valid values are:
  + **master**: Primary node.
  + **slave**: Standby node.
  + **secondary**: Log node.
  + **readreplica**: Read replica.

* `component_ids` - The list of component IDs.
