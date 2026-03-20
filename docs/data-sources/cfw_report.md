---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_report"
description: |-
  Use this data source to get the CFW report.
---

# huaweicloud_cfw_report

Use this data source to get the CFW report.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "report_profile_id" {}
variable "report_id" {}

data "huaweicloud_cfw_report" "test" {
  fw_instance_id    = var.fw_instance_id
  report_profile_id = var.report_profile_id
  report_id         = var.report_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `report_profile_id` - (Required, String) Specifies the firewall report profile ID.

* `report_id` - (Required, String) Specifies the firewall report ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attack_info` - The security incident details.

  The [attack_info](#data_attack_info_struct) structure is documented below.

* `category` - The report type.

* `internet_firewall` - The internet boundary firewall.

  The [internet_firewall](#data_internet_firewall_struct) structure is documented below.

* `send_time` - The send time.

* `statistic_period` - The statistical period.

  The [statistic_period](#data_statistic_period_struct) structure is documented below.

* `vpc_firewall` - The VPC boundary firewall.

  The [vpc_firewall](#data_vpc_firewall_struct) structure is documented below.

<a name="data_statistic_period_struct"></a>
The `statistic_period` block supports:

* `end_time` - The end time in milliseconds.

* `start_time` - The start time in milliseconds.

<a name="data_vpc_firewall_struct"></a>
The `vpc_firewall` block supports:

* `app` - The application information.

  The [app](#vpc_firewall_app_struct) structure is documented below.

* `dst_ip` - The destination IP information.

  The [dst_ip](#vpc_firewall_dst_ip_struct) structure is documented below.

* `overview` - The overview information.

  The [overview](#vpc_firewall_overview_struct) structure is documented below.

* `src_ip` - The source IP information.

  The [src_ip](#vpc_firewall_src_ip_struct) structure is documented below.

* `traffic_trend` - The traffic trend information.

  The [traffic_trend](#vpc_firewall_traffic_trend_struct) structure is documented below.

* `vpc` - The VPC protection statistics.

  The [vpc](#vpc_firewall_vpc_struct) structure is documented below.

<a name="vpc_firewall_app_struct"></a>
The `app` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="vpc_firewall_dst_ip_struct"></a>
The `dst_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="vpc_firewall_overview_struct"></a>
The `overview` block supports:

* `access_policies` - The access control policies.

  The [access_policies](#overview_access_policies_struct) structure is documented below.

* `assets` - The asset count.

  The [assets](#overview_assets_struct) structure is documented below.

* `attack_event` - The threat event.

  The [attack_event](#overview_attack_event_struct) structure is documented below.

* `traffic_peak` - The traffic peak.

  The [traffic_peak](#overview_traffic_peak_struct) structure is documented below.

<a name="overview_access_policies_struct"></a>
The `access_policies` block supports:

* `eip` - The EIP access control policies.

* `nat` - The NAT access control policies.

* `total` - The total count.

* `changed` - The changed count.

<a name="overview_assets_struct"></a>
The `assets` block supports:

* `changed` - The changed count.

* `total` - The total count.

<a name="overview_attack_event_struct"></a>
The `attack_event` block supports:

* `changed` - The changed count.

* `deny` - The blocked count.

* `total` - The total count.

<a name="overview_traffic_peak_struct"></a>
The `traffic_peak` block supports:

* `in_bps` - The inbound bps.

* `out_bps` - The outbound bps.

* `permit` - The allowed count.

* `agg_time` - The aggregation time.

* `bps` - The bandwidth.

* `deny` - The blocked count.

<a name="vpc_firewall_src_ip_struct"></a>
The `src_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="vpc_firewall_traffic_trend_struct"></a>
The `traffic_trend` block supports:

* `in_bps` - The inbound bps.

* `out_bps` - The outbound bps.

* `permit` - The allowed count.

* `agg_time` - The aggregation time.

* `bps` - The bandwidth.

* `deny` - The blocked count.

<a name="vpc_firewall_vpc_struct"></a>
The `vpc` block supports:

* `protected` - The protected count.

  The [protected](#vpc_protected_struct) structure is documented below.

* `total` - The total count.

<a name="vpc_protected_struct"></a>
The `protected` block supports:

* `changed` - The changed count.

* `total` - The total count.

* `value` - The value.

<a name="data_attack_info_struct"></a>
The `attack_info` block supports:

* `src_ip` - The TOP source IP.

  The [src_ip](#attack_info_src_ip_struct) structure is documented below.

* `trend` - The attack trend.

  The [trend](#attack_info_trend_struct) structure is documented below.

* `type` - The TOP attack distribution.

  The [type](#attack_info_type_struct) structure is documented below.

* `dst_ip` - The TOP attack destination IP.

  The [dst_ip](#attack_info_dst_ip_struct) structure is documented below.

* `ips_mode` - The intrusion prevention status.

* `level` - The attack level distribution.

  The [level](#attack_info_level_struct) structure is documented below.

* `rule` - The TOP attack rule.

  The [rule](#attack_info_rule_struct) structure is documented below.

<a name="attack_info_src_ip_struct"></a>
The `src_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="attack_info_trend_struct"></a>
The `trend` block supports:

* `in_bps` - The inbound bps.

* `out_bps` - The outbound bps.

* `permit` - The allowed count.

* `agg_time` - The aggregation time.

* `bps` - The bandwidth.

* `deny` - The blocked count.

<a name="attack_info_type_struct"></a>
The `type` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="attack_info_dst_ip_struct"></a>
The `dst_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="attack_info_level_struct"></a>
The `level` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="attack_info_rule_struct"></a>
The `rule` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="data_internet_firewall_struct"></a>
The `internet_firewall` block supports:

* `in2out` - The outbound traffic.

  The [in2out](#internet_firewall_in2out_struct) structure is documented below.

* `out2in` - The inbound traffic.

  The [out2in](#internet_firewall_out2in_struct) structure is documented below.

* `overview` - The overview.

  The [overview](#internet_firewall_overview_struct) structure is documented below.

* `traffic_trend` - The traffic trend.

  The [traffic_trend](#internet_firewall_traffic_trend_struct) structure is documented below.

* `eip` - The EIP protection status.

  The [eip](#internet_firewall_eip_struct) structure is documented below.

<a name="internet_firewall_in2out_struct"></a>
The `in2out` block supports:

* `dst_port` - The TOP access port.

  The [dst_port](#in2out_dst_port_struct) structure is documented below.

* `src_ip` - The TOP access source IP.

  The [src_ip](#in2out_src_ip_struct) structure is documented below.

* `dst_host` - The TOP access destination host.

  The [dst_host](#in2out_dst_host_struct) structure is documented below.

* `dst_ip` - The TOP access destination IP.

  The [dst_ip](#in2out_dst_ip_struct) structure is documented below.

<a name="in2out_dst_port_struct"></a>
The `dst_port` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="in2out_src_ip_struct"></a>
The `src_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="in2out_dst_host_struct"></a>
The `dst_host` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="in2out_dst_ip_struct"></a>
The `dst_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="internet_firewall_out2in_struct"></a>
The `out2in` block supports:

* `dst_ip` - The TOP access destination IP.

  The [dst_ip](#out2in_dst_ip_struct) structure is documented below.

* `dst_port` - The TOP open port.

  The [dst_port](#out2in_dst_port_struct) structure is documented below.

* `src_ip` - The TOP access source IP.

  The [src_ip](#out2in_src_ip_struct) structure is documented below.

<a name="out2in_dst_ip_struct"></a>
The `dst_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="out2in_dst_port_struct"></a>
The `dst_port` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="out2in_src_ip_struct"></a>
The `src_ip` block supports:

* `key` - The aggregation item.

* `name` - The aggregation item name.

* `value` - The statistical value.

<a name="internet_firewall_overview_struct"></a>
The `overview` block supports:

* `access_policies` - The number of access control policies.

  The [access_policies](#overview_access_policies_struct) structure is documented below.

* `assets` - The number of assets.

  The [assets](#overview_assets_struct) structure is documented below.

* `attack_event` - The number of threat events.

  The [attack_event](#overview_attack_event_struct) structure is documented below.

* `traffic_peak` - The traffic peak.

  The [traffic_peak](#overview_traffic_peak_struct) structure is documented below.

<a name="internet_firewall_traffic_trend_struct"></a>
The `traffic_trend` block supports:

* `in_bps` - The inbound bps.

* `out_bps` - The outbound bps.

* `permit` - The number of allowed requests.

* `agg_time` - The aggregation time.

* `bps` - The bandwidth.

* `deny` - The number of blocked requests.

<a name="internet_firewall_eip_struct"></a>
The `eip` block supports:

* `protected` - The protection status.

  The [protected](#eip_protected_struct) structure is documented below.

* `total` - The total number of EIPs.

<a name="eip_protected_struct"></a>
The `protected` block supports:

* `changed` - The number of changes.

* `total` - The total number of changes.

* `value` - The value.
