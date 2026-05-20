---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_content_detail"
description: |-
  Use this data source to get the content comparison detail information of a DRS job.
---

# huaweicloud_drs_compare_content_detail

Use this data source to get the content comparison detail information of a DRS job.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_job_id" {}

data "huaweicloud_drs_compare_content_detail" "test" {
  job_id         = var.job_id
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `compare_job_id` - (Required, String) Specifies the ID of the compare job.

* `target_db_name` - (Optional, String) Specifies the target database name for filtering.

* `db_name` - (Optional, String) Specifies the source database name for filtering.

* `type` - (Optional, String) Specifies the comparison type.
  The valid values are as follows:
  + **compare**: Comparable
  + **unCompare**: Not comparable

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `content_compare_result_infos` - The list of content comparison result information.

  The [content_compare_result_infos](#content_compare_result_infos_struct) structure is documented below.

<a name="content_compare_result_infos_struct"></a>
The `content_compare_result_infos` block supports:

* `source_db` - The source database name.

* `target_db` - The target database name.

* `source_table_name` - The source table name.

* `target_table_name` - The target table name.

* `source_row_num` - The number of rows in the source table.

* `target_row_num` - The number of rows in the target table.

* `difference_row_num` - The difference value between the source table and target table.

* `line_compare_result` - The line comparison result.
  The valid values are as follows:
  + **true**: Consistent
  + **false**: Inconsistent

* `content_compare_result` - The content comparison result.
  The valid values are as follows:
  + **true**: Consistent
  + **false**: Inconsistent

* `message` - The additional information.

* `compare_line_config_filter` - The filter configuration for line comparison.

* `status` - The full comparison status.

* `complete_shard_count` - The number of completed shards.

* `total_shard_count` - The total number of shards.

* `progress` - The comparison progress percentage.
