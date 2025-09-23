---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_snat_rules"
description: |-
  Use this data source to get the list of SNAT rules.
---

# huaweicloud_nat_snat_rules

Use this data source to get the list of SNAT rules.

## Example Usage

```hcl
variable "rule_id" {}

data "huaweicloud_nat_snat_rules" "test" {
  rule_id = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SNAT rules are located.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the SNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the NAT gateway to which the SNAT rule belongs.

* `floating_ip_id` - (Optional, String) Specifies the ID of the EIP associated with SNAT rule.

* `floating_ip_address` - (Optional, String) Specifies the IP of the EIP associated with SNAT rule.

* `cidr` - (Optional, String) Specifies the CIDR block to which the SNAT rule belongs.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the SNAT rule belongs.

* `source_type` - (Optional, String) Specifies the source type of the SNAT rule.
  The value can be one of the following:
  + **0** : The use scenario is VPC.
  + **1** : The use scenario is DC.

* `status` - (Optional, String) Specifies the status of the SNAT rule.
  The value can be one of the following:
  + **ACTIVE**: The SNAT rule is available.
  + **EIP_FREEZED**: The global EIP is frozen associated with SNAT rule.
  + **INACTIVE**: The SNAT rule is unavailable.

* `global_eip_id` - (Optional, String) Specifies the ID of the global EIP associated with SNAT rule.

* `global_eip_address` - (Optional, String) Specifies the IP of the global EIP associated with SNAT rule.

* `description` - (Optional, String) Specifies the description of the SNAT rule.

* `created_at` - (Optional, String) Specifies the creation time of the SNAT rule.
  The format is **yyyy-mm-dd hh:mm:ss.SSSSSS**. e.g. **2024-12-20 15:03:04.000000**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the SNAT rules.
  The [rules](#nat_snat_rules) structure is documented below.

<a name="nat_snat_rules"></a>
The `rules` block supports:

* `id` - The ID of the SNAT rule.

* `gateway_id` - The ID of the NAT gateway to which the SNAT rule belongs.

* `cidr` - The CIDR block to which the SNAT rule belongs.

* `subnet_id` - The ID of the subnet to which the SNAT rule belongs.

* `source_type` - The source type of the SNAT rule.

* `floating_ip_id` - The IDs of the EIP associated with SNAT rule, multiple EIP IDs separate by commas.
  e.g. **ID1,ID2**.

* `floating_ip_address` - The IPs of the EIP associated with SNAT rule, multiple EIP IPs separate by commas.
  e.g. **IP1,IP2**.

* `description` - The description of the SNAT rule.

* `status` - The status of the SNAT rule.

* `global_eip_id` - The IDs of the global EIP associated with SNAT rule, multiple global EIP IDs separate by commas.
  e.g. **ID1,ID2**.

* `global_eip_address` - The IPs of the global EIP associated with SNAT rule, multiple global EIP IPs separate by commas.
  e.g. **IP1,IP2**.

* `freezed_ip_address` - The IP of the frozen global EIP associated with SNAT rule.

* `created_at` - The creation time of the SNAT rule.
