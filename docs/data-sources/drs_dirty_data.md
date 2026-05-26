---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_dirty_data"
description: |-
  Use this data source to get the dirty data list of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_dirty_data

Use this data source to get the dirty data list of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_dirty_data" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

* `begin_time` - (Optional, String) Specifies the start time in UTC format, for example: 2020-09-01T18:50:20Z.

* `end_time` - (Optional, String) Specifies the end time in UTC format, for example: 2020-09-01T19:50:20Z.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dirty_data_list` - The dirty data list.

  The [dirty_data_list](#dirty_data_list_struct) structure is documented below.

<a name="dirty_data_list_struct"></a>
The `dirty_data_list` block supports:

* `db_name` - The database name.

* `schema_name` - The schema name.

* `table_name` - The table name.

* `error_sql` - The error SQL.

* `error_time` - The error occurrence time in UTC format, for example: 2023-06-10T03:01:52Z.

* `error_msg` - The error message.
