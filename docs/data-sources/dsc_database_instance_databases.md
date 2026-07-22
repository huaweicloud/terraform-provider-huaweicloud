---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_database_instance_databases"
description: |-
  Use this data source to get the list of databases under a database instance within HuaweiCloud.
---

# huaweicloud_dsc_database_instance_databases

Use this data source to get the list of databases under a database instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dsc_database_instance_databases" "test" {
  instance_id = var.instance_id
  type        = "MySQL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the database instance ID.

* `type` - (Required, String) Specifies the database type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `databases` - The database information list.
  The [databases](#dsc_databases_struct) structure is documented below.

<a name="dsc_databases_struct"></a>
The `databases` block supports:

* `id` - The database ID.

* `db_port` - The database port.

* `db_name` - The database name.

* `asset_name` - The asset name.

* `authorized` - Whether the database is authorized.

* `default` - Whether the database is the default database.
