---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_schemas"
description: |-
  Use this data source to get the list of RDS PostgreSQL schemas.
---

# huaweicloud_rds_pg_schemas

Use this data source to get the list of RDS PostgreSQL schemas.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_pg_schemas" "test" {
  instance_id = var.instance_id
  db_name     = "test_database_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the PostgreSQL instance ID.

* `db_name` - (Required, String) Specifies the database name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_schemas` - Indicates the list of database schemas.

  The [database_schemas](#database_schemas_struct) structure is documented below.

<a name="database_schemas_struct"></a>
The `database_schemas` block supports:

* `schema_name` - Indicates the schema name.

* `owner` - Indicates the schema owner.
