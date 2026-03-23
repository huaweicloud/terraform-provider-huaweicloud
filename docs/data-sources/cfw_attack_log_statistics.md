---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_attack_log_statistics"
description: |-
  Use this data source to get the CFW attack log statistics.
---

# huaweicloud_cfw_attack_log_statistics

Use this data source to get the CFW attack log statistics.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_attack_log_statistics" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "internet"
  action         = 0
  item           = "src_region_id"
  value          = "target-object"
  range          = "2"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `log_type` - (Required, String) Specifies the log type. Valid values are:
  + **internet**: North-South oriented log
  + **nat**: NAT scenario log
  + **vpc**: East-West oriented log
  + **vgw**: VGW scenario log

* `action` - (Required, Int) Specifies the action. Valid values are:
  + `0`: All
  + `1`: Intercept

* `item` - (Required, String) Specifies the aggregation type. Valid values are:
  + **src_region_id**: TOP external attack source regions
  + **attack_type**: TOP attack types
  + **in_src_ip**: TOP internal attack source IPs
  + **out_src_ip**: TOP external attack source IPs
  + **dst_port**: TOP destination ports
  + **dst_ip**: TOP destination IPs
  + **attack_rule**: TOP attack rules
  + **src_ip**: TOP source IPs

* `value` - (Required, String) Specifies the statistical objects.

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

* `data` - The data.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `app_count` - The application count.

* `attack_rule_count` - The attack rule count.

* `attack_type_count` - The attack type count.

* `count` - The number of attacks.

* `dst_ip_count` - The destination IP count.

* `dst_port_count` - The destination port count.

* `end_time` - The end time.

* `records` - The attack details.

  The [records](#data_records_struct) structure is documented below.

* `src_ip_count` - The source IP count.

* `start_time` - The start time.

* `total` - The total number.

<a name="data_records_struct"></a>
The `records` block supports:

* `action` - The action.

* `app` - The application.

* `attack_rule` - The attack rule.

* `attack_rule_id` - The attack rule ID.

* `attack_type` - The attack type.

* `direction` - The attack direction.

* `dst_ip` - The destination IP.

* `dst_port` - The destination port.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `dst_province_id` - The destination province ID.

* `dst_province_name` - The destination province name.

* `dst_city_id` - The destination city ID.

* `dst_city_name` - The destination city name.

* `event_time` - The event time.

* `level` - The risk level.

* `protocol` - The protocol.

* `source` - The source.

* `src_ip` - The source IP.

* `real_ip` - The real IP.

* `tag` - The tag.

* `src_port` - The source port.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `src_province_id` - The source province ID.

* `src_province_name` - The source province name.

* `src_city_id` - The source city ID.

* `src_city_name` - The source city name.

* `vgw_id` - The VGW ID.
