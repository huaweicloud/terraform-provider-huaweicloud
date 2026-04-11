---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_flow_log_trend"
description: |-
  Use this data source to get the CFW flow log trend within HuaweiCloud.
---

# huaweicloud_cfw_flow_log_trend

Use this data source to get the CFW flow log trend within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "log_type" {}

data "huaweicloud_cfw_flow_log_trend" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = var.log_type
  range          = "1"
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

* `range` - (Optional, String) Specifies the time range.  
  The valid values are as follows:
  + **0**: Last one hour.
  + **1**: Last one day.
  + **2**: Last seven days.

-> The parameters `range` and `start_time`/`end_time` cannot both be empty.

* `direction` - (Optional, String) Specifies the session direction.  
  The valid values are as follows:
  + **in2out**: Outbound direction.
  + **out2in**: Inbound direction.

* `start_time` - (Optional, Int) Specifies the start time in millisecond timestamp.

* `end_time` - (Optional, Int) Specifies the end time in millisecond timestamp.

* `vgw_id` - (Optional, List) Specifies the list of VGW IDs.

* `asset_type` - (Optional, String) Specifies the IP type.  
  The valid values are as follows:
  + **public**: Public IP address.
  + **private**: Private IP address.
  + **open_port**: Open port.

* `ip` - (Optional, List) Specifies the list of IP addresses.

* `vpc` - (Optional, List) Specifies the list of VPC IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The flow trend data points.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `agg_time` - The aggregation time.

* `in_bps` - The inbound bandwidth in bps.

* `out_bps` - The outbound bandwidth in bps.
