---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_logs"
description: |-
  Use this data source to query CFW logs within HuaweiCloud.
---

# huaweicloud_cfw_logs

Use this data source to query CFW logs within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "start_time" {}
variable "end_time" {}
variable "log_type" {}
variable "type" {}

data "huaweicloud_cfw_logs" "test" {
  fw_instance_id = var.fw_instance_id
  start_time     = var.start_time
  end_time       = var.end_time
  log_type       = var.log_type
  type           = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `start_time` - (Required, Int) Specifies the start time in millisecond timestamp.

* `end_time` - (Required, Int) Specifies the end time in millisecond timestamp.

* `log_type` - (Required, String) Specifies the log type.  
  The valid values are as follows:
  + **internet**: North-south traffic log.
  + **nat**: NAT scenario log.
  + **vpc**: East-west traffic log.
  + **vgw**: VGW scenario log.

* `type` - (Required, String) Specifies the log category.  
  The valid values are as follows:
  + **attack**: Attack log.
  + **acl**: Access control log.
  + **flow**: Flow log.
  + **url**: URL log.

* `filters` - (Optional, List) Specifies the filter conditions.

  The [filters](#filters_struct) structure is documented below.

* `log_id` - (Optional, String) Specifies the document ID for pagination.

  Please refer to the document link
  [reference](https://support.huaweicloud.com/intl/en-us/api-cfw/ListLogs.html#ListLogs__request_Filter) for values.

* `next_date` - (Optional, Int) Specifies the next query cursor.

  Please refer to the document link
  [reference](https://support.huaweicloud.com/intl/en-us/api-cfw/ListLogs.html#ListLogs__request_Filter) for values.

<a name="filters_struct"></a>
The `filters` block supports:

* `field` - (Required, String) Specifies the log field name, for example **src_ip**.

* `operator` - (Required, String) Specifies the operator.  `
  The valid values are as follows:
  + **equal**: Equals.
  + **not_equal**: Not equals.
  + **contain**: Contains.
  + **starts_with**: Starts with.

* `values` - (Optional, List) Specifies the field values.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The log records.

  The [records](#log_records_struct) structure is documented below.

<a name="log_records_struct"></a>
The `records` block supports:

* `app` - The application.

* `bytes` - The number of flow bytes. Present in flow logs.

* `direction` - The session direction.  
  The valid values are as follows:
  + **out2in**: From external to internal.
  + **in2out**: From internal to external.

* `dst_host` - The destination host.

* `dst_ip` - The destination IP address.

* `dst_port` - The destination port.

* `end_time` - The session end time in millisecond timestamp. Present in flow logs.

* `log_id` - The log ID used for pagination.

* `packets` - The number of flow packets. Present in flow logs.

* `protocol` - The protocol.

* `src_ip` - The source IP address.

* `src_port` - The source port.

* `start_time` - The session start time in millisecond timestamp. Present in flow logs.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `dst_province_id` - The destination province ID.

* `dst_province_name` - The destination province name.

* `dst_city_id` - The destination city ID.

* `dst_city_name` - The destination city name.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `src_province_id` - The source province ID.

* `src_province_name` - The source province name.

* `src_city_id` - The source city ID.

* `src_city_name` - The source city name.

* `vgw_id` - The virtual gateway ID.

* `sctp_verification_tag` - The SCTP verification tag. Present in flow logs.

* `sctp_is_handshake_flow` - The SCTP handshake flow flag. Present in flow logs.

* `qos_rule_id` - The QoS rule ID. Present in flow and access control logs.

* `qos_rule_name` - The QoS rule name. Present in flow and access control logs.

* `qos_channel_id` - The QoS channel ID. Present in flow logs.

* `qos_channel_name` - The QoS channel name. Present in flow logs.

* `qos_drop_packets` - The number of QoS-dropped packets. Present in flow logs.

* `qos_drop_bytes` - The number of QoS-dropped bytes. Present in flow logs.

* `qos_rule_type` - The QoS rule type. Present in flow and access control logs.

* `qos_channel_type` - The QoS channel type. Present in flow logs.

* `action` - The action. Present in attack, access control, and URL logs.

* `url` - The URL. Present in URL logs.

* `hit_time` - The hit time in millisecond timestamp. Present in access control and URL logs.

* `rule_id` - The rule ID. Present in access control and URL logs.

* `rule_name` - The rule name. Present in access control and URL logs.

* `rule_type` - The rule type. Present in access control and URL logs.

* `attack_rule` - The attack rule. Present in attack logs.

* `attack_rule_id` - The attack rule ID. Present in attack logs.

* `attack_type` - The attack type. Present in attack logs.

* `event_time` - The event time in millisecond timestamp. Present in attack logs.

* `level` - The attack level. Present in attack logs.

* `packet` - The rule payload. Present in attack logs.

* `source` - The attack source. Present in attack logs.

* `real_ip` - The real IP address. Present in attack logs.

* `tag` - The tag type. Present in attack logs. The value `1` indicates a WAF back-to-source IP.
