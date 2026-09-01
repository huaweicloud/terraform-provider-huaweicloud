---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_offline_key_analysis_nodes"
description: |-
  Use this data source to query the offline key analyses.
---

# huaweicloud_dcs_offline_key_analysis_nodes

Use this data source to query the offline key analyses

## Example Usage

```hcl
variable "instance_id" {}
variable "task_id" {}

data "huaweicloud_dcs_offline_key_analysis_nodes" "test" {
  instance_id = var.instance_id
  task_id     = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `task_id` - (Required, String) Specifies the ID of the task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The list of offline key analysis nodes.
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - The node ID.

* `name` - The node name.

* `group_name` - The group name.

* `node_ipv6` - The node IP address.
  Firstly, return the IPv6 address. If the node does not have an IPv6 configuration, then return the IPv4 address.
