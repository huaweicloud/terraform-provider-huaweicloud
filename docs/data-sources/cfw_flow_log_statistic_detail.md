---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_flow_log_statistic_detail"
description: |-
  Use this data source to get the CFW flow log statistic detail within HuaweiCloud.
---

# huaweicloud_cfw_flow_log_statistic_detail

Use this data source to get the CFW flow log statistic detail within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_flow_log_statistic_detail" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "internet"
  item           = "dst_ip"
  value          = "100.93.4.158"
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

* `value` - (Required, String) Specifies the statistic object.

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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The flow log statistic detail.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `app_count` - The application statistics count.

* `bytes` - The total bytes.

* `dst_ip_count` - The destination IP count.

* `dst_port_count` - The destination port count.

* `end_time` - The end time in the summarized result.

* `records` - The TOP session detail list.

  The [records](#records_struct) structure is documented below.

* `request_byte` - The request bytes.

* `response_byte` - The response bytes.

* `sessions` - The session count.

* `src_ip_count` - The source IP count.

* `start_time` - The start time in the summarized result.

<a name="records_struct"></a>
The `records` block supports:

* `app` - The application.

* `bytes` - The bytes.

* `dst_associate_instance_type` - The destination IP associated asset type.

* `dst_device_name` - The destination IP associated asset name.

* `dst_ip` - The destination IP.

* `dst_port` - The destination port.

* `dst_host` - The destination domain name.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `end_time` - The session end time.

* `protocol` - The protocol.

* `request_byte` - The request bytes.

* `response_byte` - The response bytes.

* `sessions` - The session count.

* `src_associate_instance_type` - The source IP associated asset type.

* `src_device_name` - The source IP associated asset name.

* `src_ip` - The source IP.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `start_time` - The session start time.
