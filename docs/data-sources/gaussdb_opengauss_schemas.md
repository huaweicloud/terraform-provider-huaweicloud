---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_schemas"
description: |-
  Use this data source to get the database schemas of a specified GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_schemas

Use this data source to get the database schemas of a specified GaussDB OpenGauss instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_schemas" "this" {
  instance_id = var.instance_id
  db_name     = "test_db_name"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instance. If omitted, the provider-level region will
  be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

* `db_name` - (Required, String) Specifies the database name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `database_schemas` - Indicates the list of the database schemas.

  The [database_schemas](#database_schemas_struct) structure is documented below.

<a name="database_schemas_struct"></a>
The `database_schemas` block supports:

* `schema_name` - Indicates the schema name.

* `owner` - Indicates the owner of the schema.
