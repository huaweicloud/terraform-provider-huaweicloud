---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_eip_protection_traffic"
description: |-
  Use this data source to query the EIP protection traffic of a specified EIP within HuaweiCloud.
---

# huaweicloud_antiddos_eip_protection_traffic

Use this data source to query the EIP protection traffic of a specified EIP within HuaweiCloud.

## Example Usage

```hcl
variable "floating_ip_id" {}
variable "eip_address" {}

data "huaweicloud_antiddos_eip_protection_traffic" "test" {
  floating_ip_id = var.floating_ip_id
  ip             = var.eip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `floating_ip_id` - (Required, String) Specifies the ID of the EIP.

* `ip` - (Optional, String) Specifies the EIP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The list of EIP protection traffic. The [data](#data_block) structure is documented below.

<a name="data_block"></a>
The `data` block supports:

* `period_start` - The start time of the statistics period.

* `bps_in` - The inbound traffic rate, in bit/s.

* `bps_attack` - The attack traffic rate, in bit/s.

* `total_bps` - The total traffic rate, in bit/s.

* `pps_in` - The inbound packet rate, in packets per second (pps).

* `pps_attack` - The attack packet rate, in packets per second (pps).

* `total_pps` - The total packet rate, in packets per second (pps).
