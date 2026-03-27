---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_traffic_log_trend"
description: |-
  Use this data source to get the CFW traffic log trend within HuaweiCloud.
---

# huaweicloud_cfw_traffic_log_trend

Use this data source to get the CFW traffic log trend within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_traffic_log_trend" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "vpc"
  agg_type       = "max"
  range          = "2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `log_type` - (Required, String) Specifies the log type.  
  The valid values are as follows:
  + **internet**: North-south traffic log.
  + **nat**: NAT scenario log.
  + **vpc**: East-west traffic log.
  + **vgw**: VGW scenario log.

* `agg_type` - (Required, String) Specifies the aggregation type.  
  The valid values are as follows:
  + **avg**: Average.
  + **max**: Maximum.

* `range` - (Optional, String) Specifies the time range.  
  The valid values are as follows:
  + **0**: Last one hour.
  + **1**: Last one day.
  + **2**: Last seven days.

-> The parameters `range` and `start_time`/`end_time` cannot both be empty.

* `start_time` - (Optional, Int) Specifies the start time in millisecond timestamp.

* `end_time` - (Optional, Int) Specifies the end time in millisecond timestamp.

* `vgw_id` - (Optional, List) Specifies the list of VGW IDs.

* `ip` - (Optional, List) Specifies the list of IP addresses to query traffic trend for.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The list of traffic trend data.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `agg_time` - The aggregation time point in millisecond timestamp.

* `bps` - The bits per second.

* `in_bps` - The inbound bits per second.

* `out_bps` - The outbound bits per second.
