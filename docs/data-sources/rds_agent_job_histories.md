---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_agent_job_histories"
description: |-
  Use this data source to query the execution history of an agent job of RDS within HuaweiCloud.
---

# huaweicloud_rds_agent_job_histories

Use this data source to query the execution history of an agent job of RDS within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "job_id" {}

data "huaweicloud_rds_agent_job_histories" "test" {
  instance_id = var.instance_id
  job_id      = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `job_id` - (Required, String) Specifies the ID of the agent job.

* `run_status` - (Optional, String) Specifies the execution status of the job to filter.
  The valid values are as follows:
  + **failed**
  + **succeeded**
  + **retrying**
  + **canceled**
  + **in_progress**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - The list of execution histories.
  The [histories](#rds_db_agent_job_histories_attr) structure is documented below.

<a name="rds_db_agent_job_histories_attr"></a>
The `histories` block supports:

* `history_id` - The history record ID.

* `run_status` - The execution status of the job.
  The valid values are as follows:
  + **failed**
  + **succeeded**
  + **retrying**
  + **canceled**
  + **in_progress**

* `run_time` - The latest execution time, in the format **yyyy-mm-ddThh:mm:ssZ**.

* `run_duration` - The job execution duration, in the format **hh:mm:ss**.

* `message` - The execution message.
