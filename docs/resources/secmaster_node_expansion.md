---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_node_expansion"
description: |-
  Manages a node expansion resource within HuaweiCloud.
---

# huaweicloud_secmaster_node_expansion

Manages a node expansion resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "node_id" {}

resource "huaweicloud_secmaster_node_expansion" "test" {
  workspace_id   = var.workspace_id
  node_id        = var.node_id
  custom_label   = "test custom label"
  data_center    = "test data center"
  description    = "test description"
  maintainer     = "admin"
  network_plane  = "business"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the node belongs.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the node.

* `custom_label` - (Optional, String) Specifies the custom label of the node.

* `data_center` - (Optional, String) Specifies the data center of the node.

* `description` - (Optional, String) Specifies the description of the node.

* `maintainer` - (Optional, String) Specifies the maintainer of the node.

* `network_plane` - (Optional, String) Specifies the network plane of the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the node ID.

## Import

The node expansion can be imported using the `workspace_id` and the `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_node_expansion.test <workspace_id>/<id>
```
