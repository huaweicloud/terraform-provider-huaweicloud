---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_replications"
description: |-
  Use this data source to query the replications information.
---

# huaweicloud_dcs_replications

Use this data source to query the replications information.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}

data "huaweicloud_dcs_replications" "test" {
  instance_id = var.instance_id
  group_id    = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `group_id` - (Required, String) Specifies the ID of the shard.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `replications` - Indicates the list of replication.
  The [replications](#dcs_replications_replications_struct) structure is documented below.

<a name="dcs_replications_replications_struct"></a>
The `replications` block supports:

* `id` - Indicates the ID of the replication.

* `role` - Indicates the role of the replication.

* `replication_ip` - Indicates the IP address of the replication.

* `is_replication` - Whether the replica is newly added.

* `node_id` - Indicates the ID of the node.

* `status` - Indicates the status of the replication.

* `az_code` - Indicates the availability zone where the replication located.

* `dimensions` - Indicates the corresponding monitoring indicator dimension information of the replication.
  The [dimensions](#dcs_replications_dimensions_struct) structure is documented below.

<a name="dcs_replications_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - Indicates the monitoring dimension name.

* `value` - Indicates the dimension value.
