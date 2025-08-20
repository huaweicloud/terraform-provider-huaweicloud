---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_bandwidth_curve"
description: |-
  Use this data source to get the list of Advanced Anti-DDos bandwidth curve within HuaweiCloud.
---

# huaweicloud_aad_bandwidth_curve

Use this data source to get the list of Advanced Anti-DDos bandwidth curve within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_bandwidth_curve" "test" {}
```

## Argument Reference

The following arguments are supported:

* `value_type` - (Required, String) Specifies the value type.  
  The valid values are as follows:
  + **mean**: Average value.
  + **peak**: Peak value.

* `domains` - (Optional, String) Specifies the domains. If not specified, all domains are queried.

* `start_time` - (Optional, String) Specifies the start time.

* `end_time` - (Optional, String) Specifies the end time.

* `recent` - (Optional, String) Specifies recent.  
  The valid values are as follows:
  + **yesterday**
  + **today**
  + **3days**
  + **1week**
  + **1month**

  `recent` cannot be empty when both `start_time` and `end_time` are empty.

* `overseas_type` - (Optional, String) Specifies instance type.
  The valid values are as follows:
  + **0**: Mainland.
  + **1**: Overseas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `curve` - The list of bandwidth curve detail.  
  The [curve](#curve_struct) structure is documented below.

<a name="curve_struct"></a>
The `curve` block supports:

* `in` - The ingress bandwidth.

* `out` - The egress bandwidth.

* `time` - The timestamp.
