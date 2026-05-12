---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscriptions"
description: |-
  Use this data source to get a list of DRS subscription jobs.
---

# huaweicloud_drs_subscriptions

Use this data source to get a list of DRS subscription jobs.

## Example Usage

```hcl
variable "job_type" {}

data "huaweicloud_drs_subscriptions" "test" { 
  job_type = var.job_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_type` - (Required, String) Specifies the job scenario.

* `engine_type` - (Optional, String) Specifies the engine type.

* `net_type` - (Optional, String) Specifies the network type.

* `name` - (Optional, String) Specifies the task ID or name.

* `status` - (Optional, String) Specifies the task status.

* `description` - (Optional, String) Specifies the description of the job.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `instance_ids` - (Optional, List) Specifies the list of database instance IDs. A maximum of 10 IDs are supported.

* `instance_ip` - (Optional, String) Specifies the database instance IP address.

* `sort_key` - (Optional, String) Specifies the sort key for the returned results. The default value is `create_time`.

* `sort_dir` - (Optional, String) Specifies the sort order, which can be **desc** (descending) or **asc** (ascending).
  The default value is `desc`.

* `service_name` - (Optional, String) Specifies the service name.

* `is_billing` - (Optional, Bool) Specifies whether billing is enabled. The values can be **true**, **false**, or
  omitted for all. The default is all.

* `begin_at` - (Optional, String) Specifies the start time.

* `tags` - (Optional, Map) Specifies the tags of the job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The list of subscription jobs.

The [subscriptions](#subscriptions_struct) structure is documented below.

<a name="subscriptions_struct"></a>
The `subscriptions` block supports:

* `id` - The job ID.

* `name` - The job name.

* `status` - The job status.

* `created_time` - The task creation time.

* `begin_time` - The consumption start time.

* `now_time` - The current time.

* `description` - The job description.

* `enterprise_project_id` - The enterprise project ID.

* `job_action` - The job action information.

The [job_action](#job_action_struct) structure is documented below.

<a name="job_action_struct"></a>
The `job_action` block supports:

* `available_actions` - The list of available actions for the task.
  The valid values are as follows:
  + **CREATE**: Create a task.
  + **CHOOSE_OBJECT**: Select objects (re-edit during incremental phase).
  + **PRE_CHECK**: Pre-check.
  + **CHANGE_MODE**: Modify task mode.
  + **FREE_RESOURCE**: Release resources.
  + **MODIFY_DB_CONFIG**: Modify database configuration.
  + **RESET_DB_PWD**: Reset database passwords (source and target databases).
  + **MODIFY_CONFIGURATION**: Modify task configuration.
  + **PAUSE**: Pause the task.
  + **START**: Start the task.
  + **CHANGE**: Modify the task.
  + **RETRY**: Retry the task.
  + **RESET**: Reset the task.
  + **DELETE**: Delete the task.
  + **QUERY_PRE_CHECK**: Query pre-check results.
  + **SWITCH_OVER**: Disaster recovery switchover.
  + **START_INCR**: Start incremental task (Cassandra).
  + **MODIFY_TASK_NUMBER**: Modify thread count configuration (Cassandra).
  + **CONTINUE_JOB**: Start failed or stopped tasks (Oracle to GaussDB Distributed).
  + **STOP_JOB**: Stop the task (Oracle to GaussDB Distributed).
  + **CONTINUE_CAPTURE**: Start capture (Oracle to GaussDB Distributed).
  + **STOP_CAPTURE**: Stop capture (Oracle to GaussDB Distributed).
  + **CONTINUE_APPLY**: Start replay (Oracle to GaussDB Distributed).
  + **STOP_APPLY**: Stop replay (Oracle to GaussDB Distributed).
  + **PAY_ORDER**: Pay for yearly/monthly subscription orders.
  + **UNSUBSCRIBE**: Unsubscribe from yearly/monthly subscriptions.
  + **TO_PERIOD**: Convert to yearly/monthly subscription.
  + **TO_RENEW**: Renew yearly/monthly subscription.
  + **ORDER_INFO**: Order details.
  + **CHANGE_FLAVOR**: Change specifications.
  + **CLONE**: Clone the task.

* `unavailable_actions` - The list of unavailable actions for the task. The values are the same as `available_actions`.

* `current_action` - The current action command of the task.
  The valid values are as follows:
  + **API_CONFIGURATION_ACTION**: Callable by tasks configured via Open API.
  + **CHANGE**: Modify the task.
  + **CHANGE_MODE**: Modify task mode.
  + **CHOOSE_OBJECT**: Select objects.
  + **CLONE**: Clone the task.
  + **CONTINUE_APPLY**: Start replay (applicable to Oracle synchronized to GaussDB Distributed).
  + **CONTINUE_CAPTURE**: Start capture (applicable to Oracle synchronized to GaussDB Distributed).
  + **CONTINUE_JOB**: Start failed or stopped tasks (applicable to Oracle synchronized to GaussDB Distributed).
  + **CREATE**: Create a task.
  + **DELETE**: Delete the task.
  + **FREE_RESOURCE**: Release resources.
  + **JUMP_RETRY**: Jump retry for the task.
  + **MODIFY_CONFIGURATION**: Modify task configuration.
  + **MODIFY_DB_CONFIG**: Modify database configuration.
  + **MODIFY_TASK_NUMBER**: Modify thread count configuration.
  + **NODE_FLAVOR_MODIFY**: Change specifications.
  + **ORDER_INFO**: Order details.
  + **PAUSE**: Pause the task.
  + **PAY_ORDER**: Pay for yearly/monthly subscription orders.
  + **PRE_CHECK**: Pre-check.
  + **QUERY_PRE_CHECK**: Query pre-check results.
  + **RESET**: Reset the task.
  + **RESET_DB_PWD**: Reset database passwords (source and target databases).
  + **RETRY**: Retry the task.
  + **START**: Start the task.
  + **START_INCR**: Start incremental task.
  + **STOP_APPLY**: Stop replay (applicable to Oracle synchronized to GaussDB Distributed).
  + **STOP_CAPTURE**: Stop capture (applicable to Oracle synchronized to GaussDB Distributed).
  + **STOP_JOB**: Stop the task (applicable to Oracle synchronized to GaussDB Distributed).
  + **SWITCH_OVER**: Disaster recovery switchover.
  + **TO_PERIOD**: Convert to yearly/monthly subscription.
  + **TO_RENEW**: Renew yearly/monthly subscription.
  + **UNSUBSCRIBE**: Unsubscribe from yearly/monthly subscriptions.
