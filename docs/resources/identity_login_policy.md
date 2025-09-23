---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_login_policy"
description: |-
  Manages the account login policy within HuaweiCloud.
---

# huaweicloud_identity_login_policy

Manages the account login policy within HuaweiCloud.

-> **NOTE:**
  You *must* have admin privileges to use this resource.  
  This resource overwrites an existing configuration, make sure one resource per account.  
  During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_login_policy" "test" {
  account_validity_period    = 20
  lockout_duration           = 30
  login_failed_times         = 10
  period_with_login_failures = 30
  session_timeout            = 120
  show_recent_login_info     = true
  custom_info_for_login      = "hello Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `account_validity_period` - (Optional, Int) Specifies the validity period (days) to disable users
  if they have not logged in within the period. The valid value is range from `0` to `240`.

* `custom_info_for_login` - (Optional, String) Specifies the custom information that will be displayed
  upon successful login.

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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of account login policy, which is the same as the account ID.

## Import

Identity login policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import huaweicloud_identity_login_policy.example <your account ID>
```
