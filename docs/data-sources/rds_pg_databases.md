---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_databases"
description: ""
---

# huaweicloud_rds_pg_databases

Use this data source to get the list of RDS PostgreSQL databases.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_pg_databases" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the PostgreSQL instance ID.

* `name` - (Optional, String) Specifies the database name.

* `owner` - (Optional, String) Specifies the database owner.

* `character_set` - (Optional, String) Specifies the character set used by the database.
  For details, see [documentation](https://www.postgresql.org/docs/16/infoschema-character-sets.html).

* `lc_collate` - (Optional, String) Specifies the database collation.
  For details, see [documentation](https://support.huaweicloud.com/intl/en-us/bestpractice-rds/rds_pg_0017.html).

* `size` - (Optional, Int) Specifies the database size, in bytes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the database list.
  The [databases](#Pg_Databases) structure is documented below.

<a name="Pg_Databases"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `owner` - Indicates the database owner.

* `character_set` - Indicates the character set used by the database.

* `lc_collate` - Indicates the database collation.

* `size` - Indicates the database size, in bytes.

* `description` - Indicates the database description.
