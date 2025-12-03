---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_operational_report_welfare"
description: |-
  Use this data source to query the information in the news and promotions area of a monthly operations report.
---

# huaweicloud_hss_operational_report_welfare

Use this data source to query the information in the news and promotions area of a monthly operations report.

## Example Usage

```hcl
data "huaweicloud_hss_operational_report_welfare" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hot_info` - The hot events information.

  The [hot_info](#hot_info_struct) structure is documented below.

* `version_update_info` - The version update information.

  The [version_update_info](#version_update_info_struct) structure is documented below.

* `activities_info` - The promotional activities information.

  The [activities_info](#activities_info_struct) structure is documented below.

<a name="hot_info_struct"></a>
The `hot_info` block supports:

* `title` - The hot event title.

* `url_json` - The hot event links.

<a name="version_update_info_struct"></a>
The `version_update_info` block supports:

* `title` - The version update titel.

* `url_json` - The version update links.

<a name="activities_info_struct"></a>
The `activities_info` block supports:

* `title` - The promotion title.

* `url_json` - The promotion links.
