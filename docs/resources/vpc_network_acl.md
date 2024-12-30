---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_acl"
description: ""
---

# huaweicloud_vpc_network_acl

Manages a VPC network ACL resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "subnet_id_1" {}
variable "subnet_id_2" {}

resource "huaweicloud_vpc_network_acl" "test" {
  name                  = var.name
  description           = "created by terraform"
  enterprise_project_id = 0
  enabled               = true

  ingress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "192.168.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  ingress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "192.168.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }

  egress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "172.16.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  egress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "172.16.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }
  
  associated_subnets {
    subnet_id = var.subnet_id_1
  }

  associated_subnets {
    subnet_id = var.subnet_id_2
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the network ACL name. The value can contain no more than 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the network ACL.

* `description` - (Optional, String) Specifies the network ACL description. The value can contain no more
  than 255 characters and cannot contain angle brackets (< or >).

* `enabled` - (Optional, Bool) Specifies whether the network ACL is enabled. The default value is **true**.

* `ingress_rules` - (Optional, List) Specifies the ingress rules of the network ACL.
  The [rules](#rules) structure is documented below.

* `egress_rules` - (Optional, List) Specifies the egress rules of the network ACL.
  The [rules](#rules) structure is documented below.

* `associated_subnets` - (Optional, List) Specifies the associated subnets of the network ACL.
  The [associated_subnets](#subnets) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the network ACL.

<a name="rules"></a>
The `ingress_rules` and `egress_rules` block supports:

* `action` - (Required, String) Specifies the rule action. The value can be: **allow** and **deny**.

* `protocol` - (Required, String) Specifies the rule protocol The value can be **tcp**, **udp**, **icmp**, **icmpv6**,
  or an IP protocol number (0–255). The value **any** indicates all protocols.

* `ip_version` - (Required, Int) Specifies the IP version of a network ACL rule.
  The value can be **4** (IPv4) and **6** (IPv6).

* `name` - (Optional, String) Specifies the network ACL rule name. The value can contain no more than 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `description` - (Optional, String) Specifies the network ACL rule description. The value can contain no more
  than 255 characters. The value cannot contain angle brackets (< or >).

* `source_ip_address` - (Optional, String) Specifies the source IP address or CIDR block of a network ACL rule.
 The `source_ip_address` and `source_address_group_id` cannot be configured at the same time.

* `source_ip_address_group_id` - (Optional, String) Specifies the source IP address group ID of a network ACL rule.
  `source_ip_address` and `source_address_group_id` cannot be configured at the same time.

* `source_port` - (Optional, String) Specifies the source ports of a network ACL rule.
  You can specify a single port or a port range. Separate every two entries with a comma.
  
* `destination_ip_address` - (Optional, String) Specifies the destination IP address or CIDR block of a network ACL rule.
  The `destination_ip_address` and `destination_address_group_id` cannot be configured at the same time.
  
* `destination_ip_address_group_id` - (Optional, String) Specifies the destination IP address group ID of a network ACL rule.
  The `destination_ip_address` and `destination_address_group_id` cannot be configured at the same time.

* `destination_port` - (Optional, String) Specifies the destination ports of a network ACL rule.
  You can specify a single port or a port range. Separate every two entries with a comma.

<a name="subnets"></a>
The `associated_subnets` block supports:

* `subnet_id` - (Required, String) Specifies the ID of the subnet to associate with the network ACL.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in uuid format.

* `status` - The status of the ACL.

* `created_at` - The created time of the ACL.

* `updated_at` - The updated time of the ACL.

* `ingress_rules` - The ingress rules of the network ACL.
  The [rules](#rules_resp) structure is documented below.

* `egress_rules` - The egress rules of the network ACL.
  The [rules](#rules_resp) structure is documented below.

<a name="rules_resp"></a>
The `ingress_rules` and `egress_rules` block supports:

* `rule_id` - The ID of the rule.

## Import

The network ACL can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_network_acl.test <id>
```
