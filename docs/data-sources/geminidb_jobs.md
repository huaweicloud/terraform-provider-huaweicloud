---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_jobs"
description: |-
  Use this data source to get the list of GeminiDB tasks.
---

# huaweicloud_geminidb_jobs

Use this data source to get the list of GeminiDB tasks.

## Example Usage

### Query all tasks

```hcl
data "huaweicloud_geminidb_jobs" "test" {}
```

### Query tasks by job ID

```hcl
variable "job_id" {}

data "huaweicloud_geminidb_jobs" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Optional, String) Specifies the task ID.

* `name` - (Optional, String) Specifies the task name.
  The valid values are as follows:
  + **CreateInstance**: Create an instance.
  + **RestoreNewInstance**: Restore data to a new instance.
  + **EnlargeInstance**: Add nodes.
  + **ReduceInstance**: Delete nodes.
  + **RestartInstance**: Restart an instance.
  + **RestartNode**: Restart a node.
  + **EnlargeInstanceVolume**: Scale up storage space of an instance.
  + **ReduceInstanceVolume**: Scale in storage space of an instance.
  + **ResizeInstance**: Change the specifications of an instance.
  + **UpgradeDbVersion**: Upgrade the engine version.
  + **BindPublicIP**: Bind an EIP to an instance.
  + **UnbindPublicIP**: Unbind an EIP from an instance.
  + **DeleteInstance**: Delete an instance.
  + **EnlargeInstanceColdVolume**: Scale up cold storage of an instance.
  + **AddInstanceColdVolume**: Enable cold storage for an instance.
  + **ModifySecurityGroup**: Modify a security group.
  + **ModifyCcmCert**: Modify a CCM certificate.
  + **ModifyPort**: Change a port.
  + **ConstructDisasterRecovery**: Establish a DR relationship.
  + **DeConstructDisasterRecovery**: Remove a DR relationship.
  + **SwitchOverDisasterRecovery**: Switch a DR relationship.
  + **BuildBiActiveInstance**: Create an instance with a dual-active DR relationship.
  + **ReleaseBiActiveInstance**: Remove a dual-active relationship from an instance.
  + **BackupInstance**: Back up an instance.

* `status` - (Optional, String) Specifies the database type.
  The valid values are as follows:
  + **Running**: Indicates a task is being executed
  + **Completed**: Indicates a task is complete.
  + **Failed**: Indicates a task failed.

* `start_time` - (Optional, String) Specifies query start time in the **yyyy-mm-ddThh:mm:ss** format.
  The default value is `30` days before the current date.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

* `end_time` - (Optional, String) Specifies query end time in the **yyyy-mm-ddThh:mm:ss** format.
  The default value is current time.
  It must be later than the start time and the time span cannot exceed `30` days.
  T is the separator between calendar and hourly notation of time. Z indicates the time zone offset.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of tasks.
  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - The tasks ID.

* `name` - The task name.

* `status` - The task status.

* `start_time` - The task start time.

* `end_time` - The task end time.

* `progress` - The task execution progress.

* `instance` - The details of the instance associated with the task.
  The [instance](#instance_struct) structure is documented below.

* `fail_reason` - The task failure information.

<a name="instance_struct"></a>
The `instance` block supports:

* `id` - The instance ID.

* `name` - The instance name.
