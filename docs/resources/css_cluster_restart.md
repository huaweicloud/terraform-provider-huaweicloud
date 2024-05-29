---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_restart"
description: |-
  Manages CSS cluster restart resource within HuaweiCloud.
---

# huaweicloud_css_cluster_restart

Manages CSS cluster restart resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_css_cluster_restart" "test" {
  cluster_id = var.cluster_id
  type       = "role"
  value      = "ess"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the CSS cluster.
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the operation type of the CSS cluster restart.
  The value can be **role** or **node**. The value can only be **role** when the `is_rolling` is **true**.
  Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the value under the operation type. If the operation
  role is node, the value is the node ID. If the operation role is role, the value is one or multiple node
  types (such as **ess**, **ess-master**, **ess-client**, **ess-cold**, and **all**). Use commas (,) to
  separate multiple node types.
  Changing this creates a new resource.

* `is_rolling` - (Optional, Bool, ForceNew) Specifies whether to roll restart.
  Changing this creates a new resource.

  -> **NOTE:** Rolling restart is only supported when the number of nodes in the cluster (including Master nodes,
  Client nodes, and cold data nodes) is greater than 3.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
