---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_event_subscriptions"
description: |-
  Use this data source to get the list of event subscriptions.
---

# huaweicloud_dws_event_subscriptions

Use this data source to get the list of event subscriptions.

## Example Usage

```hcl
variable "subscription_name" {}

data "huaweicloud_dws_event_subscriptions" "test" {
  name = var.subscription_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the event subscription.

* `notification_target_name` - (Optional, String) Specifies the name of notification target.

* `enable` - (Optional, String) Specifies whether the event subscription is enabled.
  The options are as follows:
  + **1**: enabled.
  + **0**: disabled.

* `category` - (Optional, String) Specifies the category of source event.
  The valid values are **management**, **monitor**, **security** and **system alarm**.
  If there are multiple categories, separate by commas, e.g. **management,security**.

* `severity` - (Optional, String) Specifies the severity of source event.
  The valid values are **normal** and **warning**. If there are multiple severities, separate by commas,
  e.g. **normal,warning**.

* `source_type` - (Optional, String) Specifies the type of source event.
  The valid values are **cluster**, **backup** and **disaster-recovery**. If there are multiple types,
  separate by commas, e.g. **cluster,disaster-recovery**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `event_subscriptions` - The list of event subscriptions.
  The [event_subscriptions](#attrblock_event_subscriptions) structure is documented below.

<a name="attrblock_event_subscriptions"></a>
The `event_subscriptions` block supports:

* `id` - The ID of event subscription.

* `name` - The name of the event subscription.

* `category` - The category of source event.

* `enable` - Whether the event subscription is enabled.

* `name_space` - The name space of the event subscription.

* `notification_target` - The notification target.

* `notification_target_name` - The name of notification target.

* `notification_target_type` - The type of notification target.

* `project_id` - The project ID of the event subscription.

* `severity` - The severity of source event.

* `source_id` - The ID of source event.

* `source_type` - The type of source event.

* `language` - The language of the event subscription.

* `time_zone` - The time zone of the event subscription.
