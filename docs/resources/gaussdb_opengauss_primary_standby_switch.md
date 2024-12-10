---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_primary_standby_switch"
description: |-
  Manages a GaussDB OpenGauss primary standby switch resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_primary_standby_switch

Manages a GaussDB OpenGauss primary standby switch resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

resource "huaweicloud_gaussdb_opengauss_primary_standby_switch" "test" {
  instance_id = var.instance_id

  shards {
    node_id      = var.node_id
    component_id = var.component_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance.

  Changing this parameter will create a new resource.

* `shards` - (Required, List, ForceNew) Specifies the nodes. You can switch standby DNs of multiple shards to primary DNs.
  The node information is the node IDs and component IDs of shards whose standby DNs are promoted to primary.
  The [shards](#shards_struct) structure is documented below.

  Changing this parameter will create a new resource.

<a name="shards_struct"></a>
The `shards` block supports:

* `node_id` - (Required, String, ForceNew) Specifies the ID of the node where the standby DN to be promoted to primary is
  deployed.

  Changing this parameter will create a new resource.

* `component_id` - (Required, String, ForceNew) Specifies the ID of the standby DN to be promoted to primary.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
