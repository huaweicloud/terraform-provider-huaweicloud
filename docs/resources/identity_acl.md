---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_acl"
description: |-
  Manages an ACL policy resource within HuaweiCloud.
---

# huaweicloud_identity_acl

Manages an ACL policy resource within HuaweiCloud.  
The ACL policies allowing user access only from specified IP address ranges and IPv4 CIDR blocks.

-> You **must** have admin privileges to use this resource.<br>
   The ACL take effect for IAM users under the Domain account rather than the account itself.

~> If you are managing ACL policy with **API** type, please ensure that the current execution machine's EIP is added to
   access addresses to guarantee that provider resource can correctly call the IAM APIs.

## Example Usage

### API Access

```hcl
resource "huaweicloud_identity_acl" "acl" {
  type = "api"

  ip_cidrs {
    cidr        = "159.138.39.192/32"
    description = "This is a test CIDR block"
  }
  ip_ranges {
    range       = "0.0.0.0-255.255.255.0"
    description = "This is a test IP range"
  }
}
```

### Console Access with IPv4 and IPv6 Ranges

The following example restricts console access to specific IPv4 and IPv6 address ranges.

~> When configuring console ACL, please ensure that the IP address of your current machine is included in the access
   addresses to avoid being locked out of the console.

```hcl
resource "huaweicloud_identity_acl" "console_acl" {
  type = "console"

  ip_ranges {
    range       = "10.212.1.156-10.202.1.157"
    description = "Office IPv4 address range"
  }

  ipv6_ranges {
    range       = "2001:db8:0:0:0:0:0:0-2001:db8:0:0:0:0:FFFF:FFFF"
    description = "Office IPv6 address range"
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, NonUpdatable) Specifies the type of the ACL policy.  
  The valid values are as follows:
  + **console**
  + **api**

* `ip_cidrs` - (Optional, List) Specifies the IPv4 CIDR blocks from which console access or API access is allowed.
  The [ip_cidrs](#iam_acl_ip_cidrs) structure is documented below.

* `ip_ranges` - (Optional, List) Specifies the IP address ranges from which console access or API access is allowed.
  The [ip_ranges](#iam_acl_ip_ranges) structure is documented below.

-> Up to `200` **ip_cidrs** and **ip_ranges** can be created in total for each access method.

* `ipv6_cidrs` - (Optional, List) Specifies the IPv6 CIDR blocks from which console access or API access is allowed.
  The [ipv6_cidrs](#iam_acl_ipv6_cidrs) structure is documented below.

* `ipv6_ranges` - (Optional, List) Specifies the IPv6 address ranges from which console access is allowed.
  The [ipv6_ranges](#iam_acl_ipv6_ranges) structure is documented below.

<a name="iam_acl_ip_cidrs"></a>
The `ip_cidrs` block supports:

* `cidr` - (Required, String) Specifies the IPv4 CIDR block which allow access through console or API.

  -> CIDR blocks are not allowed to conflict with or duplicate each other.

* `description` - (Optional, String) Specifies a description about an IPv4 CIDR block. This parameter can contain a
  maximum of `255` characters and the following characters are not allowed:**@#%^&*<>\\**.

<a name="iam_acl_ip_ranges"></a>
The `ip_ranges` block supports:

* `range` - (Required, String) Specifies the IPv4 address range which allow access through console or API.

  -> IP address ranges are not allowed to conflict with or duplicate each other.

* `description` - (Optional, String) Specifies a description about an IP address range. This parameter can contain a
  maximum of `255` characters and the following characters are not allowed:**@#%^&*<>\\**.

<a name="iam_acl_ipv6_cidrs"></a>
The `ipv6_cidrs` block supports:

* `cidr` - (Required, String) Specifies the IPv6 CIDR block which allow access through console.

* `description` - (Optional, String) Specifies a description about an IPv6 CIDR block.

<a name="iam_acl_ipv6_ranges"></a>
The `ipv6_ranges` block supports:

* `range` - (Required, String) Specifies the IPv6 address range which allow access through console.

* `description` - (Optional, String) Specifies a description about an IP address range.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of identity ACL.
