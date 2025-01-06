---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node_sync"
description: |-
  Use this resource to sync the CCE node within HuaweiCloud.
---

# huaweicloud_cce_node_sync

Use this resource to sync the CCE node within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "node_id" {}

resource "huaweicloud_cce_node_sync" "test" {
  cluster_id = var.cluster_id
  node_id    = var.node_id
}
```

~> Deleting node sync is not supported, it will only be removed from the state.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the node sync resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `node_id`.
