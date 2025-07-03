---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_primary_standby_switch"
description: |-
  Manages an RDS instance primary standby switch resource within HuaweiCloud.
---

# huaweicloud_rds_primary_standby_switch

Manages an RDS instance primary standby switch resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_primary_standby_switch" "test" {
  instance_id = var.target_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of instance.

* `force` - (Optional, Bool, NonUpdatable) Specifies whether to perform a forcible primary/standby switchover.
  This parameter is valid only for the PostgreSQL DB engine.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
