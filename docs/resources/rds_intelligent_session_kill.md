---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_intelligent_session_kill"
description: |-
  Manage an RDS intelligent session kill resource within HuaweiCloud.
---

# huaweicloud_rds_intelligent_session_kill

Manage an RDS intelligent session kill resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_intelligent_session_kill" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS instance.

* `auto_add_sql_limit_rule` - (Optional, String, NonUpdatable) Specifies whether to enable automatic SQL throttling.
  Value options:
  + **true**: Enable automatic SQL throttling.
  + **false**: Disable automatic SQL throttling.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID which is same as `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
