---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_available_upgrade_small_versions"
description: |-
  Use this data source to query the RDS available upgrade small versions (PostgreSQL) within HuaweiCloud.
---

# huaweicloud_rds_available_upgrade_small_versions

Use this data source to query the RDS available upgrade small versions (PostgreSQL) within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_rds_available_upgrade_small_versions" "test" {
  database_name = "postgresql"
  version       = "16"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the small versions.
  If omitted, the provider-level region will be used.

* `database_name` - (Required, String) Specifies the database engine name.
  The valid value is **postgresql**.

* `version` - (Required, String) Specifies the database version to query small versions for.
  The value can be a major version number or a minor version number, for example: **16**, **16.5**, **16.5.21**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_stores` - The list of available DB engine versions to which an instance can be upgraded
  The [data_stores](#data_stores_attr) structure is documented below.

<a name="data_stores_attr"></a>
The `data_stores` block supports:

* `id` - The database version ID.

* `name` - The available small version number.

* `favored` - Whether this is a recommended version.
