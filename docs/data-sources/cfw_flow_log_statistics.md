---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_flow_log_statistics"
description: |-
  Use this data source to get the CFW flow log statistics within HuaweiCloud.
---

# huaweicloud_cfw_flow_log_statistics

Use this data source to get the CFW flow log statistics within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_flow_log_statistics" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "internet"
  item           = "dst_host"
  range          = "2"
  direction      = "in2out"
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

* `item` - (Required, String) Specifies the aggregation type.  
  The valid values are as follows:
  + **src_ip**: Source IP.
  + **dst_ip**: Destination IP.
  + **dst_port**: Destination port.
  + **protocol**: Protocol.
  + **dst_host**: Destination domain name.
  + **app**: Application.
  + **dst_region_name**: Destination region.
  + **src_region_name**: Source region.
  + **risk_ip**: Risk IP.
  + **risk_host**: Risk domain name.
  + **open_port**: Open port.

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
  + **public**: Public IP.
  + **private**: Private IP.
  + **open_port**: Open port.

* `size` - (Optional, String) Specifies the number of aggregated records. The value ranges from **0** to **10**.
  Defaults to **5**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The list of flow log statistics records.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `apps` - The application statistics.

  The [apps](#item_vo_struct) structure is documented below.

* `associate_instance_type` - The associated asset type.

* `device_name` - The associated asset name.

* `item` - The aggregation item value.

* `last_time` - The last access time in millisecond timestamp.

* `agg_start_time` - The aggregation start time in millisecond timestamp.

* `agg_end_time` - The aggregation end time in millisecond timestamp.

* `ports` - The port statistics.

  The [ports](#item_vo_struct) structure is documented below.

* `region` - The region.

* `request_byte` - The number of request bytes.

* `response_byte` - The number of response bytes.

* `sessions` - The number of sessions.

* `tags` - The tags.

* `src_ip` - The source IP statistics.

  The [src_ip](#item_vo_struct) structure is documented below.

* `dst_ip` - The destination IP statistics.

  The [dst_ip](#item_vo_struct) structure is documented below.

* `protocol` - The protocol.

<a name="item_vo_struct"></a>
The `apps`, `ports`, `src_ip` and `dst_ip` blocks support:

* `key` - The aggregation item key.

* `name` - The aggregation item name.

* `value` - The statistics value.
