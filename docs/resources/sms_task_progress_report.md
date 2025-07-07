---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_progress_report"
description: |-
  Manages an SMS task progress report resource within HuaweiCloud.
---

# huaweicloud_sms_task_progress_report

Manages an SMS task progress report resource within HuaweiCloud.

~> Deleting task progress report resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_sms_task_progress_report" "test" {
  task_id       = var.task_id
  subtask_name  = "DETTACH_AGENT_IMAGE"
  progress      = 100
  replicatesize = 1000
  totalsize     = 100000
  process_trace = "migrate details"
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the migration task ID.

* `subtask_name` - (Required, String, NonUpdatable) Specifies the name of the subtask whose progress is reported.
  Values can be as follows:
  + **CREATE_CLOUD_SERVER**: creating a new server.
  + **SSL_CONFIG**: configuring a secure channel.
  + **ATTACH_AGENT_IMAGE**: attaching the disk that hosts the agent image.
  + **DETTACH_AGENT_IMAGE**: detaching the disk that hosts the agent image.
  + **FORMAT_DISK_LINUX**: formatting partitions on Linux.
  + **FORMAT_DISK_LINUX_FILE**: formatting partitions on Linux for a file-level migration.
  + **FORMAT_DISK_LINUX_BLOCK**: formatting partitions on Linux for a block-level migration.
  + **FORMAT_DISK_WINDOWS**: formatting partitions on Windows.
  + **MIGRATE_LINUX_FILE**: replicating files on Linux.
  + **MIGRATE_LINUX_BLOCK**: replicating blocks on Linux.
  + **MIGRATE_WINDOWS_BLOCK**: replicating blocks on Windows.
  + **CLONE_VM**: cloning the target server.
  + **SYNC_LINUX_FILE**: synchronizing files on Linux.
  + **SYNC_LINUX_BLOCK**: synchronizing blocks on Linux.
  + **SYNC_WINDOWS_BLOCK**: synchronizing blocks on Windows.
  + **CONFIGURE_LINUX**: modifying system configurations on Linux.
  + **CONFIGURE_LINUX_BLOCK**: modifying system configurations on Linux for a block-level migration.
  + **CONFIGURE_LINUX_FILE**: modifying system configurations on Linux for a file-level migration.
  + **CONFIGURE_WINDOWS**: modifying system configurations on Windows.

* `progress` - (Required, Int, NonUpdatable) Specifies the progress of the subtask, the unit is percentage.

* `replicatesize` - (Required, Int, NonUpdatable) Specifies the amount of data that has been replicated in the subtask,
  the unit is bytes.

* `totalsize` - (Required, Int, NonUpdatable) Specifies the total amount of data to be migrated in the subtask.

* `process_trace` - (Required, String, NonUpdatable) Specifies the detailed progress of the migration or synchronization.

* `migrate_speed` - (Optional, Int, NonUpdatable) Specifies the migration rate, the unit is Mbit/s.

* `compress_rate` - (Optional, Int, NonUpdatable) Specifies the file compression rate.

* `remain_time` - (Optional, Int, NonUpdatable) Specifies the remaining time.

* `total_cpu_usage` - (Optional, Int, NonUpdatable) Specifies the CPU usage of the server, the unit is percentage.

* `agent_cpu_usage` - (Optional, Int, NonUpdatable) Specifies the CPU usage of the agent, the unit is percentage.

* `total_mem_usage` - (Optional, Int, NonUpdatable) Specifies the memory usage of the server, the unit is MB.

* `agent_mem_usage` - (Optional, Int, NonUpdatable) Specifies the memory usage of the agent, the unit is MB.

* `total_disk_io` - (Optional, Int, NonUpdatable) Specifies the disk I/O of the server, the unit is MB/s.

* `agent_disk_io` - (Optional, Int, NonUpdatable) Specifies the disk I/O of the agent, the unit is MB/s.

* `need_migration_test` - (Optional, Bool, NonUpdatable) Specifies whether migration drilling is enabled.

* `agent_time` - (Optional, String, NonUpdatable) Specifies the current local time of the source server, which is used
  for overspeed detection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
