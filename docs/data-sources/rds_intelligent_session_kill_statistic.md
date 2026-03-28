---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_intelligent_session_kill_statistic"
description: |-
  Use this data source to get the real-time statistics of a DB instance's killing sessions.
---

# huaweicloud_rds_intelligent_session_kill_statistic

Use this data source to get the real-time statistics of a DB instance's killing sessions.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_intelligent_session_kill_statistic" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `node_id` - (Optional, String) Specifies the node ID. This parameter is valid only for cluster instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - Indicates the statistics based on different policies.

  The [statistics](#statistics_struct) structure is documented below.

<a name="statistics_struct"></a>
The `statistics` block supports:

* `keyword` - Indicates the throttling keyword extracted based on the statistics policy.

* `raw_sql_text` - Indicates the example SQL statement that matches the SQL throttling keyword.

* `ids` - Indicates the IDs of threads that meet the statistics policy.

* `count` - Indicates the total number of the thread IDs.

* `total_time` - Indicates the total execution time of the threads.

* `avg_time` - Indicates the average execution time of the threads.

* `max_time` - Indicates the maximum execution time of the threads.

* `strategy` - Indicates the statistics policy. The value can be:
  + **top3_time**: top 3 execution durations
  + **top3_count**: top 3 occurrence times
  + **top3_avg_time**: top 3 average execution durations

* `advice_concurrency` - Indicates the recommended maximum concurrency.

* `type` - Indicates the type of data filtered based on the statistics policy. The value can be:
  + **kill**: indicates the sessions to be terminated after session killing is delivered at the current statistical time.
  + **limit**: indicates the rules to be added when Enable SQL Throttling and Add Rule is selected at the current
    statistical time.
