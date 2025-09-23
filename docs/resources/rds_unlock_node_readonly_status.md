---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_unlock_node_readonly_status"
description: |-
  Manages an RDS unlock node readonly status resource within HuaweiCloud.
---

# huaweicloud_rds_unlock_node_readonly_status

Manages an RDS unlock node readonly status resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_unlock_node_readonly_status" "test" {
  instance_id              = var.instance_id
  status_preservation_time = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS SQLServer instance.

* `status_preservation_time` - (Required, Int, NonUpdatable) Specifies the duration (in minutes) during which the HA
  component no longer sets the instance to read-only. Value ranges: **0** to **1440**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
