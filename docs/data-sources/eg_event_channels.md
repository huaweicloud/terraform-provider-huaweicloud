---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_channels"
description: |-
  Use this data source to filter EG event channels within HuaweiCloud.
---

# huaweicloud_eg_event_channels

Use this data source to filter EG event channels within HuaweiCloud.

## Example Usage

### Query all kinds of event channels

```hcl
data "huaweicloud_eg_event_channels" "test" {}
```

### Query all official event channels

```hcl
data "huaweicloud_eg_event_channels" "test" {
  provider_type = "OFFICIAL"
}
```

### Query the custom event channel with the specified name

```hcl
variable "channel_name" {}

data "huaweicloud_eg_event_channels" "test" {
  provider_type = "CUSTOM"
  name          = var.channel_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event channels are located.  
  If omitted, the provider-level region will be used.

* `provider_type` - (Optional, String) Specifies the type of the event channels to be queried.
  + **OFFICIAL**
  + **CUSTOM**
  + **PARTNER**

* `name` - (Optional, String) Specifies the channel name used to query specified event channel.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the event channels
  belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `channels` - All event channels that match the filter parameters.  
  The [channels](#eg_event_channels_attr) structure is documented below.

<a name="eg_event_channels_attr"></a>
The `channels` block supports:

* `id` - The ID of the event channel.

* `name` - The name of the event channel.

* `description` - The description of the event channel.

* `provider_type` - The type of the event channel.

* `enterprise_project_id` - The ID of the enterprise project to which the event channel belongs.

* `cross_account_ids` - The list of domain IDs (other tenants) for the cross-account policy.

* `created_at` - The creation time of the event channel.

* `updated_at` - The latest update time of the event channel.
