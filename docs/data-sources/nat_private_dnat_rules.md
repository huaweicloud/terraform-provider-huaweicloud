---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_dnat_rules"
description: |-
  Use this data source to get the list of private DNAT rules.
---

# huaweicloud_nat_private_dnat_rules

Use this data source to get the list of private DNAT rules.

## Example Usage

```hcl
variable "backend_type" {}

data "huaweicloud_nat_private_dnat_rules" "test" {
  backend_type = var.backend_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the private DNAT rules are located.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the private DNAT rule.

* `gateway_id` - (Optional, String) Specifies the ID of the private NAT gateway to which the private DNAT rules
  belong.  

* `backend_type` - (Optional, String) Specifies the type of the backend instance to which the private DNAT rules
  belong.
  The value can be one of the following:
  + **COMPUTE**: ECS instance.
  + **VIP**: VIP.
  + **ELB**: ELB loadbalancer.
  + **ELBv3**: The ver.3 ELB loadbalancer.
  + **CUSTOMIZE**: The custom backend IP address.

* `protocol` - (Optional, String) Specifies the protocol type of the private DNAT rules.
  The value can be one of the following:
  + **tcp**
  + **udp**
  + **any**

* `internal_service_port` - (Optional, String) Specifies the port of the backend instance to which the private DNAT
  rule belongs.

* `backend_interface_id` - (Optional, String) Specifies the network interface ID of the backend instance to which the
  private DNAT rule belongs.

* `transit_ip_id` - (Optional, String) Specifies the ID of the transit IP associated with the private DNAT rules.

* `transit_service_port` - (Optional, String) Specifies the port of the transit IP associated with the private DNAT rule.

* `backend_private_ip` - (Optional, String) Specifies the private IP address of the backend instance to which the
  private DNAT rule belongs.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the private DNAT
  rules belong.

* `description` - (Optional, List) Specifies the description of the private DNAT rule.

* `external_ip_address` - (Optional, List) Specifies the transit IP address used to the private DNAT rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list ot the private DNAT rules.
  The [rules](#private_dnat_rules) structure is documented below.

<a name="private_dnat_rules"></a>
The `rules` block supports:

* `id` - The ID of the private DNAT rule.

* `gateway_id` - The ID of the private NAT gateway to which the private DNAT rule belongs.

* `backend_type` - The type of the backend instance to which the private DNAT rule belongs.

* `protocol` - The protocol type of the private DNAT rule.

* `description` - The description of the private DNAT rule.

* `status` - The status of the private DNAT rule.

* `internal_service_port` - The port of the backend instance to which the private DNAT rule belongs.

* `backend_interface_id` - The network interface ID of the backend instance to which the private DNAT rule belongs.

* `transit_ip_id` - The ID of the transit IP associated with the private DNAT rule.

* `transit_` - The port of the transit IP associated with the private DNAT rule.

* `backend_private_ip` - The private IP address of the backend instance to which the private DNAT rule belongs.

* `created_at` - The creation time of the private DNAT rule.

* `updated_at` - The latest update time of the private DNAT rule.

* `enterprise_project_id` - The ID of the enterprise project to which the private DNAT rule belongs.
