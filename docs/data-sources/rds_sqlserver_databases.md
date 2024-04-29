---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sqlserver_databases"
description: ""
---

# huaweicloud_rds_sqlserver_databases

Use this data source to get the list of RDS SQLServer databases.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_rds_sqlserver_databases" "test" {
  instance_id = var.instance_id
  name        = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS SQLServer instance.

* `name` - (Optional, String) Specifies the database name.

* `character_set` - (Optional, String) Specifies the character set used by the database.

* `state` - (Optional, String) Specifies the database status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - The list of databases.
  The [databases](#RDS_sqlserver_databases) structure is documented below.

<a name="RDS_sqlserver_databases"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `character_set` - Indicates the character set used by the database.

* `state` - Indicates the database status.
