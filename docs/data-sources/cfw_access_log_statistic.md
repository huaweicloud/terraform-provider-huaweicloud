---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_access_log_statistic"
description: |-
  Use this data source to get the CFW access log statistic within HuaweiCloud.
---

# huaweicloud_cfw_access_log_statistic

Use this data source to get the CFW access log statistic within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_access_log_statistic" "test" {
  fw_instance_id = var.fw_instance_id
  item           = "strategy_dashboard"
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
  + **strategy_hit_info**: Policy hit trend.
  + **strategy_dashboard**: Policy hit overview.
  + **top_deny_rule**: Blocking rules with the most hits.
  + **dst_ip**: Top blocked destination IP addresses.
  + **src_ip**: Top blocked source IP addresses.
  + **dst_port**: Top blocked ports.
  + **dst_region**: Top blocked destination regions.
  + **src_region**: Top blocked source regions.

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

* `rule_id` - (Optional, List) Specifies the list of rule IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The access log statistic.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `deny_count` - The block count.

* `deny_top_one_acl_id` - The ID of a blocking policy that is frequently hit.

* `deny_top_one_acl_name` - The name of a blocking policy that is frequently hit.

* `hit_count` - The number of hits.

* `in2out_deny_dst_ip_list` - The list of top blocked destination IP addresses in the outbound direction.

  The [in2out_deny_dst_ip_list](#access_top_member_vo_struct) structure is documented below.

* `in2out_deny_dst_port_list` - The list of top blocked ports in the outbound direction.

  The [in2out_deny_dst_port_list](#access_top_member_vo_struct) structure is documented below.

* `in2out_deny_dst_region_list` - The list of top blocked destination regions in the outbound direction.

  The [in2out_deny_dst_region_list](#access_top_member_vo_struct) structure is documented below.

* `in2out_deny_src_ip_list` - The list of top blocked source IP addresses in the outbound direction.

  The [in2out_deny_src_ip_list](#access_top_member_vo_struct) structure is documented below.

* `out2in_deny_dst_ip_list` - The list of top blocked destination IP addresses in the inbound direction.

  The [out2in_deny_dst_ip_list](#access_top_member_vo_struct) structure is documented below.

* `out2in_deny_dst_port_list` - The list of top blocked destination ports in the inbound direction.

  The [out2in_deny_dst_port_list](#access_top_member_vo_struct) structure is documented below.

* `out2in_deny_src_ip_list` - The list of top blocked source IP addresses in the inbound direction.

  The [out2in_deny_src_ip_list](#access_top_member_vo_struct) structure is documented below.

* `out2in_deny_src_port_list` - The list of top blocked source ports in the inbound direction.

  The [out2in_deny_src_port_list](#access_top_member_vo_struct) structure is documented below.

* `out2in_deny_src_region_list` - The list of top blocked source regions in the inbound direction.

  The [out2in_deny_src_region_list](#access_top_member_vo_struct) structure is documented below.

* `permit_count` - The allow count.

* `permit_top_one_acl_id` - The ID of an allow policy that is frequently hit.

* `permit_top_one_acl_name` - The name of an allow policy that is frequently hit.

* `records` - The list of hit trend records.

  The [records](#records_struct) structure is documented below.

* `top_deny_rule_list` - The top blocking rule list.

  The [top_deny_rule_list](#access_top_member_vo_struct) structure is documented below.

<a name="access_top_member_vo_struct"></a>
The `in2out_deny_dst_ip_list`, `in2out_deny_dst_port_list`, `in2out_deny_dst_region_list`, `in2out_deny_src_ip_list`,
`out2in_deny_dst_ip_list`, `out2in_deny_dst_port_list`, `out2in_deny_src_ip_list`, `out2in_deny_src_port_list`,
`out2in_deny_src_region_list` and `top_deny_rule_list` blocks support:

* `count` - The count.

* `item` - The item.

* `name` - The item name.

<a name="records_struct"></a>
The `records` block supports:

* `agg_time` - The aggregation time.

* `deny_access_top_counts` - The number of blocked objects.

* `permit_access_top_counts` - The number of allowed objects.

* `total_access_top_counts` - The number of hits.
