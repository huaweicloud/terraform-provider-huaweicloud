---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_extra_column_info"
description: |-
  Use this data source to query the extra column info of data processing objects of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_extra_column_info

Use this data source to query the extra column info of data processing objects of a DRS job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_extra_column_info" "test" {
  job_id = var.job_id
}
```

### Filter by Sent Status

```hcl
variable "job_id" {}

data "huaweicloud_drs_extra_column_info" "test" {
  job_id            = var.job_id
  is_only_show_sent = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the extra column info.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job to query extra column info.

* `is_only_show_sent` - (Optional, Bool) Specifies whether to query only the processed objects that have been sent.

* `fetch_all` - (Optional, Bool) Specifies whether to query all data.  
  When set to **true**, offset and limit parameters are ignored.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `column_process_objects` - The list of column process objects.  
  The [column_process_objects](#drs_extra_column_info_column_process_objects) structure is documented below.

<a name="drs_extra_column_info_column_process_objects"></a>
The `column_process_objects` block supports:

* `object_alias_name` - The alias name of the object after mapping.

* `object_source_names` - The source object names of the selected source database.

* `is_sent` - Whether the extra column has been sent.

* `extra_column_infos` - The extra column information.  
  The [extra_column_infos](#drs_extra_column_info_extra_column_infos) structure is documented below.

<a name="drs_extra_column_info_extra_column_infos"></a>
The `extra_column_infos` block supports:

* `column_name` - The name of the column.

* `column_type` - The type of the column.

* `column_value` - The value of the column.

* `data_type` - The data type of the column.
