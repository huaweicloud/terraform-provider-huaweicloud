---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_authorized_databases"
description: |-
  Use this data source to query the authorized databases of a database user in a specified RDS instance.
---

# huaweicloud_rds_mysql_authorized_databases

Use this data source to query the authorized databases of a database user in a specified RDS instance.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_name" {}

data "huaweicloud_rds_mysql_authorized_databases" "test" {
  instance_id = var.instance_id
  user_name   = var.user_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `user_name` - (Required, String) Specifies the user name of the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the list of authorized database objects.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block contains:

* `name` - Indicates the name of the authorized database.

* `readonly` - Indicates whether the permission is read-only. Value can be:
  + **true**: Indicates the permission is read-only.
  + **false**: Indicates the permission is read/write.
