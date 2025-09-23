---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_ddos_attack_protection_info"
description: |-
  Use this data source to get the list of Advanced Anti-DDos DDoS attack protection info within HuaweiCloud.
---

# huaweicloud_aad_ddos_attack_protection_info

Use this data source to get the list of Advanced Anti-DDos DDoS attack protection info within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "ip" {}
variable "type" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_aad_ddos_attack_protection_info" "test" {
  instance_id = var.instance_id
  ip          = var.ip
  type        = var.type
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

* `ip` - (Required, String) Specifies the high defense IP address.

* `type` - (Required, String) Specifies the request type. The valid values are **pps** and **bps**.

* `start_time` - (Required, String) Specifies the start time.

* `end_time` - (Required, String) Specifies the end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `flow_bps` - The BPS flow information list. This field is returned when type is **bps**.  
  The [flow_bps](#flow_bps_struct) structure is documented below.

* `flow_pps` - The PPS flow information list. This field is returned when type is **pps**.  
  The [flow_pps](#flow_pps_struct) structure is documented below.

<a name="flow_bps_struct"></a>
The `flow_bps` block supports:

* `utime` - The data time.

* `attack_bps` - The attack traffic.

* `normal_bps` - The normal traffic.

<a name="flow_pps_struct"></a>
The `flow_pps` block supports:

* `utime` - The data time.

* `attack_pps` - The attack packet rate.

* `normal_pps` - The normal packet rate.
