---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_attack_log_total"
description: |-
  Use this data source to get the CFW attack log total within HuaweiCloud.
---

# huaweicloud_cfw_attack_log_total

Use this data source to get the CFW attack log total within HuaweiCloud.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "log_type" {}
variable "range" {}

data "huaweicloud_cfw_attack_log_total" "test" {
  fw_instance_id = var.fw_instance_id
  log_type       = var.log_type
  range          = var.range
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

* `start_time` - (Optional, Int) Specifies the start time in millisecond timestamp.

* `end_time` - (Optional, Int) Specifies the end time in millisecond timestamp.

* `vgw_id` - (Optional, List) Specifies the list of VGW IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The attack overview statistics.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `attack_count` - The number of attacks.

* `deny_count` - The number of blocks.

* `permit_count` - The number of permits.

* `risk_ports` - The number of risk ports.
