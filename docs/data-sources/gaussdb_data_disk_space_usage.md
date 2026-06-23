---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_data_disk_space_usage"
description: |-
  Use this data source to query the data disk space usage of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_data_disk_space_usage

Use this data source to query the data disk space usage of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_data_disk_space_usage" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data disk space usage.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_disk_capacity` - The total capacity of the data disk, in GB.

* `data_disk_usage` - The used capacity of the data disk, in GB.

* `space_usage_growth_per_day` - The average daily growth of space usage, in GB.

* `estimated_remaining_days` - The estimated number of remaining days before the disk is full.

* `cn_components` - The list of CN node components.
  The [cn_components](#gaussdb_data_disk_space_usage_cn_components) structure is documented below.

* `dn_components` - The list of DN node components.
  The [dn_components](#gaussdb_data_disk_space_usage_dn_components) structure is documented below.

<a name="gaussdb_data_disk_space_usage_cn_components"></a>
The `cn_components` block supports:

* `node_id` - The ID of the node.

* `component_id` - The ID of the component.

* `node_name` - The name of the node.

<a name="gaussdb_data_disk_space_usage_dn_components"></a>
The `dn_components` block supports:

* `node_id` - The ID of the node.

* `component_id` - The ID of the component.

* `role` - The role of the component (e.g., slave).

* `node_name` - The name of the node.
