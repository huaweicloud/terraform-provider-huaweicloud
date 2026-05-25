---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_smn_batch_set"
description: |-
  Manages a resource to batch set SMN for DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_smn_batch_set

Manages a resource to batch set SMN for DRS jobs within HuaweiCloud.

-> This resource is a one-time action resource used to batch set SMN for DRS jobs. Deleting this resource
  will not undo the SMN configuration, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "jobs" { 
  type = list(object({
    job_id      = string
    status      = string
    engine_type = string
  }))
}
variable "topic_urn" {}

resource "huaweicloud_drs_pwd_batch_modify" "test" {
  dynamic "jobs" {
    for_each = var.jobs
    
    content {
      job_id      = jobs.value.job_id
      status      = jobs.value.status
      engine_type = jobs.value.engine_type
    }
  }

  alarm_notify_info { 
    topic_urn     = var.topic_urn 
    delay_time    = 1200 
    rto_delay     = 20 
    rpo_delay     = 20 
    alarm_to_user = false 
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the list of DRS jobs to configure SMN.
  The [jobs](#jobs_struct) structure is documented below.

* `alarm_notify_info` - (Required, List, NonUpdatable) Specifies the alarm notification information.
  The [alarm_notify_info](#alarm_notify_info_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the DRS job ID.

* `status` - (Optional, String, NonUpdatable) Specifies the job status.
  The valid values are as follows:
  + **CREATING**: Creating.
  + **CREATE_FAILED**: Creation failed.
  + **CONFIGURATION**: Configuring.
  + **STARTJOBING**: Starting.
  + **WAITING_FOR_START**: Waiting to start.
  + **START_JOB_FAILED**: Start failed.
  + **PAUSING**: Paused.
  + **FULL_TRANSFER_STARTED**: Full transfer started (initialized for disaster recovery).
  + **FULL_TRANSFER_FAILED**: Full transfer failed (initialization failed for disaster recovery).
  + **FULL_TRANSFER_COMPLETE**: Full transfer completed (initialization completed for disaster recovery).
  + **INCRE_TRANSFER_STARTED**: Incremental transfer started (disaster recovery active).
  + **INCRE_TRANSFER_FAILED**: Incremental transfer failed (disaster recovery abnormal).
  + **RELEASE_RESOURCE_STARTED**: Releasing resources.
  + **RELEASE_RESOURCE_FAILED**: Resource release failed.
  + **RELEASE_RESOURCE_COMPLETE**: Completed.
  + **REBUILD_NODE_STARTED**: Recovering from failure.
  + **REBUILD_NODE_FAILED**: Recovery from failure failed.
  + **CHANGE_JOB_STARTED**: Job changing.
  + **CHANGE_JOB_FAILED**: Job change failed.
  + **DELETED**: Deleted.
  + **CHILD_TRANSFER_STARTING**: Child task starting (re-editing).
  + **CHILD_TRANSFER_STARTED**: Child task transferring (re-editing).
  + **CHILD_TRANSFER_COMPLETE**: Child task transfer completed (re-editing).
  + **CHILD_TRANSFER_FAILED**: Child task transfer failed (re-editing).
  + **RELEASE_CHILD_TRANSFER_STARTED**: Child task releasing (re-editing).
  + **RELEASE_CHILD_TRANSFER_COMPLETE**: Child task completed (re-editing).
  + **NODE_UPGRADE_START**: Upgrade started.
  + **NODE_UPGRADE_COMPLETE**: Upgrade completed.
  + **NODE_UPGRADE_FAILED**: Upgrade failed.

* `engine_type` - (Optional, String, NonUpdatable) Specifies the DRS job engine type.
  The valid values are as follows:
  + **mysql**: MySQL to MySQL migration and synchronization.
  + **mongodb**: MongoDB to DDS migration and synchronization.
  + **cloudDataGuard-mysql**: MySQL to MySQL disaster recovery.
  + **gaussdbv5**: GaussDB synchronization.
  + **mysql-to-kafka**: MySQL to Kafka synchronization.
  + **taurus-to-kafka**: TaurusDB to Kafka synchronization.
  + **gaussdbv5ha-to-kafka**: GaussDB centralized edition to Kafka synchronization.
  + **postgresql**: PostgreSQL to PostgreSQL synchronization.
  + **oracle-to-gaussdbv5**: Oracle to GaussDB distributed edition.
  + **oracle-to-gaussdbv5ha**: Oracle to GaussDB centralized edition.
  + **gaussdbv5-to-oracle**: GaussDB distributed edition to Oracle.
  + **gaussdbv5ha-to-oracle**: GaussDB centralized edition to Oracle.
  + **gaussdbv5-to-kafka**: GaussDB distributed edition to Kafka.

<a name="alarm_notify_info_struct"></a>
The `alarm_notify_info` block supports:

* `topic_urn` - (Optional, String, NonUpdatable) Specifies the SMN topic URN.

* `delay_time` - (Optional, Int, NonUpdatable) Specifies the subscription delay time in seconds.

* `rto_delay` - (Optional, Int, NonUpdatable) Specifies the RTO (Recovery Time Objective) delay threshold in seconds.

* `rpo_delay` - (Optional, Int, NonUpdatable) Specifies the RPO (Recovery Point Objective) delay threshold in seconds.

* `alarm_to_user` - (Optional, Bool, NonUpdatable) Specifies whether to send alarm notifications to users.
  Defaults to **false**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The results of batch setting SMN.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `status` - The operation status. The valid values are as follows:
  + **success**: Success.
  + **failed**: Failed.
