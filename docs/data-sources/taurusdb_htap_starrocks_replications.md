---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_replications"
description: |-
  Use this data source to query the data replication tasks of a TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_replications

Use this data source to query the data replication tasks of a TaurusDB HTAP StarRocks instance within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_replications" "test" {
  instance_id = var.htap_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HTAP StarRocks data replication tasks.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP StarRocks instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `replications` - The list of data replication tasks.
  The [replications](#htap_starrocks_replications_attr) structure is documented below.

<a name="htap_starrocks_replications_attr"></a>
The `replications` block supports:

* `source_database` - The source TaurusDB database name.

* `target_database` - The target StarRocks database name.

* `task_name` - The replication task name.

* `status` - The current status.
  The valid values are as follows:
  + **Yes**: Normal.
  + **No**: Abnormal.

* `stage` - The synchronization stage.
  The valid values are as follows:
  + **wait**: Waiting for synchronization.
  + **incremental**: Incremental synchronization.
  + **full**: Full synchronization.
  + **cancelled**: Synchronization cancelled.
  + **paused**: Synchronization paused.

* `percentage` - The progress percentage.

* `is_need_repair` - Whether the task need to be repaired.

* `is_main_task` - Whether the task is the main task.
