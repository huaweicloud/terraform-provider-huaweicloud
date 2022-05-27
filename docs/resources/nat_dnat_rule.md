---
subcategory: "NAT Gateway (NAT)"
---

# huaweicloud_nat_dnat_rule

Manages a DNAT rule resource within HuaweiCloud.

## Example Usage

### DNAT rule in VPC scenario

```hcl
resource "huaweicloud_compute_instance" "instance_1" {
  ...
}

resource "huaweicloud_nat_dnat_rule" "dnat_1" {
  nat_gateway_id        = var.natgw_id
  floating_ip_id        = var.publicip_id
  port_id               = huaweicloud_compute_instance.instance_1.network[0].port
  protocol              = "tcp"
  internal_service_port = 23
  external_service_port = 8023
}
```

### DNAT rule in Direct Connect scenario

```hcl
resource "huaweicloud_nat_dnat_rule" "dnat_2" {
  nat_gateway_id        = var.natgw_id
  floating_ip_id        = var.publicip_id
  private_ip            = "10.0.0.12"
  protocol              = "tcp"
  internal_service_port = 80
  external_service_port = 8080
}
```

### DNAT rule in VPC scenario, allow the rds instance to provide external services

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

resource "huaweicloud_nat_dnat_rule" "dnat_rule_pgSql" {
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

* `region` - (Optional, String, ForceNew) The region in which to create the dnat rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new dnat rule.

* `nat_gateway_id` - (Required, String, ForceNew) ID of the nat gateway this dnat rule belongs to. Changing this creates
  a new dnat rule.

* `floating_ip_id` - (Required, String, ForceNew) Specifies the ID of the floating IP address. Changing this creates a
  new dnat rule.

* `protocol` - (Required, String, ForceNew) Specifies the protocol type. Currently, TCP, UDP, and ANY are supported.
  Changing this creates a new dnat rule.

* `internal_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or BMSs to provide services for
  external systems. Changing this creates a new dnat rule.

* `external_service_port` - (Required, Int, ForceNew) Specifies port used by ECSs or BMSs to provide services for
  external systems. Changing this creates a new dnat rule.

* `port_id` - (Optional, String, ForceNew) Specifies the port ID of network. This parameter is mandatory in VPC
 scenario. Use [huaweicloud_networking_port](../data-sources/networking_port) to get the port if just know a fixed IP
 addresses on the port. Changing this creates a new dnat rule.

* `private_ip` - (Optional, String, ForceNew) Specifies the private IP address of a user. This parameter is mandatory in
  Direct Connect scenario. Changing this creates a new dnat rule.

* `description` - (Optional, String, ForceNew) Specifies the description of the dnat rule.
  The value is a string of no more than 255 characters, and angle brackets (<>) are not allowed.
  Changing this creates a new dnat rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `created_at` - Dnat rule creation time.

* `status` - Dnat rule status.

* `floating_ip_address` - The actual floating IP address.

## Import

DNAT rules can be imported using the following format:

```
$ terraform import huaweicloud_nat_dnat_rule.dnat_1 f4f783a7-b908-4215-b018-724960e5df4a
```
