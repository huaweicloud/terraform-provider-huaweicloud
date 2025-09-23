---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_nodes_remove"
description: |-
  Use this resource to remove nodes from a CCE cluster within HuaweiCloud.
---

# huaweicloud_cce_nodes_remove

Use this resource to remove nodes from a CCE cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "node_id_1" {}
variable "node_id_2" {}

resource "huaweicloud_cce_nodes_remove" "test" {
  cluster_id  = var.cluster_id

  nodes {
    uid = node_id_1
  }
  nodes {
    uid = node_id_2
  }
}

```

~> Deleting nodes remove resource is not supported, it will only be removed from the state.  

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE pool nodes add resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `nodes` - (Required, List, NonUpdatable) Specifies the list of nodes to remove form the cluster.
  The [nodes](#nodes) structure is documented below.

<a name="nodes"></a>
The `nodes` block supports:

* `uid` - (Required, String, NonUpdatable) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
