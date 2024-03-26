---
subcategory: "Data Lake Insight (DLI)"
---

# huaweicloud_dli_flinksql_job

Manages a flink sql job resource within HuaweiCloud DLI.

## Example Usage

### Create a flink job

```hcl
variable "sql" {}
variable "jobName" {}

resource "huaweicloud_dli_flinksql_job" "test" {
  name = var.jobName
  type = "flink_sql_job"
  sql  = var.sql
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DLI flink job resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the job. Length range: 1 to 57 characters.
 which may consist of letters, digits, underscores (_) and hyphens (-).

* `type` - (Optional, String, ForceNew) Specifies the type of the job. The valid values are `flink_sql_job`,
 `flink_opensource_sql_job` and `flink_sql_edge_job`. Default value is `flink_sql_job`.
  Changing this parameter will create a new resource.

* `run_mode` - (Optional, String) Specifies job running mode. The options are as follows:

  + **shared_cluster**: indicates that the job is running on a shared cluster.
  + **exclusive_cluster**: indicates that the job is running on an exclusive cluster.
  + **edge_node**: indicates that the job is running on an edge node.
  
  The default value is `shared_cluster`.

* `description` - (Optional, String) Specifies job description. Length range: 1 to 512 characters.

* `queue_name` - (Optional, String) Specifies name of a queue.
  If you want to use the parameters, the `run_mode` parameter must be set to `exclusive_cluster`.

* `sql` - (Optional, String) Specifies stream SQL statement, which includes at least the following
 three parts: source, query, and sink. Length range: 1024x1024 characters.

* `cu_number` - (Optional, Int) Specifies number of CUs selected for a job. The default value is 2.

* `parallel_number` - (Optional, Int) Specifies number of parallel for a job. The default value is 1.

* `checkpoint_enabled` - (Optional, Bool) Specifies whether to enable the automatic job snapshot function.
  + **true**: indicates to enable the automatic job snapshot function.
  + **false**: indicates to disable the automatic job snapshot function.

  Default value: false

* `checkpoint_mode` - (Optional, String) Specifies snapshot mode. There are two options:
  + **exactly_once**: indicates that data is processed only once.
  + **at_least_once**: indicates that data is processed at least once.

  The default value is `exactly_once`.

* `checkpoint_interval` - (Optional, Int) Specifies snapshot interval. The unit is second.
  The default value is 10.

* `obs_bucket` - (Optional, String) Specifies OBS path. OBS path where users are authorized to save the
  snapshot. This parameter is valid only when `checkpoint_enabled` is set to `true`. OBS path where users are authorized
  to save the snapshot. This parameter is valid only when `log_enabled` is set to `true`.

* `log_enabled` - (Optional, Bool) Specifies whether to enable the function of uploading job logs to
  users' OBS buckets. The default value is false.
  
* `smn_topic` - (Optional, String) Specifies SMN topic. If a job fails, the system will send a message to
 users subscribed to the SMN topic.
  
* `restart_when_exception` - (Optional, Bool) Specifies whether to enable the function of automatically
 restarting a job upon job exceptions. The default value is false.
  
* `idle_state_retention` - (Optional, Int) Specifies retention time of the idle state. The unit is hour.
 The default value is 1.

* `edge_group_ids` - (Optional, List) Specifies edge computing group IDs.
  
* `dirty_data_strategy` - (Optional, String) Specifies dirty data policy of a job.
  + **2:obsDir**: Save the dirty data to the obs path `obsDir`. For example: `2:yourBucket/output_path`
  + **1**: Trigger a job exception
  + **0**: Ignore

  The default value is `0`.
  
* `udf_jar_url` - (Optional, String) Specifies name of the resource package that has been uploaded to the
  DLI resource management system. The UDF Jar file of the SQL job is specified by this parameter.
  
* `manager_cu_number` - (Optional, Int) Specifies number of CUs in the JobManager selected for a job.
 The default value is 1.
  
* `tm_cus` - (Optional, Int) Specifies number of CUs for each Task Manager. The default value is 1.
  
* `tm_slot_num` - (Optional, Int) Specifies number of slots in each Task Manager.
 The default value is (**parallel_number** * **tm_cus**)/(**cu_number** - **manager_cu_number**).
  
* `resume_checkpoint` - (Optional, Bool) Specifies whether the abnormal restart is recovered from the
 checkpoint.
  
* `resume_max_num` - (Optional, Int) Specifies maximum number of retry times upon exceptions. The unit is
 `times/hour`. Value range: `-1` or greater than `0`. The default value is `-1`, indicating that the number of times is
 unlimited.

* `runtime_config` - (Optional, Map) Specifies customizes optimization parameters when a Flink job is
 running.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Job ID in Int format.

* `status` - The Job status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Clusters can be imported by their `id`. For example,

```
terraform import huaweicloud_dli_flinksql_job.test 12345
```
