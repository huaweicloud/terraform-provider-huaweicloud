---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_sql_jobs"
description: |-
  Use this data source to get the list of the DLI SQL jobs.
---

# huaweicloud_dli_sql_jobs

Use this data source to get the list of the DLI SQL jobs.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_dli_sql_jobs" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Optional, String) Specifies the ID of the job to be queried.

* `type` - (Optional, String) Specifies the type of the jobs to be queried.

* `status` - (Optional, String) Specifies the status of the job to be queried.
  The valid values are **FINISHED**, **FAILED** and **CANCELED**.

* `queue_name` - (Optional, String) Specifies the queue name which the jobs to be submitted belong.

* `start_time` - (Optional, String) Specifies the time when a job is start to be queried.
  The format is `YYYY-MM-DDThh:mm:ss{timezone}`, e.g. `2024-01-01T08:00:00+08:00`.

* `end_time` - (Optional, String) Specifies the time when a job is end to be queried.
  The format is `YYYY-MM-DDThh:mm:ss{timezone}`, e.g. `2024-01-01T08:00:00+08:00`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - All jobs that match the filter parameters.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - The ID of job.

* `type` - The type of job.

* `status` - The status of the job.

* `queue_name` - The queue name which this job to be submitted belongs.

* `owner` - The user who submits the job.

* `database_name` - The database name where the table that records its operations is located.

* `sql` - The SQL statement is executed by the job.

* `start_time` - The time when a job is start, in RFC 3339 format.

* `end_time` - The time when a job is end, in RFC 3339 format.

* `duration` - The job running duration (unit: millisecond).

* `tags` - The key/value pairs to associate with the job.
