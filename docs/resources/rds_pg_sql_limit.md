---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_sql_limit"
description: |-
  Manage an RDS PostgreSQL SQL limit resource within HuaweiCloud.
---

# huaweicloud_rds_pg_sql_limit

Manage an RDS PostgreSQL SQL limit resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}

resource "huaweicloud_rds_pg_sql_limit" "instance" {
  instance_id     = var.instance_id
  db_name         = var.db_name
  query_id        = "5"
  max_concurrency = 20
  max_waiting     = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds PostgreSQL SQL limit resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS PostgreSQL instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the name of the database.

* `max_concurrency` - (Required, Int) Specifies the number of SQL statements executed simultaneously.
  Value ranges from `0` to `50000`. `0` means no limit.

* `max_waiting` - (Required, Int) Specifies the max waiting time in seconds.

* `query_id` - (Optional, String, NonUpdatable) Specifies the query ID.
  Value ranges: **-9223372036854775808~9223372036854775807**.

* `query_string` - (Optional, String, NonUpdatable) Specifies the text form of SQL statement.

  -> **NOTE:** Exactly one of `query_id`, `query_string` should be specified.

* `search_path` - (Optional, String, NonUpdatable) Specifies the query order for names that are not schema qualified.
  Defaults to **public**,

* `switch` - (Optional, String) Specifies the SQL limit switch. Value options: **open**, **close**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID. The format is `<instance_id>/<db_name>/<sql_limit_id>`.

* `sql_limit_id` - Indicates the ID of SQL limit.

* `is_effective` - Indicates whether the SQL limit is effective.

## Import

The SQL limit can be imported using the `instance_id`, `db_name` and `sql_limit_id`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_rds_pg_sql_limit.test <instance_id>/<db_name>/<sql_limit_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `query_id`, `query_string`. It is generally
recommended running `terraform plan` after importing an RDS PostgreSQL SQL limit. You can then decide if changes should
be applied to the RDS PostgreSQL SQL limit, or the resource definition should be updated to align with the RDS PostgreSQL
SQL limit. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_rds_pg_sql_limit" "test" {
  ...

  lifecycle {
    ignore_changes = [
      "query_id", "query_string",
    ]
  }
}
```
