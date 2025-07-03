---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_tasks"
description: |-
  Use this data source to get the list of SMS migration tasks.
---

# huaweicloud_sms_tasks

Use this data source to get the list of SMS migration tasks.

## Example Usage

```hcl
data "huaweicloud_sms_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `state` - (Optional, String) Specifies the migration task status.
  Values can be as follows:
  + **READY**: The migration task is ready for execution.
  + **RUNNING**: The migration task is being executed.
  + **SYNCING**: The incremental data is being synchronized.
  + **MIGRATE_SUCCESS**: The migration succeeds.
  + **MIGRATE_FAIL**: The migration fails.
  + **ABORTING**: The migration task is being stopped.
  + **ABORT**: The migration task is stopped.
  + **DELETING**: The migration task is being deleted.
  + **SYNC_F_ROLLBACKING**: The synchronization fails and the task is being rolled back.
  + **SYNC_F_ROLLBACK_SUCCESS**: The synchronization fails and the rollback is successful.

* `name` - (Optional, String) Specifies the task name.

* `task_id` - (Optional, String) Specifies the task ID.

* `source_server_id` - (Optional, String) Specifies the source server ID.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the information about the queried tasks.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the migration task ID.

* `name` - Indicates the task name.

* `type` - Indicates the task type.

* `os_type` - Indicates the OS type.

* `state` - Indicates the task status.

* `estimate_complete_time` - Indicates the estimated completion time.

* `create_date` - Indicates the task creation time.

* `priority` - Indicates the migration process priority.

* `speed_limit` - Indicates the migration rate limit.

* `migrate_speed` - Indicates the migration rate, the unit is MB/s.

* `compress_rate` - Indicates the compression rate.

* `start_target_server` - Indicates whether the target server is started after the migration is complete.

* `error_json` - Indicates the error message.

* `total_time` - Indicates the task duration.

* `migration_ip` - Indicates the IP address of the target server.

* `sub_tasks` - Indicates the information about subtasks associated with the migration task

  The [sub_tasks](#tasks_sub_tasks_struct) structure is documented below.

* `source_server` - Indicates the information about the source server associated with the migration task.

  The [source_server](#tasks_source_server_struct) structure is documented below.

* `enterprise_project_id` - Indicates the migration project ID.

* `target_server` - Indicates the information about the target server associated with the migration task.

  The [target_server](#tasks_target_server_struct) structure is documented below.

* `log_collect_status` - Indicates the log collection status.

* `clone_server` - Indicates the information about the cloned server.

  The [clone_server](#tasks_clone_server_struct) structure is documented below.

* `syncing` - Indicates whether synchronization is enabled.

* `network_check_info` - Indicates the network performance metrics and measurement results.

  The [network_check_info](#tasks_network_check_info_struct) structure is documented below.

* `special_config` - Indicates the configuration information of advanced migration options.

  The [special_config](#tasks_special_config_struct) structure is documented below.

* `total_cpu_usage` - Indicates the CPU usage of the server, the unit is percentage.

* `agent_cpu_usage` - Indicates the CPU usage of the agent, the unit is percentage.

* `total_mem_usage` - Indicates the memory usage of the server, the unit is MB.

* `agent_mem_usage` - Indicates the memory usage of the agent, the unit is MB.

* `total_disk_io` - Indicates the disk I/O of the server, the unit is MB/s.

* `agent_disk_io` - Indicates the disk I/O of the agent, the unit is MB/s.

* `need_migration_test` - Indicates whether migration drilling is enabled.

<a name="tasks_sub_tasks_struct"></a>
The `sub_tasks` block supports:

* `id` - Indicates the subtask ID.

* `name` - Indicates the subtask name.

* `progress` - Indicates the progress of the subtask.

* `start_date` - Indicates the start time of the subtask.

* `end_date` - Indicates the end time of the subtask.

* `process_trace` - Indicates the detailed progress of the migration or synchronization.

<a name="tasks_source_server_struct"></a>
The `source_server` block supports:

* `id` - Indicates the ID of the source server in the SMS database.

* `ip` - Indicates the IP address of the source server.

* `name` - Indicates the source server name in SMS.

* `os_type` - Indicates the OS type of the source server.

* `os_version` - Indicates the OS version.

* `oem_system` - Indicates whether the OS is an OEM version (Windows).

* `state` - Indicates the source server status.

<a name="tasks_target_server_struct"></a>
The `target_server` block supports:

* `id` - Indicates the ID of the target server in the SMS database.

* `vm_id` - Indicates the ID of the target server.

* `name` - Indicates the name of the target server.

* `ip` - Indicates the IP address of the target server.

* `os_type` - Indicates the OS type of the target server.

* `os_version` - Indicates the OS version.

<a name="tasks_clone_server_struct"></a>
The `clone_server` block supports:

* `vm_id` - Indicates the ID of the cloned server.

* `name` - Indicates the name of the cloned server.

<a name="tasks_network_check_info_struct"></a>
The `network_check_info` block supports:

* `domain_connectivity` - Indicates the connectivity to domain names.

* `destination_connectivity` - Indicates the connectivity to the target server.

* `network_delay` - Indicates the network latency.

* `network_jitter` - Indicates the network jitter.

* `migration_speed` - Indicates the bandwidth.

* `loss_percentage` - Indicates the packet loss rate.

* `cpu_usage` - Indicates the CPU usage.

* `mem_usage` - Indicates the memory usage.

* `evaluation_result` - Indicates the network evaluation result.

<a name="tasks_special_config_struct"></a>
The `special_config` block supports:

* `config_key` - Indicates the advanced migration option.

* `config_value` - Indicates the value specified for the advanced migration option.

* `config_status` - Indicates the reserved field that describes the configuration status.
