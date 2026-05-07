---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_agent_jobs"
description: |-
  Use this data source to query the agent jobs of RDS within HuaweiCloud.
---

# huaweicloud_rds_agent_jobs

Use this data source to query the agent jobs of RDS within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_db_agent_jobs" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `job_type` - (Optional, String) Specifies the type of the job.
  The valid values are as follows: **replication**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of DB agent jobs.
  The [jobs](#rds_db_agent_jobs_attr) structure is documented below.

<a name="rds_db_agent_jobs_attr"></a>
The `jobs` block supports:

* `job_id` - The job ID.

* `job_name` - The job name.

* `is_enabled` - Whether the job is enabled.

* `run_time` - The latest execution time, in the format **yyyy-mm-ddThh:mm:ssZ**.

* `run_status` - The execution status of the job.
  The valid values are as follows:
  + **failed**
  + **succeeded**
  + **retrying**
  + **canceled**
  + **in_progress**

* `last_failure_time` - The latest failure time, in the format **yyyy-mm-ddThh:mm:ssZ**.

* `failure_count` - The historical failure count.

* `agent_type` - The type of the job agent.
  The valid values are as follows:
  + **snapshot**
  + **log_reader**
  + **distribution**
  + **merge**
  + **queue_reader**

* `profile_id` - The profile ID. This parameter is valid when `job_type` is **replication**.

* `profile_name` - The profile name. This parameter is valid when `job_type` is **replication**.
