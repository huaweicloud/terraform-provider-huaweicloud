---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_login_policy"
description: |-
  Manages the account login policy v5 resource within HuaweiCloud.
---

# huaweicloud_identityv5_login_policy

Manages the account login policy v5 resource within HuaweiCloud.

-> **NOTE:**
  You *must* have admin privileges to use this resource.  
  This resource overwrites an existing configuration, make sure one resource per account.  
  During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

```hcl
resource "huaweicloud_identityv5_login_policy" "test" {
  user_validity_period       = 20
  custom_info_for_login      = "hello Terraform"
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true
  allow_address_netmasks {
    address_netmask = "255.0.0.0/1"
    description     = "terraform test"
  }
  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.254"
    description = "terraform test1"
  }
  allow_ip_ranges {
    ip_range    = "0.0.0.0-255.255.255.100"
    description = "terraform test2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `user_validity_period` - (Optional, Int) Specifies the validity period (days) to disable users
  if they have not logged in within the period. The valid value is range from `0` to `240`.

* `custom_info_for_login` - (Optional, String) Specifies the custom information that will be displayed
  upon successful login. It cannot contain the following special
  characters: `@`, `#`, `%`, `&`, `<`, `>`, `\`, `$`, `^` and `*`.
  Its maximum length: 64.

* `lockout_duration` - (Optional, Int) Specifies the duration (minutes) to lock users out.
  The valid value is range from `15` to `1440`, defaults to `15`.

* `login_failed_times` - (Optional, Int) Specifies the number of unsuccessful login attempts to lock users out.
  The valid value is range from `3` to `10`, defaults to `5`.

* `period_with_login_failures` - (Optional, Int) Specifies the period (minutes) to count the number of unsuccessful
  login attempts. The valid value is range from `15` to `60`, defaults to `15`.
  
* `session_timeout` - (Optional, Int) Specifies the session timeout (minutes) that will apply if you or users created
  using your account do not perform any operations within a specific period.
  The valid value is range from `15` to `1,440`, defaults to `60`.

* `show_recent_login_info` - (Optional, Bool) Specifies whether to display last login information upon successful login.
  The value can be **true** or **false**.

* `allow_address_netmasks` - (Optional, List) Specifies the IP addresses or network segments that are allowed to access.
  The [allow_address_netmasks](#IdentityV5LoginPolicy_AllowAddressNetmasks) structure is documented below.

* `allow_ip_ranges` - (Optional, List) Specifies the IP address range that can be accessed.
  The [allow_ip_ranges](#IdentityV5LoginPolicy_AllowIpRanges) structure is documented below.

* `allow_ip_ranges_ipv6` - (Optional, List) Specifies the IPv6 address range that can be accessed.
  The [allow_ip_ranges_ipv6](#IdentityV5LoginPolicy_AllowIpRangesIpv6) structure is documented below.

<a name="IdentityV5LoginPolicy_AllowAddressNetmasks"></a>
The `allow_address_netmasks` block contains:

* `address_netmask` - (Required, String) Specifies the IP address or network segment, for example, `192.168.0.1/24`.

* `description` - (Optional, String) Specifies the description information, which cannot contain the following special
  characters: `@`, `#`, `%`, `&`, `<`, `>`, `\`, `$`, `^` and `*`.

<a name="IdentityV5LoginPolicy_AllowIpRanges"></a>
The `allow_ip_ranges` block contains:

* `ip_range` - (Required, String) Specifies the IP address range, for example, `0.0.0.0-255.255.255.255`.

* `description` - (Optional, String) Specifies the description information, which cannot contain the following special
  characters: `@`, `#`, `%`, `&`, `<`, `>`, `\`, `$`, `^` and `*`.

<a name="IdentityV5LoginPolicy_AllowIpRangesIpv6"></a>
The `allow_ip_ranges_ipv6` block contains:

* `ip_range` - (Required, String) Specifies the IPv6 address range.

* `description` - (Optional, String) Specifies the description information, which cannot contain the following special
  characters: `@`, `#`, `%`, `&`, `<`, `>`, `\`, `$`, `^` and `*`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of account login policy, which is the same as the account ID.

## Import

Identity login policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import huaweicloud_identityv5_login_policy.example <your account ID>
```
