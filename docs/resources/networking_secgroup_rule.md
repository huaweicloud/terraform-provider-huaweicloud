---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_secgroup_rule"
description: ""
---

# huaweicloud_networking_secgroup_rule

Manages a Security Group Rule resource within HuaweiCloud.

## Example Usage

### Create an ingress rule that opens TCP port 8080 with port range parameters

```hcl
variable "security_group_id" {}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = var.security_group_id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 8080
  port_range_max    = 8080
  remote_ip_prefix  = "0.0.0.0/0"
}
```

### Create an egress rule that opens TCP port 8080 with port range parameters

```hcl
variable "security_group_id" {}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = var.security_group_id
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 8080
  port_range_max    = 8080
  remote_ip_prefix  = "0.0.0.0/0"
}
```

### Create an ingress rule that enable the remote address group and open some TCP ports

```hcl
variable "group_name" {}
variable "security_group_id" {}

resource "huaweicloud_vpc_address_group" "test" {
  name = var.group_name

  addresses = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id       = var.security_group_id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  ports                   = "80,500,600-800"
  protocol                = "tcp"
  priority                = 5
  remote_address_group_id = huaweicloud_vpc_address_group.test.id
}
```

### Create an egress rule that enable the remote address group and open some TCP ports

```hcl
variable "group_name" {}
variable "security_group_id" {}

resource "huaweicloud_vpc_address_group" "test" {
  name = var.group_name

  addresses = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id       = var.security_group_id
  direction               = "egress"
  action                  = "allow"
  ethertype               = "IPv4"
  ports                   = "80,500,600-800"
  protocol                = "tcp"
  priority                = 5
  remote_address_group_id = huaweicloud_vpc_address_group.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the security group rule resource. If
  omitted, the provider-level region will be used. Changing this creates a new security group rule.

* `security_group_id` - (Required, String, ForceNew) Specifies the security group ID the rule should belong to. Changing
  this creates a new security group rule.

* `direction` - (Required, String, ForceNew) Specifies the direction of the rule, valid values are **ingress** or
  **egress**. Changing this creates a new security group rule.

* `ethertype` - (Required, String, ForceNew) Specifies the layer 3 protocol type, valid values are **IPv4** or **IPv6**.
  Changing this creates a new security group rule.

* `description` - (Optional, String, ForceNew) Specifies the supplementary information about the networking security
  group rule. This parameter can contain a maximum of 255 characters and cannot contain angle brackets (< or >).
  Changing this creates a new security group rule.

* `protocol` - (Optional, String, ForceNew) Specifies the layer 4 protocol type, valid values are **tcp**, **udp**,
  **icmp** and **icmpv6**. If omitted, the protocol means that all protocols are supported.
  This is required if you want to specify a port range. Changing this creates a new security group rule.

* `port_range_min` - (Optional, Int, ForceNew) Specifies the lower part of the allowed port range, valid integer value
  needs to be between `1` and `65,535`. Changing this creates a new security group rule.
  This parameter and `ports` are alternative.

* `port_range_max` - (Optional, Int, ForceNew) Specifies the higher part of the allowed port range, valid integer value
  needs to be between `1` and `65,535`. Changing this creates a new security group rule.
  This parameter and `ports` are alternative.

* `ports` - (Optional, String, ForceNew) Specifies the allowed port value range, which supports single port (80),
  continuous port (1-30) and discontinuous port (22, 3389, 80) The valid port values is range form `1` to `65,535`.
  Changing this creates a new security group rule.

* `remote_ip_prefix` - (Optional, String, ForceNew) Specifies the remote CIDR, the value needs to be a valid CIDR (i.e.
  192.168.0.0/16). If not specified, the empty value means all IP addresses, which is same as the value `0.0.0.0/0`.
  Changing this creates a new security group rule.

* `remote_group_id` - (Optional, String, ForceNew) Specifies the remote group ID. Changing this creates a new security
  group rule.

* `remote_address_group_id` - (Optional, String, ForceNew) Specifies the remote address group ID.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `action` - (Optional, String, ForceNew) Specifies the effective policy. The valid values are **allow** and **deny**.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

* `priority` - (Optional, Int, ForceNew) Specifies the priority number.
  The valid value is range from `1` to `100`. The default value is `1`.
  This parameter is not used with `port_range_min` and `port_range_max`.
  Changing this creates a new security group rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Default is 10 minutes.

## Import

Security Group Rules can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_networking_secgroup_rule.secgroup_rule_1 aeb68ee3-6e9d-4256-955c-9584a6212745
```
