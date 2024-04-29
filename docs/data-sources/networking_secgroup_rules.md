---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_networking_secgroup_rules"
description: ""
---

# huaweicloud_networking_secgroup_rules

Use this data source to get the list of the available HuaweiCloud security group rules.

## Example Usage

```hcl
variable "security_group_id" {}

data "huaweicloud_networking_secgroup_rules" "test" {
  security_group_id = var.security_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the security group rules. If omitted, the
  provider-level region will be used.

* `security_group_id` - (Required, String) Specifies the security group ID that the rule should belong to.

* `rule_id` - (Optional, String) Specifies the security group rule ID used for query.

* `protocol` - (Optional, String) Specifies the security group rule protocol type used for query.  
  The value can be **tcp**, **udp**, **icmp**, **icmpv6** or IP protocol number, if empty, it indicates support for
  all protocols.

* `description` - (Optional, String) Specifies the security group rule description used for query.

* `remote_group_id` - (Optional, String) Specifies the remote security group ID used for query.

* `direction` - (Optional, String) Specifies the direction of the security group rule used for query.  
  The valid values are as follows:
  + **ingress**
  + **egress**

* `action` - (Optional, String) Specifies the effective policy of the security group rule used for query.  
  The valid values are as follows:
  + **allow**
  + **deny**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `rules` - All security group rules that match the filter parameters.
  The [rules](#secgroup_rules) structure is documented below.

<a name="secgroup_rules"></a>
The `rules` block supports:

* `id` - The ID of the security group rule.

* `description` - The description of the security group rule.

* `security_group_id` - The security group ID that the rule should belong to.

* `direction` - The direction of the security group rule.

* `protocol` - The security group rule protocol type.

* `ethertype` - The security group rule IP address protocol type. The value can be **IPv4** or **IPv6**.

* `ports` - The range of port values for security group rule. Which supports single port (80), continuous port (1-30)
  and discontinuous port (22, 3389, 80).

* `action` - The effective policy of the security group rule.

* `priority` - The priority of security group rule. The valid value ranges from `1` to `100`, `1` represents the
  highest priority.

* `remote_group_id` - The remote security group ID.  
  This field is mutually exclusive with `remote_ip_prefix` and `remote_address_group_id`.

* `remote_ip_prefix` - The remote IP address. The value can be in the CIDR format or IP addresses.  
  This field is mutually exclusive with `remote_group_id` and `remote_address_group_id`.

* `remote_address_group_id` - The remote address group ID.  
  This field is mutually exclusive with `remote_group_id` and `remote_ip_prefix`.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The latest update time, in UTC format.
