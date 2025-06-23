---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_read_replica_restorable_databases"
description: |-
  Use this data source to query the databases, schemas, and tables that can be restored from a read replica in 
  HuaweiCloud RDS for PostgreSQL.
---

# huaweicloud_rds_read_replica_restorable_databases

Use this data source to query the databases, schemas, and tables that can be restored from a read replica in
HuaweiCloud RDS for PostgreSQL.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_read_replica_restorable_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_limit` - Indicates the maximum number of databases that can be returned in a single response.

* `total_tables` - Indicates the total number of tables returned in the response.

* `table_limit` - Indicates the maximum number of tables that can be returned in a single response.

* `databases` - Indicates the list of databases that can be restored to the primary instances

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the name of the database.

* `total_tables` - Indicates the total number of tables in the database.

* `schemas` - Indicates the list of schemas in the database.

  The [schemas](#schemas_struct) structure is documented below.

<a name="schemas_struct"></a>
The `schemas` block supports:

* `name` - Indicates the name of the schema.

* `total_tables` - Indicates the total number of tables in the schema.

* `tables` - Indicates the list of tables in the schema.

  The [tables](#tables_struct) structure is documented below.

<a name="tables_struct"></a>
The `tables` block supports:

* `name` - Indicates the name of the table.
