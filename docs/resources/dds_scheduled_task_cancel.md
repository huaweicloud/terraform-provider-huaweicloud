---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_scheduled_task_cancel"
description: |-
  Manages a DDS scheduled task cancel resource within HuaweiCloud.
---

# huaweicloud_dds_scheduled_task_cancel

Manages a DDS scheduled task cancel resource within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_dds_scheduled_task_cancel" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `job_id` - (Required, String, ForceNew) Specifies the task ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - Indicates the create time.

* `end_time` - Indicates the end time.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the instance status.

* `job_name` - Indicates the task name.

* `job_status` - Indicates the task execution status.

* `start_time` - Indicates the start time.
