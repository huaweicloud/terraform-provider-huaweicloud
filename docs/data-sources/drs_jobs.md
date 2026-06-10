---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_jobs"
description: |-
  Use this data source to get the list of DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_jobs

Use this data source to get the list of DRS jobs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_jobs" "test" {
  job_type = "migration"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_type` - (Required, String) Specifies the task scenario.
  The valid values are as follows:
  + **migration**: Real-time migration.
  + **sync**: Real-time synchronization.
  + **cloudDataGuard**: Real-time disaster recovery.

* `name` - (Optional, String) Specifies the task ID or name.
  Multiple task IDs can be separated by commas, up to 10.

* `status` - (Optional, String) Specifies the task status.
  The valid values are as follows:
  + **CREATING**: Creating.
  + **CREATE_FAILED**: Creation failed.
  + **CONFIGURATION**: Configuring.
  + **STARTJOBING**: Starting.
  + **WAITING_FOR_START**: Waiting to start.
  + **START_JOB_FAILED**: Task start failed.
  + **FULL_TRANSFER_STARTED**: Full migration in progress (initialization in disaster recovery scenario).
  + **FULL_TRANSFER_FAILED**: Full migration failed (initialization failed in disaster recovery scenario).
  + **FULL_TRANSFER_COMPLETE**: Full migration completed (initialization completed in disaster recovery scenario).
  + **INCRE_TRANSFER_STARTED**: Incremental migration in progress.
  + **INCRE_TRANSFER_FAILED**: Incremental migration failed (disaster recovery abnormal in disaster recovery scenario).
  + **RELEASE_RESOURCE_STARTED**: Ending task.
  + **RELEASE_RESOURCE_FAILED**: End task failed.
  + **RELEASE_RESOURCE_COMPLETE**: Ended.
  + **CHANGE_JOB_STARTED**: Task changing.
  + **CHANGE_JOB_FAILED**: Task change failed.
  + **CHILD_TRANSFER_STARTING**: Subtask starting.
  + **CHILD_TRANSFER_STARTED**: Subtask migrating.
  + **CHILD_TRANSFER_COMPLETE**: Subtask migration completed.
  + **CHILD_TRANSFER_FAILED**: Subtask migration failed.
  + **RELEASE_CHILD_TRANSFER_STARTED**: Subtask ending.
  + **RELEASE_CHILD_TRANSFER_COMPLETE**: Subtask ended.

* `engine_type` - (Optional, String) Specifies the engine type.
  The valid values are as follows:
  + **oracle-to-gaussdbv5**: Oracle to GaussDB distributed, used in real-time synchronization scenario.
  + **mysql-to-mysql**: MySQL to MySQL, used in real-time migration and synchronization scenarios.
  + **redis-to-gaussredis**: Redis to GeminiDB Redis, used in real-time migration scenario.
  + **rediscluster-to-gaussredis**: Redis cluster to GeminiDB Redis, used in real-time migration scenario.

* `net_type` - (Optional, String) Specifies the network type.
  The valid values are as follows:
  + **eip**: Public network.
  + **vpc**: VPC network.
  + **vpn**: VPN or dedicated line network.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  Defaults to "", meaning to query tasks of all enterprise projects.

* `sort_key` - (Optional, String) Specifies the keyword for sorting the result.
  The valid values are as follows: **name**, **status**, **create_time**, **net_type**, **job_direction**, **pay_mode**.
  Defaults to **create_time**.

* `sort_dir` - (Optional, String) Specifies the sorting direction.
  The valid values are as follows: **desc** (descending) and **asc** (ascending).
  Defaults to **desc**.

* `instance_ids` - (Optional, List) Specifies the database instance ID list.
  Defaults to null, meaning not to filter by database instance ID.

* `instance_ip` - (Optional, String) Specifies the IP address of the database instance bound to DRS.
  Defaults to "", meaning not to filter by DRS bound database IP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of DRS jobs.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - The task ID.

* `name` - The task name.

* `status` - The task status.

* `description` - The task description.

* `create_time` - The task creation time.

* `engine_type` - The engine type.

* `net_type` - The network type.

* `charging_mode` - The charging mode.
  The valid values are as follows:
  + **period**: Yearly/Monthly.
  + **on_demand**: On-demand.

* `billing_tag` - Whether the task is billed.

* `job_direction` - The task direction.
  The valid values are as follows:
  + **up**: Into the cloud (in disaster recovery scenario, the local cloud is the standby).
  + **down**: Out of the cloud (in disaster recovery scenario, the local cloud is the primary).
  + **non-dbs**: Self-built.

* `job_type` - The task scenario.

* `task_type` - The task mode.
  The valid values are as follows:
  + **FULL_TRANS**: Full.
  + **FULL_INCR_TRANS**: Full + incremental.
  + **INCR_TRANS**: Incremental.

* `enterprise_project_id` - The enterprise project ID.

* `job_mode` - The task mode.
  The valid values are as follows:
  + **single**: Single task.
  + **sync_child**: Synchronization subtask.
  + **multi_to_single**: Many-to-one task.

* `job_mode_role` - The task role.
  The valid values are as follows:
  + **parent**: Parent task.
  + **child**: Child task.
  + **master**: Primary task.
  + **slave**: Standby task.

* `is_multi_az` - Whether the task is a primary/standby task.

* `node_role` - The task node role.

* `node_new_framework` - Whether the task uses the new framework.

* `job_action` - The task action collection.

  The [job_action](#job_action_struct) structure is documented below.

* `children` - The subtask list.

  The [children](#children_struct) structure is documented below.

<a name="children_struct"></a>
The `children` block supports:

* `id` - The subtask ID.

* `name` - The subtask name.

* `status` - The subtask status.

* `description` - The subtask description.

* `create_time` - The subtask creation time.

* `engine_type` - The engine type.

* `net_type` - The network type.

* `charging_mode` - The charging mode.

* `billing_tag` - Whether the subtask is billed.

* `job_direction` - The task direction.

* `job_type` - The task scenario.

* `task_type` - The task mode.

* `enterprise_project_id` - The enterprise project ID.

* `job_mode` - The task mode.

* `job_mode_role` - The task role.

* `is_multi_az` - Whether the subtask is a primary/standby task.

* `node_role` - The task node role.

* `node_new_framework` - Whether the task uses the new framework.

* `job_action` - The task action collection.

  The [job_action](#job_action_struct) structure is documented above.

<a name="job_action_struct"></a>
The `job_action` block supports:

* `available_actions` - The list of available actions for the task.

* `unavailable_actions` - The list of unavailable actions for the task.

* `current_action` - The current action of the task.
  The valid values are as follows:
  + **API_CONFIGURATION_ACTION**: Task in OPEN API configuration.
  + **CHANGE**: Modify task.
  + **CHANGE_MODE**: Modify task mode.
  + **CHOOSE_OBJECT**: Select object.
  + **CLONE**: Clone task.
  + **CONTINUE_APPLY**: Start replay (for Oracle to GaussDB distributed).
  + **CONTINUE_CAPTURE**: Start capture (for Oracle to GaussDB distributed).
  + **CONTINUE_JOB**: Start failed or stopped task (for Oracle to GaussDB distributed).
  + **CREATE**: Create task.
  + **DELETE**: Delete task.
  + **FREE_RESOURCE**: Release resource.
  + **JUMP_RETRY**: Jump and retry task.
  + **MODIFY_CONFIGURATION**: Modify task configuration.
  + **MODIFY_DB_CONFIG**: Modify database configuration.
  + **MODIFY_TASK_NUMBER**: Modify thread count configuration.
  + **NODE_FLAVOR_MODIFY**: Specification change.
  + **ORDER_INFO**: Order details.
  + **PAUSE**: Pause task.
  + **PAY_ORDER**: Pay yearly/monthly order.
  + **PRE_CHECK**: Pre-check.
  + **QUERY_PRE_CHECK**: Query pre-check result.
  + **RESET**: Reset task.
  + **RESET_DB_PWD**: Reset database password (source and target databases).
  + **RETRY**: Retry task.
  + **START**: Start task.
  + **START_INCR**: Start incremental task.
  + **STOP_APPLY**: Stop replay (for Oracle to GaussDB distributed).
  + **STOP_CAPTURE**: Stop capture (for Oracle to GaussDB distributed).
  + **STOP_JOB**: Stop task (for Oracle to GaussDB distributed).
  + **SWITCH_OVER**: Disaster recovery switchover.
  + **TO_PERIOD**: Convert to yearly/monthly task.
  + **TO_RENEW**: Renew yearly/monthly task.
  + **UNSUBSCRIBE**: Unsubscribe yearly/monthly task.
