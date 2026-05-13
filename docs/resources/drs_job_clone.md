---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_job_clone"
description: |-
  Manages a DRS job clone resource within HuaweiCloud.
---

# huaweicloud_drs_job_clone

Manages a DRS job clone resource within HuaweiCloud.

-> 1. This resource is a one-time action resource used to clone a DRS job. Deleting this resource will not delete the
  cloned job from the cloud, but will only remove the resource information from the tf state file.
  <br/>2. Cloning is not supported for tasks in the following states: creating, creation failed, configuring, waiting
  to start, starting, or deleting.

## Example Usage

```hcl
variable "job_id" {} 
variable "name" {}

resource "huaweicloud_drs_job_clone" "test" { 
  job_id = var.job_id 
  name   = var.cname
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the source DRS job to be cloned.

* `name` - (Required, String, ForceNew) Specifies the name of the cloned DRS job.
  The name must be between 4 and 50 characters long, start with a letter, and can contain letters, digits,
  hyphens (-), or underscores (_). It cannot contain other special characters and must be unique.
  Changing this creates a new resource.

* `task_version` - (Optional, String, NonUpdatable) Specifies the task version. For new UX tasks, the value is **2.0**.
  The default value is empty, and this parameter can be ignored.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the cloned job.

* `is_clone_job` - The indication of whether the job is a cloned job.

* `create_time` - The creation time of the cloned job.
