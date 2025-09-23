---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_alarm_subscription"
description: |-
  Manages a GaussDB(DWS) alarm subscription resource within HuaweiCloud.
---

# huaweicloud_dws_alarm_subscription

Manages a GaussDB(DWS) alarm subscription resource within HuaweiCloud.

## Example Usage

```hcl
variable "smn_urn" {}
variable "smn_name" {}

resource "huaweicloud_dws_alarm_subscription" "test" {
  name                     = "demo"
  enable                   = 1
  notification_target      = var.smn_urn
  notification_target_name = var.smn_name
  notification_target_type = "SMN"
  alarm_level              = "urgent,important"
  time_zone                = "GMT+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the alarm subscription.

* `enable` - (Required, Int) Whether the alarm subscription is enabled.  
  The options are as follows:
    + **1**: enable.
    + **0**: closed.

* `notification_target` - (Required, String) The notification target.  

* `notification_target_name` - (Required, String) The name of notification target.  

* `notification_target_type` - (Required, String) The type of notification target. Currently only **SMN** is supported.

* `time_zone` - (Required, String, ForceNew) The time_zone of the alarm subscription.  

  Changing this parameter will create a new resource.

* `alarm_level` - (Optional, String) The level of alarm. separate multiple alarm levels with commas (,).
  The valid values are **urgent**, **important**, **minor**, and **prompt**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The DWS alarm subscription can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dws_alarm_subscription.test <id>
```
