---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_object_mappings"
description: |-
  Use this data source to get object mappings of a DRS job.
---

# huaweicloud_drs_object_mappings

Use this data source to get object mappings of a DRS job.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_object_mappings" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `db_name` - (Optional, String) Specifies the database name used to filter the object mappings.

* `schema_name` - (Optional, String) Specifies the schema name used to filter the object mappings.

* `table_name` - (Optional, String) Specifies the table name used to filter the object mappings.

* `has_column_info` - (Optional, String) Specifies whether to query the column mapping information.
  Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `object_mapping_list` - The list of object mappings.

  The [object_mapping_list](#object_mapping_list_struct) structure is documented below.

<a name="object_mapping_list_struct"></a>
The `object_mapping_list` block supports:

* `source_db_name` - The source database name.

* `source_schema_name` - The source schema name.

* `source_table_name` - The source table name.

* `target_db_name` - The target database name.

* `target_schema_name` - The target schema name.

* `target_table_name` - The target table name.

* `has_column_info` - Whether the object mapping includes column information.
