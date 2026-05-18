---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_actions"
description: |-
  Use this data source to get the list of DRS job actions within HuaweiCloud.
---

# huaweicloud_drs_actions

Use this data source to get the list of DRS job actions within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_actions" "test" { 
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

* `job_action` - The collection of allowed and disallowed operations.

  The [job_action](#job_action_struct) structure is documented below.

<a name="job_action_struct"></a>
The `job_action` block supports:

* `available_actions` - The list of available operations for the task.
  The valid values are as follows:

  + **CREATE**: Create a task.
  + **CHOOSE_OBJECT**: Select objects for editing during incremental synchronization.
  + **PRE_CHECK**: Perform pre-check.
  + **CHANGE_MODE**: Change task mode.
  + **FREE_RESOURCE**: Release resources.
  + **MODIFY_DB_CONFIG**: Modify database configuration.
  + **RESET_DB_PWD**: Reset database passwords (source and destination databases).
  + **MODIFY_CONFIGURATION**: Modify task configuration.
  + **PAUSE**: Pause the task.
  + **START**: Start the task.
  + **CHANGE**: Modify the task.
  + **RETRY**: Retry the task.
  + **RESET**: Reset the task.
  + **DELETE**: Delete the task.
  + **QUERY_PRE_CHECK**: Query pre-check results.
  + **SWITCH_OVER**: Perform disaster recovery switchover.
  + **START_INCR**: Start incremental task for Cassandra.
  + **MODIFY_TASK_NUMBER**: Modify thread count configuration for Cassandra.
  + **CONTINUE_JOB**: Continue failed or stopped tasks for Oracle-GaussDB distributed.
  + **STOP_JOB**: Stop tasks for Oracle-GaussDB distributed.
  + **CONTINUE_CAPTURE**: Start capture.
  + **STOP_CAPTURE**: Stop capture.
  + **CONTINUE_APPLY**: Start apply.
  + **STOP_APPLY**: Stop apply.
  + **PAY_ORDER**: Pay for yearly/monthly subscription orders.
  + **UNSUBSCRIBE**: Unsubscribe from yearly/monthly subscription.
  + **TO_PERIOD**: Convert to periodic subscription.
  + **TO_RENEW**: Renew yearly/monthly subscription.
  + **ORDER_INFO**: View order details.
  + **CHANGE_FLAVOR**: Change specifications.
  + **CLONE**: Clone the task.

* `unavailable_actions` - The list of unavailable operations for the task.
  The valid values are as follows:

  + **CREATE**: Create a task.
  + **CHOOSE_OBJECT**: Select objects for editing during incremental synchronization.
  + **PRE_CHECK**: Perform pre-check.
  + **CHANGE_MODE**: Change task mode.
  + **FREE_RESOURCE**: Release resources.
  + **MODIFY_DB_CONFIG**: Modify database configuration.
  + **RESET_DB_PWD**: Reset database passwords (source and destination databases).
  + **MODIFY_CONFIGURATION**: Modify task configuration.
  + **PAUSE**: Pause the task.
  + **START**: Start the task.
  + **CHANGE**: Modify the task.
  + **RETRY**: Retry the task.
  + **RESET**: Reset the task.
  + **DELETE**: Delete the task.
  + **QUERY_PRE_CHECK**: Query pre-check results.
  + **SWITCH_OVER**: Perform disaster recovery switchover.
  + **START_INCR**: Start incremental task for Cassandra.
  + **MODIFY_TASK_NUMBER**: Modify thread count configuration for Cassandra.
  + **CONTINUE_JOB**: Continue failed or stopped tasks for Oracle-GaussDB distributed.
  + **STOP_JOB**: Stop tasks for Oracle-GaussDB distributed.
  + **CONTINUE_CAPTURE**: Start capture.
  + **STOP_CAPTURE**: Stop capture.
  + **CONTINUE_APPLY**: Start apply.
  + **STOP_APPLY**: Stop apply.
  + **PAY_ORDER**: Pay for yearly/monthly subscription orders.
  + **UNSUBSCRIBE**: Unsubscribe from yearly/monthly subscription.
  + **TO_PERIOD**: Convert to periodic subscription.
  + **TO_RENEW**: Renew yearly/monthly subscription.
  + **ORDER_INFO**: View order details.
  + **CHANGE_FLAVOR**: Change specifications.
  + **CLONE**: Clone the task.

* `current_action` - The current operation of the task.
