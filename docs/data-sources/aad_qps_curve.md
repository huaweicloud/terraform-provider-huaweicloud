---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_qps_curve"
description: |-
  Use this data source to get the list of Advanced Anti-DDoS QPS curve within HuaweiCloud.
---

# huaweicloud_aad_qps_curve

Use this data source to get the list of Advanced Anti-DDoS QPS curve within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_qps_curve" "test" {
  value_type = "mean"
  recent     = "today"
}
```

## Argument Reference

The following arguments are supported:

* `value_type` - (Required, String) Specifies the value type.
  The valid values are as follows:
  + **mean**: Average QPS value.
  + **peak**: Peak QPS value.
  + **source**: Response status code from origin server.
  + **proxy**: Response status code from AAD protection.

* `domains` - (Optional, String) Specifies the domain name to query. If not specified, data for all domains will be returned.

* `start_time` - (Optional, String) Specifies the start time of the query.

* `end_time` - (Optional, String) Specifies the end time of the query.

* `recent` - (Optional, String) Specifies the recent time range.
  The valid values are:
  + **yesterday**
  + **today**
  + **3days**
  + **1week**
  + **1month**

* `overseas_type` - (Optional, String) Specifies the instance type.
  The valid values are:
  + **0**: Mainland China.
  + **1**: Outside Mainland China.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `curve` - The list of QPS curve detail.
  The [curve](#curve_struct) structure is documented below.

<a name="curve_struct"></a>
The `curve` block supports:

* `time` - The timestamp of the QPS curve.

* `total` - The total number of requests.

* `attack` - The number of attack requests.

* `basic` - The number of requests processed by web basic protection.

* `cc` - The number of requests processed by CC attack protection.

* `custom_custom` - The number of requests processed by precise protection.
