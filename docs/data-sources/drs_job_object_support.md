---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_object_support"
description: |-
  Use this data source to get object selection and column mapping support information for a specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_job_object_support

Use this data source to get object selection and column mapping support information for a specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_job_object_support" "test" { 
  job_id = var.job_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_full_trans_support_object` - The indication of whether full migration tasks support object selection.

* `is_incre_trans_support_object` - The indication of whether incremental migration tasks support object selection.

* `is_full_incre_trans_support_object` - The indication of whether full plus incremental migration tasks support
  object selection.

* `support_object_import_engine` - The list of engines that support object import.

* `is_support_column_mapping` - The indication of whether column mapping is supported.

* `is_database_support_search` - The indication of whether database-level search is supported.

* `is_schema_support_search` - The indication of whether schema-level search is supported.

* `is_table_support_search` - The indication of whether table-level search is supported.

* `file_size` - The maximum file size supported for object import, in MB.

* `previous_select` - The method used to select migration or synchronization objects in the previous operation.
  If empty, no method has been selected yet. Valid values include:
  + **srcImportObject**: The previous selection method was import.

* `import_level` - The level of object import. Valid values are:
  + **table**: Table level.
  + **database**: Database level.

* `is_import_cloumn` - The indication of whether the column processing method was set to import in the previous
  operation.
  + **true**: The previous column processing method was import.
  + **false** or empty: The previous column processing method was manual selection.

* `import_mapping_type` - The file import mapping scenario. Valid values are:
  + **table_mapping**
  + **topic_mapping**

* `is_import_unique_key` - The indication of whether unique key information is imported.
