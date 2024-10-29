---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_primary_standby_switch"
description: |-
  Manages a DDS primary standby switch resource within HuaweiCloud.
---

# huaweicloud_dds_primary_standby_switch

Manages a DDS primary standby switch resource within HuaweiCloud.

## Example Usage

### Perform switch for a replica set instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dds_primary_standby_switch" "test" {
  instance_id = var.instance_id
}
```

### Promote standby node to primary for replica set node, shard node or config node

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dds_primary_standby_switch" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `node_id` - (Optional, String, ForceNew) Specifies the ID of replica set node, shard node or config node.
  Changing this creates a new resource.

  -> If `node_id` is not specified, perform a primary/secondary switchover in a replica set instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
