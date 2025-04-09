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

* `transit_ip_ids` - (Required, List) Specifies the IDs of the transit IPs associated with the private SNAT rule.
  The maximum of `20` transit IPs can be bound, and the transit IPs must belong the same transit subnet.

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

* `transit_ip_associations` - The transit IP list associate with the private SNAT rule.
  The [transit_ip_associations](#snat_transit_ip_associations) structure is documented below.

* `enterprise_project_id` - The ID of the enterprise project to which the private SNAT rule belongs.

<a name="snat_transit_ip_associations"></a>
The `transit_ip_associations` block supports:

* `transit_ip_id` - The ID of the transit IP associated with the private SNAT rule.

* `transit_ip_address` - The IP address of the transit IP associated with the private SNAT rule.

## Import

The private SNAT rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_snat_rule.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `transit_ip_ids`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_nat_private_snat_rule" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      transit_ip_ids,
    ]
  }
}
```
