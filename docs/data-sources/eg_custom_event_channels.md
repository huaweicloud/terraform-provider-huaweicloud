---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_custom_event_channels"
description: ""
---

# huaweicloud_eg_custom_event_channels

Use this data source to filter EG custom event channels within HuaweiCloud.

## Example Usage

```hcl
variable "channel_name" {}

data "huaweicloud_eg_custom_event_channels" "test" {
  name = var.channel_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the custom event channels are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the channel name used to query specified custom event channel.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the custom event
  channels belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `channels` - The filtered custom event channels.
  The [channels](#eg_custom_event_channels) structure is documented below.

<a name="eg_custom_event_channels"></a>
The `channels` block supports:

* `id` - The ID of the custom event channel.

* `name` - The name of the custom event channel.

* `description` - The description of the custom event channel.

* `provider_type` - The type of the custom event channel.

* `enterprise_project_id` - The ID of the enterprise project to which the custom event channel belongs.

* `cross_account_ids` - The list of domain IDs (other tenants) for the cross-account policy.

* `created_at` - The creation time of the custom event channel.

* `updated_at` - The latest update time of the custom event channel.
