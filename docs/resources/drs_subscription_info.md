---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_subscription_info"
description: |-
  Use this resource to update the information of a DRS subscription task within HuaweiCloud.
---

# huaweicloud_drs_subscription_info

Use this resource to update the information of a DRS subscription task within HuaweiCloud.

-> This resource is a one-time action resource for updating the information of a DRS subscription task. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_subscription_info" "test" {
  job_id      = var.job_id
  name        = "DRS-123123123"
  description = "TEST"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to update the subscription task info.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS subscription task.

* `name` - (Optional, String, NonUpdatable) Specifies the name of the subscription task.  
  The name must be between 4 and 50 characters, and can contain letters, digits, hyphens or underscores.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the subscription task.

* `consume_time` - (Optional, Int, NonUpdatable) Specifies the consumption time point, in timestamp format.  
  After modifying the consumption time point, the incremental data pulled will start from the modified time point.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the task ID.
