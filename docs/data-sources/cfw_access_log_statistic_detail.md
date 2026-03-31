---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_access_log_statistic_detail"
description: |-
  Use this data source to get the CFW access log statistic detail within HuaweiCloud.
---

# huaweicloud_cfw_access_log_statistic_detail

Use this data source to get the CFW access log statistic detail within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_access_log_statistic_detail" "test" {
  fw_instance_id = var.fw_instance_id
  item           = "dst_ip"
  item_id        = "192.0.2.1"
  range          = "0"
  direction      = "in2out"
  log_type       = "internet"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `item` - (Required, String) Specifies the aggregation type.  
  The valid values are as follows:
  + **top_deny_rule**: Blocking rules with the most hits.
  + **dst_ip**: Blocked destination IP addresses.
  + **src_ip**: Blocked source IP addresses.
  + **dst_port**: Blocked ports.
  + **dst_region**: Blocked destination regions.
  + **src_region**: Blocked source regions.

* `item_id` - (Required, String) Specifies the aggregation object.

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

* `log_type` - (Optional, String) Specifies the log type.  
  The valid values are as follows:
  + **internet**: North-south traffic log.
  + **nat**: NAT scenario log.
  + **vpc**: East-west traffic log.
  + **vgw**: VGW scenario log.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The access log statistic detail.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `dst_ip_count` - The number of destination IP addresses.

* `dst_port_count` - The number of destination ports.

* `hit_count` - The number of hits.

* `protocol_count` - The number of protocols.

* `recent_end_time` - The end time of the recent window.

* `recent_start_time` - The start time of the recent window.

* `record_total` - The number of records.

* `records` - The hit detail list.

  The [records](#access_log_info_struct) structure is documented below.

* `rule_hit_count` - The number of hit rules.

* `src_ip_count` - The number of source IP addresses.

<a name="access_log_info_struct"></a>
The `records` block supports:

* `action` - The action.

* `app` - The application.

* `url` - The URL.

* `dst_host` - The destination domain name.

* `dst_ip` - The destination IP address.

* `dst_port` - The destination port.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `dst_province_id` - The destination province ID.

* `dst_province_name` - The destination province name.

* `dst_city_id` - The destination city ID.

* `dst_city_name` - The destination city name.

* `hit_time` - The hit time.

* `log_id` - The log ID.

* `protocol` - The protocol.

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `rule_type` - The rule type.

* `src_ip` - The source IP address.

* `src_port` - The source port.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `src_province_id` - The source province ID.

* `src_province_name` - The source province name.

* `src_city_id` - The source city ID.

* `src_city_name` - The source city name.

* `vgw_id` - The VGW ID.

* `qos_rule_id` - The QoS rule ID.

* `qos_rule_name` - The QoS rule name.

* `qos_rule_type` - The QoS rule type.
