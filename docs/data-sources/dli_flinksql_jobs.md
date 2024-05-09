---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinksql_jobs"
description: |-
  Use this data source to get the list of the DLI flinksql jobs.
---

# huaweicloud_dli_flinksql_jobs

Use this data source to get the list of the DLI flinksql jobs.

## Example Usage

```hcl
variable "queue_name" {}

data "huaweicloud_dli_flinksql_jobs" "test" {
  queue_name = var.queue_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `queue_name` - (Optional, String) Specifies the name of DLI queue which this job to be queried.

* `job_id` - (Optional, String) Specifies the ID of the job to be queried.

* `tags` - (Optional, Map) Specifies the key/value pairs to be queried.

* `manager_cu_num` - (Optional, Int) Specifies number of CUs in the job manager to be queried.

* `cu_num` - (Optional, Int) Specifies number of CUs to be queried.

* `parallel_num` - (Optional, Int) Specifies number of parallel to be queried.

* `tm_cu_num` - (Optional, Int) Specifies number of CUs occupied by a single TM to be queried.

* `tm_slot_num` - (Optional, Int) Specifies number of single TM slots to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - All jobs that match the filter parameters.

  The [jobs](#job_list_jobs_struct) structure is documented below.

<a name="job_list_jobs_struct"></a>
The `jobs` block supports:

* `id` - The ID of the job.

* `name` - The name of the job.

* `flink_version` - The version of the flink.

* `type` - The type of the job.

* `status` - The status of the job.

* `run_mode` - The job running mode.

* `description` - The description of the job.

* `queue_name` - The name of DLI queue which this job run in.

* `sql` - The stream SQL statement.

* `cu_num` - The number of CUs selected for a job.

* `parallel_num` - The number of parallel for a job.

* `checkpoint_enabled` - The whether to enable the automatic job snapshot function.

* `checkpoint_mode` - The snapshot mode.

* `checkpoint_interval` - The snapshot interval. The unit is second.

* `obs_bucket` - The name of OBS bucket.

* `log_enabled` - The whether to enable the function of uploading job logs to users' OBS buckets.

* `smn_topic` - The SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.

* `restart_when_exception` - The whether to enable the function of restart upon exceptions.

* `idle_state_retention` - The retention time of the idle state. The unit is hour.

* `edge_group_ids` - The edge computing group IDs.

* `dirty_data_strategy` - The dirty data policy of a job.

* `udf_jar_url` - The name of the resource package that has been uploaded to the
  DLI resource management system. The UDF Jar file of the SQL job is specified by this parameter.

* `manager_cu_num` - The number of CUs in the job manager selected for a job.

* `tm_cu_num` - The number of CUs occupied by a single TM.

* `tm_slot_num` - The number of single TM slots.

* `resume_checkpoint` - The whether the abnormal restart is recovered from the checkpoint.

* `resume_max_num` - The maximum number of retry times upon exceptions.

* `runtime_config` - The customize optimization parameters during flink job runtime.

* `operator_config` - The degree of parallelism configuration of an operator, in JSON format.

* `static_estimator_config` - The static flow chart resource estimation parameters, in JSON format.
