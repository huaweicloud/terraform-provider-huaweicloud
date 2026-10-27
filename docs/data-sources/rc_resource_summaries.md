---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_summaries"
description: |-
  Use this data source to get the resource summaries within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_summaries

Use this data source to get the resource summaries within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_resource_summaries" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the ResourceQL name.

* `type` - (Optional, List) Specifies the resource type (provider type).

* `region_id` - (Optional, List) Specifies the region ID list.

* `ep_id` - (Optional, List) Specifies the resource enterprise project ID list.

* `tags` - (Optional, List) Specifies the resource tag list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `summaries` - Resource summaries grouped by provider, type and region.

  The [summaries](#summaries_struct) structure is documented below.

<a name="summaries_struct"></a>
The `summaries` block supports:

* `provider` - The cloud service name.

* `types` - The resource type list.

  The [types](#types_struct) structure is documented below.

<a name="types_struct"></a>
The `types` block supports:

* `type` - The resource type name.

* `regions` - The region list.

  The [regions](#regions_struct) structure is documented below.

<a name="regions_struct"></a>
The `regions` block supports:

* `region_id` - The region ID.

* `count` - The number of this resource type in the current region.
