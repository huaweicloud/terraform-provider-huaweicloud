---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_acls"
description: |-
  Use this data source to get the list of VPC network ACLs.
---

# huaweicloud_vpc_network_acls

Use this data source to get the list of VPC network ACLs.

## Example Usage

```hcl
variable "network_acl_name" {}
variable "enterprise_project_id" {}

data "huaweicloud_vpc_network_acls" "basic" {
  name                  = var.network_acl_name
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to obtain the network ACLs.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the network ACL name. The value can contain no more than 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `network_acl_id` - (Optional, String) Specifies the network ACL ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the network ACL.

* `enabled` - (Optional, String) Specifies whether the network ACL is enabled. The value can be **true** or **false**.

* `status` - (Optional, String) Specifies the status of the network ACL.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in uuid format.

* `network_acls` - The list of VPC network ACLs.
  The [network_acls](#network_acls) structure is documented below.

<a name="network_acls"></a>
The `network_acls` block supports:

* `name` - The network ACL name.

* `id` - The network ACL ID.

* `enterprise_project_id` - The enterprise project ID of the network ACL.

* `description` - The network ACL description.

* `enabled` - Whether the network ACL is enabled.

* `ingress_rules` - The ingress rules of the network ACL.
  The [rules](#rules) structure is documented below.

* `egress_rules` - The egress rules of the network ACL.
  The [rules](#rules) structure is documented below.

* `associated_subnets` - The associated subnets of the network ACL.
  The [associated_subnets](#subnets) structure is documented below.

* `status` - The status of the ACL.

* `created_at` - The created time of the ACL.

* `updated_at` - The updated time of the ACL.

<a name="rules"></a>
The `ingress_rules` and `egress_rules` block supports:

* `rule_id` - The ID of the rule.

* `action` - The rule action.

* `protocol` - The rule protocol.

* `ip_version` - The IP version of a network ACL rule.

* `name` - The network ACL rule name.

* `description` - The network ACL rule description.

* `source_ip_address` - The source IP address or CIDR block of a network ACL rule.

* `source_ip_address_group_id` - The source IP address group ID of a network ACL rule.

* `source_port` - The source ports of a network ACL rule.
  
* `destination_ip_address` - The destination IP address or CIDR block of a network ACL rule.
  
* `destination_ip_address_group_id` - The destination IP address group ID of a network ACL rule.

* `destination_port` - The destination ports of a network ACL rule.

<a name="subnets"></a>
The `associated_subnets` block supports:

* `subnet_id` - The ID of the subnet to associate with the network ACL.
