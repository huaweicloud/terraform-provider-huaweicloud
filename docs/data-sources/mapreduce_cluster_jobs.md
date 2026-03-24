---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_jobs"
description: |-
  Use this data source to query the job list under the specified cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_jobs

Use this data source to query the job list under the specified cluster within HuaweiCloud.

## Example Usage

### Query all jobs under the specified cluster

```hcl
variable "cluster_id" {}

data "huaweicloud_mapreduce_cluster_jobs" "test" {
  cluster_id = var.cluster_id
}
```

### Query job by job ID

```hcl
variable "cluster_id" {}
variable "job_id" {}

data "huaweicloud_mapreduce_cluster_jobs" "test" {
  cluster_id = var.cluster_id
  job_id     = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster jobs are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

* `job_id` - (Optional, String) Specifies the ID of the job.

* `job_name` - (Optional, String) Specifies the name of the job.  
  Fuzzy search is supported.

* `user` - (Optional, String) Specifies the user name of the job submitter.

* `job_type` - (Optional, String) Specifies the type of the job.

* `job_state` - (Optional, String) Specifies the execution status of the job.

* `job_result` - (Optional, String) Specifies the execution result of the job.

* `queue` - (Optional, String) Specifies the resource queue name of the job.

* `submitted_time_begin` - (Optional, String) Specifies the begin time of the submitted jobs, in RFC3339 format.

* `submitted_time_end` - (Optional, String) Specifies the end time of the submitted jobs, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of jobs that match the filter parameters.  
  The [jobs](#cluster_jobs) structure is documented below.

<a name="cluster_jobs"></a>
The `jobs` block supports:

* `job_id` - The ID of the job.

* `user` - The user name of the job submitter.

* `job_name` - The name of the job.

* `job_result` - The execution result of the job.

* `job_state` - The execution status of the job.

* `job_progress` - The execution progress of the job.

* `job_type` - The type of the job.

* `started_time` - The start time of the job, in RFC3339 format.

* `submitted_time` - The submit time of the job, in RFC3339 format.

* `finished_time` - The finish time of the job, in RFC3339 format.

* `elapsed_time` - The elapsed time of the job, in milliseconds.

* `arguments` - The runtime arguments of the job.

* `launcher_id` - The launcher ID of the job.

* `properties` - The properties of the job.

* `app_id` - The application ID of the job.

* `tracking_url` - The tracking URL of the job logs.

* `queue` - The resource queue name of the job.
