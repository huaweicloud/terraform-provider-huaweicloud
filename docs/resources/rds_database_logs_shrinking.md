---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_database_logs_shrinking"
description: |-
  Manages an RDS database logs shrinking resource within HuaweiCloud.
---

# huaweicloud_rds_database_logs_shrinking

Manages an RDS database logs shrinking resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}

resource "huaweicloud_rds_database_logs_shrinking" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of instance.

  Changing this creates a new resource.

* `db_name` - (Required, String, ForceNew) Specifies the name of the database.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the instance ID.
