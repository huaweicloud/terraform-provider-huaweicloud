---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_count"
description: |-
  Use this data source to get the count of resources within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_count

Use this data source to get the count of resources within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_resource_count" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `name` - (Optional, String) Specifies the ResourceQL name.

* `type` - (Optional, List) Specifies the resource type (provider type).

* `region_id` - (Optional, List) Specifies the region ID list.

* `ep_id` - (Optional, List) Specifies the enterprise project ID list.

* `tags` - (Optional, List) Specifies the resource tag list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_count` - Total number of resources.
