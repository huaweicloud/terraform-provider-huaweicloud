---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_engine_versions"
description: |-
  Use this data source to get the database specifications of a specified DB engine.
---

# huaweicloud_gaussdb_mysql_engine_versions

Use this data source to get the database specifications of a specified DB engine.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_engine_versions" "test" {
  database_name = "gaussdb-mysql"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `database_name` - (Required, String) Specifies the DB engine.
  Value options: **gaussdb-mysql**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datastores` - Indicates the DB version list.

  The [datastores](#datastores_struct) structure is documented below.

<a name="datastores_struct"></a>
The `datastores` block supports:

* `id` - Indicates the DB version ID.

* `name` - Indicates the DB version number.
  Only the major version number with two digits is returned.

* `version` - Indicates the compatible open-source DB version.
  A three-digit open-source version is returned.

* `kernel_version` - Indicates the DB version.
  A complete four-digit version is returned.
