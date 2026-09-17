---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_replay_records"
description: |-
  Use this data source to query the replay SQL execution records of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_replay_records

Use this data source to query the replay SQL execution records of a DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_replay_records" "test" {
  job_id     = var.job_id
  start_time = 1700421943
  end_time   = 1700439961
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the replay SQL execution records.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `start_time` - (Required, Int) Specifies the start time of the SQL execution records, in seconds timestamp format.

* `end_time` - (Required, Int) Specifies the end time of the SQL execution records, in seconds timestamp format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of items in the list, which is independent of pagination.

* `sql_records` - The SQL execution records of the replay task.

  The [sql_records](#sql_records_struct) structure is documented below.

<a name="sql_records_struct"></a>
The `sql_records` block supports:

* `execute_time` - The execution time of the record, in seconds timestamp format.

* `finished_sql` - The number of finished SQL statements.

* `abnormal_sql` - The number of abnormal SQL statements.

* `slow_sql` - The number of slow SQL statements.

* `total_sql` - The total number of SQL statements in the execution record.
