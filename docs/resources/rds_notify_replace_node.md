---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_notify_replace_node"
description: |-
  Manages an RDS notify replace node resource within HuaweiCloud.
---

# huaweicloud_rds_notify_replace_node

Manages an RDS notify replace node resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_rds_notify_replace_node" "test" {
  instance_id    = var.instance_id
  node_id        = var.node_id
  replace_action = "REPLACE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS read replica instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the node ID of the RDS read replica instance.

* `replace_action` - (Required, String, NonUpdatable) Specifies the replacement action. Value options:
  + **REPLACE**: node replace
  + **REPLACE_ROLLBACK**: node replace rollback

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<instance_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
