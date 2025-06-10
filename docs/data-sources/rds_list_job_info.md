---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_list_job_info"
description: |-
  Use this data source to query a single asynchronous task (job) in the RDS task center by its job ID.
---

# huaweicloud_rds_list_job_info

Use this data source to query a single asynchronous task (job) in the RDS task center by its job ID.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_rds_list_job_info" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `job_id` - (Required, String) Specifies the identifier of the asynchronous task. This value is returned by the RDS
  API when an asynchronous task is created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `job` - Indicates the list of containing job information.

  The [job](#job_struct) structure is documented below.

<a name="job_struct"></a>
The `job` block contains:

* `id` - Indicates the job ID.

* `name` - Indicates the name of the job.

* `status` - Indicates the status of the job.

* `created` -  Indicates the creation time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.
  T is the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `ended` -  Indicates the end time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format. T is
  the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `process` - Indicates the task execution progress. The execution progress (such as "60",
  indicating the task execution progress is 60%) is displayed only when the task is being
  executed. Otherwise, "" is returned.

* `entities` - Indicates the job-specific information. The content varies depending on the job type.

* `fail_reason` - Indicates the failure reason for the job, if applicable.

* `instances` - Indicates the list containing instance information.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block contains:

* `id` - Indicates the ID of the RDS instance associated with the job.

* `name` - Indicates the name of the RDS instance associated with the job.
