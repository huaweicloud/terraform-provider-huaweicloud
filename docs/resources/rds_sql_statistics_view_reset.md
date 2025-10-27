---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sql_statistics_view_reset"
description: |-
  Manages an RDS SQL statistics view reset resource within HuaweiCloud.
---

# huaweicloud_rds_sql_statistics_view_reset

Manages an RDS SQL statistics view reset resource within HuaweiCloud.

-> **NOTE:** Deleting RDS SQL statistics view reset is not supported. If you destroy a resource of SQL statistics view
reset, the resource is only removed from the state, but it remains in the cloud. And the instance doesn't return to the
state before upgrade.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_sql_statistics_view_reset" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
