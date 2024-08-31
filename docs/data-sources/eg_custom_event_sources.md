---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_custom_event_sources"
description: ""
---

# huaweicloud_eg_custom_event_sources

Use this data source to filter EG custom event sources within HuaweiCloud.

## Example Usage

```hcl
variable "source_name" {}

data "huaweicloud_eg_custom_event_sources" "test" {
  name = var.source_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the custom event sources are located.  
  If omitted, the provider-level region will be used.

* `channel_id` - (Optional, String) Specifies the ID of the custom event channel to which the custom event sources
  belong.

* `source_id` - (Optional, String) Specifies the event source ID used to query specified custom event source.

* `name` - (Optional, String) Specifies the event source name used to query specified custom event source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sources` - The filtered custom event source.
  The [object](#eg_custom_event_sources) structure is documented below.

<a name="eg_custom_event_sources"></a>
The `sources` block supports:

* `id` - The ID of the custom event source.

* `channel_id` - The ID of the custom event channel to which the custom event source belong.

* `channel_name` - The name of the custom event channel to which the custom event source belong.

* `name` - The name of the custom event source.

* `type` - The type of the custom event source.

* `description` - The description of the custom event source.

* `status` - The status of the custom event source.
  + **CREATE_FAILED**
  + **RUNNING**
  + **ERROR**

* `created_at` - The creation time of the custom event source.

* `updated_at` - The update time of the custom event source.
