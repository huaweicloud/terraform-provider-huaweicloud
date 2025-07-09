---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_engine_versions"
description: |-
  Use this data source to obtain all version information of the specified engine type of RDS.
---

# huaweicloud_rds_engine_versions

Use this data source to obtain all version information of the specified engine type of RDS.

## Example Usage

```hcl
data "huaweicloud_rds_engine_versions" "test" {
  type = "SQLServer"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the RDS engine versions.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the RDS engine type.
  The valid values are **MySQL**, **PostgreSQL**, **SQLServer** and **MariaDB**, default to **MySQL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the list of RDS versions.
  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `id` - Indicates the version ID.

* `name` - Indicates the version name.
