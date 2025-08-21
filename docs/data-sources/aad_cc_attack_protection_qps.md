---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_cc_attack_protection_qps"
description: |-
  Use this data source to get the QPS information of CC attack protection within HuaweiCloud.
---

# huaweicloud_aad_cc_attack_protection_qps

Use this data source to get the QPS information of CC attack protection within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_cc_attack_protection_qps" "test" {
  recent = "1month"
}
```

## Argument Reference

The following arguments are supported:

* `recent` - (Required, String) Specifies the time range for querying CC attack protection QPS data.
  Valid values are **yesterday**, **today**, **3days**, **1week**, **1month**.

* `domains` - (Optional, String) Specifies the domain name to query.

* `start_time` - (Optional, String) Specifies the start time for querying CC attack protection QPS data.

* `end_time` - (Optional, String) Specifies the end time for querying CC attack protection QPS data.

* `overseas_type` - (Optional, String) Specifies the protection region.
  The options are as follows:
  + `0`: Mainland China.
  + `1`: Outside Mainland China.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `qps` - The QPS value of CC attack protection.
