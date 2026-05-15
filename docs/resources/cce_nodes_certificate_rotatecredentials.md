---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_nodes_certificate_rotatecredentials"
description: |-
  Use this resource to rotate certificate of specified nodes in a cluster within HuaweiCloud.
---

# huaweicloud_cce_nodes_certificate_rotatecredentials

Use this resource to rotate certificate of specified nodes in a cluster within HuaweiCloud.

-> This resource is a one-time action resource for rotating certificate of specified nodes. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_id" {}

resource "huaweicloud_cce_nodes_certificate_rotatecredentials" "test" {
  cluster_id  = var.cluster_id
  api_version = "v3"
  kind        = "List"
  node_list {
    node_id = var.node_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to rotate the node credentials.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NoneUpdatable) Specifies the ID of the cluster whose node credentials need to be rotated.

* `api_version` - (Required, String, NoneUpdatable) Specifies the API version.
  The value is fixed and cannot be changed. Value options: **v3**.

* `kind` - (Required, String, NoneUpdatable) Specifies the API type.
  The value is fixed and cannot be changed. Value options: **List**.

* `node_list` - (Required, List, NoneUpdatable) Specifies the list of nodes whose credentials need to be rotated.
  The [node_list](#node_list_attr) structure is documented below.

<a name="node_list_attr"></a>
The `node_list` block supports:

* `node_id` - (Required, String, NoneUpdatable) Specifies the node ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the cluster ID.
