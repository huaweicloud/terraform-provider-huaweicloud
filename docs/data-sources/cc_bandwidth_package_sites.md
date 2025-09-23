---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_bandwidth_package_sites"
description: |-
  Use this data source to get the list of CC bandwidth sites.
---

# huaweicloud_cc_bandwidth_package_sites

Use this data source to get the list of CC bandwidth sites.

## Example Usage

```hcl
data "huaweicloud_cc_bandwidth_package_sites" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `site_code` - (Optional, String) Specifies the site code.

* `region_id` - (Optional, String) Specifies the region ID.

* `name` - (Optional, String) Specifies the name used for fuzzy search.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidth_package_sites` - Indicates the site list.

  The [bandwidth_package_sites](#bandwidth_package_sites_struct) structure is documented below.

<a name="bandwidth_package_sites_struct"></a>
The `bandwidth_package_sites` block supports:

* `id` - Indicates the instance ID.

* `site_code` - Indicates the site code.

* `site_type` - Indicates the site type. The default type is **region**.

* `name_cn` - Indicates the instance Chinese name.

* `name_en` - Indicates the  instance English name.

* `description` - Indicates the description.

* `region_id` - Indicates the region ID.

* `created_at` - Indicates the creation time.
  The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.

* `updated_at` - Indicates the update time.
  The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.
