---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddosv2_eip_defense_statuses"
description: |-
  Use this data source to query the defense statuses of EIPs protected by Anti-DDoS V2.
---

# huaweicloud_antiddosv2_eip_defense_statuses

Use this data source to query the defense statuses of EIPs protected by Anti-DDoS V2.

## Example Usage

```hcl
variable "eip_status" {}
variable "eip_address" {}

data "huaweicloud_antiddosv2_eip_defense_statuses" "test" {
  status = var.eip_status
  ips    = var.eip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the protection status of the EIP. Options are:
  + **normal**: Normal.
  + **configging**: Configuring.
  + **notConfig**: Not configured.
  + **packetcleaning**: Packet cleaning.
  + **packetdropping**: Packet dropping.

* `ips` - (Optional, String) Specifies the IP address for filtering, supports partial matching.
  For example, if you enter **192.168**, it will match IPs like **192.168.111.1** and **10.192.168.8**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `ddos_status` - The list of EIP defense statuses.
  The [ddos_status](#ddos_status_struct) structure is documented below.

<a name="ddos_status_struct"></a>
The `ddos_status` block supports:

* `floating_ip_id` - The ID of the EIP.

* `floating_ip_address` - The IP address of the EIP.

* `product_type` - The type of the EIP protection service. Options are:
  + **Anti-DDoS**: EIP belongs to Anti-DDoS traffic cleaning.
  + **CNAD**: EIP belongs to DDoS native advanced protection.

* `status` - The protection status of the EIP. Valid values are:
  + **normal**: Normal.
  + **configging**: Configuring.
  + **notConfig**: Not configured.
  + **packetcleaning**: Packet cleaning.
  + **packetdropping**: Packet dropping.

* `clean_threshold` - The cleaning threshold.

* `block_threshold` - The blackhole threshold.
