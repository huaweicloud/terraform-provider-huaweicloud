---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_primary_standby_switch"
description: |-
  Use this resource to promote a read replica to primary for a TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_primary_standby_switch

Use this resource to promote a read replica to primary for a TaurusDB instance within HuaweiCloud.

-> This resource is a one-time action resource for promoting a read replica to the primary node. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_taurusdb_primary_standby_switch" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the read replica that will be promoted to the primary.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
