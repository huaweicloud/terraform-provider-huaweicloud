---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_charts"
description: |-
  Use this data source to get the list of CCE charts within HuaweiCloud.
---

# huaweicloud_cce_charts

Use this data source to get the list of CCE charts within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cce_charts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE charts. If omitted, the
  provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `charts` - The list of charts.
  The [charts](#CCE_charts) structure is documented below.

<a name="CCE_charts"></a>
The `charts` block supports:

* `id` - The chart ID.

* `name` - The chart name.

* `values` - The values of the chart.

* `translate` - The traslate source of the chart.

* `instruction` - The instruction of the chart.

* `version` - The chart version.

* `description` - The description of the chart.

* `source` - The source of the chart.

* `icon_url` - The icon URL.

* `public` - Whether the chart is public.

* `chart_url` - The chart URL.

* `created_at` - The create time.

* `updated_at` - The update time.
