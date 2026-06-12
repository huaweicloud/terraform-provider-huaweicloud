---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_scheduled_task_cancel"
description: |-
  Manages a resource to cancel a scheduled job for GeminiDB within HuaweiCloud.
---

# huaweicloud_geminidb_scheduled_task_cancel

Manages a resource to cancel a scheduled job for GeminiDB within HuaweiCloud.

-> This resource is only a one-time action resource for stopping a backup of GeminDB.
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_geminidb_scheduled_task_cancel" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds instance resource. If omitted, the
  provider-level region will be used. Changing this creates a new geminidb instance resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the scheduled job ID to cancel. The job ID can be obtained
  from the scheduled job list API. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as the job ID.
