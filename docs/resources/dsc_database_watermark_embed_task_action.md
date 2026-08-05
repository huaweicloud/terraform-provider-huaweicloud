---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_database_watermark_embed_task_action"
description: |-
  Use this resource to operate a DSC database watermark embed task within HuaweiCloud.
---

# huaweicloud_dsc_database_watermark_embed_task_action

Use this resource to operate a DSC database watermark embed task within HuaweiCloud.

-> This resource is a one-time action resource used to operate a DSC database watermark embed task. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_dsc_database_watermark_embed_task_action" "test" {
  task_id = var.task_id
  action  = "START"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the database watermark embed task is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the ID of the database watermark embed task.

* `action` - (Required, String, NonUpdatable) Specifies the operation type of the database watermark embed task.  
  The valid values are as follows:
  + **ENABLE**
  + **DISABLE**
  + **START**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
