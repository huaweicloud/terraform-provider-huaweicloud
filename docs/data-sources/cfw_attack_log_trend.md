---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_attack_log_trend"
description: |-
  Use this data source to get the CFW attack log trend.
---

# huaweicloud_cfw_attack_log_trend

Use this data source to get the CFW attack log trend.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_attack_log_trend" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "internet"
  range          = "2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `log_type` - (Required, String) Specifies the log type. Valid values are:
  + **internet**: North-South oriented log
  + **nat**: NAT scenario log
  + **vpc**: East-West oriented log
  + **vgw**: VGW scenario log

* `range` - (Optional, String) Specifies the time range. Valid values are:
  + **0**: one hour
  + **1**: one day
  + **2**: seven days

* `start_time` - (Optional, Int) Specifies the start time, in millisecond timestamp.

* `end_time` - (Optional, Int) Specifies the end time, in millisecond timestamp.

* `vgw_id` - (Optional, List) Specifies the VGW ID list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The attack trends.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `deny_count` - The number of blocks.

* `permit_count` - The release times.

* `time` - The aggregation time.
