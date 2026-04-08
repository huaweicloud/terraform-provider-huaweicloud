---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_chart_values"
description: |-
  Use this data source to get CCE Autopilot chart values.
---

# huaweicloud_cce_autopilot_chart_values

Use this data source to get CCE Autopilot chart values.

## Example Usage

```hcl
variable "chart_id" {}

data "huaweicloud_cce_autopilot_chart_values" "test" {
  chart_id = var.chart_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `chart_id` - (Required, String) Specifies the ID of CCE autopilot chart.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `values` - The values of CCE autopilot chart template.
