---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_databases"
description: ""
---

# huaweicloud_rds_mysql_databases

Use this data source to get the list of RDS MySQL databases.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_rds_mysql_databases" "test" {
  instance_id   = var.instance_id
  name          = "test"
  character_set = "utf8"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `name` - (Optional, String) Specifies the database name.

* `character_set` - (Optional, String) Specifies the character set used by the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The list of databases.
  The [databases](#RDS_mysql_databases) structure is documented below.

<a name="RDS_mysql_databases"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `character_set` - Indicates the character set used by the database.

* `description` - Indicates the database description.
