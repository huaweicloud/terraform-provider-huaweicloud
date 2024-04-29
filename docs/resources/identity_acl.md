---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_acl"
description: ""
---

# huaweicloud_identity_acl

Manages an ACL resource within HuaweiCloud IAM service. The ACL allowing user access only from specified IP address
ranges and IPv4 CIDR blocks. The ACL take effect for IAM users under the Domain account rather than the account itself.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_acl" "acl" {
  type = "console"

  ip_cidrs {
    cidr        = "159.138.39.192/32"
    description = "This is a test ip address"
  }
  ip_ranges {
    range       = "0.0.0.0-255.255.255.0"
    description = "This is a test ip range"
  }
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) Specifies the ACL is created through the Console or API.
  Valid values are **console** and **api**. Changing this parameter will create a new ACL.

* `ip_cidrs` - (Optional, List) Specifies the IPv4 CIDR blocks from which console access or api access is allowed.
  The `ip_cidrs` cannot repeat. The [object](#ip_cidrs_object) structure is documented below.

* `ip_ranges` - (Optional, List) Specifies the IP address ranges from which console access or api access is allowed.
  The `ip_ranges` cannot repeat. The [object](#ip_ranges_object) structure is documented below.

-> **NOTE:** Up to 200 `ip_cidrs` and `ip_ranges` can be created in total for each access method.

<a name="ip_cidrs_object"></a>
The `ip_cidrs` block supports:

* `cidr` - (Required, String) Specifies the IPv4 CIDR block, for example, **192.168.0.0/24**.

* `description` - (Optional, String) Specifies a description about an IPv4 CIDR block. This parameter can contain a
  maximum of 255 characters and the following characters are not allowed:**@#%^&*<>\\**.

<a name="ip_ranges_object"></a>
The `ip_ranges` block supports:

* `range` - (Required, String) Specifies the Ip address range, for example, **0.0.0.0-255.255.255.0**.

* `description` - (Optional, String) Specifies a description about an IP address range. This parameter can contain a
  maximum of 255 characters and the following characters are not allowed:**@#%^&*<>\\**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of identity ACL.
