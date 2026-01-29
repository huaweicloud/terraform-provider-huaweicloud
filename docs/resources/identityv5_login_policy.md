---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_login_policy"
description: |-
  Manages the account login policy resource within HuaweiCloud.
---

# huaweicloud_identityv5_login_policy

Manages the account login policy resource within HuaweiCloud.

-> You *must* have admin privileges to use this resource.  
  This resource overwrites an existing configuration, make sure one resource per account.  
  During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

```hcl
variable "allow_address_netmasks" {
  type = list(object({
    address_netmask = string
    description     = optional(string)
  }))
}
variable "allow_ip_ranges" {
  type = list(object({
    ip_range    = string
    description = optional(string)
  }))
}

resource "huaweicloud_identityv5_login_policy" "test" {
  user_validity_period       = 20
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true

  dynamic "allow_address_netmasks" {
    for_each = var.allow_address_netmasks

    content {
      address_netmask = allow_address_netmasks.value.address_netmask
      description     = allow_address_netmasks.value.description
    }
  }

  dynamic "allow_ip_ranges" {
    for_each = var.allow_ip_ranges

    content {
      ip_range    = allow_ip_ranges.value.ip_range
      description = allow_ip_ranges.value.description
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `user_validity_period` - (Optional, Int) Specifies the validity period to disable users, in days.  
  If users do not log in within the validity period, their account will be disabled, and this value does not apply
  to the root user.  
  The valid value ranges from `0` to `240`.

* `custom_info_for_login` - (Optional, String) Specifies the custom information that will be displayed
  upon successful login.  
  It cannot contain the following special characters: `@#%&<>\$^*`.  
  The maximum length is `64` characters.

* `lockout_duration` - (Optional, Int) Specifies the lockout duration after multiple failed login
  attempts, in minutes.  
  The valid value ranges from `15` to `1,440`, defaults to `15`.

* `login_failed_times` - (Optional, Int) Specifies number of consecutive failed login attempts before
  the account is locked.  
  The valid value ranges from `3` to `10`, defaults to `5`.

* `period_with_login_failures` - (Optional, Int) Specifies The period to reset the account lockout
  counter, in minutes.  
  The valid value ranges from `15` to `60`, defaults to `15`.
  
* `session_timeout` - (Optional, Int) Specifies the session timeout duration after user login, in minutes.  
  The valid value ranges from `15` to `1,440`, defaults to `60`.

* `show_recent_login_info` - (Optional, Bool) Specifies whether to display the last login information.  
  Defaults to **false**.

* `allow_address_netmasks` - (Optional, List) Specifies the IP address list or network segment list that are
  allowed to access.  
  The [allow_address_netmasks](#IdentityV5LoginPolicy_AllowAddressNetmasks) structure is documented below.

* `allow_ip_ranges` - (Optional, List) Specifies the IP address range list that are allowed to access.  
  The [allow_ip_ranges](#IdentityV5LoginPolicy_AllowIpRanges) structure is documented below.

* `allow_ip_ranges_ipv6` - (Optional, List) Specifies the IPv6 address range list that are allowed to access.  
  The [allow_ip_ranges_ipv6](#IdentityV5LoginPolicy_AllowIpRangesIpv6) structure is documented below.

<a name="IdentityV5LoginPolicy_AllowAddressNetmasks"></a>
The `allow_address_netmasks` block supports:

* `address_netmask` - (Required, String) Specifies the IP address or network segment, for example, `192.168.0.1/24`.

* `description` - (Optional, String) Specifies the description information.  
  It cannot contain the following special characters: `@#%&<>\$^*`.

<a name="IdentityV5LoginPolicy_AllowIpRanges"></a>
The `allow_ip_ranges` block supports:

* `ip_range` - (Required, String) Specifies the IP address range, for example, `0.0.0.0-255.255.255.255`.

* `description` - (Optional, String) Specifies the description information.  
  It cannot contain the following special characters: `@#%&<>\$^*`.

<a name="IdentityV5LoginPolicy_AllowIpRangesIpv6"></a>
The `allow_ip_ranges_ipv6` block supports:

* `ip_range` - (Required, String) Specifies the IPv6 address range.

* `description` - (Optional, String) Specifies the description information.  
  It cannot contain the following special characters: `@#%&<>\$^*`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the domain ID.

## Import

The login policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_login_policy.test <id>
```
