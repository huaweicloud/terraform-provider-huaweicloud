---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_report_history"
description: |-
  Use this data source to get the CFW security report send history within HuaweiCloud.
---

# huaweicloud_cfw_report_history

Use this data source to get the CFW security report send history within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "report_profile_id" {}

data "huaweicloud_cfw_report_history" "test" {
  fw_instance_id    = var.fw_instance_id
  report_profile_id = var.report_profile_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `report_profile_id` - (Required, String) Specifies the security report profile ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The report send history list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `report_id` - The report ID.

* `send_time` - The send time in millisecond timestamp.

* `statistic_period` - The statistics period for custom reports.

  The [statistic_period](#statistic_period_struct) structure is documented below.

<a name="statistic_period_struct"></a>
The `statistic_period` block supports:

* `end_time` - The end time in millisecond timestamp.

* `start_time` - The start time in millisecond timestamp.
