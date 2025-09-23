---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_snat_rules"
description: |-
  Use this data source to get the list of private SNAT rules.
---

# huaweicloud_nat_private_snat_rules

Use this data source to get the list of private SNAT rules.

## Example Usage

```hcl
variable "cidr" {}

data "huaweicloud_nat_private_snat_rules" "test" {
  cidr = var.cidr
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the private SNAT rules are located.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the private SNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the private NAT gateway to which the private SNAT rules
  belong.  

* `cidr` - (Optional, String) Specifies the CIDR block of the private SNAT rule.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the private SNAT rule belongs.

* `transit_ip_id` - (Optional, String) Specifies the ID of the transit IP associated with the private SNAT rule.

* `transit_ip_address` - (Optional, String) Specifies the IP address of the transit IP associated with the private
  SNAT rule.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the private SNAT
  rules belong.

* `description` - (Optional, List) Specifies the description of the private SNAT rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list ot the private SNAT rules.
  The [rules](#snatRules) structure is documented below.

<a name="snatRules"></a>
The `rules` block supports:

* `id` - The ID of the private SNAT rule.

* `gateway_id` - The ID of the private NAT gateway to which the private SNAT rule belongs.

* `cidr` - The CIDR block of the private SNAT rule.

* `subnet_id` - The ID of the subnet to which the private SNAT rule belongs.

* `description` - The description of the private SNAT rule.

* `status` - The status of the private SNAT rule.

* `transit_ip_associations` - The transit IP list associate with the private SNAT rule.
  The [transit_ip_associations](#rules_transit_ip_associations) structure is documented below.

* `created_at` - The creation time of the private SNAT rule.

* `updated_at` - The latest update time of the private SNAT rule.

* `enterprise_project_id` - The ID of the enterprise project to which the private SNAT rule belongs.

<a name="rules_transit_ip_associations"></a>
The `transit_ip_associations` block supports:

* `transit_ip_id` - The ID of the transit IP associated with the private SNAT rule.

* `transit_ip_address` - The IP address of the transit IP associated with the private SNAT rule.
