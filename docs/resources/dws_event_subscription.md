---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_event_subscription"
description: |-
  Manages a GaussDB(DWS) event subscription resource within HuaweiCloud.  
---

# huaweicloud_dws_event_subscription

Manages a GaussDB(DWS) event subscription resource within HuaweiCloud.  

## Example Usage

```hcl
  variable "smn_urn" {}
  variable "smn_name" {}

  resource "huaweicloud_dws_event_subscription" "test" {
    name                     = "demo"
    enable                   = 1
    notification_target      = var.smn_urn
    notification_target_name = var.smn_name
    notification_target_type = "SMN"
    category                 = "management,monitor,security"
    severity                 = "normal,warning"
    source_type              = "cluster,backup,disaster-recovery"
    time_zone                = "GMT+08:00"
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the event subscription.

* `enable` - (Required, String) Whether the event subscription is enabled.  
  The options are as follows:
    + **1**: open.
    + **0**: closed.

* `notification_target` - (Required, String) The notification target.  

* `notification_target_name` - (Required, String) The name of notification target.  

* `notification_target_type` - (Required, String) The type of notification target. Currently only **SMN** is supported.

* `source_id` - (Optional, String) ID of source event.

* `source_type` - (Optional, String) The type of source event.  
  The valid values are **cluster**, **backup**, and **disaster-recovery**.

* `category` - (Optional, String) The category of source event.  
  The valid values are **management**, **monitor**, **security** and **system alarm**.

* `severity` - (Optional, String) The severity of source event.  
  The valid values are **normal**, and **warning**.

* `time_zone` - (Optional, String, ForceNew) The time_zone of the event subscription.  

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The DWS event subscription can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dws_event_subscription.test <id>
```
