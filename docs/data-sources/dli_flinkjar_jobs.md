---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinkjar_jobs"
description: |-
  Use this data source to get the list of the DLI flinkjar jobs.
---

# huaweicloud_dli_flinkjar_jobs

Use this data source to get the list of the DLI flinkjar jobs.

## Example Usage

```hcl
variable "queue_name" {}

data "huaweicloud_dli_flinkjar_jobs" "test" {
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

* `status` - The status of the job.

* `description` - The description of the job.

* `queue_name` - The name of DLI queue which this job run in.

* `main_class` - The jar package main class.

* `entrypoint_args` - The job entry arguments.

* `flink_version` - The flink version.

* `obs_bucket` - The name of OBS bucket.

* `log_enabled` - The whether to enable the function of uploading job logs to users' OBS buckets.

* `smn_topic` - The SMN topic. If a job fails, the system will send a message to users subscribed to the SMN topic.

* `manager_cu_num` - The number of CUs in the job manager selected for a job.

* `cu_num` - The number of CUs selected for a job.

* `parallel_num` - The number of parallel for a job.

* `restart_when_exception` - The whether to enable the function of restart upon exceptions.

* `entrypoint` - The JAR file where the job main class is located.
  It is the name of the package that has been uploaded to the DLI.

* `dependency_jars` - The other dependency jars. It is the name of the package that has been uploaded to the DLI.

* `dependency_files` - The dependency files. It is the name of the package that has been uploaded to the DLI.

* `resume_checkpoint` - The whether the abnormal restart is recovered from the checkpoint.

* `runtime_config` - The customize optimization parameters during flink job runtime.

* `resume_max_num` - The maximum number of retry times upon exceptions.

* `checkpoint_path` - The checkpoint save path.

* `tm_cu_num` - The number of CUs occupied by a single TM.

* `tm_slot_num` - The number of single TM slots.

* `image` - The custom image.

* `feature` - The custom job features. Type of the Flink image used by a job.
    + **basic**: indicates that the basic Flink image provided by DLI is used.
    + **custom**: indicates that the user-defined Flink image is used.
