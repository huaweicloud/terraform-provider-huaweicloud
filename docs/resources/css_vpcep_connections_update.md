---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_vpcep_connections_update"
description: |-
  Use this resource to update the VPCEP service connections of a CSS cluster.
---

# huaweicloud_css_vpcep_connections_update

Use this resource to update the VPCEP service connections of a CSS cluster.

-> This resource is only a one-time action resource for updating the VPCEP service connections of a CSS cluster.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "action" {}
variable "endpoint_ids" {
  type = list(string)
}

resource "huaweicloud_css_vpcep_connections_update" "test" {
  cluster_id       = var.cluster_id
  action           = var.action
  endpoint_id_list = var.endpoint_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the CSS cluster.

* `action` - (Required, String, NonUpdatable) Specifies the expected behavior.
  The valid values are as follows:
  + **receive**: Accept the VPC endpoint.
  + **reject**: Reject the VPC endpoint.

* `endpoint_id_list` - (Required, List, NonUpdatable) Specifies the list of VPC endpoint IDs.
  The value is a list of strings.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
