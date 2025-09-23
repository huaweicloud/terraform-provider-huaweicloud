---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_classes"
description: |-
  Use this data source to get the list of CC bandwidth classes.
---

# huaweicloud_cc_bandwidth_package_classes

Use this data source to get the list of CC bandwidth classes.

## Example Usage

```hcl
data "huaweicloud_cc_bandwidth_package_classes" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_levels` - Indicates the list of bandwidth package classes.

  The [bandwidth_package_levels](#bandwidth_package_levels_struct) structure is documented below.

<a name="bandwidth_package_levels_struct"></a>
The `bandwidth_package_levels` block supports:

* `id` - Indicates the instance ID.

* `level` - Indicates the bandwidth package class.

* `name_cn` - Indicates the instance Chinese name.

* `name_en` - Indicates the instance English name.

* `display_priority` - Indicates the priority of the bandwidth package. A smaller value indicates a higher priority.
  The value can be:
  + **Platinum**: 1 to 50
  + **Gold**: 51 to 100
  + **Silver**: 101 to 150
  + **Other**: greater than 151

* `description` - Indicates the description.

* `created_at` - Indicates the creation time.
  The UTC time is in the **yyy-MM-ddTHH:mm:ss** format.

* `updated_at` - Indicates the update time.
  The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.
