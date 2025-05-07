---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_dnat_rule"
description: |-
  Manages a DNAT rule resource of the **public** NAT within HuaweiCloud.
---

# huaweicloud_nat_dnat_rule

Manages a DNAT rule resource of the **public** NAT within HuaweiCloud.

## Example Usage

### DNAT rule in VPC scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = var.gateway_id
  floating_ip_id        = var.publicip_id
  port_id               = huaweicloud_compute_instance.test.network[0].port
  protocol              = "tcp"
  internal_service_port = 23
  external_service_port = 8023
}
```

```hcl
variable "gateway_id" {}
variable "geip_id" {}

resource "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_global_eip_associate" "test" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.test
  ]

  nat_gateway_id        = var.gateway_id
  global_eip_id         = var.geip_id
  port_id               = huaweicloud_compute_instance.test.network[0].port
  protocol              = "tcp"
  internal_service_port = 23
  external_service_port = 8023
}
```

### DNAT rule in VPC scenario and specify the port ranges

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id              = var.gateway_id
  floating_ip_id              = var.publicip_id
  port_id                     = huaweicloud_compute_instance.test.network[0].port
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
```

```hcl
variable "gateway_id" {}
variable "geip_id" {}

resource "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_global_eip_associate" "test" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.test
  ]

  nat_gateway_id              = var.gateway_id
  global_eip_id               = var.geip_id
  port_id                     = huaweicloud_compute_instance.test.network[0].port
  protocol                    = "tcp"
  internal_service_port_range = "23-823"
  external_service_port_range = "8023-8823"
}
```

### DNAT rule in Direct Connect scenario

```hcl
variable "gateway_id" {}
variable "publicip_id" {}

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = var.gateway_id
  floating_ip_id        = var.publicip_id
  private_ip            = "10.0.0.12"
  protocol              = "any"
  internal_service_port = 0
  external_service_port = 0
}
```

```hcl
variable "gateway_id" {}
variable "geip_id" {}

resource "huaweicloud_global_eip_associate" "test" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "test" {
  depends_on = [
    huaweicloud_global_eip_associate.test
  ]

  nat_gateway_id        = var.gateway_id
  global_eip_id         = var.geip_id
  private_ip            = "10.0.0.12"
  protocol              = "any"
  internal_service_port = 0
  external_service_port = 0
}
```

### DNAT rule in VPC scenario, allow the RDS instance to provide external services

```hcl
variable "subnet_id" {}
variable "natgw_id" {}
variable "publicip_id" {}

resource "huaweicloud_rds_instance" "db_pgSql" {
  ...
}

data "huaweicloud_networking_port" "pgSql_network_port" {
  network_id = var.subnet_id
  fixed_ip   = huaweicloud_rds_instance.db_pgSql.fixed_ip
}

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = var.natgw_id
  floating_ip_id        = var.publicip_id
  port_id               = data.huaweicloud_networking_port.pgSql_network_port.port_id
  protocol              = "tcp"
  internal_service_port = huaweicloud_rds_instance.db_pgSql.db.0.port
  external_service_port = 5432
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DNAT rule is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `nat_gateway_id` - (Required, String, ForceNew) Specifies the ID of the NAT gateway to which the DNAT rule belongs.  
  Changing this will create a new resource.

* `floating_ip_id` - (Optional, String) Specifies the ID of the floating IP address.

* `global_eip_id` - (Optional, String) Specifies the ID of the global EIP connected by the DNAT rule.

-> Fields `floating_ip_id` and `global_eip_id` cannot be set or empty simultaneously.

* `protocol` - (Required, String) Specifies the protocol type.  
  The valid values are **tcp**, **udp**, and **any**.

* `internal_service_port` - (Optional, Int) Specifies port used by Floating IP provide services for external
  systems.  
  Exactly one of `internal_service_port` and `internal_service_port_range` must be set.

* `external_service_port` - (Optional, Int) Specifies port used by ECSs or BMSs to provide services for
  external systems.  
  Exactly one of `external_service_port` and `external_service_port_range` must be set.  
  Required if `internal_service_port` is set.

* `internal_service_port_range` - (Optional, String) Specifies port range used by Floating IP provide services
  for external systems.  
  This parameter and `external_service_port_range` are mapped **1:1** in sequence(, ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.

* `external_service_port_range` - (Optional, String) Specifies port range used by ECSs or BMSs to provide
  services for external systems.  
  This parameter and `internal_service_port_range` are mapped **1:1** in sequence(, ranges must have the same length).
  The valid value for range is **1~65535** and the port ranges can only be concatenated with the `-` character.  
  Required if `internal_service_port_range` is set.

* `port_id` - (Optional, String) Specifies the port ID of network. This parameter is mandatory in VPC scenario.  
  Use [huaweicloud_networking_port](../data-sources/networking_port) to get the port if just know a fixed IP addresses
  on the port.

* `private_ip` - (Optional, String) Specifies the private IP address of a user. This parameter is mandatory in
  Direct Connect scenario.

* `description` - (Optional, String) Specifies the description of the DNAT rule.  
  The value is a string of no more than `255` characters, and angle brackets (<>) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The creation time of the DNAT rule.

* `status` - The current status of the DNAT rule.

* `floating_ip_address` - The actual floating IP address.

* `global_eip_address` - The global EIP address connected by the DNAT rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

DNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_dnat_rule.test <id>
```
