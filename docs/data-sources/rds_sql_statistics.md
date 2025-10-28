---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sql_statistics"
description: |-
  Use this data source to get the list of SQL statement statistics of a DB instance.
---

# huaweicloud_rds_sql_statistics

Use this data source to get the list of SQL statement statistics of a DB instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_sql_statistics" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `list` - Indicates the list of SQL statistics.

  The [list](#list_struct) structure is documented below.

<a name="list_struct"></a>
The `list` block supports:

* `query` - Indicates the text format of an SQL statement.

* `rows` - Indicates the scanned rows.

* `can_use` - Indicates whether SQL throttling can be applied.

* `user_name` - Indicates the username.

* `database` - Indicates the database name.

* `query_id` - Indicates the internal hash code calculated by the SQL parse tree.

* `calls` - Indicates the number of calls.
