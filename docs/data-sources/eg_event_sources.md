---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_sources"
description: |-
  Use this data source to filter EG event sources within HuaweiCloud.
---

# huaweicloud_eg_event_sources

Use this data source to filter EG event sources within HuaweiCloud.

## Example Usage

### Query all kinds of event sources

```hcl
data "huaweicloud_eg_event_sources" "test" {}
```

### Query all official event sources

```hcl
data "huaweicloud_eg_event_sources" "test" {
  provider_type = "OFFICIAL"
}
```

### Query the custom event source with the specified name

```hcl
variable "source_name" {}

data "huaweicloud_eg_event_sources" "test" {
  provider_type = "CUSTOM"
  name          = var.source_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event sources are located.  
  If omitted, the provider-level region will be used.

* `provider_type` - (Optional, String) Specifies the type of the event sources to be queried.  
  The valid values are as follows:
  + **OFFICIAL**
  + **CUSTOM**
  + **PARTNER**

* `channel_id` - (Optional, String) Specifies the ID of the event channel to which the event sources belong.

* `name` - (Optional, String) Specifies the name of the event source to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sources` - All event sources that match the filter parameters.
  The [sources](#eg_event_sources) structure is documented below.

<a name="eg_event_sources"></a>
The `sources` block supports:

* `id` - The ID of the event source.

* `channel_id` - The ID of the event channel to which the event source belong.

* `channel_name` - The name of the event channel to which the event source belong.

* `name` - The name of the event source.

* `label` - The display name of the official event source.

* `provider_type` - The provider type of the event source.
  + **OFFICIAL**
  + **CUSTOM**
  + **PARTNER**

* `event_types` - The event types that official event source provided.
The [event_types](#eg_source_event_types) structure is documented below.

* `type` - The type of the event source.

* `description` - The description of the event source.

* `status` - The status of the event source.
  + **CREATE_FAILED**
  + **RUNNING**
  + **ERROR**

* `created_at` - The creation time of the event source.

* `updated_at` - The update time of the event source.

* `detail` - The message instance link information encapsulated in json format.

<a name="eg_source_event_types"></a>
The `event_types` block supports:

* `name` - The name of the event type.

* `description` - The description of the event type.
