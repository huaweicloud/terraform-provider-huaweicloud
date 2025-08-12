---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_virtual_interface_switchover"
description: |-
  Manages a DC virtual interface switchover resource within HuaweiCloud.
---

# huaweicloud_dc_virtual_interface_switchover

Manages a DC virtual interface switchover resource within HuaweiCloud.

## Example Usage

```hcl
variable "resource_id" {}

resource "huaweicloud_dc_virtual_interface_switchover" "test" {
  resource_id = var.resource_id
  operation   = "shutdown"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the virtual gateway is located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `resource_id` - (Required, String, NonUpdatable) Specifies the ID of the resource on which the switchover test is to
  be performed.

* `operation` - (Required, String, NonUpdatable) Specifies whether to perform a switchover test.
  Value options: **shutdown** and **undo_shutdown**.

* `resource_type` - (Optional, String, NonUpdatable) Specifies the type of the resource on which the switchover test is
  to be performed. Defaults to **virtual_interface**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `start_time` - Indicates the start time of the switchover test.

* `end_time` - Indicates the end time of the switchover test.

* `operate_status` - Indicates the status of the switchover test.
