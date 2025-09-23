---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_attack_logs"
description: |-
  Use this data source to get the list of CFW attack logs.
---

# huaweicloud_cfw_attack_logs

Use this data source to get the list of CFW attack logs.

-> **NOTE:** Up to 1000 logs can be retrieved. Set filter criteria to narrow down the search scope.

## Example Usage

```hcl
variable start_time {}
variable end_time {}
variable fw_instance_id {}

data "huaweicloud_cfw_attack_logs" "test" {
  fw_instance_id = var.fw_instance_id
  start_time     = var.start_time
  end_time       = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `start_time` - (Required, String) Specifies the start time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `end_time` - (Required, String) Specifies the end time. The time is in UTC.
  The format is **yyyy-MM-dd HH:mm:ss**.

* `src_ip` - (Optional, String) Specifies the source IP address.

* `src_port` - (Optional, Int) Specifies the source port.

* `dst_ip` - (Optional, String) Specifies the destination IP address.

* `dst_port` - (Optional, Int) Specifies the destination port.

* `app` - (Optional, String) Specifies the application protocol.

* `attack_type` - (Optional, String) Specifies the intrusion event type.

* `attack_rule` - (Optional, String) Specifies the intrusion event rule.

* `level` - (Optional, String) Specifies the threat level.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `log_type` - (Optional, String) Specifies the log type.
  The valid values are **internet**, **nat** and **vpc**.

* `attack_rule_id` - (Optional, String) Specifies the attack rule ID.

* `src_region_name` - (Optional, String) Specifies the source region name.

* `dst_region_name` - (Optional, String) Specifies the destination region name.

* `src_province_name` - (Optional, String) Specifies the source province name.

* `dst_province_name` - (Optional, String) Specifies the destination province name.

* `src_city_name` - (Optional, String) Specifies the source city name.

* `dst_city_name` - (Optional, String) Specifies the destination city name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The attack log records.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `packet` - The attack log packet.

* `source` - The source.

* `src_ip` - The source IP address.

* `src_port` - The source port.

* `direction` - The direction.

* `dst_port` - The destination port.

* `app` - The application protocol.

* `attack_rule_id` - The attack rule ID.

* `protocol` - The protocol.

* `action` - The action.

* `event_time` - The event time.

* `attack_rule` - The attack rule.

* `log_id` - The log ID.

* `dst_ip` - The destination IP address.

* `packet_messages` - The packet messages.

  The [packet_messages](#records_packet_messages_struct) structure is documented below.

* `attack_type` - The attack type.

* `level` - The threat level.

* `packet_length` - The packet length.

* `src_region_id` - The source region ID.

* `src_region_name` - The source region name.

* `dst_region_id` - The destination region ID.

* `dst_region_name` - The destination region name.

* `src_province_id` - The source province ID.

* `src_province_name` - The source province name.

* `src_city_id` - The source city ID.

* `src_city_name` - The source city name.

* `dst_province_id` - The destination province ID.

* `dst_province_name` - The destination province name.

* `dst_city_id` - The destination city ID.

* `dst_city_name` - The destination city name.

<a name="records_packet_messages_struct"></a>
The `packet_messages` block supports:

* `utf8_string` - The utf-8 string.

* `hex_index` - The hexadecimal index.

* `hexs` - The hexadecimal series.
