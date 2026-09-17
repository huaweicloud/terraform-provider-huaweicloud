---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_export_jobs"
description: |-
  Use this data source to export the list of DRS tasks matching the filter parameters within HuaweiCloud.
---

# huaweicloud_drs_export_jobs

Use this data source to export the list of DRS tasks matching the filter parameters within HuaweiCloud.

-> This data source is used to trigger an asynchronous export task. After a successful export, an async job ID will be
   returned, which can be used to query the export progress and obtain the export result.

## Example Usage

### Export backup migration tasks

```hcl
data "huaweicloud_drs_export_jobs" "test" {
  job_type = "backupMigration"
}
```

### Export real-time migration tasks with filter

```hcl
variable "job_name" {}

data "huaweicloud_drs_export_jobs" "test" {
  job_type = "migration"
  name     = var.job_name
  status   = "FULL_TRANSFER_COMPLETE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to export the DRS tasks.
  If omitted, the provider-level region will be used.

* `job_type` - (Required, String) Specifies the business scenario type of the DRS tasks to be exported.  
  The valid values are as follows:
  + **migration**: Real-time migration.
  + **sync**: Real-time synchronization.
  + **cloudDataGuard**: Real-time disaster recovery.
  + **backupMigration**: Backup migration.
  + **subscription**: Subscription.

* `name` - (Optional, String) Specifies the name of the DRS tasks to be exported.

* `status` - (Optional, String) Specifies the status of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, the valid values are as follows:
  + **CREATING**: The task is being created.
  + **CREATE_FAILED**: The task fails to be created.
  + **CONFIGURATION**: The task is being configured.
  + **STARTJOBING**: The task is being started.
  + **WAITING_FOR_START**: The task is waiting to be started.
  + **START_JOB_FAILED**: The task fails to be started.
  + **FULL_TRANSFER_STARTED**: The full migration is in progress, and the disaster recovery is being initialized.
  + **FULL_TRANSFER_FAILED**: The full migration fails, and the disaster recovery fails to be initialized.
  + **FULL_TRANSFER_COMPLETE**: The full migration is complete, and the disaster recovery is initialized.
  + **INCRE_TRANSFER_STARTED**: The incremental migration is in progress, and the disaster recovery is in progress.
  + **INCRE_TRANSFER_FAILED**: The incremental migration fails, and the disaster recovery is abnormal.
  + **RELEASE_RESOURCE_STARTED**: The task is being ended.
  + **RELEASE_RESOURCE_FAILED**: The task fails to be ended.
  + **RELEASE_RESOURCE_COMPLETE**: The task has ended.
  + **CHANGE_JOB_STARTED**: The task is being changed.
  + **CHANGE_JOB_FAILED**: The task fails to be changed.
  + **CHILD_TRANSFER_STARTING**: The subtask is being started.
  + **CHILD_TRANSFER_STARTED**: The subtask is being migrated.
  + **CHILD_TRANSFER_COMPLETE**: The subtask migration is complete.
  + **CHILD_TRANSFER_FAILED**: The subtask migration fails.
  + **RELEASE_CHILD_TRANSFER_STARTED**: The subtask is being ended.
  + **RELEASE_CHILD_TRANSFER_COMPLETE**: The subtask has ended.
  
  When the `job_type` is set to **backupMigration**, the valid values are as follows:
  + **SUBSCRIPTION_STARTED**: The task is being created.
  + **CREATE_FAILED**: The task fails to be created.
  + **SUBSCRIPTION_FAILED**: The subscription fails.
  + **CONFIGURATION**: The task is being configured.
  + **STARTJOBING**: The task is being started.
  + **REBUILD_NODE_STARTED**: The task is in fault recovery.
  + **REBUILD_NODE_FAILED**: The task fails to recover from a fault.
  + **START_JOB_FAILED**: The task fails to be started.
  
  When the `job_type` is set to **subscription**, the valid values are as follows:
  + **SUCCESS**: The subscription is successful.
  + **TRANSFERRING**: The backup is being migrated.
  + **FAILED**: The subscription fails.
  + **PRECHECK FAILED**: The precheck fails.

* `dbs_instance_ids` - (Optional, List) Specifies the list of database instance IDs, with a maximum of 10 IDs.

* `description` - (Optional, String) Specifies the description of the DRS tasks to be exported.

* `create_at` - (Optional, String) Specifies the creation time range of the DRS tasks to be exported, the format is
  **YYYY-MM-DDTHH:MM:SS.mmmZ|YYYY-MM-DDTHH:MM:SS.mmmZ**, for example, **2026-05-25T13:41:00.127Z|2026-05-28T13:41:00.999Z**.
  When the `job_type` is set to **subscription**, this parameter indicates the consumption time point.

* `completed_at` - (Optional, String) Specifies the completion time range of the DRS tasks to be exported, the format is
  **YYYY-MM-DDTHH:MM:SS.mmmZ|YYYY-MM-DDTHH:MM:SS.mmmZ**.
  This parameter is only supported when the `job_type` is set to **backupMigration**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the DRS tasks to be exported.

* `engine_type` - (Optional, String) Specifies the engine type of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.
  When the `job_type` is set to **backupMigration**, the default value is **sqlserver**.

* `net_type` - (Optional, String) Specifies the network type of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.  
  The valid values are as follows:
  + **eip**: Public network.
  + **vpc**: VPC network.
  + **vpn**: VPN or dedicated network.

* `billing_tag` - (Optional, String) Specifies the billing flag of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync**, **cloudDataGuard** or **subscription**, this parameter is
  supported.  
  The valid values are as follows:
  + **0**: Billing.
  + **1**: Not billed.

* `billing_mode` - (Optional, String) Specifies the billing mode of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.  
  The valid values are as follows:
  + **0**: Pay-per-use.
  + **1**: Yearly/Monthly.

* `public_ip` - (Optional, String) Specifies the public IP bound to the DRS tasks to be exported.
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.

* `instance_ip` - (Optional, String) Specifies the IP of the database instance bound to the DRS tasks to be exported.
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.

* `inner_ip` - (Optional, String) Specifies the inner IP bound to the DRS tasks to be exported.
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.

* `spec_type` - (Optional, String) Specifies the specification type of the DRS tasks to be exported.  
  When the `job_type` is set to **migration**, **sync** or **cloudDataGuard**, this parameter is supported.  
  The valid values are as follows:
  + **2xlarge**: Extra-large specification.
  + **xlarge**: Large specification.
  + **high**: Medium-high specification.
  + **medium**: Medium specification.
  + **small**: Small specification.
  + **micro**: Minimum specification.

* `direction` - (Optional, String) Specifies the data flow direction of the DRS tasks to be exported.  
  When the `job_type` is set to **migration** or **sync**, the valid values are as follows:
  + **up**: Upstream to cloud.
  + **down**: Downstream from cloud.
  + **non-dbs**: Self-built.
  + **multiWrite**: Bidirectional synchronization.
  
  When the `job_type` is set to **cloudDataGuard**, the valid values are as follows:
  + **up**: This cloud is the standby.
  + **down**: This cloud is the primary.
  + **multiWrite**: Dual-active disaster recovery.

* `task_type` - (Optional, String) Specifies the synchronization type of the DRS tasks to be exported.  
  When the `job_type` is set to **migration** or **sync**, this parameter is supported.  
  The valid values are as follows:
  + **FULL_INCR_TRANS**: Full-incremental migration.
  + **FULL_TRANS**: Full migration.
  + **INCR_TRANS**: Incremental migration.

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The default value is **create_time**. The supported keywords are as follows:
  + **name**
  + **status**
  + **create_time**
  + **net_type**
  + **job_direction**
  + **pay_mode**

* `sort_dir` - (Optional, String) Specifies the sorting order of the returned results.
  The default value is **desc**.
  The valid values are as follows:
  + **desc**: Descending order.
  + **asc**: Ascending order.

* `tag` - (Optional, String) Specifies the tag of the DRS tasks to be exported, in JSON format, for example,
  **{"key": "value"}**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `async_job_id` - The ID of the asynchronous export task.
