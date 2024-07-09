---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_function_events"
description: |-
  Use this data source to get the list of function events within HuaweiCloud.
---

# huaweicloud_fgs_function_events

Use this data source to get the list of function events within HuaweiCloud.

## Example Usage

```hcl
variable "function_urn" {}

data "huaweicloud_fgs_function_events" "test" {
  function_urn = var.function_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `function_urn` - (Required, String) Specifies the function URN to which the events belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - All events that match the filter parameters.

  The [events](#events_struct) structure is documented below.

<a name="events_struct"></a>
The `events` block supports:

* `id` - The event ID.

* `name` - The event name.

* `updated_at` - The latest update time of the function event, in RFC3339 format.
