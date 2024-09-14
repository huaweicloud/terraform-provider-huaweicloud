---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinksql_job"
description: ""
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

* `name` - (Required, String) Specifies the name of the job. Length range: `1` to `57` characters.
 which may consist of letters, digits, underscores (_) and hyphens (-).

* `type` - (Optional, String, ForceNew) Specifies the type of the job. The valid values are `flink_sql_job`,
 `flink_opensource_sql_job` and `flink_sql_edge_job`. Default value is `flink_sql_job`.
  Changing this parameter will create a new resource.

* `run_mode` - (Optional, String) Specifies job running mode. The options are as follows:

  + **shared_cluster**: indicates that the job is running on a shared cluster.
  + **exclusive_cluster**: indicates that the job is running on an exclusive cluster.
  + **edge_node**: indicates that the job is running on an edge node.
  
  The default value is `shared_cluster`.

* `description` - (Optional, String) Specifies job description. Length range: `1` to `512` characters.

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

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the resource.

* `flink_version` - (Optional, String) Specifies the version of the flink.
  The valid values are `1.10` and `1.12`, defalut value is `1.10`.
  
* `operator_config` - (Optional, String) Specifies degree of parallelism (DOP) configuration of an operator, in
  JSON format.

* `static_estimator` - (Optional, Bool) Specifies whether to estimate static resources. Default value is `false`.

* `static_estimator_config` - (Optional, String) Specifies the traffic or hit rate configuration of each operator, in
  JSON format.

* `graph_type` - (Optional, String) Specifies the type of stream graph to be generated by the Flink SQL job.
  The valid values are `simple_graph` and `job_graph`. The default value is `simple_graph`.

  -> When `type` is set to `flink_opensource_sql_job`, the `flink_version`, `operator_config`, `static_estimator`,
     `static_estimator_config` and `graph_type` parameters are valid.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Job ID in Int format.

* `status` - The Job status.

* `stream_graph` - The simplified stream graph or static stream graph information of the Flink SQL job.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Clusters can be imported by their `id`. For example,

```bash
terraform import huaweicloud_dli_flinksql_job.test 12345
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `static_estimator`, `graph_type`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cae_component" "test" {
  ...
  lifecycle {
    ignore_changes = [
      static_estimator, graph_type,
    ]
  }
}
```
