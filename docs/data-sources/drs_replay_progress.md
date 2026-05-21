---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_replay_progress"
description: |-
  Use this data source to get the replay progress information for a DRS job within HuaweiCloud.
---

# huaweicloud_drs_replay_progress

Use this data source to get the replay progress information for a DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {
  type = string
}

data "huaweicloud_drs_replay_progress" "test" { 
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

* `progress` - The replay progress percentage.

* `parse_count` - The total number of SQL statements to be parsed.

* `replay_count` - The total number of replayed SQL statements.

* `task_mode` - The migration mode.

* `process_time` - The migration time, timestamp in milliseconds.

* `transfer_status` - The migration status.

* `max_time` - The maximum replay time, timestamp in milliseconds.

* `min_time` - The minimum replay time, timestamp in milliseconds.

* `now_time` - The current replay time, timestamp in milliseconds.

* `min_export_time` - The minimum time of the replay report, timestamp in milliseconds.

* `max_export_time` - The maximum time of the replay report, timestamp in milliseconds.

* `replay_sql_now_list` - The list of SQL statements currently being replayed.

  The [replay_sql_now_list](#replay_sql_now_list_struct) structure is documented below.

<a name="replay_sql_now_list_struct"></a>
The `replay_sql_now_list` block supports:

* `thread_id` - The session ID.

* `created_at` - The creation time.

* `modified_at` - The modification time.

* `shard_id` - The shard ID.

* `schema_name` - The schema name.

* `sql_statement` - The SQL statement content.

* `latency` - The original latency.

* `execute_latency` - The execution latency.

* `target_type` - The target database type.

* `target_name` - The target database identifier.

* `status` - The replay status.
  The valid values are as follows:
  + **running**: Executing.
  + **finish**: Completed.
