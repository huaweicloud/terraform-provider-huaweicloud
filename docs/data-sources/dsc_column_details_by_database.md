---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_column_details_by_database"
description: |-
  Use this data source to get the column details of sensitive data by database dimension within HuaweiCloud.
---

# huaweicloud_dsc_column_details_by_database

Use this data source to get the column details of sensitive data by database dimension within HuaweiCloud.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_dsc_column_details_by_database" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the group label ID for filtering.
  Either `label_id` or `type_id` must be specified.

* `type_id` - (Optional, String) Specifies the type ID for filtering.
  Either `label_id` or `type_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `results` - The column details list by database dimension.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `asset_name` - The asset name.

* `count` - The match count.

* `db_name` - The database name.

* `db_type` - The database type.

* `tables` - The table information list.

  The [tables](#results_tables_struct) structure is documented below.

<a name="results_tables_struct"></a>
The `tables` block supports:

* `count` - The match count.

* `table_name` - The table name.

* `columns` - The column information list.

  The [columns](#results_tables_columns_struct) structure is documented below.

<a name="results_tables_columns_struct"></a>
The `columns` block supports:

* `column_name` - The column name.

* `classification_tags` - The classification tags.

* `level_tags` - The level tags.
