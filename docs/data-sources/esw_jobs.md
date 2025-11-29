---
subcategory: "Enterprise Switch (ESW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_esw_jobs"
description: |-
  Use this data source to get the list of ESW jobs.
---

# huaweicloud_esw_jobs

Use this data source to get the list of ESW jobs.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_esw_jobs" "test" {
  resource_id = var.resource_id
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the jobs. If omitted, the provider-level region will be
  used.

* `resource_id` - (Required, String) Specifies the ID of the instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the list of jobs.
  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - Indicates the ID of the job.

* `name` - Indicates the name of the job.

* `status` - Indicates the status of the job.

* `begin_time` - Indicates the beginning time of the job.

* `end_time` - Indicates the end time of the job.

* `process` - Indicates the current progress of the job.

* `fail_reason` - Indicates the fail reason of the job.

* `resource_id` - Indicates the ID of the resource associated with the job.

* `resource_name` - Indicates the name ID of the resource associated with the job.

* `resource_type` - Indicates the type ID of the resource associated with the job.

* `project_id` - Indicates the project ID.
