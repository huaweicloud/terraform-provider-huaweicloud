---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_charts"
description: |- 
  Use this data source to get  the list of CCE Autopilot charts within huaweicloud.
---

# huaweicloud_cce_autopilot_charts

Use this data source to get  the list of CCE Autopilot charts within huaweicloud.

## Example Usage

```hcl
data "huaweicloud_cce_autopilot_charts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `charts` - The charts data in the cce cluster.

  The [object](#charts) structure is documented below.

<a name="charts"></a>
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

* `create_at` - The create time.

* `update_at` - The update time.
