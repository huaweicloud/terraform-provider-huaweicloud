---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_dnat_rules"
description: |-
  Use this data source to get the list of DNAT rules.
---

# huaweicloud_nat_dnat_rules

Use this data source to get the list of DNAT rules.

## Example Usage

```hcl
variable "protocol" {}

data "huaweicloud_nat_dnat_rules" "test" {
  protocol = var.protocol
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DNAT rules are located.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the DNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the NAT gateway to which the DNAT rule belongs.

* `protocol` - (Optional, String) Specifies the protocol type of the DNAT rule.
  The value can be one of the following:
  + **tcp**
  + **udp**
  + **any**

* `port_id` - (Optional, String) Specifies the port ID of the backend instance to which the DNAT rule belongs.

* `private_ip` - (Optional, String) Specifies the private IP address of the backend instance to which the DNAT rule
  belongs.

* `status` - (Optional, String) Specifies the status of the DNAT rule.
  The value can be one of the following:
  + **ACTIVE**: The SNAT rule is available.
  + **EIP_FREEZED**: The EIP is frozen associated with SNAT rule.
  + **INACTIVE**: The SNAT rule is unavailable.

* `internal_service_port` - (Optional, String) Specifies the port of the backend instance to which the DNAT rule
  belongs.

* `external_service_port` - (Optional, String) Specifies the port of the EIP associated with the DNAT rule.

* `floating_ip_id` - (Optional, String) Specifies the ID of the EIP associated with the DNAT rule.

* `floating_ip_address` - (Optional, String) Specifies the IP address of the EIP associated with the DNAT rule.

* `global_eip_id` - (Optional, String) Specifies the ID of the global EIP associated with the DNAT rule.

* `global_eip_address` - (Optional, String) Specifies the IP address of the global EIP associated with the DNAT rule.

* `description` - (Optional, String) Specifies the description of the DNAT rule.

* `created_at` - (Optional, String) Specifies the creation time of the DNAT rule.
  The format is **yyyy-mm-dd hh:mm:ss.SSSSSS**. e.g. **2024-12-20 15:03:04.000000**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list ot the DNAT rules.
  The [rules](#nat_dnat_rules) structure is documented below.

<a name="nat_dnat_rules"></a>
The `rules` block supports:

* `id` - The ID of the DNAT rule.

* `gateway_id` - The ID of the NAT gateway to which the DNAT rule belongs.

* `protocol` - The protocol type of the DNAT rule.

* `port_id` - The port ID of the backend instance to which the DNAT rule belongs.

* `private_ip` - The private IP address of the backend instance to which the DNAT rule belongs.

* `internal_service_port` - The port of the backend instance to which the DNAT rule belongs.

* `external_service_port` - The port of the EIP associated with the DNAT rule belongs.

* `floating_ip_id` - The ID of the EIP associated with the DNAT rule.

* `floating_ip_address` - The IP address of the EIP associated with the DNAT rule.

* `global_eip_id` - The ID of the global EIP associated with the DNAT rule.

* `global_eip_address` - The IP address of the global EIP associated with the DNAT rule.

* `internal_service_port_range` - The port range of the backend instance to which the DNAT rule belongs.

* `external_service_port_range` - The port range of the EIP associated with the DNAT rule.

* `description` - The description of the DNAT rule.

* `status` - The status of the DNAT rule.

* `created_at` - The creation time of the DNAT rule.
