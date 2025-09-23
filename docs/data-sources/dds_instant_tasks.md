---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instant_tasks"
description: |-
  Use this data source to get the list of DDS instant tasks.
---

# huaweicloud_dds_instant_tasks

Use this data source to get the list of DDS instant tasks.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dds_instant_tasks" "test" {
  start_time = var.start_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `start_time` - (Required, String) Specifies the start time. The format of the start time is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - (Required, String) Specifies the end time. The format of the end time is **yyyy-mm-ddThh:mm:ssZ**
  and the end time must be later than the start time. The time span cannot be longer than 30 days.

* `status` - (Optional, String) Specifies the task status.
  + **Running**: Indicates that the task is being executed.
  + **Completed**: Indicates that the task is completed.
  + **Failed**: Indicates that the task fails.

* `name` - (Optional, String) Specifies the task name. The value can be:
  + **CreateMongoDB**: Create a cluster instance.
  + **CreateMongoDBReplica**: Create a replica set instance.
  + **CreateMongoDBReplicaSingle**: Create a single node instance.
  + **EnlargeMongoDBVolume**: Scale up the storage capacity of a DB instance.
  + **ResizeMongoDBInstance**: Change the class of a DB instance of Community Edition.
  + **ResizeDfvMongoDBInstance**: Change the class of a DB instance of Enhanced Edition.
  + **EnlargeMongoDBGroup**: Add a node.
  + **ReplicaSetEnlargeNode**: Add a standby node to a replica set instance.
  + **AddReadonlyNode**: Add a read replica.
  + **RestartInstance**: Restart a cluster instance.
  + **RestartGroup**: Restart a cluster node group.
  + **RestartNode**: Restart a cluster node.
  + **RestartReplicaSetInstance**: Restart a replica set instance.
  + **RestartReplicaSingleInstance**: Restart a single node instance.
  + **SwitchPrimary**: Perform a primary/standby switchover.
  + **ModifyIp**: Change the private IP address.
  + **ModifySecurityGroup**: Modify a security group.
  + **ModifyPort**: Change the database port.
  + **BindPublicIP**: Bind an EIP.
  + **UnbindPublicIP**: Unbind an EIP.
  + **SwitchInstanceSSL**: Switch the SSL.
  + **AzMigrate**: Migrate a DB instance from one AZ to another.
  + **CreateIp**: Enable the shard/config IP address.
  + **ModifyOpLogSize**: Change the oplog size.
  + **RestoreMongoDB**: Restore a cluster instance to a new DB instance.
  + **RestoreMongoDB_Replica**: Restore a replica set instance to a new DB instance.
  + **RestoreMongoDB_Replica_Single**: Restore a single node instance to a new DB instance.
  + **RestoreMongoDB_Replica_PITR**: Restore a replica set instance to a specified point in time.
  + **MongodbSnapshotBackup**: Create a physical backup.
  + **MongodbSnapshotEBackup**: Create a snapshot backup.
  + **MongodbRestoreData2CurrentInstance**: Restore a backup to the original DB instance.
  + **MongodbRestoreData2NewInstance**: Restore a backup to a new DB instance.
  + **MongodbPitr2CurrentInstance**: Restore a backup to a specified time point of the original DB instance.
  + **MongodbPitr2NewInstance**: Restore a backup to a specified time point of a new DB instance.
  + **MongodbRecycleBackup**: Restore a backup from the recycle bin.
  + **MongodbRestoreTable**: Restore databases and tables to a specified point in time.
  + **UpgradeDatabaseVersion**: Upgrade the database patch.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the tasks list.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `fail_reason` - Indicates the task failure information.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `progress` - Indicates the task execution progress.

* `status` - Indicates the task status.

* `created_at` - Indicates the task creation time.

* `ended_at` - Indicates the task end time.
