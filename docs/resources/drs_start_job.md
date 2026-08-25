---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_start_job"
description: |-
  Manages a resource to start a DRS job within HuaweiCloud.
---

# huaweicloud_drs_start_job

Manages a resource to start a DRS job within HuaweiCloud.

-> 1. This resource is a one-time action resource used to start a DRS job. Deleting this resource will not stop the
  job or undo the start action, but will only remove the resource information from the tf state file.
  <br/>2. The resource performs a pre-check before starting the job.
  <br/>3. You must specify an existing DRS job ID. If the job does not exist or is not in a pending start status,
  the operation may fail.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_start_job" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `start_time` - (Optional, String, NonUpdatable) Specifies the time to start the job. The time format is a timestamp
  precise to the millisecond, e.g. **1684466549755**, which indicates **2023-05-19 11:22:29.755**.
  Start immediately by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `job_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
