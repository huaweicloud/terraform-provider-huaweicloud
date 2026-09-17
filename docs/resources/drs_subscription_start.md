---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscription_start"
description: |-
  Use this resource to start a DRS subscription task within HuaweiCloud.
---

# huaweicloud_drs_subscription_start

Use this resource to start a DRS subscription task within HuaweiCloud.

-> This resource is a one-time action resource for starting a DRS subscription task. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_subscription_start" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to start the subscription task.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS subscription task to start.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the task ID.
