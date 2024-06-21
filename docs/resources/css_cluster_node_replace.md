---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_node_replace"
description: |-
  Manages CSS cluster node replace resource within HuaweiCloud
---

# huaweicloud_css_cluster_node_replace

Manages CSS cluster node replace resource within HuaweiCloud

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_id" {}

resource "huaweicloud_css_cluster_node_replace" "test" {
  cluster_id   = var.cluster_id
  node_id      = var.node_id
  agency       = "ess_replace_agency"
  migrate_data = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS cluster.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS cluster node.

* `agency` - (Optional, String, NonUpdatable) Specifies the IAM agency used to access CSS.

* `migrate_data` - (Optional, Bool, NonUpdatable) Specifies whether to migrate data. Defaults to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
