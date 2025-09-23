---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_job_action"
description:|-
  Manages a job action resource of DataArts Factory within HuaweiCloud.
---

# huaweicloud_dataarts_factory_job_action

Manages a job action resource of DataArts Factory within HuaweiCloud.

-> Destroying resources does not change the current action status of the job.

## Example Usage

```hcl
variable "workspace_id" {}
variable "job_name" {}

resource "huaweicloud_dataarts_factory_job_action" "test" {
  workspace_id = var.workspace_id
  action       = "start"
  job_name     = var.job_name
  process_type = "BATCH"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Optional, String, ForceNew) Specified the ID of the workspace to which the job belongs.
  Changing this creates a new resource.
  If this parameter is not set, the default workspace is used by default.

* `job_name` - (Required, String, ForceNew) Specified the name of the job.
  Changing this creates a new resource.
  The name contains a maximum of  `128` characters, including only letters, numbers, hyphens (-),
  underscores (_), and periods (.).

* `process_type` - (Required, String, ForceNew) Specified the type of the job.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **REAL_TIME**: Real-time processing.
  + **BATCH**: Batch processing.

* `action` - (Required, String) Specified the action type of the job.  
  The valid values are as follows:
  + **start**
  + **stop**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `job_name`.

* `status` - The current status of the job.
  + **NORMAL**
  + **STOPPED**
  + **SCHEDULING**
  + **PAUSED**
  + **EXCEPTION**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
