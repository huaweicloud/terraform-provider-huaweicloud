---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_snat_rule"
description: ""
---

# huaweicloud_nat_private_snat_rule

Manages an SNAT rule resource of the **private** NAT within HuaweiCloud.

## Example Usage

### Create an SNAT rule via subnet ID

```hcl
variable "gateway_id" {}
variable "transit_ip_id" {}
variable "subnet_id" {}

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id    = var.gateway_id
  transit_ip_id = var.transit_ip_id
  subnet_id     = var.subnet_id
}
```

### Create an SNAT rule via CIDR

```hcl
variable "gateway_id" {}
variable "transit_ip_id" {}
variable "cidr_block" {}

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id    = var.gateway_id
  transit_ip_id = var.transit_ip_id
  cidr          = var.cidr_block
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SNAT rule is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `gateway_id` - (Required, String, ForceNew) Specifies the private NAT gateway ID to which the SNAT rule belongs.  
  Changing this will create a new resource.

* `transit_ip_id` - (Required, String) Specifies the ID of the transit IP associated with SNAT rule.

* `cidr` - (Optional, String, ForceNew) Specifies the CIDR block of the match rule.  
  Changing this will create a new resource.  
  Exactly one of `cidr` and `subnet_id` must be set.

-> SNAT rules under the same private NAT gateway cannot have the same CIDR, but they can be proper subsets of other
   CIDRs.

* `subnet_id` - (Optional, String, ForceNew) Specifies the subnet ID of the match rule.  
  Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the SNAT rule, which contain maximum of `255`
  characters, and angle brackets (< and >) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The creation time of the SNAT rule.

* `updated_at` - The latest update time of the SNAT rule.

* `transit_ip_address` - The address of the transit IP.

* `enterprise_project_id` - The ID of the enterprise project to which the private SNAT rule belongs.

## Import

SNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_snat_rule.test df9b61e9-79c1-4a75-bfab-736e224ced71
```
