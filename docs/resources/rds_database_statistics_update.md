---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_database_statistics_update"
description: |-
  Manage an RDS database statistics update resource within HuaweiCloud.
---

# huaweicloud_rds_database_statistics_update

Manage an RDS database statistics update resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}

resource "huaweicloud_rds_database_statistics_update" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the name of the database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
