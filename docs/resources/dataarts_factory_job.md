---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_job"
description: ""
---

# huaweicloud_dataarts_factory_job

Manages a job resource of DataArts Factory within HuaweiCloud.

A job consists of one or more nodes, such as Hive SQL and CDM Job nodes.
DLF supports two types of jobs: batch jobs and real-time jobs.

## Example Usage

```hcl
variable "workspace_id" {}
variable "cmd_name" {}

resource "huaweicloud_dataarts_factory_job" "test" {
  workspace_id = var.workspace_id
  name         = "demo"
  process_type = "REAL_TIME"

  nodes {
    name = "Rest_client_demo_1"
    type = "RESTAPI"
    location {
      x = 10
      y = 10
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = var.cmd_name
    }

  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Job name.  
  The name contains a maximum of 128 characters, including only letters, numbers, hyphens (-),
  underscores (_), and periods (.). The job name must be unique.

* `nodes` - (Required, List) Node definition.
  The [nodes](#job_Node) structure is documented below.

* `schedule` - (Required, List) Scheduling configuration.
  The [schedule](#job_Schedule) structure is documented below.

* `process_type` - (Required, String) Job type.  
  The valid values are as follows:
    - **REAL_TIME**: real-time processing.
    - **BATCH**: batch processing.

* `workspace_id` - (Optional, String, ForceNew) The workspace ID.
  If this parameter is not set, the default workspace is used by default.
  Changing this parameter will create a new resource.

* `params` - (Optional, List) Job parameter definition.
  The [params](#job_Param) structure is documented below.

* `directory` - (Optional, String) Path of a job in the directory tree.  
  If the directory of the path does not exist during job creation, a directory is automatically
  created in the root directory /, for example, /dir/a/.

* `log_path` - (Optional, String) The OBS path where job execution logs are stored.

* `basic_config` - (Optional, List) Baisc job information.
  The [basic_config](#job_BasicConfig) structure is documented below.

<a name="job_Node"></a>
The `nodes` block supports:

* `name` - (Required, String) Node name.  
  The name contains a maximum of 128 characters, including only letters, numbers, hyphens (-),
  underscores (_), and periods (.). Names of the nodes in a job must be unique.

* `type` - (Required, String) Node type.  
  The options are as follows:
    - **HiveSQL**: Runs Hive SQL scripts.
    - **SparkSQL**: Runs Spark SQL scripts.
    - **DWSSQL**: Runs DWS SQL scripts.
    - **DLISQL**: Runs DLI SQL scripts.
    - **Shell**: Runs shell SQL scripts.
    - **CDMJob**: Runs CDM jobs.
    - **DISTransfer Task**: Creates DIS dump tasks.
    - **CloudTableManager**: Manages CloudTable tables, including creating and deleting tables.
    - **OBSManager**: Manages OBS paths, including creating and deleting paths.
    - **RESTClient**: Sends REST API requests.
    - **SMN**: Sends short messages or emails.
    - **MRSSpark**: Runs Spark jobs of MRS.
    - **MapReduce**: Runs MapReduce jobs of MRS.
    - **MRSFlinkJob**: Runs Flink jobs of MRS.
    - **MRSHetuEngine**: Runs HetuEngine jobs of MRS.
    - **DLISpark**: Runs Spark jobs of DLF.
    - **RDSSQL**: Transfers SQL statements to RDS for execution.

* `location` - (Required, List) Location of a node on the job canvas
  The [location](#job_Location) structure is documented below.

* `pre_node_name` - (Optional, List) Name of the previous node on which the current node depends.

* `conditions` - (Optional, List) Node execution condition.  
  Whether the node is executed or not depends on the calculation result of the EL expression saved
  in the expression field of condition.
  The [conditions](#job_Condition) structure is documented below.

* `properties` - (Required, List) Node property. Each type of node has its own property definition.  
  - **HiveSQL**: For details, see [Table 14](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **SparkSQL**: For details, see [Table 15](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **DWSSQL**: For details, see [Table 16](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **DLISQL**: For details, see [Table 17](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **Shell**: For details, see [Table 18](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **CDMJob**: For details, see [Table 19](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **DISTransferTask**: For details, see [Table 20](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **CloudTableManager**: For details, see [Table 21](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **OBSManager**: For details, see [Table 22](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **RESTClient**: For details, see [Table 23](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **SMN**: For details, see [Table 24](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **MRSSpark**: For details, see [Table 25](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **MapReduce**: For details, see [Table 26](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **DLISpark**: For details, see [Table 27](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **MRSFlinkJob**: For details, see [Table 29](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).
  - **MRSHetuEngine**: For details, see [Table 30](https://support.huaweicloud.com/intl/en-us/api-dataartsstudio/dataartsstudio_02_0084.html).

  The [properties](#job_NodeProperty) structure is documented below.

* `polling_interval` - (Optional, Int) Interval at which node running results are checked.  
  Unit: second; value range: `1` to `60`.
  Default value: `10`.

* `max_execution_time` - (Optional, Int) Maximum execution time of a node.  
  If a node is not executed within the maximum execution time, the node is set to the failed state.  
  Unit: minute; value range: `5` to `1,440`.
  Default value: `60`.

* `retry_times` - (Optional, Int) Number of the node retries.  
  The value ranges from `0` to `5`. `0` indicates no retry.  
  Default value: `0`.

* `retry_interval` - (Optional, Int) Interval at which a retry is performed upon a failure.  
  The value ranges from `5` to `120`.  
  Unit: second.  
  Default value: `120`.

* `fail_policy` - (Optional, String) Node failure policy.  
  - **FAIL**: Terminate the execution of the current job.
  - **IGNORE**: Continue to execute the next node.
  - **SUSPEND**: Suspend the execution of the current job.  
  - **FAIL_CHILD**: Terminate the execution of the subsequent node.
  
  The default value is **FAIL**.

* `event_trigger` - (Optional, List) Event trigger for the real-time job node.
  The [event_trigger](#job_NodeEventTrigger) structure is documented below.

* `cron_trigger` - (Optional, List) Cron trigger for the real-time job node
  The [cron_trigger](#job_NodeCronTrigger) structure is documented below.

<a name="job_Location"></a>
The `location` block supports:

* `x` - (Required, Int) Position of the node on the horizontal axis of the job canvas.

* `y` - (Required, Int) Position of the node on the vertical axis of the job canvas.

<a name="job_Condition"></a>
The `conditions` block supports:

* `pre_node_name` - (Required, String) Name of the previous node on which the current node depends.

* `expression` - (Required, String) EL expression.  
  If the calculation result of the EL expression is true, this node is executed.

<a name="job_NodeProperty"></a>
The `properties` block supports:

* `name` - (Optional, String) Property name.

* `value` - (Optional, String) Property value.

<a name="job_NodeEventTrigger"></a>
The `event_trigger` block supports:

* `event_type` - (Required, String) Event type.  
  The valid values are as follows:
    - **KAFKA**: Select the corresponding connection name and topic. When a new Kafka message is
        received, the job is triggered.
    - **OBS**: Select the OBS path to be listened to. If new files exist in the path, scheduling is
       triggered. The path name can be referenced using variable Job.trigger.obsNewFiles. The
        prerequisite is that DIS notifications have been configured for the OBS path.
    - **DIS**: Currently, only newly reported data events from the DIS stream can be monitored.
        Each time a data record is reported, the job runs once.

* `channel` - (Required, String) DIS stream name.  
  Perform the following operations to obtain the stream name:  
    - Log in to the management console.  
    - Click **Data Ingestion Service** and select **Stream Management** from the left navigation pane.  
    - The stream management page lists the existing streams.

* `fail_policy` - (Optional, String) Job failure policy.  
  The valid values are as follows:
    - **SUSPEND**: Suspend the event.
    - **IGNORE**: Ignore the failure and process with the next event.

    The default value is **SUSPEND**.

* `concurrent` - (Optional, Int) Number of the concurrently scheduled jobs.  
  Value range: `1` to `128`.  
  Default value: `1`.

* `read_policy` - (Optional, String) Access policy.  
  The valid values are as follows:
    - **LAST**: Access data from the last location.
    - **NEW**: Access data from a new location.  

  The default value is **LAST**.

<a name="job_NodeCronTrigger"></a>
The `cron_trigger` block supports:

* `start_time` - (Required, String) Scheduling start time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**,
  which is an ISO 8601 time format.  
  For example, 2018-10-22T23:59:59+08, which indicates that a job starts to be scheduled at 23:59:59
  on October 22nd, 2018.

* `end_time` - (Optional, String) Scheduling end time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**,
  which is an ISO 8601 time format.  
  For example, 2018-10-22T23:59:59+08, which indicates that a job stops to be scheduled at 23:59:59
  on October 22nd, 2018.  
  If the end time is not set, the job will continuously be executed based on the scheduling period.

* `expression` - (Required, String) Cron expression in the format of **`<second><minute><hour><day><month><week>`**.

* `expression_time_zone` - (Optional, String) Time zone corresponding to the Cron expression.  
  Default value: time zone where DataArts Studio is located

* `period` - (Optional, String) Job execution interval consisting of a time and time unit.  
  Example: 1 hours, 1 days, 1 weeks, 1 months.  
  The value must match the value of expression.

* `depend_pre_period` - (Optional, Bool) Indicates whether to depend on the execution result of the current
  job's dependent job in the previous scheduling period.  
  Default value: **false**.

* `depend_jobs` - (Optional, List) Job dependency configuration.
  The [depend_jobs](#job_DependJobs) structure is documented below.

* `concurrent` - (Optional, Int) Number of concurrent executions allowed.

<a name="job_DependJobs"></a>
The `depend_jobs` block supports:

* `jobs` - (Required, List) A list of dependent jobs. Only the existing jobs can be depended on.

* `depend_period` - (Optional, String) Dependency period.  
  The valid values are as follows:
    - **SAME_PERIOD**: To run a job or not depends on the execution result of its depended job in
        the current scheduling period.
    - **PRE_PERIOD**: To run a job or not depends on the execution result of its depended job in
        the previous scheduling period.  

  The default value is **SAME_PERIOD**.

* `depend_fail_policy` - (Optional, String) Dependency job failure policy.  
  The valid values are as follows:
    - **FAIL**: Stop the job and set the job to the failed state.
    - **IGNORE**: Continue to run the job.
    - **SUSPEND**: Suspend the job.  

  The default value is **FAIL**.

<a name="job_Schedule"></a>
The `schedule` block supports:

* `type` - (Required, String) Scheduling type.  
  - **EXECUTE_ONCE**: The job runs immediately and runs only once.
  - **CRON**: The job runs periodically.
  - **EVENT**: The job is triggered by events.

* `cron` - (Optional, List) When `type` is set to **CRON**, configure the scheduling frequency and start time.
  The [cron](#job_ScheduleCron) structure is documented below.

* `event` - (Optional, List) When `type` is set to **EVENT**, configure information such as the event source.
  The [Event](#job_ScheduleEvent) structure is documented below.

<a name="job_ScheduleCron"></a>
The `cron` block supports:

* `start_time` - (Required, String) Scheduling start time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**,
  which is an ISO 8601 time format.  
  For example, 2018-10-22T23:59:59+08, which indicates that a job starts to be scheduled at 23:59:59
  on October 22nd, 2018.

* `end_time` - (Optional, String) Scheduling end time in the format of **yyyy-MM-dd'T'HH:mm:ssZ**,
  which is an ISO 8601 time format.  
  For example, 2018-10-22T23:59:59+08, which indicates that a job stops to be scheduled at 23:59:59
  on October 22nd, 2018.  
  If the end time is not set, the job will continuously be executed based on the scheduling period.

* `expression` - (Required, String) Cron expression in the format of **`<second><minute><hour><day><month><week>`**.

* `expression_time_zone` - (Optional, String) Time zone corresponding to the Cron expression.  
  Default value: time zone where DataArts Studio is located

* `depend_pre_period` - (Optional, Bool) Indicates whether to depend on the execution result of the current
  job's dependent job in the previous scheduling period.  
  Default value: **false**.

* `depend_jobs` - (Optional, List) Job dependency configuration.
  The [depend_jobs](#job_ScheduleCronDependJobs) structure is documented below.

<a name="job_ScheduleCronDependJobs"></a>
The `depend_jobs` block supports:

* `jobs` - (Required, List) A list of dependent jobs. Only the existing jobs can be depended on.

* `depend_period` - (Optional, String) Dependency period.  
  The valid values are as follows:
    - **SAME_PERIOD**: To run a job or not depends on the execution result of its depended job in
        the current scheduling period.
    - **PRE_PERIOD**: To run a job or not depends on the execution result of its depended job in
        the previous scheduling period.  

  The default value is **SAME_PERIOD**.

* `depend_fail_policy` - (Optional, String) Dependency job failure policy.  
  The valid values are as follows:
    - **FAIL**: Stop the job and set the job to the failed state.
    - **IGNORE**: Continue to run the job.
    - **SUSPEND**: Suspend the job.  

  The default value is **FAIL**.

<a name="job_ScheduleEvent"></a>
The `event` block supports:

* `event_type` - (Required, String) Event type.  
  The valid values are as follows:
    - **KAFKA**: Select the corresponding connection name and topic. When a new Kafka message is
        received, the job is triggered.
    - **OBS**: Select the OBS path to be listened to. If new files exist in the path, scheduling is
       triggered. The path name can be referenced using variable Job.trigger.obsNewFiles. The
        prerequisite is that DIS notifications have been configured for the OBS path.
    - **DIS**: Currently, only newly reported data events from the DIS stream can be monitored.
        Each time a data record is reported, the job runs once.

* `channel` - (Required, String) DIS stream name.  
  Perform the following operations to obtain the stream name:  
    - Log in to the management console.  
    - Click **Data Ingestion Service** and select **Stream Management** from the left navigation pane.  
    - The stream management page lists the existing streams.

* `fail_policy` - (Optional, String) Job failure policy.  
  The valid values are as follows:
    - **SUSPEND**: Suspend the event.
    - **IGNORE**: Ignore the failure and process with the next event.

    The default value is **SUSPEND**.

* `concurrent` - (Optional, Int) Number of the concurrently scheduled jobs.  
  Value range: `1` to `128`.  
  Default value: `1`.

* `read_policy` - (Optional, String) Access policy.  
  The valid values are as follows:
    - **LAST**: Access data from the last location.
    - **NEW**: Access data from a new location.  

  The default value is **LAST**.

<a name="job_Param"></a>
The `params` block supports:

* `name` - (Required, String) Name of a parameter.  
  The name contains a maximum of 64 characters, including only letters, numbers, hyphens (-), and
  underscores (_).

* `value` - (Required, String) Value of the parameter.  
  It cannot exceed 1024 characters.

* `type` - (Optional, String) Parameter type.  
  The valid values are as follows:
    - **variable**
    - **constants**

  Defaults to **variable**.

<a name="job_BasicConfig"></a>
The `basic_config` block supports:

* `owner` - (Optional, String) Job owner.

* `priority` - (Optional, Int) Job priority.  
  The value ranges from `0` to `2`.  
  `0` indicates a top priority, `1` indicates a medium priority, and `2` indicates a low priority.
  Default value: `0`.

* `execute_user` - (Optional, String) Job execution user. The value must be an existing user.

* `instance_timeout` - (Optional, Int) Maximum execution time of a job instance.  
  Unit: minute; value range: `5` to `1440`.  
  Default value: `60`.

* `custom_fields` - (Optional, Map) Custom fields.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `name`.

## Import

The job can be imported using `workspace_id`, `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_factory_job.test <workspace_id>/<name>
```
