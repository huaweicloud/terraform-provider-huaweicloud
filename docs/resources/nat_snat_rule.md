---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud_nat_snat_rule

Manages an SNAT rule resource of the **public** NAT within HuaweiCloud.

## Example Usage

### SNAT rule in VPC scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}
variable "subent_id" {}

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = var.gateway_id
  floating_ip_id = var.publicip_id
  subnet_id      = var.subent_id
}
```

### SNAT rule in DC (Direct Connect) scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = var.gateway_id
  floating_ip_id = var.publicip_id
  source_type    = 1
  cidr           = "192.168.10.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SNAT rule is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `nat_gateway_id` - (Required, String, ForceNew) Specifies the ID of the gateway to which the SNAT rule belongs.  
  Changing this will create a new resource.

* `floating_ip_id` - (Required, String) Specifies the IDs of floating IPs connected by SNAT rule.  
  Multiple floating IPs are separated using commas (,). The number of floating IP IDs cannot exceed `20`.

* `subnet_id` - (Optional, String, ForceNew) Specifies the network IDs of subnet connected by SNAT rule (VPC side).  
  This parameter and `cidr` are alternative. Changing this will create a new resource.

* `cidr` - (Optional, String, ForceNew) Specifies the CIDR block connected by SNAT rule (DC side).  
  This parameter and `subnet_id` are alternative. Changing this will create a new resource.

* `source_type` - (Optional, Int, ForceNew) Specifies the resource scenario.  
  The valid values are **0** (VPC scenario) and **1** (Direct Connect scenario), and the default value is `0`.
  Only `cidr` can be specified over a Direct Connect connection. Changing this will create a new resource.

* `description` - (Optional, String) Specifies the description of the SNAT rule.
  The value is a string of no more than `255` characters, and angle brackets (<>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `floating_ip_address` - The actual floating IP address.

* `status` - The status of the SNAT rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

SNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_snat_rule.test 9e0713cb-0a2f-484e-8c7d-daecbb61dbe4
```
