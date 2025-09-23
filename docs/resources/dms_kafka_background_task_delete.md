---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_background_task_delete"
description: |-
  Manages a DMS kafka background task delete resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_background_task_delete

Manages a DMS kafka background task delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "task_id" {}

resource "huaweicloud_dms_kafka_background_task_delete" "test" {
  instance_id = var.instance_id
  task_id     = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `task_id` - (Required, String, ForceNew) Specifies the task ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
