---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_attack_log_statistic"
description: |-
  Use this data source to get the CFW attack log statistic within HuaweiCloud.
---

# huaweicloud_cfw_attack_log_statistic

Use this data source to get the CFW attack log statistic within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_attack_log_statistic" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = "internet"
  item           = "src"
  range          = "0"
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
  + **dst**: Top attacked destinations.
  + **src**: Top attack sources.

* `range` - (Optional, String) Specifies the time range.  
  The valid values are as follows:
  + **0**: Last one hour.
  + **1**: Last one day.
  + **2**: Last seven days.

* `direction` - (Optional, String) Specifies the session direction.  
  The valid values are as follows:
  + **in2out**: Outbound direction.
  + **out2in**: Inbound direction.

* `start_time` - (Optional, Int) Specifies the start time in millisecond timestamp.

* `end_time` - (Optional, Int) Specifies the end time in millisecond timestamp.

* `vgw_id` - (Optional, List) Specifies the list of VGW IDs.

* `size` - (Optional, String) Specifies the aggregation size. The value ranges from **0** to **100**.
  The default value is **50**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The attack log statistics.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `apps` - The application list.

  The [apps](#attack_statistic_top_info_struct) structure is documented below.

* `associated_name` - The name of the associated resource.

* `associated_type` - The type of the associated resource.  
  The valid values are as follows:
  + **PORT**: IPv4 ECS.
  + **IPV6_PORT**: IPv6 ECS.

* `attack_count` - The number of attacks.

* `attack_type` - The attack type.  
  The valid values are as follows:
  + **Access Control**: Access control.
  + **Vulnerability scanning**: Vulnerability scanning.
  + **Email attack**: Email attack.
  + **Vulnerability Attack**: Vulnerability attack.
  + **Web attack**: Web attack.
  + **password attack**: Password attack.
  + **Hijacking attack**: Hijacking attack.
  + **Protocol exception**: Protocol exception.
  + **Trojan horse**: Trojan horse.

* `deny_count` - The number of interceptions.

* `dst_ports` - The list of destination ports.

  The [dst_ports](#attack_statistic_top_info_struct) structure is documented below.

* `ip` - The IP address.

* `latest_time` - The latest attack time in millisecond timestamp.

* `region_id` - The region ID.

* `region_name` - The region name.

* `src_type` - The attack source type.

* `vgw_id` - The VGW ID.

<a name="attack_statistic_top_info_struct"></a>
The `apps` and `dst_ports` blocks support:

* `count` - The count.

* `item` - The item.

* `item_id` - The item ID.
