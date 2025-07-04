---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_ies_availability_zones"
description: |-
  Use this data source to get IES availability zones list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_ies_availability_zones

Use this data source to get IES availability zones list of the Workspace APP within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_ies_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - The list of availability zones.

  The [availability_zones](#workspace_ies_availability_zones_azs) structure is documented below.

<a name="workspace_ies_availability_zones_azs"></a>
The `availability_zones` block supports:

* `availability_zone` - The ID of the availability zone, such as **cn-north-4a**.

* `display_name` - The display name of the availability zone.

* `i18n` - The internationalization information of the availability zone.

* `sold_out` - The sold out information for the availability zone.

  The [sold_out](#workspace_ies_availability_zones_sold_out) structure is documented below.

* `product_ids` - The list of custom supported product IDs for the availability zone.

* `visible` - Whether the availability zone is visible.

* `default_availability_zone` - Whether this is the default availability zone.

<a name="workspace_ies_availability_zones_sold_out"></a>
The `sold_out` block supports:

* `products` - The list of sold out product IDs.
