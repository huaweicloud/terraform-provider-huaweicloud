---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_sites"
description: |-
  Use this data source to get the list of CC global connection bandwidth sites.
---

# huaweicloud_cc_global_connection_bandwidth_sites

Use this data source to get the list of CC global connection bandwidth sites.

## Example Usage

```hcl
variable "site_id" {}

data "huaweicloud_cc_global_connection_bandwidth_sites" "test" {
  site_id = var.site_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `site_id` - (Optional, String) Specifies site ID.

* `site_code` - (Optional, String) Specifies site code.

* `site_type` - (Optional, String) Specifies site type.
  + **Area**: a geographic region site.
  + **SubArea**: a region site.
  + **Region**: a multi-city site.

* `name_en` - (Optional, String) Specifies the site name in English.

* `name_cn` - (Optional, String) Specifies the site name in Chinese.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `site_infos` - The site list.

  The [site_infos](#site_infos_struct) structure is documented below.

<a name="site_infos_struct"></a>
The `site_infos` block supports:

* `id` - The site ID.

* `description` - The site description.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `name_en` - User-defined site name in English.

* `name_cn` - User-defined site name in Chinese.

* `site_code` - The site code.

* `site_type` - The site type.

* `service_list` - The list of services supported at the site. Multiple services are separated by commas (,).

* `group_list` - The site group list.

  The [group_list](#site_infos_group_list_struct) structure is documented below.

* `region_id` - The site region ID.

* `public_border_group` - Whether the site is a central site or an edge site.

<a name="site_infos_group_list_struct"></a>
The `group_list` block supports:

* `id` - The site group list ID.

* `description` - The site group list description.

* `name_en` - User-defined site group name in English.

* `name_cn` - User-defined site group name in Chinese.
